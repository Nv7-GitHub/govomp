package govomp

import vk "github.com/vulkan-go/vulkan"

var instance vk.Instance

// Init initializes vulkan and the vulkan instance.
func Init() error {
	err := vk.SetDefaultGetInstanceProcAddr()
	if err != nil {
		return err
	}

	err = vk.Init()
	if err != nil {
		return err
	}

	instanceCreateInfo := &vk.InstanceCreateInfo{
		SType: vk.StructureTypeInstanceCreateInfo,
		PApplicationInfo: &vk.ApplicationInfo{
			SType:              vk.StructureTypeApplicationInfo,
			ApiVersion:         vk.MakeVersion(1, 0, 0),
			ApplicationVersion: vk.MakeVersion(1, 0, 0),
		},
	}

	err = vk.Error(vk.CreateInstance(instanceCreateInfo, nil, &instance))
	if err != nil {
		return err
	}

	err = vk.InitInstance(instance)
	if err != nil {
		return err
	}

	return nil
}
