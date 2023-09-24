package database

import (
	"context"
	"gin-rest-api-example/internal/cache"
	"gin-rest-api-example/internal/game/model"

	"github.com/jinzhu/gorm"
)

type GameDB interface {
	RunInTx(ctx context.Context, f func(ctx context.Context) error) error
	FindPopularGameModesByAreaCode(ctx context.Context, areaCode string) ([]*model.GameMode, error)
}

type gameDB struct {
	db *gorm.DB
}

func NewGameDB(db *gorm.DB, cacher cache.Cacher) GameDB {
	return &gameDB{db: db}
}

func (gdb *gameDB) RunInTx(ctx context.Context, f func(ctx context.Context) error) error {
	// Implement the RunInTx method here
	// You can use transactions with the provided DB instance (gdb.db) and call the function f
	return f(ctx)
}

func (gdb *gameDB) FindPopularGameModesByAreaCode(ctx context.Context, areaCode string) ([]*model.GameMode, error) {
	sqlQuery := `
		SELECT gm.id, gm.mode
		FROM GameMode gm
		JOIN GameModePopularity gmp ON gm.id = gmp.gamemodeid
		WHERE gmp.areacode = ? 
		ORDER BY gmp.popularity DESC;
	`

	// Execute the raw SQL query and map the results to the GameMode struct.
	rows, err := gdb.db.Raw(sqlQuery, areaCode).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var popularGameModes []*model.GameMode

	for rows.Next() {
		var gameMode model.GameMode
		if err := gdb.db.ScanRows(rows, &gameMode); err != nil {
			return nil, err
		}
		popularGameModes = append(popularGameModes, &gameMode)
	}

	return popularGameModes, nil
}
