package helpers

import (
	"SDT_ApiServices/common"
)

// Msg_Response represents a standard JSON success response.
// Example: {"status":"SUCCESS", "errMsg":"Operation completed"}
type Msg_Response struct {
	Status string `json:"status"`
	ErrMsg string `json:"errMsg"`
}

// GetSuccessString returns a JSON string for success responses.
// It takes a message and wraps it in Msg_Response format.
// Example: GetSuccessString("User created") → {"status":"SUCCESS","errMsg":"User created"}
func GetSuccessString(Msg string) Msg_Response {

	var Msg_Res Msg_Response

	Msg_Res.Status = common.SUCCESSCODE
	Msg_Res.ErrMsg = Msg

	return Msg_Res

}

// GetErrorString returns a JSON string for error responses.
// It takes an error code, description, and the original error for logging.
// Example: GetErrorString("400", "Invalid Email", err)
// → {"status":"ERROR","errMsg":"400/Invalid Email"}
func GetErrorString(ErrCode string, ErrDescription string, pErr error) Msg_Response {
	var Err_Response Msg_Response
	Err_Response.Status = common.ErrorCode
	Err_Response.ErrMsg = ErrCode + "/" + ErrDescription
	return Err_Response
}
