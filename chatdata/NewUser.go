package chatdata

import (
	"database/sql"
	"encoding/json"
	"fcs23pkg/ftdb"
	"fcs23pkg/helpers"
	"fmt"
	"io/ioutil"
	"net/http"
)

type NewUserStruct struct {
	UserName string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
	ClientId string `json:"clientid"`
}

// type ResponseStruct struct{
// 	Status
// }

// Create User fro chat
func Signup(w http.ResponseWriter, r *http.Request) {
	// fmt.Println("asdfgasdf")
	lDebug := new(helpers.HelperStruct)
	lDebug.SetUid(r)
	lDebug.Log(helpers.Statement, "Signup (+)")
	(w).Header().Set("Access-Control-Allow-Origin", "*")
	(w).Header().Set("Access-Control-Allow-Credentials", "false")
	(w).Header().Set("Access-Control-Allow-Methods", "PUT,OPTIONS")
	(w).Header().Set("Access-Control-Allow-Headers", "Accept,Content-Type,Content-Length,Accept-Encoding,X-CSRF-Token,Authorization")
	(w).Header().Set("Content-Type", "application/json")

	// if r.Method == http.MethodOptions {
	// 	w.WriteHeader(http.StatusNoContent)
	// 	fmt.Println("Handled OPTIONS request")
	// 	return
	// }

	switch r.Method {
	case http.MethodPut:

		lDb, lErr := ftdb.LocalDbConnect(ftdb.MariaFTPRD)

		if lErr != nil {
			lDebug.Log(helpers.Elog, "S001"+lErr.Error())
			fmt.Fprint(w, helpers.GetError_String("S001", "Something went wrong. Please try agin later."))
			return
		}
		defer lDb.Close()

		lErr = HandleNewUser(lDebug, lDb, r)

		if lErr != nil {
			lDebug.Log(helpers.Elog, "S002"+lErr.Error())

			if lErr.Error() == "olduser" {
				fmt.Fprint(w, helpers.GetError_String("S002", "Exist"))
				return

			}
			fmt.Fprint(w, helpers.GetError_String("S002", "Something went wrong. Please try agin later."))
			return
		}

		fmt.Fprint(w, helpers.GetMsg_String("S", "Account Created Successfully."))

		lDebug.Log(helpers.Statement, "Signup (-)")
	}
}

func HandleNewUser(pdebug *helpers.HelperStruct, pDb *sql.DB, r *http.Request) error {
	pdebug.Log(helpers.Statement, "HandleNewUser (+)")
	var lUserData NewUserStruct
	lCreatedBy := "FTC"

	lByte, lErr := ioutil.ReadAll(r.Body)
	if lErr != nil {
		pdebug.Log(helpers.Elog, "GNU001"+lErr.Error())
		return fmt.Errorf("GNU001 : %v", lErr)
	}

	lErr = json.Unmarshal(lByte, &lUserData)
	if lErr != nil {
		pdebug.Log(helpers.Elog, "GNU002"+lErr.Error())
		return fmt.Errorf("GNU002 : %v", lErr)
	}

	lExistingUser, lErr := CheckUser(pdebug, pDb, lUserData.ClientId, lUserData.Email)
	if lErr != nil {
		pdebug.Log(helpers.Elog, "GNU003"+lErr.Error())
		return fmt.Errorf("GNU003 : %v", lErr)
	}

	if lExistingUser == "Y" {
		return fmt.Errorf("olduser")
	}

	lCoreString := `insert into userconfig (clientid,username,email,password,createdby,createddate,updatedby,updateddate) 
	values (?,?,?,?,?,now(),?,now())`

	_, lErr = pDb.Exec(lCoreString, lUserData.ClientId, lUserData.UserName, lUserData.Email, lUserData.Password, lCreatedBy, lCreatedBy)
	if lErr != nil {
		pdebug.Log(helpers.Elog, "GNU004"+lErr.Error())
		return fmt.Errorf("GNU004 : %v", lErr)
	}

	pdebug.Log(helpers.Statement, "HandleNewUser (-)")
	return nil
}

