package inspect

import "fmt"
import "github.com/go-gl/gl/v4.6-core/gl"
import "github.com/dmarychev/gazebo/core"

type iUniformVariable struct {
	Name     string
	Location uint32
}

func (uvi iUniformVariable) String() string {
	return fmt.Sprintf("%v(location=%v) uniform variable", uvi.Name, uvi.Location)
}

func uniformVariables(t *core.Technique) ([]iUniformVariable, error) {

	var numUniforms int32
	gl.GetProgramInterfaceiv(uint32(*t), gl.UNIFORM, gl.ACTIVE_RESOURCES, &numUniforms)
	core.CheckError()

	uniformSet := make([]iUniformVariable, 0, numUniforms)
	if cap(uniformSet) == 0 {
		return uniformSet, nil
	}

	for uniformIndex := uint32(0); uniformIndex < uint32(numUniforms); uniformIndex++ {

		// retrieve name length
		var nameLen int32
		nameLenProp := uint32(gl.NAME_LENGTH)
		gl.GetProgramResourceiv(uint32(*t), gl.UNIFORM, uint32(uniformIndex), 1, &nameLenProp, 1, nil, &nameLen)
		core.CheckError()

		// retrieve name
		name := make([]uint8, nameLen)
		gl.GetProgramResourceName(uint32(*t), gl.UNIFORM, uint32(uniformIndex), int32(len(name)), &nameLen, &name[0])
		name = name[:nameLen]
		core.CheckError()

		// location
		var location int32
		locationProp := uint32(gl.LOCATION)
		gl.GetProgramResourceiv(uint32(*t), gl.UNIFORM, uint32(uniformIndex), 1, &locationProp, 1, nil, &location)
		core.CheckError()

		uniformSet = append(uniformSet, iUniformVariable{
			Name:     string(name),
			Location: uint32(location),
		})
	}

	return uniformSet, nil
}

type iBufferVariable struct {
	Name   string
	Index  uint32
	Offset uint32
}

func (bvi iBufferVariable) String() string {
	return fmt.Sprintf("%v(offset=%v) buffer variable", bvi.Name, bvi.Offset)
}

func bufferVariable(t *core.Technique, varIndex uint32) (*iBufferVariable, error) {

	// retrieve name length
	var nameLen int32
	nameLenProp := uint32(gl.NAME_LENGTH)
	gl.GetProgramResourceiv(uint32(*t), gl.BUFFER_VARIABLE, uint32(varIndex), 1, &nameLenProp, 1, nil, &nameLen)
	core.CheckError()

	// retrieve name
	name := make([]uint8, nameLen)
	gl.GetProgramResourceName(uint32(*t), gl.BUFFER_VARIABLE, uint32(varIndex), int32(len(name)), &nameLen, &name[0])
	name = name[:nameLen]
	core.CheckError()

	// retrieve offset
	var offset int32
	offsetProp := uint32(gl.OFFSET)
	gl.GetProgramResourceiv(uint32(*t), gl.BUFFER_VARIABLE, uint32(varIndex), 1, &offsetProp, 1, nil, &offset)
	core.CheckError()

	return &iBufferVariable{
		Index:  varIndex,
		Name:   string(name),
		Offset: uint32(offset),
	}, nil
}

type iShaderStorageBuffer struct {
	Name      string
	Binding   uint32
	Variables []iBufferVariable
}

func (ssbi iShaderStorageBuffer) String() string {
	return fmt.Sprintf("%v(binding=%v variables=%v) ssbo", ssbi.Name, ssbi.Binding, len(ssbi.Variables))
}

func shaderStorageBuffers(t *core.Technique) ([]iShaderStorageBuffer, error) {
	var numSsb int32
	gl.GetProgramInterfaceiv(uint32(*t), gl.SHADER_STORAGE_BLOCK, gl.ACTIVE_RESOURCES, &numSsb)
	core.CheckError()

	ssbiSet := make([]iShaderStorageBuffer, 0, numSsb)
	if cap(ssbiSet) == 0 {
		return ssbiSet, nil
	}

	for ssbIndex := uint32(0); ssbIndex < uint32(numSsb); ssbIndex++ {

		// retrieve name length
		var nameLen int32
		nameLenProp := uint32(gl.NAME_LENGTH)
		gl.GetProgramResourceiv(uint32(*t), gl.SHADER_STORAGE_BLOCK, ssbIndex, 1, &nameLenProp, 1, nil, &nameLen)

		// retrieve name
		name := make([]uint8, nameLen)
		gl.GetProgramResourceName(uint32(*t), gl.SHADER_STORAGE_BLOCK, ssbIndex, int32(len(name)), &nameLen, &name[0])
		name = name[:nameLen]

		// retrieve binding
		var binding int32
		bindingProp := uint32(gl.BUFFER_BINDING)
		gl.GetProgramResourceiv(uint32(*t), gl.SHADER_STORAGE_BLOCK, ssbIndex, 1, &bindingProp, 1, nil, &binding)

		// retrieve number of variables
		var numVariables int32
		numVariablesProp := uint32(gl.NUM_ACTIVE_VARIABLES)
		gl.GetProgramResourceiv(uint32(*t), gl.SHADER_STORAGE_BLOCK, ssbIndex, 1, &numVariablesProp, 1, nil, &numVariables)

		// retrieve variable indices
		varIndices := make([]int32, numVariables)
		varIndicesProp := uint32(gl.ACTIVE_VARIABLES)
		gl.GetProgramResourceiv(uint32(*t), gl.SHADER_STORAGE_BLOCK, ssbIndex, 1, &varIndicesProp, numVariables, nil, &varIndices[0])

		variableInfos := make([]iBufferVariable, 0, numVariables)
		for _, varIndex := range varIndices {

			varInfo, err := bufferVariable(t, uint32(varIndex))
			if err != nil {
				return nil, err
			}

			variableInfos = append(variableInfos, *varInfo)
		}

		ssbiSet = append(ssbiSet, iShaderStorageBuffer{
			Name:      string(name),
			Binding:   uint32(binding),
			Variables: variableInfos,
		})
	}

	return ssbiSet, nil
}

type iTechnique struct {
	UniformVariables     []iUniformVariable
	ShaderStorageBuffers []iShaderStorageBuffer
}

func InspectTechnique(t *core.Technique) (*iTechnique, error) {
	ssboInfoSet, err := shaderStorageBuffers(t)
	if err != nil {
		return nil, err
	}

	uniformInfoSet, err := uniformVariables(t)
	if err != nil {
		return nil, err
	}

	return &iTechnique{
		UniformVariables:     uniformInfoSet,
		ShaderStorageBuffers: ssboInfoSet,
	}, nil
}
