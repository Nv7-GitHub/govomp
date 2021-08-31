package main

import (
	"fmt"
	"unsafe"

	vk "github.com/vulkan-go/vulkan"
)

func getBuffer(physicalDevice vk.PhysicalDevice, device vk.Device, data []float32) vk.Buffer {
	bufSize := len(data) * int(unsafe.Sizeof(data[0]))

	// LEFT OF HERE
	var properties vk.PhysicalDeviceMemoryProperties
	vk.GetPhysicalDeviceMemoryProperties(physicalDevice, &properties)
	// the properties' c ref has all the values but the Go object doesnt, can access C vals (kind of) with reflect
	// since the properties Go object isn't filled in, it doesn't work

	mem := allocMemory(bufSize, properties, device)
	fmt.Println(mem)

	return nil
}

func allocMemory(size int, props vk.PhysicalDeviceMemoryProperties, device vk.Device) vk.DeviceMemory {
	memTypeIndex := uint32(vk.MaxMemoryTypes)

	fmt.Println(props.MemoryTypeCount)

	for i := uint32(0); i < props.MemoryTypeCount; i++ {
		fmt.Println(vk.MemoryPropertyFlags(vk.MemoryPropertyHostVisibleBit)&props.MemoryTypes[i].PropertyFlags != 0)
		if (vk.MemoryPropertyFlags(vk.MemoryPropertyHostVisibleBit)&props.MemoryTypes[i].PropertyFlags != 0) &&
			(vk.MemoryPropertyFlags(vk.MemoryPropertyHostCoherentBit)&props.MemoryTypes[i].PropertyFlags != 0) &&
			(vk.MemoryPropertyFlags(vk.MemoryPropertyHostCachedBit)&props.MemoryTypes[i].PropertyFlags != 0) && // For read performance
			(vk.DeviceSize(size) < props.MemoryHeaps[props.MemoryTypes[i].HeapIndex].Size) {
			memTypeIndex = i
			break
		}
	}

	if memTypeIndex == vk.MaxMemoryTypes {
		panic("govomp: no suitable memory type found")
	}

	memAllocInfo := &vk.MemoryAllocateInfo{
		SType:           vk.StructureTypeMemoryAllocateInfo,
		AllocationSize:  vk.DeviceSize(size),
		MemoryTypeIndex: memTypeIndex,
	}

	var memory vk.DeviceMemory
	err := vk.Error(vk.AllocateMemory(device, memAllocInfo, nil, &memory))
	handle(err)

	return memory
}
