// Package log пакет предназначен для упрощения ведения логов в программах.
// Представляет надстройку над go-kit/kit/log. 
// Этот пакет создан для ведения важных логов, т.е. без нормально работающего логгера продолжение работы программы не имеет смысла, т.к. разрабочик не получит важную информацию в случае откраза.
// Потому этот парет предоставялет возможность использовать важный логгер (см. `NewImportantLogger()`), который следит за количеством ошибок при выполнение записи в носитель лога (например, файл). И если число ошибок превышает заданное - вызывает панику.
package log

// Logger интерфейс логгера базовый для препятствования сильному связыванию
// см. https://dave.cheney.net/2015/11/05/lets-talk-about-logging
// https://docs.google.com/document/d/1oTjtY49y8iSxmM9YBaz2NrZIlaXtQsq3nQMd-E0HtwM/edit
// https://groups.google.com/forum/#!topic/golang-dev/F3l9Iz1JX4g
/*
Это позволит выводить любой тип данных, например:

// common formatted string output
Log(fmt.Sprintf(“some custom print string: %s”, value))

// no format string, just values
Log(someText, struct1, struct2, err)

// Custom key value style structs that can be interrogated with the logger that satisfied the interface to create structured logging
Log(key, value, key, value)

// log levels (with a colog style implementation)
Log(“info: starting up system”)
Log(fmt.Sprintf(“err: kaboom %s”, err))
 */
