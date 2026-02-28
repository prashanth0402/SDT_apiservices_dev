package helpers

import (
	"SDT_ApiServices/common"
	"encoding/json"
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
func GetSuccessString(Msg string) string {

	var Msg_Res Msg_Response

	Msg_Res.Status = common.SUCCESSCODE
	Msg_Res.ErrMsg = Msg

	result, err := json.Marshal(Msg_Res)

	if err != nil {
		log.Println(err)
	}

	return string(result)

}

// GetErrorString returns a JSON string for error responses.
// It takes an error code, description, and the original error for logging.
// Example: GetErrorString("400", "Invalid Email", err)
// → {"status":"ERROR","errMsg":"400/Invalid Email"}
func GetErrorString(ErrCode string, ErrDescription string, pErr error) string {
	log.Println(ErrCode, pErr)
	var Err_Response Msg_Response
	Err_Response.Status = common.ErrorCode
	Err_Response.ErrMsg = ErrCode + "/" + ErrDescription

	lResult, err := json.Marshal(Err_Response)

	if err != nil {
		log.Println(err)
	}

	return string(lResult)

}

// CheckError returns true if error exists and prints error code + message
func CheckError(err error) bool {
	return err != nil
}
