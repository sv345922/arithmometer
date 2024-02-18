package server

import (
	"arithmometer/orchestr/inter/tasker"
	"context"
	"log"
	"net/http"
)

func stateContext(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ws, err := tasker.RunTasker()
		defer ws.Save()
		if err != nil {
			log.Println("mw: ", err)
		}
		ctx := context.WithValue(context.Background(), "ws", ws)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
