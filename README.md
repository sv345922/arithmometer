# arithmometer
@VladimirViSi (telegram)


Принимаемой выражение должно состоять из чередующихся последовательно чисел и операторов: )(+-*/
пробелы между символами не учитываются.
Первое число в выражении может быть отрицательным. В скобках первое число не должно быть 
отрицательным, в операторах умножения число не может быть отрицательным

Выражение парсится по алгоритму Дейкстры, через постфиксную запись
в базе данных выражение сохраняется в виде списка символов в постфиксной записи.

Сервер работает по адресу "127.0.0.1:8000"
Выражение передается клиентом POST запросом в следующем виде:
например 127.0.0.1:8000/newexpression?expr=-1+(2-3)/4+5
в теле запроса передается JSON таймингов операций в секундах
{"plus":1,"minus":1,"mult":1,"div":1}
клиенту возвращается id выражения

Клиент спрашивает о завершении расчетов передавая id выражения:
получает ответ (результат или "ожидайте")
! Любой клиент, который знает id, может получить ответ

! Сервер может отрабатывать только одно задание клиента


**запуск сервера**
go run ./cmd/orch/

**запуск клиента**
go run ./client/

