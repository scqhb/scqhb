package until

import "testing"

func BenchmarkSc(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Sc()
	}
}
