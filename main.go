package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

// UserResponse: JSONレスポンスの構造体を定義し、型安全性を確保
type UserResponse struct {
	Message string `json:"message"`
	ID      string `json:"id"`
	Detail  string `json:"detail"`
}

func main() {
	port := os.Getenv("PORT") // 環境変数からポートを取得
	if port == "" {
		port = "8080" // 環境変数が設定されていない場合のデフォルト値
	}
	serverAddr := fmt.Sprintf(":%s", port)

	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// ルート（/）のハンドラー
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)

		// --- 問題点1: w.Writeのエラーハンドリング ---
		_, err := w.Write([]byte("Hello! Go軽量APIサーバーが起動しました。"))
		if err != nil {
			log.Printf("Error writing root response: %v", err)
		}
	})

	// ユーザー関連のエンドポイント
	r.Route("/users", func(r chi.Router) {
		r.Get("/{userID}", getUser)
	})

	// サーバーの起動
	log.Printf("💡 Server listening on http://localhost%s", serverAddr) // log.Printfに統一
	log.Fatal(http.ListenAndServe(serverAddr, r))
}

// ユーザーIDを取得してレスポンスを返すハンドラー関数
func getUser(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "userID")
	response := UserResponse{
		Message: "ユーザー情報を取得しました。",
		ID:      userID,
		Detail:  "このAPIは軽量ルーターchiを使っています",
	}

	// JSONエンコードを試みる。成功したらヘッダーとステータスを送信。
	// エンコードエラーに備えて、一度バッファに書き出すか、エラーチェックを先に済ませる。
	// ここではシンプルに、エンコード結果を一度バイトスライスに変換するアプローチを提案。
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		log.Printf("Error marshalling JSON response for user %s: %v", userID, err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(jsonResponse) // エンコード済みのバイトスライスを書き込む
	if err != nil {
		log.Printf("Error writing JSON response for user %s: %v", userID, err)
	}
}
