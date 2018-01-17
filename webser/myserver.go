package webser

import (
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
		//io.WriteString(h, strconv.FormatInt(time,10))
		token := fmt.Sprintf("s%", hex.EncodeToString(h.Sum(nil)))
		t, _ := template.ParseFiles("./view/login.ctpl")
		t.Execute(w, token)
	}else if r.Method == "POST"{
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

func Start(){
	mux := &MyMux{}
	err := http.ListenAndServe(":9090", mux)
	if err != nil{
		log.Fatal("ListenAndServe: ", err)
	}
}