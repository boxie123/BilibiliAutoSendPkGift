package main

import (
	"fmt"
	"log"
	"sync"

	"github.com/boxie123/BilibiliAutoSendPkGift/utils"
)

func main() {
	var filePath = utils.GetSettingFilePath()
	_, cookie, roomId := utils.ReaderSettingMode(filePath)

	client := utils.MakeClient()
	bagGiftList := utils.GetBagList(client, cookie)

	var wg sync.WaitGroup
	var count int = 0
	var mu sync.Mutex
	for _, bagGiftInfo := range bagGiftList {
		if bagGiftInfo.GiftName != "PK票" {
			continue
		}
		wg.Add(1)
		go func(bagGiftInfo utils.BagGiftInfo) {
			defer wg.Done()

			err := utils.SendGiftFromBag(client, cookie, bagGiftInfo, roomId)
			if err != nil {
				log.Println("发送礼物失败")
			} else {
				mu.Lock()
				count = count + bagGiftInfo.GiftNum
				mu.Unlock()
			}
		}(bagGiftInfo)
	}
	wg.Wait()
	fmt.Printf("共送出 %d 张PK票\n", count)
}
