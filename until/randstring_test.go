package until

import "testing"

const n = 16

func BenchmarkRunes(b *testing.B) {
	for i := 0; i < b.N; i++ {
		RandStringBytesMaskImprSrcUnsafe(n)
	}
}
