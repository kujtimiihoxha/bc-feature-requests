package models

// Predefined model error codes.
const (
	ErrDatabase = -1
	ErrSystem   = -2
	ErrUnAuthorized = -3
	ErrNotFound = -4
	ErrPasswordMismatch = -5
	ErrEmailExists   = -6
	ErrUsernameExists   = -7
)

// CodeInfo.
// The code info structure.
// Code: int code value (lower than 0 if error).
// Info: additional information.
type CodeInfo struct {
	Code int    `json:"code"`
	Info string `json:"info"`
}

// NewErrorInfo return a CodeInfo represents error.
func ErrorInfo(code int, info string) *CodeInfo {
	return &CodeInfo{code, info}
}

// NewNormalInfo return a CodeInfo represents OK.
func OkInfo(info string) *CodeInfo {
	return &CodeInfo{0, info}
}

