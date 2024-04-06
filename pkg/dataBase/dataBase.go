package dataBase

import (
	"arithmometer/pkg/expressions"
	"arithmometer/pkg/taskQueue"
	"arithmometer/pkg/timings"
	"encoding/json"
	"log"
	"os"
)

type DataBase struct {
	// список выражений (с таймингами)
	Tasks       *taskQueue.Tasks          `json:"tasks"`
	Expressions []*expressions.Expression `json:"expressions"` // []Expression
	Timings     timings.Timings           `json:"timings"`
	AllNodes    map[uint64]NodeDB         `json:"allNodes"` // map[uint64]NodeDB
}

type NodeDB struct {
	// разница с treeExpression.Node в том, что вместо ссылок на дочернии узлы,
	// хранятся идентификаторы узлов
	NodeId     uint64  `json:"nodeId"`
	Op         string  `json:"op"` // оператор
	XId        uint64  `json:"x"`
	YId        uint64  `json:"y"`
	Val        float64 `json:"Val"`   // значение узла
	Sheet      bool    `json:"sheet"` // флаг листа
	Calculated bool    `json:"calculated"`
	ParentId   uint64  `json:"parent"` // узел родитель
}

func NewDB() *DataBase {
	result := DataBase{
		Tasks:       taskQueue.NewTasks(),
		Expressions: []*expressions.Expression{},
		Timings:     timings.Timings{},
		AllNodes:    make(map[uint64]NodeDB),
	}
	return &result
}

// SafeJSON Сохраняет структуру в базе данных, в папке db
func SafeJSON(name string, expr *DataBase) error {
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

// CreateEmptyDb Cоздает файл пустой БД
func CreateEmptyDb() error {

	// Создаем файл с пустой БД
	err := SafeJSON[DataBase]("db", NewDB())
	if err != nil {
		log.Println("Ошибка создания пустой БД")
	}
	return nil
}
