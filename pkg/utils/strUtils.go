package utils

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
	"unicode"
)

// TrimWhitespace 去除字符串的前后空格
func TrimWhitespace(s string) string {
	return strings.TrimSpace(s)
}

// IsEmpty 检查字符串是否为空
func IsEmpty(s string) bool {
	return len(s) == 0
}

// IsNotEmpty 检查字符串是否非空
func IsNotEmpty(s string) bool {
	return len(s) > 0
}

// ToUpperCase 将字符串转换为大写
func ToUpperCase(s string) string {
	return strings.ToUpper(s)
}

// ToLowerCase 将字符串转换为小写
func ToLowerCase(s string) string {
	return strings.ToLower(s)
}

// Capitalize 首字母大写，其余小写
func Capitalize(s string) string {
	if len(s) == 0 {
		return s
	}
	return strings.ToUpper(string(s[0])) + strings.ToLower(s[1:])
}

// SnakeToCamel 将 snake_case 转换为 camelCase
func SnakeToCamel(s string) string {
	parts := strings.Split(s, "_")
	for i := 1; i < len(parts); i++ {
		parts[i] = Capitalize(parts[i])
	}
	return strings.Join(parts, "")
}

// CamelToSnake 将 camelCase 转换为 snake_case
func CamelToSnake(s string) string {
	var result []rune
	for i, r := range s {
		if unicode.IsUpper(r) && i > 0 {
			result = append(result, '_')
		}
		result = append(result, unicode.ToLower(r))
	}
	return string(result)
}

// Substring 截取字符串，支持负数索引
func Substring(s string, start, length int) (string, error) {
	runes := []rune(s)
	strLen := len(runes)

	if start < 0 || start >= strLen {
		return "", errors.New("start index out of range")
	}
	if length < 0 {
		return "", errors.New("length must be non-negative")
	}

	end := start + length
	if end > strLen {
		end = strLen
	}

	return string(runes[start:end]), nil
}

// Contains 判断字符串是否包含指定子字符串
func Contains(s, substr string) bool {
	return strings.Contains(s, substr)
}

// ReplaceAll 替换字符串中的所有子字符串
func ReplaceAll(s, old, new string) string {
	return strings.ReplaceAll(s, old, new)
}

// IsAlpha 判断字符串是否只包含字母
func IsAlpha(s string) bool {
	for _, r := range s {
		if !unicode.IsLetter(r) {
			return false
		}
	}
	return true
}

// IsNumeric 判断字符串是否只包含数字
func IsNumeric(s string) bool {
	for _, r := range s {
		if !unicode.IsDigit(r) {
			return false
		}
	}
	return true
}

// IsAlphanumeric 判断字符串是否只包含字母和数字
func IsAlphanumeric(s string) bool {
	for _, r := range s {
		if !unicode.IsLetter(r) && !unicode.IsDigit(r) {
			return false
		}
	}
	return true
}

// IsEmail 检查字符串是否为有效的电子邮件格式
func IsEmail(s string) bool {
	// 简单的正则表达式来验证邮箱格式
	const emailRegex = `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	re := regexp.MustCompile(emailRegex)
	return re.MatchString(s)
}

// IsPhoneNumber 验证字符串是否是一个有效的中国手机号
func IsPhoneNumber(s string) bool {
	// 中国手机号正则
	phoneRegex := `^1[0-9]\d{9}$`
	re := regexp.MustCompile(phoneRegex)
	return re.MatchString(s)
}

// Reverse 反转字符串
func Reverse(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

// ToInt 将字符串转换为整数，如果转换失败则返回错误
func ToInt(s string) (int, error) {
	var result int
	_, err := fmt.Sscanf(s, "%d", &result)
	if err != nil {
		return 0, errors.New("invalid integer format")
	}
	return result, nil
}

// ToFloat 将字符串转换为浮点数，如果转换失败则返回错误
func ToFloat(s string) (float64, error) {
	var result float64
	_, err := fmt.Sscanf(s, "%f", &result)
	if err != nil {
		return 0, err
	}
	return result, nil
}

// PadLeft 补全字符串，确保它的长度至少为指定的长度，短的部分用指定字符填充
func PadLeft(s string, length int, padChar rune) string {
	if len(s) >= length {
		return s
	}
	return string(padChar) + PadLeft(s, length-1, padChar)
}

// PadRight 补全字符串，确保它的长度至少为指定的长度，短的部分用指定字符填充
func PadRight(s string, length int, padChar rune) string {
	if len(s) >= length {
		return s
	}
	return s + string(padChar)
}
