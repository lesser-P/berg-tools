package utils

import (
	"fmt"
	"reflect"
	"strings"
)

func CopyStruct[T any, E any](dst T, src E) (err error) {
	dstVal := reflect.ValueOf(dst)
	srcVal := reflect.ValueOf(src)

	if dstVal.Kind() != reflect.Ptr || dstVal.IsNil() {
		return fmt.Errorf("dst is not a pointer or is nil")
	}

	if srcVal.Kind() == reflect.Ptr {
		srcVal = srcVal.Elem()
	}
	dstVal = dstVal.Elem()
	return copyByJSONTags(dstVal, srcVal)
}

func copyByJSONTags(dst, src reflect.Value) error {
	srcType := src.Type()
	// 构建源结构体的标签映射
	srcTag := make(map[string]reflect.StructField)
	for i := 0; i < srcType.NumField(); i++ {
		field := srcType.Field(i)
		tag := field.Tag.Get("json")
		if tag != "" && tag != "-" {
			if strings.Split(tag, ",")[0] != "" {
				srcTag[tag] = field
			}
		}
	}

	dstType := dst.Type()
	for i := 0; i < dstType.NumField(); i++ {
		dstField := dstType.Field(i)
		dstTag := pareJsonTag(dstField)
		if dstTag == "" {
			continue
		}

		srcField, exit := srcTag[dstTag]
		if !exit {
			continue
		}

		srcValue := src.FieldByIndex(srcField.Index)
		dstValue := dst.Field(i)
		err := processField(dstValue, srcValue)
		if err != nil {
			return fmt.Errorf("字段%s拷贝失败%v", dstTag, err)
		}
	}
	return nil
}

func processField(dst, src reflect.Value) error {
	if src.Kind() == reflect.Ptr {
		if src.IsNil() {
			return nil
		}
		src = src.Elem()
	}

	if dst.Kind() == reflect.Ptr {
		if dst.IsNil() {
			if dst.CanSet() {
				dst.Set(reflect.New(dst.Type().Elem()))
			}
		}
		dst = dst.Elem()
	}

	if src.Kind() == reflect.Struct && dst.Kind() == reflect.Struct {
		return copyByJSONTags(dst, src)
	}

	if src.Type() != dst.Type() {
		return nil
	}

	if dst.CanSet() {
		dst.Set(src)
	}
	return nil
}

func pareJsonTag(field reflect.StructField) string {
	tag := field.Tag.Get("json")
	if tag == "" || tag == "-" {
		return ""
	}
	return strings.Split(tag, ",")[0]
}
