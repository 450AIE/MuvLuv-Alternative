package mysql

import (
	"Web/model"
	"log"
	"strconv"
)

func AddThumbsUpToComment(userID string, postID string) error {
	operation := "update book_comment set praise_count=praise_count+1 where post_id=?"
	_, err := DB.Exec(operation, postID)
	if err != nil {
		return err
	}
	operation = "insert into is_thumbsup_comment (user_id,post_id) value(?,?)"
	_, err = DB.Exec(operation, userID, postID)
	if err != nil {
		return err
	}
	return nil
}

func CancelThumbsUpToComment(userID string, postID string) error {
	operation := "update book_comment set praise_count=praise_count-1 where post_id=?"
	_, err := DB.Exec(operation, postID)
	if err != nil {
		return err
	}
	operation = "delete from is_thumbsup_comment where user_id=? and post_id=?"
	_, err = DB.Exec(operation, userID, postID)
	if err != nil {
		return err
	}
	return nil
}

// 如果该用户已经点赞了该评论，返回true，否则返回false，其他错误返回err
func CheckTheUserWhetherAlreadyThumbsUpTheComment(userID string, postID string) (bool, error) {
	operation := "select count(ID) from is_thumbsup_comment where post_id=? and user_id =?"
	row := DB.QueryRow(operation, postID, userID)
	var num int
	err := row.Scan(&num)
	if err != nil {
		log.Println(err)
		return false, err
	}
	if num == 1 { //已经点赞过了
		return true, nil
	} else {
		return false, nil
	}
}

func ChangeBookCommentNow(in *model.ChangeBookCommentForMySQL) error {
	operation := "update book_comment set content = ? where post_id=?"
	_, err := DB.Exec(operation, in.Content, in.PostID)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

// 是true，不是false，其他情况err
func CheckCommenterWheterIsChanger(postID string, userID string) (bool, error) {
	operation := "select user_id from book_comment where post_id=?"
	row := DB.QueryRow(operation, postID)
	var commenterID int
	err := row.Scan(&commenterID)
	if err != nil {
		log.Println(err)
		return false, err
	}
	userid, err := strconv.ParseInt(userID, 10, 64)
	if err != nil {
		return false, nil
	}
	if int(userid) == commenterID {
		return true, nil
	} else {
		return false, nil
	}
}

func DeleteTheCommentAllStar(postID string) error {
	operation := "delete from is_thumbsup_comment where post_id=?"
	_, err := DB.Exec(operation, postID)
	if err != nil {
		return err
	}
	return nil
}

func GetBookIDWithPostID(postID string) (string, error) {
	operation := "select book_id from book_comment where post_id=?"
	row := DB.QueryRow(operation, postID)
	var bookID string
	err := row.Scan(&bookID)
	if err != nil {
		return "", err
	}
	return bookID, nil
}

func CheckTheBookCommentExist(postID string) (bool, error) {
	operation := "select count(post_id) from book_comment where post_id=?"
	row := DB.QueryRow(operation, postID)
	var num int
	err := row.Scan(&num)
	if err != nil {
		return false, err
	}
	if num == 1 {
		return true, nil
	} else {
		return false, nil
	}
}
