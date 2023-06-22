package utils

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strconv"
	"strings"
)

const (
	ApiBagList  = "https://api.live.bilibili.com/xlive/web-room/v1/gift/bag_list"
	ApiSendGift = "https://api.live.bilibili.com/xlive/revenue/v1/gift/sendBag"
)

func parseCookie(cookieStr string) ([]*http.Cookie, error) {
	// 创建一个包含cookie字符串的简单HTTP请求字符串
	requestStr := fmt.Sprintf("GET / HTTP/1.1\r\nHost: bilibili.com\r\nCookie: %s\r\n\r\n", cookieStr)

	// 使用http.ReadRequest函数解析HTTP请求字符串
	request, err := http.ReadRequest(bufio.NewReader(strings.NewReader(requestStr)))
	if err != nil {
		return nil, err
	}

	// 从请求中获取解析后的cookie
	return request.Cookies(), nil
}

func GetUIDFromCookie(cookieStr string) (int, error) {
	// 从cookie字符串解析uid
	cookieParts := strings.Split(cookieStr, "; ")
	for _, part := range cookieParts {
		kv := strings.Split(part, "=")
		switch kv[0] {
		case "DedeUserID":
			uid, err := strconv.Atoi(kv[1])
			if err != nil {
				return 0, err
			}
			return uid, nil
		}
	}
	return 0, errors.New("cookie中无uid信息")
}

func MakeClient(cookieStr string) *http.Client {
	// 创建一个新的http.Client实例
	client := &http.Client{}

	// 创建一个cookiejar，用于存储cookie
	jar, _ := cookiejar.New(nil)

	// 使用parseCookie函数解析cookie字符串
	cookies, err := parseCookie(cookieStr)
	if err != nil {
		panic(err)
	}

	// 将解析后的cookie添加到cookiejar中
	u, _ := url.Parse("https://api.live.bilibili.com")
	jar.SetCookies(u, cookies)
	client.Jar = jar

	return client
}

func GetBagList(client *http.Client) []BagGiftInfo {
	resp, err := client.Get(ApiBagList)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	var apiResponse ApiResponseBagList
	err = json.NewDecoder(resp.Body).Decode(&apiResponse)
	if err != nil {
		panic(fmt.Sprintf("Error: %v", err))
	}

	if apiResponse.Code != 0 {
		panic(fmt.Sprintf("Error: %s", apiResponse.Message))
	}

	return apiResponse.Data.List
}

func SendGiftFromBag(client *http.Client, bagGiftInfo BagGiftInfo, uid int, roomId int) error {
	params := map[string]interface{}{
		"uid":           uid,
		"bag_id":        bagGiftInfo.BagID,
		"gift_id":       bagGiftInfo.GiftID,
		"gift_num":      bagGiftInfo.GiftNum,
		"platform":      "pc",
		"send_ruid":     0,
		"storm_beat_id": 0,
		"price":         0,
		"biz_code":      "live",
		"biz_id":        roomId,
		"ruid":          roomId,
	}

	jsonData, err := json.Marshal(params)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", ApiSendGift, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	fmt.Println("Response:", string(body))
	return nil
}
