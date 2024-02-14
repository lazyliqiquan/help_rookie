package models

const (
	SuperAdmin uint64 = iota //超级管理员(第一个注册的用户是超级管理员,有且只有一个)
	Admin                    //普通管理员
)

type Users struct {
	ID           int `gorm:"primaryKey"`
	Name         string
	Email        string
	Password     string
	Score        int
	RegisterTime string
	Ban          uint64
	SeekHelps    []SeekHelps `gorm:"foreignKey:UsersID"`
	LendHands    []LendHands `gorm:"foreignKey:UsersID"`
}
