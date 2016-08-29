package controllers

// Controller Error.
// Error response structure.
// Status: tells the error status ex. 404, 500 etc.
// Code: tells the error Code ex. 1000, 1001 etc. this is application specific codes.
// Message: error message.
// DevInfo: if in development mode will send additional information needed from the developers.
// MoreInfo: additional information for the user.
type ControllerError struct {
	Status   int    `json:"status"`
	Code     int    `json:"code"`
	Message  string `json:"message"`
	DevInfo  string `json:"dev_info,omitempty"`
	MoreInfo string `json:"more_info"`
}

// ErrorController.
// ErrorController handlers 404 error and returns errors from other controllers.
type ErrorController struct {
	BaseController
}
// Predefined const error strings.
const (
	ErrInputData    = "Data entry errors"
	ErrDatabase     = "Database operation errors"
	ErrDupUser      = "User information already exists"
	ErrNoUser       = "User information does not exist"
	ErrPass         = "Incorrect password"
	ErrNoUserPass   = "User information does not exist or the password is incorrect"
	ErrNoUserChange = "User information does not exist or has not changed the data"
	ErrInvalidUser  = "User information is incorrect"
	ErrOpenFile     = "Open File Error"
	ErrWriteFile    = "Error writing file"
	ErrSystem       = "System error"
	ErrInputDataValidation      = "Error validating data"
)
// Predefined error messages
var (
	err404          = &ControllerError{404, 404, "Not Found", "", "Api endpoint not found, please see documentations at:[TODO]"}
	errInputData    = &ControllerError{400, 10001, ErrInputData, "Client parameter error", ""}
	errDatabase     = &ControllerError{500, 10002, ErrDatabase, ErrDatabase, ""}
	errDupUser      = &ControllerError{400, 10003, ErrDupUser, "Duplicate database records", ""}
	errNoUser       = &ControllerError{400, 10004, ErrNoUser, "Database record does not exist", ""}
	errPass         = &ControllerError{400, 10005, ErrNoUserPass, ErrPass, ""}
	errNoUserPass   = &ControllerError{400, 10006, ErrNoUserPass, "Database records do not exist or the password is incorrect", ""}
	errNoUserChange = &ControllerError{400, 10007, ErrNoUserChange, "Database records do not exist or data are not changed", ""}
	errInvalidUser  = &ControllerError{400, 10008, ErrInvalidUser, "Session information is incorrect", ""}
	errOpenFile     = &ControllerError{500, 10009, "Server Error", ErrOpenFile, ""}
	errWriteFile    = &ControllerError{500, 10010, "Server Error", ErrWriteFile, ""}
	errSystem       = &ControllerError{500, 10011, ErrSystem, ErrSystem, ""}
	errExpired      = &ControllerError{400, 10012, "Login expired", "Verification token expired", ""}
	errPermission   = &ControllerError{400, 10013, "Permission denied", "You are not authorized", ""}
	errInputDataValidation   = &ControllerError{400, 10014, ErrInputDataValidation, ErrInputDataValidation, ""}
	errDataNotFound          = &ControllerError{404, 404, "Data not found", "No data found", ""}
)

// Error404
// Default error handler for Not Found error.
func (c *ErrorController) Error404() {
	c.RetError(err404)
}
