package gocrypto

import (
	"crypto"
	"testing"
)

var hashRawData = []byte("hash|abcdefghijklmnopqrstuvwxyz1234567890")

func TestMd5(t *testing.T) {
	val := Md5(hashRawData)
	want := "25bf4c81ab3eca5a78287ce8bdbeb34f"
	if val != want {
		t.Fatalf("got %v, want %v", val, want)
	}
	t.Log(val)
}

func TestSha1(t *testing.T) {
	val := Sha1(hashRawData)
	want := "f21388e39122fbb04b37a479a604282e9874a995"
	if val != want {
		t.Fatalf("got %v, want %v", val, want)
	}
	t.Log(val)
}

func TestSha256(t *testing.T) {
	val := Sha256(hashRawData)
	want := "8fe0ace6a1aef456a9223f3b29298f5d992838cd9d35e4fc5dfb5d65baff9486"
	if val != want {
		t.Fatalf("got %v, want %v", val, want)
	}
	t.Log(val)
}

func TestSha512(t *testing.T) {
	val := Sha512(hashRawData)
	want := "9c19e2b14668f8f8479e8c554f69a4d65fa2537039764f8af5aa730f60698d3e114e1fc9a1cac11bcae65fd6437ab121f9e6971adda9c3142b05bccda3bed82b"
	if val != want {
		t.Fatalf("got %v, want %v", val, want)
	}
	t.Log(val)
}

func BenchmarkMd5(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Md5(hashRawData)
	}
}

func BenchmarkSha1(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Sha1(hashRawData)
	}
}

func BenchmarkSha256(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Sha256(hashRawData)
	}
}

func BenchmarkSha512(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Sha512(hashRawData)
	}
}

func TestHash(t *testing.T) {
	type args struct {
		hashType crypto.Hash
		rawData  []byte
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "MD4",
			args: args{
				hashType: crypto.MD4,
				rawData:  hashRawData,
			},
			wantErr: false,
		},
		{
			name: "MD5",
			args: args{
				hashType: crypto.MD5,
				rawData:  hashRawData,
			},
			wantErr: false,
		},
		{
			name: "SHA1",
			args: args{
				hashType: crypto.SHA1,
				rawData:  hashRawData,
			},
			wantErr: false,
		},
		{
			name: "SHA224",
			args: args{
				hashType: crypto.SHA224,
				rawData:  hashRawData,
			},
			wantErr: false,
		},
		{
			name: "SHA256",
			args: args{
				hashType: crypto.SHA256,
				rawData:  hashRawData,
			},
			wantErr: false,
		},
		{
			name: "SHA384",
			args: args{
				hashType: crypto.SHA384,
				rawData:  hashRawData,
			},
			wantErr: false,
		},
		{
			name: "SHA512",
			args: args{
				hashType: crypto.SHA512,
				rawData:  hashRawData,
			},
			wantErr: false,
		},
		{
			name: "MD5SHA1",
			args: args{
				hashType: crypto.MD5SHA1,
				rawData:  hashRawData,
			},
			wantErr: false,
		},
		{
			name: "RIPEMD160",
			args: args{
				hashType: crypto.RIPEMD160,
				rawData:  hashRawData,
			},
			wantErr: false,
		},
		{
			name: "SHA3_224",
			args: args{
				hashType: crypto.SHA3_224,
				rawData:  hashRawData,
			},
			wantErr: false,
		},
		{
			name: "SHA3_256",
			args: args{
				hashType: crypto.SHA3_256,
				rawData:  hashRawData,
			},
			wantErr: false,
		},
		{
			name: "SHA3_384",
			args: args{
				hashType: crypto.SHA3_384,
				rawData:  hashRawData,
			},
			wantErr: false,
		},
		{
			name: "SHA3_512",
			args: args{
				hashType: crypto.SHA3_512,
				rawData:  hashRawData,
			},
			wantErr: false,
		},
		{
			name: "SHA512_224",
			args: args{
				hashType: crypto.SHA512_224,
				rawData:  hashRawData,
			},
			wantErr: false,
		},
		{
			name: "SHA512_256",
			args: args{
				hashType: crypto.SHA512_256,
				rawData:  hashRawData,
			},
			wantErr: false,
		},
		{
			name: "BLAKE2s_256",
			args: args{
				hashType: crypto.BLAKE2s_256,
				rawData:  hashRawData,
			},
			wantErr: false,
		},
		{
			name: "BLAKE2b_256",
			args: args{
				hashType: crypto.BLAKE2b_256,
				rawData:  hashRawData,
			},
			wantErr: false,
		},
		{
			name: "BLAKE2b_384",
			args: args{
				hashType: crypto.BLAKE2b_384,
				rawData:  hashRawData,
			},
			wantErr: false,
		},
		{
			name: "BLAKE2b_512",
			args: args{
				hashType: crypto.BLAKE2b_512,
				rawData:  hashRawData,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Hash(tt.args.hashType, tt.args.rawData)
			if (err != nil) != tt.wantErr {
				t.Errorf("Hash() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			t.Log(got)
		})
	}
}

func BenchmarkHash(b *testing.B) {
	b.Run("MD4", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			Hash(crypto.MD4, hashRawData)
		}
	})

	b.Run("MD5", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			Hash(crypto.MD5, hashRawData)
		}
	})

	b.Run("SHA1", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			Hash(crypto.SHA1, hashRawData)
		}
	})

	b.Run("SHA224", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			Hash(crypto.SHA224, hashRawData)
		}
	})

	b.Run("SHA256", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			Hash(crypto.SHA256, hashRawData)
		}
	})

	b.Run("SHA384", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			Hash(crypto.SHA384, hashRawData)
		}
	})

	b.Run("SHA512", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			Hash(crypto.SHA512, hashRawData)
		}
	})

	b.Run("MD5SHA1", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			Hash(crypto.MD5SHA1, hashRawData)
		}
	})

	b.Run("RIPEMD160", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			Hash(crypto.RIPEMD160, hashRawData)
		}
	})

	b.Run("SHA3_224", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			Hash(crypto.SHA3_224, hashRawData)
		}
	})

	b.Run("SHA3_256", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			Hash(crypto.SHA3_256, hashRawData)
		}
	})

	b.Run("SHA3_384", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			Hash(crypto.SHA3_384, hashRawData)
		}
	})

	b.Run("SHA3_512", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			Hash(crypto.SHA3_512, hashRawData)
		}
	})

	b.Run("SHA512_224", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			Hash(crypto.SHA512_224, hashRawData)
		}
	})

	b.Run("SHA512_256", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			Hash(crypto.SHA512_256, hashRawData)
		}
	})

	b.Run("BLAKE2s_256", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			Hash(crypto.BLAKE2s_256, hashRawData)
		}
	})

	b.Run("BLAKE2b_256", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			Hash(crypto.BLAKE2b_256, hashRawData)
		}
	})

	b.Run("BLAKE2b_384", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			Hash(crypto.BLAKE2b_384, hashRawData)
		}
	})

	b.Run("BLAKE2b_512", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			Hash(crypto.BLAKE2b_512, hashRawData)
		}
	})
}

