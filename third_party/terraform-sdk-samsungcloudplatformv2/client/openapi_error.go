package scpsdk

// GenericOpenAPIError Provides access to the body, error and model on returned errors.
type GenericOpenAPIError struct {
	ResponseBody []byte
	ErrorMessage string
	ErrorModel   interface{}
}

// Error returns non-empty string if there was an error.
func (e GenericOpenAPIError) Error() string {
	return e.ErrorMessage
}

// Body returns the raw bytes of the response
func (e GenericOpenAPIError) Body() []byte {
	return e.ResponseBody
}

// Model returns the unpacked model of the error
func (e GenericOpenAPIError) Model() interface{} {
	return e.ErrorModel
}
