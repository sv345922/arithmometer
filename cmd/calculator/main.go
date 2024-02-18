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

const URL = "http://127.0.0.1:8000/"

// запрашивает задачу у оркестратора
func getTask(calcId string) (*calc.TaskContainer, error) {
	container := &calc.TaskContainer{}
	url := URL + "gettask?id=" + calcId
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	body, err := io.ReadAll(resp.Body)
	defer resp.Body.Close()
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

// Отправляем ответ, если не отправилось, возвращаем ошибку
func SendAnswer(id string, answer calc.Answer) error {
	url := URL + "giveanswer"
	container := calc.AnswerContainer{
		Id:      id,
		AnswerN: answer,
	}
	data, _ := json.Marshal(container) //ошибку пропускаем
	r := bytes.NewReader(data)

	resp, err := http.Post(url, "application/json", r)
	if err != nil {
		fmt.Printf("обшибка отправки запроса POST", err) //TODO delete
		return err
	}
	if resp.StatusCode == http.StatusOK {
		return nil
	}
	return fmt.Errorf("ошибка отправки ответа")
}

// Выполняет запросы оркестратору и вычисляет выражение
// TODO периодическое подтверждения работы
func main() {
	log.Print("calculator is runing")
	calcId := int(time.Now().UnixNano())
	result := make(chan calc.Answer)
	for {
		container, err := getTask(fmt.Sprintf("%d", calcId))
		if err != nil {
			log.Println("ошибка получения задачи", err)
			time.Sleep(5 * time.Second)
			continue
		}
		// Окрестратор не дал задание
		if container == nil {
			log.Println("нет задач для вычислителя")
			time.Sleep(5 * time.Second)
			continue
		}
		log.Println("задача принята")
		// запускаем задачу в горутине
		go func(container *calc.TaskContainer) {
			res, err := calc.Calculate(container)
			result <- calc.Answer{
				Result: res,
				Err:    err,
			}
		}(container)
		answer := <-result
		log.Printf("задача %v выполнена, результат %f\n", container.TaskN, answer.Result)
		// отправляем ответ, до тех пор пока он не будет принят
		for {
			err = SendAnswer(container.Id, answer)
			if err == nil {
				break
			}
		}
		log.Println("отправлен ответ оркестратору")
		time.Sleep(time.Second)
	}
}
