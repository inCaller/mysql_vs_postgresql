package main

import "reflect"

func getFieldByName(data *QueryData, fieldName string) interface{} {
	return reflect.ValueOf(*data).FieldByName(fieldName).Interface()
}
