package utils

import (
	"encoding/json"
	"fmt"
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
func ReaderSetting(filePath string) (string, string, int) {
	var ConfigData, _ = os.ReadFile(filePath)
	var configContent = ConfigInfo{}

	err := json.Unmarshal(ConfigData, &configContent)
	if err != nil {
		panic("读取登录信息失败")
	}

	var cookie = configContent.Cookie
	var accessKey = configContent.AccessKey
	var roomId = configContent.RoomId

	if roomId == 0 {
		log.Println("配置项缺失：roomId")
		fmt.Printf("请输入要送礼物的直播间号：")
		_, err = fmt.Scanf("%d\n", &roomId)
		if err != nil {
			log.Fatal("输入的直播间号错误")
		}
	}

	return accessKey, cookie, roomId
}
