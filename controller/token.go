package controller

import (
	"Web/utility"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
	"time"
)

// 产生Accesstoken
func GenAccessToken(username string) (string, error) {
	j := JWT{
		username,
		time.Now().Add(AccessTokenDuration).Unix(),
		time.Now().Unix(),
		"G-string",
		"A",
	}
	payload, err := json.Marshal(j)
	if err != nil {
		return "payload序列化失败", err
	}
	header := `{"alg": "HS256","typ": "JWT"}"`
	//只有signature部分经过了HS256运算，header和payload是只是经过了base64编码，解码就可以直接看到
	signature := HmacSha256(base64.URLEncoding.EncodeToString([]byte(header))+"."+base64.URLEncoding.EncodeToString(payload), secret)
	token := base64.URLEncoding.EncodeToString([]byte(header)) + "." + base64.URLEncoding.EncodeToString(payload) + "." + base64.URLEncoding.EncodeToString(signature)
	return token, nil
}

// 产生Refresh Token
func GenRefreshToken(username string) (string, error) {
	j := JWT{
		username,
		time.Now().Add(RefreshTokenDuration).Unix(),
		time.Now().Unix(),
		"G-string",
		"R",
	}
	payload, err := json.Marshal(j)
	if err != nil {
		return "payload序列化失败", err
	}
	header := `{"alg": "HS256","typ": "JWT"}"`
	signature := HmacSha256(base64.URLEncoding.EncodeToString([]byte(header))+"."+base64.URLEncoding.EncodeToString(payload), secret)
	token := base64.URLEncoding.EncodeToString([]byte(header)) + "." + base64.URLEncoding.EncodeToString(payload) + "." + base64.URLEncoding.EncodeToString(signature)
	return token, nil
}

func GenRefreshTokenWithTimeToCheck(username string, issuetime int64) (string, error) {
	j := JWT{
		username,
		issuetime + (time.Now().Add(RefreshTokenDuration).Unix() - time.Now().Unix()),
		issuetime,
		"G-string",
		"R",
	}
	payload, err := json.Marshal(j)
	if err != nil {
		return "payload序列化失败", err
	}
	header := `{"alg": "HS256","typ": "JWT"}"`
	signature := HmacSha256(base64.URLEncoding.EncodeToString([]byte(header))+"."+base64.URLEncoding.EncodeToString(payload), secret)
	token := base64.URLEncoding.EncodeToString([]byte(header)) + "." + base64.URLEncoding.EncodeToString(payload) + "." + base64.URLEncoding.EncodeToString(signature)
	return token, nil
}

func GenAccessTokenWithTimeToCheck(username string, issueTime int64) (string, error) {
	j := JWT{
		username,
		issueTime + (time.Now().Add(AccessTokenDuration).Unix() - time.Now().Unix()),
		issueTime,
		"G-string",
		"A",
	}
	payload, err := json.Marshal(j)
	if err != nil {
		return "payload序列化失败", err
	}
	header := `{"alg": "HS256","typ": "JWT"}"`
	//只有signature部分经过了HS256运算，header和payload是只是经过了base64编码，解码就可以直接看到
	signature := HmacSha256(base64.URLEncoding.EncodeToString([]byte(header))+"."+base64.URLEncoding.EncodeToString(payload), secret)
	token := base64.URLEncoding.EncodeToString([]byte(header)) + "." + base64.URLEncoding.EncodeToString(payload) + "." + base64.URLEncoding.EncodeToString(signature)
	return token, nil
}

// 没有验证用户是否存在
// 注意，如果是到期了那么err为nil，是其他错误err就非nil
// 先检查是否到期，没到期就把得到的token再计算生成一边Token，如何和获得的Token一样，就有效.
func VerifyAccessToken(token string) (bool, error, string) {
	//首先是base64解码变为3个部分，拿出前2个部分再次生成token与该token对比
	//base64编码里面没有’.‘这个点的字符，所以要把三部分分开解码
	parts := strings.SplitN(token, ".", 3)
	//把中间token信息部分通过base64解码
	decode, err := base64.URLEncoding.DecodeString(parts[1])
	if err != nil {
		return false, errors.New(TokenErrMap[TokenErrDecode]), ""
	}
	var jwt JWT
	//把中间部分的信息赋值到结构体内
	err = json.Unmarshal(decode, &jwt)
	if jwt.Subject != "A" {
		return false, nil, ""
	}
	//判断是否到期
	//如果到期
	if time.Now().Unix() > jwt.ExpiresAt {
		return false, nil, ""
	} else { //没有过期,就计算一遍签名
		verifyToken, err := GenAccessTokenWithTimeToCheck(jwt.Audience, jwt.IssuedAt)
		if err != nil {
			return false, errors.New(TokenErrMap[TokenGenErr]), ""
		}
		//开始验证
		//现在是说明Token是伪造的
		if verifyToken != token {
			return false, errors.New(TokenErrMap[TokenFalse]), ""
		} else { //token正确
			return true, nil, jwt.Audience
		}
	}
}

// 通用于access和refresh
// 从请求体中获取token并且提取出需要的token部分
// 没有token则err为nil，string为“”
func GetTheTokenFromHeader(c *gin.Context) (string, error) {
	//c.Request是包含从URL中获取到的所有数据的结构体。.GET可以指定获取哪一个数据
	//token就是放在请求头的Authorization的Bearer中。一般形式为Authorization: Bearer xxxxxx.xxx.xxx，我们会GET:以后的字符串，要切分
	tokenstring := c.Request.Header.Get("Authorization")
	if tokenstring == "" {
		//c.JSON(http.StatusOK, gin.H{
		//	"status": TokenNon,
		//	"info":   TokenErrMap[TokenNon],
		//})
		c.Abort()
		return "", errors.New(TokenErrMap[TokenNon])
	}
	//fmt.Printf("tokenString为%v\n", tokenstring)
	//token可能为空，为空时怎么处理，定义状态？为空时不就是没有token吗？
	parts := strings.SplitN(tokenstring, " ", 2)
	if !(parts[0] == "Bearer" && len(parts) == 2) {
		c.JSON(http.StatusOK, gin.H{
			"status": TokenErrFrom,
			"info":   TokenErrMap[TokenErrFrom],
		})
		c.Abort()
		return "", errors.New(TokenErrMap[TokenErrFrom])
	}
	return parts[1], nil
}

// 第三者是该用户的名字
func VerifyRefreshToken(token string) (bool, error, string) {
	parts := strings.SplitN(token, ".", 3)
	var jwt JWT
	decode, err := base64.URLEncoding.DecodeString(parts[1])
	if err != nil {
		return false, errors.New(utility.StatusInfoMap[utility.StatusServerBusy]), ""
	}
	err = json.Unmarshal(decode, &jwt)
	fmt.Println(jwt)
	if jwt.Subject != "R" {
		return false, nil, ""
	}
	if err != nil {
		return false, errors.New(utility.StatusInfoMap[utility.StatusServerBusy]), ""
	}
	verifytoken, err := GenRefreshTokenWithTimeToCheck(jwt.Audience, jwt.IssuedAt)
	if err != nil {
		return false, errors.New(utility.StatusInfoMap[utility.StatusServerBusy]), ""
	}
	//用户传递的refresh token不正确
	if verifytoken != token {
		return false, nil, ""
	} else {
		return true, nil, jwt.Audience
	}
}
