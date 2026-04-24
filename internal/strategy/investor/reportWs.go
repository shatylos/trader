package investor

import (
	"bytes"
	"fmt"
	"github.com/gorilla/websocket"
	_struct "github.com/shatylos/trader/internal/strategy/investor/struct"
	"github.com/shatylos/trader/tools/apperrors"
	"github.com/shatylos/trader/tools/logger"
	"github.com/shatylos/trader/tools/math"
	"github.com/shatylos/trader/tools/trading"
	"github.com/shatylos/trader/web/helper"
	"html/template"
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
		err = apperrors.Wrap(err, "web socket handler error")
		logger.PrintError(err)
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
		err = apperrors.Wrap(err, "error getting template for web socket")
		logger.PrintError(err)
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
		err = apperrors.Wrap(err, "error executing 'current-price' template for web socket")
		logger.PrintError(err)
		return
	}

	for client := range ws.Clients {
		err = client.WriteMessage(websocket.TextMessage, buf.Bytes())
		if err != nil {
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
		err = apperrors.Wrap(err, "error getting template for web socket")
		logger.PrintError(err)
		return
	}
	var buf bytes.Buffer
	err = tmpl.ExecuteTemplate(&buf, "timeframe-item-status", timeframe)
	if err != nil {
		err = apperrors.Wrap(err, "error executing 'timeframe-item-status' template for web socket")
		logger.PrintError(err)
		return
	}

	for client := range ws.Clients {
		err = client.WriteMessage(websocket.TextMessage, buf.Bytes())
		if err != nil {
			client.Close()
			delete(ws.Clients, client)
		}
	}
}
