package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
)

// Задача для вычисления
var expr = "-1+2-3/(4+5) * 6 -7"

func SendNewExpression(exprString string) (string, bool) {
	//errTotal := errors.New("ошибка отправки нового выражения")
	// Создать запрос
	url := "http://127.0.0.1:8000/newexpression"
	// Задать тайминги вычислений
	timing := &Timings{
		Plus:  6,
		Minus: 6,
		Mult:  6,
		Div:   6,
	}
	//timing = nil
	var exprst = NewExp{
		Expr:    exprString,
		Timings: timing,
	}
	data, _ := json.Marshal(exprst) //ошибку пропускаем
	r := bytes.NewReader(data)
	resp, err := http.Post(url, "application/json", r)
	if err != nil {
		return "", false
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", false
	}
	fmt.Printf("Постановка задачи\nStatus: %s\nBody:\n%s\n", resp.Status, string(body))
	id := string(body)
	fmt.Println("Задача отправлена")
	return id, true
}

// Получает результат вычислений
func GetResult(id string) (string, string, error) {
	errTotal := errors.New("ошибка получения результата")
	// Создать запрос
	url := "http://127.0.0.1:8000/getresult" + "?id=" + id
	resp, err := http.Get(url)
	if err != nil {
		return "", "", err
	}
	body, err := io.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return "", "", errTotal
	}
	fmt.Printf("Получение результата\nStatus: %s\nBody:\n%s\n", resp.Status, string(body))
	return resp.Status, string(body), nil
}
func main() {
	// отправка выражения
	//flag.Parse()
	if len(os.Args) > 1 {
		expr = os.Args[1]
	}
	id, _ := SendNewExpression(expr)
	fmt.Println()
	_ = id
	//
	//time.Sleep(10 * time.Second)
	//// получение ответа
	//_, _, err := GetResult(id)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//GetResult("250")
}
