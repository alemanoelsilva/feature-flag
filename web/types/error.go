package types

type ErrorCreateFeatureFlagForm struct {
	HasError           bool
	IsNameError        bool
	IsDescriptionError bool
	IsRequestError     bool
	ErrorMessage       string
}
