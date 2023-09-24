package model

type GameModePopularity struct {
	AreaCode   string `gorm:"column:areacode"`
	GameModeID uint   `gorm:"column:gamemodeid"`
	popularity uint   `gorm:"column:popularity"`
}

type GameMode struct {
	ID   uint   `gorm:"column:id"`
	Mode string `gorm:"column:mode"`
}
