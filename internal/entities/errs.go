package entities

import "errors"

var ErrExpression = errors.New("invalid expression")

var ZeroDiv = errors.New("zero division")

var NoTaskInQueue = errors.New("no task in queue")

var QueueError = errors.New("queue error")

// невозможная ошибка калькулятора - неправильный оператор
var OpEr = errors.New("invalid operator")
