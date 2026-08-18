# バックエンド構成

Go アプリケーションは単一の `main` パッケージで動かす。初心者が変更箇所を見つけやすいよう、
ファイルは責務ごとに分ける。パッケージを細かく分けるのは、独立した再利用や境界が必要になったときだけにする。

```text
application/
├── main.go             # 設定の読込、依存関係の生成、HTTPサーバーの起動
├── app.go              # ハンドラが使う依存関係
├── router.go           # URLとハンドラ、JSON応答、静的ファイル配信
├── auth.go             # Firebaseトークン検証とメール許可リスト
├── manga_handlers.go   # 漫画・巻・ページAPIのHTTP処理
├── manga.go            # カタログから漫画・巻を取得する処理
├── page_urls.go        # ページ画像の署名付きURLをバッチで発行
├── media.go            # GCS署名処理
└── catalog_server.go   # 実行時カタログの準備と読込
```

## 変更時の指針

- 新しいAPIは、まず `router.go` にURLを追加し、HTTP処理を `manga_handlers.go` のような対象別ファイルに置く。
- SQLで取得する処理は、HTTP応答を組み立てる処理から分けて `manga.go` に置く。
- GCSやFirebaseなど外部サービスの呼出しは、対象サービスのファイルに閉じ込める。
- `main.go` にHTTPハンドラやSQLを追加しない。起動手順だけを残す。
- 構造を変える変更と、機能を変える変更は別々にして、各変更後にテストを実行する。
