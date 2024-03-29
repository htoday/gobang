package mysql

import (
	"gobang/app/api/global"
	"gobang/app/api/internal/model"
)

/*
CREATE TABLE `gobangUsers` (
	`id` BIGINT(20) NOT NULL AUTO_INCREMENT,
	`username` VARCHAR(30) DEFAULT 'unnamed',
	`password` VARCHAR(30) DEFAULT 'unset',
	`nickname` VARCHAR(30) DEFAULT 'player',
	`email` VARCHAR(40) DEFAULT 'null',
	`starAmount` BIGINT(20) DEFAULT '0',
	PRIMARY KEY(`id`)
)ENGINE=InnoDB AUTO_INCREMENT=1 CHARSET=utf8mb4;
`avatar` LONGBLOB,
CREATE TABLE user_followers (
id INT AUTO_INCREMENT PRIMARY KEY,
user_id BIGINT,
follower_id BIGINT,
FOREIGN KEY (user_id) REFERENCES gobangUsers(id),
FOREIGN KEY (follower_id) REFERENCES gobangUsers(id)
);
*/

func CheckUser(username1 string) (bool, error) {
	var userExists bool
	query := "SELECT EXISTS(SELECT 1 FROM gobangUsers WHERE username = ?) AS user_exists"
	err := global.MysqlDB.Get(&userExists, query, username1)
	if err != nil {
		global.Logger.Warn("mysql CheckUser failed" + err.Error())
	}
	return userExists, err
}
func AddNewUser(username string, password string, nickname string) (error, int) {
	sqlStr := "insert into gobangUsers(username,password,nickname) values (?,?,?)"
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
	sqlstr := "SELECT * FROM gobangUsers where username = ?"
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
func GetUserInformation(username string) (model.User, error) {
	sqlstr := "SELECT * FROM gobangUsers where username = ?"
	var user1 model.User
	err := global.MysqlDB.Get(&user1, sqlstr, username)
	if err != nil {
		return user1, err
	}
	return user1, nil
}
func GetNickname(username string) (string, error) {
	sqlstr := "SELECT nickname FROM gobangUsers where username = ?"
	var nickname string
	err := global.MysqlDB.Get(&nickname, sqlstr, username)
	return nickname, err
}
func GetUid(username string) (int, error) {
	sqlstr := "SELECT id FROM gobangUsers where username = ?"
	var id int
	err := global.MysqlDB.Get(&id, sqlstr, username)
	return id, err
}
func AddFriend(targetUid int, followerUid int) error {
	sqlStr := "insert into user_followers(user_id,follower_id) values (?,?)"
	_, err := global.MysqlDB.Exec(sqlStr, targetUid, followerUid)
	if err != nil {
		return err
	}
	return nil
}
func DelFriend(targetUid int, followerUid int) error {
	sqlStr := "DELETE FROM user_followers WHERE user_id = ? AND follower_id=?"
	_, err := global.MysqlDB.Exec(sqlStr, targetUid, followerUid)
	if err != nil {
		return err
	}
	return nil
}
func GetUserInformationFromID(id int) (model.User, error) {
	sqlstr := "SELECT * FROM gobangUsers where id = ?"
	var user1 model.User
	err := global.MysqlDB.Get(&user1, sqlstr, id)
	if err != nil {
		return user1, err
	}
	return user1, nil
}
func UpdateStar(winner string, loser string) {
	_, _ = global.MysqlDB.Exec("UPDATE gobangUsers SET starAmount = starAmount + 1 WHERE username = ?", winner)
	sqlstr := "SELECT starAmount FROM gobangUsers where username = ?"
	var loserStar int
	_ = global.MysqlDB.Get(&loserStar, sqlstr, loser)
	if loserStar == 0 {
		return
	} else {
		_, _ = global.MysqlDB.Exec("UPDATE gobangUsers SET starAmount = starAmount + 1 WHERE username = ?", loser)
	}
}
