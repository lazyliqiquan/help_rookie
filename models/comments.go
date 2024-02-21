package models

type Comments struct {
	ID          int `gorm:"primaryKey"`
	Text        string
	SendTime    string
	Status      bool //false 可以浏览 true 不可浏览
	UsersId     int
	Users       Users
	SeekHelpsID int
	LendHandsID int
}
