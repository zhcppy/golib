package util

import "testing"

//go:generate go test -bench=. -run=none -benchmem
// -bench=. 运行所有基准测试
// -run=none 匹配一个从来没有的单元测试方法，过滤掉单元测试的输出
// -benchmem 输出每次操作分配的字节数，每次操作分配内存的次数

/*
goos: darwin
goarch: amd64
pkg: github.com/zhcppy/golib/common/util
BenchmarkBytes-12                        3132924               382 ns/op              32 B/op          2 allocs/op
BenchmarkRunes-12                        2380344               506 ns/op              96 B/op          2 allocs/op
BenchmarkBytesRmndr-12                   3422845               325 ns/op              32 B/op          2 allocs/op
BenchmarkBytesMask-12                    3504278               333 ns/op              32 B/op          2 allocs/op
BenchmarkBytesMaskImpr-12               13709001                86.3 ns/op            32 B/op          2 allocs/op
BenchmarkBytesMaskImprSrc-12            13768291                86.7 ns/op            32 B/op          2 allocs/op
PASS
ok      github.com/zhcppy/golib/common/util     8.846s
 */

const n = 16

func BenchmarkBytes(b *testing.B) {
	for i := 0; i < b.N; i++ {
		RandStringBytes(n)
	}
}

func BenchmarkRunes(b *testing.B) {
	for i := 0; i < b.N; i++ {
		RandStringRunes(n)
	}
}

func BenchmarkBytesRmndr(b *testing.B) {
	for i := 0; i < b.N; i++ {
		RandStringBytesRmndr(n)
	}
}

func BenchmarkBytesMask(b *testing.B) {
	for i := 0; i < b.N; i++ {
		RandStringBytesMask(n)
	}
}

func BenchmarkBytesMaskImpr(b *testing.B) {
	for i := 0; i < b.N; i++ {
		RandStringBytesMaskImpr(n)
	}
}

func BenchmarkBytesMaskImprSrc(b *testing.B) {
	for i := 0; i < b.N; i++ {
		RandStringBytesMaskImprSrc(n)
	}
}
