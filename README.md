# Light API Server (Go + chi/v5)

[![Language](https://img.shields.io/badge/Language-Go-blue)](https://golang.org/)
[![Go Version](https://img.shields.io/github/go-mod/go-version/shouni/go-light-api)](https://golang.org/)
[![GitHub tag (latest by date)](https://img.shields.io/github/v/tag/shouni/go-light-api)](https://github.com/shouni/go-light-api/tags)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)


# Light API Server (Go + chi/v5 + SQLite3)

## 🌟 プロジェクト概要

このプロジェクトは、Go言語の標準ライブラリ（`net/http`）と軽量ルーター\*\*`go-chi/chi`\*\*を使用して構築された、**データ永続化機能を持つ**シンプルなAPIサーバーです。

データ層にはサーバーレスで動作する**SQLite3**を採用しており、外部データベースサーバーなしでデータの\*\*CRUD操作（作成・読み取り）\*\*を実現しています。**学習とベストプラクティス**（安全なJSON処理、環境設定、堅牢なエラーハンドリング）の適用に焦点を当てています。

-----

## 🚀 技術スタック

| 要素 | 技術 | 備考 |
| :--- | :--- | :--- |
| **言語** | **Go** (Golang) | 高速でシンプル、コンカレンシーに強い |
| **データベース** | **SQLite3** | 外部サーバー不要の軽量ファイルベースDB |
| **DBドライバ** | **mattn/go-sqlite3** | GoからSQLiteを操作するためのドライバ |
| **ルーティング** | **go-chi/chi/v5** | 軽量でモジュール化されたルーター |
| **データ処理** | **encoding/json** | 型安全なJSONレスポンス生成 |
| **設定管理** | **os** | 環境変数からのポート読み込み (`$PORT`) |

-----

## 🛠️ セットアップと実行方法

### 1\. 依存関係のインストール

プロジェクトに必要なGoモジュールをインストールします。

```bash
# プロジェクトフォルダへ移動
cd go-light-api

# ルーターとSQLiteドライバをインストール
go get github.com/go-chi/chi/v5
go get github.com/mattn/go-sqlite3

# 依存関係を整理
go mod tidy
```

### 2\. ビルドと起動

実行可能ファイルをビルドし、サーバーを起動します。起動時に\*\*`users.db`\*\*というSQLiteファイルが自動で作成され、`users`テーブルが初期化されます。

```bash
# 実行可能ファイルを生成
go build -o bin/light-api

# サーバーを起動 (デフォルトポート: 8080)
./bin/light-api
# 💡 Server listening on http://localhost:8080...
```

-----

## 🌐 APIエンドポイント (CRUD)

現在、ユーザーリソースに対する**作成 (Create)** と **読み取り (Read)** が実装されています。

| HTTPメソッド | パス | 説明 | 備考 |
| :--- | :--- | :--- | :--- |
| **GET** | `/` | サーバーの稼働状態確認 | DB接続成功も確認します |
| **POST** | `/users` | 新しいユーザーを登録 (**Create**) | Request Bodyに `{id, name, email}` が必要 |
| **GET** | `/users/{userID}` | 指定されたユーザーIDの情報をDBから取得 (**Read**) | ユーザーが見つからない場合は **404 Not Found** を返します |

### 動作確認例 (curl)

#### 1\. ユーザー登録 (POST)

```bash
curl -X POST http://localhost:8080/users \
     -H "Content-Type: application/json" \
     -d '{"id": "john.doe", "name": "John Doe", "email": "john@example.com"}'
```

**成功レスポンス (201 Created)**:

```json
{"message":"ユーザーが正常に登録されました。","user":{"id":"john.doe","name":"John Doe","email":"john@example.com"}}
```

#### 2\. ユーザー情報取得 (GET)

```bash
curl http://localhost:8080/users/john.doe
```

**成功レスポンス (200 OK)**:

```json
{"message":"ユーザー情報を取得しました。","user":{"id":"john.doe","name":"John Doe","email":"john@example.com"}}
```

-----

## ⚙️ ベストプラクティスと堅牢性への配慮

このプロジェクトは、Go言語における堅牢なアプリケーション開発の学習を目的としています。

* **データ永続化**: SQLite3を採用し、サーバーを再起動してもデータが失われない設計。
* **堅牢なDB操作**: `database/sql`と**プリペアドステートメント**を使用し、**SQLインジェクション脆弱性**に対応。
* **安全なJSON処理**: `fmt.Sprintf` ではなく `encoding/json` を使用し、構造体ベースで型安全なJSONレスポンスを生成。
* **適切なエラーハンドリング**: DBクエリ、JSONエンコード/デコード、HTTP書き込みなど、各段階でエラーをチェックし、適切にロギング。
* **設定の柔軟性**: サービス設定値（ポート番号）をハードコードせず、**環境変数** (`PORT`) から読み込む仕組みを導入。

### 📜 ライセンス (License)

このプロジェクトは [MIT License](https://opensource.org/licenses/MIT) の下で公開されています。
