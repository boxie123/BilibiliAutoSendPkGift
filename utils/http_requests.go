package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

const (
	ApiBagList  = "https://api.live.bilibili.com/xlive/web-room/v1/gift/bag_list"
	ApiSendGift = "https://api.live.bilibili.com/xlive/revenue/v1/gift/sendBag"
)

func getInfoFromCookie(cookieStr string) (int, string, error) {
	// 从cookie字符串解析 uid 和 csrf
	cookieParts := strings.Split(cookieStr, "; ")
	var uid int
	var err error
	var csrf string
	for _, part := range cookieParts {
		kv := strings.Split(part, "=")
		switch kv[0] {
		case "DedeUserID":
			uid, err = strconv.Atoi(kv[1])
		case "bili_jct":
			csrf = kv[1]
		}
	}
	return uid, csrf, err
}

func MakeClient() *http.Client {
	// 创建一个新的http.Client实例
	client := &http.Client{}

	return client
}

func GetBagList(client *http.Client, cookie string) []BagGiftInfo {
	req, _ := http.NewRequest("GET", ApiBagList, nil)
	req.Header.Add("Cookie", cookie)

	resp, err := client.Do(req)
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

func SendGiftFromBag(client *http.Client, cookie string, bagGiftInfo BagGiftInfo, roomId int) error {
	uid, csrf, err := getInfoFromCookie(cookie)
	if err != nil {
		return err
	}
	paramsMap := map[string]interface{}{
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
		"ruid":          1485569,
		"csrf":          csrf,
		"csrf_token":    csrf,
	}

	params := url.Values{}
	for key, value := range paramsMap {
		params.Set(key, fmt.Sprintf("%v", value))
	}

	req, err := http.NewRequest("POST", ApiSendGift, strings.NewReader(params.Encode()))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Cookie", cookie)

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var apiResponse ApiResponseCommon
	err = json.NewDecoder(resp.Body).Decode(&apiResponse)
	if err != nil {
		return err
	}

	if apiResponse.Code != 0 {
		return fmt.Errorf("response error when send gift: %s", apiResponse.Message)
	}

	return nil
}
