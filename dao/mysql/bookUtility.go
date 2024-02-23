package mysql

import (
	"Web/model"
	"log"
)

func SearchAndGetBookWithName(in *model.BookSearchToMySQL) (bool, error, model.Book) {
	operation := "select * from books where name=?"
	rows, err := DB.Query(operation, in.Name)
	if err != nil {
		log.Println(err)
		return false, err, model.Book{}
	}
	var book = model.Book{}
	for rows.Next() {
		err := rows.Scan(&book.Book_id, &book.Name, &book.Is_star, &book.Author, &book.Comment_num, &book.Score, &book.Cover, &book.Publish_time, &book.Link, &book.Lable)
		if err != nil {
			log.Println(err)
			return false, err, model.Book{}
		}
	}
	//如果没有填充，说明没有找到
	if book.Author == "" {
		return false, nil, model.Book{}
	}
	return true, nil, book
}

// 得到的是自增ID，是未经过雪花算法加密过的ID
func GetUserIDWithName(username string) (string, error) {
	operation := "select id from users where nickname=?"
	row := DB.QueryRow(operation, username)
	var userID string
	err := row.Scan(&userID)
	if err != nil {
		log.Println(err)
		return "", err
	}
	return userID, nil
}

// true表示已存在，false表示不存在,err出现表示有查询错误
func CheckIfAlreadyStar(userID string, bookID string) (bool, error) {
	operation := "select count(id) from is_starbook where user_id=? and book_id=?"
	row := DB.QueryRow(operation, userID, bookID)
	var num int
	err := row.Scan(&num)
	if err != nil {
		return false, err
	}
	if num == 0 {
		return false, nil
	} else {
		return true, nil
	}
}

func GetUserStarBookID(in *model.CheckUserStarBookToMySQL) ([]int, error) {
	operation := "select book_id from is_starbook where user_id=?"
	rows, err := DB.Query(operation, in.UserID)
	if err != nil {
		return nil, err
	}
	var ID int
	var starbooklist []int
	for rows.Next() {
		err := rows.Scan(&ID)
		if err != nil {
			return nil, err
		}
		starbooklist = append(starbooklist, ID)
	}
	return starbooklist, nil
}

// 书籍不存在返回false,nil
func CheckBookExist(bookID string) (bool, error) {
	operation := "select count(?) from books"
	var num int
	row := DB.QueryRow(operation, bookID)
	err := row.Scan(&num)
	if err != nil {
		return false, err
	}
	if num == 0 {
		return false, nil
	}
	return true, nil
}

func FillTheRestUserInfo(in *model.BookComment) error {
	operation := "select avatar,nickname from users where id=?"
	row := DB.QueryRow(operation, in.User_id)
	err := row.Scan(&in.Avatar, &in.Nickname)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func GetTheBookInfoWithID(bookID string) model.Book {
	operation := "select * from books where book_id=?"
	row := DB.QueryRow(operation, bookID)
	var bookinfo model.Book
	row.Scan(&bookinfo.Book_id, &bookinfo.Name, &bookinfo.Is_star, &bookinfo.Author, &bookinfo.Comment_num, &bookinfo.Score, &bookinfo.Cover, &bookinfo.Publish_time, &bookinfo.Link, &bookinfo.Lable)
	return bookinfo
}
