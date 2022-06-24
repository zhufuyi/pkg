package mysql

import (
	"testing"
)

var (
	gLog    = newGormLogger()
	gValues = []interface{}{
		"sql",
		"user.go:59",
		6512100,
		"SELECT * FROM `user`  WHERE `user`.`deleted_at` IS NULL AND ((name LIKE ?)) ORDER BY id desc LIMIT 15 OFFSET 0",
		[]interface{}{
			"%张%",
		},
		3,
	}
)

func BenchmarkPrint(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		gLog.Print(gValues)
	}
}
