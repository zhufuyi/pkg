package krand

import (
	"fmt"
	"testing"
)

func TestRandInt(t *testing.T) {
	l := 100

	fmt.Println("随机生成默认的随机数: [0, 100]")
	for i := 0; i < l; i++ {
		fmt.Printf("%d ", RandInt())
	}

	fmt.Println("\n\n", "随机生成数: [0, max]")
	for i := 0; i < l; i++ {
		fmt.Printf("%d ", RandInt(20))
	}

	fmt.Println("\n\n", "随机生成数: [min, max]")
	for i := 0; i < l; i++ {
		fmt.Printf("%d ", RandInt(10, 20))
	}

	fmt.Println("\n\n", "随机生成数: [max, min]")
	for i := 0; i < l; i++ {
		fmt.Printf("%d ", RandInt(2000, 1000))
	}
}

func TestKRand(t *testing.T) {
	//fmt.Printf("%p, %s\n", kinds, kinds)

	// 随机纯数字
	fmt.Printf("%s\n", KRand(R_NUM))
	fmt.Printf("%s\n", KRand(R_NUM, 100))

	// 随机大写字母
	fmt.Printf("%s\n", KRand(R_UPPER))
	fmt.Printf("%s\n", KRand(R_UPPER, 100))

	// 随机小写字母
	fmt.Printf("%s\n", KRand(R_LOWER))
	fmt.Printf("%s\n", KRand(R_LOWER, 100))

	// 随机数字、大写字母
	fmt.Printf("%s\n", KRand(R_NUM|R_UPPER))
	fmt.Printf("%s\n", KRand(R_NUM|R_UPPER, 100))

	// 随机数字、小写字母
	fmt.Printf("%s\n", KRand(R_NUM|R_LOWER))
	fmt.Printf("%s\n", KRand(R_NUM|R_LOWER, 100))

	// 随机数字、大写字母、小写字母
	fmt.Printf("%s\n", KRand(R_All))
	fmt.Printf("%s\n", KRand(R_All, 100))
}

func BenchmarkKRand_NUM_6(b *testing.B) {
	for i := 0; i < b.N; i++ {
		KRand(R_NUM)
	}
}

func BenchmarkKRand_UPPER_6(b *testing.B) {
	for i := 0; i < b.N; i++ {
		KRand(R_UPPER)
	}
}

func BenchmarkKRand_LOWER_6(b *testing.B) {
	for i := 0; i < b.N; i++ {
		KRand(R_LOWER)
	}
}

func BenchmarkKRand_NUM_UPPER_16(b *testing.B) {
	for i := 0; i < b.N; i++ {
		KRand(R_NUM|R_UPPER, 16)
	}
}

func BenchmarkKRand_NUM_LOWER_16(b *testing.B) {
	for i := 0; i < b.N; i++ {
		KRand(R_NUM|R_LOWER, 16)
	}
}

func BenchmarkKRand_UPPER_LOWER_16(b *testing.B) {
	for i := 0; i < b.N; i++ {
		KRand(R_UPPER|R_LOWER, 16)
	}
}

func BenchmarkKRand_ALL_16(b *testing.B) {
	for i := 0; i < b.N; i++ {
		KRand(R_All, 16)
	}
}
