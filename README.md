# Light API

[![Language](https://img.shields.io/badge/Language-Go-blue)](https://golang.org/)
[![Go Version](https://img.shields.io/github/go-mod/go-version/shouni/go-light-api)](https://golang.org/)
[![GitHub tag (latest by date)](https://img.shields.io/github/v/tag/shouni/go-light-api)](https://github.com/shouni/go-light-api/tags)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

# Light API Server (Go + chi/v5 + SQLite3)

## 🌟 プロジェクト概要

このプロジェクトは、Go言語の標準ライブラリ（`net/http`）と軽量ルーター\*\*`go-chi/chi`\*\*を使用して構築された、**データ永続化機能を持つ**シンプルなAPIサーバーです。

データ層にはサーバーレスで動作する**SQLite3**を採用しており、外部データベースサーバーなしでデータの**CRUD操作（作成・読み取り）を実現しています。設計段階でリポジトリパターン**、**依存性注入 (DI)**、**Goの慣習（`internal`パッケージ）を徹底的に適用し、高い保守性**と**テスト容易性**を確保しています。

-----

## 🚀 技術スタック

| 要素 | 技術 | 備考 |
| :--- | :--- | :--- |
| **言語** | **Go** (Golang) | 高速でシンプル、コンカレンシーに強い |
| **データベース** | **SQLite3** | 外部サーバー不要の軽量ファイルベースDB |
| **DBドライバ** | **mattn/go-sqlite3** | GoからSQLiteを操作するためのドライバ |
| **ルーティング** | **go-chi/chi/v5** | 軽量でモジュール化されたルーター |
| **設計パターン** | **リポジトリパターン** / **依存性注入 (DI)** | 責務分離とテスト容易性の確保 |

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

リファクタリングにより、実行可能ファイルは `cmd/main.go` からビルドされます。起動時に\*\*`users.db`\*\*というSQLiteファイルが自動で作成され、`users`テーブルが初期化されます。

```bash
# 実行可能ファイルを生成
go build -o bin/light-api ./cmd

# サーバーを起動 (デフォルトポート: 8080)
./bin/light-api
# 💡 Server listening on http://localhost:8080...
```

-----

## ⚙️ 内部設計とアーキテクチャ

このプロジェクトは、単なるプロトタイプではなく、Goのベストプラクティスに従ったアーキテクチャを採用しています。

### ファイル構成 (`internal` パッケージの利用)

| パス | 責務 |
| :--- | :--- |
| `cmd/main.go` | エントリポイント、ルーティング設定、依存性の組み立て |
| `internal/model/` | データの構造体定義 (`User`, `HealthCheckResponse`) |
| `internal/repository/` | DB操作のロジック層。**すべてのSQL文をここに集約** |

### 依存性注入 (DI)

グローバル変数を使用せず、`main`関数内で初期化された`UserRepository`インスタンスを**ファクトリ関数**を通じて各ハンドラーに渡し、**テスト容易性**を確保しています。

-----

## 🌐 APIエンドポイント (CRUD)

現在、ユーザーリソースに対する**作成 (Create)** と **読み取り (Read)** が実装されています。

| HTTPメソッド | パス | 説明 | 備考 |
| :--- | :--- | :--- | :--- |
| **GET** | `/` | サーバーの稼働状態確認 | レスポンス形式はJSONで統一 |
| **POST** | `/users` | 新しいユーザーを登録 (**Create**) | ID重複時は **409 Conflict** を返します |
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

#### 2\. ユーザー登録失敗 (ID重複時)

```bash
curl -X POST http://localhost:8080/users \
     -H "Content-Type: application/json" \
     -d '{"id": "john.doe", "name": "New Name", "email": "new@example.com"}'
```

**失敗レスポンス (409 Conflict)**:

```json
{"message":"User with ID 'john.doe' already exists"}
```

-----

## 🌟 ベストプラクティスと堅牢性への配慮

| 項目 | 詳細 |
| :--- | :--- |
| **API一貫性** | 正常系・エラー系問わず、すべてのレスポンスを**JSON形式**で統一。 |
| **セマンティクス** | ID重複エラーを正確に検出し、HTTPステータスコード\*\*`409 Conflict`\*\*で通知。 |
| **エラーチェーン** | リポジトリ層でエラーにコンテキストを追加し、`fmt.Errorf` の **`%w`** でラップ。 |
| **DIの適用** | `main`関数からハンドラーへ依存オブジェクトを注入し、**単体テストを容易**に。 |
| **セキュリティ** | **プリペアドステートメント**を使用し、SQLインジェクション脆弱性に対応。 |

## 📜 ライセンス (License)

このプロジェクトは [MIT License](https://opensource.org/licenses/MIT) の下で公開されています。

