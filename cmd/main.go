package main

import (
	"database/sql"
	"encoding/json"
	"errors" // errorsパッケージをインポート
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

func main() {
	// データベースの初期化
	db, err := initDB()
	if err != nil {
		log.Fatalf("Database initialization failed: %v", err)
	}
	defer db.Close()

	// リポジトリの初期化
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
	r.Get("/", healthCheckHandler)

	r.Route("/users", func(r chi.Router) {
		r.Post("/", makeCreateUserHandler(userRepo))
		r.Get("/{userID}", makeGetUserHandler(userRepo))
	})

	// --- サーバー起動 ---
	log.Printf("💡 Server listening on http://localhost%s", serverAddr)
	log.Fatal(http.ListenAndServe(serverAddr, r))
}

// ------------------------------------
// DB接続初期化関数 (変更なし)
// ------------------------------------

func initDB() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "./users.db")
	if err != nil {
		return nil, fmt.Errorf("error opening database: %w", err)
	}

	if err = db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("error connecting to database: %w", err)
	}
	log.Println("✅ Successfully connected to SQLite database.")
	return db, nil
}

// ------------------------------------
// ヘルパー関数 (変更なし)
// ------------------------------------

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	if payload == nil {
		return
	}

	if err := json.NewEncoder(w).Encode(payload); err != nil {
		log.Printf("Error encoding JSON response: %v", err)
	}
}

// ------------------------------------
// APIハンドラー生成関数 (依存性注入)
// ------------------------------------

// healthCheckHandler (変更なし)
func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	resp := model.HealthCheckResponse{
		Status:  "ok",
		Message: "Hello! Go軽量APIサーバーが起動しました。DB接続もOKです。",
	}
	respondWithJSON(w, http.StatusOK, resp)
}

// makeCreateUserHandler: ユーザー登録ハンドラーを生成 (エラーハンドリングを修正)
func makeCreateUserHandler(repo repository.UserRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var u model.User

		// 1. JSONデコードエラー
		if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
			respondWithJSON(w, http.StatusBadRequest, model.UserResponse{Message: "Invalid request body"})
			return
		}

		// 2. バリデーションエラー
		if u.ID == "" || u.Name == "" {
			respondWithJSON(w, http.StatusBadRequest, model.UserResponse{Message: "ID and Name are required"})
			return
		}

		// 3. DB作成エラー
		if err := repo.Create(&u); err != nil {
			log.Printf("Error creating user %s via repository: %v", u.ID, err)

			// ErrDuplicateEntry をチェックし、409 Conflict を返す
			if errors.Is(err, repository.ErrDuplicateEntry) {
				respondWithJSON(w, http.StatusConflict, model.UserResponse{Message: fmt.Sprintf("User with ID '%s' already exists", u.ID)}) // 409 Conflict
				return
			}

			// その他のエラーは 500 Internal Server Error
			respondWithJSON(w, http.StatusInternalServerError, model.UserResponse{Message: "Failed to create user due to internal error"})
			return
		}

		// 成功レスポンス
		respondWithJSON(w, http.StatusCreated, model.UserResponse{Message: "ユーザーが正常に登録されました。", User: &u})
	}
}

// makeGetUserHandler: ユーザー情報取得ハンドラーを生成 (変更なし)
func makeGetUserHandler(repo repository.UserRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := chi.URLParam(r, "userID")

		u, err := repo.FindByID(userID)

		// 1. DBクエリエラー
		if err != nil {
			log.Printf("Handler Error querying user %s: %v", userID, err)
			respondWithJSON(w, http.StatusInternalServerError, model.UserResponse{Message: "Internal database error"})
			return
		}

		// 2. データNotFound
		if u == nil {
			response := model.UserResponse{Message: fmt.Sprintf("ユーザーID '%s' は見つかりませんでした。", userID)}
			respondWithJSON(w, http.StatusNotFound, response)
			return
		}

		// 成功レスポンス
		respondWithJSON(w, http.StatusOK, model.UserResponse{Message: "ユーザー情報を取得しました。", User: u})
	}
}
