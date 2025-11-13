package investor

import (
	"bytes"
	"fmt"
	"github.com/gorilla/websocket"
	_struct "github.com/shatylos/trader/internal/strategy/investor/struct"
	"github.com/shatylos/trader/tools/logger"
	"github.com/shatylos/trader/tools/math"
	"github.com/shatylos/trader/tools/trading"
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

func (ws *WebSocket) SendCurrentPrice(i *Investor) {
	var err error
	var tmpl *template.Template
	tmpl, err = helper.GetTemplate("web/template/investor/report.html")
	if err != nil {
		logger.Error(fmt.Sprintf("Error getting template for web socket: %v", err))
		return
	}
	availableMain := trading.CurrencyAmountAvailable(i.State.Wallet, i.Config.MainCurrency)
	availableTrade := trading.CurrencyAmountAvailable(i.State.Wallet, i.Config.TradeCurrency)
	var buf bytes.Buffer
	err = tmpl.ExecuteTemplate(&buf, "current-price", ReportTemplateData{
		TradeCurrency:      ws.Config.TradeCurrency,
		MainCurrency:       ws.Config.MainCurrency,
		CurrentPrice:       i.State.CurrentPrice,
		PricePrecision:     int(ws.Config.PricePrecision),
		AvailableMain:      availableMain,
		AvailableTrade:     availableTrade,
		AvailableTotalMain: availableMain + math.Mul(availableTrade, i.State.CurrentPrice),
		QtyPrecision:       int(i.Config.QtyPrecision),
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

func (ws *WebSocket) SendTimeframeItemStatus(timeframe _struct.TimeframeItem) {
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
		logger.Error(fmt.Sprintf("Error executing 'timeframe-item-status' template for web socket: %v", err))
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

func (ws *WebSocket) SendHeapTimeframeStatus(heapTimeframe _struct.HeapTimeframe) {
	var err error
	var tmpl *template.Template
	tmpl, err = helper.GetTemplate("web/template/investor/report.html")
	if err != nil {
		logger.Error(fmt.Sprintf("Error getting template for web socket: %v", err))
		return
	}
	var buf bytes.Buffer
	err = tmpl.ExecuteTemplate(&buf, "heap-timeframe-status", heapTimeframe)
	if err != nil {
		logger.Error(fmt.Sprintf("Error executing 'heap-timeframe-status' template for web socket: %v", err))
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
