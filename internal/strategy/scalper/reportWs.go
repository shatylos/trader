package scalper

import (
	"bytes"
	"github.com/gorilla/websocket"
	"github.com/shatylos/trader/tools/apperrors"
	"github.com/shatylos/trader/tools/logger"
	"github.com/shatylos/trader/web/helper"
	"html/template"
	"net/http"
	"sync"
	"time"
)

type WebSocket struct {
	Clients map[*websocket.Conn]bool
	mu      sync.Mutex
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

	ws.mu.Lock()
	ws.Clients[conn] = true
	ws.mu.Unlock()

	for {
		_, _, err = conn.ReadMessage()
		if err != nil {
			ws.mu.Lock()
			delete(ws.Clients, conn)
			ws.mu.Unlock()
			break
		}
	}
}

func (ws *WebSocket) SendState(s *Scalper) {
	ws.sendBlock(s, "current-state")
}

func (ws *WebSocket) SendPositions(s *Scalper) {
	ws.sendBlock(s, "positions")
}

func (ws *WebSocket) SendPNL(s *Scalper) {
	ws.sendBlock(s, "pnl")
}

// sendBlock renders the given template block with the current period report
// data and broadcasts the fragment to all connected clients. The fragment is
// swapped on the page by the element id, so only blocks rendered for the
// current period (id suffixed with -ws) get replaced.
func (ws *WebSocket) sendBlock(s *Scalper, blockName string) {
	if !ws.hasClients() {
		return
	}

	var err error
	var tmpl *template.Template
	tmpl, err = helper.GetTemplate("web/template/scalper/report.html")
	if err != nil {
		err = apperrors.Wrap(err, "error getting template for web socket")
		logger.PrintError(err)
		return
	}

	now := time.Now()
	firstDayOfNextMonth := time.Date(now.Year(), now.Month()+1, 1, 23, 59, 59, 0, now.Location())
	to := firstDayOfNextMonth.AddDate(0, 0, -1)
	from := time.Date(to.Year(), to.Month(), 1, 0, 0, 0, 0, now.Location())

	var data ReportTemplateData
	data, err = s.GetReportData(from, to)
	if err != nil {
		err = apperrors.Wrap(err, "error get report data for web socket")
		logger.PrintError(err)
		return
	}

	var buf bytes.Buffer
	err = tmpl.ExecuteTemplate(&buf, blockName, data)
	if err != nil {
		err = apperrors.Wrap(err, "error executing '%s' template for web socket", blockName)
		logger.PrintError(err)
		return
	}

	ws.mu.Lock()
	defer ws.mu.Unlock()
	for client := range ws.Clients {
		err = client.WriteMessage(websocket.TextMessage, buf.Bytes())
		if err != nil {
			client.Close()
			delete(ws.Clients, client)
		}
	}
}

func (ws *WebSocket) hasClients() bool {
	ws.mu.Lock()
	defer ws.mu.Unlock()
	return len(ws.Clients) > 0
}
