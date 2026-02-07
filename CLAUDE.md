# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## 專案概述

Mailbox 是一個簡易電子報發送系統，使用 Golang 實作。主要功能包括：建立發送 campaign 資訊、匯入訂閱者資訊（群組標記）、發送 HTML 格式電子報、開信追蹤與連結點擊追蹤。系統設計為在 Docker 容器中運行。

## 核心架構

### 資料流與組件互動

1. **Campaign 管理流程**：使用者首先透過 `campaign.Create()` 建立 campaign，每個 campaign 都有唯一的 ID 和加密種子（seed）。種子用於產生 HMAC 簽章，確保追蹤連結的安全性。

2. **使用者管理流程**：透過 CSV 檔案匯入訂閱者資料到資料庫的 `user` 表。每個使用者有 email、email_uni（標準化後的唯一信箱）、groups（群組標記）、f_name 和 l_name。

3. **郵件發送流程**：
   - `mails.ProcessSend()` 從資料庫查詢使用者清單
   - 對每個使用者，替換 HTML 模板中的個人化標籤（{{FNAME}}, {{LNAME}}, {{READER}}）
   - 如果啟用 `replaceLink`，會掃描所有 `<a>` 標籤和 `{{WASHI}}` 標籤，將連結替換為追蹤連結
   - 使用 goroutine 並發發送，透過 channel 控制並發數量（預設 7）
   - 透過 AWS SES API 發送郵件

4. **追蹤機制**：
   - **開信追蹤**：在郵件中插入 `{{READER}}` 標籤，轉換為包含 HMAC 簽章的圖片連結，當使用者開信時請求該圖片，server 記錄到 `reader` 表
   - **點擊追蹤**：掃描 HTML 中的 `<a>` 標籤，將原始 URL 存入 `links` 表，替換為追蹤連結（/door/ 路徑），點擊時記錄到 `doors` 表後重導向到原始 URL
   - **Washi 標籤**：特殊的連結類型，支援延遲替換個人化參數（如 `{{WASHI}}http://example.com/?name={{FNAME}}{{/WASHI}}`），在點擊時才進行個人化替換

5. **安全性設計**：所有追蹤連結都使用 HMAC-SHA256 驗證，基於 campaign seed、campaign ID 和 user ID 生成。Server 端驗證 HMAC 後才記錄追蹤資料，防止偽造追蹤請求。

### 主要套件職責

- **campaign**：管理 campaign 的建立、種子儲存與 HMAC 生成/驗證
- **mails**：處理郵件發送、HTML 模板替換、追蹤連結生成、AWS SES 整合
- **reader**：記錄開信追蹤資料
- **utils**：提供資料庫連線、HMAC 工具函數、email 標準化函數
- **cmd**：Cobra CLI 指令實作，包含 campaign、user、send、server 子指令

### 資料庫架構

- **campaign**：儲存 campaign ID 與加密種子
- **user**：訂閱者資訊，包含 email、email_uni、groups、名字、alive 狀態
- **reader**：開信紀錄，記錄 uid、cid、IP、user agent
- **doors**：點擊紀錄，記錄 uid、cid、linkid、IP、user agent
- **links**：連結對應表，儲存 campaign 中使用的原始 URL 及其 hash
- **vote**：投票紀錄表

連線字串固定為：`root:mailboxdbs@tcp(MARIADB:3306)/mailbox?parseTime=true`

## 常用開發指令

### 建置與測試

```bash
# 建置專案
go build -v ./...

# 執行所有測試
go test -v ./...

# 執行 race detector 測試
go test -race ./campaign...
go test -race ./mails...
go test -race ./reader...
go test -race ./utils...

# 執行測試覆蓋率檢查
sh ./goclean.sh

# 執行效能測試
go test -v -bench=Bench -benchmem -run='Bench' ./mails
go test -v -bench=Bench -benchmem -run='Bench' ./utils
```

### Docker 映像建置

```bash
# 建置基礎映像（包含程式碼與相依套件）
sh ./build-base.sh

# 建置最小化映像（只包含編譯後的執行檔）
sh ./build-min.sh
```

