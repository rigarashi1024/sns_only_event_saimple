package auth

import "golang.org/x/crypto/bcrypt"

const passwordHashCost = bcrypt.DefaultCost

// HashPassword は平文パスワードを bcrypt でハッシュ化します。
// パスワードは復号する必要がないため、暗号化ではなく一方向ハッシュとして保存します。
func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), passwordHashCost)
	if err != nil {
		return "", err
	}

	return string(hashedPassword), nil
}

// VerifyPassword は平文パスワードと保存済みハッシュが一致するかを検証します。
func VerifyPassword(password string, passwordHash string) bool {
	return bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(password)) == nil
}
