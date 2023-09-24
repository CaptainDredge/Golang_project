package game

import (
	"gin-rest-api-example/internal/config"
	"gin-rest-api-example/internal/database"
	gamedb "gin-rest-api-example/internal/game/database"
	"gin-rest-api-example/internal/middleware"
	"gin-rest-api-example/internal/middleware/handler"
	"gin-rest-api-example/pkg/logging"
	"gin-rest-api-example/pkg/validate"
	"net/http"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func NewHandler(gameDB gamedb.GameDB) *Handler {
	return &Handler{
		gameDB: gameDB,
	}
}

type Handler struct {
	gameDB gamedb.GameDB
}

// articleBySlug handles GET /v1/api/articles/:slug
func (h *Handler) popularityByAreacode(c *gin.Context) {
	handler.HandleRequest(c, func(c *gin.Context) *handler.Response {
		logger := logging.FromContext(c)
		// bind
		type RequestUri struct {
			Areacode string `uri:"areacode"`
		}
		var uri RequestUri
		if err := c.ShouldBindUri(&uri); err != nil {
			logger.Errorw("article.handler.articleBySlug failed to bind", "err", err)
			var details []*validate.ValidationErrDetail
			if vErrs, ok := err.(validator.ValidationErrors); ok {
				details = validate.ValidationErrorDetails(&uri, "uri", vErrs)
			}
			return handler.NewErrorResponse(http.StatusBadRequest, handler.InvalidUriValue, "invalid article request in uri", details)
		}

		// find
		gamemodes, err := h.gameDB.FindPopularGameModesByAreaCode(c.Request.Context(), uri.Areacode)
		if err != nil {
			if database.IsRecordNotFoundErr(err) {
				return handler.NewErrorResponse(http.StatusNotFound, handler.NotFoundEntity, "not found area code", nil)
			}
			return handler.NewInternalErrorResponse(err)
		}
		return handler.NewSuccessResponse(http.StatusOK, NewGameModesResponse(gamemodes, len(gamemodes)))
	})
}

func RouteV1(cfg *config.Config, h *Handler, r *gin.Engine, auth *jwt.GinJWTMiddleware) {
	v1 := r.Group("v1/api")
	v1.Use(middleware.RequestIDMiddleware(), middleware.TimeoutMiddleware(cfg.ServerConfig.WriteTimeout))

	articleV1 := v1.Group("popularity")
	// anonymous
	articleV1.Use()
	{
		articleV1.GET(":areacode", h.popularityByAreacode)
	}
}
