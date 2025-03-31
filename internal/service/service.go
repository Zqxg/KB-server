package service

import (
	"fmt"
	"projectName/internal/repository"
	"projectName/pkg/jwt"
	"projectName/pkg/log"
	"projectName/pkg/sid"
	"reflect"
)

type Service struct {
	Logger *log.Logger
	Sid    *sid.Sid
	Jwt    *jwt.JWT
	Tm     repository.Transaction
}

func NewService(
	tm repository.Transaction,
	logger *log.Logger,
	sid *sid.Sid,
	jwt *jwt.JWT,
) *Service {
	return &Service{
		Logger: logger,
		Sid:    sid,
		Jwt:    jwt,
		Tm:     tm,
	}
}

// page初始化
func InitPage(pageIndex int, pageSize int) (int, int) {
	if pageIndex < 1 {
		pageIndex = 1
	}
	if pageSize < 10 {
		pageSize = 10
	}
	if pageSize > 100 {
		pageSize = 100
	}
	return pageIndex, pageSize
}

// GenerateEsMapping 自动生成 Elasticsearch 映射
func GenerateEsMapping(input interface{}) (string, error) {
	// 获取结构体类型
	t := reflect.TypeOf(input)
	if t.Kind() != reflect.Struct {
		return "", fmt.Errorf("expected a struct type")
	}

	// 构建映射字符串
	mapping := `{
		"settings": {
			"number_of_shards": 1,
			"number_of_replicas": 1
		},
		"mappings": {
			"properties": {`

	// 遍历结构体字段
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)

		// 获取字段名称（json 标签）
		fieldName := field.Tag.Get("json")
		if fieldName == "" {
			fieldName = field.Name
		}

		// 根据字段类型生成对应的 Elasticsearch 类型
		fieldType := field.Type.Name()
		var esType string

		switch fieldType {
		case "string":
			esType = "text"
		case "int", "int32", "int64":
			esType = "integer"
		case "bool":
			esType = "boolean"
		case "time.Time":
			esType = "date"
		case "uint":
			esType = "keyword" // 假设 uint 类型为精确匹配
		default:
			esType = "text"
		}

		// 将每个字段添加到映射
		mapping += fmt.Sprintf(`"%s": {
			"type": "%s"
		},`, fieldName, esType)
	}

	// 移除最后一个逗号并关闭 JSON 格式
	mapping = mapping[:len(mapping)-1] + `
			}
		}
	}`

	return mapping, nil
}
