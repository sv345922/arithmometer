package tasker

import (
	"encoding/json"
	"log"
	"os"
)

type DataBase struct {
	// список выражений (с таймингами)
	Expressions *Expressions `json:"expressions"`
	Tasks       *Tasks       `json:"tasks"`
	//Timings     *Timings     `json:"timings"`
	//mu sync.Mutex
}

type additiveJSON interface {
	DataBase
}

func NewDB() *DataBase {
	result := DataBase{
		Expressions: NewExpressions(),
		Tasks:       NewTasks(),
		//mu:          sync.Mutex{},
	}
	return &result
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
	path := wd + "\\orchestr\\db\\" + name + ".json"
	err = os.WriteFile(path, jsonBytes, 0666)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

// Загружает структуру из db и возвращает её
func LoadDB() (*DataBase, error) {
	wd, err := os.Getwd()
	if err != nil {
		log.Println(err)
		return nil, err
	}
	path := wd + "\\orchestr\\db\\" + "db.json"
	data, err := os.ReadFile(path)
	if err != nil {
		log.Println("ошибка открытия json", err)
		return nil, err
	}
	var result DataBase
	err = json.Unmarshal(data, &result)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return &result, nil
}
