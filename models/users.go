package models

// 从低位开始使用
const (
	Root             int = iota //超级管理员(第一个注册的用户是超级管理员,有且只有一个)
	Login                       //登录
	PublishSeekHelp             //发表求助
	ViewSeekHelp                //浏览求助
	EditSeekHelp                //编辑求助
	PublishLendHand             //发表帮助
	ViewLendHand                //浏览帮助
	EditLendHand                //编辑帮助
	PublishComment              //发表评论
	ViewComment                 //浏览评论
	PublishShareCode            //发表共享代码
	ViewShareCode               //浏览共享代码
)

type Users struct {
	ID           int `gorm:"primaryKey"`
	Name         string
	Email        string
	Password     string
	Avatar       string
	Score        int
	RegisterTime string
	Ban          int
	Message      GormIntList
	SeekHelps    []SeekHelps `gorm:"foreignKey:UsersID"`
	LendHands    []LendHands `gorm:"foreignKey:UsersID"`
}

func UserPermit(option, ban int) bool {
	return (ban & (1 << option)) == 0
}
