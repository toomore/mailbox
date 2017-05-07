Mailbox
========
簡易電子報發送系統，使用 Golang 實作

Cmd
----
1. mailbox_campaign：建立 campaign，包含產生該 campaign 的亂數種子。
2. mailbox_import_csv：匯入訂閱者的資訊。
3. mailbox_sender：發送電子報，以 HTML 檔案發送。
4. mailbox_server：接收開信訊息。

相關的操作請參考 `-h` 的說明，但可能什麼都沒說 XD

Docker
-------
### `toomore/mailbox:base`
將基本的程式碼放入，相關的套件也一併下載。

    sh ./build-alpine.sh

### `toomore/mailbox:cmd`
只將編譯過的 `cmd` 程式放入。

    cd ./cmd ; sh ./make.sh;

Required
---------
1. AWS SES `KEY`, `Token`.
2. Update `./prod-run-cmd.sh`, `./dev-run-cmd.sh` files.
  * `mailbox_ses_api`, `mailbox_ses_key`, `mailbox_ses_sender`, `mailbox_web_site`

Nginx config
-------------
    proxy_set_header X-Real-IP $remote_addr;
    proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
    proxy_set_header X-Args $query_string;
    proxy_set_header X-Uri $uri;
