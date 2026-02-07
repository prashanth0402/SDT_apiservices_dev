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

type CreateGrpReqStruct struct {
	GroupName string `json:"groupname"`
	FTCcode   string `json:"ftccode"`
}

type GetGroupName struct {
	Groupinfo []JoinGrpReqStruct `json:"groupinfo"`
	Status    string             `json:"status"`
}

// Create a group in ChatRoom and add record in chatroom_membership.
func CreateGroup(w http.ResponseWriter, r *http.Request) {

	lDebug := new(helpers.HelperStruct)
	lDebug.SetUid(r)
	lDebug.Log(helpers.Statement, "CreateGroup (+)")
	(w).Header().Set("Access-Control-Allow-Origin", common.EKYCAllowedOrigin)
	(w).Header().Set("Access-Control-Allow-Credentials", "true")
	(w).Header().Set("Access-Control-Allow-Methods", "PUT, OPTIONS")
	(w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

	switch r.Method {
	case "PUT":

		// lFTCcode, lErr := appsession.KycReadCookie(r, lDebug, common.CHATName)

		// if lErr != nil {
		// 	lDebug.Log(helpers.Elog, "CG001"+lErr.Error())
		// 	fmt.Fprint(w, helpers.GetError_String("Log001", "Something went wrong. Please try agin later."))
		// 	return
		// }

		lDb, lErr := ftdb.LocalDbConnect(ftdb.MariaFTPRD)

		if lErr != nil {
			lDebug.Log(helpers.Elog, "CG002"+lErr.Error())
			fmt.Fprint(w, helpers.GetError_String("CG002", "Something went wrong. Please try agin later."))
			return
		}
		defer lDb.Close()

		lErr = HandleCreateGrp(lDebug, lDb, r)
		fmt.Println("dfgfg", lErr)

		if lErr != nil {
			fmt.Println("dfgfg", lErr)
			if lErr.Error() == "exist" {
				lDebug.Log(helpers.Elog, "CG003"+lErr.Error())
				fmt.Fprint(w, helpers.GetError_String("CG003", "Exist"))
				return
			}
			lDebug.Log(helpers.Elog, "CG003"+lErr.Error())
			fmt.Fprint(w, helpers.GetError_String("CG003", "Invalid Credential."))
			return
		}
		fmt.Fprint(w, helpers.GetMsg_String("S", "Valid Credential."))

		lDebug.Log(helpers.Statement, "CreateGroup (-)")

	}

}

func HandleCreateGrp(pdebug *helpers.HelperStruct, pDb *sql.DB, r *http.Request) error {
	pdebug.Log(helpers.Statement, "HandleCreateGrp (+)")
	var lUserData CreateGrpReqStruct
	var lChatID int
	// lCreatedBy := "FTC"

	lByte, lErr := ioutil.ReadAll(r.Body)
	if lErr != nil {
		pdebug.Log(helpers.Elog, "HCG001"+lErr.Error())
		return fmt.Errorf("HCG001 : %v", lErr)
	}

	lErr = json.Unmarshal(lByte, &lUserData)
	if lErr != nil {
		pdebug.Log(helpers.Elog, "HCG002"+lErr.Error())
		return fmt.Errorf("HCG002 : %v", lErr)
	}

	// fmt.Println("lUserData ", lUserData)

	lExistingUser, lErr := CheckExistingGrp(pdebug, pDb, lUserData.GroupName)
	if lErr != nil {
		pdebug.Log(helpers.Elog, "HCG003"+lErr.Error())
		return fmt.Errorf("HCG003 : %v", lErr)
	}

	if lExistingUser == "Y" {
		pdebug.Log(helpers.Elog, "HCG003 : Group Already Exist...")
		// fmt.Println("lExistingUser ", lExistingUser)
		return fmt.Errorf("exist")
	}

	// fmt.Println(" after lExistingUser ", lExistingUser)

	lCoreString := `insert into chatroom (roomname,roomcreatedby,createdby,createddate,updatedby,updateddate,isactive)
values ( ?,?,?,now(),?,now(),"Y" )`

	lRowId, lErr := pDb.Exec(lCoreString, lUserData.GroupName, lUserData.FTCcode, lUserData.FTCcode, lUserData.FTCcode)
	if lErr != nil {
		pdebug.Log(helpers.Elog, "HCG004"+lErr.Error())
		return fmt.Errorf("HCG004 : %v", lErr)
	}

	// fmt.Println(" after lCoreString ", lExistingUser)

	lTempChatID, lErr := lRowId.LastInsertId()
	if lErr != nil {
		pdebug.Log(helpers.Elog, "HCG005"+lErr.Error())
		return fmt.Errorf("HCG005 : %v", lErr)
	}

	lChatID = int(lTempChatID)

	lErr = JoinGroup(pdebug, pDb, lUserData.FTCcode, lChatID)
	if lErr != nil {
		pdebug.Log(helpers.Elog, "HCG006"+lErr.Error())
		return fmt.Errorf("HCG006 : %v", lErr)
	}

	pdebug.Log(helpers.Statement, "HandleCreateGrp (-)")
	return nil
}

func CheckExistingGrp(pDebug *helpers.HelperStruct, pDb *sql.DB, pGrpName string) (string, error) {
	pDebug.Log(helpers.Statement, "CheckExistingGrp (+)")
	var lGroupExtFlag string

	lCorestring := `SELECT CASE WHEN count(*) > 0 THEN 'Y' ELSE 'N' END AS Flag
	FROM chatroom
	WHERE roomname  = ?`

	rows, lErr := pDb.Query(lCorestring, pGrpName)
	if lErr != nil {
		return "", helpers.ErrReturn(lErr)
	}

	for rows.Next() {
		lErr := rows.Scan(&lGroupExtFlag)
		if lErr != nil {
			return "", helpers.ErrReturn(lErr)
		}
	}
	pDebug.Log(helpers.Statement, "CheckExistingGrp (-)")
	return lGroupExtFlag, nil
}

func GetActiveGroup(w http.ResponseWriter, r *http.Request) {
	lDebug := new(helpers.HelperStruct)
	lDebug.SetUid(r)
	lDebug.Log(helpers.Statement, "GetActiveGroup (+)")
	(w).Header().Set("Access-Control-Allow-Origin", "*")
	(w).Header().Set("Access-Control-Allow-Credentials", "true")
	(w).Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	(w).Header().Set("Access-Control-Allow-Headers", "FTcode,Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

	switch r.Method {
	case "GET":

		var lResp GetGroupName

		lResp.Status = common.SuccessCode

		lFTCcode := r.Header.Get("FTcode")

		// lFTCcode, lErr := appsession.KycReadCookie(r, lDebug, common.CHATName)

		// if lErr != nil {
		// 	lDebug.Log(helpers.Elog, "CG001"+lErr.Error())
		// 	fmt.Fprint(w, helpers.GetError_String("Log001", "Something went wrong. Please try again later."))
		// 	return
		// }

		lDb, lErr := ftdb.LocalDbConnect(ftdb.MariaFTPRD)

		if lErr != nil {
			lDebug.Log(helpers.Elog, "CG002"+lErr.Error())
			fmt.Fprint(w, helpers.GetError_String("CG002", "Something went wrong. Please try again later."))
			return
		}
		defer lDb.Close()

		lResp.Groupinfo, lErr = FetchGrp(lDebug, lDb, lFTCcode)
		if lErr != nil {
			lDebug.Log(helpers.Elog, "CG003"+lErr.Error())
			fmt.Fprint(w, helpers.GetError_String("CG003", "Invalid Credential."))
			return
		}
		lData, lErr := json.Marshal(lResp)
		if lErr != nil {
			lDebug.Log(helpers.Elog, "CG004"+lErr.Error())
			fmt.Fprint(w, helpers.GetError_String("CG004", "Something went wrong. Please try again later."))
			return
		}
		fmt.Fprint(w, string(lData))
		lDebug.Log(helpers.Statement, "GetActiveGroup (-)")

	}
}

func FetchGrp(pDebug *helpers.HelperStruct, pDb *sql.DB, pFtcode string) ([]JoinGrpReqStruct, error) {
	pDebug.Log(helpers.Statement, "FetchGrp (+)")

	var lGroupDet JoinGrpReqStruct
	var lAllGrp []JoinGrpReqStruct

	lCorestring := ` select c.roomname,c.id,c2.id 
	from chatroom c , chatroommembership c2 
	where c.id = c2.chatroomid and c.isactive = 'Y' and c2.isActive = 'Y' and c2.user= ? `

	rows, lErr := pDb.Query(lCorestring, pFtcode)
	if lErr != nil {
		return lAllGrp, helpers.ErrReturn(lErr)
	}

	for rows.Next() {
		lErr := rows.Scan(&lGroupDet.GroupName, &lGroupDet.Groupid, &lGroupDet.MemberID)
		if lErr != nil {
			return lAllGrp, helpers.ErrReturn(lErr)
		}
		lAllGrp = append(lAllGrp, lGroupDet)
	}
	pDebug.Log(helpers.Statement, "FetchGrp (-)")
	return lAllGrp, nil
}

func GetAllGroup(w http.ResponseWriter, r *http.Request) {
	lDebug := new(helpers.HelperStruct)
	lDebug.SetUid(r)
	lDebug.Log(helpers.Statement, "GetAllGroup (+)")
	(w).Header().Set("Access-Control-Allow-Origin", "*")
	(w).Header().Set("Access-Control-Allow-Credentials", "true")
	(w).Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	(w).Header().Set("Access-Control-Allow-Headers", "FTcode,Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

	switch r.Method {
	case "GET":

		var lResp GetGroupName

		lResp.Status = common.SuccessCode

		lFTCcode := r.Header.Get("FTcode")

		lDb, lErr := ftdb.LocalDbConnect(ftdb.MariaFTPRD)

		if lErr != nil {
			lDebug.Log(helpers.Elog, "CG002"+lErr.Error())
			fmt.Fprint(w, helpers.GetError_String("CG002", "Something went wrong. Please try again later."))
			return
		}
		defer lDb.Close()

		lResp.Groupinfo, lErr = FetchAllGrp(lDebug, lDb, lFTCcode)
		if lErr != nil {
			lDebug.Log(helpers.Elog, "CG003"+lErr.Error())
			fmt.Fprint(w, helpers.GetError_String("CG003", "Something went wrong. Please try again later."))
			return
		}
		lData, lErr := json.Marshal(lResp)
		if lErr != nil {
			lDebug.Log(helpers.Elog, "CG004"+lErr.Error())
			fmt.Fprint(w, helpers.GetError_String("CG004", "Something went wrong. Please try again later."))
			return
		}
		fmt.Fprint(w, string(lData))
		lDebug.Log(helpers.Statement, "GetAllGroup (-)")

	}
}

func FetchAllGrp(pDebug *helpers.HelperStruct, pDb *sql.DB, pFtcode string) ([]JoinGrpReqStruct, error) {
	pDebug.Log(helpers.Statement, "FetchAllGrp (+)")

	var lGroupDet JoinGrpReqStruct
	var lAllGrp []JoinGrpReqStruct

	lCorestring := `select c.roomname,c.id,u.username 
	from chatroom c,userconfig u 
	where c.roomcreatedby = u.clientid and c.isactive = 'Y'`

	rows, lErr := pDb.Query(lCorestring)
	if lErr != nil {
		return lAllGrp, helpers.ErrReturn(lErr)
	}

	for rows.Next() {
		lErr := rows.Scan(&lGroupDet.GroupName, &lGroupDet.Groupid, &lGroupDet.MemberName)
		if lErr != nil {
			return lAllGrp, helpers.ErrReturn(lErr)
		}
		lAllGrp = append(lAllGrp, lGroupDet)
	}
	pDebug.Log(helpers.Statement, "FetchAllGrp (-)")
	return lAllGrp, nil
}

func DeleteGroup() {
	//
}
