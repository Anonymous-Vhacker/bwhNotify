package task

import (
	"bwhNotify/logger"
	"bwhNotify/service"
	"time"
)

type getBwhStatus struct {
}

func (s getBwhStatus) Run() {
	logger.Log().Info("start to get VPS status")
	infos := service.GetBwhApiService().GetServiceInfos()
	msg := service.NewMsgWithHostInfos(infos)
	logger.Log().Info("start to send VPS status to dingTalk bot")
	err := service.GetDingTalkService().SendText(msg)
	if err != nil {
		logger.Log().Error("failed to send to dingTalk bot", "err", err.Error())
		for i := 0; i < 3; i++ {
			time.Sleep(5 * time.Second)
			if err = service.GetDingTalkService().SendText(msg); err != nil {
				logger.Log().Error("failed to send to dingTalk bot", "err", err.Error())
			} else {
				break
			}
		}
	} else {
		logger.Log().Info("success to send to dingTalk bot")
	}
}
