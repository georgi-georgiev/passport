package notifications

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

type NotificationHandlers struct {
	service *NotificationService
	log     *zap.Logger
}

func NewNotificationHandlers(service *NotificationService, log *zap.Logger) *NotificationHandlers {
	return &NotificationHandlers{service: service, log: log}
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:    4096,
	WriteBufferSize:   4096,
	EnableCompression: true,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (h *NotificationHandlers) Reader(c *gin.Context) {
	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer ws.Close()
	for {
		//mt, message, err := ws.ReadMessage()
		// if err != nil {
		// 	fmt.Println(err)
		// 	break
		// }

		// fmt.Println("read", string(message))

		notifications, err := h.service.repository.GetAllNotRead(c.Request.Context())
		if err != nil {
			fmt.Println(err)
			break
		}

		for _, notification := range notifications {
			n, err := json.Marshal(notification)
			if err != nil {
				fmt.Println(err)
				break
			}

			err = ws.WriteMessage(websocket.TextMessage, n)
			if err != nil {
				fmt.Println(err)
				break
			}
		}

		time.Sleep(5 * time.Second)
	}
}
