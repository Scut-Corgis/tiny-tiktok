package service

func GenerateToken(name string) string {
	return name
}
func VerifyToken(name string, token string) bool {
	return name == token
}
