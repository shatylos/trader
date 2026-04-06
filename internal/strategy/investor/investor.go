package investor

import (
	"context"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/shatylos/trader/internal/domain"
	domainStructs "github.com/shatylos/trader/internal/domain/structs"
	"github.com/shatylos/trader/internal/strategy/investor/storage"
	_struct "github.com/shatylos/trader/internal/strategy/investor/struct"
	"github.com/shatylos/trader/tools/apperrors"
	"github.com/shatylos/trader/tools/logger"
	"github.com/shatylos/trader/tools/trading"
	"net/http"
	"time"
)

type Investor struct {
	Config        Config
	provider      domain.SpotDomainInterface
	Timeframes    []_struct.TimeframeItem
	HeapTimeframe _struct.HeapTimeframe
	State         State
	prevState     State
	Storage       storage.Storage
	WebSocket     WebSocket
}

type State struct {
	CurrentPrice float64
	Wallet       *domainStructs.DomainWallet
}

func (i *Investor) Init(mux *http.ServeMux) {

	i.WebSocket = WebSocket{
		Clients:   make(map[*websocket.Conn]bool),
		Config:    &i.Config,
		Broadcast: make(chan string),
	}
	mux.HandleFunc(fmt.Sprintf("GET /%s/ws-report", i.GetId()), i.WebSocket.WsHandler)
}

func (i *Investor) GetId() string {
	return i.Config.Id
}

func (i *Investor) IsEnabled() bool {
	return i.Config.Enabled
}

func (i *Investor) GetTitle() string {
	if !i.Config.Enabled {
		return fmt.Sprintf("Investor: %s (%s) (DISABLED)", i.Config.Id, i.Config.CoinPare)
	}
	return fmt.Sprintf("Investor: %s (%s)", i.Config.Id, i.Config.CoinPare)
}

func (i *Investor) DoAction() (err error) {
	if !i.Config.Enabled {
		if i.Config.Verbose {
			logger.Info("Setup is disabled. Skip the action")
		}
		return
	}

	ctx := i.getContext()

	if i.State.Wallet == nil {
		err = i.updateWalletInfo()
		if err != nil {
			err = apperrors.Wrap(err, "error update wallet info")
			return
		}
	}

	for key := range i.Timeframes {
		err = i.handleTimeframe(ctx, &(i.Timeframes[key]))
		if err != nil {
			err = apperrors.Wrap(err, "error handle timeframe")
			return
		}
		if i.Timeframes[key].IsStatusChanged() {
			i.WebSocket.SendTimeframeItemStatus(i.Timeframes[key])
		}
	}
	//err = i.handleHeapTimeframe(ctx, &i.HeapTimeframe)
	//if err != nil {
	//	err = apperrors.Wrap(err, "error handle heap timeframe")
	//	return
	//}
	//if i.HeapTimeframe.IsStatusChanged() {
	//	i.WebSocket.SendHeapTimeframeStatus(i.HeapTimeframe)
	//}
	if i.IsBalanceChanged() {
		i.WebSocket.SendCurrentPrice(i)
	}
	return
}

func (i *Investor) Wait() {
	time.Sleep(i.Config.TimeoutDuration)
}

func (i *Investor) getContext() (ctx context.Context) {
	ctx = context.Background()
	ctx = context.WithValue(ctx, _struct.CtxSetupKey, i)
	ctx = context.WithValue(ctx, _struct.CtxMainCurrencyKey, i.Config.MainCurrency)
	ctx = context.WithValue(ctx, _struct.CtxTradeCurrencyKey, i.Config.TradeCurrency)
	return
}

func (i *Investor) IsBalanceChanged() (isChanged bool) {
	availableMain := trading.CurrencyAmountAvailable(i.State.Wallet, i.Config.MainCurrency)
	availableTrade := trading.CurrencyAmountAvailable(i.State.Wallet, i.Config.TradeCurrency)

	prevAvailableMain := trading.CurrencyAmountAvailable(i.State.Wallet, i.Config.MainCurrency)
	prevAvailableTrade := trading.CurrencyAmountAvailable(i.State.Wallet, i.Config.TradeCurrency)

	if availableMain != prevAvailableMain {
		isChanged = true
	}
	if availableTrade != prevAvailableTrade {
		isChanged = true
	}
	if i.State.CurrentPrice != i.prevState.CurrentPrice {
		isChanged = true
	}
	i.prevState = i.State
	return
}
