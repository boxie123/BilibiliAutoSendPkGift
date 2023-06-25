package utils

// 配置信息
type ConfigInfo struct {
	AccessKey string `json:"accessKey"`
	Cookie    string `json:"cookie"`
	RoomId    int    `json:"roomId"`
}

// bilibili api 普遍返回数据格式
type ApiResponseCommon struct {
	Code    int                    `json:"code"`
	Message string                 `json:"message"`
	Data    map[string]interface{} `json:"data"`
}

// 包裹礼物信息
type BagListData struct {
	List []BagGiftInfo `json:"list"`
	Time int           `json:"time"`
}

type BagGiftInfo struct {
	BagID       int64  `json:"bag_id"`
	GiftID      int    `json:"gift_id"`
	GiftName    string `json:"gift_name"`
	GiftNum     int    `json:"gift_num"`
	GiftType    int    `json:"gift_type"`
	ExpireAt    int64  `json:"expire_at"`
	CornerMark  string `json:"corner_mark"`
	CornerColor string `json:"corner_color"`
	CountMap    []struct {
		Num   int    `json:"num"`
		Text  string `json:"text"`
		Flags []int  `json:"flags"`
	} `json:"count_map"`
	BindRoomID   int    `json:"bind_roomid"`
	BindRoomText string `json:"bind_room_text"`
	Type         int    `json:"type"`
	CardImage    string `json:"card_image"`
	CardGif      string `json:"card_gif"`
	CardID       int    `json:"card_id"`
	CardRecordID int    `json:"card_record_id"`
	IsShowSend   bool   `json:"is_show_send"`
	ExpireText   string `json:"expire_text"`
	MaxSendLimit int    `json:"max_send_limit"`
	DiyCountMap  int    `json:"diy_count_map"`
}
