package main

import "C"

import (
	"fmt"
	"math"
	"unsafe"

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

	// Get Buffer
	dat := []float32{1, 2, 3}
	inp := getBuffer(physicalDevice, device, dat)
	out, outMem := allocBuffer(int(unsafe.Sizeof(dat[0]))*len(dat), device, physicalDevice)
	uniform := createUniformData(physicalDevice, device, len(dat))

	// Create Shader
	shader := createShader(shader, device)

	// Create Bindings
	bindings := []vk.DescriptorSetLayoutBinding{
		{
			Binding:         1,
			DescriptorType:  vk.DescriptorTypeStorageBuffer,
			DescriptorCount: 1,
			StageFlags:      vk.ShaderStageFlags(vk.ShaderStageComputeBit),
		},
		{
			Binding:         2,
			DescriptorType:  vk.DescriptorTypeStorageBuffer,
			DescriptorCount: 1,
			StageFlags:      vk.ShaderStageFlags(vk.ShaderStageComputeBit),
		},
		{
			Binding:         0,
			DescriptorType:  vk.DescriptorTypeUniformBuffer,
			DescriptorCount: 1,
			StageFlags:      vk.ShaderStageFlags(vk.ShaderStageComputeBit),
		},
	}

	layoutCreateInfo := vk.DescriptorSetLayoutCreateInfo{
		SType:        vk.StructureTypeDescriptorSetLayoutCreateInfo,
		BindingCount: uint32(len(bindings)),
		PBindings:    bindings,
	}

	var descriptorSetLayout vk.DescriptorSetLayout
	err = vk.Error(vk.CreateDescriptorSetLayout(device, &layoutCreateInfo, nil, &descriptorSetLayout))
	handle(err)

	// Create Pipeline
	var pipelineLayout vk.PipelineLayout
	pipelineLayoutCreateInfo := vk.PipelineLayoutCreateInfo{
		SType:          vk.StructureTypePipelineLayoutCreateInfo,
		SetLayoutCount: 1,
		PSetLayouts:    []vk.DescriptorSetLayout{descriptorSetLayout},
	}
	err = vk.Error(vk.CreatePipelineLayout(device, &pipelineLayoutCreateInfo, nil, &pipelineLayout))
	handle(err)

	pipelineCreateInfo := vk.ComputePipelineCreateInfo{
		SType: vk.StructureTypeComputePipelineCreateInfo,
		Stage: vk.PipelineShaderStageCreateInfo{
			SType:  vk.StructureTypePipelineShaderStageCreateInfo,
			Flags:  vk.PipelineShaderStageCreateFlags(vk.ShaderStageComputeBit),
			Module: shader,
			PName:  "main\x00",
		},
		Layout: pipelineLayout,
	}

	pipelines := make([]vk.Pipeline, 1)
	err = vk.Error(vk.CreateComputePipelines(device, vk.PipelineCache(vk.NullHandle), 1, []vk.ComputePipelineCreateInfo{pipelineCreateInfo}, nil, pipelines))
	handle(err)

	// Create Descriptor Pool
	descriptorPoolSize := vk.DescriptorPoolSize{
		Type:            vk.DescriptorTypeStorageBuffer,
		DescriptorCount: 2,
	}
	descriptorPoolCreateInfo := vk.DescriptorPoolCreateInfo{
		SType:         vk.StructureTypeDescriptorPoolCreateInfo,
		MaxSets:       1,
		PoolSizeCount: 1,
		PPoolSizes:    []vk.DescriptorPoolSize{descriptorPoolSize},
	}

	var descriptorPool vk.DescriptorPool
	err = vk.Error(vk.CreateDescriptorPool(device, &descriptorPoolCreateInfo, nil, &descriptorPool))
	handle(err)

	// Create Descriptor Set
	descriptorSetAllocateInfo := vk.DescriptorSetAllocateInfo{
		SType:              vk.StructureTypeDescriptorSetAllocateInfo,
		DescriptorPool:     vk.DescriptorPool(descriptorPool),
		DescriptorSetCount: 1,
		PSetLayouts:        []vk.DescriptorSetLayout{descriptorSetLayout},
	}

	var descriptorSet vk.DescriptorSet
	err = vk.Error(vk.AllocateDescriptorSets(device, &descriptorSetAllocateInfo, &descriptorSet))
	handle(err)

	// Get Buffers Ready
	inputBufferInfo := vk.DescriptorBufferInfo{
		Buffer: inp,
		Range:  vk.DeviceSize(vk.WholeSize),
	}
	outBufferInfo := vk.DescriptorBufferInfo{
		Buffer: out,
		Range:  vk.DeviceSize(vk.WholeSize),
	}
	uniformBufferInfo := vk.DescriptorBufferInfo{
		Buffer: uniform,
		Range:  vk.DeviceSize(vk.WholeSize),
	}
	writeDescriptorSet := []vk.WriteDescriptorSet{
		{
			SType:           vk.StructureTypeWriteDescriptorSet,
			DstSet:          vk.DescriptorSet(descriptorSet),
			DstBinding:      1,
			DescriptorCount: 1,
			DescriptorType:  vk.DescriptorTypeStorageBuffer,
			PBufferInfo:     []vk.DescriptorBufferInfo{inputBufferInfo},
		},
		{
			SType:           vk.StructureTypeWriteDescriptorSet,
			DstSet:          vk.DescriptorSet(descriptorSet),
			DstBinding:      2,
			DescriptorCount: 1,
			DescriptorType:  vk.DescriptorTypeStorageBuffer,
			PBufferInfo:     []vk.DescriptorBufferInfo{outBufferInfo},
		},
		{
			SType:           vk.StructureTypeWriteDescriptorSet,
			DstSet:          vk.DescriptorSet(descriptorSet),
			DstBinding:      0,
			DescriptorCount: 1,
			DescriptorType:  vk.DescriptorTypeStorageBuffer,
			PBufferInfo:     []vk.DescriptorBufferInfo{uniformBufferInfo},
		},
	}
	vk.UpdateDescriptorSets(device, uint32(len(writeDescriptorSet)), writeDescriptorSet, 0, nil)

	// Create Command Pool
	commandPoolCreateInfo := vk.CommandPoolCreateInfo{
		SType: vk.StructureTypeCommandPoolCreateInfo,
	}

	var commandPool vk.CommandPool
	err = vk.Error(vk.CreateCommandPool(device, &commandPoolCreateInfo, nil, &commandPool))
	handle(err)

	// Create Command Buffer
	commandBufferAllocateInfo := vk.CommandBufferAllocateInfo{
		SType:              vk.StructureTypeCommandBufferAllocateInfo,
		CommandPool:        commandPool,
		Level:              vk.CommandBufferLevelPrimary,
		CommandBufferCount: 1,
	}

	commandBuffers := make([]vk.CommandBuffer, 1)
	err = vk.Error(vk.AllocateCommandBuffers(device, &commandBufferAllocateInfo, commandBuffers))
	handle(err)

	// Create Command Buffer
	err = vk.Error(vk.BeginCommandBuffer(commandBuffers[0], &vk.CommandBufferBeginInfo{
		SType: vk.StructureTypeCommandBufferBeginInfo,
		Flags: vk.CommandBufferUsageFlags(vk.CommandBufferUsageOneTimeSubmitBit),
	}))
	handle(err)

	vk.CmdBindPipeline(commandBuffers[0], vk.PipelineBindPointCompute, pipelines[0])
	vk.CmdBindDescriptorSets(commandBuffers[0], vk.PipelineBindPointCompute, pipelineLayout, 0, 1, []vk.DescriptorSet{descriptorSet}, 0, nil)
	workGroupSize := 1
	vk.CmdDispatch(commandBuffers[0], uint32(math.Ceil(float64(len(dat))/float64(workGroupSize))), 1, 1)

	err = vk.Error(vk.EndCommandBuffer(commandBuffers[0]))
	handle(err)

	// Get Device Queue
	var queue vk.Queue
	vk.GetDeviceQueue(device, 0, 0, &queue)

	// Submit Queue & Wait
	submitInfo := vk.SubmitInfo{
		SType:              vk.StructureTypeSubmitInfo,
		CommandBufferCount: 1,
		PCommandBuffers:    commandBuffers,
	}

	err = vk.Error(vk.QueueSubmit(queue, 1, []vk.SubmitInfo{submitInfo}, vk.NullFence))
	handle(err)

	err = vk.Error(vk.QueueWaitIdle(queue))
	handle(err)

	// Read Data Back
	var data unsafe.Pointer
	err = vk.Error(vk.MapMemory(device, outMem, 0, vk.DeviceSize(vk.WholeSize), 0, &data))
	handle(err)

	outData := unsafe.Slice((*float32)(data), len(dat))

	vk.UnmapMemory(device, outMem)

	fmt.Println(outData)
}
