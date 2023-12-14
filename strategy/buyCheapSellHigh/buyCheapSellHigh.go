package buyCheapSellHigh

import (
	"bitbucket.org/shatylos/trader/domain"
	"bitbucket.org/shatylos/trader/domain/constant"
	tradeConst "bitbucket.org/shatylos/trader/trading/constant"
	"bitbucket.org/shatylos/trader/utils"
	"fmt"
	"time"
)

type BuyCheapSellHigh struct {
	isInit                    bool
	Id                        string
	CoinPare                  string
	Domain                    domain.DomainInterface
	MainCurrency              string
	TradeCurrency             string
	Resolution                string
	TimeoutSeconds            time.Duration
	CostRanges                []int64
	PercentRanges             []int64
	LongTermMaxPrice          float64
	LongTermMinPrice          float64
	LongTermPercentBuffer     float64
	PurchaseVolumePrecision   int64
	PurchasePricePrecision    int64
	MinutesToReducePriceRange int64
	AvgPriceCandleLimit       int64
	AvgPriceCandleOffset      int64
	ManualBuyPriceBeforeStart float64
}

func (s *BuyCheapSellHigh) IsInit() bool {
	return s.isInit
}

func (s *BuyCheapSellHigh) Initialise() error {
	if s.Domain.GetType() != constant.DomainTypeSpot {
		return utils.AppError{
			Message: "Strategy buyCheapSellHigh works only with spot domain type",
		}
	}
	s.isInit = true
	return nil
}

func (s *BuyCheapSellHigh) DoAction() error {

	openOrders, err := s.getOpenOrders()
	if err != nil {
		return err
	}

	if len(openOrders) >= 2 {
		err := s.cancelOldOrdersWithBigRanges(openOrders)
		if err != nil {
			return err
		}
		return nil
	}

	if len(openOrders) == 1 {
		err = s.cancelAllOrder(openOrders[0])
		if err != nil {
			return err
		}
		err = s.fillPrices()
		if err != nil {
			return err
		}
		//err = s.calculateHistoryOrderValues()
		//if err != nil {
		//	return err
		//}
	}

	historyOrders, err := s.getHistoryOrders()
	if err != nil {
		return err
	}

	baseCurrencyBalance, tradeCurrencyBalance, err := s.getBalances()
	if err != nil {
		return err
	}

	buyPrice, buyQty, sellPrice, sellQty, err := s.getPricesAndQtysToNewOrders(historyOrders, baseCurrencyBalance, tradeCurrencyBalance)
	if err != nil {
		return err
	}

	if buyPrice > 0 && buyQty > 0 {
		orderId, err := s.setLimitOrder(buyPrice, buyQty, tradeConst.SideBuy)
		if err != nil {
			return err
		}
		err = s.setOrderToStorage(orderId, baseCurrencyBalance, tradeCurrencyBalance)
		if err != nil {
			return err
		}
	} else {
		utils.LogInfo(fmt.Sprintf("[Buy Cheap Sell High] unexpected values for buy orders: buyPrice: %f, buyQty: %f", buyPrice, buyQty))
	}
	if sellPrice > 0 && sellQty > 0 {
		orderId, err := s.setLimitOrder(sellPrice, sellQty, tradeConst.SideSell)
		if err != nil {
			return err
		}
		err = s.setOrderToStorage(orderId, baseCurrencyBalance, tradeCurrencyBalance)
		if err != nil {
			return err
		}
	} else {
		utils.LogInfo(fmt.Sprintf("[Buy Cheap Sell High] unexpected values for sell orders: sellPrice: %f, sellQty: %f", sellPrice, sellQty))
	}

	return nil
}

func (s *BuyCheapSellHigh) Wait() {
	time.Sleep(time.Second * s.TimeoutSeconds)
}
