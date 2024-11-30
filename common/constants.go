package common

const (
	SUCCESSCODE = "S"
	ERRORCODE   = "E"
	YES         = "Y"
	INVALID     = "I"
	NO          = "N"
	EMAIL       = "EMAIL"
	MOBILE      = "MOBILE"
)

var AllowOrgin string

type ResponseStruct struct {
	ClientId   string `json:"clientId"`
	ProfilePic string `json:"profilepic"`
	Status     string `json:"status"`
	Errmsg     string `json:"errMsg"`
}

const (
	//--------------WALL APPLICATION CONSTANTS ------------------------

	ABHICookieName       = "ftab_pt"
	ABHIClientCookieName = "ftab_ud"
	//--------------OTHER COMMON CONSTANTS ------------------------
	CookieMaxAge = 300

	TechExcelPrefix = "TECHEXCELPROD.capsfo.dbo."

	SuccessCode  = "S" //success
	ErrorCode    = "E" //error
	LoginFailure = "I" //??
	NcbEnable    = "Y"

	StatusPending = "P" //pending
	StatusApprove = "A" //Approve
	StatusReject  = "R" //Reject
	StatusNew     = "N" //new
	Statement     = "1"
	Detail        = "2"
	Panic         = "P"
	NoPanic       = "NP"
	INSERT        = "INSERT"
	UPDATE        = "UPDATE"
	SUCCESS       = "success"
	FAILED        = "failed"
	PENDING       = "pending"
	AUTOBOT       = "AUTOBOT"
	Mobile        = "M"
	Web           = "W"
)

var ABHIAllowOrigin []string

// var ABHIBrokerId = 0

var ABHIBrokerId = 0

var ABHIFlag string
