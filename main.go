package main

import (
	"strings"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"mime/multipart"
	"bytes"
	"time"
	"fmt"
	"gostu_demo/orm"
	"gostu_demo/webser"
	_ "github.com/go-sql-driver/mysql"
)


func main(){
    go uploadfile("d:/photo/cut.png", "http://localhost:9090/upload")
	webser.Start()
}

func uploadfile(filename string, url string){
    time.Sleep(30 * time.Second)
	bodyBuf := &bytes.Buffer{}
	bodyWriter :=multipart.NewWriter(bodyBuf)
	//模拟创建form表单字段
	strs := strings.Split(filename, "/")
	destname := strs[len(strs) - 1]
	filewriter, err := bodyWriter.CreateFormFile("uploadfile", destname)
	if err != nil{
		fmt.Println("error writing to buffer")
		return
	}
    //打开文件句柄操作
	fh, err := os.Open(filename)
	if err != nil{
		fmt.Println("error open file")
		return
	}
	defer fh.Close()

	//拷贝文件
	_, err = io.Copy(filewriter, fh)
	if err != nil{
		fmt.Println("error copy file")
		return
	}
	contentType := bodyWriter.FormDataContentType()
	bodyWriter.Close()

	resp, err := http.Post(url, contentType, bodyBuf)
	if err != nil{
		fmt.Println("error post buffer")
		return
	}
	defer resp.Body.Close()
	resp_body, err := ioutil.ReadAll(resp.Body)
	if err != nil{
		fmt.Println("error read all")
		return
	}
	fmt.Println(resp.Status)
	fmt.Println(string(resp_body))
}


type UserInfo struct{
	TableName orm.TableName "userinfo"
	UserName string `name:"username"`
	Uid int `name:"uid"PK:"true"auto:"true"`
	DepartName string `name:"departname"`
	Created string `name:"created"`
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