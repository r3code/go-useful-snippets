package transactor

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"
	"github.com/juju/errors"

	multierror "github.com/hashicorp/go-multierror"
)

// Transaction интерфейс моделирующий интерфейс стандартной транзакции в
// `database/sql`.
//
// Позволяет скрыть от функции `TxFunc` методы commit и rollback (они
// вызываются функцией `WithTransaction`,
// чтобы они не были вызваны случайно, потому эти методы сюда не включены.
type Transaction interface {
	Exec(query string, args ...interface{}) (sql.Result, error)
	Prepare(query string) (*sql.Stmt, error)
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(query string, args ...interface{}) *sql.Row
	QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
}

// TransactionX интерфейс моделирующий методы транзакции sqlx,
// но без Commit и Rollback, чтобы внутри `TxFunc` нельзя было вызвать их,
// т.к. работой транзакии управляет функция `WithTransaction`
type TransactionX interface {
	Transaction
	Rebind(query string) string
	BindNamed(query string, arg interface{}) (string, []interface{}, error)
	NamedQuery(query string, arg interface{}) (*sqlx.Rows, error)
	NamedExec(query string, arg interface{}) (sql.Result, error)
	Select(dest interface{}, query string, args ...interface{}) error
	Queryx(query string, args ...interface{}) (*sqlx.Rows, error)
	QueryRowx(query string, args ...interface{}) *sqlx.Row
	Get(dest interface{}, query string, args ...interface{}) error
	MustExec(query string, args ...interface{}) sql.Result
	Preparex(query string) (*sqlx.Stmt, error)
	Stmtx(stmt interface{}) *sqlx.Stmt
	NamedStmt(stmt *sqlx.NamedStmt) *sqlx.NamedStmt
	PrepareNamed(query string) (*sqlx.NamedStmt, error)
}

// TxFunc будет вызвано с инициализированным объектом `Transaction`
// которые может быть использоват для исполнения выражений и запросов к БД.
type TxFunc func(tx TransactionX) error

// WithTransaction создает новую транзакцию, выполняет функцию `TxFunc`,
// если она выполнена с ошибками, то выполняется откат транзакции (rollback),
// иначе фиксируется транзакция (commit)
func WithTransaction(db *sqlx.DB, txFunc TxFunc) (err error) {
	tx, err := db.Beginx()
	if err != nil {
		return errors.Annotate(err, "Transaction Begin failed")
	}
	err = execute(tx, txFunc)
	if err != nil {
		return err
	}
	return nil
}

func execute(tx *sqlx.Tx, txFunc TxFunc) (err error) {
	var handleErrors = func() {
		if p := recover(); p != nil { // err = nil
			// где-то случилась паника
			err = errors.Errorf("panic in WithTransaction(): %v", p)
			rollErr := tx.Rollback()
			if rollErr != nil {
				err = multierror.Append(err, errors.Annotate(rollErr, "Failed Transaction Rollback after panic"))
			}
			return // err заполнено
		}

		if err != nil {
			err = errors.Annotate(err, "Wrapped function exit with error")
			rollErr := tx.Rollback()
			if rollErr != nil {
				err = multierror.Append(err, errors.Annotate(rollErr, "Failed Transaction Rollback after error in wrapped function"))
			}
			return // err заполнено
		}
		// Если все в проядке, то зафиксируем
		err = tx.Commit()
		if err != nil {
			err = errors.Annotate(err, "Failed Transaction Commit")
		}
	}

	defer handleErrors()
	err = txFunc(tx)
	return err
}

// WithTransactionMany запустить несколько функций в транзакции
func WithTransactionMany(db *sqlx.DB, funcs ...TxFunc) error {
	for _, fn := range funcs {
		if err := WithTransaction(db, fn); err != nil {
			return err
		}
	}
	return nil
}

// WithCtxTransaction creates a new transaction with ctx and handles rollback/commit based on the
// error object returned by the `txFunc`
func WithCtxTransaction(ctx context.Context, opt *sql.TxOptions, db *sqlx.DB, fn TxFunc) (err error) {
	tx, err := db.BeginTxx(ctx, opt)
	if err != nil {
		return errors.Annotate(err, "Transaction BeginTx failed")
	}
	return execute(tx, fn)
}
