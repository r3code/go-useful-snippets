package transactor_test

import (
	"github.com/r3code/go-useful-snippets/dbutils/transactor"
	"flag"
	"fmt"
	
	"os"
	"testing"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/juju/errors"
	_ "github.com/lib/pq"
)

var (
	dbHost     = flag.String("test.dbhost", "localhost", "Postgres Database host address")
	dbPort     = flag.Int("test.dbport", 5432, "Postgres Database port number")
	dbUser     = flag.String("test.dbuser", "test", "Postgres Database user name")
	dbPassword = flag.String("test.dbpass", "test", "Postgres Database user password")
	dbName     = flag.String("test.dbname", "test", "Postgres Database name")
)

var cachedDBConn *sqlx.DB

// GetConnection возвращает экземпляр пула соединений *sqlx.DB, если экземпляр создан ранее, то вернет его иначе, вернет уже созданный ранее. Сама проверяет возможность подключиться раз в секунду в течение установленного таймаута
func GetConnection(datasource string, connectTimeoutSec int8) (*sqlx.DB, error) {
	if cachedDBConn != nil {
		return cachedDBConn, nil
	}
	cachedDBConn, dbErr := sqlx.Open("postgres", datasource)
	if dbErr != nil {
		return nil, errors.Annotate(dbErr, "Database pool opening error")
	}
	pingErr := pingDatabase(cachedDBConn, connectTimeoutSec)
	if pingErr != nil {
		return nil, errors.Annotate(pingErr, "Database connect check error. ")
	}

	return cachedDBConn, nil
}

// проверяем возможность подключения к БД раз в секунду за установленный таймаут, если удалось то выйдет без ошибки
func pingDatabase(db *sqlx.DB, connectTimeoutSec int8) (err error) {
	for i := int8(0); i < connectTimeoutSec; i++ {
		err = db.Ping()

		if err == nil {
			return nil
		}
		time.Sleep(time.Second)
	}
	return err
}

// MustGetTestConnection получить соединение с тестовой базой или вызвать panic при ошибках
func MustGetTestConnection() (conn *sqlx.DB) {
	dbConnStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		*dbHost, uint16(*dbPort), *dbUser, *dbPassword, *dbName, "disable")

	db, err := GetConnection(dbConnStr, 5)
	if err != nil {
		panic(fmt.Sprintf("Failed to create DB connection: %v", err))
	}
	return db
}

func MustMigrateDB(db *sqlx.DB) {
	db.Exec("DROP TABLE public.wtt")
	db.MustExec(`
	CREATE TABLE public.wtt
	(
		name character varying COLLATE pg_catalog."default" NOT NULL,
		CONSTRAINT wtt_pkey PRIMARY KEY (name)
	)
	WITH (
		OIDS = FALSE
	)
	TABLESPACE pg_default;
	
	ALTER TABLE public.wtt
		OWNER to postgres;`)
}

func TestWithTransaction_MustSuccess(t *testing.T) {
	if testing.Short() {
		t.Skip("Skip integration test")
	}
	okFunc := func(tx transactor.TransactionX) error {
		rows, err := tx.Queryx("SELECT * from public.wtt")
		defer rows.Close()
		return err // ошибки быть не должно
	}
	err := transactor.WithTransaction(MustGetTestConnection(), okFunc)
	if err != nil {
		t.Error("Error returned istead of return `nil`", err)
	}
}

func TestWithTransaction_PanicMustRollback(t *testing.T) {
	if testing.Short() {
		t.Skip("Skip integration test")
	}
	panicFunc := func(tx transactor.TransactionX) error {
		tx.MustExec("INSERT INTO public.wtt VALUES('test_panic')")
		panic("test panic")
	}
	db := MustGetTestConnection()
	err := transactor.WithTransaction(db, panicFunc)
	count := -1
	db.Get(&count, "SELECT count(*) FROM public.wtt WHERE 'name' = 'test_panic'")
	if count != 0 {
		t.Error("No rollback happened after panic")
	}
	if err == nil {
		t.Error("No error raised after panic", err)
	}
}

func TestWithTransaction_TxFuncErrMustRollback(t *testing.T) {
	if testing.Short() {
		t.Skip("Skip integration test")
	}
	funcError := errors.New("@txFunc test error@")

	errorringFunc := func(tx transactor.TransactionX) error {
		tx.MustExec("INSERT INTO public.wtt VALUES('test_error')")
		return funcError
	}
	db := MustGetTestConnection()
	err := transactor.WithTransaction(db, errorringFunc)
	count := -1
	db.Get(&count, "SELECT count(*) FROM public.wtt WHERE 'name' = 'test_error'")
	if count != 0 {
		t.Error("No rollback happened after txFunc error")
	}
	orgiginErr := errors.Cause(err)
	if orgiginErr != funcError {
		t.Errorf("Not raised an error after txFunc error, have %v, want %v", err, funcError)
	}
}

func setup() {
	db := MustGetTestConnection()
	MustMigrateDB(db)
}
func shutdown() {
	// no
}

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	shutdown()
	os.Exit(code)
}
