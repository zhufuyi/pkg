package gjwt

import (
	"fmt"
	"testing"
)

func TestGenerateToken(t *testing.T) {
	//token, err := GenerateToken("123", "456")
	//token, err := GenerateToken("", "")
	token, err := GenerateTokenID("123456")
	if err != nil {
		t.Error(err)
		return
	}

	fmt.Println(token)
}

func TestParseToken(t *testing.T) {
	//token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhcHBfa2V5IjoiMjAyY2I5NjJhYzU5MDc1Yjk2NGIwNzE1MmQyMzRiNzAiLCJhcHBfc2VjcmV0IjoiMjUwY2Y4YjUxYzc3M2YzZjhkYzhiNGJlODY3YTlhMDIiLCJleHAiOjE2MTQxNzI3MTl9.A1PUVcafb8ZkjhkxwS4ku8i3h3rBYTGbqrhV7naUdgA"
	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6IjEyMzQ1NiIsImFwcF9rZXkiOiIiLCJhcHBfc2VjcmV0IjoiIiwiZXhwIjoxNjE0MjI1MDIwfQ.pY6-6aHeuomp2ifTlIl1JUjaxuaFmP-9YLl0vfWLSvA"

	c, err := ParseToken(token)
	if err != nil {
		t.Error(err)
		return
	}

	fmt.Println(c)
}
