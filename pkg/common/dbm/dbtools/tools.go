package dbtools

import (
	"fmt"
	"github.com/beego/beego/orm"
	"main/common/tools"
	"reflect"
)

func ApplyCondition(model interface{}) *orm.Condition {
	elem := reflect.ValueOf(model).Elem()
	return ApplyConditionByReflect(elem)
}

func ApplyConditionByReflect(val reflect.Value) *orm.Condition {
	cond := orm.NewCondition()
	for i := 0; i < val.NumField(); i++ {
		typeField := val.Type().Field(i)
		valueField := val.Field(i)
		// 跳过零值字段和未导出的字段
		if !valueField.IsValid() || !valueField.CanInterface() || reflect.DeepEqual(valueField.Interface(), reflect.Zero(typeField.Type).Interface()) {
			continue
		}
		// 根据字段类型构建条件
		format := "%s__%s"
		switch valueField.Kind() {
		case reflect.String:
			cond = cond.And(fmt.Sprintf(format, tools.ToSnakeCase(typeField.Name), "contains"), valueField.String())
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			cond = cond.And(fmt.Sprintf(format, tools.ToSnakeCase(typeField.Name), "exact"), valueField.Int())
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			cond = cond.And(fmt.Sprintf(format, tools.ToSnakeCase(typeField.Name), "exact"), valueField.Uint())
		case reflect.Bool:
			cond = cond.And(fmt.Sprintf(format, tools.ToSnakeCase(typeField.Name), "exact"), valueField.Int())
		case reflect.Float32, reflect.Float64:
			cond = cond.And(fmt.Sprintf(format, tools.ToSnakeCase(typeField.Name), "exact"), valueField.Float())
		default:
			continue
		}
	}
	return cond
}

func ApplyUpdateParams(model interface{}) map[string]interface{} {
	val := reflect.ValueOf(model).Elem()
	return ApplyUpdateParamsByReflect(val)
}

func ApplyUpdateParamsByReflect(val reflect.Value) map[string]interface{} {
	params := make(map[string]interface{})
	for i := 0; i < val.NumField(); i++ {
		typeField := val.Type().Field(i)
		valueField := val.Field(i)
		// 跳过零值字段和未导出的字段
		if !valueField.IsValid() || !valueField.CanInterface() || reflect.DeepEqual(valueField.Interface(), reflect.Zero(typeField.Type).Interface()) {
			continue
		}
		switch valueField.Kind() {
		case reflect.String:
			params[tools.ToSnakeCase(typeField.Name)] = valueField.String()
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			params[tools.ToSnakeCase(typeField.Name)] = valueField.Int()
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			params[tools.ToSnakeCase(typeField.Name)] = valueField.Uint()
		case reflect.Bool:
			params[tools.ToSnakeCase(typeField.Name)] = valueField.Bool()
		case reflect.Float32, reflect.Float64:
			params[tools.ToSnakeCase(typeField.Name)] = valueField.Float()
		default:
			continue
		}
	}
	return params
}
