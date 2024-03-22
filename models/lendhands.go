package models

// 在新增帖子的时候，按钮上的文字是发布(需要上传代码文件)
// 在修改帖子的时候，按钮上的文字是保存(上传的代码文件将会覆盖原本的)
// 在帮助帖子被求助者认可后，不能再修改
const (
	LView        int = iota //浏览
	LEdit                   //修改
	LViewComment            //浏览评论
	LAddComment             //添加评论
)

type LendHands struct {
	ID          int `gorm:"primaryKey"`
	CreateTime  string
	UpdateTime  string
	CodePath    string
	DiffPath    string
	Status      bool
	Like        GormIntList
	Document    string
	ImagePath   GormStrList
	Tags        GormStrList
	Ban         int
	LikeSum     int
	CommentSum  int
	UsersID     int
	Users       Users
	SeekHelpsID int
	SeekHelps   SeekHelps
	Comments    []Comments //`gorm:"foreignKey:LendHandsID"`
}
