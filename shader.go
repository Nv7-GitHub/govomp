package govomp

import (
	"unsafe"

	vk "github.com/vulkan-go/vulkan"
)

type Shader struct {
	shader vk.ShaderModule
	device vk.Device
}

// NewShader accepts SPIR-V shader code and returns a shader object.
func (d *Device) NewShader(shader []byte) (*Shader, error) {
	buf := make([]uint32, len(shader)/4)
	vk.Memcopy(unsafe.Pointer((*struct {
		Data uintptr
		Len  int
		Cap  int
	})(unsafe.Pointer(&buf)).Data), shader)

	var out vk.ShaderModule
	info := vk.ShaderModuleCreateInfo{
		SType:    vk.StructureTypeShaderModuleCreateInfo,
		CodeSize: uint(len(shader)),
		PCode:    buf,
	}

	err := vk.Error(vk.CreateShaderModule(d.device, &info, nil, &out))
	if err != nil {
		return nil, err
	}

	return &Shader{
		shader: out,
		device: d.device,
	}, nil
}
