package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	_ "github.com/mattn/go-sqlite3"

	"go-light-api/internal/model"
	"go-light-api/internal/repository"
)

// グローバル変数としてUserRepositoryのインターフェースインスタンスを保持
// これにより、ハンドラー関数でリポジトリインスタンスを直接利用できる
var userRepo repository.UserRepository

func main() {
	// データベースの初期化
	db := initDB()
	defer db.Close()

	// リポジトリの初期化
	userRepo = repository.NewUserRepository(db)
	if err := userRepo.InitTable(); err != nil {
		log.Fatalf("Error initializing users table: %v", err)
	}
	log.Println("✅ Users table ready.")

	// --- サーバー設定 ---
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
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
// DB接続初期化関数
// ------------------------------------

func initDB() *sql.DB {
	db, err := sql.Open("sqlite3", "./users.db")
	if err != nil {
		log.Fatalf("Error opening database: %v", err)
	}

	if err = db.Ping(); err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}
	log.Println("✅ Successfully connected to SQLite database.")
	return db
}

// ------------------------------------
// APIハンドラー関数 (リポジトリを利用)
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
	var u model.User
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if u.ID == "" || u.Name == "" {
		http.Error(w, "ID and Name are required", http.StatusBadRequest)
		return
	}

	// リポジトリメソッドを呼び出す（SQL文の詳細は関知しない）
	if err := userRepo.Create(&u); err != nil {
		log.Printf("Error creating user %s via repository: %v", u.ID, err)
		http.Error(w, "Failed to create user (ID likely exists or DB error)", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated) // 201 Created

	response := model.UserResponse{Message: "ユーザーが正常に登録されました。", User: &u}
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Error encoding JSON response after create: %v", err)
	}
}

// getUserHandler: ユーザー情報取得 (GET /users/{userID})
func getUserHandler(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "userID")

	// リポジトリメソッドを呼び出す
	u, err := userRepo.FindByID(userID)

	if err != nil {
		// リポジトリから返されたDBエラーをここでロギングする (責務の分離)
		log.Printf("Handler Error querying user %s: %v", userID, err)
		http.Error(w, "Internal database error", http.StatusInternalServerError)
		return
	}

	if u == nil {
		// データが見つからなかった場合 (リポジトリから nil, nil が返された)
		response := model.UserResponse{Message: fmt.Sprintf("ユーザーID '%s' は見つかりませんでした。", userID)}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound) // 404 Not Found
		if jsonErr := json.NewEncoder(w).Encode(response); jsonErr != nil {
			log.Printf("Error encoding 404 response: %v", jsonErr)
		}
		return
	}

	// 成功 (データが見つかった場合)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK) // 200 OK

	response := model.UserResponse{Message: "ユーザー情報を取得しました。", User: u}
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Error encoding JSON response for user %s: %v", userID, err)
	}
}
