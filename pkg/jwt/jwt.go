package jwt

import (
	"time"

	"nautilus/pkg/conf"

	"github.com/golang-jwt/jwt"
)

// CreateToken jwt 生成 token
// token 15分钟有效
// 需要配置 TOKEN_SECRET
func CreateToken(uid uint64) (token string, err error) {
	atClaims := jwt.MapClaims{}
	atClaims["authorized"] = true
	atClaims["user_id"] = uid
	atClaims["exp"] = time.Now().Add(time.Minute * 15).Unix()
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	token, err = at.SignedString([]byte(conf.Get("TOKEN_SECRET")))
	return
}
