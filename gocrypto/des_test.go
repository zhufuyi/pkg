package gocrypto

import "testing"

var (
	desRawData = []byte("des|abcdefghijklmnopqrstuvwxyz1234567890")
	deskey     = []byte("ABCDEFGH")
)

func TestDes(t *testing.T) {
	want := desRawData

	// ECB default mod and key
	t.Run("default des ebc", func(t *testing.T) {
		cypherData, _ := DesEncrypt(desRawData) // 加密
		got, _ := DesDecrypt(cypherData)        // 解密
		if string(got) != string(want) {
			t.Fatalf("got [%s], want [%s]", got, want)
		}
		t.Logf("[%s]  <=>  [%x]", desRawData, cypherData)
	})

	// ECB
	t.Run("des ecb", func(t *testing.T) {
		cypherData, _ := DesEncrypt(desRawData, WithDesKey(deskey), WithDesModeECB()) // 加密
		got, _ := DesDecrypt(cypherData, WithDesKey(deskey), WithDesModeECB())        // 解密
		if string(got) != string(want) {
			t.Fatalf("got [%s], want [%s]", got, want)
		}
		t.Logf("[%s]  <=>  [%x]", desRawData, cypherData)
	})

	// CBC
	t.Run("des cbc", func(t *testing.T) {
		cypherData, _ := DesEncrypt(desRawData, WithDesKey(deskey), WithDesModeCBC()) // 加密
		got, _ := DesDecrypt(cypherData, WithDesKey(deskey), WithDesModeCBC())        // 解密
		if string(got) != string(want) {
			t.Fatalf("got [%s], want [%s]", got, want)
		}
		t.Logf("[%s]  <=>  [%x]", desRawData, cypherData)
	})

	// CFB
	t.Run("des cfb", func(t *testing.T) {
		cypherData, _ := DesEncrypt(desRawData, WithDesKey(deskey), WithDesModeCFB()) // 加密
		got, _ := DesDecrypt(cypherData, WithDesKey(deskey), WithDesModeCFB())        // 解密
		if string(got) != string(want) {
			t.Fatalf("got [%s], want [%s]", got, want)
		}
		t.Logf("[%s]  <=>  [%x]", desRawData, cypherData)
	})

	// CTR
	t.Run("des ctr", func(t *testing.T) {
		cypherData, _ := DesEncrypt(desRawData, WithDesKey(deskey), WithDesModeCTR()) // 加密
		got, _ := DesDecrypt(cypherData, WithDesKey(deskey), WithDesModeCTR())        // 解密
		if string(got) != string(want) {
			t.Fatalf("got [%s], want [%s]", got, want)
		}
		t.Logf("[%s]  <=>  [%x]", desRawData, cypherData)
	})
}

func BenchmarkDes(b *testing.B) {
	b.Run("des ecb encrypt", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			DesEncrypt(desRawData, WithDesModeECB())
		}
	})
	b.Run("des ecb decrypt", func(b *testing.B) {
		cypherData, err := DesEncrypt(desRawData, WithDesModeECB())
		if err != nil {
			b.Fatal(err)
		}
		var tmp []byte
		copy(tmp, cypherData)
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			copy(cypherData, tmp)
			DesDecrypt(cypherData, WithDesModeECB())
		}
	})

	b.Run("des cbc encrypt", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			DesEncrypt(desRawData, WithDesModeCBC())
		}
	})
	b.Run("des cbc decrypt", func(b *testing.B) {
		cypherData, err := DesEncrypt(desRawData, WithDesModeCBC())
		if err != nil {
			b.Fatal(err)
		}
		var tmp []byte
		copy(tmp, cypherData)
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			copy(cypherData, tmp)
			DesDecrypt(cypherData, WithDesModeCBC())
		}
	})

	b.Run("des cfb encrypt", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			DesEncrypt(desRawData, WithDesModeCFB())
		}
	})
	b.Run("des cfb decrypt", func(b *testing.B) {
		cypherData, err := DesEncrypt(desRawData, WithDesModeCFB())
		if err != nil {
			b.Fatal(err)
		}
		var tmp []byte
		copy(tmp, cypherData)
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			copy(cypherData, tmp)
			DesDecrypt(cypherData, WithDesModeCFB())
		}
	})

	b.Run("des ctr encrypt", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			DesEncrypt(desRawData, WithDesModeCTR())
		}
	})
	b.Run("des ctr decrypt", func(b *testing.B) {
		cypherData, err := DesEncrypt(desRawData, WithDesModeCTR())
		if err != nil {
			b.Fatal(err)
		}
		var tmp []byte
		copy(tmp, cypherData)
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			copy(cypherData, tmp)
			DesDecrypt(cypherData, WithDesModeCTR())
		}
	})
}

func TestDesHex(t *testing.T) {
	want := string(desRawData)

	t.Run("default des ecb", func(t *testing.T) {
		cypherStr, _ := DesEncryptHex(string(desRawData)) // 加密
		got, _ := DesDecryptHex(cypherStr)                // 解密
		if got != want {
			t.Fatalf("got [%s], want [%s]", got, want)
		}
		t.Logf("[%s]  <=>  [%s]", desRawData, cypherStr)
	})
}
