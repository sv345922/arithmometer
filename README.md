# arithmometer
@VladimirViSi (telegram)


Принимаемой выражение должно состоять из чередующихся последовательно чисел и операторов: )(+-*/
пробелы между символами могут присутствовать.
Первое число в выражении может быть отрицательным. 
В скобках первое число не должно быть отрицательным, в операторах умножения число 
не может быть отрицательным (проверить)

Выражение парсится по алгоритму Дейкстры, через постфиксную запись
в базе данных выражение сохраняется в виде списка символов в постфиксной записи.

Сервер работает по адресу "127.0.0.1:8000"
Выражение передается клиентом POST запросом по адресу 
127.0.0.1:8000/newexpression
в теле запроса передается JSON со строкой выражения и таймингами операций в секундах
{
Expr:    "-1+(2-3)/4+5",
Timings: {"plus":1,"minus":1,"mult":1,"div":1},
}
клиенту возвращается id выражения

Клиент спрашивает о завершении расчетов передавая id выражения в Get запросе
получает ответ (результат или "ожидайте")
например 
curl -v localhost:8000/getresult?id=1708192194295205100
конце id выражения(может меняться)

! Любой клиент, который знает id, может получить ответ


Запускать все из папки arithmometr

#### Запуск сервера

_go run ./cmd/orch/_

#### Запуск клиента

_go run ./client/_

#### Запуск вычислителя

_go run ./cmd/calculator/_

сколько раз запустишь, столько и будет вычислителей. Запускать либо в фоне, либо
в отдельных сеансах командной строки


curl -v localhost:8000/getresult?id=1708192194295205100
цифры в конце id выражения

