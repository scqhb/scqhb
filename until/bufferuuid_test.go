package until

import "testing"

func BenchmarkBufferUuidv4(b *testing.B) {
	for i := 0; i < b.N; i++ {
	}
}
