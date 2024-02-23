package mysql

import (
	"Web/model"
	"Web/utility"
	"errors"
	"fmt"
	"log"
)

var (
	ErrorUserExist    = errors.New("该用户已存在")
	ErrorUserNotExist = errors.New("该用户不存在")
	//ErrorParams = errors.New("请填写完整用户名和密码")
	ErrorServerBusy = errors.New("服务繁忙")
	ErrorPassword   = errors.New("密码错误")
)

func LoginUser(in *model.UserLoginToMySQL) error {
	exist, err := CheckUserExistWithNameAndPassword(in.Username, in.Password)
	if errors.Is(err, ErrorPassword) {
		return ErrorPassword
	}
	if !exist {
		return ErrorUserNotExist
	} else {
		return nil
	}
}

// 将用户加入MySQL
func AddUser(in *model.UserRegisterToMySQL) error {
	operation := "insert into users(nickname,password,user_id) values(?,?,?)"
	//不可以存储明文的密码，要对密码加密
	encpassword := encrypt(in.Password)
	fmt.Println(encpassword)
	_, err := DB.Exec(operation, in.Username, encpassword, GenUID())
	if err != nil {
		log.Println(err)
		return ErrorServerBusy
	}
	return nil
}

// 往MySQL里面加入注册的用户
func RegisterUser(in *model.UserRegisterToMySQL) error {
	//已有该用户
	exist, err := CheckUserExistWithName(in.Username)
	if err != nil { //查询错误
		return err
	}
	if exist { //用户已存在
		return ErrorUserExist
	} else { //用户不存在，可以开始注册
		err = AddUser(in)
		return err
	}
}

func PasswordChange(in *model.UserChangePassWordInToMySQL) error {
	//验证用户是否存在(感觉可有可无一样)
	//修改密码
	//返回响应
	ok, err := CheckUserExistWithNameAndPassword(in.Username, in.Old_password)
	if err != nil {
		return err
	}
	//如果用户不存在
	if !ok {
		return ErrorUserNotExist
	}
	//如果用户存在且密码正确，进行修改
	err = passwordChang(in)
	if err != nil {
		return err
	}
	return nil
}

func passwordChang(in *model.UserChangePassWordInToMySQL) error {
	operation := "update users set password = ? where nickname = ?"
	newpassword := encrypt(in.New_password)
	_, err := DB.Exec(operation, newpassword, in.Username)
	if err != nil {
		return ErrorServerBusy
	}
	return nil
}

// 如果用户想把某个值设置为空，那么可能要指明填写某个特殊字符了。我们就简单粗暴的吧""的不修改
func InfoChange(in *model.UserToChangeInfoForMySQL) error {
	//因为用户可能只修改某些值，有些值是“”，要修改不为“”的值
	if in.Nickname != "" {
		err := NameChange(in)
		if err != nil {
			return err
		}
	}
	if in.Avatar != "" {
		err := AvatarChange(in)
		if err != nil {
			return err
		}
	}
	if in.Introduction != "" {
		err := IntroductionChange(in)
		if err != nil {
			return err
		}
	}
	if in.TelePhone != "" {
		err := PhoneChange(in)
		if err != nil {
			return err
		}
	}
	if in.QQ != "" {
		err := QQChange(in)
		if err != nil {
			return err
		}
	}
	if in.Gender != "" {
		err := GenderChange(in)
		if err != nil {
			return err
		}

	}
	if in.Email != "" {
		err := EmailChange(in)
		if err != nil {
			return err
		}
	}
	if in.Birthday != "" {
		err := BirthdayChange(in)
		if err != nil {
			return err
		}
	}
	return nil
}

func GetUserInfo(user *model.UserInfoGetForMySQL) (model.User, error) {
	//检查用户是否存在
	exist, err := CheckUserExistWithID(user.UserID)
	//查找错误
	if err != nil {
		return model.User{}, errors.New(utility.StatusInfoMap[utility.StatusServerBusy])
	}
	//用户不存在
	if !exist {
		return model.User{}, errors.New(utility.StatusInfoMap[utility.StatusUserNotExist])
	}
	UserInfo, err := GetUserInfoNow(user)
	if err != nil {
		return model.User{}, err
	}
	return UserInfo, nil
}

func FocusUser(in *model.UserFocusForMySQL) error {
	//先确保一下这个被关注的用户存在
	exist, err := CheckUserExistWithID(in.CaredUserID)
	if err != nil {
		log.Println(err)
		return err
	}
	if !exist {
		return errors.New(utility.StatusInfoMap[utility.StatusCaredUserNotExist])
	}
	operation := "insert into user_care (caring_user_id,cared_user_id)value(?,?)"
	res, err := DB.Exec(operation, in.CaringUserID, in.CaredUserID)
	//什么都没有修改时，res为nil
	if res == nil { //说明之前已经关注了，现在取消关注
		err = CancelTheFocus(in)
		if err != nil {
			return err
		}
		return nil
	}
	if err != nil {
		return err
	}
	return nil
}

func CancelTheFocus(in *model.UserFocusForMySQL) error {
	operation := "delete from user_care where caring_user_id=? and cared_user_id=?"
	_, err := DB.Exec(operation, in.CaringUserID, in.CaredUserID)
	if err != nil {
		return err
	}
	return nil
}
