package nats

import "testing"

func TestRegister(t *testing.T) {
	if err := Register(); err != nil {
		t.Fatalf("first register failed: %v", err)
	}
	if err := Register(); err != nil {
		t.Fatalf("second register should be idempotent: %v", err)
	}
}
