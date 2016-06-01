package main

import "reflect"

func getFieldByName(data *QueryData, query_param *Param) interface{} {
	return reflect.ValueOf(*data).FieldByName(query_param.ParamName).Interface()
}
