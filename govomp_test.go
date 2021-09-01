package govomp

import (
	_ "embed"
	"fmt"
	"testing"
)

//go:embed shader.spv
var shader []byte

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
		fmt.Println(device.Name)

		// Buffers
		buf, err := device.NewArrayBuffer(dat)
		if err != nil {
			t.Fatal(err)
		}

		out, err := device.AllocateArrayBuffer(len(dat))
		if err != nil {
			t.Fatal(err)
		}

		uniform, err := device.NewArrayBuffer([]float32{float32(len(dat))})
		if err != nil {
			t.Fatal(err)
		}

		// Create & Run Shader
		shader, err := device.NewShader(shader)
		if err != nil {
			t.Fatal(err)
		}

		err = device.RunShader(shader, len(dat), 1, uniform, buf, out)
		if err != nil {
			t.Fatal(err)
		}

		// Read output
		outDat, err := out.Read()
		if err != nil {
			t.Fatal(err)
		}
		fmt.Println(outDat)
	}
}
