package until

import (
	"fmt"
	"github.com/google/uuid"
	"testing"
)

func BenchmarkUuidv4(b *testing.B) {
	for i := 0; i < b.N; i++ {
		uuid.NewRandom()
		// uuid.New()
	}
}
func BenchmarkUuidv5(b *testing.B) {
	I := 0
	for i := 0; i < b.N; i++ {
		uuid.New()
		I++

		// uuid.New() 1429945
	}
	fmt.Println("I::::", I)
}
