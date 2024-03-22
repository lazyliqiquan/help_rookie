package models

// 如果该求助已经有用户帮助，那么该求助帖子不能再修改
const (
	SView        int = iota //浏览
	SEdit                   //修改
	SViewComment            //浏览评论
	SAddComment             //添加评论
	SAddLendHand            //添加帮助
)

type SeekHelps struct {
	ID          int `gorm:"primaryKey"`
	Reward      int
	CreateTime  string
	UpdateTime  string
	CodePath    string
	Language    string
	Like        GormIntList
	PageView    int
	Status      bool
	Document    string
	ImagePath   GormStrList
	Tags        string
	Ban         int
	LikeSum     int
	LendHandSum int
	CommentSum  int
	UsersID     int
	Users       Users
	LendHands   []LendHands //`gorm:"foreignKey:SeekHelpsID"`
	Comments    []Comments  //`gorm:"foreignKey:SeekHelpsID"`
}
