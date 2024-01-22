package mysql

import (
	"gobang/app/api/global"
	"gobang/app/api/internal/model"
)

func CheckUser(username1 string) (bool, error) {
	var userExists bool
	query := "SELECT EXISTS(SELECT 1 FROM users WHERE username = ?) AS user_exists"
	err := global.MysqlDB.Get(&userExists, query, username1)
	if err != nil {
		global.Logger.Warn("mysql CheckUser failed" + err.Error())
	}
	return userExists, err
}
func AddNewUser(username string, password string, nickname string) (error, int) {
	sqlStr := "insert into users(username,password,nickname) values (?,?,?)"
	ret, err := global.MysqlDB.Exec(sqlStr, username, password, nickname)
	if err != nil {
		return err, -1
	}
	newID, err := ret.LastInsertId() // 新插入数据的id
	if err != nil {
		return err, -1
	}
	return err, int(newID)
}
func CheckPassword(username string, password string) (bool, error) {
	sqlstr := "SELECT * FROM users where username = ?"
	var user1 model.User
	err := global.MysqlDB.Get(&user1, sqlstr, username)
	if err != nil {
		return false, err
	}
	if password == user1.Password {
		return true, err
	}
	return false, err
}
