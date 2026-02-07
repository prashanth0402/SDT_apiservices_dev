package chatdata

import (
	"database/sql"
	"fcs23pkg/ftdb"
	"fcs23pkg/helpers"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

type WebSocketServer struct {
	upgrader websocket.Upgrader
}

var clients = make(map[*websocket.Conn]bool)

var broadcast = make(chan Message, 1)

// type WebSocketServerChan struct {
// 	clients    map[*websocket.Conn]bool
// 	broadcast  chan Message
// 	addClient  chan *websocket.Conn
// 	removeClient chan *websocket.Conn
// }

type Message struct {
	SenderId       string `json:"senderid"`
	Receiverid     string `json:"receiverid"`
	Msg            string `json:"msg"`
	Attachmentid   string `json:"attachmentid"`
	ChatRoomId     int    `json:"chatroomid"`
	MsgType        string `json:"msgtype"`
	AttachmentName string `json:"attachmentname"`
	AttachmentType string `json:"attachmenttype"`
}

func NewWebSocketServer() *WebSocketServer {
	return &WebSocketServer{
		upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin: func(r *http.Request) bool {
				return true // You can adjust this to check the origin
			},
		},
	}
}

func (wss *WebSocketServer) Connect(w http.ResponseWriter, r *http.Request, pDb *sql.DB) error {
	// log.Println("Connect(+)")

	lDebug := new(helpers.HelperStruct)
	lDebug.SetUid(r)
	lDebug.Log(helpers.Statement, "Connect (+)")

	// params, err := url.ParseQuery(r.URL.RawQuery)
	// if err != nil {
	// 	lDebug.Log(helpers.Statement, "Failed to parse query parameters:", err)
	// 	return err
	// }
	// paramValue := params.Get("param")
	// log.Println("Received query parameter:", paramValue)
	// lGrpId, _ := strconv.Atoi(paramValue)
	// lGrpId := 5
	conn, err := wss.upgrader.Upgrade(w, r, nil)
	if err != nil {
		lDebug.Log(helpers.Statement, "Failed to Upgrade:", err)
		return err
	}

	clients[conn] = true

	// go wss.HandleMessages(lDebug, pDb, lGrpId)

	go wss.Read(lDebug, pDb, conn)

	// log.Println("Connect(-)")
	lDebug.Log(helpers.Statement, "Connect (-)")

	return nil
}

func (wss *WebSocketServer) Read(pDebug *helpers.HelperStruct, pDb *sql.DB, conn *websocket.Conn) {
	defer func() {
		delete(clients, conn)
		conn.Close()
	}()

	for {
		var msg Message
		lErr := conn.ReadJSON(&msg)
		if lErr != nil {
			pDebug.Log(helpers.Statement, "WebSocket read error:", lErr)
			break
		}

		// fmt.Println(" msg.ChatRoomId ", msg.ChatRoomId)

		switch msg.MsgType {
		case "Store":
			lErr = StoreMsg(pDebug, pDb, msg)
			if lErr != nil {
				pDebug.Log(helpers.Statement, "Read001 ", lErr)
			}
			broadcast <- msg
		case "Delete":
			lErr = DeleteMsg(pDebug, pDb, msg.ChatRoomId)
			if lErr != nil {
				pDebug.Log(helpers.Statement, "Read002 ", lErr)
			}
			broadcast <- msg
		case "Change_Group":
			go wss.HandleMessages(pDebug, pDb)
			broadcast <- msg

		}
	}
}

func (wss *WebSocketServer) HandleMessages(pDebug *helpers.HelperStruct, pDb *sql.DB) {
	// log.Println(msg.ChatRoomId)
	for {
		select {
		case msg := <-broadcast:
			// fmt.Println("HandleMessages  ", msg)
			for client := range clients {
				// fmt.Println("gggggg ", msg.ChatRoomId)
				lRespdata, lErr := ChatHistory(pDebug, pDb, msg.ChatRoomId)
				if lErr != nil {
					pDebug.Log(helpers.Statement, "HM0001:", lErr)
					continue
				}
				// fmt.Println("asgdjsadg ", lRespdata.GroupId)
				lErr = client.WriteJSON(lRespdata)
				if lErr != nil {
					pDebug.Log(helpers.Statement, "HM0003:", lErr)
					client.Close()
					delete(clients, client)
				}
			}
			time.Sleep(5 * time.Second)
		}
	}
}

func WebSocketHandler(w http.ResponseWriter, r *http.Request) {
	web := NewWebSocketServer()

	lDb, lErr := ftdb.LocalDbConnect(ftdb.MariaFTPRD)

	if lErr != nil {
		log.Println("Error on LocaldbConnection", lErr)
		fmt.Fprint(w, helpers.GetError_String("JG002", "Something went wrong. Please try agin later."))
		return
	}

	// fmt.Println(web)
	err := web.Connect(w, r, lDb)
	if err != nil {
		log.Println(err.Error())
	}
	// go web.PeriodicWriter(10 * time.Second)

}
