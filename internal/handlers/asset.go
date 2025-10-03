package handlers

import (
	"fmt"
	"github.com/shatylos/trader/internal/domain/structs"
	"github.com/shatylos/trader/internal/setup"
	"github.com/shatylos/trader/tools/logger"
	_type "github.com/shatylos/trader/tools/type"
	"github.com/shatylos/trader/web/helper"
	"net/http"
	"strconv"
	"time"
)

type AssetResponseData struct {
	Message string
	SetupId string
}

func AssetFormHandler(w http.ResponseWriter, r *http.Request) {
	setupInst, err := setup.GetSetupByID(r.PathValue("setup_id"))
	if err != nil {
		http.NotFound(w, r)
		return
	}

	template, err := helper.GetTemplate("web/template/assets.html")
	if err != nil {
		logger.Error(err.Error())
		return
	}
	err = template.Execute(w, AssetResponseData{
		SetupId: setupInst.ID,
	})
	if err != nil {
		logger.Error(err.Error())
		return
	}
}

func AssetAddHandler(w http.ResponseWriter, r *http.Request) {

	var err error
	setupInst, err := setup.GetSetupByID(r.PathValue("setup_id"))
	if err != nil {
		http.NotFound(w, r)
		return
	}

	amount, err := _type.ToFloat64(r.FormValue("amount"))
	if err != nil {
		logger.Warning(fmt.Sprintf("Invalid amount value: %s", err.Error()))
		http.Error(w, "Invalid amount value", 400)
		return
	}

	transactionType := r.FormValue("type")
	if transactionType != structs.TransactionTypeDeposit && transactionType != structs.TransactionTypeWithdraw {
		logger.Warning(fmt.Sprintf("Invalid deposit value: %s", transactionType))
		http.Error(w, "Invalid deposit value", 400)
		return
	}

	timezoneOffsetStr := r.FormValue("timezoneOffset")
	offsetMinutes, err := strconv.Atoi(timezoneOffsetStr)
	if err != nil {
		logger.Warning(fmt.Sprintf("Invalid timezone value: %s", err.Error()))
		http.Error(w, "Invalid timezone value", 400)
		return
	}
	loc := time.FixedZone("User Timezone", -offsetMinutes*60)
	dateTime, err := time.ParseInLocation("2006-01-02T15:04", r.FormValue("datetime"), loc)
	if err != nil {
		logger.Warning(fmt.Sprintf("Invalid Date/Time value: %s", err.Error()))
		http.Error(w, "Invalid Date/Time value", 400)
		return
	}

	err = setupInst.Strategy.AddAssetTransaction(amount, dateTime, transactionType)

	if err != nil {
		http.Error(w, "Error adding the record", 400)
		return
	}

	template, err := helper.GetTemplate("web/template/assets.html")
	if err != nil {
		logger.Error(err.Error())
		return
	}

	err = template.Execute(w, AssetResponseData{
		Message: "Success",
		SetupId: setupInst.ID,
	})
}
