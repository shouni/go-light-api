package main

import (
	"fmt"
	"log"
	"net/http"

	// 軽量ルーターとして go-chi/chi を使用
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	// 1. ルーターの作成（標準ライブラリの http.Handler インターフェースを実装）
	r := chi.NewRouter()

	// 2. ミドルウェアの設定（ロギングやリカバリーなどの便利機能）
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// 3. ルーティングの定義
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		// ルート（/）にアクセスがあった場合の処理
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Hello! Go軽量APIサーバーが起動しました。"))
	})

	// ユーザー関連のエンドポイント
	r.Route("/users", func(r chi.Router) {
		// GET /users/{userID} の定義
		r.Get("/{userID}", getUser)
	})

	// 4. サーバーの起動
	port := ":8080"
	fmt.Printf("💡 Server listening on http://localhost%s\n", port)
	// chi.Mux (r) は http.Handler のため、標準ライブラリの ListenAndServe に渡せます
	log.Fatal(http.ListenAndServe(port, r))
}

// ユーザーIDを取得してレスポンスを返すハンドラー関数
func getUser(w http.ResponseWriter, r *http.Request) {
	// chi.URLParam を使って、URLのパスから動的なパラメータ（userID）を簡単に取得
	userID := chi.URLParam(r, "userID")

	// レスポンスの書き込み
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// 今回はDB接続などがないため、静的なJSONを返す
	response := fmt.Sprintf(`{"message": "ユーザー情報を取得しました。", "id": "%s", "detail": "このAPIは軽量ルーターchiを使っています"}`, userID)
	w.Write([]byte(response))
}
