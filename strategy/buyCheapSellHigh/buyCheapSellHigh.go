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
	isInit                  bool
	CoinPare                string        `yaml:"coin_pare"`
	DomainCode              string        `yaml:"domain_code"`
	MainCurrency            string        `yaml:"main_currency"`
	TradeCurrency           string        `yaml:"trade_currency"`
	Resolution              string        `yaml:"resolution"`
	TimeoutSeconds          time.Duration `yaml:"timeout_seconds"`
	CostRanges              []int64       `yaml:"cost_ranges"`
	PercentRanges           []int64       `yaml:"percent_ranges"`
	LongTermMaxPrice        float64       `yaml:"long_term_max_price"`
	LongTermMinPrice        float64       `yaml:"long_term_min_price"`
	LongTermPercentBuffer   float64       `yaml:"long_term_percent_buffer"`
	PurchaseVolumePrecision int64         `yaml:"purchase_volume_precision"`
}

func (s *BuyCheapSellHigh) IsInit() bool {
	return s.isInit
}

func (s *BuyCheapSellHigh) Initialise() error {
	domainItem, err := domain.GetDomainInterface(s.DomainCode)
	if err != nil {
		return err
	}
	if domainItem.GetType() != constant.DomainTypeSpot {
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
		return nil
	}

	if len(openOrders) == 1 {
		err = s.cancelAllOrder(openOrders[0])
		if err != nil {
			return err
		}
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
		err = s.setLimitOrder(buyPrice, buyQty, tradeConst.SideBuy)
		if err != nil {
			return err
		}
	} else {
		utils.LogInfo(fmt.Sprintf("[Buy Cheap Sell High] unexpected values for buy orders: buyPrice: %f, buyQty: %f", buyPrice, buyQty))
	}
	if sellPrice > 0 && sellQty > 0 {
		err = s.setLimitOrder(sellPrice, sellQty, tradeConst.SideSell)
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
