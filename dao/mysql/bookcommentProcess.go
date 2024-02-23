package mysql

import (
	"Web/model"
	"Web/utility"
	"errors"
	"fmt"
	"log"
	"strconv"
	"time"
)

func WriteBookComment(in *model.WriteBookCommentToMySQL) (int, error) {
	operation := "insert into book_comment (user_id,book_id,content,publish_time) value(?,?,?,?)"
	_, err := DB.Exec(operation, in.UserID, in.BookID, in.Content, time.Now().Unix())
	if err != nil {
		log.Println(err)
		return -1, err
	}
	//获得这次书评的序列号
	operation = "select post_id from book_comment where user_id=? and book_id=? and content=?"
	row := DB.QueryRow(operation, in.UserID, in.BookID, in.Content)
	var ID int
	err = row.Scan(&ID)
	if err != nil {
		log.Println(err)
		return -1, err
	}
	//把对应书籍ID的评论数+1
	operation = "update books set comment_num= comment_num+1 where book_id=?"
	_, err = DB.Exec(operation, in.BookID)
	if err != nil {
		return -1, err
	}
	return ID, nil
}

func ChangeBookComment(in *model.ChangeBookCommentForMySQL) error {
	//通过post_id检查一下用户是否一致
	ok, err := CheckCommenterWheterIsChanger(in.PostID, in.UserID)
	if err != nil {
		return nil
	}
	if !ok {
		return errors.New(utility.StatusInfoMap[utility.StatusBookChangerAndCommenterNotEqual])
	} else { //用户一致，开始修改
		err = ChangeBookCommentNow(in)
		if err != nil {
			return err
		}
		return nil
	}
}

// 我真的日了，首先，要检查查看的这个用户是否关注评论用户，是否点赞评论
// 其次，要检查评论用户，获取他的信息，填充到BookComment里面返回给查看用户
// 如果username为空，说明没有该用户
func GetBookComment(in *model.GetBookCommentForMySQL) ([]model.BookComment, error) {
	//选择查看哪本书的评论
	//operation := "select (book_id,publish_time,content,user_id,praise_count,is_praised,is_focus) from book_comment where book_id=?"
	operation := "select post_id,publish_time,content,user_id,praise_count,is_praised,is_focus from book_comment where book_id=?"
	rows, err := DB.Query(operation, in.BookID)
	if err != nil {
		fmt.Println("这里出错")
		panic(err)
		log.Println(err)
		return nil, err
	}
	var comment model.BookComment
	var Commentlist []model.BookComment
	for rows.Next() {
		comment.Book_id = in.BookID
		err = rows.Scan(&comment.Post_id, &comment.Publish_time, &comment.Content, &comment.User_id, &comment.Praise_count, &comment.Is_praise, &comment.Is_focus)
		if err != nil {
			return nil, err
		}
		//检查一下用户是否点赞这个评论
		is, err := CheckTheUserWhetherAlreadyThumbsUpTheComment(in.UserID, strconv.Itoa(comment.Post_id))
		if err != nil {
			log.Println(err)
			return nil, err
		}
		if is {
			comment.Is_praise = true
		}
		//再填充头像等评论用户的信息
		err = FillTheRestUserInfo(&comment)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		//检查查看用户是否关注了评论用户,userID为空就不用管了，因为MySQL里默认就是false
		if in.UserID != "" {
			ifcare, err := CheckUserCare(in.UserID, strconv.Itoa(comment.User_id))
			if err != nil {
				log.Println(err)
				return nil, err
			}
			comment.Is_focus = ifcare
		}
		Commentlist = append(Commentlist, comment)
	}
	return Commentlist, nil
}

func DeleteBookComment(in *model.DeleteBookCommentToMySQL) error {
	//要删除点赞信息，删除评论,注意book表的评论数-1
	err := DeleteTheCommentAllStar(in.PostID)
	bookID, err := GetBookIDWithPostID(in.PostID)
	if err != nil {
		return err
	}
	if err != nil {
		return err
	}
	operation := "delete from book_comment where post_id=?"
	_, err = DB.Exec(operation, in.PostID)
	if err != nil {
		log.Println(err)
		return err
	}
	operation = "update books set comment_num=comment_num-1 where book_id=?"
	_, err = DB.Exec(operation, bookID)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}
