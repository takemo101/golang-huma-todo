package shared

const Token = "token"

// トークンが一致しないかどうかを判定する
func IsInvalidToken(token string) bool {
	return token != Token
}
