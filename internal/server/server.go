package server

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"websockets_htmx_sysinfo/internal/service"

	"github.com/gorilla/websocket"
)

type Server struct {
	subscriberMessageBuffer int
	mux                     http.ServeMux
	subscribersMu           sync.Mutex
	subscribers             map[*subscriber]struct{}
	hardwareService         *service.HardwareService
	upgrader                websocket.Upgrader
}

type subscriber struct {
	conn *websocket.Conn
	msgs chan []byte
}

func NewServer(hardwareService *service.HardwareService) *Server {
	s := &Server{
		subscriberMessageBuffer: 10,
		subscribers:             make(map[*subscriber]struct{}),
		hardwareService:         hardwareService,
		upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin: func(r *http.Request) bool {
				return true // Разрешить все источники
			},
		},
	}
	s.mux.Handle("/", http.FileServer(http.Dir("./htmx")))
	s.mux.HandleFunc("/ws", s.subscribeHandler)
	return s
}

func (s *Server) Router() *http.ServeMux {
	return &s.mux
}

func (s *Server) subscribeHandler(w http.ResponseWriter, r *http.Request) {

	conn, err := s.upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Upgrade error:", err)
		return
	}
	defer func() {
		fmt.Println("Closing connection")
		conn.Close()
	}()

	subscriber := &subscriber{
		conn: conn,
		msgs: make(chan []byte, s.subscriberMessageBuffer),
	}
	s.addSubscriber(subscriber)
	defer func() {
		fmt.Println("Removing subscriber")
		s.removeSubscriber(subscriber)
	}()

	s.handleSubscriber(subscriber)
}

func (s *Server) addSubscriber(subscriber *subscriber) {
	s.subscribersMu.Lock()
	s.subscribers[subscriber] = struct{}{}
	s.subscribersMu.Unlock()
}

func (s *Server) removeSubscriber(subscriber *subscriber) {
	s.subscribersMu.Lock()
	delete(s.subscribers, subscriber)
	s.subscribersMu.Unlock()
}

func (s *Server) handleSubscriber(subscriber *subscriber) {
	for msg := range subscriber.msgs {
		err := subscriber.conn.WriteMessage(websocket.TextMessage, msg)
		if err != nil {
			fmt.Println("Write message error:", err)
			return
		}
	}
}

func (s *Server) publishMsg(msg []byte) {
	s.subscribersMu.Lock()
	defer s.subscribersMu.Unlock()

	for sub := range s.subscribers {
		select {
		case sub.msgs <- msg:
		default:
			close(sub.msgs)
			delete(s.subscribers, sub)
		}
	}
}

func FormatUpdateTimestamp(timestamp time.Time) string {
	formattedTime := timestamp.Format("2006-01-02 15:04:05")
	return `<div id="update-timestamp"><p><i class="fa fa-circle" style="color: red;"></i> Last update: ` + formattedTime + `</p></div>`
}

func (s *Server) PublishSystemData() error {
	systemData, err := s.hardwareService.GetSystemSection()
	if err != nil {
		return err
	}
	diskData, err := s.hardwareService.GetDiskSection()
	if err != nil {
		return err
	}
	cpuData, err := s.hardwareService.GetCPUSection()
	if err != nil {
		return err
	}

	timestampHTML := FormatUpdateTimestamp(time.Now())

	// Отправка одним сообщением
	// msg := fmt.Sprintf(
	// 	"%s<div class='row'><div class='col-md-6'>%s</div><div class='col-md-6'>%s</div><div class='col-md-6'>%s</div></div>",
	// 	timestampHTML, systemData, diskData, cpuData,
	// )
	// s.publishMsg([]byte(msg))

	// Отправка обновлений по отдельности
	s.publishMsg([]byte("UPDATE_TIMESTAMP:" + timestampHTML))
	s.publishMsg([]byte("UPDATE_SYSTEM_DATA:" + systemData))
	s.publishMsg([]byte("UPDATE_DISK_DATA:" + diskData))
	s.publishMsg([]byte("UPDATE_CPU_DATA:" + cpuData))

	return nil
}
