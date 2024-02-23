package service

import (
	"Web/dao/mysql"
	"Web/model"
	"Web/utility"
	"errors"
	"log"
)

func RegisterServe(in *model.UserRegisterIn) error {
	var u = model.UserRegisterToMySQL{in.Username, in.Password}
	return mysql.RegisterUser(&u)
}

func LoginServe(in *model.UserLoginIn) error {
	var u = model.UserLoginToMySQL{in.Username, in.Password}
	return mysql.LoginUser(&u)
}

func PasswordChangeServe(in *model.UserChangePassWordIn) error {
	var u = model.UserChangePassWordInToMySQL{in.Username, in.Old_password, in.New_password}
	return mysql.PasswordChange(&u)
}

func InfoChangeServe(oldUsername string, in *model.UserToChangeInfo) error {
	var u = model.UserToChangeInfoForMySQL{oldUsername, *in}
	return mysql.InfoChange(&u)
}

func GetUserInfoServe(userID string) (model.User, error) {
	var u = model.UserInfoGetForMySQL{UserID: userID}
	return mysql.GetUserInfo(&u)
}

func GetBookListServe() ([]model.Book, error) {
	return mysql.GetBookList()
}

func SearchBookWithName(name string) (bool, error, model.Book) {
	var u = model.BookSearchToMySQL{name}
	return mysql.SearchAndGetBookWithName(&u)
}

func ChangeBookStar(userID string, bookID string) error {
	var u = model.BookChangeStar{userID, bookID}
	return mysql.ChangeBookStar(&u)
}

// 得到的是未经雪花算法加密过的自增ID
func GetUserIDWithName(name string) (string, error) {
	return mysql.GetUserIDWithName(name)
}

func GetTheUserStarBookID(username string) ([]int, error) {
	ID, err := mysql.GetUserIDWithName(username)
	if err != nil {
		return nil, err
	}
	var u = model.CheckUserStarBookToMySQL{ID}
	return mysql.GetUserStarBookID(&u)
}

func ChangeBookStarToTrueInBookList(booklist *[]model.Book, starbookID *[]int) *[]model.Book {
	for _, v := range *starbookID {
		for k, v1 := range *booklist {
			if v == v1.Book_id {
				(*booklist)[k].Is_star = true
			}
		}
	}
	return booklist
}

func GetTheLabeledBook(label string) (*[]model.Book, error) {
	var u = model.GetLabelBookForMySQL{label}
	booklist, err := mysql.GetLabeledBook(&u)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return &booklist, nil
}

// 如果username为空，说明没有该用户
func GetBookComment(username string, bookID string) (*[]model.BookComment, error) {
	if username == "" {
		var u = model.GetBookCommentForMySQL{"", bookID}
		bookcomment, err := mysql.GetBookComment(&u)
		return &bookcomment, err
	}
	userID, err := GetUserIDWithName(username)
	if err != nil {
		return nil, err
	}
	var u = model.GetBookCommentForMySQL{userID, bookID}
	bookcomment, err := mysql.GetBookComment(&u)
	return &bookcomment, err
}

// 返回该书评的ID号
func WriteBookComment(username string, bookID string, content string) (int, error) {
	userID, err := GetUserIDWithName(username)
	if err != nil {
		log.Println(err)
		return -1, err
	}
	var u = model.WriteBookCommentToMySQL{userID, bookID, content}
	return mysql.WriteBookComment(&u)
}

// 只需要返回是否成功即可
func ChangeBookComment(postID string, username string, content string) error {
	//检查修改用户和原评论用户是否是同一个用户,书籍ID和书评对应的书籍ID是否对应。直接找到post_id对应的用户ID和书籍ID对照即可
	userID, err := GetUserIDWithName(username)
	if err != nil {
		log.Println(err)
		return err
	}
	var u = model.ChangeBookCommentForMySQL{postID, userID, content}
	return mysql.ChangeBookComment(&u)
}

func ThumbsUp(username string, model1 string, targetID string) error {
	userID, err := GetUserIDWithName(username)
	if err != nil {
		return err
	}
	var u = model.ThumbsUpForMySQL{userID, model1, targetID}
	return mysql.ThumbsUp(&u)
}

func GetUserCollectedBookList(username string) ([]model.UserCollectedBookList, error) {
	userID, err := GetUserIDWithName(username)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	var u = model.GetUserCollectedBookListFromMySQL{userID}
	return mysql.GetUserCollectedBookList(&u)
}

// 两次关注就是取消关注
func UserFocus(caringUserName string, caredUserID string) error {
	caringUserID, err := GetUserIDWithName(caringUserName)
	if err != nil {
		return err
	}
	var u = model.UserFocusForMySQL{model.UserFocus{caringUserID, caredUserID}}
	return mysql.FocusUser(&u)
}

// 要检验发布者和删除者是否是用一个人,以及这个评论是否存在
func DeleteBookComment(username string, postID string) error {
	exist, err := CheckCommentWhetherExist(postID)
	if err != nil {
		return err
	}
	if !exist {
		return errors.New(utility.StatusInfoMap[utility.StatusBookCommentNotExist])
	}
	userID, err := GetUserIDWithName(username)
	if err != nil {
		return err
	}
	is, err := mysql.CheckWhetherUserIsCommenter(userID, postID)
	if err != nil {
		return err
	}
	if !is {
		return errors.New(utility.StatusInfoMap[utility.StatusUserIsNotTheCommenter])
	}
	var u = model.DeleteBookCommentToMySQL{postID}
	return mysql.DeleteBookComment(&u)
}

func CheckCommentWhetherExist(postID string) (bool, error) {
	return mysql.CheckTheBookCommentExist(postID)
}
