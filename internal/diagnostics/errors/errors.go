package errors

// Не менять массив!!
var errors = []Error{}

type Error struct {
	Code        uint16
	Message     string // сообщение об ошибке, кратокое описание
	Arrow       string // строка с указанием места ошибки и пояснением. если равен "BLOCK=FALSE" то не будет блока с кодом, а только указание на строку и столбец
	Description string // подробное описание ошибки
}
