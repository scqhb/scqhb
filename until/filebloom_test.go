package until

import (
	"testing"
)

func BenchmarkReadfiletoMap(b *testing.B) {
	for i := 0; i < b.N; i++ {
		//Uuidv4()
		//uuid.New()
		ReadfiletoMap()
	}
}
