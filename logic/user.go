package logic

import (
	"go-jichu/dao/mysql"
	"go-jichu/models"
	"go-jichu/pkg/jwt"
	"go-jichu/pkg/snowflake"
)

func SingUp(p *models.ParamSignUp) (err error) {
	//判断用户是否存在
	err = mysql.CheckUserExist(p.Username)
	if err != nil {
		return err
	}
	//生成uid
	userId := snowflake.GetID()
	//构造一个user实例
	user := &models.User{
		UserID:   userId,
		UserName: p.Username,
		Password: p.Password,
	}
	//保存数据库
	return mysql.InsertUser(user)
}

func Login(p *models.ParamLogin) (user *models.User, err error) {
	user = &models.User{
		UserName: p.Username,
		Password: p.Password,
		Token:    "",
	}
	//传递的是一个指针，就可以拿到userid
	if err := mysql.Login(user); err != nil {
		return nil, err
	}
	//生成jwt的token
	token, err := jwt.GenToken(user.UserID, user.UserName)
	if err == nil {
		user.Token = token
	}
	return
}
