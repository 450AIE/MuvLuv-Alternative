package controller

import (
	"Web/model"
	"Web/service"
	"Web/utility"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
)

// 这里的逻辑乱了，屎山
func GetBookCommentHandler(c *gin.Context) {
	//如果有Access token,那么就找到对应用户的对该评论的操作信息
	token, err := GetTheTokenFromHeader(c)
	if err != nil && token != "" {
		log.Println(err)
		utility.ResponseErr(c, utility.StatusServerBusy)
		return
	}
	if token != "" { //要注意哪些评论是用户点了赞的要针对性显示
		exist, err, username := VerifyAccessToken(token)
		if err != nil {
			log.Println(err)
			if err.Error() == utility.StatusInfoMap[utility.StatusAccessTokenOutOfRange] {
				utility.ResponseErr(c, utility.StatusAccessTokenOutOfRange)
				return
			}
			utility.ResponseErr(c, utility.StatusServerBusy)
			return
		}
		//提取参数
		bookID := c.Param("book_id")
		if exist {
			bookComment, err := service.GetBookComment(username, bookID)
			if err != nil {
				log.Println(err)
				utility.ResponseErr(c, utility.StatusBookCommentGetErr)
				return
			}
			fmt.Println(bookComment)
			utility.ResponseSuccess(c, bookComment)
			return
		} else {
			bookComment, err := service.GetBookComment("", bookID)
			if err != nil {
				log.Println(err)
				utility.ResponseErr(c, utility.StatusBookCommentGetErr)
				return
			}
			utility.ResponseSuccess(c, bookComment)
			return
		}
	} else { //token为空时
		bookID := c.Param("book_id")
		bookComment, err := service.GetBookComment("", bookID)
		if err != nil {
			log.Println(err)
			utility.ResponseErr(c, utility.StatusBookCommentGetErr)
			return
		}
		fmt.Println(bookComment)
		utility.ResponseSuccess(c, bookComment)
		return
	}
}

func WriteBookCommentHandler(c *gin.Context) {
	//检验token
	token, err := GetTheTokenFromHeader(c)
	if err != nil {
		log.Println(err)
		utility.ResponseErr(c, utility.StatusServerBusy)
		return
	}
	exist, err, username := VerifyAccessToken(token)
	if err != nil {
		log.Println(err)
		utility.ResponseErr(c, utility.StatusServerBusy)
		return
	}
	if !exist {
		utility.ResponseErr(c, utility.StatusAccessTokenOutOfRange)
		return
	} else { //access token没过期
		//获取body和path的参数
		var bookcomment model.WriteBookComment
		err := c.ShouldBindJSON(&bookcomment)
		bookcomment.Book_id = c.Param("book_id")
		if err != nil {
			log.Println(err)
			utility.ResponseErr(c, utility.StatusServerBusy)
			return
		}
		//开始写书评
		commentID, err := service.WriteBookComment(username, bookcomment.Book_id, bookcomment.Content)
		if err != nil {
			log.Println(err)
			utility.ResponseErr(c, utility.StatusBookWriteCommentErr)
			return
		}
		utility.ResponseSuccess(c, commentID)
		return
	}
}

// 删除帖子要删除对应用户点赞信息，那么还需要一张表存储点赞情况
func DeleteBookCommentHandler(c *gin.Context) {
	//验证token
	token, err := GetTheTokenFromHeader(c)
	if err != nil {
		utility.ResponseErr(c, utility.StatusServerBusy)
		return
	}
	exist, err, username := VerifyAccessToken(token)
	if err != nil {
		utility.ResponseErr(c, utility.StatusServerBusy)
		return
	}
	if !exist {
		utility.ResponseErr(c, utility.StatusAccessTokenOutOfRange)
		return
	} else { //有access token
		//解析参数
		postID := c.Param("comment_id")
		err = service.DeleteBookComment(username, postID)
		if err != nil {
			if err.Error() == utility.StatusInfoMap[utility.StatusUserIsNotTheCommenter] {
				utility.ResponseErr(c, utility.StatusUserIsNotTheCommenter)
				return
			}
			if err.Error() == utility.StatusInfoMap[utility.StatusBookCommentNotExist] {
				utility.ResponseErr(c, utility.StatusBookCommentNotExist)
				return
			}
		}
		utility.ResponseSuccess(c, nil)
		return
	}
}

func ChangeBookCommentHandler(c *gin.Context) {
	//先检验token
	token, err := GetTheTokenFromHeader(c)
	if err != nil {
		log.Println(err)
		utility.ResponseErr(c, utility.StatusServerBusy)
		return
	}
	exist, err, username := VerifyAccessToken(token)
	if err != nil {
		log.Println(err)
		utility.ResponseErr(c, utility.StatusAccessTokenErr)
		return
	}
	if !exist {
		utility.ResponseErr(c, utility.StatusAccessTokenOutOfRange)
		return
	} else {
		//获取参数
		var comment model.ChangeBookComment
		err := c.ShouldBindJSON(&comment)
		if err != nil {
			utility.ResponseErr(c, utility.StatusServerBusy)
			return
		}
		comment.Comment_id = c.Param("comment_id")
		//记得判断一下申请修改的用户和当初写这个的用户是否是同一个人,以及书籍ID和书评对应的ID是否一样
		err = service.ChangeBookComment(comment.Comment_id, username, comment.Content)
		if err != nil {
			if err.Error() == utility.StatusInfoMap[utility.StatusBookChangerAndCommenterNotEqual] {
				utility.ResponseErr(c, utility.StatusBookChangerAndCommenterNotEqual)
				return
			}
			utility.ResponseErr(c, utility.StatusServerBusy)
			return
		}
		utility.ResponseSuccess(c, nil)
		return
	}
}
