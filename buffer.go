package main

/*
typedef struct UniformData {
	int totalCount;
} UniformData;

void* getPtr(UniformData* value) {
	return (void*)value;
}

size_t getSize(UniformData value) {
	return sizeof(value);
}
*/
import "C"

import (
	"unsafe"

	vk "github.com/vulkan-go/vulkan"
)

func getBuffer(physicalDevice vk.PhysicalDevice, device vk.Device, data []float32) vk.Buffer {
	bufSize := len(data) * int(unsafe.Sizeof(data[0]))

	var properties vk.PhysicalDeviceMemoryProperties
	vk.GetPhysicalDeviceMemoryProperties(physicalDevice, &properties)
	properties.Deref()

	// Alloc
	mem := allocMemory(bufSize, properties, device)

	// Map and transfer
	var payload unsafe.Pointer
	err := vk.Error(vk.MapMemory(device, mem, 0, vk.DeviceSize(bufSize), 0, &payload))
	handle(err)

	byteArr := unsafe.Slice((*byte)(unsafe.Pointer(&data[0])), len(data)*int(unsafe.Sizeof(data[0])))
	n := vk.Memcopy(payload, byteArr)
	if n != len(byteArr) {
		panic("govomp: failed to copy memory")
	}

	vk.UnmapMemory(device, mem)

	// Create buffer
	return createBuffer(mem, device, bufSize)
}

func createBuffer(mem vk.DeviceMemory, device vk.Device, size int) vk.Buffer {
	var buffer vk.Buffer
	err := vk.Error(vk.CreateBuffer(device, &vk.BufferCreateInfo{
		SType:                 vk.StructureTypeBufferCreateInfo,
		Size:                  vk.DeviceSize(size),
		Usage:                 vk.BufferUsageFlags(vk.BufferUsageStorageBufferBit),
		SharingMode:           vk.SharingModeExclusive,
		QueueFamilyIndexCount: 1,
	}, nil, &buffer))
	handle(err)

	err = vk.Error(vk.BindBufferMemory(device, buffer, mem, 0))
	handle(err)

	return buffer
}

func allocMemory(size int, props vk.PhysicalDeviceMemoryProperties, device vk.Device) vk.DeviceMemory {
	memTypeIndex := uint32(vk.MaxMemoryTypes)

	for i := uint32(0); i < props.MemoryTypeCount; i++ {
		props.MemoryTypes[i].Deref()
		props.MemoryHeaps[i].Deref()
		if (vk.MemoryPropertyFlags(vk.MemoryPropertyHostVisibleBit)&props.MemoryTypes[i].PropertyFlags != 0) &&
			(vk.MemoryPropertyFlags(vk.MemoryPropertyHostCoherentBit)&props.MemoryTypes[i].PropertyFlags != 0) &&
			//(vk.MemoryPropertyFlags(vk.MemoryPropertyHostCachedBit)&props.MemoryTypes[i].PropertyFlags != 0) && // For read performance
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

func allocBuffer(size int, device vk.Device, physicalDevice vk.PhysicalDevice) (vk.Buffer, vk.DeviceMemory) {
	var properties vk.PhysicalDeviceMemoryProperties
	vk.GetPhysicalDeviceMemoryProperties(physicalDevice, &properties)
	properties.Deref()

	mem := allocMemory(size, properties, device)
	buf := createBuffer(mem, device, size)
	return buf, mem
}

func createUniformData(physicalDevice vk.PhysicalDevice, device vk.Device, size int) vk.Buffer {
	data := C.UniformData{totalCount: (C.int)(size)}
	bufSize := C.getSize(data)

	var properties vk.PhysicalDeviceMemoryProperties
	vk.GetPhysicalDeviceMemoryProperties(physicalDevice, &properties)
	properties.Deref()

	// Alloc
	mem := allocMemory(int(bufSize), properties, device)

	// Map and transfer
	var payload unsafe.Pointer
	err := vk.Error(vk.MapMemory(device, mem, 0, vk.DeviceSize(bufSize), 0, &payload))
	handle(err)

	ptr := C.getPtr(&data)
	byteArray := unsafe.Slice((*byte)(ptr), int(bufSize))
	n := vk.Memcopy(payload, byteArray)
	if n != len(byteArray) {
		panic("govomp: failed to copy memory")
	}

	vk.UnmapMemory(device, mem)

	// Create buffer
	return createBuffer(mem, device, int(bufSize))
}
