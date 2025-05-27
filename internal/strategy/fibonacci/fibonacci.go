package fibonacci

import (
	"fmt"
	"github.com/shatylos/trader/internal/domain"
	"github.com/shatylos/trader/internal/domain/constant"
	domainStructs "github.com/shatylos/trader/internal/domain/structs"
	strategyStorage "github.com/shatylos/trader/internal/strategy/fibonacci/storage"
	"github.com/shatylos/trader/internal/strategy/fibonacci/storage/mongo"
	"github.com/shatylos/trader/internal/strategy/fibonacci/structs"
	"github.com/shatylos/trader/tools"
	"github.com/shatylos/trader/tools/logger"
	"time"
)

type Fibonacci struct {
	isInit   bool
	config   Config
	provider domain.DomainInterface
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
	if f.provider.GetType() != constant.DomainTypeMargin {
		return tools.AppError{
			Message: "Strategy fibonacci works only with Derivatives provider type",
		}
	}
	f.isInit = true
	return nil
}

func (f *Fibonacci) DoAction() (err error) {

	if !f.config.Enabled {
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

	return
}

func (f *Fibonacci) Wait() {
	time.Sleep(time.Second * f.config.TimeoutSeconds)
}

func (f *Fibonacci) ResetOrderData() error {
	//TODO implement me
	panic("implement me")
}
