package checker

import (
	"testing"

	"github.com/syed-imam/elle-go/generator"
)

func BenchmarkCheckSerial(b *testing.B) {
	h := generator.Serial(generator.Config{Keys: []string{"x", "y"}, Txns: 300, MaxMops: 4}, 1)
	b.ResetTimer()
	for b.Loop() {
		if _, err := Check(h); err != nil {
			b.Fatal(err)
		}
	}
}
