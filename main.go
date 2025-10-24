package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	// SQLiteドライバをインポート (_ で登録のみ行う)
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	_ "github.com/mattn/go-sqlite3"
)

// グローバル変数としてデータベース接続を保持
var db *sql.DB

// User: ユーザーデータ構造体（DBとのやり取り、リクエストボディのデコードに使用）
type User struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

// UserResponse: ユーザー情報取得時のレスポンス構造体
type UserResponse struct {
	Message string `json:"message"`
	User    *User  `json:"user,omitempty"` // データがない場合は省略
}

func main() {
	// データベースの初期化
	initDB()
	defer db.Close() // サーバー終了時にDB接続を閉じる

	// --- サーバー設定 ---
	port := os.Getenv("PORT") // 環境変数からポートを取得
	if port == "" {
		port = "8080" // 環境変数が設定されていない場合のデフォルト値
	}
	serverAddr := fmt.Sprintf(":%s", port)

	// --- ルーター設定 ---
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// --- エンドポイント定義 ---
	r.Get("/", healthCheckHandler)

	r.Route("/users", func(r chi.Router) {
		r.Post("/", createUserHandler)     // ユーザー登録
		r.Get("/{userID}", getUserHandler) // ユーザー情報取得
	})

	// --- サーバー起動 ---
	log.Printf("💡 Server listening on http://localhost%s", serverAddr)
	log.Fatal(http.ListenAndServe(serverAddr, r))
}

// ------------------------------------
// DB初期化関数
// ------------------------------------

func initDB() {
	var err error
	// データベースファイルを開く（存在しない場合は新規作成）
	db, err = sql.Open("sqlite3", "./users.db")
	if err != nil {
		log.Fatalf("Error opening database: %v", err)
	}

	// 接続が確立されているか確認
	if err = db.Ping(); err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}
	log.Println("✅ Successfully connected to SQLite database: ./users.db")

	// テーブルを作成するSQL文
	createTableSQL := `
	CREATE TABLE IF NOT EXISTS users (
		id TEXT PRIMARY KEY,
		name TEXT NOT NULL,
		email TEXT
	);`

	// 実行
	_, err = db.Exec(createTableSQL)
	if err != nil {
		log.Fatalf("Error creating users table: %v", err)
	}
	log.Println("✅ Users table ready.")
}

// ------------------------------------
// APIハンドラー関数
// ------------------------------------

// healthCheckHandler: 稼働確認用
func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte("Hello! Go軽量APIサーバーが起動しました。DB接続もOKです。"))
	if err != nil {
		log.Printf("Error writing health check response: %v", err)
	}
}

// createUserHandler: ユーザー登録 (POST /users)
func createUserHandler(w http.ResponseWriter, r *http.Request) {
	var u User
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// 必須フィールドの簡易チェック
	if u.ID == "" || u.Name == "" {
		http.Error(w, "ID and Name are required", http.StatusBadRequest)
		return
	}

	// データをDBに挿入（プリペアドステートメントでSQLインジェクション対策）
	_, err := db.Exec("INSERT INTO users(id, name, email) VALUES(?, ?, ?)", u.ID, u.Name, u.Email)
	if err != nil {
		log.Printf("Error inserting user %s: %v", u.ID, err)
		http.Error(w, "Failed to create user (ID likely exists)", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated) // 201 Created

	// 登録されたユーザー情報をJSONで返す
	response := UserResponse{Message: "ユーザーが正常に登録されました。", User: &u}
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Error encoding JSON response after create: %v", err)
	}
}

// getUserHandler: ユーザー情報取得 (GET /users/{userID})
func getUserHandler(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "userID")

	var u User

	// DBからデータを取得
	row := db.QueryRow("SELECT id, name, email FROM users WHERE id = ?", userID)

	// スキャン（取得したデータをGoの構造体にマッピング）
	err := row.Scan(&u.ID, &u.Name, &u.Email)

	if err == sql.ErrNoRows {
		// 該当データがない場合
		response := UserResponse{Message: fmt.Sprintf("ユーザーID '%s' は見つかりませんでした。", userID)}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound) // 404 Not Found
		if jsonErr := json.NewEncoder(w).Encode(response); jsonErr != nil {
			log.Printf("Error encoding 404 response: %v", jsonErr)
		}
		return
	} else if err != nil {
		// その他のDBエラー
		log.Printf("Error querying user %s: %v", userID, err)
		http.Error(w, "Internal database error", http.StatusInternalServerError) // 500 Internal Server Error
		return
	}

	// 成功 (データが見つかった場合)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK) // 200 OK

	response := UserResponse{Message: "ユーザー情報を取得しました。", User: &u}
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Error encoding JSON response for user %s: %v", userID, err)
	}
}
