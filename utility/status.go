package utility

import (
	"errors"
)

// 把各种错误都写成一个状态码
const (
	StatusSuccess = 10000 + iota
	StatusUserExist
	StatusUserNotExist
	StatusUserPasswordErr
	StatusParamErr
	StatusServerBusy
	StatusAccessTokenErr
	StatusAccessTokenOutOfRange
	StatusRefreshTokenErr
	StatusRefreshTokenOutOfRange
	StatusNameChangeErr
	StatusGenderChangeErr
	StatusAvaterChangeErr
	StatusIntroductionChangeErr
	StatusPhoneChangeErr
	StatusQQChangeErr
	StatusEmailChangeErr
	StatusBirthdayChangeErr
	StatusBookQueryErr
	StatusBookNotExist
	StatusBookAlreadyStar
	StatusBookNoSuchLabel
	StatusBookCommentGetErr
	StatusBookWriteCommentErr
	StatusBookChangerAndCommenterNotEqual
	StatusGetFormDateErr
	StatusNowNotAchieveEmbededComment
	StatusGetCollectedBookErr
	StatusCaredUserNotExist
	StatusUserIsNotTheCommenter
	StatusBookCommentNotExist
)

// 由于各个地方都有相应号码，不统一，想统一的时候，已经积重难返了
var StatusInfoMap = map[int]string{
	StatusSuccess:                         "success",
	StatusUserExist:                       "该用户已存在",
	StatusUserNotExist:                    "该用户不存在",
	StatusUserPasswordErr:                 "密码错误",
	StatusParamErr:                        "请填写完整用户名或密码",
	StatusServerBusy:                      "服务繁忙",
	StatusAccessTokenErr:                  "AccessToken错误",
	StatusAccessTokenOutOfRange:           "AccessToken过期",
	StatusRefreshTokenErr:                 "RefreshToken错误",
	StatusRefreshTokenOutOfRange:          "RefreshToken过期",
	StatusNameChangeErr:                   "用户名修改错误",
	StatusGenderChangeErr:                 "性别修改错误",
	StatusAvaterChangeErr:                 "头像修改错误",
	StatusIntroductionChangeErr:           "简介修改错误",
	StatusPhoneChangeErr:                  "电话号码修改错误",
	StatusQQChangeErr:                     "QQ号码修改错误",
	StatusEmailChangeErr:                  "电子邮箱修改错误",
	StatusBirthdayChangeErr:               "生日修改错误",
	StatusBookQueryErr:                    "书籍列表查询错误",
	StatusBookNotExist:                    "书籍不存在",
	StatusBookAlreadyStar:                 "书籍已经收藏过了",
	StatusBookNoSuchLabel:                 "没有该标签的书籍",
	StatusBookCommentGetErr:               "获取书籍评论错误",
	StatusBookWriteCommentErr:             "纂写书评错误",
	StatusBookChangerAndCommenterNotEqual: "申请更改书评者和发布书评者不同",
	StatusGetFormDateErr:                  "获取Form表单数据错误",
	StatusNowNotAchieveEmbededComment:     "暂未实现嵌套评论",
	StatusGetCollectedBookErr:             "获得收藏书籍错误",
	StatusCaredUserNotExist:               "被关注的用户不存在",
	StatusUserIsNotTheCommenter:           "请求删除评论的用户不是评论发布者，无权删除",
	StatusBookCommentNotExist:             "评论不存在",
}

// 传递进来状态码status，获得对应的info
func GetInfo(status int) string {
	info, ok := StatusInfoMap[status]
	if !ok {
		return StatusInfoMap[StatusServerBusy]
	}
	return info
}

// 摘抄自mysql里定义的错误，因为会发生循环引用，只能这么做了
var (
	ErrorUserExist    = errors.New("该用户已存在")
	ErrorUserNotExist = errors.New("该用户不存在")
	//ErrorParams = errors.New("请填写完整用户名和密码")
	ErrorServerBusy = errors.New("服务繁忙")
	ErrorPassword   = errors.New("密码错误")
)

func GetErrorStatus(err error) int {
	switch err {
	case ErrorUserExist:
		return StatusUserExist
	case ErrorUserNotExist:
		return StatusUserNotExist
	case ErrorServerBusy:
		return StatusServerBusy
	case ErrorPassword:
		return StatusUserPasswordErr
	}
	return StatusServerBusy
}

//func GetTokenStatus(status int)int{
//	info , err:=TokenErrMap[status]
//	if err!=nil{
//		return StatusServerBusy
//	}
//	return info
//}
