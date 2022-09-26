package test

import (
	"fmt"
	uuid "github.com/satori/go.uuid"
	"testing"
)

func TestGenerateUUID(t *testing.T) {
	s := uuid.NewV4().String()
	fmt.Printf("s: %v\n", s)
}
