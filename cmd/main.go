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

// グローバル変数 userRepo は削除されました。

func main() {
	// ------------------------------------
	// 1. データベースの初期化とエラーハンドリング
	// ------------------------------------
	db, err := initDB() // エラーを受け取る形に変更
	if err != nil {
		log.Fatalf("Database initialization failed: %v", err) // ここで致命的なエラーとして処理
	}
	defer db.Close()

	// ------------------------------------
	// 2. リポジトリの初期化 (ローカル変数として初期化)
	// ------------------------------------
	userRepo := repository.NewUserRepository(db)
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
	// healthCheckHandler は依存性がないため、直接登録
	r.Get("/", healthCheckHandler)

	r.Route("/users", func(r chi.Router) {
		// 依存性注入: ファクトリ関数を通じてリポジトリを渡す
		r.Post("/", makeCreateUserHandler(userRepo))
		r.Get("/{userID}", makeGetUserHandler(userRepo))
	})

	// --- サーバー起動 ---
	log.Printf("💡 Server listening on http://localhost%s", serverAddr)
	log.Fatal(http.ListenAndServe(serverAddr, r))
}

// ------------------------------------
// DB接続初期化関数 (エラーを返すように変更)
// ------------------------------------

func initDB() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "./users.db")
	if err != nil {
		return nil, fmt.Errorf("error opening database: %w", err)
	}

	if err = db.Ping(); err != nil {
		db.Close() // 接続失敗時は開いた接続を確実に閉じる
		return nil, fmt.Errorf("error connecting to database: %w", err)
	}
	log.Println("✅ Successfully connected to SQLite database.")
	return db, nil
}

// ------------------------------------
// ヘルパー関数
// ------------------------------------

// respondWithJSON: 共通のJSONレスポンス送信ロジックをカプセル化
func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	if payload == nil {
		return // payload が nil の場合はエンコードしない
	}

	if err := json.NewEncoder(w).Encode(payload); err != nil {
		log.Printf("Error encoding JSON response: %v", err)
		// 既にヘッダーが送信されているため、ロギングに留める
	}
}

// ------------------------------------
// APIハンドラー生成関数 (依存性注入)
// ------------------------------------

// healthCheckHandler: 稼働確認用 (依存性なし)
func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	// ここは respondWithJSON よりもシンプルなので w.Write を維持
	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte("Hello! Go軽量APIサーバーが起動しました。DB接続もOKです。"))
	if err != nil {
		log.Printf("Error writing health check response: %v", err)
	}
}

// makeCreateUserHandler: ユーザー登録ハンドラーを生成
func makeCreateUserHandler(repo repository.UserRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var u model.User
		if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		if u.ID == "" || u.Name == "" {
			http.Error(w, "ID and Name are required", http.StatusBadRequest)
			return
		}

		if err := repo.Create(&u); err != nil { // 依存オブジェクト(repo)を利用
			log.Printf("Error creating user %s via repository: %v", u.ID, err)
			http.Error(w, "Failed to create user (ID likely exists or DB error)", http.StatusInternalServerError)
			return
		}

		respondWithJSON(w, http.StatusCreated, model.UserResponse{Message: "ユーザーが正常に登録されました。", User: &u})
	}
}

// makeGetUserHandler: ユーザー情報取得ハンドラーを生成
func makeGetUserHandler(repo repository.UserRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := chi.URLParam(r, "userID")

		u, err := repo.FindByID(userID) // 依存オブジェクト(repo)を利用

		if err != nil {
			// リポジトリ層からラップされたエラーが返される (コンテキスト情報を含む)
			log.Printf("Handler Error querying user %s: %v", userID, err)
			http.Error(w, "Internal database error", http.StatusInternalServerError)
			return
		}

		if u == nil {
			// データが見つからなかった場合
			response := model.UserResponse{Message: fmt.Sprintf("ユーザーID '%s' は見つかりませんでした。", userID)}
			respondWithJSON(w, http.StatusNotFound, response)
			return
		}

		// 成功 (データが見つかった場合)
		respondWithJSON(w, http.StatusOK, model.UserResponse{Message: "ユーザー情報を取得しました。", User: u})
	}
}
