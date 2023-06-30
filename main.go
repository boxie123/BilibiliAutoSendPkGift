package main

import (
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/boxie123/BilibiliAutoSendPkGift/utils"
	login "github.com/boxie123/GoBilibiliLogin"
)

func main() {
	_, _, filePath := login.Login()
	_, cookie, roomId := utils.ReaderSetting(filePath)

	client := &http.Client{}
	bagGiftList := utils.GetBagList(client, cookie)

	var wg sync.WaitGroup
	var count int = 0
	var mu sync.Mutex
	giftName := "PK票"
	for _, bagGiftInfo := range bagGiftList {
		if bagGiftInfo.GiftName != giftName {
			continue
		}
		wg.Add(1)
		go func(bagGiftInfo utils.BagGiftInfo) {
			defer wg.Done()

			err := utils.SendGiftFromBag(client, cookie, bagGiftInfo, roomId)
			if err != nil {
				log.Printf("发送礼物失败: %v", err)
			} else {
				mu.Lock()
				count = count + bagGiftInfo.GiftNum
				mu.Unlock()
			}
		}(bagGiftInfo)
	}
	wg.Wait()
	fmt.Printf("共送出 %d 份 %s\n", count, giftName)
}
