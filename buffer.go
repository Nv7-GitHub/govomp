package govomp

import (
	"unsafe"

	vk "github.com/vulkan-go/vulkan"
)

type Buffer struct {
	mem    *Memory
	buf    vk.Buffer
	device vk.Device
}

// Memory represents memory space on a device
type Memory struct {
	device vk.Device
	Memory vk.DeviceMemory
	Size   int
}

// AllocMemory allocates memory on the device
func (d *Device) AllocMemory(size int) (*Memory, error) {
	memTypeIndex := uint32(vk.MaxMemoryTypes)

	for i := uint32(0); i < d.memProperties.MemoryTypeCount; i++ {
		memFlags := d.memProperties.MemoryTypes[i].PropertyFlags
		if (vk.MemoryPropertyFlags(vk.MemoryPropertyHostVisibleBit)&memFlags != 0) &&
			(vk.MemoryPropertyFlags(vk.MemoryPropertyHostCoherentBit)&memFlags != 0) &&
			//(vk.MemoryPropertyFlags(vk.MemoryPropertyHostCachedBit)&memFlags.PropertyFlags != 0) && // For read performance, doesn't affect performance much tho
			(vk.DeviceSize(size) < d.memProperties.MemoryHeaps[d.memProperties.MemoryTypes[i].HeapIndex].Size) {
			memTypeIndex = i
			break
		}
	}

	if memTypeIndex == vk.MaxMemoryTypes {
		return nil, ErrNoSuitableMemoryTypeFound
	}

	memAllocInfo := &vk.MemoryAllocateInfo{
		SType:           vk.StructureTypeMemoryAllocateInfo,
		AllocationSize:  vk.DeviceSize(size),
		MemoryTypeIndex: memTypeIndex,
	}

	var memory vk.DeviceMemory
	err := vk.Error(vk.AllocateMemory(d.device, memAllocInfo, nil, &memory))
	if err != nil {
		return nil, err
	}
	return &Memory{
		device: d.device,
		Memory: memory,
		Size:   size,
	}, nil
}

// CreateVulkanBuffer creates a vulkan buffer object from memory and binds it to the device
func (m *Memory) GetVulkanBuffer() (vk.Buffer, error) {
	var buffer vk.Buffer
	err := vk.Error(vk.CreateBuffer(m.device, &vk.BufferCreateInfo{
		SType:                 vk.StructureTypeBufferCreateInfo,
		Size:                  vk.DeviceSize(m.Size),
		Usage:                 vk.BufferUsageFlags(vk.BufferUsageStorageBufferBit),
		SharingMode:           vk.SharingModeExclusive,
		QueueFamilyIndexCount: 1,
	}, nil, &buffer))
	return buffer, err
}

// TODO: Use generics for the array (same as below) comment
// WriteArray writes an array to the memory
func (m *Memory) WriteArray(data []float32) error {
	bufSize := len(data) * int(unsafe.Sizeof(data[0]))

	var ptr unsafe.Pointer
	err := vk.Error(vk.MapMemory(m.device, m.Memory, 0, vk.DeviceSize(bufSize), 0, &ptr))
	if err != nil {
		return err
	}

	byteArr := unsafe.Slice((*byte)(unsafe.Pointer(&data[0])), len(data)*int(unsafe.Sizeof(data[0])))
	n := vk.Memcopy(ptr, byteArr)
	if n != len(byteArr) {
		return ErrFailedToCopyMemory
	}

	vk.UnmapMemory(m.device, m.Memory)

	return nil
}

// TODO: Use generics for the data, and support float32 and int32 arrays [waiting for go 1.18 release]
// NewBuffer creates a buffer on the target device with the provided data
func (d *Device) NewBufferArray(data []float32) (*Buffer, error) {
	bufSize := len(data) * int(unsafe.Sizeof(data[0]))

	mem, err := d.AllocMemory(bufSize)
	if err != nil {
		return nil, err
	}

	err = mem.WriteArray(data)
	if err != nil {
		return nil, err
	}

	buf, err := mem.GetVulkanBuffer()
	if err != nil {
		return nil, err
	}

	return &Buffer{
		mem:    mem,
		device: d.device,
		buf:    buf,
	}, nil
}
