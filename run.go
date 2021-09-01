package govomp

import (
	"math"

	vk "github.com/vulkan-go/vulkan"
)

// RunShader runs a shader with the given buffers
func (d *Device) RunShader(shader *Shader, threadCount int, workGroupSize int, args ...Buffer) error {
	// Create Descriptor Layout
	bindings := make([]vk.DescriptorSetLayoutBinding, len(args))
	for i := uint32(0); i < uint32(len(args)); i++ {
		bindings[i] = vk.DescriptorSetLayoutBinding{
			Binding:         i,
			DescriptorType:  vk.DescriptorTypeUniformBuffer,
			DescriptorCount: 1,
			StageFlags:      vk.ShaderStageFlags(vk.ShaderStageComputeBit),
		}
	}
	layoutCreateInfo := vk.DescriptorSetLayoutCreateInfo{
		SType:        vk.StructureTypeDescriptorSetLayoutCreateInfo,
		BindingCount: uint32(len(bindings)),
		PBindings:    bindings,
	}
	var descriptorSetLayout vk.DescriptorSetLayout
	err := vk.Error(vk.CreateDescriptorSetLayout(d.device, &layoutCreateInfo, nil, &descriptorSetLayout))
	if err != nil {
		return err
	}

	// Create Pipeline
	var pipelineLayout vk.PipelineLayout
	pipelineLayoutCreateInfo := vk.PipelineLayoutCreateInfo{
		SType:          vk.StructureTypePipelineLayoutCreateInfo,
		SetLayoutCount: 1,
		PSetLayouts:    []vk.DescriptorSetLayout{descriptorSetLayout},
	}
	err = vk.Error(vk.CreatePipelineLayout(d.device, &pipelineLayoutCreateInfo, nil, &pipelineLayout))
	if err != nil {
		return err
	}

	pipelineCreateInfo := vk.ComputePipelineCreateInfo{
		SType: vk.StructureTypeComputePipelineCreateInfo,
		Stage: vk.PipelineShaderStageCreateInfo{
			SType:  vk.StructureTypePipelineShaderStageCreateInfo,
			Flags:  vk.PipelineShaderStageCreateFlags(vk.ShaderStageComputeBit),
			Module: shader.shader,
			PName:  "main\x00", // Runs "main" function
		},
		Layout: pipelineLayout,
	}

	pipelines := make([]vk.Pipeline, 1)
	err = vk.Error(vk.CreateComputePipelines(d.device, vk.PipelineCache(vk.NullHandle), 1, []vk.ComputePipelineCreateInfo{pipelineCreateInfo}, nil, pipelines))
	if err != nil {
		return err
	}

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
	err = vk.Error(vk.CreateDescriptorPool(d.device, &descriptorPoolCreateInfo, nil, &descriptorPool))
	if err != nil {
		return err
	}

	// Create Descriptor Set
	descriptorSetAllocateInfo := vk.DescriptorSetAllocateInfo{
		SType:              vk.StructureTypeDescriptorSetAllocateInfo,
		DescriptorPool:     vk.DescriptorPool(descriptorPool),
		DescriptorSetCount: 1,
		PSetLayouts:        []vk.DescriptorSetLayout{descriptorSetLayout},
	}

	var descriptorSet vk.DescriptorSet
	err = vk.Error(vk.AllocateDescriptorSets(d.device, &descriptorSetAllocateInfo, &descriptorSet))
	if err != nil {
		return err
	}

	// Get Buffers Ready
	writeDescriptorSet := make([]vk.WriteDescriptorSet, len(args))
	for i := range writeDescriptorSet {
		writeDescriptorSet[i] = vk.WriteDescriptorSet{
			SType:           vk.StructureTypeWriteDescriptorSet,
			DstSet:          vk.DescriptorSet(descriptorSet),
			DstBinding:      uint32(i),
			DescriptorCount: 1,
			DescriptorType:  vk.DescriptorTypeStorageBuffer,
			PBufferInfo: []vk.DescriptorBufferInfo{{
				Buffer: args[i].getBuffer(),
				Range:  vk.DeviceSize(vk.WholeSize),
			}},
		}
	}
	vk.UpdateDescriptorSets(d.device, uint32(len(writeDescriptorSet)), writeDescriptorSet, 0, nil)

	// Create Command Pool
	commandPoolCreateInfo := vk.CommandPoolCreateInfo{
		SType: vk.StructureTypeCommandPoolCreateInfo,
	}

	var commandPool vk.CommandPool
	err = vk.Error(vk.CreateCommandPool(d.device, &commandPoolCreateInfo, nil, &commandPool))
	if err != nil {
		return err
	}

	// Create Command Buffer
	commandBufferAllocateInfo := vk.CommandBufferAllocateInfo{
		SType:              vk.StructureTypeCommandBufferAllocateInfo,
		CommandPool:        commandPool,
		Level:              vk.CommandBufferLevelPrimary,
		CommandBufferCount: 1,
	}

	commandBuffers := make([]vk.CommandBuffer, 1)
	err = vk.Error(vk.AllocateCommandBuffers(d.device, &commandBufferAllocateInfo, commandBuffers))
	if err != nil {
		return err
	}

	// Add commands to command buffer
	err = vk.Error(vk.BeginCommandBuffer(commandBuffers[0], &vk.CommandBufferBeginInfo{
		SType: vk.StructureTypeCommandBufferBeginInfo,
		Flags: vk.CommandBufferUsageFlags(vk.CommandBufferUsageOneTimeSubmitBit),
	}))
	if err != nil {
		return err
	}

	vk.CmdBindPipeline(commandBuffers[0], vk.PipelineBindPointCompute, pipelines[0])
	vk.CmdBindDescriptorSets(commandBuffers[0], vk.PipelineBindPointCompute, pipelineLayout, 0, 1, []vk.DescriptorSet{descriptorSet}, 0, nil)
	vk.CmdDispatch(commandBuffers[0], uint32(math.Ceil(float64(threadCount)/float64(workGroupSize))), 1, 1)

	err = vk.Error(vk.EndCommandBuffer(commandBuffers[0]))
	if err != nil {
		return err
	}

	// Setup Device Queue
	var queue vk.Queue
	vk.GetDeviceQueue(d.device, 0, 0, &queue)

	// Run Commands in Command buffer
	submitInfo := vk.SubmitInfo{
		SType:              vk.StructureTypeSubmitInfo,
		CommandBufferCount: 1,
		PCommandBuffers:    commandBuffers,
	}

	err = vk.Error(vk.QueueSubmit(queue, 1, []vk.SubmitInfo{submitInfo}, vk.NullFence))
	if err != nil {
		return err
	}

	// Wait for it to finish
	return vk.Error(vk.QueueWaitIdle(queue))
}
