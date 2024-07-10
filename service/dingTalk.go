package service

import (
	"bwhNotify/config"
	"bwhNotify/model"
	"bwhNotify/util"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"sync"
	"time"
)

const (
	dingTalkApiPrefix = "https://oapi.dingtalk.com/robot/send"
	dingTalkMsgTitle  = `VPS日报 %s
========================`
	dingTalkMsgVpsInfo = `
【主机名称】 %s
【 IP  地址 】 %s
【地理位置】 %s
【每月流量】 %s
【已用流量】 %s
【已用比例】 %s
【昨日用量】 %s
【下个周期】 %s`
	dingTalkMsgVpsErr = `
！获取VPS信息失败！
【veid】 %s
【api_key】 %s`
	dingTalkMsgWarn = `
！此次记录为新的计费周期，故昨日用量可能稍有不准！`
	dingTalkMsgDivider = `
========================`
)

type DingTalkService struct {
	AccessToken string `json:"accessToken"`
	Secret      string `json:"secret"`
}

var (
	dinTalkOnce sync.Once
	dingTalkSrv *DingTalkService
)

func InitDingTalkService(conf config.DingTalkConf) error {
	if conf.AccessToken == "" {
		return fmt.Errorf("access token of dingTalk bot cannot be empty")
	}
	dinTalkOnce.Do(func() {
		dingTalkSrv = &DingTalkService{
			AccessToken: conf.AccessToken,
			Secret:      conf.Secret,
		}
	})
	return nil
}

// GetDingTalkService 获取钉钉机器人服务
func GetDingTalkService() *DingTalkService {
	return dingTalkSrv
}

func NewMsgWithHostInfos(infos []model.HostInfo) string {
	var result string
	result += fmt.Sprintf(dingTalkMsgTitle, time.Now().In(util.CST).Format("2006-01-02 15:04:05"))
	for _, info := range infos {
		if info.Error != 0 {
			result += fmt.Sprintf(dingTalkMsgVpsErr, info.Veid, info.ApiKey)
		} else {
			result += fmt.Sprintf(
				dingTalkMsgVpsInfo,
				info.Hostname,
				info.IpAddress,
				info.NodeLocation,
				util.FormatSize(info.PlanMonthlyData),
				util.FormatSize(info.DataCounter),
				fmt.Sprintf("%.2f%%", float64(info.DataCounter)/float64(info.PlanMonthlyData)*100),
				util.FormatSize(info.LastUsedData),
				time.Unix(info.DataNextReset, 0).In(util.CST).Format("2006-01-02 15:04:05"),
			)
			if info.NewDataReset {
				result += dingTalkMsgWarn
			}
		}
		result += dingTalkMsgDivider
	}
	return result
}

// SendText 给钉钉机器人发送文本消息msg
func (s *DingTalkService) SendText(msg string) error {
	if s.AccessToken == "" {
		return fmt.Errorf("access token of dingTalk bot cannot be empty")
	}
	bodyToSend := genDingTalkBotRequestContent(msg)
	// 生成加签字符串
	timestamp := time.Now().UnixMilli()
	var sign string
	if s.Secret != "" {
		sign = util.GenerateDingTalkSign(timestamp, s.Secret)
	}

	res := model.DingTalkResponseBody{}
	req, err := http.NewRequest(http.MethodPost, dingTalkApiPrefix, bytes.NewReader(bodyToSend))
	if err != nil {
		return err
	}
	// 请求query参数
	q := req.URL.Query()
	q.Add("access_token", s.AccessToken)
	if s.Secret != "" {
		q.Add("timestamp", strconv.FormatInt(timestamp, 10))
		q.Add("sign", sign)
	}
	req.URL.RawQuery = q.Encode()
	req.Header.Add("Content-Type", "application/json")
	// 发送请求
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil || resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to access dingtalk bot api. error: [%d] %s", resp.StatusCode, string(body))
	}
	if err = json.Unmarshal(body, &res); err != nil {
		return err
	}
	if res.Errcode != 0 {
		return fmt.Errorf("dingtalk error code: %d. err: %s", res.Errcode, res.Errmsg)
	}
	return nil
}

// 生成钉钉机器人api的请求体
func genDingTalkBotRequestContent(msg string) []byte {
	requestBody := model.DingTalkRequestBody{
		Msgtype: "text",
		At:      model.DingTalkAtObj{},
		Text: model.DingTalkTextObj{
			Content: msg,
		},
	}
	body, _ := json.Marshal(requestBody)
	return body
}
