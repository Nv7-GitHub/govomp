package main

import (
	"fmt"

	vk "github.com/vulkan-go/vulkan"
)

func runCompute(physicalDevice vk.PhysicalDevice, instance vk.Instance) {
	// Get queue family
	var queueFamilyPropertyCount uint32
	vk.GetPhysicalDeviceQueueFamilyProperties(physicalDevice, &queueFamilyPropertyCount, nil)

	queueFamilyProperties := make([]vk.QueueFamilyProperties, queueFamilyPropertyCount)
	vk.GetPhysicalDeviceQueueFamilyProperties(physicalDevice, &queueFamilyPropertyCount, queueFamilyProperties)

	// Create Device Queue
	queueCreateInfo := vk.DeviceQueueCreateInfo{
		SType:            vk.StructureTypeDeviceQueueCreateInfo,
		QueueCount:       1,
		PQueuePriorities: []float32{1.0},
	}
	deviceCreateInfo := vk.DeviceCreateInfo{
		SType:                vk.StructureTypeDeviceCreateInfo,
		QueueCreateInfoCount: 1,
		PQueueCreateInfos:    []vk.DeviceQueueCreateInfo{queueCreateInfo},
	}

	var device vk.Device
	err := vk.Error(vk.CreateDevice(physicalDevice, &deviceCreateInfo, nil, &device))
	handle(err)

	var queue vk.Queue
	vk.GetDeviceQueue(device, 0, 0, &queue)

	// Get Buffer
	mem := getBuffer(physicalDevice, device, []float32{1, 2, 3})
	fmt.Println(mem)
}
