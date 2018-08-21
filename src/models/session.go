package models

import (
	"encoding/pem"
	"errors"
	"crypto/x509"
	"crypto/rsa"
	"crypto/rand"
	"github.com/gin-gonic/gin"
	"strings"
	"strconv"
	"time"
	"crypto/md5"
	"fmt"
)

type Session struct{
	UserName string
	Email string
	UserId string
	DecSecert string
	Password string
}

type SessionCookie struct {
	UserId string
	Secert string
}
func NewSessionCookie(c *gin.Context)*SessionCookie{
	return &SessionCookie{
		UserId:GetCUserId(c),
		Secert:GetCSecert(c),
	}
}

// 刷新用户登录态
func RefreshUserLoginInfo(c *gin.Context,userSession *Session) error{
	userPwd,err := RasEncryUserPwd(userSession.UserId,userSession.Password)
	if err!=nil{
		return err
	}

	c.SetCookie("b_secert",userPwd,86400,"/","/",false,true)
	c.SetCookie("b_uid",userSession.UserId,86400,"/","/",false,true)
	return nil
}

// 检测用户登录态信息
/*
	用户秘钥是根据user_id+md5(user_password)+timestamp组成的
*/
func CheckUserLoginInfo(c *gin.Context,sessionCookie *SessionCookie)(int,error){
	decSecert, err := RsaDecrypt([]byte("dashdausdasdasdads"), []byte(sessionCookie.Secert))
	if err!=nil{
		return 1,err
	}
	// 从secert解析用户信息
	userId,pwd,loginTime,_ := GetUserInfoFromSecert(string(decSecert))
	// 如果上次登录的时间戳有点久的话就重新登录
	if time.Now().Unix()-loginTime>GetLoginExpire(){
		return 2,errors.New("登录态已经失效")
	}
	// 验证用户登录信息
	userModel := DbUserModel{}
	userData,err := userModel.FindUserInfo(userId)
	if err!=nil{
		return 3,err
	}

	md5Instance := md5.New()
	tmpMd5Pwd := md5Instance.Sum([]byte(userData.Password))
	if  string(tmpMd5Pwd) != pwd{
		return 4,err
	}

	return 0,nil
}
// 解析secert
func GetUserInfoFromSecert(secert string) (userId,pwd string,loginTime int64,err error){
	userInfo := strings.Split(secert, "/")
	if len(userInfo) != 3 {
		return "","",0,errors.New("dec userinfo failed")
	}

	loginTime,_ = strconv.ParseInt(userInfo[2],10,64)
	return userInfo[0],userInfo[1],loginTime,nil
}
// 获取登录态失效时间
func GetLoginExpire()int64{
	// 目前返回的是3小时
	return 21600
}
// 使用公钥加密数据
func RsaEncrypt(publicKey,origData []byte)([]byte,error){
	block,_ := pem.Decode(publicKey)
	if block == nil{
		return nil,errors.New("public key error")
	}

	pubInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	pub := pubInterface.(*rsa.PublicKey)
	return rsa.EncryptPKCS1v15(rand.Reader, pub, origData)
}
// 使用私钥解密数据
func RsaDecrypt(privateKey,cipherText []byte)([]byte, error){
	block, _ := pem.Decode(privateKey)
	if block == nil {
		return nil, errors.New("private key error!")
	}
	priv, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	return rsa.DecryptPKCS1v15(rand.Reader, priv, cipherText)
}
// 加密用户秘钥
func RasEncryUserPwd(userId,pwd string) (string,error){
	md5Instance := md5.New()
	needEncryStr := fmt.Sprintf("%s_%s_%d",userId, string(md5Instance.Sum([]byte(pwd))),time.Now().Unix())

	encryPwd,err := RsaEncrypt([]byte("qweqweqweq"), []byte(needEncryStr))
	if err!=nil{
		return "",err
	}

	return string(encryPwd),nil
}
// 获取cookie中的userId
func GetCUserId(c *gin.Context)string{
	if cookie,err :=c.Request.Cookie("b_secert");err==nil{
		return cookie.Value
	}
	return ""
}
// 获取cookie中的秘钥
func GetCSecert(c *gin.Context)string{
	if cookie,err := c.Request.Cookie("b_secert");err==nil{
		return cookie.Value
	}

	return ""
}
