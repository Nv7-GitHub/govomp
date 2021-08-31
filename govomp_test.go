package govomp

// Remember to do -gcflags=-G=3!

import (
	"fmt"
	"testing"
)

func TestSquares(t *testing.T) {
	err := Init()
	if err != nil {
		t.Fatal(err)
	}

	devices, err := GetDevices()
	if err != nil {
		t.Fatal(err)
	}

	for _, device := range devices {
		fmt.Println(device.Name, device.Type)
	}
}
