package chat

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

// https://github.dev/tinkerbaj/chat-websocket-gin/tree/main
const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// Client struct for websocket connection and message sending
type Client struct {
	RoomID string
	Conn   *websocket.Conn
	send   chan Message
	hub    *Hub
}

// NewClient creates a new client
func NewClient(id string, conn *websocket.Conn, hub *Hub) *Client {
	return &Client{RoomID: id, Conn: conn, send: make(chan Message, 256), hub: hub}
}

// Client goroutine to read messages from client
func (c *Client) Read() {

	defer func() {
		c.hub.unregister <- c
		err := c.Conn.Close()
		if err != nil {
			return
		}
	}()

	c.Conn.SetReadLimit(maxMessageSize)
	err := c.Conn.SetReadDeadline(time.Now().Add(pongWait))
	if err != nil {
		return
	}
	c.Conn.SetPongHandler(func(string) error {
		err := c.Conn.SetReadDeadline(time.Now().Add(pongWait))
		if err != nil {
			return err
		}
		return nil
	})
	for {
		var msg Message
		err := c.Conn.ReadJSON(&msg)
		if err != nil {
			fmt.Println("Error: ", err)
			break
		}
		c.hub.broadcast <- msg
	}
}

// Client goroutine to write messages to client
func (c *Client) Write() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		err := c.Conn.Close()
		if err != nil {
			return
		}
	}()

	for {
		select {
		case message, ok := <-c.send:
			err := c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err != nil {
				return
			}
			if !ok {
				// The hub closed the channel.
				err := c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				if err != nil {
					return
				}
				return
			} else {
				err := c.Conn.WriteJSON(message)
				if err != nil {
					fmt.Println("Error: ", err)
					break
				}
			}
		case <-ticker.C:
			err := c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err != nil {
				return
			}
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}

	}
}

// Close Client closing channel to unregister client
func (c *Client) Close() {
	close(c.send)
}

// ServeWS Function to handle websocket connection and register client to hub and start goroutines
func ServeWS(ctx *gin.Context, roomId string, hub *Hub) {
	fmt.Print(roomId)
	ws, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	client := NewClient(roomId, ws, hub)

	hub.register <- client
	go client.Write()
	go client.Read()
}

