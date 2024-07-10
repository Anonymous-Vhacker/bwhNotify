package util

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"os"
	"time"
)

// GenerateDingTalkSign 根据时间戳和密钥获取钉钉加签字符串
func GenerateDingTalkSign(timestamp int64, secret string) string {
	// 创建待签名的字符串
	stringToSign := fmt.Sprintf("%d\n%s", timestamp, secret)
	// 创建 HMAC-SHA256 签名
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(stringToSign))
	signData := mac.Sum(nil)
	// Base64 编码
	sign := base64.StdEncoding.EncodeToString(signData)
	// URL 编码
	//encodedSign := url.QueryEscape(sign) // 不需要。交给外部http请求的query自动处理
	return sign
}

// FormatSize 格式化字节数大小
func FormatSize(bytesNum int64) string {
	if bytesNum >= 1024*1024*1024 {
		return fmt.Sprintf("%.2f GB", float64(bytesNum)/1024/1024/1024)
	} else if bytesNum >= 1024*1024 {
		return fmt.Sprintf("%.2f MB", float64(bytesNum)/1024/1024)
	} else if bytesNum >= 1024 {
		return fmt.Sprintf("%.2f KB", float64(bytesNum)/1024)
	} else {
		return fmt.Sprintf("%d B", bytesNum)
	}
}

// FileExists 判断文件是否存在
func FileExists(path string) (exists bool, err error) {
	info, err := os.Stat(path)
	if err == nil {
		if info.IsDir() {
			return false, errors.New("path is a directory")
		}
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

var CST = time.FixedZone("CST", 8*3600)
