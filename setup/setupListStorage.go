package setup

import (
	setupStructs "bitbucket.org/shatylos/trader/setup/structs"
	strategyStructs "bitbucket.org/shatylos/trader/strategy/structs"
	"time"
)

var setupList []*setupStructs.Setup

func init() {
	setupListInit()
}

func setupListInit() {
	setupList = make([]*setupStructs.Setup, 0)

	setupList = append(setupList, &setupStructs.Setup{
		DomainCode: "exmo",
		Strategy:   strategyStructs.Strategy{},
	})
}

func LoadNextSetupStep(setupChanelContext chan *setupStructs.Setup) {
	var setupItemResult *setupStructs.Setup

	for {
		for _, setupListItem := range setupList {
			if setupListItem.GetStatus() == setupStructs.StatusReadyForNext {
				setupListItem.SetStatus(setupStructs.StatusInProgress)
				setupItemResult = setupListItem
				break
			}
		}

		if setupItemResult != nil {
			setupChanelContext <- setupItemResult
			return
		}
		time.Sleep(time.Second / 4)
	}
}
