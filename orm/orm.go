package orm

import (
	"fmt"
	"database/sql"
	"errors"
	"strings"
	"reflect"
)

/*
  表信息
*/
type TableInfo struct{
	Name string  //表名
	Fields []FieldInfo //表字段信息
	TMMap map[string]string //表字段与实体字段名映射关系，key为表字段名，val为实体字段名
}

/*
  表字段详细信息
*/
type FieldInfo struct{
	Name string
	IsPrimaryKey bool
	IsAutoGenerate bool
	Valve reflect.Value
}

/*
  实体对象信息
*/
type ModelInfo struct{
	TableInfo  // 实体对应的表信息
	TbName string // 表名称
	Model interface{} //实体实例
}

//表名
type TableName string
//表名类型
var typeTableName TableName
var tableNameType reflect.Type = reflect.TypeOf(typeTableName)

//实体映射，key为表名，val为实体信息
var ModelMapping map[string]ModelInfo

/*
  注册实体，每当有一个实体时，需要调用该方法注册。
  注册到 ModelMapping
*/
func Register(model interface{}){
	if ModelMapping == nil{
		ModelMapping = make(map[string]ModelInfo)
	}
	tbInfo, _ := getTableInfo(model)
	ModelMapping[tbInfo.Name] = ModelInfo{TbName:tbInfo.Name, Model:model}
}

/*
  根据实体通过反射获取表信息
  返回表信息
*/
func getTableInfo(model interface{})(tabInfo *TableInfo, err error){
	defer func(){
        if e := recover(); err != nil{
			tabInfo = nil
			err = e.(error)
		}
	}()

	err = nil
	tabInfo = &TableInfo{}
	tabInfo.TMMap = make(map[string]string)
	rt := reflect.TypeOf(model)
	rv := reflect.ValueOf(model)

	tabInfo.Name = rt.Name()
	if rt.Kind() == reflect.Ptr{
		rt = rt.Elem()
		rv = rv.Elem()
	}
	//字段解析
	for i, j := 0, rt.NumField(); i < j; i++{
		rtf := rt.Field(i)
		rvf := rv.Field(i)
		if rtf.Type == tableNameType{
			tabInfo.Name = string(rtf.Tag)
			continue
		}
		if rtf.Tag == "-"{
			continue
		}
		//解析字段的tag
		var f FieldInfo
		//没有tag,表字段名与实体字段ing一致
		if rtf.Tag == ""{
			f = FieldInfo{Name:rtf.Name, IsAutoGenerate:false, IsPrimaryKey:false, Valve:rvf}
			tabInfo.TMMap[rtf.Name] = rtf.Name
		}else{
			strTag := string(rtf.Tag)
			if strings.Index(strTag, ":") == -1{
				//tag中没有":"时，表字段名与实体字段ing一致
				f = FieldInfo{Name:rtf.Name, IsAutoGenerate:false, IsPrimaryKey:false, Valve:rvf}
				tabInfo.TMMap[rtf.Name] = rtf.Name
			}else{
				//解析tag中的name值为表字段名
			    strName := rtf.Tag.Get("name")
			    if strName == ""{
				    strName = rtf.Name
				}
				//解析tag中的PK
			    isPk := false
			    strIspk := rtf.Tag.Get("PK")
			    if strIspk == "true"{
				    isPk = true
				}
				//解析tag中的auto
			    isAuto := false
			    strIsauto := rtf.Tag.Get("auto")
			    if strIsauto == "true"{
                    isAuto = true
			    }
				f = FieldInfo{Name:strName, IsPrimaryKey:isPk, IsAutoGenerate:isAuto, Valve:rvf}
				tabInfo.TMMap[strName] = rtf.Name
		    }
		}
		tabInfo.Fields = append(tabInfo.Fields, f)
	}
	return
}

/*
  根据实体生成插入语句
*/
func generateInsertSql(model interface{})(string, []interface{}, *TableInfo, error){
	//获取表信息
	tbInfo, err := getTableInfo(model)
	if err != nil{
		return "", nil, nil, err
	}
	if len(tbInfo.Fields) == 0 {
		return "", nil, nil, errors.New(tbInfo.Name + "结构体中没有字段")
	}

	//根据字段信息拼Sql语句，以及参数值
	strSql := "insert into " + tbInfo.Name
	strFileds := ""
	strValues := ""
	var params []interface{}
	for _, v := range tbInfo.Fields{
		if v.IsAutoGenerate {
			continue
		}
		strFileds += v.Name + ","
		strValues += "?,"
		params = append(params, v.Valve.Interface())
	}
	if strFileds == ""{
		return "", nil, nil, errors.New(tbInfo.Name + "结构体中没有字段，或只有自增字段")
	}
	strFileds = strings.TrimRight(strFileds, ",")
	strValues = strings.TrimRight(strValues, ",")
	strSql += " (" + strFileds + ") values(" + strValues + ")"
	fmt.Println("sql: ",strSql)
	fmt.Println("params: ",params)
	return strSql, params, tbInfo, nil
}

/*
  根据实体生成修改的sql语句
*/
func generateUpdateSql(model interface{})(string, []interface{}, error){
	//获取表信息
	tbInfo, err := getTableInfo(model)
	if err != nil{
		return "", nil, err
	}
	if len(tbInfo.Fields) == 0 {
		return "", nil, errors.New(tbInfo.Name + "结构体中没有字段")
	}
	//根据字段信息拼Sql语句，以及参数值
	strSql := "update " + tbInfo.Name + " set "
	strFileds := ""
	strWhere := ""
	var p interface{}
	var params []interface{}
	for _, v := range tbInfo.Fields{
		if v.IsAutoGenerate && !v.IsPrimaryKey{
			continue
		}
		if v.IsPrimaryKey{
			strWhere += v.Name + "=?"
			p = v.Valve.Interface()
			continue
		}
        strFileds += v.Name + "=?,"
		params = append(params, v.Valve.Interface())
	}
	params = append(params, p)
	strFileds = strings.TrimRight(strFileds, ",")
	strSql += strFileds + " where " + strWhere
	fmt.Println("update sql: ", strSql)
	fmt.Println("update params: ", params)
	return strSql, params, nil
}

/*
  自动生成删除的sql语句，以主键为删除条件
*/
func generateDeleteSql(model interface{})(string, []interface{}, error){
	//获取表信息
	tbInfo, err := getTableInfo(model)
	if err != nil{
		return "", nil, err
	}
	//根据字段信息拼Sql语句，以及参数值
	strSql := "delete from " + tbInfo.Name + " where "
	var idVal interface{}
	for _, v := range tbInfo.Fields{
		if v.IsPrimaryKey{
			strSql += v.Name + "=?"
			idVal = v.Valve.Interface()
		}
	}
	params := []interface{}{idVal}
	fmt.Println("update sql: ", strSql)
	fmt.Println("update params: ", params)
	return strSql, params, nil
}

/*
  设置自增长字段的值
*/
func setAuto(result sql.Result, tbInfo *TableInfo)(err error){
    defer func(){
        if e := recover(); e != nil{
			err = e.(error)
		}
	}()
	id, err := result.LastInsertId()
	if id == 0{
		return
	}
	if err != nil{
		return
	}
	for _, v := range tbInfo.Fields{
		if v.IsAutoGenerate && v.Valve.CanSet(){
			v.Valve.SetInt(id)
			break
		}
	}
	return
}