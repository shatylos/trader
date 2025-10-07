package investor

import (
	"context"
	"fmt"
	"github.com/shatylos/trader/internal/domain"
	domainStructs "github.com/shatylos/trader/internal/domain/structs"
	"github.com/shatylos/trader/internal/strategy/investor/storage"
	"github.com/shatylos/trader/tools/logger"
	"time"
)

type Investor struct {
	config     Config
	provider   domain.SpotDomainInterface
	Timeframes []Timeframe
	State      State
	Storage    storage.Storage
	Wallet     *domainStructs.DomainWallet
}

type State struct {
	CurrentPrice float64
}

type ctxKey string

const CtxSetupKey ctxKey = "setup"

func (i *Investor) GetId() string {
	return i.config.Id
}

func (i *Investor) GetTitle() string {
	return fmt.Sprintf("Investor: %s (%s)", i.config.Id, i.config.CoinPare)
}

func (i *Investor) DoAction() (err error) {
	if !i.config.Enabled {
		if i.config.Verbose {
			logger.Info("Setup is disabled. Skip the action")
		}
		return
	}

	ctx := context.Background()
	ctx = context.WithValue(ctx, CtxSetupKey, i)

	if i.Wallet == nil {
		err = i.updateWalletInfo()
		if err != nil {
			return
		}
	}

	for key := range i.Timeframes {
		err = i.handleTimeframe(ctx, &(i.Timeframes[key]))
		if err != nil {
			return
		}
	}
	return
}

func (i *Investor) Wait() {
	time.Sleep(time.Second * i.config.TimeoutSeconds)
}

func (i *Investor) getTimeframeItemByDeal(deal *storage.Deal) (timeFrameItem *Timeframe) {
	for _, timeFrame := range i.Timeframes {
		if timeFrame.Config.Resolution == deal.Timeframe {
			timeFrameItem = &timeFrame
			return
		}
	}
	return
}