/*
cpu: Intel(R) Core(TM) i7-8700 CPU @ 3.20GHz
BenchmarkHash
BenchmarkHash/MD4
BenchmarkHash/MD4-12             2896478               413.7 ns/op
BenchmarkHash/MD5
BenchmarkHash/MD5-12             5462320               217.5 ns/op
BenchmarkHash/SHA1
BenchmarkHash/SHA1-12            4606863               261.0 ns/op
BenchmarkHash/SHA224
BenchmarkHash/SHA224-12          3577316               339.1 ns/op
BenchmarkHash/SHA256
BenchmarkHash/SHA256-12          3522372               339.2 ns/op
BenchmarkHash/SHA384
BenchmarkHash/SHA384-12          2648484               454.6 ns/op
BenchmarkHash/SHA512
BenchmarkHash/SHA512-12          2426125               486.1 ns/op
BenchmarkHash/MD5SHA1
BenchmarkHash/MD5SHA1-12         3083764               389.0 ns/op
BenchmarkHash/RIPEMD160
BenchmarkHash/RIPEMD160-12       2202613               545.2 ns/op
BenchmarkHash/SHA3_224
BenchmarkHash/SHA3_224-12        1374512               862.6 ns/op
BenchmarkHash/SHA3_256
BenchmarkHash/SHA3_256-12        1401825               861.7 ns/op
BenchmarkHash/SHA3_384
BenchmarkHash/SHA3_384-12        1352902               873.8 ns/op
BenchmarkHash/SHA3_512
BenchmarkHash/SHA3_512-12        1344568               884.1 ns/op
BenchmarkHash/SHA512_224
BenchmarkHash/SHA512_224-12      2845068               427.3 ns/op
BenchmarkHash/SHA512_256
BenchmarkHash/SHA512_256-12      2807102               430.7 ns/op
BenchmarkHash/BLAKE2s_256
BenchmarkHash/BLAKE2s_256-12     2989849               381.1 ns/op
BenchmarkHash/BLAKE2b_256
BenchmarkHash/BLAKE2b_256-12     2453086               492.0 ns/op
BenchmarkHash/BLAKE2b_384
BenchmarkHash/BLAKE2b_384-12     2326148               516.4 ns/op
BenchmarkHash/BLAKE2b_512
BenchmarkHash/BLAKE2b_512-12     2190686               544.3 ns/op
*/
