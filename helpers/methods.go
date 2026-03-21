package helpers

import (
	"SDT_ApiServices/common"
	"log"
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
	log.Println("Err_Response", Err_Response)
	return Err_Response
}

// func Error_Response(c *gin.Context, ErrCode string, ErrDescription string, pErr error) {
// 	c.JSON(http.StatusBadRequest, GetErrorString(ErrCode, ErrDescription, pErr))
// 	return
// }
// func success_Response(c *gin.Context, Msg string) {
// 	c.JSON(http.StatusOK, GetSuccessString(Msg))
// 	return
// }

// func InternalServerError(c *gin.Context, ErrCode string, ErrDescription string, pErr error) {
// 	c.JSON(http.StatusInternalServerError, GetErrorString(ErrCode, ErrDescription, pErr))
// 	return
// }
