package chatdata

import (
	"database/sql"
	"encoding/json"
	"fcs23pkg/common"
	"fcs23pkg/ftdb"
	"fcs23pkg/helpers"
	"fmt"
	"io/ioutil"
	"net/http"
)

type JoinGrpReqStruct struct {
	GroupName  string `json:"groupname"`
	Groupid    int    `json:"groupid"`
	MemberID   int    `json:"memberid,omitempty"`
	MemberName string `json:"membername,omitempty"`
	FtCode     string `json:"ftcode,omitempty"`
}

func JoinGrpReq(w http.ResponseWriter, r *http.Request) {
	lDebug := new(helpers.HelperStruct)
	lDebug.SetUid(r)
	lDebug.Log(helpers.Statement, "JoinGrpReq (+)", common.EKYCAllowedOrigin)
	(w).Header().Set("Access-Control-Allow-Origin", common.EKYCAllowedOrigin)
	(w).Header().Set("Access-Control-Allow-Credentials", "true")
	(w).Header().Set("Access-Control-Allow-Methods", "PUT, OPTIONS")
	(w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

	switch r.Method {

	case "PUT":
		lDb, lErr := ftdb.LocalDbConnect(ftdb.MariaFTPRD)

		if lErr != nil {
			lDebug.Log(helpers.Elog, "JG002"+lErr.Error())
			fmt.Fprint(w, helpers.GetError_String("JG002", "Something went wrong. Please try agin later."))
			return
		}
		defer lDb.Close()

		lErr = HandelJoinGroup(lDebug, lDb, r)
		if lErr != nil {
			lDebug.Log(helpers.Elog, "JG003"+lErr.Error())
			fmt.Fprint(w, helpers.GetError_String("JG003", lErr.Error()))
			return
		}

		fmt.Fprint(w, helpers.GetMsg_String("S", "Joined Successfully."))

		lDebug.Log(helpers.Statement, "JoinGrpReq (-)")

	}
}

func HandelJoinGroup(pDebug *helpers.HelperStruct, pDb *sql.DB, r *http.Request) error {
	pDebug.Log(helpers.Statement, "HandelJoinGroup (+)")

	var lJoinGroupReq JoinGrpReqStruct
	lBody, lErr := ioutil.ReadAll(r.Body)
	if lErr != nil {
		pDebug.Log(helpers.Elog, "GNU001"+lErr.Error())
		return fmt.Errorf("GNU001 : %v", lErr)
	}

	lErr = json.Unmarshal(lBody, &lJoinGroupReq)
	if lErr != nil {
		pDebug.Log(helpers.Elog, "GNU002"+lErr.Error())
		return fmt.Errorf("GNU002 : %v", lErr)
	}

	lUserFlag, lErr := CheckUserinGrp(pDebug, pDb, lJoinGroupReq.FtCode, lJoinGroupReq.Groupid)
	if lErr != nil || lUserFlag == "Y" {
		if lErr == nil {
			lErr = fmt.Errorf("already member")
		}
		pDebug.Log(helpers.Elog, "GNU003 "+lErr.Error(), lUserFlag)
		return fmt.Errorf("%v", lErr)
	}

	lErr = JoinGroup(pDebug, pDb, lJoinGroupReq.FtCode, lJoinGroupReq.Groupid)
	if lErr != nil {
		pDebug.Log(helpers.Elog, "GNU004 "+lErr.Error(), lUserFlag)
		return fmt.Errorf("GNU004 : %v", lErr)
	}

	pDebug.Log(helpers.Statement, "HandelJoinGroup (-)")
	return nil

}

func JoinGroup(pDebug *helpers.HelperStruct, pDb *sql.DB, pUser string, pGroupid int) error {
	pDebug.Log(helpers.Statement, "JoinGroup (+)")

	lCorestring := `insert into chatroommembership (user,chatroomid,isActive,createdby,createddate,updatedby,updateddate)
    values ( ?,?,'Y',?,now(),?,now())`

	_, lErr := pDb.Exec(lCorestring, pUser, pGroupid, pUser, pUser)
	if lErr != nil {
		pDebug.Log(helpers.Elog, "JG001"+lErr.Error())
		return fmt.Errorf("JG001 : %v", lErr)
	}

	pDebug.Log(helpers.Statement, "JoinGroup (-)")
	return nil
}

func CheckUserinGrp(pDebug *helpers.HelperStruct, pDb *sql.DB, pUser string, pGroupid int) (string, error) {
	pDebug.Log(helpers.Statement, "CheckUserinGrp (+)")
	var lUserFlag string

	lCorestring := `SELECT CASE WHEN count(*) > 0 THEN 'Y' ELSE 'N' END AS Flag
		FROM chatroommembership
		WHERE user = ? and chatroomid=?`

	rows, lErr := pDb.Query(lCorestring, pUser, pGroupid)
	if lErr != nil {
		return "", helpers.ErrReturn(lErr)
	}

	for rows.Next() {
		lErr := rows.Scan(&lUserFlag)
		if lErr != nil {
			return "", helpers.ErrReturn(lErr)
		}
	}

	pDebug.Log(helpers.Statement, "CheckUserinGrp (-)")
	return lUserFlag, nil
}

func LeaveGroup() {
	//
}
