package demo

import (
	"fmt"
	"reflect"
)

/*
  获取对象的信息
*/
func StructInfo(o interface{}){
	//获取对象的类型
	t := reflect.TypeOf(o)
	fmt.Println(t.Name(), "object type: ", t.Name())

	if k := t.Kind(); k != reflect.Struct{
		fmt.Println("the object is not a struct, but it is", t.Kind())
		return
	}

	//获取对象的值
	v := reflect.ValueOf(o)
	fmt.Println(t.Name(), "object value: ", v)

	//获取对象的字段
	fmt.Println(t.Name(), "fields: ")
	for i := 0; i < t.NumField(); i++{
		f := t.Field(i)
		val := v.Field(i).Interface()
		fmt.Printf("%6s:%v = %v \n", f.Name, f.Type, val)
		//通过递归调用获取子类型的信息
		t1 := reflect.TypeOf(val)
		if k := t1.Kind(); k == reflect.Struct{
			StructInfo(val)
		}
	}
    //获取对象的函数
	fmt.Println(t.Name(), "methods: ", t.NumMethod())
	for i := 0; i < t.NumMethod(); i++{
		m := t.Method(i)
		fmt.Printf("%10s:%v \n", m.Name, m.Type)
	}
}

/*
  匿名字段的反射
*/
func Annoy(o interface{}){
	t := reflect.TypeOf(o)
	for i := 0; i < t.NumField(); i++{
		f := t.Field(i)
		fmt.Printf("%10s:%#v \n", f.Name, f)
	}
}

/*
  通过反射设置字段
*/
func ReflectSet(o interface{}){
	v := reflect.ValueOf(o)
	if v.Kind() == reflect.Ptr && !v.Elem().CanSet(){
		fmt.Println("修改失败")
		return
	}
	v = v.Elem()
	//获取字段
	f := v.FieldByName("Name")
	if !f.IsValid(){
		fmt.Println("修改失败")
		return
	}
	//设置值
	if f.Kind() == reflect.String{
		f.SetString("chairis")
	}
}

/* 
  通过反射调用函数
*/
func ReflectMethod(o interface{}){
	v := reflect.ValueOf(o)
	//无参函数调用
	m1:= v.MethodByName("Say")
	m1.Call([]reflect.Value{})

	//有参函数调用
	m2 := v.MethodByName("Hello")
	m2.Call([]reflect.Value{reflect.ValueOf("iris")})
}