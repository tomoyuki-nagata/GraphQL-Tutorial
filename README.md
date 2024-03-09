# GraphQL-Tutorial
[Goで学ぶGraphQLサーバーサイド入門](https://zenn.dev/hsaki/books/golang-graphql)の内容に従う

## 変更事項

### 3章 自作のスキーマを使ってGraphQLサーバーを作ろう
#### 今回のお題 - 簡略版GitHub API v4
最初に作成する段階では、ユーザー取得の時に認証のディレクティブがついているため、4章の実装後の動作確認で`directive isAuthenticated is not implemented`というエラーが発生する。
そのためこの段階ではディレクティブのアノテーションをつけないようにする
```diff
type Query {
  repository(
    name: String!
    owner: String!
  ): Repository

  user(
    name: String!
- ): User @isAuthenticated
+ ): User 
  node(
    id: ID!
  ): Node
}
```

### 4章 リゾルバの実装 - 基本編 
#### データを格納するDBの準備
setup.shにてcreated_atの型はDATETIMEにしないとsqlboilerの自動生成コードがエラーになる。
https://github.com/saki-engineering/graphql-sample/issues/1#issue-1725020469
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

## 起動方法
setup.shを起動しDBを作成する
```sh
./setup.sh
```

※ デフォルトではShebang(1行目の#!)にzshを指定しているため、必要に応じてを変更すること
※ Macではデフォルトでsqlite3が導入されている。WindowsやLinuxは必要に応じてインストールすること