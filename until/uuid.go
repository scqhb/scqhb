package until

import (
	"github.com/google/uuid"
)

// he Uuidv4 is google uuid
func Uuidv4() string {
	u4 := uuid.New()
	return u4.String()
}
