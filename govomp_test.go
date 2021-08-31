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

	dat := []float32{1, 2, 3}

	for _, device := range devices {
		buf, err := device.NewArrayBuffer(dat)
		if err != nil {
			t.Fatal(err)
		}

		fmt.Println(device.Name, device.Type)
		fmt.Println(buf.Read())
	}
}
