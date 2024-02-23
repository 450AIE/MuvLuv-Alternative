package model

// 用户信息
type User struct {
	Id           int    `json:"id"`
	Gender       string `json:"gender"`
	Nickname     string `json:"nickname"`
	QQ           int    `json:"QQ"`
	Birthday     string `json:"birthday"`
	Email        string `json:"email"`
	Avatar       string `json:"avatar"`
	Introduction string `json:"introduction"`
	Phone        int    `json:"phone"`
}

// sb前端送来修改用户信息的结构体
type UserToChangeInfo struct {
	Gender       string `json:"gender"`
	Nickname     string `json:"nickname"`
	QQ           string `json:"QQ"`
	Birthday     string `json:"birthday"`
	Email        string `json:"email"`
	Avatar       string `json:"avatar"`
	Introduction string `json:"introduction"`
	TelePhone    string `json:"telePhone"`
}

type UserToChangeInfoForMySQL struct {
	OldNickname string `json:"oldNickname"`
	UserToChangeInfo
}

// 书籍信息
type Book struct {
	Book_id      int    `json:"book_id"`
	Name         string `json:"name"`
	Is_star      bool   `json:"is_star"`
	Author       string `json:"author"`
	Comment_num  int    `json:"comment_num"`
	Score        int    `json:"score"`
	Cover        string `json:"cover"`
	Publish_time string `json:"publish_Time"`
	Link         string `json:"link"`
	Lable        string `json:"lable"`
}

type BookComment struct {
	Post_id      int    `json:"post_id"`
	Book_id      string `json:"book_id"`
	Publish_time int    `json:"publish_time"`
	Content      string `json:"content"`
	User_id      int    `json:"user_id"`
	Avatar       string `json:"avatar"`
	Nickname     string `json:"nickname"`
	Praise_count int    `json:"praise_count"`
	Is_praise    bool   `json:"is_praise"`
	Is_focus     bool   `json:"is_focus"`
}

type GetBookCommentForMySQL struct {
	UserID string `json:"user_id"`
	BookID string `json:"book_id"`
}

// 注册业务，从body获取到的结构体
type UserRegisterIn struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type UserRegisterToMySQL struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type UserLoginIn struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type UserLoginToMySQL struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type UserChangePassWordIn struct {
	Username     string `json:"username"`
	Old_password string `json:"old_password"`
	New_password string `json:"new_password"`
}

type UserChangePassWordInToMySQL struct {
	Username     string `json:"username"`
	Old_password string `json:"old_password"`
	New_password string `json:"new_password"`
}

type UserInfoGetForMySQL struct {
	UserID string `json:"userID"`
}

type BookSearchToMySQL struct {
	Name string `json:"book_name"`
}

type BookSearchIn struct {
	Book_name string `json:"book_name,omitempty"`
}

type BookChangeStar struct {
	UserID string `json:"userID"`
	BookID string `json:"bookID"`
}

//type BookID struct {
//	Book_id string `json:"book_id"`
//}

type CheckUserStarBookToMySQL struct {
	UserID string `json:"userID"`
}

type GetLabelBookForMySQL struct {
	Label string `json:"label"`
}

type WriteBookComment struct {
	Book_id string `json:"book_id"`
	Content string `json:"content"`
}

type WriteBookCommentToMySQL struct {
	UserID  string `json:"user_id"`
	BookID  string `json:"book_id"`
	Content string `json:"content"`
}

type ChangeBookComment struct {
	Comment_id string `json:"comment_id"`
	Content    string `json:"content"`
}
type ChangeBookCommentForMySQL struct {
	PostID  string `json:"postID"`
	UserID  string `json:"userID"`
	Content string `json:"content"`
}

type ThumbsUpForMySQL struct {
	UserID   string `json:"userID"` //点赞者
	Model    string `json:"model"`
	TargetID string `json:"targetID"` //被点赞的评论
}

type UserCollectedBookList struct {
	BookID       string `json:"bookID"`
	Name         string `json:"name"`
	Publish_time string `json:"publish_Time"`
	Link         string `json:"link"`
}

type GetUserCollectedBookListFromMySQL struct {
	UserID string `json:"userID"`
}

type UserFocus struct {
	CaringUserID string `json:"caringUserID"`
	CaredUserID  string `json:"caredUserID"`
}

type UserFocusForMySQL struct {
	UserFocus
}

type DeleteBookCommentToMySQL struct {
	PostID string `json:"postID"`
}
