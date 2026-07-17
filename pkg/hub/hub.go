package hub

import "github.com/gorilla/websocket"

type Client struct {
	Conn *websocket.Conn
	Send chan []byte
	RoomId string 
}
type Room struct{
	Members map[*Client]bool
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
