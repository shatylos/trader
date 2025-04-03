package fibonacci

import (
	"fmt"
	"github.com/shatylos/trader/internal/domain"
	"github.com/shatylos/trader/internal/domain/constant"
	domainStructs "github.com/shatylos/trader/internal/domain/structs"
	strategyStorage "github.com/shatylos/trader/internal/strategy/fibonacci/storage"
	"github.com/shatylos/trader/internal/strategy/fibonacci/storage/mongo"
	structs "github.com/shatylos/trader/internal/strategy/fibonacci/structs"
	"github.com/shatylos/trader/internal/strategy/struct"
	"github.com/shatylos/trader/tools"
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

	var internalPosition structs.Position
	var storage mongo.MongoStorage
	storage, err = strategyStorage.GetStorage(f.config.Id)
	if err != nil {
		return
	}
	internalPosition, err = storage.GetLastInternalPosition()

	var providerPosition domainStructs.DomainPosition
	providerPosition, err = f.provider.GetPosition(f.config.CoinPare)
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
		internalPosition, err = f.calculateNewPosition(internalPosition)
	}

	err = f.actionByPosition(internalPosition, providerPosition)

	return
}

func (f *Fibonacci) GetReport(from time.Time, to time.Time) (*_struct.Report, error) {
	//TODO implement me
	panic("implement me")
}

func (f *Fibonacci) Wait() {
	time.Sleep(time.Second * f.config.TimeoutSeconds)
}

func (f *Fibonacci) ResetOrderData() error {
	//TODO implement me
	panic("implement me")
}

func (f *Fibonacci) AddAssetTransaction(amount float64, dateTime time.Time, transactionType string) error {
	//TODO implement me
	panic("implement me")
}
