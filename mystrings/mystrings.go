package mystrings

// reverse a string from left to righ
func Revese(s string) string {
	result := ""
	for _, v := range s {
		result = string(v) + result
	}
	return result
}
