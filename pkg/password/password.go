package password

import (
	"fmt"
	"regexp"

	"golang.org/x/crypto/bcrypt"
)

// Bcrypt cost，使用默认值 = 10
const passwordCost = bcrypt.DefaultCost

// 使用bcrypt加密用户密码
func EncryptPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), passwordCost)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// 检查密码的Hash和密码是否匹配
func CompareHashAndPassword(hash, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}

// 检查密码是否符合格式要求
func ValidatePassword(password, confirmPassword string) error {
	if password != confirmPassword {
		return fmt.Errorf("password and confirm password do not match")
	}
	if !(len(password) >= 8 && len(password) <= 32) {
		return fmt.Errorf("password must be between 8 and 32 characters")
	}
	// 包含数字
	numberExp := `[0-9]+`
	// 包含字母
	strExp := `[a-zA-Z]+`
	// 包含特殊字符:非空非字母非数字
	symbolExp := `[^\w\s]+`
	// 8～32位字符，英文、数字和符号必须同时存在
	if !(regexp.MustCompile(numberExp).MatchString(password) &&
		regexp.MustCompile(strExp).MatchString(password) &&
		regexp.MustCompile(symbolExp).MatchString(password)) {
		return fmt.Errorf("password must contain at least one number, one letter, and one special character")
	}
	return nil
}