### 開發環境設定

```bash
# 啟動 MariaDB 容器（開發用）
sh ./dev-run-mariadb.sh

# 啟動 MariaDB 客戶端
sh ./dev-run-mariadb-client.sh

# 執行應用程式容器
sh ./dev-run-docker.sh

# 初始化資料庫（CI 環境）
mysql -h MARIADB -uroot -pmailboxdbs < ./sql/database.sql
mysql -h MARIADB -uroot -pmailboxdbs mailbox < ./sql/tables.sql
```

### CLI 指令使用

```bash
# 建立新的 campaign
mailbox campaign create

# 列出所有 campaign
mailbox campaign list

# 產生追蹤連結
mailbox campaign hash --cid [campaign_id] --uid [user_id]

# 匯入使用者資料（dry run）
mailbox user import ./list.csv -d

# 匯入使用者資料
mailbox user import ./list.csv

# 更新使用者資料
mailbox user update ./list.csv

# 顯示群組使用者
mailbox user show [group_name]

# 發送電子報（dry run）
mailbox send -p [html_path] -t [text_path] -s "Subject" -g [group] --cid [cid] -d

# 發送電子報
mailbox send -p [html_path] -t [text_path] -s "Subject" -g [group] --cid [cid]

# 發送給特定使用者
mailbox send -p [html_path] -t [text_path] -s "Subject" --uid="6,12" --cid [cid]

# 啟動追蹤 server
mailbox server -p :8801

# 查看 campaign 開信狀況
mailbox campaign open [group] [cid1] [cid2] ...

# 查看開信次數統計
mailbox campaign opencount [group] [cid]

# 查看開信歷史
mailbox campaign openhistory [group] [cid]

# 查看連結點擊紀錄
mailbox campaign doors [group] [cid]
```

## 環境變數設定

執行應用程式需要以下環境變數：

- `mailbox_ses_key`：AWS SES Access Key
- `mailbox_ses_token`：AWS SES Secret Token
- `mailbox_ses_sender`：發件者信箱（格式：`Sender Name <email@example.com>`）
- `mailbox_web_site`：追蹤連結的網域（不含 https 與結尾斜線，如：`open.example.com`）
- `mailbox_ses_replyto`：（選用）回信信箱

## HTML 模板標籤

郵件 HTML 檔案支援以下替換標籤：

- `{{FNAME}}`：訂閱者的 first name
- `{{LNAME}}`：訂閱者的 last name
- `{{READER}}`：開信追蹤連結（通常放在 `<img src="{{READER}}">`）
- `{{WASHI}}...{{/WASHI}}`：延遲替換的個人化連結（如：`{{WASHI}}http://example.com/?name={{FNAME}}{{/WASHI}}`）

## Nginx 設定需求

追蹤 server 需要 Nginx 反向代理，並在設定檔中加入以下 headers：

```nginx
proxy_set_header X-Real-IP $remote_addr;
proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
proxy_set_header X-Args $query_string;
proxy_set_header X-Uri $uri;
```

## 重要實作細節

### Email 標準化邏輯

`utils.FormatEmail()` 函數會標準化 email 地址：
1. 移除 `+` 符號及其後到 `@` 之前的所有字元
2. 移除本地部分（@ 之前）的所有 `.` 符號
3. 轉換為小寫並去除空白

這用於 `email_uni` 欄位，確保同一個實體信箱的不同變體（如 `user+tag@gmail.com` 和 `user@gmail.com`）被識別為同一使用者。

### 並發控制

郵件發送使用 channel 作為信號量（semaphore）控制並發數量，預設為 7。每次發送前會往 channel 送入一個結構體，完成後取出，確保同時最多只有 N 個 goroutine 在發送郵件。

### HMAC 快取機制

`campaign.GetSeed()` 會快取已查詢過的 campaign seed，避免重複查詢資料庫。快取存放在 `cacheSeed` map 中。

### 錯誤重試

`mails.Send()` 在發送失敗時會自動重試最多 5 次。
