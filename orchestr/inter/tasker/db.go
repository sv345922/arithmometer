package tasker

import (
	"arithmometer/orchestr/parsing"
	"encoding/json"
	"log"
	"os"
	"sync"
)

type DataBase struct {
	// список выражений (с таймингами)
	Expressions *Expressions `json:"expressions"`
	Tasks       *Tasks       `json:"tasks"`
	//Timings     *Timings     `json:"timings"`
	mu sync.Mutex
}

type additiveJSON interface {
	*Expression | []*parsing.Node | *parsing.Node | []*parsing.Symbol | Timings | map[string]*Expression | DataBase
}

// SafeJSON Сохраняет структуру в базе данных, в папке db
func SafeJSON[T additiveJSON](name string, expr T) error {
	jsonBytes, err := json.Marshal(expr)
	if err != nil {
		log.Println(err)
		return err
	}
	wd, err := os.Getwd()
	if err != nil {
		log.Println(err)
		return err
	}
	path := wd + "/orchestr/db/" + name + ".json"
	err = os.WriteFile(path, jsonBytes, 0666)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

// Загружает структуру из db и возвращает её
func LoadJSON[T additiveJSON](name string) (*T, error) {
	var result T
	wd, err := os.Getwd()
	if err != nil {
		log.Println(err)
		return nil, err
	}
	path := wd + "/orchestr/db/" + name + ".json"
	data, err := os.ReadFile(path)
	if err != nil {
		log.Println("ошибка открытия json", err)
		return nil, err
	}
	err = json.Unmarshal(data, result)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return &result, nil
}
