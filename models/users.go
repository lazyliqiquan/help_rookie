package models

// 从低位开始使用
const (
	SuperAdmin int = iota //超级管理员(第一个注册的用户是超级管理员,有且只有一个)
	Admin                 //普通管理员
	Login                 //登录
	SeekHelp              //求助
	LendHand              //帮助
	Comment               //评论
)

type Users struct {
	ID           int `gorm:"primaryKey"`
	Name         string
	Email        string
	Password     string
	Score        int
	RegisterTime string
	Ban          int
	SeekHelps    []SeekHelps `gorm:"foreignKey:UsersID"`
	LendHands    []LendHands `gorm:"foreignKey:UsersID"`
}

func IsLogin(ban int) bool {
	return (ban & (1 << Login)) == 0
}

func IsAdmin(ban int) bool {
	return ((ban & (1 << Admin)) == 0) || IsRoot(ban)
}

func IsRoot(ban int) bool {
	return (ban & (1 << SuperAdmin)) == 0
}
