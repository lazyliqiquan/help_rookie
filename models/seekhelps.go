package models

type SeekHelps struct {
	ID         int `gorm:"primaryKey"`
	Score      int
	CreateTime string
	CodePath   string
	Language   string
	Like       GormIntList
	Ban        int
	Status     bool
	UsersID    int
	Document   Documents   `gorm:"foreignKey:SeekHelpsID"`
	LendHands  []LendHands `gorm:"foreignKey:SeekHelpsID"`
	Comments   []Comments  `gorm:"foreignKey:SeekHelpsID"`
}
