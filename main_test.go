package main

import (
	"math/rand"
	"testing"
)

func TestHex(t *testing.T) {
	type args struct {
		binary []byte
		upper  bool
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"[00,00,00,00]", args{[]byte{0, 0, 0, 0}, false}, "00000000"},
		{"[ff,ff,ff,ff]", args{[]byte{0xff, 0xff, 0xff, 0xff}, false},
			"ffffffff"},
		{"[01,23,45,67,89,ab,cd,ef]",
			args{[]byte{0x01, 0x23, 0x45, 0x67, 0x89, 0xab, 0xcd, 0xef}, false},
			"0123456789abcdef"},
		{"[00,00,00,00]", args{[]byte{0, 0, 0, 0}, true}, "00000000"},
		{"[ff,ff,ff,ff]", args{[]byte{0xff, 0xff, 0xff, 0xff}, true},
			"FFFFFFFF"},
		{"[01,23,45,67,89,ab,cd,ef]",
			args{[]byte{0x01, 0x23, 0x45, 0x67, 0x89, 0xab, 0xcd, 0xef}, true},
			"0123456789ABCDEF"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := hex(tt.args.binary, tt.args.upper); got != tt.want {
				t.Errorf("hex() = %v, want %v", got, tt.want)
			}
		})
	}
}

func BenchmarkHexLower(b *testing.B) {
	binary := make([]byte, 512)
	for i := 0; i < 512; i++ {
		binary[i] = byte(rand.Intn(0x100))
	}
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			hex(binary, false)
		}
	})
}

func BenchmarkHexUpper(b *testing.B) {
	binary := make([]byte, 512)
	for i := 0; i < 512; i++ {
		binary[i] = byte(rand.Intn(0x100))
	}
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			hex(binary, true)
		}
	})
}
