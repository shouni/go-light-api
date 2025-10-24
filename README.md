# Light API Server (Go + chi/v5)

[![Language](https://img.shields.io/badge/Language-Go-blue)](https://golang.org/)
[![Go Version](https://img.shields.io/github/go-mod/go-version/shouni/go-light-api)](https://golang.org/)
[![GitHub tag (latest by date)](https://img.shields.io/github/v/tag/shouni/go-light-api)](https://github.com/shouni/go-light-api/tags)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)


## 🌟 プロジェクト概要

このプロジェクトは、Go言語の標準ライブラリ（`net/http`）をベースに、軽量ルーター`go-chi/chi`を使用して構築されたシンプルなAPIサーバーのプロトタイプです。

**学習とベストプラクティス**の適用に焦点を当てており、環境変数からの設定読み込み、構造体を使用した安全なJSONレスポンス、および適切なエラーハンドリング（HTTP書き込みエラー、JSONエンコードエラー）が実装されています。

## 🚀 技術スタック

| 要素 | 技術 | 備考 |
| :--- | :--- | :--- |
| **言語** | **Go** (Golang) | 高速でシンプル、コンカレンシーに強い |
| **Webフレームワーク** | **net/http** (Standard Library) | サーバーの基盤 |
| **ルーティング** | **go-chi/chi/v5** | 軽量でモジュール化されたルーター |
| **ミドルウェア** | **chi/middleware** | ロギング (`Logger`)、パニックからの回復 (`Recoverer`) |
| **データ処理** | **encoding/json** | 型安全なJSONレスポンス生成 |
| **設定管理** | **os** | 環境変数からのポート読み込み (`$PORT`) |

-----

## 🛠️ セットアップと実行方法

### 1\. リポジトリのクローンと依存関係の取得

```bash
# リポジトリをクローン
git clone <your-repository-url>
cd go-light-api

# 依存関係をインストール（go mod init済みの前提）
go mod tidy
```

### 2\. ビルド

実行可能ファイルを作成します。今回は `bin/light-api` というパスに出力します。

```bash
go build -o bin/light-api
```

### 3\. サーバーの起動

ポート番号は環境変数 `PORT` から読み込まれます（未設定の場合はデフォルトで **`8080`** を使用）。

```bash
# デフォルトポート (8080) で起動
./bin/light-api
# 💡 Server listening on http://localhost:8080

# 特定のポート (例: 3000) を指定して起動
PORT=3000 ./bin/light-api
# 💡 Server listening on http://localhost:3000
```

-----

## 🌐 APIエンドポイント

現在定義されているエンドポイントは以下の通りです。

| HTTPメソッド | パス | 説明 | 備考 |
| :--- | :--- | :--- | :--- |
| **GET** | `/` | サーバーの稼働確認 | "Hello\!" メッセージを返します |
| **GET** | `/users/{userID}` | 指定されたユーザーIDの情報を取得 | `userID` はパスパラメータとして抽出されます |

### 動作確認例 (curl)

#### 1\. ルートアクセス

```bash
curl http://localhost:8080/
# => Hello! Go軽量APIサーバーが起動しました。
```

#### 2\. ユーザー情報取得

```bash
curl http://localhost:8080/users/alice_dev
```

**レスポンス (JSON)**

```json
{
  "message": "ユーザー情報を取得しました。",
  "id": "alice_dev",
  "detail": "このAPIは軽量ルーターchiを使っています"
}
```

-----

## ⚙️ ベストプラクティスと堅牢性への配慮

このプロジェクトには、お客様が過去に開発された「Action Perfect Get On Go」と同様に、堅牢性向上のための学習要素が含まれています。

* **安全なJSON処理**: `fmt.Sprintf` ではなく `encoding/json` を使用し、構造体ベースで型安全なJSONを生成。
* **適切なエラーハンドリング**: HTTPレスポンス書き込み (`w.Write`) やJSONエンコード時のエラーをチェックし、`log.Printf` で記録。
* **設定の柔軟性**: サービス設定値（ポート番号）をハードコードせず、**環境変数**から読み込む仕組みを導入。
* **ログの一貫性**: サーバー起動メッセージを含むすべての出力を `log` パッケージに統一。

### 📜 ライセンス (License)

このプロジェクトは [MIT License](https://opensource.org/licenses/MIT) の下で公開されています。
