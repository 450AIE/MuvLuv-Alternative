package controller

import (
	"Web/model"
	"Web/service"
	"Web/utility"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
)

// 没有token就是没登录状态，is_star就默认false，有token就验证，过期了返回响应
// 我们似乎可以这么处理书籍根据用户是否收藏显示true或者false，就是书籍默认false，
// 弄第三张表，如果是登陆用户，就根据这张表找到用户的收藏情况来针对性显示
func GetBookListHandler(c *gin.Context) {
	//查看有无token
	token, _ := GetTheTokenFromHeader(c)
	if token == "" {
		booklist, err := service.GetBookListServe()
		if err != nil {
			if err.Error() == utility.StatusInfoMap[utility.StatusBookQueryErr] {
				utility.ResponseErr(c, utility.StatusBookQueryErr)
				return
			}
			utility.ResponseErr(c, utility.StatusServerBusy)
			return
		}
		utility.ResponseSuccess(c, booklist)
		return
	}
	exist, err, username := VerifyAccessToken(token)
	if err != nil {
		utility.ResponseErr(c, utility.StatusServerBusy)
		return
	}
	if exist { //登录用户，先找到该用户的收藏情况,我们不是后来会返回一个切片吗，我们就把切片里面的is_star改了就可以
		StarBookList, err := service.GetTheUserStarBookID(username)
		booklist, err := service.GetBookListServe()
		if err != nil {
			if err.Error() == utility.StatusInfoMap[utility.StatusBookQueryErr] {
				utility.ResponseErr(c, utility.StatusBookQueryErr)
				return
			}
			utility.ResponseErr(c, utility.StatusServerBusy)
			return
		}
		//这里需要把收藏了的改为true
		custombooklist := service.ChangeBookStarToTrueInBookList(&booklist, &StarBookList)
		//
		utility.ResponseSuccess(c, custombooklist)
		return
	} else {
		utility.ResponseErr(c, utility.StatusAccessTokenOutOfRange)
		return
	}
}

func SearchBookHandler(c *gin.Context) {
	var bookname = model.BookSearchIn{}
	bookname.Book_name = c.Query("book_name")
	exist, err, book := service.SearchBookWithName(bookname.Book_name)
	if err != nil {
		log.Println(err)
		utility.ResponseErr(c, utility.StatusServerBusy)
		return
	}
	fmt.Println(exist)
	if !exist {
		utility.ResponseErr(c, utility.StatusBookNotExist)
		return
	}
	utility.ResponseSuccess(c, book)
	return
}

func BookStarHandler(c *gin.Context) {
	//先验证有没有access token，有的话再检查有没有这本书，有的话改is_star为true
	token, err := GetTheTokenFromHeader(c)
	if err != nil {
		log.Println(err)
		utility.ResponseErr(c, utility.StatusAccessTokenErr)
		return
	}
	ok, err, username := VerifyAccessToken(token)
	if err != nil {
		log.Println(err)
		utility.ResponseErr(c, utility.StatusAccessTokenErr)
		return
	}
	if !ok {
		log.Println(err)
		utility.ResponseErr(c, utility.StatusAccessTokenOutOfRange)
		return
	}
	//获取userid
	userID, err := service.GetUserIDWithName(username)
	if err != nil {
		log.Println(err)
		utility.ResponseErr(c, utility.StatusServerBusy)
		return
	}
	//获取form的bookid
	//在获取Foem表单前，必须加这一行语句，该语句会先读取body解析为url和表单数据
	err = c.Request.ParseForm()
	if err != nil {
		log.Println(err)
		utility.ResponseErr(c, utility.StatusServerBusy)
		return
	}
	bookid := c.Request.Form.Get("book_id")
	log.Println(bookid)
	if err != nil {
		log.Println(err)
		utility.ResponseErr(c, utility.StatusServerBusy)
		return
	}
	err = service.ChangeBookStar(userID, bookid)
	if err != nil {
		if err.Error() == utility.StatusInfoMap[utility.StatusBookAlreadyStar] {
			utility.ResponseErr(c, utility.StatusBookAlreadyStar)
			return
		}
		utility.ResponseErr(c, utility.StatusServerBusy)
		return
	}
	utility.ResponseSuccess(c, nil)
	return
}

// 我们获取label标签相关书籍的ID，然后调用书籍列表Handler得到所有书籍再根据ID排除，
// 或者直接数据库找，还是后者吧，前者太傻逼了
func GetLabledBookHandler(c *gin.Context) {
	//先获取query
	label := c.Query("label")
	//如果label为空就跳转所有书籍显示
	if label == "" {
		GetBookListHandler(c)
	} else {
		booklist, err := service.GetTheLabeledBook(label)
		if err != nil {
			log.Println(err)
			utility.ResponseErr(c, utility.StatusServerBusy)
			return
		}
		//没有错误且为nil说明找不到这种书籍
		if *booklist == nil {
			utility.ResponseErr(c, utility.StatusBookNoSuchLabel)
			return
		}
		utility.ResponseSuccess(c, booklist)
		return
	}
}
