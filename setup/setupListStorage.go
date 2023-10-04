package setup

import (
	"bitbucket.org/shatylos/trader/setup/reader"
	setupStructs "bitbucket.org/shatylos/trader/setup/structs"
	"time"
)

var setupList []*setupStructs.Setup

func init() {
	err := setupListInit()
	if err != nil {
		panic(err)
	}
}

func getReader() ReaderInterface {
	return &reader.YamlReader{}
}

func setupListInit() error {
	setupList = make([]*setupStructs.Setup, 0)
	err := error(nil)

	setupList, err = getReader().GetSetupList()

	if err != nil {
		return err
	}

	return nil
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
