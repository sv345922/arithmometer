package tasker

import (
	"arithmometer/orchestr/parsing"
	"context"
	"fmt"
	"log"
	"sync"
)

// Достает рабочее пространство из контекста
func GetWs(ctx context.Context) (*WorkingSpace, bool) {
	mu := sync.Mutex{}
	mu.Lock()
	defer mu.Unlock()
	ws, ok := ctx.Value("ws").(*WorkingSpace)
	return ws, ok
}

// Создает новый список выражений
// .Dict - словарь с ссылками на выражения
// .ListExpr - список с ссылками на выраженя, повторяет .Dict
func NewExpressions() *Expressions {
	res := Expressions{}
	res.Dict = make(map[uint64]*Expression)
	res.ListExpr = make([]*Expression, 0)
	return &res
}

// Добавляет выражение в список выражений
func (e *Expressions) Add(expression *Expression) {
	l_prev := len(e.Dict)
	e.Dict[expression.IdExpression] = expression
	if l_prev+1 == len(e.Dict) {
		e.ListExpr = append(e.ListExpr, expression)
	}
}

// возвращает выражение из списка задач
func FindExpression(id uint64, e *Expressions) *Expression {
	e.mu.RLock()
	defer e.mu.RUnlock()
	if task, ok := e.Dict[id]; ok {
		return task
	}
	return nil
}

// Создает новый список задач
// .Dict - словарь с задачами
// .Queue - очередь задач
func NewTasks() *Tasks {
	res := Tasks{}
	res.Dict = make(map[uint64]*TaskContainer)
	res.Queue = NewDequeue()
	return &res
}

// возвращает задачу из списка (мапы) задач
func FindNodes(id uint64, n map[uint64]*parsing.Node) *parsing.Node {
	if node, ok := n[id]; ok {
		return node
	}
	return nil
}

// Создает рабочее пространство и сохраняет базу данных
func RunTasker() (*WorkingSpace, error) {

	ws := &WorkingSpace{
		Tasks:       NewTasks(),
		Expressions: NewExpressions(),
		Timings:     &Timings{},
		AllNodes:    make(map[uint64]*parsing.Node), // ключом является id узла
	}
	// Восстанавливаем выражения и задачи из базы данных
	err := restoreTaskExpr(ws)
	if err != nil {
		log.Println("ошибка восстановления из БД", err)
	}

	// обновляем ws
	ws.Update()
	err = ws.Save()
	if err != nil {
		err = fmt.Errorf("ошибка сохранения БД: %v", err)
	}
	return ws, err
}

// Восстанавливает рабочее пространство из сохраненной базы данных
func restoreTaskExpr(ws *WorkingSpace) error {
	// Проверка существования БД  и создание пустой бд при необходимости
	//err := checkDb()
	//if err != nil {
	//	log.Println("Ошибка проверки/создания бд", err)
	//}

	// Загрузка сохраненной БД
	savedDb, err := LoadDB()
	if err != nil {
		log.Println("ошибка загрузки бд", err)
	}
	ws.Expressions = savedDb.Expressions
	ws.Tasks = savedDb.Tasks
	return nil
}

// Cоздает файл пустой БД
func CreateEmptyDb() error {
	// получить рабочую папку
	//wd, err := os.Getwd()
	//if err != nil {
	//	log.Println(err)
	//	return err
	//}
	// путь файла базы данных
	//path := wd + "\\orchestr\\db\\db.json"

	//// проверка на существование файла
	//_, fileError := os.Stat(path)
	//// если он существует выходим
	//if os.IsNotExist(fileError) {
	//	fmt.Printf("бд существует") // TODO удалить
	//	return nil
	//}
	// Создаем файл с пустой БД
	err := SafeJSON[DataBase]("db", *NewDB())
	if err != nil {
		log.Println("Ошибка создания пустой БД")
	}
	return nil
}
