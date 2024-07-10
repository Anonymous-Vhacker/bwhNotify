package service

import (
	"bwhNotify/config"
	"bwhNotify/logger"
	"bwhNotify/model"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"
)

const (
	bwhApiPrefix         = "https://api.64clouds.com/v1/"
	bwhApiGetServiceInfo = "getServiceInfo" // 获取vps信息
)

type BwhApiService struct {
	Hosts []*model.HostInfo // vps主机
}

var (
	bwhApiOnce sync.Once
	bwhApiSrv  *BwhApiService
	client     = &http.Client{
		Timeout: time.Second * 10,
	}
)

func InitBwhApiService(conf []config.BWHostConf) error {
	bwhApiOnce.Do(func() {
		bwhApiSrv = &BwhApiService{}
		for _, host := range conf {
			bwhApiSrv.Hosts = append(bwhApiSrv.Hosts, &model.HostInfo{
				Veid:   host.Veid,
				ApiKey: host.ApiKey,
			})
		}
	})
	return nil
}

func GetBwhApiService() *BwhApiService {
	return bwhApiSrv
}

func (bwh *BwhApiService) GetServiceInfos() []model.HostInfo {
	infos := make([]model.HostInfo, 0)
	for _, host := range bwh.Hosts {
		res, err := apiGetServiceInfo(*host)
		if err != nil { // 获取失败，5秒后再次获取
			logger.Log().Error("GetServiceInfo failed.", "err:", err, "veid:", host.Veid)
			for i := 0; i < 3; i++ {
				time.Sleep(5 * time.Second)
				if res, err = apiGetServiceInfo(*host); err != nil {
					logger.Log().Error("GetServiceInfo failed.", "err:", err, "veid:", host.Veid)
				} else {
					break
				}
			}
		}
		host.UpdateHostInfo(res)
		infos = append(infos, *host)
	}
	return infos
}

func apiGetServiceInfo(h model.HostInfo) (model.BwhApiGetServiceInfoResponse, error) {
	res := model.BwhApiGetServiceInfoResponse{Error: -1}
	req, err := http.NewRequest(http.MethodGet, bwhApiPrefix+bwhApiGetServiceInfo, nil)
	if err != nil {
		return res, err
	}
	// 请求query参数
	q := req.URL.Query()
	q.Add("veid", h.Veid)
	q.Add("api_key", h.ApiKey)
	req.URL.RawQuery = q.Encode()
	// 发送请求
	resp, err := client.Do(req)
	if err != nil {
		return res, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil || resp.StatusCode != http.StatusOK {
		return res, fmt.Errorf("failed to get service info. error: [%d] %s", resp.StatusCode, string(body))
	}
	err = json.Unmarshal(body, &res)
	if err != nil {
		return res, err
	}
	if res.Error != 0 {
		return res, fmt.Errorf("error code: %d. %s", res.Error, res.Message)
	}
	return res, nil
}
