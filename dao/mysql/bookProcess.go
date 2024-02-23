package mysql

import (
	"Web/model"
	"Web/utility"
	"errors"
	"log"
	"strconv"
)

//要注意的是，没有登录看所有的is_star都是false，登录后不同的人有不同的is_star，艹，怎么办

func GetBookList() ([]model.Book, error) {
	//查询所有的书籍，不只是一行的所有数据，该怎么办呢？
	operation := "select * from books "
	rows, err := DB.Query(operation)
	if err != nil {
		log.Println(err)
		return nil, errors.New(utility.StatusInfoMap[utility.StatusBookQueryErr])
	}
	var book = model.Book{}
	var booklist []model.Book
	for rows.Next() {
		err := rows.Scan(&book.Book_id, &book.Name, &book.Is_star, &book.Author, &book.Comment_num, &book.Score, &book.Cover, &book.Publish_time, &book.Link, &book.Lable)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		booklist = append(booklist, book)
	}
	return booklist, nil
}

func ChangeBookStar(in *model.BookChangeStar) error {
	//先检查书是否存在，存在再改变
	exist, err := CheckBookExist(in.BookID)
	if err != nil {
		return err
	}
	if !exist {
		return errors.New(utility.StatusInfoMap[utility.StatusBookNotExist])
	}
	//开始改变,先检查是否收藏
	exist, err = CheckIfAlreadyStar(in.UserID, in.BookID)
	if err != nil {
		return errors.New(utility.StatusInfoMap[utility.StatusBookAlreadyStar])
	}
	if exist {
		return errors.New(utility.StatusInfoMap[utility.StatusBookAlreadyStar])
	}
	err = bookStarChange(in.UserID, in.BookID)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func bookStarChange(userID string, bookID string) error {
	operation := "insert into is_starbook (user_id,book_id) value(?,?)"
	_, err := DB.Exec(operation, userID, bookID)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func GetLabeledBook(in *model.GetLabelBookForMySQL) ([]model.Book, error) {
	operation := "select * from books where label like ? "
	rows, err := DB.Query(operation, "%"+in.Label+"%") //模糊匹配
	if err != nil {
		log.Println(err)
		return nil, err
	}
	var book model.Book
	var booklist []model.Book
	for rows.Next() {
		err = rows.Scan(&book.Book_id, &book.Name, &book.Is_star, &book.Author, &book.Comment_num, &book.Score, &book.Cover, &book.Publish_time, &book.Link, &book.Lable)
		if err != nil {
			return nil, err
		}
		booklist = append(booklist, book)
	}
	return booklist, nil
}

// 没有实现嵌套评论功能，所以model当前只会是1
func ThumbsUp(in *model.ThumbsUpForMySQL) error {
	if in.Model == "1" {
		//检查一下当前用户是否点赞过该评论
		already, err := CheckTheUserWhetherAlreadyThumbsUpTheComment(in.UserID, in.TargetID)
		if err != nil {
			return err
		}
		if !already { //没有点赞过，该为点赞
			err = AddThumbsUpToComment(in.UserID, in.TargetID)
			if err != nil {
				return err
			}
			return nil
		} else { //点赞过了，取消点赞
			err = CancelThumbsUpToComment(in.UserID, in.TargetID)
			if err != nil {
				return err
			}
			return nil
		}
	} else { //model==2时，暂未实现
		return errors.New("暂未实现嵌套评论")
	}
}

func GetUserCollectedBookList(in *model.GetUserCollectedBookListFromMySQL) ([]model.UserCollectedBookList, error) {
	//得到的book_id就是用户收藏的书籍的id
	operation := "select book_id from is_starbook where user_id = ?"
	rows, err := DB.Query(operation, in.UserID)
	//查询不到数据是null,不是报错
	if err != nil {
		log.Println(err)
		return nil, err
	}
	var book model.UserCollectedBookList
	var booklist []model.UserCollectedBookList
	for rows.Next() {
		//把每一个书籍ID对应信息找出来，针对性填充
		var bookID string
		err := rows.Scan(&bookID)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		bookinfo := GetTheBookInfoWithID(bookID)
		book.BookID = strconv.Itoa(bookinfo.Book_id)
		book.Link = bookinfo.Link
		book.Name = bookinfo.Name
		book.Publish_time = bookinfo.Publish_time
		booklist = append(booklist, book)
	}
	return booklist, nil
}
