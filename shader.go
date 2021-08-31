package main

import (
	_ "embed"
	"unsafe"

	vk "github.com/vulkan-go/vulkan"
)

//go:embed shader.spv
var shader []byte

func createShader(shader []byte, device vk.Device) vk.ShaderModule {
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
	err := vk.Error(vk.CreateShaderModule(vk.Device(device), &info, nil, &out))
	handle(err)

	return out
}
