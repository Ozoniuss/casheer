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

type InvalidModel struct {
	modelName       string
	validationError error
}

func NewInvalidModelError(resource string, validationError error) InvalidModel {
	return InvalidModel{
		modelName:       resource,
		validationError: validationError,
	}
}

func (r InvalidModel) Error() string {
	return fmt.Sprintf("Invalid %s resource: %s", r.modelName, r.validationError.Error())
}

type InvalidQueryParams struct {
	orig error
}

func NewInvalidQueryParamsError(orig error) InvalidQueryParams {
	return InvalidQueryParams{
		orig: orig,
	}
}

func (q InvalidQueryParams) Error() string {
	return q.orig.Error()
}

func (q InvalidQueryParams) Unwrap() error {
	return q.orig
}
