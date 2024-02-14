package models

type LendHands struct {
	ID          int `gorm:"primaryKey"`
	CreateTime  string
	CodePath    string
	DiffPath    string
	Status      string
	Like        GormIntList
	Ban         int
	UsersID     int
	SeekHelpsID int
	Document    Documents  `gorm:"foreignKey:LendHandsID"`
	Comments    []Comments `gorm:"foreignKey:LendHandsID"`
}
