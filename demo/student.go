package demo

import (
	"fmt"
)

type Address struct{
	City string
	Area string
}

type Student struct{
	Address
	Name string
	Age int
}

func (this Student) Say(){
	fmt.Println("hello, i am ", this.Name, "and i am ", this.Age)
}

func (this Student) Hello(word string){
	fmt.Println("hello", word, ". i am ", this.Name)
}
