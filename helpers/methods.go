package helpers

import (
	"SDT_ApiServices/common"
	"encoding/json"
	"log"
)

type Error_Response struct {
	Status string `json:"status"`
	ErrMsg string `json:"errMsg"`
}

type Msg_Response struct {
	Status string `json:"status"`
	ErrMsg string `json:"errMsg"`
}

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

func GetErrorString(ErrCode string, ErrDescription string, pErr error) string {
	log.Println(ErrCode, pErr)
	var Err_Response Error_Response
	Err_Response.Status = common.ErrorCode
	Err_Response.ErrMsg = ErrCode + "/" + ErrDescription

	lResult, err := json.Marshal(Err_Response)

	if err != nil {
		log.Println(err)
	}

	return string(lResult)

}
