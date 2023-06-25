package utils

import (
	"encoding/json"
	"log"
	"os"
)

// 获取配置文件路径
func GetSettingFilePath() string {
	var FilePath string
	if len(os.Args) <= 1 {
		log.Fatalf("请选择一个配置文件\n")
	} else {
		FilePath = os.Args[len(os.Args)-1]
	}
	_, err := os.Lstat(FilePath)
	if err != nil {
		log.Fatalf("[%v]不存在\n", FilePath)
	}
	log.Printf("配置文件:[%v]\n", FilePath)
	return FilePath
}

// 读取配置文件
func ReaderSettingMode(filePath string) (string, string, int) {
	var LoginData, _ = os.ReadFile(filePath)
	var loginContent = ConfigInfo{}

	err := json.Unmarshal(LoginData, &loginContent)
	if err != nil {
		panic("读取登录信息失败")
	}

	var cookie = loginContent.Cookie
	var accessKey = loginContent.AccessKey
	var roomId = loginContent.RoomId

	return accessKey, cookie, roomId
}
