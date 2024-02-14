package models

type Comments struct {
	ID          int `gorm:"primaryKey"`
	Text        string
	SendTime    string
	UsersId     int
	Users       Users
	SeekHelpsID int
	LendHandsID int
}
