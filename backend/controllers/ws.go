package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	"reat-time-forum/middleware"
	"reat-time-forum/models"
	"reat-time-forum/utils"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// Struct WS
type Manager struct {
	clients ClientList
	sync.RWMutex
}

type ClientList map[*Client]bool

type Client struct {
	Connection       *websocket.Conn
	NbrOfConnections int
	manager          *Manager
	Client_id        int
	Nickname         string
}

func NewManager() *Manager {
	return &Manager{
		clients: make(ClientList),
	}
}

func (m *Manager) ServeWebsocket(w http.ResponseWriter, r *http.Request) {
	cookie := utils.CheckCookie(r)
	if cookie == nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	user_id := models.GetUserId(cookie.Value)
	user_nickname := models.GetUserNickname(user_id)
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println(err)
	}

	client := NewClient(conn, m, user_id, user_nickname)
	m.addClient(client)
	go client.readmessages()
}

func (m *Manager) addClient(client *Client) {
	m.Lock()
	defer m.Unlock()
	m.clients[client] = true
	for c := range m.clients {
		if c != client {
			c.Connection.WriteJSON(map[string]string{"user": "online", "nickname": client.Nickname})
			client.Connection.WriteJSON(map[string]string{"user": "online", "nickname": c.Nickname})
			continue
		}
	}
}

func (m *Manager) removeClient(client *Client) {
	m.Lock()
	defer m.Unlock()

	if _, ok := m.clients[client]; ok {
		client.Connection.Close()
		delete(m.clients, client)
		if client.NbrOfConnections == 1 {
			for c := range m.clients {
				if c.Client_id != client.Client_id {
					c.Connection.WriteJSON(map[string]string{"user": "offline", "nickname": client.Nickname})
				}
			}
		}
	}
}

var rateLimit middleware.RateLimit

func (c *Client) readmessages() {
	defer func() {
		c.manager.removeClient(c)
	}()
	for {

		_, payload, err := c.Connection.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				fmt.Println(err)
			}
			break
		}
		allowed := rateLimit.Allow(c.Connection.RemoteAddr().String())
		if !allowed {
			break
		}
		var msg utils.Message

		if err := json.Unmarshal(payload, &msg); err != nil {
			fmt.Println(err)
		}

		// var sender_nickname string

		receiverClient := getClient(c, msg.To)
		if receiverClient == nil {
			c.Connection.WriteMessage(websocket.TextMessage, []byte(`{"status": "failed", "message": "User is offline"}`))
			continue
		} else {
			err := models.Insertmsg(c.Client_id, receiverClient.Client_id, msg.Content)
			if err != nil {
				fmt.Println(err)
				c.Connection.WriteMessage(websocket.TextMessage, []byte(`{"status": "failed", "message": "internale server error"}`))
			} else {
				for reciever := range c.manager.clients {
					if reciever.Client_id == receiverClient.Client_id {
						reciever.Connection.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf(
							`{"status": "success", "from": "%s", "content": "%s"}`, c.Nickname, msg.Content)))
						c.Connection.WriteMessage(websocket.TextMessage, []byte(`{"status": "successe", "message": "Your message is delevered"}`))
					}
				}
			}
		}
	}
}

// func Insertmsg(sender_id, receiver_id int, content string) error {
// 	query := "INSERT INTO messages (sender_id, reciever_id, content, creation_date) VALUES (?,?,?, strftime('%s', 'now'))"
// 	_, err := db.Exec(query, sender_id, receiver_id, content)
// 	return err
// }

func getClient(c *Client, clinetNickname string) *Client {
	for client := range c.manager.clients {
		if client.Nickname == clinetNickname {
			return client
		}
	}
	return nil
}

func NewClient(conn *websocket.Conn, manager *Manager, user_id int, nickname string) *Client {
	return &Client{
		Connection: conn,
		manager:    manager,
		Client_id:  user_id,
		Nickname:   nickname,
	}
}
