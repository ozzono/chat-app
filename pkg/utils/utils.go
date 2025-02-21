package utils

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	netWS "golang.org/x/net/websocket"
)

func NewSocketConnection(w http.ResponseWriter, r *http.Request) (*websocket.Conn, error) {
	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return nil, errors.Wrap(err, "error upgrading connection")
	}

	return conn, nil

}

// NewEchoWebsocketConnection creates a new WebSocket connection using Echo context
func NewEchoWebsocketConnection(c echo.Context) (*netWS.Conn, error) {
	handler := netWS.Handler(WSHandlerfunc)
	handler.ServeHTTP(c.Response(), c.Request())

	conn, err := netWS.Dial(c.Request().RequestURI, "", "http://"+c.Request().Host)
	if err != nil {
		return nil, errors.Wrap(err, "failed to establish nws connection")
	}
	return conn, nil
}

func WSHandlerfunc(ws *netWS.Conn) {
	defer ws.Close()
	for {
		// Read message from client
		var msg string
		err := netWS.Message.Receive(ws, &msg)
		if err != nil {
			log.Println("Read error:", err)
			break
		}
		// Write message back to client
		err = netWS.Message.Send(ws, msg)
		if err != nil {
			log.Println("Write error:", err)
			break
		}
	}
}
