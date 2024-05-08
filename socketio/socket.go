package socketio

import (
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"sync"
	"time"
)

const (
	// Time allowed to write a message to the peer.
	writeTimeout = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongTimeout = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingTimeout = (pongTimeout * 9) / 10

	maxMessageSize = 512
)

type Socket[T any] struct {
	ID             string
	WriteTimeout   time.Duration
	PingTimeout    time.Duration
	PongTimeout    time.Duration
	MaxMessageSize int64
	conn           *websocket.Conn
	done           chan struct{}
	quit           sync.Once
	wg             sync.WaitGroup
	errCh          chan *SocketError
	readCh         chan T
	writeCh        chan any
}

func NewSocket[T any](conn *websocket.Conn) (*Socket[T], func()) {
	socket := &Socket[T]{
		ID:             uuid.New().String(),
		WriteTimeout:   writeTimeout,
		PongTimeout:    pongTimeout,
		PingTimeout:    pingTimeout,
		MaxMessageSize: maxMessageSize,
		conn:           conn,
		done:           make(chan struct{}),
		errCh:          make(chan *SocketError),
		readCh:         make(chan T),
		writeCh:        make(chan any),
	}

	socket.wg.Add(1)
	go func() {
		defer socket.wg.Done()
		socket.writer()
	}()
	socket.wg.Add(1)
	go func() {
		defer socket.wg.Done()
		socket.reader()
	}()
	return socket, socket.close
}

func (s *Socket[T]) Emit(msg T) bool {
	select {
	case <-s.done:
		return false
	case s.writeCh <- msg:
		return true
	}
}

func (s *Socket[T]) EmitAny(msg any) bool {
	select {
	case <-s.done:
		return false
	case s.writeCh <- msg:
		return true
	}

}

func (s *Socket[T]) Listen() <-chan T {
	return s.readCh
}

func (s *Socket[T]) close() {
	s.quit.Do(func() {
		close(s.done)
		s.wg.Wait()
		err := s.conn.Close()
		if err != nil {
			return
		}
	})
}

func (s *Socket[T]) Error(err *SocketError) bool {
	select {
	case <-s.done:
		return false
	case s.errCh <- err:
		return true
	}
}

func (s *Socket[T]) writer() {
	ticker := time.NewTicker(s.PingTimeout)
	defer ticker.Stop()

	for {
		select {
		case <-s.done:
			writeClose(s.conn)
			return

		case err := <-s.errCh:
			writeError(s.conn, err.Code, err)
			return

		case msg := <-s.writeCh:
			_ = s.conn.SetWriteDeadline(time.Now().Add(s.WriteTimeout))

			if err := s.conn.WriteJSON(msg); err != nil {
				writeError(s.conn, websocket.CloseInternalServerErr, err)
				return
			}

		case <-ticker.C:
			_ = s.conn.SetWriteDeadline(time.Now().Add(s.WriteTimeout))
			if err := writePing(s.conn); err != nil {
				return
			}
		}
	}
}

func (s *Socket[T]) reader() {
	defer close(s.readCh)

	for {
		_ = s.conn.SetReadDeadline(time.Now().Add(s.PongTimeout))
		s.conn.SetReadLimit(s.MaxMessageSize)
		s.conn.SetPongHandler(func(string) error {
			return s.conn.SetReadDeadline(time.Now().Add(s.PongTimeout))
		})

		var msg T
		if err := s.conn.ReadJSON(&msg); err != nil {
			return
		}

		select {
		case <-s.done:
			return
		case s.readCh <- msg:
		}

	}
}
