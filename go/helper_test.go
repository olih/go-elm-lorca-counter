package main

import "testing"

// TestIncrement should add the increment to the number for positive increment
func TestIncrement(t *testing.T) {
	total := Increment(10, 2)
	if total != 12 {
		t.Errorf("Increment was incorrect, got: %d, want: %d.", total, 12)
	}
}
