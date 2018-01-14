package webser

import (
	"strings"
	"fmt"
	"net/http"
	"log"
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

func Start(){
	mux := &MyMux{}
	err := http.ListenAndServe(":9090", mux)
	if err != nil{
		log.Fatal("ListenAndServe: ", err)
	}
}