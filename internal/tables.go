package forum

// SingInUsers ..
type SingInUsers struct {
	ID       int
	Name     string
	Login    string
	Password string
}

// Post ..
type Post struct {
	ID              int
	Title           string
	Text            string
	Author          string
	Tags            []string
	Comments        []Comment
	CountOfLikes    int
	CountOfDisLikes int
	ErrorVal        Error
}

type ChosenPost struct {
	Post         Post
	IsPostOfUser bool
}

type EditComments struct {
	Comments []Comment
	ErrorVal Error
}

// Comment ..
type Comment struct {
	ID              int
	Text            string
	Author          string
	CountOfLikes    int
	CountOfDisLikes int
	IsCommentOfUser bool
}

//Error ..
type Error struct {
	MSG string
	Err bool
}

//Homepage ..
type Homepage struct {
	Posts                    []Post
	Categories               []string
	Category                 string
	InSession                bool
	LikedPostNotification    []PostNotification
	DislikedPostNotification []PostNotification
	CommentNotification      []CommentNotification
	IsNotification           bool
	CountOfNotifications     int
}

//UserProfile ..
type UserProfile struct {
	User           User
	CreatedPosts   []Post
	LikedPosts     []Post
	DislikedPosts  []Post
	CommentedPosts []CommentedPosts
	Notifications  []Notification
	IsNotification bool
}

//UnionEnterPost ..
type UnionEnterPost struct {
	Post       Post
	ErrorVal   Error
	Categories []string
}

//UpdatePost ..
type UpdatePost struct {
	PostId     int
	Post       CreatePost
	ErrorVal   Error
	Categories []string
}

//UpdateComment ..
type UpdateComment struct {
	PostId   string
	Comments []Comment
	Comment  Comment
	ErrorVal Error
}

//UnionEnterComment ..
type UnionEnterComment struct {
	ErrorVal Error
	PostId   int
	Comment  string
}

//CreatePost ..
type CreatePost struct {
	Title      string
	Text       string
	Categories []string
}

//Categories ..
type Categories struct {
	Tags      []string
	InSession bool
}

//PostNotification ..
type PostNotification struct {
	Post  Post
	Liker SingInUsers
}

//CommentNotification ..
type CommentNotification struct {
	Post        Post
	Commentator SingInUsers
}

//CommentedPosts ..
type CommentedPosts struct {
	Post    Post
	Comment Comment
}

//allnotification
type Notification struct {
	ID      int
	Title   *string // '*' chtobi poluchit' nil esli pole mozhet byt' null i s nim porabotat'
	Comment *string
	// Like      *bool
	// DisLike   *bool
	LikerName       *string
	CommentatorName *string
	CommentedPost   *string
}
