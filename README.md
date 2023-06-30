# BilibiliAutoSendPkGift
 一张一张点太麻烦了, 所以 自动送出所有PK票

## 用法

### 下载可执行文件

```cmd
BilibiliAutoSendPkGift.exe your-config-file.json
```
或直接双击运行，会自动生成配置文件

### 手动构建

```cmd
go run main.go your-config-file.json
```

### 配置文件格式
```json
{
    "accessKey": "",    // 非必要
    "cookie": "",       // 登录信息
    "roomId": 1184275   // 要送礼物的房间号
}
```

> 可在自动生成登录信息后添加 `roomId` 项

或直接从浏览器开发者工具中复制cookie字符串
