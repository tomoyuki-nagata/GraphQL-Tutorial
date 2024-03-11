# GraphQL-Tutorial
[Goで学ぶGraphQLサーバーサイド入門](https://zenn.dev/hsaki/books/golang-graphql)の内容に従う

## 起動方法
1. setup.shを起動しDBを作成する
  ```sh
  ./setup.sh
  ```

※ デフォルトではShebang(1行目の#!)にzshを指定しているため、必要に応じてを変更すること

※ Macではデフォルトでsqlite3が導入されている。WindowsやLinuxは必要に応じてインストールすること

2. サーバーを起動する
  ```sh
  go run ./server.go
  ```

3. `http://localhost:8080`にアクセスし、下記のクエリを投げてみる
  ```gql
  query {
    node(id: "REPO_1") {
      ... on Repository {
        name
        issues(first: 2) {
          nodes {
            number
            author {
              name
            }
          }
        }
      }
    }
  }
  ```

## チュートリアル変更事項

### 3章 自作のスキーマを使ってGraphQLサーバーを作ろう

#### 今回のお題 - 簡略版GitHub API v4

最初に作成する段階では、ユーザー取得の時に認証のディレクティブがついているため、4章の実装後の動作確認で`directive isAuthenticated is not implemented`というエラーが発生する。

そのためこの段階ではディレクティブのアノテーションをつけないようにする。
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

issuesテーブルにauthorカラムがないため後続の6章の中でうまくいかない箇所が出てくるため下記のように修正する。
```diff
CREATE TABLE IF NOT EXISTS issues(\
	id TEXT PRIMARY KEY NOT NULL,\
	url TEXT NOT NULL,\
	title TEXT NOT NULL,\
 	title TEXT NOT NULL,\
 	closed INTEGER NOT NULL DEFAULT 0,\
 	number INTEGER NOT NULL,\
+	author TEXT NOT NULL,\
 	repository TEXT NOT NULL,\
 	CHECK (closed IN (0, 1)),\
-	FOREIGN KEY (repository) REFERENCES repositories(id)\
+	FOREIGN KEY (repository) REFERENCES repositories(id),\
+	FOREIGN KEY (author) REFERENCES users(id)\
 );
 
 CREATE TABLE IF NOT EXISTS projects(\
 	id TEXT PRIMARY KEY NOT NULL,\
 	title TEXT NOT NULL,\
 	url TEXT NOT NULL,\
+	number INTEGER NOT NULL,\
 	owner TEXT NOT NULL,\
 	FOREIGN KEY (owner) REFERENCES users(id)\
 );
-INSERT INTO issues(id, url, title, closed, number, repository) VALUES\
-	('ISSUE_1', 'http://example.com/repo1/issue/1', 'First Issue', 1, 1, 'REPO_1'),\
-	('ISSUE_2', 'http://example.com/repo1/issue/2', 'Second Issue', 0, 2, 'REPO_1'),\
-	('ISSUE_3', 'http://example.com/repo1/issue/3', 'Third Issue', 0, 3, 'REPO_1')\
+INSERT INTO issues(id, url, title, closed, number, author, repository) VALUES\
+	('ISSUE_1', 'http://example.com/repo1/issue/1', 'First Issue', 1, 1, 'U_1', 'REPO_1'),\
+	('ISSUE_2', 'http://example.com/repo1/issue/2', 'Second Issue', 0, 2, 'U_1', 'REPO_1'),\
+	('ISSUE_3', 'http://example.com/repo1/issue/3', 'Third Issue', 0, 3, 'U_1', 'REPO_1'),\
+	('ISSUE_4', 'http://example.com/repo1/issue/4', '', 0, 4, 'U_1', 'REPO_1'),\
+	('ISSUE_5', 'http://example.com/repo1/issue/5', '', 0, 5, 'U_1', 'REPO_1'),\
+	('ISSUE_6', 'http://example.com/repo1/issue/6', '', 0, 6, 'U_1', 'REPO_1'),\
+	('ISSUE_7', 'http://example.com/repo1/issue/7', '', 0, 7, 'U_1', 'REPO_1')\
 ;
 
-INSERT INTO projects(id, title, url, owner) VALUES\
-	('PJ_1', 'My Project', 'http://example.com/project/1', 'U_1')\
+INSERT INTO projects(id, title, url, number, owner) VALUES\
+	('PJ_1', 'My Project', 'http://example.com/project/1', 1, 'U_1'),\
+	('PJ_2', 'My Project 2', 'http://example.com/project/2', 2, 'U_1')\
 ;
 
```

### 5章 リゾルバの実装 - 応用編

#### リゾルバを分割する前の状況確認

下記のクエリを実行すると、issuesとpullRequestsがnullのためエラーになる。
```gql
query {
  repository(name: "repo1", owner: "hsaki"){
    id
    name
    createdAt
    owner {
      name
    }
    issue(number:1) {
      url
    }
    issues(first: 2) {
      nodes{
        title
      }
    }
    pullRequest(number:1) {
      baseRefName
      closed
      headRefName
    }
    pullRequests(last:2) {
      nodes{
        url
        number
      }
    }
  }
}
```

一旦Repositoryのスキーマ定義にてNOT NULL制約を外しておくとテキスト通りの挙動になる
```diff
type Repository implements Node {
  id: ID!
  owner: User!
  name: String!
  createdAt: DateTime!
  issue(
    number: Int!
  ): Issue
  issues(
    after: String
    before: String
    first: Int
    last: Int
-  ): IssueConnection!
+  ): IssueConnection
  pullRequest(
    number: Int!
  ): PullRequest
  pullRequests(
    after: String
    before: String
    first: Int
    last: Int
-  ): PullRequestConnection!
+  ): PullRequestConnection
  
}
```

### 10章 ディレクティブを利用した認証機構の追加

#### ディレクティブを利用したGraphQL層での認証機構の追加

3章で変更した認証ディレクティブを再度設定する
```diff
type Query {
  repository(
    name: String!
    owner: String!
  ): Repository

  user(
    name: String!
- ): User
+ ): User @isAuthenticated 
  node(
    id: ID!
  ): Node
}
```