func CheckUser(pDebug *helpers.HelperStruct, pDb *sql.DB, pClientid, pEmail string) (string, error) {
	pDebug.Log(helpers.Statement, "CheckUser (+)")
	var flag string

	lCorestring := `SELECT CASE WHEN count(*) > 0 THEN 'Y' ELSE 'N' END AS Flag
		FROM userconfig
		WHERE clientid  = ? or email=?`

	rows, lErr := pDb.Query(lCorestring, pClientid, pEmail)
	if lErr != nil {
		return "", helpers.ErrReturn(lErr)

	}

	for rows.Next() {
		lErr := rows.Scan(&flag)
		if lErr != nil {
			return "", helpers.ErrReturn(lErr)
		}
	}
	pDebug.Log(helpers.Statement, "CheckUser (-)")
	return flag, nil
}

type LogReqStruct struct {
	Password string `json:"password"`
	ClientId string `json:"clientid"`
}

func Login(w http.ResponseWriter, r *http.Request) {
	lDebug := new(helpers.HelperStruct)
	lDebug.SetUid(r)
	lDebug.Log(helpers.Statement, "Login (+)")
	(w).Header().Set("Access-Control-Allow-Origin", "*")
	(w).Header().Set("Access-Control-Allow-Credentials", "true")
	(w).Header().Set("Access-Control-Allow-Methods", "PUT, OPTIONS")
	(w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

	switch r.Method {
	case "PUT":

		var lUserData LogReqStruct
		lDb, lErr := ftdb.LocalDbConnect(ftdb.MariaFTPRD)

		if lErr != nil {
			lDebug.Log(helpers.Elog, "Log001"+lErr.Error())
			fmt.Fprint(w, helpers.GetError_String("Log001", "Something went wrong. Please try agin later."))
			return
		}
		defer lDb.Close()

		lByte, lErr := ioutil.ReadAll(r.Body)
		if lErr != nil {
			lDebug.Log(helpers.Elog, "Log002"+lErr.Error())
			return
		}

		lErr = json.Unmarshal(lByte, &lUserData)
		if lErr != nil {
			lDebug.Log(helpers.Elog, "Log003"+lErr.Error())
			return
		}

		lDebug.Log(helpers.Statement, lUserData, "lUserData    ")
		lErr = HandleLoginReq(lDebug, lDb, lUserData)
		if lErr != nil {
			lDebug.Log(helpers.Elog, "Log004"+lErr.Error())
			fmt.Fprint(w, helpers.GetError_String("Log002", "Invalid Credential."))
			return
		}
		lErr = appsession.KycSetcookie(w, lDebug, common.CHATName, lUserData.ClientId, common.CookieMaxAge)
		if lErr != nil {
			lDebug.Log(helpers.Elog, "Log004"+lErr.Error())
			fmt.Fprint(w, helpers.GetError_String("Log002", "Invalid Credential."))
			return
		}

		fmt.Fprint(w, helpers.GetMsg_String("S", "Valid Credential."))
		lDebug.Log(helpers.Statement, "Login (-)")

	}
}

func HandleLoginReq(pdebug *helpers.HelperStruct, pDb *sql.DB, lUserData LogReqStruct) error {
	pdebug.Log(helpers.Statement, "HandleLoginReq (+)")
	// var lUserData LogReqStruct

	lExistingUser, lErr := ValidateUser(pdebug, pDb, lUserData.ClientId, lUserData.Password)
	if lErr != nil {
		pdebug.Log(helpers.Elog, "HLR003"+lErr.Error())
		return fmt.Errorf("HLR003 : %v", lErr)
	}
	if lExistingUser == "N" {
		// pdebug.Log(helpers.Elog, "HLR004"+lErr.Error())
		// fmt.Println("sdfgsdf", lExistingUser)
		return fmt.Errorf("HLR004 : %v", lErr)
	}

	pdebug.Log(helpers.Statement, "HandleLoginReq (-)")
	return nil
}

func ValidateUser(pDebug *helpers.HelperStruct, pDb *sql.DB, pClientid, pPassword string) (string, error) {
	pDebug.Log(helpers.Statement, "ValidateUser (+)")
	var flag string

	lCorestring := `SELECT CASE WHEN count(*) > 0 THEN 'Y' ELSE 'N' END AS Flag
		FROM userconfig
		WHERE clientid  = ? and password= ?`

	rows, lErr := pDb.Query(lCorestring, pClientid, pPassword)
	if lErr != nil {
		return "", helpers.ErrReturn(lErr)
	}

	for rows.Next() {
		lErr := rows.Scan(&flag)
		if lErr != nil {
			return "", helpers.ErrReturn(lErr)
		}
	}
	pDebug.Log(helpers.Statement, "ValidateUser (-)")
	return flag, nil
}
