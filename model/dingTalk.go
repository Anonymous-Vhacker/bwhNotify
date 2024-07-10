package model

type DingTalkTextObj struct {
	Content string `json:"content"` // 文本消息的内容
}

type DingTalkAtObj struct {
	IsAtAll   bool     `json:"isAtAll"`   // 是否@所有人
	AtMobiles []string `json:"atMobiles"` // 被@的群成员手机号
	AtUserIds []string `json:"atUserIds"` // 被@的群成员userId
}

type DingTalkRequestQuery struct { // 钉钉机器人api请求query参数
	AccessToken string `query:"access_token"`
	Timestamp   string `query:"timestamp"`
	Sign        string `query:"sign"`
}

type DingTalkRequestBody struct { // 钉钉机器人api请求消息体
	Msgtype string          `json:"msgtype"` // 消息类型 text:文本类型 link:链接消息 markdown:Markdown类型 actionCard:actionCard类型 feedCard:feedCard类型
	Text    DingTalkTextObj `json:"text"`    // 文本类型消息
	At      DingTalkAtObj   `json:"at"`      // 被@的群成员信息
}

type DingTalkResponseBody struct { // 钉钉机器人api回复消息体
	Errcode int    `json:"errcode"`
	Errmsg  string `json:"errmsg"`
}
