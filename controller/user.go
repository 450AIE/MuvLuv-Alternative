package controller

import (
	"Web/dao/mysql"
	"Web/model"
	"Web/service"
	"Web/utility"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
)

// 路由收到请求后，进行处理
func RegisterHandler(c *gin.Context) {
	//1.参数校验，接收到的肯定是JSON数据，所以一定要一个结构体
	var p model.UserRegisterIn
	err := c.ShouldBindJSON(&p) //获取body的JSON数据并且反序列化给p
	if err != nil {
		//请求参数有误,当传递的格式不是JSON，或者和结构体字段对不上，就会返回err
		log.Println(err)
		utility.ResponseErr(c, utility.StatusServerBusy)
		return
	}
	//校验数据
	if len(p.Username) == 0 || len(p.Password) == 0 {
		utility.ResponseErr(c, utility.StatusParamErr)
		return
	}
	//2.业务处理(service处理)
	err = service.RegisterServe(&p)
	if err != nil {
		log.Println(err)
		utility.ResponseErr(c, utility.GetErrorStatus(err))
		return
	}
	//3.返回响应
	utility.ResponseSuccess(c, nil)
}

func LoginHandler(c *gin.Context) {
	var p model.UserLoginIn
	p.Username = c.Query("username")
	p.Password = c.Query("password")
	if len(p.Username) == 0 || len(p.Password) == 0 {
		utility.ResponseErr(c, utility.StatusParamErr)
		return
	}
	//会检查用户密码是否正确，用户是否存在
	err := service.LoginServe(&p)
	if err != nil {
		log.Println(err)
		utility.ResponseErr(c, utility.GetErrorStatus(err))
		return
	}
	//到了这里，说明用户存在，可以进行签发token或者鉴权token
	//如何判断是否有有效的token呢？
	token, err := GetTheTokenFromHeader(c)
	//如果没有token，签发access 和 refresh token
	if err != nil && err.Error() == TokenErrMap[TokenNon] { //token为空时，就签发token
		//签发access token
		token, err = GenAccessToken(p.Username)
		if err != nil {
			log.Println(err)
			utility.ResponseErr(c, utility.StatusServerBusy)
			return
		}
		//签发refresh token
		rtoken, err := GenRefreshToken(p.Username)
		if err != nil {
			log.Println(err)
			utility.ResponseErr(c, utility.StatusServerBusy)
			return
		}
		utility.ResponseSuccess(c, gin.H{
			"refresh_token": rtoken,
			"token":         token,
		})
		return
	} else if err != nil && err.Error() == TokenErrMap[TokenErrFrom] {
		//格式错误，估计要返回系统繁忙
		utility.ResponseErr(c, utility.StatusServerBusy)
		return
	}
	//走到这里说明携带了正常的token
	//这里验证Token，如果right直接登陆，err要看是错误的token还是到期了进行不同处理
	ok, err, _ := VerifyAccessToken(token)
	if err != nil { //token不正确或者其他的错误
		//伪造的Token
		if err.Error() == TokenErrMap[TokenFalse] {
			log.Println(err)
			utility.ResponseErr(c, utility.StatusAccessTokenErr)
			return
		} else { //其他的错误
			utility.ResponseErr(c, utility.StatusServerBusy)
			return
		}
	} else if !ok { //到期了，需要用户上传refresh token请求access token
		utility.ResponseErr(c, utility.StatusAccessTokenOutOfRange)
		return
	}
	//到这里说明AccessToken正常且没有到期
	utility.ResponseSuccess(c, token)
}

// 重新获得一个refresh token
func RefreshTokenHandler(c *gin.Context) {
	//1.检查用户的access token是否合法
	token, err := GetTheTokenFromHeader(c)
	if err != nil {
		utility.ResponseErr(c, utility.StatusServerBusy)
		return
	}
	ok, err, username := VerifyAccessToken(token)
	if err != nil { //这个err是其他错误
		utility.ResponseErr(c, utility.StatusServerBusy)
		return
	}
	//如果refresh不存在
	if !ok {
		utility.ResponseErr(c, utility.StatusAccessTokenErr)
		return
	}
	//2.若合法就签发一个refresh token附带Access Token
	rtoken, _ := GenRefreshToken(username)
	atoken, _ := GenAccessToken(username)
	utility.ResponseSuccess(c, gin.H{
		"refresh_token": rtoken,
		"token":         atoken,
	})
	return
}

func PasswordChangeHandler(c *gin.Context) {
	//先校验access token
	//有效的话。再提取token里面的信息
	//去mysql修改密码
	token, err := GetTheTokenFromHeader(c)
	if err != nil {
		fmt.Println(1)
		utility.ResponseErr(c, utility.StatusAccessTokenErr)
		return
	}
	ok, err, username := VerifyAccessToken(token)
	if err != nil {
		utility.ResponseErr(c, utility.StatusAccessTokenErr)
		return
	}
	//Access Token到期
	if !ok {
		utility.ResponseErr(c, utility.StatusAccessTokenOutOfRange)
		return
	}
	//没到期的话，再去提取请求参数
	var p model.UserChangePassWordIn
	p.Old_password = c.Query("old_password")
	p.New_password = c.Query("new_password")
	if err != nil {
		utility.ResponseErr(c, utility.StatusServerBusy)
		return
	}
	p.Username = username
	fmt.Println(p)
	err = service.PasswordChangeServe(&p)
	if err != nil {
		utility.ResponseErr(c, utility.GetErrorStatus(err))
		return
	}
	utility.ResponseSuccess(c, nil)
	return
}

// 肯定是登陆了有token才会来修改信息，所以不用设置一个if来判断是否登陆
func InfoChangeHandler(c *gin.Context) {
	//1.检查Token，然后根据Token携带的username检查用户是否存在
	//得到token
	token, err := GetTheTokenFromHeader(c)
	if err != nil {
		utility.ResponseErr(c, utility.StatusServerBusy)
		return
	}
	//验证token
	ok, err, username := VerifyAccessToken(token)
	if err != nil {
		utility.ResponseErr(c, utility.StatusServerBusy)
		return
	}
	if !ok {
		utility.ResponseErr(c, utility.StatusAccessTokenOutOfRange)
		return
	}
	//检查用户是否存在
	exist, err := mysql.CheckUserExistWithName(username)
	if err != nil {
		utility.ResponseErr(c, utility.GetErrorStatus(err))
		return
	}
	if !exist {
		utility.ResponseErr(c, utility.StatusUserNotExist)
		return
	}
	//2.获取修改后的信息，验证，都成功再去修改
	var user model.UserToChangeInfo
	err = c.ShouldBindJSON(&user)
	if err != nil {
		utility.ResponseErr(c, utility.StatusServerBusy)
		return
	}
	//3.修改的时候，如果传递的“”，则在MYSQL改为null
	err = service.InfoChangeServe(username, &user)
	if err != nil {
		utility.ResponseErr(c, utility.StatusServerBusy)
		return
	}
	utility.ResponseSuccess(c, nil)
}

func GetUserInfoHandler(c *gin.Context) {
	userID := c.Param("user_id")
	user, err := service.GetUserInfoServe(userID)
	if err != nil {
		if err.Error() == utility.StatusInfoMap[utility.StatusUserNotExist] {
			utility.ResponseErr(c, utility.StatusUserNotExist)
			return
		}
		utility.ResponseErr(c, utility.StatusServerBusy)
		return
	}
	utility.ResponseSuccess(c, user)
	return
}
