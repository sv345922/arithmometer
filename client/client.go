package main

import (
	"arithmometer/internal/configs"
	"arithmometer/internal/entities"
	"arithmometer/internal/useCases/newExpression"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

// Задача для вычисления
var expr = "-1+2-3/(4+5) * 6 -7 * 8"

func SendNewExpression(exprString string) (string, bool) {
	// Создать запрос
	url := "http://127.0.0.1:" + configs.Port + "/newexpression"
	// Задать тайминги вычислений
	timing := &entities.Timings{
		Plus:  1,
		Minus: 1,
		Mult:  1,
		Div:   1,
	}
	//timing = nil
	var expression = newExpression.NewExpr{
		Expression: exprString,
		Timings:    timing,
	}
	data, _ := json.Marshal(expression) //ошибку пропускаем
	r := bytes.NewReader(data)
	resp, err := http.Post(url, "application/json", r)
	if err != nil {
		return "", false
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", false
	}
	fmt.Printf("Постановка задачи\nStatus: %s\tBody:\t%s\n", resp.Status, string(body))
	id := string(body)
	//fmt.Println("Задача отправлена")
	return id, true
}

// Получает результат вычислений
func GetResult(id string) (string, string, error) {
	errTotal := errors.New("ошибка получения результата")
	// Создать запрос
	url := "http://127.0.0.1:" + configs.Port + "/getresult" + "?id=" + id
	resp, err := http.Get(url)
	if err != nil {
		return "", "", err
	}
	body, err := io.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return "", "", errTotal
	}
	fmt.Printf("Получение результата\nStatus: %s\tBody:\t%s\n", resp.Status, string(body))
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
	//fmt.Println(id)

	time.Sleep(10 * time.Second)
	// получение ответа
	_, answer, err := GetResult(id)
	if err != nil {
		log.Fatal(err)
	}
	for !strings.Contains(answer, "результат выражения") {
		time.Sleep(3 * time.Second)
		_, answer, err = GetResult(id)
	}

}
