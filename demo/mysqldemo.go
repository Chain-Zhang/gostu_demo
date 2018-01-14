package demo

import (
	"time"
	"fmt"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

var (
	dbhost = "106.14.224.234:3306"
	dbusername = "root"
	dbpassword = "chain.123"
	dbname = "demo"
)

func Trans(){
	db := GetDB()
	defer db.Close();

	tx, err := db.Begin()
	if err != nil{
		fmt.Println("db.Begin error: ", err.Error())
		return
	}
    isCommit := true
	defer func(){
		if isCommit{
			tx.Commit()
			fmt.Println("commit")
		}else{
			tx.Rollback()
			fmt.Println("Rollback")
		}
	}()
	_, err = tx.Exec("insert into userinfo(username,departname,created) values(?,?,?)","username","departname",time.Now())
	if err != nil{
        isCommit = false
	}
	_, err = tx.Exec("insert into userinfo(username,departname,created) values(?,?,?)","username","departname",time.Now())
	if err != nil{
        isCommit = false
	}
	_, err = tx.Exec("insert into userinfo(username,departname,created) values(?,?,?)","username","departname",time.Now())
	if err != nil{
        isCommit = false
	}
}

/*
  获取sql.DB对象
*/
func GetDB() *sql.DB{
    db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8", dbusername, dbpassword, dbhost, dbname))
	CheckErr(err)
	return db
}

/*
  插入数据
*/
func Insert(username, departname, method string)bool{
	db := GetDB()
	defer db.Close()

	if method == "1"{
		_, err := db.Exec("insert into userinfo(username,departname,created) values(?,?,?)",username,departname,time.Now())
		if err != nil{
			fmt.Println("insert err: ", err.Error())
			return false
		}
		fmt.Println("insert success!")
		return true
	}else if method == "2"{
		stmt, err := db.Prepare("INSERT userinfo SET username=?,departname=?,created=?")
		if err != nil{
			fmt.Println("insert prepare error: ", err.Error())
			return false
		}
		_, err = stmt.Exec(username, departname, time.Now())
		if err != nil{
			fmt.Println("insert exec error: ", err.Error())
			return false
		}
		fmt.Println("insert success!")
		return true
	}
	return false
}

/*
  根据id，修改名称
*/
func UpdateName(id int, name string)bool{
	db := GetDB()
	defer db.Close()

	stmt, err := db.Prepare("update userinfo set username=? where uid=?")
    if err != nil{
		fmt.Println("update name prepare error: ", err.Error())
		return false
	}
	_, err = stmt.Exec(name, id)
	if err != nil{
		fmt.Println("update name exec error: ", err.Error())
		return false
	}
	fmt.Println("update name success!")
	return true
}

/*
  根据id删除数据
*/
func Delete(id int) bool {
	db := GetDB()
	defer db.Close()

	stmt, err := db.Prepare("delete from userinfo where uid=?")
    if err != nil{
		fmt.Println("delete prepare error: ", err.Error())
		return false
	}
	_, err = stmt.Exec(id)
	if err != nil{
		fmt.Println("delete exec error: ", err.Error())
		return false
	}
	fmt.Println("delete success!")
	return true
}

func GetOne(id int){
	db := GetDB()
	defer db.Close()
	var username, departname, created string
	err := db.QueryRow("select username, departname, created from userinfo where uid=?",id).Scan(&username, &departname, &created)
    if err != nil{
		fmt.Println("get one error: ", err.Error())
		return
	}
	fmt.Println("username: ", username, "departname: ", departname, "created: ", created)
}

func GetAll(){
	db := GetDB()
	defer db.Close()
	
	rows, err := db.Query("select username, departname, created from userinfo")
    if err != nil{
		fmt.Println("get all error: ", err.Error())
		return
	}
	for rows.Next(){
		var username, departname, created string
		if err := rows.Scan(&username, &departname, &created); err == nil{
            fmt.Println("username: ", username, "departname: ", departname, "created: ", created)
		}
	}
	
}

func CheckErr(err error){
	if err != nil{
		fmt.Println("err: ", err.Error())
		panic(err)
	}
}
