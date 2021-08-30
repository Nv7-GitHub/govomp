package main

import (
	"fmt"

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

	// Get Physical Device
	var deviceCount uint32
	err = vk.Error(vk.EnumeratePhysicalDevices(instance, &deviceCount, nil))
	handle(err)
	fmt.Println(deviceCount)
}
