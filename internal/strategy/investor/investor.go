package investor

import (
	"context"
	"fmt"
	"github.com/shatylos/trader/internal/domain"
	domainStructs "github.com/shatylos/trader/internal/domain/structs"
	"github.com/shatylos/trader/internal/strategy/investor/entity"
	"github.com/shatylos/trader/internal/strategy/investor/storage"
	_struct "github.com/shatylos/trader/internal/strategy/investor/struct"
	"github.com/shatylos/trader/tools/logger"
	"time"
)

type Investor struct {
	Config     Config
	provider   domain.SpotDomainInterface
	Timeframes []_struct.Timeframe
	State      State
	Storage    storage.Storage
}

type State struct {
	CurrentPrice float64
	Heap         *entity.Heap
	Wallet       *domainStructs.DomainWallet
}

func (i *Investor) GetId() string {
	return i.Config.Id
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
			return
		}
	}

	for key := range i.Timeframes {
		if i.Timeframes[key].Config.IsHeap {
			err = i.handleHeapTimeframe(ctx, &(i.Timeframes[key]))
			if err != nil {
				return
			}
		} else {
			err = i.handleTimeframe(ctx, &(i.Timeframes[key]))
			if err != nil {
				return
			}
		}
	}
	return
}

func (i *Investor) Wait() {
	time.Sleep(i.Config.TimeoutDuration)
}

func (i *Investor) GetTimeframeItemByDeal(deal *entity.Deal) (timeFrameItem *_struct.Timeframe) {
	for _, timeFrame := range i.Timeframes {
		if timeFrame.Config.Resolution == deal.Timeframe {
			timeFrameItem = &timeFrame
			return
		}
	}
	return
}

func (i *Investor) getContext() (ctx context.Context) {
	ctx = context.Background()
	ctx = context.WithValue(ctx, _struct.CtxSetupKey, i)
	ctx = context.WithValue(ctx, _struct.CtxMainCurrencyKey, i.Config.MainCurrency)
	ctx = context.WithValue(ctx, _struct.CtxTradeCurrencyKey, i.Config.TradeCurrency)

	var heapTimeframeResolution string
	for _, timeframe := range i.Timeframes {
		if timeframe.Config.IsHeap {
			heapTimeframeResolution = timeframe.Config.Resolution
		}
	}
	ctx = context.WithValue(ctx, _struct.CtxHeapTimeframeKey, heapTimeframeResolution)
	return
}
