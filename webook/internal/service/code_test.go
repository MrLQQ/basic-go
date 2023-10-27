package service

import (
	"fmt"
	"math/rand"
	"testing"
)

func TestCodeGenerate(t *testing.T) {
	code := rand.Intn(1000000)
	t.Log(fmt.Sprintf("%06d", code))
}
