package models

const (
	SOthersView  int = iota //非所有者浏览
	SUserView               //所有者浏览
	SEdit                   //修改
	SViewComment            //浏览评论
	SAddComment             //添加评论
	SAddLendHand            //添加帮助
)

type SeekHelps struct {
	ID         int `gorm:"primaryKey"`
	Score      int
	CreateTime string
	UpdateTime string
	CodePath   string
	Language   string
	Like       GormIntList
	PageView   int
	Status     bool
	Document   string
	ImagePath  GormStrList
	Ban        int
	UsersID    int
	LendHands  []LendHands `gorm:"foreignKey:SeekHelpsID"`
	Comments   []Comments  `gorm:"foreignKey:SeekHelpsID"`
}
