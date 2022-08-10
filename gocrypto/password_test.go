package gocrypto

import "testing"

func TestComparePasswords(t *testing.T) {
	pwd := "123"

	hashStr, err := HashAndSaltPassword(pwd)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(hashStr)

	ok := VerifyPassword(pwd, hashStr)
	if !ok {
		t.Fatal("passwords mismatch")
	}
}
