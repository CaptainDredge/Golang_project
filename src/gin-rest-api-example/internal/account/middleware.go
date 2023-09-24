package account

import (
	accountDB "gin-rest-api-example/internal/account/database"
	"gin-rest-api-example/internal/account/model"
	"gin-rest-api-example/internal/cache"
	"gin-rest-api-example/internal/config"
	"gin-rest-api-example/pkg/logging"
	"net/http"
	"time"

	"github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

var identityKey = "id"

type signIn struct {
	User struct {
		Email    string `form:"email" json:"email" binding:"email"`
		Password string `form:"password" json:"password" binding:"required"`
	} `json:"user"`
}

func CurrentUser(c *gin.Context) (*model.Account, bool) {
	data, ok := c.Get(identityKey)
	if !ok {
		return nil, false
	}
	acc, ok := data.(*model.Account)
	return acc, ok
}

func MustCurrentUser(c *gin.Context) *model.Account {
	acc, ok := CurrentUser(c)
	if ok {
		return acc
	}
	panic("no account in gin.Context")
}

func NewAuthMiddleware(cfg *config.Config, accountDB accountDB.AccountDB) (*jwt.GinJWTMiddleware, error) {
	return jwt.New(&jwt.GinJWTMiddleware{
		Realm:       "test zone",
		Key:         []byte(cfg.JwtConfig.Secret),
		Timeout:     cfg.JwtConfig.SessionTime,
		MaxRefresh:  time.Hour,
		IdentityKey: identityKey,
		PayloadFunc: func(data interface{}) jwt.MapClaims {
			if v, ok := data.(*model.Account); ok {
				return jwt.MapClaims{
					identityKey: v.Email,
				}
			}
			return jwt.MapClaims{}
		},
		IdentityHandler: func(c *gin.Context) interface{} {
			claims := jwt.ExtractClaims(c)
			email, ok := claims[identityKey].(string)
			if !ok {
				return nil
			}
			acc, err := accountDB.FindByEmail(c.Request.Context(), email)
			if err != nil {
				return nil
			}
			return acc
		},
		Authenticator: func(c *gin.Context) (interface{}, error) {
			var req signIn
			if err := c.ShouldBindJSON(&req); err != nil {
				return "", jwt.ErrMissingLoginValues
			}

			ctx := cache.WithCacheSkip(c.Request.Context(), true)
			acc, err := accountDB.FindByEmail(ctx, req.User.Email)
			if err != nil || acc.Disabled {
				return nil, jwt.ErrFailedAuthentication
			}
			err = MatchesPassword(acc.Password, req.User.Password)
			if err != nil {
				if err != bcrypt.ErrMismatchedHashAndPassword {
					logging.FromContext(c).Warnw("middleware.jwt.Authenticator found unknown error when matches password", "err", err)
				}
				return nil, jwt.ErrFailedAuthentication
			}
			return &model.Account{
				ID:        acc.ID,
				Username:  acc.Username,
				Email:     acc.Email,
				Password:  "",
				Bio:       acc.Bio,
				Image:     acc.Image,
				CreatedAt: acc.CreatedAt,
				UpdatedAt: acc.UpdatedAt,
				Disabled:  acc.Disabled,
			}, nil
		},
		Authorizator: func(data interface{}, c *gin.Context) bool {
			logging.FromContext(c).Info("middleware.jwt.Authorizator", "data", data)
			if _, ok := data.(*model.Account); ok {
				return true
			}
			return false
		},
		Unauthorized: func(c *gin.Context, code int, message string) {
			logging.FromContext(c).Info("middleware.jwt.Unauthorized", "code", code, "message", message)
			c.JSON(code, gin.H{
				"code":    code,
				"message": message,
			})
		},
		LoginResponse: func(c *gin.Context, code int, token string, expire time.Time) {
			c.JSON(http.StatusOK, gin.H{
				"code":   code,
				"token":  token,
				"expire": expire,
			})
		},
		TokenLookup:   "header: Authorization",
		TokenHeadName: "Bearer",
		TimeFunc:      time.Now,
	})
}
