package utils

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

const (
	ApiBagList         = "https://api.live.bilibili.com/xlive/web-room/v1/gift/bag_list"         // 获取包裹礼物列表
	ApiSendGift        = "https://api.live.bilibili.com/xlive/revenue/v1/gift/sendBag"           // 发送包裹礼物
	ApiGetRoomPlayInfo = "https://api.live.bilibili.com/xlive/web-room/v1/index/getRoomPlayInfo" // 根据直播间号获取uid
)

// 根据直播间房间号获取uid
func getRoomPlayInfo(client *http.Client, roomId int) int {
	data, err := GetApiResponseData(client, "", fmt.Sprintf("%s?room_id=%d", ApiGetRoomPlayInfo, roomId))
	if err != nil {
		panic(err)
	}

	uid, ok := data["uid"].(float64)
	if !ok {
		panic("Error: uid is not an float64")
	}
	return int(uid)
}

// 从cookie字符串解析获取 uid 和 csrf (bili_jct)
func getInfoFromCookie(cookieStr string) (int, string, error) {
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

// 获取包裹中有的礼物信息列表
func GetBagList(client *http.Client, cookie string) []BagGiftInfo {
	data, err := GetApiResponseData(client, cookie, ApiBagList)
	if err != nil {
		panic(err)
	}

	giftInfoList, err := parseBagGiftInfo(data)
	if err != nil {
		panic(err)
	}

	return giftInfoList
}

// 解析包裹礼物信息
func parseBagGiftInfo(data map[string]interface{}) ([]BagGiftInfo, error) {
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	var bagListData BagListData
	err = json.Unmarshal(jsonBytes, &bagListData)
	if err != nil {
		return nil, err
	}

	return bagListData.List, nil
}

// 解析常见返回值格式
func parseApiResponseCommen(resp *http.Response) (map[string]interface{}, error) {
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}
	var apiResponse ApiResponseCommon
	err := json.NewDecoder(resp.Body).Decode(&apiResponse)
	if err != nil {
		return nil, err
	}

	if apiResponse.Code != 0 {
		return nil, fmt.Errorf("post response error: %s", apiResponse.Message)
	}

	return apiResponse.Data, nil
}

// 通过Get请求api并返回其中的data字段数据
func GetApiResponseData(client *http.Client, cookie string, apiUrl string) (map[string]interface{}, error) {
	req, _ := http.NewRequest("GET", apiUrl, nil)
	if cookie != "" {
		req.Header.Set("Cookie", cookie)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	data, err := parseApiResponseCommen(resp)
	if err != nil {
		return nil, err
	}

	return data, nil
}

// 向api发送 application/x-www-form-urlencoded 格式的 POST 请求, 并获取返回值
func PostApiResponseData(client *http.Client, cookie string, apiUrl string, paramsMap map[string]interface{}) (map[string]interface{}, error) {
	params := url.Values{}
	for key, value := range paramsMap {
		params.Set(key, fmt.Sprintf("%v", value))
	}

	req, err := http.NewRequest("POST", apiUrl, strings.NewReader(params.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Cookie", cookie)

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	data, err := parseApiResponseCommen(resp)
	if err != nil {
		return nil, err
	}

	return data, nil
}

// 发送包裹礼物到特定直播间
func SendGiftFromBag(client *http.Client, cookie string, bagGiftInfo BagGiftInfo, roomId int) error {
	uid, csrf, err := getInfoFromCookie(cookie)
	if err != nil {
		return err
	}
	ruid := getRoomPlayInfo(client, roomId)

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
		"ruid":          ruid,
		"csrf":          csrf,
		"csrf_token":    csrf,
	}

	_, err = PostApiResponseData(client, cookie, ApiSendGift, paramsMap)
	if err != nil {
		return err
	}

	log.Printf("已送出：%s x %d", bagGiftInfo.GiftName, bagGiftInfo.GiftNum)
	return nil
}
