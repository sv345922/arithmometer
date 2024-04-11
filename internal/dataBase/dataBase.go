package dataBase

import (
	"arithmometer/internal/configs"
	"arithmometer/internal/entities"
	"arithmometer/internal/taskQueue"
	"encoding/json"
	"log"
	"os"
)

type DataBase struct {
	// список выражений (с таймингами)
	Queue       *taskQueue.Queue                `json:"tasks"`
	Expressions map[uint64]*entities.Expression `json:"expressions"` // []Expression
	AllNodes    map[uint64]*entities.Node       `json:"allNodes"`
	Timings     *entities.Timings               `json:"timings"`
}

func NewDB() *DataBase {
	db := DataBase{
		Queue:       taskQueue.NewQueue(),
		Expressions: make(map[uint64]*entities.Expression),
		AllNodes:    make(map[uint64]*entities.Node),
		Timings:     &entities.Timings{},
	}
	return &db
}
func (db *DataBase) Save() error {
	jsonBytes, err := json.Marshal(db)
	if err != nil {
		log.Println(err)
		return err
	}
	wd, err := os.Getwd()
	if err != nil {
		log.Println(err)
		return err
	}
	path := wd + "/db/" + configs.NameDataBase + ".json"
	err = os.WriteFile(path, jsonBytes, 0666)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

// Загружает в пустую структуру DataBase сохраненные данные
func (db *DataBase) Load() error {
	wd, err := os.Getwd()
	if err != nil {
		log.Println(err)
		return err
	}
	path := wd + "/db/" + configs.NameDataBase + ".json"
	data, err := os.ReadFile(path)
	if err != nil {
		log.Println("ошибка открытия json", err)
		return err
	}
	err = json.Unmarshal(data, &db)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

// SafeJSON Сохраняет структуру в базе данных, в папке db
func SafeJSON(name string, dataBase *DataBase) error {
	jsonBytes, err := json.Marshal(dataBase)
	if err != nil {
		log.Println(err)
		return err
	}
	wd, err := os.Getwd()
	if err != nil {
		log.Println(err)
		return err
	}
	path := wd + "/db/" + configs.NameDataBase + ".json"
	err = os.WriteFile(path, jsonBytes, 0666)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

// CreateEmptyDb Cоздает файл пустой БД
func CreateEmptyDb() error {
	// Создаем файл с пустой БД
	err := SafeJSON(configs.NameDataBase, NewDB())
	if err != nil {
		log.Println("Ошибка создания пустой БД")
	}
	return nil
}
