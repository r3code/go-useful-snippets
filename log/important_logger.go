package log

import (
	gklog "github.com/go-kit/kit/log"
	"github.com/juju/errors"
)

/* See https://github.com/go-kit/kit/issues/164#issuecomment-274185353

Пример использования:
func main() {
    // ...
    var logger log.Logger
    logger = importantLogger{
        logger: gklog.NewLogfmtLogger(log.NewSyncWriter(os.Stdout)),
        thresh: 10,
        metricErrCounter: logErrorCounter,
    }
    logger = log.NewContext(logger).With("ts",log.DefaultTimestampUTC,
		"caller", log.DefaultCaller)
    // ...
    app.AppLogic(logger, 5, "hello, world")
}

*/

type importantLogger struct {
	logger gklog.Logger
	// metricErrCounter metrics.Counter
	errorCount uint8
	maxErrors  uint8
}

// importantLogger отслеживает число ошибок записи писателя лога, а когда число ошибок переваливает за указанное число `maxErrors`, то он вызывает панику.
// Реализует интерфейс log.Logger. 
// Всегда возвращает только nil или вызывает panic.
func (il *importantLogger) Log(keyvals ...interface{}) error {
	if err := il.logger.Log(keyvals...); err != nil {
		il.errorCount++
		// TODO(dsin): добавить счетичик для вывода метрики "логгер.число_ошибок_записи" - бекэнд или expvar или prometheus
		// il.metricErrCounter.Add(1)
		tooMuchLogFails := il.errorCount > il.maxErrors
		if tooMuchLogFails {
			panic(errors.Annotate(err, "Logger write failed"))
		}
	}
	return nil
}

// NewImportantLogger созает экземпляр важного логгера из обычного, отказ которого вызывает панику и остановку программы.
// Параметр `maxErrors` указывает предельное количество ошибок запси в лог после достижения которого будет вызвана паника и выполнение программы прервется.
// Если `maxErrors` = 0, то используется значение `maxErrors` по умолчанию (=1). 
func NewImportantLogger(logger gklog.Logger, maxErrors uint8) Logger {
	if maxErrors == 0 {
		maxErrors = 1
	}
	return &importantLogger{
		logger:     logger,
		maxErrors:  maxErrors,
		errorCount: 0,
	}
}
