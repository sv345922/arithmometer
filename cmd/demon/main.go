package main

import (
	"arithmometer/calc"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

// запрашивает задачу у оркестратора
func getTask() (*calc.TaskContainer, error) {
	container := &calc.TaskContainer{}
	url := "http://127.0.0.1:8000/gettask"
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	// Если оркестратор не дал задачу возвращаем nil
	if len(body) == 0 {
		return nil, nil
	}
	// Анмаршалим body в контейнер
	err = json.Unmarshal(body, container)
	if err != nil {
		return nil, err
	}
	return container, nil
}

// Выполняет запросы оркестратору и вычисляет выражение
// TODO периодическое подтверждения работы
func main() {
	log.Print("calculator is runing")
	result := make(chan calc.Answer)
	for {
		container, err := getTask()
		if err != nil {
			log.Println("ошибка получения задачи", err)
			time.Sleep(5 * time.Second)
			continue
		}
		// Окрестратор не дал задание
		if container == nil {
			time.Sleep(5 * time.Second)
			continue
		}
		log.Println("Задача принята")
		// запускаем задачу в горутине
		go func(container *calc.TaskContainer) {
			res, err := calc.Calculate(container)
			result <- calc.Answer{
				Result: res,
				Err:    err,
			}
		}(container)
		answer := <-result
		// отправляем ответ, до тех пор пока он не будет принят
		for err != nil {
			err = SendAnswer(container.Id, answer)
			time.Sleep(time.Second)
		}
	}
}

// Отправляем ответ, если не отправилось, возвращаем ошибку
func SendAnswer(id string, answer calc.Answer) error {
	url := "http://127.0.0.1:8000/giveanswer"
	container := calc.AnswerContainer{
		Id:      id,
		AnswerN: answer,
	}
	data, _ := json.Marshal(container) //ошибку пропускаем
	r := bytes.NewReader(data)

	resp, err := http.Post(url, "application/json", r)
	if err != nil {
		return err
	}
	if resp.StatusCode == http.StatusOK {
		return nil
	}
	return fmt.Errorf("ошибка отправки ответа")
}
