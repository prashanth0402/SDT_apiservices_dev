package chatdata

import (
	"SDT_apiservices_dev/common"
	"database/sql"
	"fcs23pkg/helpers"
	"fmt"
	"strconv"
)

type MsgStruct struct {
	MsgId           string `json:"msgid"`
	Senderid        string `json:"senderid"`
	SenderName      string `json:"sendername"`
	SenderInitial   string `json:"senderinitial"`
	Receiverid      string `json:"receiverid"`
	ReceiverName    string `json:"receivername"`
	ReceiverInitial string `json:"receiverinitial"`
	Time            string `json:"time"`
	Message         string `json:"msg"`
	Attachmentid    string `json:"attachmentid"`
	AttachmentName  string `json:"filename"`
	AttachmentType  string `json:"filetype"`
}

type SocketResp struct {
	MessageArr []MsgStruct `json:"messagearr"`
	StatusCode string      `json:"statuscode"`
	GroupId    string      `json:"groupid"`
}

func ChatHistory(pDebug *helpers.HelperStruct, pDb *sql.DB, pGroupId int) (SocketResp, error) {
	// pDebug.Log(helpers.Statement, "ChatHistory (+)")

	var lMsgDet MsgStruct
	var lSocketResp SocketResp

	lSocketResp.GroupId = strconv.Itoa(pGroupId)
	lSocketResp.StatusCode = common.SuccessCode

	lCorestring := `
	SELECT DISTINCT 
	c.id,
    NVL(c.senderid, '') AS senderid,    
    NVL(u_sender.username, '') AS sender_name,
	SUBSTR(NVL(u_sender.username, ''), 1, 1) AS sender_initial,
    NVL(c.msgcontent, '') AS msgcontent,
    NVL(c.recieverid, '') AS recieverid,   
    NVL(u_receiver.username, '') AS receiver_name,
	SUBSTR(NVL(u_receiver.username, ''), 1, 1) AS receiver_initial,
    NVL(c.attachmentdocid, '') AS attachmentdocid,
    NVL(TO_CHAR(c.createdtime, 'HH24:MI'), '') AS createdtime,
    NVL(attachmentname , '') AS filename,
    NVL(attachmenttype , '') AS filetype
FROM 
    chathistory c
JOIN 
    userconfig u_sender
ON 
    c.senderid = u_sender.clientid
left JOIN 
    userconfig u_receiver
ON 
    c.recieverid = u_receiver.clientid
JOIN 
    chatroom c2 
ON 
    c.chatroomid = c2.id 
JOIN 
    chatroommembership c3 
ON 
    c3.chatroomid = c2.id
WHERE  
    c3.isActive = 'Y'  
    AND c.isdeleted = 'N'
    AND c.chatroomid = ?;
	`

	lRows, lErr := pDb.Query(lCorestring, pGroupId)

	if lErr != nil {
		pDebug.Log(helpers.Elog, "CH001"+lErr.Error())
		lSocketResp.StatusCode = common.ErrorCode
		return lSocketResp, fmt.Errorf("CH001 : %v", lErr)
	}

	for lRows.Next() {
		lErr = lRows.Scan(&lMsgDet.MsgId, &lMsgDet.Senderid, &lMsgDet.SenderName, &lMsgDet.SenderInitial, &lMsgDet.Message, &lMsgDet.Receiverid, &lMsgDet.ReceiverName, &lMsgDet.ReceiverInitial, &lMsgDet.Attachmentid, &lMsgDet.Time, &lMsgDet.AttachmentName, &lMsgDet.AttachmentType)
		if lErr != nil {
			pDebug.Log(helpers.Elog, "CH002 :"+lErr.Error())
			lSocketResp.StatusCode = common.ErrorCode
			return lSocketResp, fmt.Errorf("CH002 : %v", lErr)
		} else {
			lSocketResp.MessageArr = append(lSocketResp.MessageArr, lMsgDet)
		}
	}

	// pDebug.Log(helpers.Statement, "ChatHistory (-)")
	return lSocketResp, nil
}

// func GetChatHistory(pDebug *helpers.HelperStruct, pDb *sql.DB){
// 	pDebug.Log(helpers.Statement, "GetChatHistory (+)")

// 	pDebug.Log(helpers.Statement, "GetChatHistory (+)")

// }

func StoreMsg(pDebug *helpers.HelperStruct, pDb *sql.DB, pMsg Message) error {
	pDebug.Log(helpers.Statement, "StoreMsg (+)")

	lCorestring := `insert into chathistory (senderid,recieverid,msgcontent,attachmentdocid,createdtime,isdeleted,createdby,createddate,updatedby,updateddate,chatroomid,attachmentname,attachmenttype)
  values(?,?,?,?,now(),'N',?,now(),?,now(),?,?,?)`

	_, lErr := pDb.Exec(lCorestring, pMsg.SenderId, pMsg.Receiverid, pMsg.Msg, pMsg.Attachmentid, pMsg.SenderId, pMsg.SenderId, pMsg.ChatRoomId, pMsg.AttachmentName, pMsg.AttachmentType)

	if lErr != nil {
		pDebug.Log(helpers.Elog, "SM001"+lErr.Error())
		return fmt.Errorf("SM001 : %v", lErr)
	}

	pDebug.Log(helpers.Statement, "StoreMsg (-)")
	return nil
}

func DeleteMsg(pDebug *helpers.HelperStruct, pDb *sql.DB, pMsgId int) error {
	pDebug.Log(helpers.Statement, "StoreMsg (+)")

	lCorestring := `  update chathistory set isdeleted = 'Y' where id = ? `

	_, lErr := pDb.Exec(lCorestring, pMsgId)

	if lErr != nil {
		return fmt.Errorf("SM001 : %v", lErr)
	}

	pDebug.Log(helpers.Statement, "StoreMsg (-)")
	return nil
}
