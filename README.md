# Mailbox

[![GitHub release](https://img.shields.io/github/release/toomore/mailbox.svg)](https://github.com/toomore/mailbox/releases)
[![license](https://img.shields.io/github/license/toomore/mailbox.svg)](https://github.com/toomore/mailbox/blob/master/LICENSE)

---

## 專案概述 / Overview

**中文：** Mailbox 是一個簡易電子報發送系統，使用 Golang 實作。主要功能包括：建立發送 campaign 資訊、匯入訂閱者資訊（群組標記）、發送 HTML 格式電子報、開信追蹤、連結點擊追蹤與退訂處理。系統設計為在 Docker 容器中運行。

**English:** Mailbox is a simple newsletter sending system built with Go. It supports campaign management, subscriber import (with group tags), HTML email sending, open tracking, link click tracking, and unsubscribe handling. Designed to run in Docker containers.

---

## 快速開始 / Quick Start

**中文：** 執行前需準備 AWS SES、MariaDB、Nginx。基本流程：建立 campaign → 匯入訂閱者 → 設定環境變數 → 發送電子報 → 啟動追蹤 server → 處理退訂通知。

**English:** Requirements: AWS SES, MariaDB, Nginx. Basic flow: create campaign → import subscribers → set env vars → send newsletter → run tracking server → handle unsubscribe notifications.

---

## 安裝與建置 / Installation & Build

### Docker 映像 / Docker Images

**中文：**

- `toomore/mailbox:base`：包含程式碼與相依套件
- `toomore/mailbox:cmd`：僅包含編譯後的執行檔

```bash
sh ./build-base.sh   # 基礎映像
sh ./build-min.sh    # 最小化映像
```

**English:**

- `toomore/mailbox:base`: Includes source and dependencies
- `toomore/mailbox:cmd`: Contains only the compiled binary

### 開發環境 / Development

```bash
sh ./dev-run-mariadb.sh        # 啟動 MariaDB 容器
sh ./dev-run-docker.sh         # 執行應用程式容器
sh ./dev-run-mariadb-client.sh # MariaDB 客戶端
```

---

## 環境變數 / Environment Variables

| 變數 Variable | 說明 Description |
|---------------|------------------|
| `mailbox_ses_key` | AWS SES Access Key |
| `mailbox_ses_token` | AWS SES Secret Token |
| `mailbox_ses_sender` | 發件者信箱，格式：`Name <email@example.com>` |
| `mailbox_web_site` | 追蹤連結網域（不含 https 與結尾斜線），如：`open.example.com` |
| `mailbox_ses_replyto` | （選用）回信信箱 |
| `mailbox_unsubscribe_mailto` | （選用）List-Unsubscribe 的 mailto 信箱，未設定時 fallback `mailbox_ses_replyto` |
| `mailbox_unsubscribe_one_click` | （選用）`true`/`1` 時加上 `List-Unsubscribe-Post: List-Unsubscribe=One-Click` |
| `MAILBOX_DB_DSN` | （選用）資料庫連線字串，覆寫預設 DSN |

---

## CLI 指令 / Commands

完整說明請執行 `mailbox -h` 或參考 [cmd/docs](cmd/docs/mailbox.md)。

### Campaign

```bash
mailbox campaign create                    # 建立 campaign
mailbox campaign list                      # 列出所有 campaign
mailbox campaign hash --cid [cid] --uid [uid]  # 產生追蹤連結
mailbox campaign open [group] [cid1] [cid2]   # 開信狀況
mailbox campaign opencount [group] [cid]      # 開信次數統計
mailbox campaign openhistory [group] [cid]    # 開信歷史
mailbox campaign doors [group] [cid]          # 連結點擊紀錄
```

### User

```bash
mailbox user import ./list.csv     # 匯入訂閱者（CSV 需含 email, groups, f_name, l_name）
mailbox user import ./list.csv -d  # 預覽模式（dry run）
mailbox user update ./list.csv    # 更新訂閱者（需含 alive 欄位）
mailbox user show [group]         # 顯示群組使用者
mailbox user unsubscribe --email user@example.com              # 標記退訂（所有群組）
mailbox user unsubscribe --email user@example.com --group weekly  # 標記退訂（指定群組）
mailbox user unsubscribe --email "a@x.com,b@x.com" --reason "unsub via gmail"  # 多筆 + 原因
mailbox user unsubscribed [group] # 顯示群組已退訂（alive=0）名單與數量
```

### Send

```bash
# 依群組發送
mailbox send -p [html] -t [text] -s "Subject" -g [group] --cid [cid]

# 發送給特定使用者
mailbox send -p [html] -t [text] -s "Subject" --uid="6,12" --cid [cid]

# 預覽模式
mailbox send -p [html] -t [text] -s "Subject" -g [group] --cid [cid] -d
```

### Unsubscribe（Phase 1: manual / 手動退訂流程）

**中文：**

目前採用 `List-Unsubscribe` 郵件標頭，讓 Gmail、Apple Mail 等客戶端顯示內建退訂按鈕。收到退訂通知後由人工標記，後續寄送自動排除。

**English:**

Currently uses the `List-Unsubscribe` mail header so that Gmail, Apple Mail, and other clients display a built-in unsubscribe button. After receiving unsubscribe notifications, recipients are manually marked and automatically excluded from future sends.

#### 設定 / Setup

```bash
# 設定退訂信箱（接收郵件客戶端的退訂通知）
# Set the unsubscribe mailbox (receives unsubscribe notifications from mail clients)
export mailbox_unsubscribe_mailto="sender+unsubscribe@example.com"

# 可選：加上 one-click 退訂標頭（仍為人工處理，但提升客戶端顯示退訂按鈕的機率）
# Optional: add one-click unsubscribe header (still manual processing, improves client UI)
export mailbox_unsubscribe_one_click="true"
```

若未設定 `mailbox_unsubscribe_mailto`，會 fallback 使用 `mailbox_ses_replyto`；兩者都未設定時不加標頭。

If `mailbox_unsubscribe_mailto` is not set, falls back to `mailbox_ses_replyto`; if neither is set, no header is added.

#### 操作流程 / Workflow

```
1. 發信 → 郵件自動帶 List-Unsubscribe header
   Send  → emails include List-Unsubscribe header automatically

2. 收件者透過客戶端退訂 → 退訂通知寄到 mailbox_unsubscribe_mailto
   Recipient unsubscribes via client → notification sent to mailbox_unsubscribe_mailto

3. 人工檢視退訂通知信箱，取得退訂者 email
   Manually review the unsubscribe mailbox, collect unsubscribed emails

4. 標記退訂（擇一）/ Mark as unsubscribed (pick one):
   a. 單筆/多筆即時標記：
      mailbox user unsubscribe --email user@example.com
   b. 批次 CSV 更新（alive=0）：
      mailbox user update ./unsubscribe_list.csv

5. 驗證退訂結果 / Verify:
   mailbox user unsubscribed [group]

6. 後續 send 自動排除 alive=0 的使用者
   Future sends automatically exclude users with alive=0
```

#### Flags 說明 / Flag Reference

| Flag | 說明 Description |
|------|------------------|
| `--email` | 退訂目標 email，支援逗號分隔或重複傳入（會自動 normalize 與去重）。Target email(s), supports comma-separated or repeated values (auto-normalized and deduplicated). |
| `--group` | （選用）限定群組，未帶則套用該 email 的所有群組。Optional group filter; if omitted, applies to all groups. |
| `--reason` | （選用）退訂原因，記錄於 log 輸出，不寫入 DB。Optional reason, logged but not stored in DB. |

### Server

```bash
mailbox server -p :8801   # 啟動追蹤 server，接收開信／點擊紀錄
```

---

## 模板標籤 / Template Tags

| 標籤 Tag | 說明 Description |
|----------|------------------|
| `{{FNAME}}` | 訂閱者 first name |
| `{{LNAME}}` | 訂閱者 last name |
| `{{READER}}` | 開信追蹤連結，例：`<img src="{{READER}}">` |
| `{{WASHI}}...{{/WASHI}}` | 點擊時才替換的個人化連結，例：`{{WASHI}}http://example.com/?lname={{LNAME}}{{/WASHI}}` |

---

## Nginx 設定 / Nginx Config

追蹤 server 需經 Nginx 反向代理，並加入以下 headers：

```nginx
proxy_set_header X-Real-IP $remote_addr;
proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
proxy_set_header X-Args $query_string;
proxy_set_header X-Uri $uri;
```

---

## 授權 / License

[MIT](LICENSE)
