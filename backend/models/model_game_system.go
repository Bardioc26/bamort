package models

type GameSystem struct {
	ID          uint   `gorm:"primaryKey"`
	Code        string `gorm:"uniqueIndex;size:100;not null"`
	Name        string `gorm:"size:255;not null"`
	Description string `gorm:"type:text"`
	IsActive    bool   `gorm:"default:true;not null"`
	CreatedAt   int64  `gorm:"autoCreateTime"`
	ModifiedAt  int64  `gorm:"autoUpdateTime"`
}

// TableName sets the table name for SchemaVersion
func (GameSystem) TableName() string {
	return "game_systems"
}
