package govomp

import (
	vk "github.com/vulkan-go/vulkan"
)

// Device represents a Vulkan device
type Device struct {
	Name   string
	Type   vk.PhysicalDeviceType
	Limits vk.PhysicalDeviceLimits

	physicalDevice vk.PhysicalDevice
	device         vk.Device

	properties    vk.PhysicalDeviceProperties
	memProperties vk.PhysicalDeviceMemoryProperties
}

// init initializes the device from the physicalDevice
func (d *Device) init() error {
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
	err := vk.Error(vk.CreateDevice(d.physicalDevice, &deviceCreateInfo, nil, &device))
	if err != nil {
		return err
	}
	d.device = device

	var props vk.PhysicalDeviceProperties
	vk.GetPhysicalDeviceProperties(d.physicalDevice, &props)
	d.properties = props

	var memProps vk.PhysicalDeviceMemoryProperties
	vk.GetPhysicalDeviceMemoryProperties(d.physicalDevice, &memProps)
	d.memProperties = memProps

	return nil
}

func (d *Device) deref() {
	d.properties.Deref()
	d.properties.Limits.Deref()

	d.memProperties.Deref()
	for i := range d.memProperties.MemoryTypes {
		d.memProperties.MemoryTypes[i].Deref()
	}
	for i := range d.memProperties.MemoryHeaps {
		d.memProperties.MemoryHeaps[i].Deref()
	}

	d.Name = string(d.properties.DeviceName[:])
	d.Type = d.properties.DeviceType
	d.Limits = d.properties.Limits
}

// Free frees the device's memory
func (d *Device) Free() {
	d.properties.Limits.Free()
	d.properties.Free()
	vk.DestroyDevice(d.device, nil)
}

// GetDevices gets the available devices
func GetDevices() ([]*Device, error) {
	var deviceCount uint32
	err := vk.Error(vk.EnumeratePhysicalDevices(instance, &deviceCount, nil))
	if err != nil {
		return nil, err
	}

	devices := make([]vk.PhysicalDevice, deviceCount)
	err = vk.Error(vk.EnumeratePhysicalDevices(instance, &deviceCount, devices))
	if err != nil {
		return nil, err
	}

	out := make([]*Device, len(devices))
	for i, physicalDevice := range devices {
		out[i].physicalDevice = physicalDevice
		err = out[i].init()
		if err != nil {
			return nil, err
		}
		out[i].deref()
	}
	return out, nil
}
