Mailbox CMD [![GoDoc](https://godoc.org/github.com/toomore/mailbox/cmd?status.svg)](https://godoc.org/github.com/toomore/mailbox/cmd)
============

四個主要的小程式來運作 **Mailbox**，`mailbox_campaign`, `mailbox_import_csv`,
`mailbox_sender`, `mailbox_server`

也提供將程式編譯後再放入 `alpine`。

CMD
----

### `mailbox_campaign`
建立 `campaign`，包含產生該 `campaign` 的亂數種子。

### `mailbox_import_csv`
匯入訂閱者的資訊。

### `mailbox_sender`
發送電子報，以 **HTML** 格式發送。

### `mailbox_server`
接收開信訊息。

**Docs:** https://godoc.org/github.com/toomore/mailbox/cmd

Docker
-------

### `toomore/mailbox:cmd`
只將編譯過的 `cmd` 程式放入。

    sh ./make.sh;

**Required:** `toomore/mailbox:base`, run `sh ../build-alpine.sh` first.

... and Run

    docker run -it --rm toomore/mailbox:cmd [sh or mailbox's cmd]
