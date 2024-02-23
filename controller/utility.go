package controller

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"time"
)

// 我们就在备注这里区分AccessToken和RefreshToken吧????
// AccessToken带用户名，Refresh Token不带用户名来区分吧
type JWT struct {
	Audience  string `json:"aud"` //用户名
	ExpiresAt int64  `json:"exp"` //到期时间
	IssuedAt  int64  `json:"iat"` //签发日期
	Issuer    string `json:"iss"` //发行人
	Subject   string `json:"sub"` //备注
}

// 计算HS256
func HmacSha256(key string, data string) []byte {
	mac := hmac.New(sha256.New, []byte(key))
	_, _ = mac.Write([]byte(data))
	return mac.Sum(nil)
}

func HmacSha256ToHex(key string, data string) string {
	return hex.EncodeToString(HmacSha256(key, data))
}

func HmacSua256ToBase64(key string, data string) string {
	return base64.URLEncoding.EncodeToString(HmacSha256(key, data))
}

const (
	TokenErrFrom = 20000 + iota
	TokenNon
	TokenErrDecode
	TokenGenErr
	TokenFalse
)

var TokenErrMap = map[int]string{
	TokenNon:       "Token为空",
	TokenErrFrom:   "Token格式错误",
	TokenErrDecode: "Token解码错误",
	TokenGenErr:    "Token产生错误",
	TokenFalse:     "Token伪造",
}

var secret = "我不会否定过去的自己，即使过去充满失败，那也是组成现在的我的一部分"

const AccessTokenDuration = time.Minute * 10
const RefreshTokenDuration = time.Hour * 2
