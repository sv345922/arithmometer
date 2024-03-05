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
	Tasks       *Tasks            `json:"tasks"`
	Expressions []*Expression     `json:"expressions"` // []Expression
	Timings     Timings           `json:"timings"`
	AllNodes    map[uint64]NodeDB `json:"allNodes"` // map[uint64]NodeDB
}

type NodeDB struct {
	NodeId     uint64  `json:"nodeId"`
	Op         string  `json:"op"` // оператор
	XId        uint64  `json:"x"`
	YId        uint64  `json:"y"`     // потомки
	Val        float64 `json:"Val"`   // значение узла
	Sheet      bool    `json:"sheet"` // флаг листа
	Calculated bool    `json:"calculated"`
	ErrZeroDiv error   `json:"err"`
	ParentId   uint64  `json:"parent"` // узел родитель
}

type additiveJSON interface {
	DataBase
}

func NewDB() *DataBase {
	result := DataBase{
		Tasks:       NewTasks(),
		Expressions: []*Expression{},
		Timings:     Timings{},
		AllNodes:    make(map[uint64]NodeDB),
	}
	return &result
}

func SaveWS(ws *WorkingSpace) error {
	ws.mu.RLock()
	db := NewDB()
	// Заполняем структуру БД
	// заполняем тайминги
	db.Timings = *ws.Timings
	// заполняем очередь задач
	db.Tasks = ws.Tasks
	// создаем список существующих выражений
	// и заполняем список выражений
	for _, expression := range ws.Expressions.ListExpr {
		db.Expressions = append(db.Expressions, expression)
	}
	// создаем мапу узлов для сохранения и заносим ее в базу данных
	for key, val := range *ws.AllNodes {
		node := NodeDB{
			NodeId:     key,
			Op:         val.Op,
			XId:        0,
			YId:        0,
			Val:        val.Val,
			Sheet:      val.Sheet,
			Calculated: val.Calculated,
			ErrZeroDiv: val.ErrZeroDiv,
			ParentId:   0,
		}
		// заполняем id дочерних узлов и родителей, если их нет (лист или корень дерева)
		// оставляем значение по умолчанию (0)
		if val.X != nil {
			node.XId = val.X.NodeId
		}
		if val.Y != nil {
			node.YId = val.Y.NodeId
		}
		if val.Parent != nil {
			node.ParentId = val.Parent.NodeId
		}
		db.AllNodes[key] = node
	}
	ws.mu.RUnlock()
	err := SafeJSON[DataBase]("db", *db)
	if err != nil {
		return err
	}
	log.Println("DB saved")
	return nil
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
func LoadDB() (*WorkingSpace, error) {
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
	// Создаем рабочее пространство
	ws := WorkingSpace{
		Tasks:       result.Tasks,
		Expressions: NewExpressions(),
		Timings:     &result.Timings,
		AllNodes:    NewNodes(),
		mu:          sync.RWMutex{},
	}
	ws.mu.Lock()
	// Заполняем список выражений
	for _, value := range result.Expressions {
		ws.Expressions.Add(value)
	}
	// Заполняем список узлов
	for key, value := range result.AllNodes {
		node := parsing.Node{
			NodeId:     value.NodeId,
			Op:         value.Op,
			X:          nil,
			Y:          nil,
			Val:        value.Val,
			Sheet:      value.Sheet,
			Calculated: value.Calculated,
			ErrZeroDiv: value.ErrZeroDiv,
			Parent:     nil,
			Mu:         sync.RWMutex{},
		}
		(*ws.AllNodes)[key] = &node
	}
	// Заполняем связи деревьев узлов
	for key, value := range result.AllNodes {
		if X, ok := (*ws.AllNodes)[value.XId]; ok {
			(*ws.AllNodes)[key].X = X
		}
		if Y, ok := (*ws.AllNodes)[value.YId]; ok {
			(*ws.AllNodes)[key].Y = Y
		}
		if parent, ok := (*ws.AllNodes)[value.ParentId]; ok {
			(*ws.AllNodes)[key].Parent = parent
		}
	}
	ws.mu.Unlock()
	return &ws, nil
}

// Cоздает файл пустой БД
func CreateEmptyDb() error {

	// Создаем файл с пустой БД
	err := SafeJSON[DataBase]("db", *NewDB())
	if err != nil {
		log.Println("Ошибка создания пустой БД")
	}
	return nil
}
