package models

// 从低位开始使用(就不设置浏览权限了，因为浏览权限都没有，相当于封禁该用户)
const (
	Admin            int = iota //管理员
	Login                       //登录
	PublishSeekHelp             //发表求助
	EditSeekHelp                //编辑求助
	PublishLendHand             //发表帮助
	EditLendHand                //编辑帮助
	PublishComment              //发表评论
	PublishShareCode            //发表共享代码
)

type Users struct {
	ID              int `gorm:"primaryKey"`
	Name            string
	Email           string
	Password        string
	Avatar          string
	Score           int
	RegisterTime    string
	CommentSurplus  int
	LastPublishDate string
	Ban             int
	Message         GormIntList
	SeekHelpCollect GormIntList
	LendHandCollect GormIntList
	SeekHelps       []SeekHelps `gorm:"foreignKey:UsersID"`
	LendHands       []LendHands `gorm:"foreignKey:UsersID"`
}

func JudgePermit(option, ban int) bool {
	return (ban & (1 << option)) == 0
}
