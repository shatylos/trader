package vwapReversion

import (
	"fmt"
	domainStructs "github.com/shatylos/trader/internal/domain/structs"
	strategyStorage "github.com/shatylos/trader/internal/strategy/vwapReversion/storage"
	"github.com/shatylos/trader/internal/strategy/vwapReversion/storage/mongo"
	"github.com/shatylos/trader/internal/strategy/vwapReversion/structs"
	"github.com/shatylos/trader/tools/apperrors"
	"github.com/shatylos/trader/tools/logger"
	"time"
)

func (v *VwapReversion) AddAssetTransaction(amount float64, dateTime time.Time, transactionType string) (err error) {
	var storage mongo.MongoStorage
	storage, err = strategyStorage.GetStorage(v.GetId())
	if err != nil {
		err = apperrors.Wrap(err, "error get storage")
		return
	}

	modifyAmount := amount
	if transactionType == domainStructs.TransactionTypeWithdraw {
		modifyAmount *= -1
	}

	var positionDateTime time.Time
	positionDateTime, err = modifyPosition(storage, dateTime, modifyAmount)
	if err != nil {
		err = apperrors.Wrap(err, "error modify position")
		return
	}

	err = storage.AddAssetTransaction(structs.AssetTransaction{
		TransactionType: transactionType,
		Amount:          amount,
		CreatedTime:     positionDateTime,
	})
	if err != nil {
		err = apperrors.Wrap(err, "error add asset transaction")
		return
	}

	return
}

func modifyPosition(storage mongo.MongoStorage, dateTime time.Time, amount float64) (positionDateTime time.Time, err error) {
	var position structs.Position
	position, err = storage.GetInternalPositionCreatedBeforeTime(dateTime)
	if err != nil {
		err = apperrors.Wrap(err, "error get internal position created before time")
		return
	}
	positionDateTime = dateTime

	if position.Id != nil &&
		(position.Status == structs.StatusActive ||
			(position.Status == structs.StatusClosed && position.ClosedTime.After(dateTime))) {

		positionDateTime = position.CreatedTime
		newBalanceBefore := position.BalanceBefore + amount
		logger.Info(fmt.Sprintf("Changed position %s. Old BalanceBefore %g. New BalanceBefore %g", *position.Id, position.BalanceBefore, newBalanceBefore))
		position.BalanceBefore = newBalanceBefore
		position.TotalClosePnl = position.BalanceAfter - position.BalanceBefore
		_, err = storage.SaveInternalPosition(position)
		if err != nil {
			err = apperrors.Wrap(err, "error save internal position")
			return
		}
	}

	return
}
