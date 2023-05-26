package utils

import "golang.org/x/crypto/bcrypt"

func EncryptPassword(password string) (string, error) {
	// 使用 bcrypt 算法加密密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

func MatchPassword(password string, encryptedPassword string) error {
	// 将已加密的密码转换为 byte 类型
	byteHashedPassword := []byte(encryptedPassword)

	// 使用 bcrypt.VerifyHashAndPassword() 函数验证明文密码是否与数据库中的密码匹配
	err := bcrypt.CompareHashAndPassword(byteHashedPassword, []byte(password))
	if err != nil {
		return err
	}
	return nil
}
