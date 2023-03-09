package logger

import (
	"fmt"
	"testing"
)

var aa = []int64{1, 2, 3, 4, 5, 6, 7, 8, 9, 0, 7, 6, 4, 3, 45, 65, 5, 4, 43, 4, 5, 56, 67, 7, 5, 4, 6, 565}

func bv(p int, aa *[]int64) int {
	//	aa := "asdajkshdkajsfhduiksdhdaishdaksjhdkasjhdiuewreoieuoiqweui234lk43nuc90zxucj032n4lkanczixuc90qejdeqlkdmnlaskdjlaksjd"
	le := len(*aa)
	re := 0
	for i := 0; i < 100000; i++ {
		re = i + le
	}
	return re
}

func bv2(p int, aa *[]int64) int {
	//	aa := "asdajkshdkajsfhduiksdhdaishdaksjhdkasjhdiuewreoieuoiqweui234lk43nuc90zxucj032n4lkanczixuc90qejdeqlkdmnlaskdjlaksjd"
	re := 0
	for i := 0; i < 100000; i++ {
		re = i + len(*aa)
	}
	return re
}

func TestAA(t *testing.T) {
	var success bool
	defer func() {
		fmt.Println(success)
	}()
	success = true
}

func BenchmarkLen(b *testing.B) {
	aa := []int64{1, 2, 3, 4, 5, 6, 7, 8, 9, 0, 7, 6, 4, 3, 45, 65, 5, 4, 43, 4, 5, 56, 67, 7, 5, 4, 6, 565}
	for i := 0; i < b.N; i++ {
		bv(10, &aa)
	}
}

func BenchmarkLen2(b *testing.B) {
	aa := []int64{1, 2, 3, 4, 5, 6, 7, 8, 9, 0, 7, 6, 4, 3, 45, 65, 5, 4, 43, 4, 5, 56, 67, 7, 5, 4, 6, 565}
	for i := 0; i < b.N; i++ {
		bv2(10, &aa)
	}
}
