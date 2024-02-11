package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
)

func SendHTTPRequest(url string) (string, string, error) {
	errTotal := errors.New("Something went wrong...")
	// Создать запрос

	exprString := "-1+2-3/4+5"
	timing := &Timings{
		Plus:  5,
		Minus: 5,
		Mult:  5,
		Div:   5,
	}
	//timing = nil
	var expr = NewExp{
		Expr:    exprString,
		Timings: timing,
	}
	data, _ := json.Marshal(expr) //ошибку пропускаем
	r := bytes.NewReader(data)
	resp, err := http.Post("http://127.0.0.1:8000/newexpression", "application/json", r)
	if err != nil {
		return "", "", err
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", "", errTotal
	}
	return resp.Status, string(body), nil
}

func main() {
	url := "http://127.0.0.1:8000/newexpression?expr=-1+2-3/4+5"
	respStatus, respBody, err := SendHTTPRequest(url)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Status: %s\nBody:\n%s\n", respStatus, respBody)
}
