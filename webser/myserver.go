package webser

import (
	"io"
	"os"
	"encoding/hex"
	//"io"
	"crypto/md5"
	"fmt"
	"strings"
	"strconv"
	"regexp"
	"log"
	"time"

	"net/http"
	"html/template"
)

type MyMux struct{
}

func (p *MyMux)ServeHTTP(w http.ResponseWriter, r *http.Request){
	if r.URL.Path == "/"{
		sayHelloName(w, r)
		return
	}
	if r.URL.Path == "/about"{
		about(w, r)
		return
	}
	if r.URL.Path == "/login"{
		login(w,r)
		return
	}
	if r.URL.Path == "/upload"{
		upload(w,r)
		return
	}
	http.NotFound(w,r)
	return
}

func sayHelloName(w http.ResponseWriter, r *http.Request){
	r.ParseForm()
	fmt.Println(r.Form)
	fmt.Println("path: ", r.URL.Path)
	fmt.Println("scheme: ", r.URL.Scheme)
	fmt.Println(r.Form["url_long"])
	for k, v := range r.Form{
		fmt.Println("key: ", k)
		fmt.Println("val: ", strings.Join(v, " "))
	}
	fmt.Fprintf(w, "hello chain!")
}

func about(w http.ResponseWriter, r *http.Request){
	fmt.Fprintf(w, "i am chain, from shanghai")
}

func login(w http.ResponseWriter, r *http.Request){
	r.ParseForm() //解析form
	fmt.Println("method: ", r.Method)
	if r.Method == "GET"{
		time := time.Now().Unix()
		h := md5.New()
		h.Write([]byte(strconv.FormatInt(time,10)))
		token := hex.EncodeToString(h.Sum(nil))
		t, _ := template.ParseFiles("./view/login.ctpl")
		t.Execute(w, token)
	}else if r.Method == "POST"{
		token := r.Form.Get("token")
        if token != "" {
            //验证token的合法性
        } else {
            //不存在token报错
        }
		if len(r.Form["username"][0])==0{
			fmt.Fprintf(w, "username: null or empty \n")
		}
		age, err := strconv.Atoi(r.Form.Get("age"))
		if err != nil{
			fmt.Fprintf(w, "age: The format of the input is not correct \n")
		}
		if age < 18{
			fmt.Fprintf(w, "age: Minors are not registered \n")
		}

		if m, _ := regexp.MatchString(`^([\w\.\_]{2,10})@(\w{1,}).([a-z]{2,4})$`,
		    r.Form.Get("email")); !m {    
				fmt.Fprintf(w, "email: The format of the input is not correct \n")
		}
	}
}

func upload(w http.ResponseWriter, r *http.Request){
	r.ParseForm()
	if r.Method == "GET"{
        time := time.Now().Unix()
		h := md5.New()
		h.Write([]byte(strconv.FormatInt(time,10)))
		token := hex.EncodeToString(h.Sum(nil))
		t, _ := template.ParseFiles("./view/upload.ctpl")
		t.Execute(w, token)
	}else if r.Method == "POST"{
		token := r.Form.Get("token")
        if token != "" {
			//验证token的合法性
			fmt.Println("token: ", token)
        } else {
			//不存在token报错
			fmt.Println("token is nil or empty")
        }
		//把上传的文件存储在内存和临时文件中
		r.ParseMultipartForm(32 << 20)
		//获取文件句柄，然后对文件进行存储等处理
		file, handler, err := r.FormFile("uploadfile")
		if err != nil{
			fmt.Println("form file err: ", err)
			return
		}
		defer file.Close()
		fmt.Fprintf(w, "%v", handler.Header)
		//创建上传的目的文件
		f, err := os.OpenFile("./files/" + handler.Filename, os.O_WRONLY | os.O_CREATE, 0666)
		if err != nil{
			fmt.Println("open file err: ", err)
			return
		}
		defer f.Close()
		//拷贝文件
		io.Copy(f, file)
	}
}

func Start(){
	mux := &MyMux{}
	err := http.ListenAndServe(":9090", mux)
	if err != nil{
		log.Fatal("ListenAndServe: ", err)
	}
}