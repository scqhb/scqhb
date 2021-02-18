package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"strings"
	"time"
)

const (
	USERNAME = "root"
	PASSWORD = "oracle"
	IP       = "9.9.9.44"
	dbName   = "mytest"
	PORT     = "3306"
)

type User struct {
	Id       int
	UserName string
	PassWord string
}

var DB *sql.DB

func InitDB() {
	path := strings.Join([]string{USERNAME, ":", PASSWORD, "@tcp(", IP, ":", PORT, ")/", dbName, "?charset=utf8"}, "")
	fmt.Println("DSN:", path)
	DB, _ = sql.Open("mysql", path)
	DB.SetConnMaxLifetime(time.Second * 180)
	DB.SetMaxIdleConns(10)
	if err := DB.Ping(); err != nil {
		fmt.Println("open database fail")
		return
	}
	fmt.Println("connect success!!")
}

func InsertDB(user User) bool {
	tx, err := DB.Begin()
	if err != nil {
		fmt.Println("fx fail")
		return false
	}
	stmt, err := tx.Prepare("INSERT INTO nk_user(`name`,`password`) VALUES (?,?)")
	if err != nil {
		fmt.Println("Prepare fail")
		return false
	}
	res, err := stmt.Exec(user.UserName, user.PassWord)
	if err != nil {
		fmt.Println("exec fail")
		return false
	}
	tx.Commit()
	fmt.Println(res.LastInsertId())
	fmt.Println(res.RowsAffected())
	return true
}

func DeleteUser(user User) bool {
	tx, err := DB.Begin()
	if err != nil {
		fmt.Println("tx fail")
		return false
	}
	stmt, err := tx.Prepare("delete from nk_user where id=?")
	if err != nil {
		fmt.Println("Prepare fail")
		return false
	}
	res, err := stmt.Exec(user.Id)
	if err != nil {
		fmt.Println("exec fail")
		return false
	}
	tx.Commit()
	fmt.Println(res.LastInsertId())
	fmt.Println(res.RowsAffected())
	return true

}

func UpdateUser(user User) bool {
	tx, err := DB.Begin()
	if err != nil {
		fmt.Println(err)
		return false
	}
	stmt, err := tx.Prepare("update nk_user set name=?,password=? where id=?")
	if err != nil {
		fmt.Println(err)
		return false
	}
	result, err := stmt.Exec(user.UserName, user.PassWord, user.Id)
	if err != nil {
		fmt.Println(err)
		return false
	}
	tx.Commit()
	fmt.Println(result.LastInsertId())
	fmt.Println(result.RowsAffected())
	return true
}

func SelectDb(user User) {
	sqlStr := "select id,name,password from nk_user"
	queryRow := DB.QueryRow(sqlStr, 1)
	queryRow.Scan(&user.Id, &user.UserName, &user.PassWord)
	fmt.Println("user:%#v\n", user)
}

func main() {
	/*	InitDB()
		user := User{
			Id:       2,
			UserName: "qinhaibo222",
			PassWord: "oracle2222",
		}*/

	//UpdateUser(user)
	//InsertDB(user)
	//DeleteUser(user)
	//	SelectDb(user)

}
