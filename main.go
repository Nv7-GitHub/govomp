package main

import (
	vk "github.com/vulkan-go/vulkan"
)

var appInfo = &vk.ApplicationInfo{
	SType:              vk.StructureTypeApplicationInfo,
	ApiVersion:         vk.MakeVersion(1, 0, 0),
	ApplicationVersion: vk.MakeVersion(1, 0, 0),
}

func handle(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	// Get Instance
	err := vk.SetDefaultGetInstanceProcAddr()
	handle(err)

	err = vk.Init()
	handle(err)

	var instance vk.Instance
	instanceCreateInfo := &vk.InstanceCreateInfo{
		SType:            vk.StructureTypeInstanceCreateInfo,
		PApplicationInfo: appInfo,
	}
	err = vk.Error(vk.CreateInstance(instanceCreateInfo, nil, &instance))
	handle(err)
	err = vk.InitInstance(instance)
	handle(err)

	// Get Devices
	var deviceCount uint32
	err = vk.Error(vk.EnumeratePhysicalDevices(instance, &deviceCount, nil))
	handle(err)

	devices := make([]vk.PhysicalDevice, deviceCount)
	err = vk.Error(vk.EnumeratePhysicalDevices(instance, &deviceCount, devices))
	handle(err)

	for _, device := range devices {
		runCompute(device, instance)
	}
}
