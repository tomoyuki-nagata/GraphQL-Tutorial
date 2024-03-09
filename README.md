# GraphQL-Tutorial
[Goで学ぶGraphQLサーバーサイド入門](https://zenn.dev/hsaki/books/golang-graphql)の内容に従う

## 起動方法
setup.shを起動しDBを作成する
```sh
./setup.sh
```

※ デフォルトではShebang(1行目の#!)にzshを指定しているため、必要に応じてを変更すること
※ Macではデフォルトでsqlite3が導入されている。WindowsやLinuxは必要に応じてインストールすること

## 変更事項
setup.shにてcreated_atの型はDATETIMEにしないとsqlboilerの自動生成コードがエラーになる
https://stackoverflow.com/questions/77796849/go-error-time-time-does-not-implement-driver-valuer-missing-method-value
```diff
CREATE TABLE IF NOT EXISTS repositories(\
	id TEXT PRIMARY KEY NOT NULL,\
	owner TEXT NOT NULL,\
	name TEXT NOT NULL,\
-	created_at TIMESTAMP NOT NULL DEFAULT (DATETIME('now','localtime')),\
+	created_at DATETIME NOT NULL DEFAULT (DATETIME('now','localtime')),\
	FOREIGN KEY (owner) REFERENCES users(id)\
);

```
