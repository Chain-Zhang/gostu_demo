package orm

import (
	"database/sql"
	"fmt"
	"strings"
)

/*
  MysqlDB对象
*/
type MysqlDB struct{
	*sql.DB
	Params
}

/*
lt  : <
lte : <=
gt  : >
gte : >=
*/

/*
  参数字段
*/
type ParamField struct{
	name string
	value interface{}
}

/*
  参数
*/
type Params struct{
	from string  //表名
	fields []string  //select的字段名
	where []ParamField // where条件，暂时只支持and连接
	values []interface{} //查询的参数值
}

/*
  获取数据库对象
*/
func NewDb(driverName, dataSourceName string)(*MysqlDB, error){
	db, err := sql.Open(driverName, dataSourceName)
	if err != nil{
		return nil, err
	}
	mydb := &MysqlDB{DB:db}
	return mydb, nil
}

/*
  插入函数
*/
func (this *MysqlDB)Insert(model interface{}) error{
	strSql, params, tbInfo, err := generateInsertSql(model)
	if err != nil{
		return err
	}
	var result sql.Result
	result, err = this.Exec(strSql, params...)
	if err != nil{
		return err
	}
	setAuto(result, tbInfo)
	return nil
}

/*
  修改函数
*/
func (this *MysqlDB)Update(model interface{}) error{
	strSql, params, err := generateUpdateSql(model)
	if err != nil{
		return err
	}
	_, err = this.Exec(strSql, params...)
	if err != nil{
		return err
	}
	return nil
}

/*
  删除函数
*/
func (this *MysqlDB)Delete(model interface{}) error{
	strSql, params, err := generateDeleteSql(model)
	if err != nil{
		return err
	}
	_, err = this.Exec(strSql, params...)
	if err != nil{
		return err
	}
	return nil
}

/*
  原生Sql查询语句查询
*/
func (this *MysqlDB)Query(sql string, params []interface{}) (*MyRows, error){
	rows, err := this.DB.Query(sql, params...)
	if err != nil{
		return nil, err
	}
	myrows := &MyRows{Rows:rows, Values:make(map[string]interface{})}
	return myrows, nil
}

//设置表名
func (this *MysqlDB) From(from string) *MysqlDB{
	this.from = from
	return this
}

//设置查询条件
func (this *MysqlDB) Where(name string, val interface{}) *MysqlDB{
	this.where = append(this.where, ParamField{name:name, value:val})
	return this
}

//设置select字段
func (this *MysqlDB) Select(args ...string) *MysqlDB{
    for _, v := range args{
		this.fields = append(this.fields, v)
	}
	return this
}

//获取参数
func (this MysqlDB)getValues()[]interface{}{
	return this.values
}

//获取查询的sql语句
func (this *MysqlDB)getSelectSql()string{
	strSql := "select "
	if this.fields == nil || len(this.fields) < 1{
		strSql += "* "
	}else{
		for _, v := range this.fields{
			strSql += v + ","
		}
		strSql = strings.TrimRight(strSql, ",")
	}
	strSql += " from " + this.from
	if this.where != nil && len(this.where) > 0{
		 strSql += " where "
		 for _, v := range this.where{
			 this.values = append(this.values, v.value)
			 if strings.Contains(v.name, "__"){
				 nameArgs := strings.Split(v.name, "__")
				 if len(nameArgs) == 2{
					 switch nameArgs[1]{
					 case "lt":
						strSql += nameArgs[0] + "<? and "
						break
					case "lte":
						strSql += nameArgs[0] + "<=? and "
						break
					case "gt":
						strSql += nameArgs[0] + " > ? and "
						break
					case "gte":
						strSql += nameArgs[0] + ">=? and "
						break
					case "eq":
						strSql += nameArgs[0] + "=? and "
						break
					 }
				 }else{
					strSql += v.name + "=? and "
				 }
			 }else{
				strSql += v.name + "=? and "
			 }
		 }
		 strSql = strings.TrimRight(strSql, " and ")
	}
	return strSql
}

//执行查询语句，并映射结果到实体切片
func (this *MysqlDB)Get()([]interface{}, error){
	sql :=  this.getSelectSql()
	fmt.Println("sql: ", sql)
	vals := this.getValues()
	fmt.Println("params: ", vals)
	myrows, err := this.Query(sql, vals)
	if err != nil{
		return nil, err
	}
	return myrows.To(this.from)
}