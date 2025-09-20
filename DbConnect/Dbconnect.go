package DbConnect

// func ConnectDb() (*sql.DB, error) {
// 	log.Println("DB connection (+)")
// 	config := common.ReadTomlConfig("toml/Dbconfig.toml")
// 	NAME := fmt.Sprintf("%v", config.(map[string]interface{})["NAME"])
// 	USER := fmt.Sprintf("%v", config.(map[string]interface{})["USER"])
// 	HOST := fmt.Sprintf("%v", config.(map[string]interface{})["HOST"])
// 	PASSWORD := fmt.Sprintf("%v", config.(map[string]interface{})["PASSWORD"])
// 	PORT := fmt.Sprintf("%v", config.(map[string]interface{})["PORT"])

// 	ConnString := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", USER, PASSWORD, HOST, PORT, NAME)
// 	db, err := sql.Open("mysql", ConnString)
// 	if err != nil {
// 		log.Println("open db connection is failed: " + err.Error())
// 		return db, err
// 	} else {
// 		log.Println("LocalDbconnected  (-)")
// 		return db, nil
// 	}
// }
