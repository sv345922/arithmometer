package server

import (
	"arithmometer/orchestr/inter/tasker"
	"context"
	"fmt"
	"log"
	"net/http"
)

func stateContext(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ws, err := tasker.RunTasker()
		fmt.Printf("%v\n", ws)
		if err != nil {
			log.Println("main: %v", err)
		}
		ctx := context.WithValue(context.Background(), "ws", ws)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
