package binarypack

import "testing"

func TestUnpack(t *testing.T) {
	raw := []byte("")
	Unpack(raw, nil)
}
