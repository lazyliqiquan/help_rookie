package models

const (
	CView int = iota //浏览
)

type Comments struct {
	ID          int `gorm:"primaryKey"`
	Text        string
	SendTime    string
	Like        GormIntList
	Ban         int
	SeekHelpsID int
	LendHandsID int
	UsersId     int
	Users       Users
}
