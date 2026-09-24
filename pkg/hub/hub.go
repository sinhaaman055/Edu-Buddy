package hub

import (
	"edubuddy/pkg/models"
	"encoding/json"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type Client struct {
	Conn *websocket.Conn
	Send chan []byte
	RoomId string 
}
type Room struct{
	mutex sync.Mutex
	Members map[*Client]bool
	Session models.TestSession
}
type Hub struct{
   Rooms map[string]*Room
   Broadcast   chan []byte
   Register   chan *Client
   Unregister chan *Client
}
var RoomHub = Hub{
	Rooms:      make(map[string]*Room),
	Broadcast:  make(chan []byte),
	Register:   make(chan *Client),
	Unregister: make(chan *Client),
}
func(h *Hub) StartRoomTimer(roomID string){
   ticker:=time.NewTicker(1*time.Second)
   room:=h.Rooms[roomID]
   go func(){
	   for {
            select{
			case <-ticker.C:
				room.mutex.Lock()
				if room.Session.TimeRemaiming<=0{
					room.Session.TestIsActive=false
					room.mutex.Unlock()
					ticker.Stop()
                    h.BroadcastToRoom(roomID, "TEST_ENDED", map[string]string{
                     "message": "Test completed successfully!",
                    })
					return
				}
				room.Session.TimeRemaiming--
				room.mutex.Unlock()
			case <-room.Session.TickerDone:
				ticker.Stop()
				return
			}   
	   }
   }()
}
func (h *Hub) BroadcastToRoom(roomID string, event string, data interface{}) {
	room, exist:=h.Rooms[roomID]
	if !exist {
		return
	}
	msg := map[string]interface{}{
		"event": event,
		"data":  data,
	}
	jsonBytes, err := json.Marshal(msg)
	if err != nil {
		return
	}
	for client := range room.Members {
		select {
		case client.Send <- jsonBytes:
		default:
			close(client.Send)
			delete(room.Members, client)
		}
	}
}
func(h *Hub)Run(){
	for{
		select{
		case client:= <-h.Register:
			room,exist:=h.Rooms[client.RoomId]
			if !exist{
				room=&Room{
					Members: make(map[*Client]bool),
				}
				h.Rooms[client.RoomId]=room
			}
			room.Members[client]=true

		case client:=<-h.Unregister:
			room,exist:=h.Rooms[client.RoomId]
			if exist{
				_,ok:=room.Members[client]
				if ok{
					delete(room.Members,client)
					close(client.Send)
				}
				if len(room.Members)==0{
					delete(h.Rooms,client.RoomId)
				}
			}
			
		case msg:= <-h.Broadcast:
			for _, room := range h.Rooms {
				for client := range room.Members {
					select {
					case client.Send <- msg:
					default:
						close(client.Send)
						delete(room.Members, client)
					}
				
				}
			}
		}

	}

}