//
//import (
//	"context"
//	"fmt"
//	"github.com/imrishuroy/legal-referral/socketio"
//	"github.com/redis/go-redis/v9"
//	"github.com/rs/zerolog/log"
//	"net/http"
//	"sync"
//)
//
//var friends = map[string][]string{
//	"john":  {"alice", "bob"},
//	"alice": {"john"},
//	"bob":   {"john"},
//}
//
////type authorizer interface {
////	Issue(subject string) (string, error)
////	Verify(token string) (string, error)
////}
//
//type Chat struct {
//	remote *socketio.IORedis[Message]
//	local  *socketio.IO[Message]
//	wg     sync.WaitGroup
//	done   chan struct{}
//	evtCh  chan event
//	//authz  authorizer
//	mu sync.RWMutex
//}
//
//func New(channel string, client *redis.Client) (*Chat, func()) {
//	io := socketio.NewIO[Message]()
//	ioredis, closeChannel := socketio.NewIORedis[Message](channel, client)
//
//	chat := &Chat{
//		remote: ioredis,
//		local:  io,
//		done:   make(chan struct{}),
//		evtCh:  make(chan event),
//		//authz:  authz,
//	}
//
//	chat.loopAsync()
//
//	return chat, closeChannel
//}
//
//func (c *Chat) ServeWS(w http.ResponseWriter, r *http.Request) {
//	//token := r.URL.Query().Get("token")
//	//dd
//	//username, err := c.authz.Verify(token)
//	//if err != nil {
//	//	http.Error(w, err.Error(), http.StatusUnauthorized)
//	//	return
//	//}
//
//	username := r.URL.Query().Get("username")
//	log.Info().Msgf("username: %s", username)
//
//	socket, err, flush := c.local.ServeWS(w, r)
//	if err != nil {
//		http.Error(w, err.Error(), http.StatusBadRequest)
//		return
//	}
//
//	defer flush()
//
//	c.evtCh <- Connected{
//		username:  username,
//		sessionID: socket.ID,
//	}
//
//	defer func() {
//		c.evtCh <- Disconnected{
//			username:  username,
//			sessionID: socket.ID,
//		}
//	}()
//
//	for msg := range socket.Listen() {
//		msg.From = username
//
//		users, _ := friends[username]
//		log.Info().Msgf("users: %v", users)
//		for _, user := range users {
//			m := msg
//			m.To = user
//
//			if err := c.emitRemote(m); err != nil {
//				panic(err)
//			}
//		}
//	}
//}
//
//func (c *Chat) loop() {
//	ch, stop := c.remote.Subscribe()
//	defer stop()
//
//	for {
//		select {
//		case <-c.done:
//			return
//		case evt := <-c.evtCh:
//			c.eventProcessor(evt)
//		case msg := <-ch:
//			c.emitLocal(msg)
//		}
//	}
//}
//
//func (c *Chat) loopAsync() {
//	c.wg.Add(1)
//
//	go func() {
//		defer c.wg.Done()
//
//		c.loop()
//	}()
//}
//
//func (c *Chat) eventProcessor(evt event) {
//	switch e := evt.(type) {
//	case Connected:
//		c.connected(e)
//	case Disconnected:
//		c.disconnected(e)
//	default:
//		panic(fmt.Errorf("chat: unhandled event processor: %+v", evt))
//	}
//}
//
//func (c *Chat) addSession(username, sessionID string) error {
//	return c.remote.Client.SAdd(context.Background(), username, []string{sessionID}).Err()
//}
//
//func (c *Chat) removeSession(username, sessionID string) error {
//	return c.remote.Client.SRem(context.Background(), username, []string{sessionID}).Err()
//}
//
//func (c *Chat) getSessions(username string) ([]string, error) {
//	return c.remote.Client.SMembers(context.Background(), username).Result()
//}
//
//func (c *Chat) fetchFriends(username string) {
//	users := friends[username]
//
//	statuses := make([]Friend, len(users))
//	for i, user := range users {
//		statuses[i] = Friend{
//			Username: user,
//			Online:   c.isUserOnline(user),
//		}
//	}
//
//	c.emitLocal(Message{
//		From:    username,
//		To:      username,
//		Type:    MessageTypeFriends,
//		Friends: statuses,
//	})
//}
//
//func (c *Chat) notifyPresence(username string, online bool) {
//	users := friends[username]
//
//	for _, user := range users {
//		msg := Message{
//			From:     username,
//			To:       user,
//			Type:     MessageTypePresence,
//			Presence: &online,
//		}
//
//		if err := c.emitRemote(msg); err != nil {
//			panic(err)
//		}
//	}
//}
//
//func (c *Chat) connected(evt Connected) {
//	if err := c.addSession(evt.username, evt.sessionID); err != nil {
//		c.emitError(evt.sessionID, err)
//
//		return
//	}
//
//	c.fetchFriends(evt.username)
//	c.notifyPresence(evt.username, true)
//}
//
//func (c *Chat) disconnected(evt Disconnected) {
//	if err := c.removeSession(evt.username, evt.sessionID); err != nil {
//		c.emitError(evt.sessionID, err)
//
//		return
//	}
//
//	c.notifyPresence(evt.username, c.isUserOnline(evt.username))
//}
//
//func (c *Chat) isUserOnline(username string) bool {
//	sessions, err := c.getSessions(username)
//	return err == nil && len(sessions) > 0
//}
//
//func (c *Chat) emitError(sessionID string, err error) bool {
//	return c.local.Error(sessionID, err)
//}
//
//func (c *Chat) emitLocal(msg Message) bool {
//	log.Info().Msgf("emitLocal: %v", msg)
//	sessionIDs, _ := c.getSessions(msg.To)
//
//	for _, sid := range sessionIDs {
//		c.local.EmitAny(sid, msg)
//	}
//
//	return true
//}
//
//func (c *Chat) emitRemote(msg any) error {
//	return c.remote.Publish(context.Background(), msg)
//}
