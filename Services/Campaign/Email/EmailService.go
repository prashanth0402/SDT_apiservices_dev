package email

import (
	"SDT_ApiServices/common"
	"log"
	"net/http"
)

func GetCurVersion(w http.ResponseWriter, r *http.Request) {
	log.Println("GetCurVersion (+)", r.Method)
	(w).Header().Set("Access-Control-Allow-Origin", common.AllowOrgin)
	(w).Header().Set("Access-Control-Allow-Credentials", "true")
	(w).Header().Set("Access-Control-Allow-Methods", "POST,OPTIONS")
	(w).Header().Set("Access-Control-Allow-Headers", "Accept,Content-Type,Content-Length,Accept-Encoding,X-CSRF-Token,Authorization,userdevice")
	if r.Method == http.MethodPost {

	}
}
