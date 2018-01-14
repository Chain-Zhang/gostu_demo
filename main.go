package main

import (
	"io"
	"bufio"
	"os"
	"time"
	"fmt"
	"gostu_demo/orm"
	_ "github.com/go-sql-driver/mysql"
	"os/exec"
	"strings"
)

type UserInfo struct{
	TableName orm.TableName "userinfo"
	UserName string `name:"username"`
	Uid int `name:"uid"PK:"true"auto:"true"`
	DepartName string `name:"departname"`
	Created string `name:"created"`
}

func main(){
	name := "cmd"
	params := []string{"/c", "adb", "devices"}
	 if exeCommand(name, params){
		 fmt.Println("success")
	 }else{
		 fmt.Println("error")
	 }
}

func exeCommand(cmdName string, params []string) bool{
	cmd := exec.Command(cmdName, params...)
	fmt.Printf("执行命令: %s\n", strings.Join(cmd.Args[2:], " ")) 
	stdout, err := cmd.StdoutPipe()
	if err != nil{
		fmt.Println(os.Stderr, "error => ", err.Error())
		return false
	}
	cmd.Start()
	reader := bufio.NewReader(stdout)
	var index int

	var contents []string
	for {
		line, err := reader.ReadString('\n')
		if err != nil || err == io.EOF{
			break
		}
		fmt.Println(line)
		index ++
		contents = append(contents, line)
	}
	cmd.Wait()
	return true
}

func ormTest(){
	ui := UserInfo{UserName:"CHAIN", DepartName:"TEST", Created:time.Now().String()}
	orm.Register(new(UserInfo))
	db, err := orm.NewDb("mysql", "root:password@tcp(xxx.xx.xxx.xxxx:3306)/demo?charset=utf8")
	if err != nil {
		fmt.Println("打开SQL时出错:", err.Error())
		return
	}
	defer db.Close()
	
	//插入测试
	err = db.Insert(&ui)
	if err != nil {
		fmt.Println("插入时错误:", err.Error())
	}
	fmt.Println("插入成功")
    //修改测试
	ui.UserName = "BBBB"
	err = db.Update(ui)
	if err != nil {
		fmt.Println("修改时错误:", err.Error())
	}
	fmt.Println("修改成功")
    //删除测试
	err = db.Delete(ui)
	if err != nil {
		fmt.Println("删除时错误:", err.Error())
	}
	fmt.Println("删除成功")
	//查询测试
	res, err := db.From("userinfo").
	Select("username", "departname", "uid").
	Where("uid__gt", 20).
	Where("username", "chain").Get()
	if err != nil{
		fmt.Println("err: ", err.Error())
	}
	fmt.Println(res)
}