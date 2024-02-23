package controller

import (
	"Web/service"
	"Web/utility"
	"github.com/gin-gonic/gin"
	"log"
)

// 我们没有实现嵌套评论的功能，所以model只会为1
// 1是点赞评论，2是点赞的评论的评论
// 要注意检查这个人是否已经点赞过该评论了，因为在查看评论时会显示是否点赞过了
// 点赞两次就是取消
func ThumbsUpHandler(c *gin.Context) {
	//检验token
	token, err := GetTheTokenFromHeader(c)
	if err != nil {
		utility.ResponseErr(c, utility.StatusAccessTokenErr)
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
	} else {
		//获取参数
		err := c.Request.ParseForm()
		if err != nil {
			utility.ResponseErr(c, utility.StatusGetFormDateErr)
			return
		}
		model := c.Request.Form.Get("model")
		TargetID := c.Request.Form.Get("target_id")
		err = service.ThumbsUp(username, model, TargetID)
		if err != nil {
			if err.Error() == "暂未实现嵌套评论" {
				utility.ResponseErr(c, utility.StatusNowNotAchieveEmbededComment)
				return
			}
			utility.ResponseErr(c, utility.StatusServerBusy)
			return
		}
		utility.ResponseSuccess(c, nil)
		return
	}
}

func GetUserCollectListHandler(c *gin.Context) {
	//验证tokne
	token, err := GetTheTokenFromHeader(c)
	if err != nil {
		log.Println(err)
		utility.ResponseErr(c, utility.StatusAccessTokenErr)
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
	} else {
		CollectList, err := service.GetUserCollectedBookList(username)
		if err != nil {
			log.Println(err)
			utility.ResponseErr(c, utility.StatusGetCollectedBookErr)
			return
		}
		if CollectList == nil {
			utility.ResponseSuccess(c, "您未收藏任何书籍")
			return
		}
		utility.ResponseSuccess(c, CollectList)
		return
	}
}

func UserFocusHandler(c *gin.Context) {
	//验证tokne
	token, err := GetTheTokenFromHeader(c)
	if err != nil {
		log.Println(err)
		utility.ResponseErr(c, utility.StatusAccessTokenErr)
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
	} else {
		//获取参数
		var caredUserID string
		err := c.Request.ParseForm()
		if err != nil {
			log.Println(err)
			utility.ResponseErr(c, utility.StatusServerBusy)
			return
		}
		caredUserID = c.Request.Form.Get("user_id")
		err = service.UserFocus(username, caredUserID)
		if err != nil {
			if err.Error() == utility.StatusInfoMap[utility.StatusCaredUserNotExist] {
				utility.ResponseErr(c, utility.StatusCaredUserNotExist)
				return
			}
		}
		utility.ResponseSuccess(c, nil)
		return
	}
}
