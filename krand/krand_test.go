package krand

import (
	"fmt"
	"testing"
)

func TestInt(t *testing.T) {
	l := 100

	fmt.Println("随机生成默认的随机数: [0, 100]")
	for i := 0; i < l; i++ {
		fmt.Printf("%d ", Int())
	}

	fmt.Println("\n\n", "随机生成数: [0, max]")
	for i := 0; i < l; i++ {
		fmt.Printf("%d ", Int(20))
	}

	fmt.Println("\n\n", "随机生成数: [min, max]")
	for i := 0; i < l; i++ {
		fmt.Printf("%d ", Int(10, 20))
	}

	fmt.Println("\n\n", "随机生成数: [max, min]")
	for i := 0; i < l; i++ {
		fmt.Printf("%d ", Int(2000, 1000))
	}
}

func TestFloat64(t *testing.T) {
	l := 100

	fmt.Println("随机生成默认的随机数: [0, 100]")
	for i := 0; i < l; i++ {
		fmt.Printf("%f ", Float64(0))
	}

	fmt.Println("\n\n", "随机生成数: [0, max]")
	for i := 0; i < l; i++ {
		fmt.Printf("%f ", Float64(1, 20))
	}

	fmt.Println("\n\n", "随机生成数: [min, max]")
	for i := 0; i < l; i++ {
		fmt.Printf("%f ", Float64(2, 10, 20))
	}

	fmt.Println("\n\n", "随机生成数: [max, min]")
	for i := 0; i < l; i++ {
		fmt.Printf("%f ", Float64(4, 2000, 1000))
	}
}

func TestString(t *testing.T) {
	//fmt.Printf("%p, %s\n", kinds, kinds)

	// 随机纯数字
	fmt.Printf("%s\n", String(R_NUM))
	fmt.Printf("%s\n", String(R_NUM, 100))

	// 随机大写字母
	fmt.Printf("%s\n", String(R_UPPER))
	fmt.Printf("%s\n", String(R_UPPER, 100))

	// 随机小写字母
	fmt.Printf("%s\n", String(R_LOWER))
	fmt.Printf("%s\n", String(R_LOWER, 100))

	// 随机数字、大写字母
	fmt.Printf("%s\n", String(R_NUM|R_UPPER))
	fmt.Printf("%s\n", String(R_NUM|R_UPPER, 100))

	// 随机数字、小写字母
	fmt.Printf("%s\n", String(R_NUM|R_LOWER))
	fmt.Printf("%s\n", String(R_NUM|R_LOWER, 100))

	// 随机数字、大写字母、小写字母
	fmt.Printf("%s\n", String(R_All))
	fmt.Printf("%s\n", String(R_All, 100))
}

func BenchmarkString_NUM_6(b *testing.B) {
	for i := 0; i < b.N; i++ {
		String(R_NUM)
	}
}

func BenchmarkString_UPPER_6(b *testing.B) {
	for i := 0; i < b.N; i++ {
		String(R_UPPER)
	}
}

func BenchmarkString_LOWER_6(b *testing.B) {
	for i := 0; i < b.N; i++ {
		String(R_LOWER)
	}
}

func BenchmarkString_NUM_UPPER_16(b *testing.B) {
	for i := 0; i < b.N; i++ {
		String(R_NUM|R_UPPER, 16)
	}
}

func BenchmarkString_NUM_LOWER_16(b *testing.B) {
	for i := 0; i < b.N; i++ {
		String(R_NUM|R_LOWER, 16)
	}
}

func BenchmarkString_UPPER_LOWER_16(b *testing.B) {
	for i := 0; i < b.N; i++ {
		String(R_UPPER|R_LOWER, 16)
	}
}

func BenchmarkString_ALL_16(b *testing.B) {
	for i := 0; i < b.N; i++ {
		String(R_All, 16)
	}
}
