package investor

import (
	"bytes"
	"fmt"
	"github.com/gorilla/websocket"
	_struct "github.com/shatylos/trader/internal/strategy/investor/struct"
	"github.com/shatylos/trader/tools/logger"
	"github.com/shatylos/trader/web/helper"
	"html/template"
	"log"
	"net/http"
)

type WebSocket struct {
	Clients   map[*websocket.Conn]bool
	Config    *Config
	Broadcast chan string
}

func (ws *WebSocket) WsHandler(w http.ResponseWriter, r *http.Request) {
	var upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { return true },
	}
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		logger.Error(fmt.Sprintf("Web Socket Handler Error: %v", err))
		return
	}
	defer conn.Close()
	ws.Clients[conn] = true

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			delete(ws.Clients, conn)
			break
		}
		fmt.Printf("Got message: %s\n", msg)
		ws.Broadcast <- string(msg)
	}
}

func (ws *WebSocket) SendCurrentPrice(price float64) {
	var err error
	var tmpl *template.Template
	tmpl, err = helper.GetTemplate("web/template/investor/report.html")
	if err != nil {
		logger.Error(fmt.Sprintf("Error getting template for web socket: %v", err))
		return
	}
	var buf bytes.Buffer
	err = tmpl.ExecuteTemplate(&buf, "current-price", ReportTemplateData{
		TradeCurrency:  ws.Config.TradeCurrency,
		MainCurrency:   ws.Config.MainCurrency,
		CurrentPrice:   price,
		PricePrecision: int(ws.Config.PricePrecision),
	})
	if err != nil {
		logger.Error(fmt.Sprintf("Error executing 'current-price' template for web socket: %v", err))
		return
	}

	for client := range ws.Clients {
		err = client.WriteMessage(websocket.TextMessage, buf.Bytes())
		if err != nil {
			log.Println("write:", err)
			client.Close()
			delete(ws.Clients, client)
		}
	}
}

func (ws *WebSocket) SendTimeframeStatus(timeframe _struct.Timeframe) {
	var err error
	var tmpl *template.Template
	tmpl, err = helper.GetTemplate("web/template/investor/report.html")
	if err != nil {
		logger.Error(fmt.Sprintf("Error getting template for web socket: %v", err))
		return
	}
	var buf bytes.Buffer
	err = tmpl.ExecuteTemplate(&buf, "timeframe-item-status", timeframe)
	if err != nil {
		logger.Error(fmt.Sprintf("Error executing 'current-price' template for web socket: %v", err))
		return
	}

	for client := range ws.Clients {
		err = client.WriteMessage(websocket.TextMessage, buf.Bytes())
		if err != nil {
			log.Println("write:", err)
			client.Close()
			delete(ws.Clients, client)
		}
	}
}
