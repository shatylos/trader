package fibonacci

import (
	"fmt"
	"github.com/shatylos/trader/internal/domain"
	domainStructs "github.com/shatylos/trader/internal/domain/structs"
	strategyStorage "github.com/shatylos/trader/internal/strategy/fibonacci/storage"
	"github.com/shatylos/trader/internal/strategy/fibonacci/storage/mongo"
	"github.com/shatylos/trader/internal/strategy/fibonacci/structs"
	"github.com/shatylos/trader/tools/logger"
	"github.com/shatylos/trader/tools/trading"
	"time"
)

type Fibonacci struct {
	isInit   bool
	config   Config
	provider domain.FuturesDomainInterface
	state    FibState
}

type FibState struct {
	LtTrend           string
	StTrend           string
	MinPriceReview    float64
	MaxPriceReview    float64
	CurrentPrice      float64
	NewOrderCondition string
}

func (f *Fibonacci) GetId() string {
	return f.config.Id
}

func (f *Fibonacci) GetTitle() string {
	return fmt.Sprintf("Fibonacci: %s (%s)", f.config.Id, f.config.CoinPare)
}

func (f *Fibonacci) IsInit() bool {
	return f.isInit
}

func (f *Fibonacci) Initialise() error {
	f.isInit = true
	return nil
}

func (f *Fibonacci) DoAction() (err error) {

	fibState := FibState{}
	if !f.config.Enabled {
		f.state = fibState
		if f.config.Verbose {
			logger.Info("The setup is disabled. Set enabled:1 in config file to enable it.")
		}
		return
	}

	var internalPosition structs.Position
	var storage mongo.MongoStorage
	storage, err = strategyStorage.GetStorage(f.config.Id)
	if err != nil {
		return
	}
	internalPosition, err = storage.GetLastInternalPosition(0)

	var providerPosition domainStructs.DomainPosition
	providerPosition, err = f.provider.GetPosition(f.config.CoinPare)
	if err != nil {
		return
	}

	var currentPrice float64
	currentPrice, err = f.getCurrentPrice()
	if err != nil {
		return
	}

	if providerPosition.Size == 0 {
		if internalPosition.Id != nil && internalPosition.Status != structs.StatusClosed {
			err = f.closeInternalPosition(internalPosition)
			if err != nil {
				return
			}
		}
		internalPosition, err = f.calculateNewPosition()
		if err != nil {
			return
		}
		internalPosition.ProviderPosition = providerPosition
	} else if internalPosition.Status == structs.StatusActive {
		internalPosition.ProviderPosition = providerPosition
		internalPosition, err = storage.SaveInternalPosition(internalPosition)
		if err != nil {
			return
		}
	} else {
		logger.Info("Wait for close current provider position")
		return
	}

	err = f.actionByPosition(internalPosition, currentPrice)

	fibState.MinPriceReview = internalPosition.FibonacciChart.SourceMinPrice
	fibState.MaxPriceReview = internalPosition.FibonacciChart.SourceMaxPrice
	fibState.CurrentPrice = currentPrice
	fibState.LtTrend = internalPosition.LtTrend
	fibState.StTrend = internalPosition.StTrend
	fibState.NewOrderCondition = f.newOrderCondition(internalPosition)

	f.state = fibState
	return
}

func (f *Fibonacci) Wait() {
	time.Sleep(time.Second * f.config.TimeoutSeconds)
}

func (f *Fibonacci) newOrderCondition(internalPosition structs.Position) (condition string) {
	if internalPosition.Trend == trading.TrendShort {
		if internalPosition.Orders.Order1.OrderId == "" &&
			internalPosition.FibonacciChart.EntryPoint1 > 0.0 &&
			internalPosition.FibonacciChart.EntryPoint2 > 0.0 {

			condition = fmt.Sprintf(">%g and <%g", internalPosition.FibonacciChart.EntryPoint1, internalPosition.FibonacciChart.EntryPoint2)
			return
		}
		if internalPosition.Orders.Order2.OrderId == "" &&
			internalPosition.FibonacciChart.EntryPoint2 > 0.0 &&
			internalPosition.FibonacciChart.EntryPoint3 > 0.0 {

			condition = fmt.Sprintf(">%g and <%g", internalPosition.FibonacciChart.EntryPoint2, internalPosition.FibonacciChart.EntryPoint3)
			return
		}
		if internalPosition.Orders.Order3.OrderId == "" &&
			internalPosition.FibonacciChart.EntryPoint3 > 0.0 {

			condition = fmt.Sprintf(">%g", internalPosition.FibonacciChart.EntryPoint3)
			return
		}
	}
	if internalPosition.Trend == trading.TrendLong {
		if internalPosition.Orders.Order1.OrderId == "" &&
			internalPosition.FibonacciChart.EntryPoint1 > 0.0 &&
			internalPosition.FibonacciChart.EntryPoint2 > 0.0 {

			condition = fmt.Sprintf("<%g and >%g", internalPosition.FibonacciChart.EntryPoint1, internalPosition.FibonacciChart.EntryPoint2)
			return
		}
		if internalPosition.Orders.Order2.OrderId == "" &&
			internalPosition.FibonacciChart.EntryPoint2 > 0.0 &&
			internalPosition.FibonacciChart.EntryPoint3 > 0.0 {

			condition = fmt.Sprintf("<%g and >%g", internalPosition.FibonacciChart.EntryPoint2, internalPosition.FibonacciChart.EntryPoint3)
			return
		}
		if internalPosition.Orders.Order3.OrderId == "" &&
			internalPosition.FibonacciChart.EntryPoint3 > 0.0 {

			condition = fmt.Sprintf("<%g", internalPosition.FibonacciChart.EntryPoint3)
			return
		}
	}

	return
}

func (f *Fibonacci) ResetOrderData() error {
	//TODO implement me
	panic("implement me")
}
