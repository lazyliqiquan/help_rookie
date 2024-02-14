package models

type Documents struct {
	ID          int `gorm:"primaryKey"`
	ImagePaths  GormStrList
	Text        string
	SeekHelpsID int
	LendHandsID int
}
