package fibonacci

import (
	"fmt"
	strategyStorage "github.com/shatylos/trader/internal/strategy/fibonacci/storage"
	"github.com/shatylos/trader/internal/strategy/fibonacci/storage/mongo"
	"github.com/shatylos/trader/internal/strategy/fibonacci/structs"
	"github.com/shatylos/trader/tools/logger"
	"time"
)

func (f *Fibonacci) AddAssetTransaction(amount float64, dateTime time.Time, transactionType string) (err error) {
	var storage mongo.MongoStorage
	storage, err = strategyStorage.GetStorage(f.GetId())
	if err != nil {
		return
	}

	modifyAmount := amount
	if transactionType == "withdraw" {
		modifyAmount *= -1
	}

	var positionDateTime time.Time
	positionDateTime, err = modifyPosition(storage, dateTime, modifyAmount)
	if err != nil {
		return
	}

	err = storage.AddAssetTransaction(structs.AssetTransaction{
		TransactionType: transactionType,
		Amount:          amount,
		CreatedTime:     positionDateTime.Unix(),
	})
	if err != nil {
		return
	}

	return
}

func modifyPosition(storage mongo.MongoStorage, dateTime time.Time, amount float64) (positionDateTime time.Time, err error) {
	var position structs.Position
	position, err = storage.GetInternalPositionCreatedBeforeTime(dateTime)
	if err != nil {
		return
	}
	positionDateTime = dateTime

	if position.Id != nil &&
		(position.Status == structs.StatusActive ||
			(position.Status == structs.StatusClosed && position.ClosedTime > dateTime.Unix())) {

		positionDateTime = time.Unix(position.CreatedTime, 0)
		newBalanceBefore := position.BalanceBefore + amount
		logger.Info(fmt.Sprintf("Changed position %s. Old BalanceBefore %g. New BalanceBefore %g", *position.Id, position.BalanceBefore, newBalanceBefore))
		position.BalanceBefore = newBalanceBefore
		position.TotalClosePnl = position.BalanceAfter - position.BalanceBefore
		_, err = storage.SaveInternalPosition(position)
	}

	return
}
