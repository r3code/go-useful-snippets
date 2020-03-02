package log

import (
	"io"
	"os"
	"strings"

	kitlog "github.com/go-kit/kit/log"
)

// Logger переименовывает  (go-kit/log).Logger интерфейс для удобства
type Logger kitlog.Logger

var (
	// DefaultCaller указывает на место вызова, добавляет свойство "caller" в лог
	DefaultCaller = kitlog.Caller(5)
	// DefaultTimestampUTC определяет временую метку, используется в свойстве "ts" в логе
	DefaultTimestampUTC = kitlog.DefaultTimestampUTC
)

// NewDefaultStdOutLogger создает логгер с настройками по умолчанию для вывода в STDOUT
func NewDefaultStdOutLogger(disableTimestamp bool) Logger {
	var maxErrors uint8 = 10
	return NewStdOutLogger(maxErrors, disableTimestamp)
}

// NewStdOutLogger создает логгер для вывода сообщений в STDOUT
func NewStdOutLogger(maxErrors uint8, disableTimestamp bool) Logger {
	lg := NewLogger(os.Stdout, maxErrors, disableTimestamp)

	return NewImportantLogger(lg, maxErrors)
}

// NewLogger создает новый логер с указанным писателем.
// Параметр `maxErrors` указывает предельное количество ошибок запси в лог после достижения которого будет вызвана паника и выполнение программы прервется
func NewLogger(writer io.Writer, maxErrors uint8, disableTimestamp bool) Logger {
	lg := kitlog.NewLogfmtLogger(kitlog.NewSyncWriter(writer))

	il := NewImportantLogger(lg, maxErrors)

	il = kitlog.WithPrefix(il, "caller", DefaultCaller)
	if !disableTimestamp {
		il = kitlog.WithPrefix(il, "time", DefaultTimestampUTC)
	}

	return il
}

// With добавляет новые постоянно добавляемые поля со значениями в сообщение
func With(l Logger, keyvals ...interface{}) Logger {
	return kitlog.With(l, keyvals...)
}

// MustCreateComponentLog создает новый логгер для компонента или паникует, если имя не указано
func MustCreateComponentLog(l Logger, componentName string) Logger {
	if strings.TrimSpace(componentName) == "" {
		panic("Can not create named logger. Empty component name passed")
	}

	return With(l, "component", componentName)
}
