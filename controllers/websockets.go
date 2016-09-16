package controllers

import (
	"fmt"
	"github.com/astaxie/beego"
	"github.com/gorilla/websocket"
	"github.com/kujtimiihoxha/bc-feature-requests/models"
	"net/http"
)

// WebSocketController handles WebSocket requests.
type WebSocketController struct {
	BaseController
}

// Join method handles WebSocket requests for WebSocketController.
func (wsc *WebSocketController) Join() {
	wsc.Ctx.Request.Header.Add("Authorization", "Bearer "+wsc.GetString("tkn"))
	MustBeAuthenticated(wsc.Ctx)
	res := wsc.NoClientAccessOnly()
	if res != nil {
		beego.Debug("No Access", res)
		wsc.RetError(res)
		return
	}
	uname := wsc.GetString("uname")
	// Upgrade from http request to WebSocket.
	ws, err := websocket.Upgrade(wsc.Ctx.ResponseWriter, wsc.Ctx.Request, nil, 1024, 1024)
	fmt.Println(uname, err)
	if _, ok := err.(websocket.HandshakeError); ok {
		http.Error(wsc.Ctx.ResponseWriter, "Not a websocket handshake", 400)
		return
	} else if err != nil {
		beego.Error("Cannot setup WebSocket connection:", err)
		return
	}
	models.Join(uname, ws)
	defer models.Leave(uname)
	for {
		_, _, err := ws.ReadMessage()
		if err != nil {
			return
		}
	}
}
