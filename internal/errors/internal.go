package apierrors

import "fmt"

type MissingGinContextParam struct {
	paramName string
}

func NewMissingGinContextParamError(paramName string) MissingGinContextParam {
	return MissingGinContextParam{paramName: paramName}
}

func (m MissingGinContextParam) Error() string {
	return fmt.Sprintf("context parameter %s was not found.", m.paramName)
}

type InvalidGinContextParamType struct {
	paramName string
}

func NewInvalidContextParamTypeError(paramName string) InvalidGinContextParamType {
	return NewInvalidContextParamTypeError(paramName)
}

func (i InvalidGinContextParamType) Error() string {
	return fmt.Sprintf("context parameter %s has invalid type.", i.paramName)
}
