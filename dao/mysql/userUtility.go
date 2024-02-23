package mysql

import (
	"Web/model"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"log"
	"strconv"
)

// 有该用户了返回ture，没有返回false
func CheckUserExistWithName(name string) (bool, error) {
	check := "select count(id) from users where nickname = ?"
	row := DB.QueryRow(check, name)
	var num int
	err := row.Scan(&num) //把所得到的数据赋值给num
	if err != nil {
		log.Println(err)
		return false, err
	}
	if num == 1 {
		return true, nil
	} else {
		return false, nil
	}
}

// 这个ID是自增ID，不是雪花算法的ID
func CheckUserExistWithID(userID string) (bool, error) {
	operation := "select count(id) from users where id = ?"
	row := DB.QueryRow(operation, userID)
	var num int
	err := row.Scan(&num)
	if err != nil {
		log.Println(err)
		return false, err
	}
	if num == 1 {
		return true, nil
	} else {
		return false, nil
	}
}

func encrypt(password string) string {
	hash := md5.New()                    //创建一个新的md5加密哈希对象
	hash.Write([]byte(password))         //把要加密的字符串传入
	hashValue := hash.Sum(nil)           //计算该字符串的哈希值
	return hex.EncodeToString(hashValue) //这段代码就是把
	//哈希值转换为16进制字符串
}

func CheckUserExistWithNameAndPassword(name string, password string) (bool, error) {
	//通过名字找md5加密后的密码
	operation := "select password from users where nickname =?"
	var SQLpassword string
	row := DB.QueryRow(operation, name)
	err := row.Scan(&SQLpassword)
	if err != nil { //这里应该就是没有该用户
		return false, nil
	}
	//如果用户输入的密码加密后和数据库存储的密码一致，证明密码正确
	if encrypt(password) == SQLpassword {
		return true, nil
	} else { //用户存在，但密码错误
		return false, ErrorPassword
	}
}

func NameChange(user *model.UserToChangeInfoForMySQL) error {
	operation := "update users set nickname=? where nickname=?"
	_, err := DB.Exec(operation, user.Nickname, user.OldNickname)
	if err != nil {
		return err
	}
	return nil
}

func GenderChange(user *model.UserToChangeInfoForMySQL) error {
	operation := "update users set gender=? where nickname=?"
	_, err := DB.Exec(operation, user.Gender, user.OldNickname)
	if err != nil {
		return err
	}
	return nil
}

func QQChange(user *model.UserToChangeInfoForMySQL) error {
	operation := "update users set QQ=? where nickname=?"
	QQnumber, err := strconv.ParseInt(user.QQ, 10, 64)
	if err != nil {
		return err
	}
	_, err = DB.Exec(operation, QQnumber, user.OldNickname)
	if err != nil {
		return err
	}
	return nil
}

func EmailChange(user *model.UserToChangeInfoForMySQL) error {
	operation := "update users set email=? where nickname=?"
	_, err := DB.Exec(operation, user.Email, user.OldNickname)
	if err != nil {
		return err
	}
	return nil
}

func BirthdayChange(user *model.UserToChangeInfoForMySQL) error {
	operation := "update users set birthday =? where nickname=?"
	_, err := DB.Exec(operation, user.Birthday, user.OldNickname)
	if err != nil {
		return err
	}
	return nil
}

func PhoneChange(user *model.UserToChangeInfoForMySQL) error {
	operation := "update users set phone =? where nickname=?"
	number, err := strconv.ParseInt(user.TelePhone, 10, 64)
	if err != nil {
		return err
	}
	_, err = DB.Exec(operation, number, user.OldNickname)
	if err != nil {
		return err
	}
	return nil
}

func IntroductionChange(user *model.UserToChangeInfoForMySQL) error {
	operation := "update users set introduction =? where nickname=?"
	_, err := DB.Exec(operation, user.Introduction, user.OldNickname)
	if err != nil {
		return err
	}
	return nil
}

func AvatarChange(user *model.UserToChangeInfoForMySQL) error {
	operation := "update users set avatar =? where nickname=?"
	_, err := DB.Exec(operation, user.Avatar, user.OldNickname)
	if err != nil {
		return err
	}
	return nil
}

// 通过自增ID，非雪花算法加密的ID获取
func GetUserInfoNow(user *model.UserInfoGetForMySQL) (model.User, error) {
	operation := "select gender,nickname,QQ,birthday,email,avatar,introduction,phone from users where id=?"
	rows, err := DB.Query(operation, user.UserID)
	if err != nil {
		log.Println(err)
		return model.User{}, err
	}
	var UserInfo = model.User{}
	for rows.Next() {
		err := rows.Scan(&UserInfo.Gender, &UserInfo.Nickname, &UserInfo.QQ, &UserInfo.Birthday, &UserInfo.Email, &UserInfo.Avatar, &UserInfo.Introduction, &UserInfo.Phone)
		if err != nil {
			fmt.Println(err)
			return UserInfo, err
		}
	}
	ID, err := strconv.ParseInt(user.UserID, 10, 64)
	if err != nil {
		return model.User{}, err
	}
	UserInfo.Id = int(ID)
	return UserInfo, nil
}

// 检查用户是否关注另一个用户，关注返回true，没有返回false，出错才有err
func CheckUserCare(caringUserID string, caredUserID string) (bool, error) {
	operation := "select count(ID) from user_care where cared_user_id=? and caring_user_id=?"
	row := DB.QueryRow(operation, caredUserID, caringUserID)
	var num int
	err := row.Scan(&num)
	if err != nil {
		return false, err
	}
	if num == 1 {
		return true, nil
	} else {
		return false, nil
	}
}

// 是返回true，不是返回false，其他返回err
func CheckWhetherUserIsCommenter(userID string, postID string) (bool, error) {
	operation := "select user_id from book_comment where post_id=?"
	row := DB.QueryRow(operation, postID)
	var user_id string
	err := row.Scan(&user_id)
	if err != nil {
		return false, err
	}
	if user_id == userID {
		return true, nil
	} else {
		return false, nil
	}
}
