package mysql

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"sync"

	dockerinitiator "github.com/Storytel/go-docker-initiator"
	mysqldrv "github.com/go-sql-driver/mysql"
)

var _ dockerinitiator.Probe = MysqlProbe{}

// MysqlProbe implements the IProbe interface for mysql instances
type MysqlProbe struct {
	MysqlConfig
}

var logmtx = sync.Mutex{}

// DoProbe will probe by waiting for log messages
func (i MysqlProbe) DoProbe(instance *dockerinitiator.Instance) error {

	silent := log.New(ioutil.Discard, "", 0)
	logmtx.Lock()
	mysqldrv.SetLogger(silent)
	logmtx.Unlock()

	defer func() {
		logmtx.Lock()
		defer logmtx.Unlock()
		// This is the default logger for mysql
		mysqldrv.SetLogger(log.New(os.Stderr, "[mysql] ", log.Ldate|log.Ltime|log.Lshortfile))
	}()

	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s)/%s", "root", i.Password, instance.GetHost(), i.DbName))
	defer db.Close()
	if err != nil {
		return err
	}

	var version string
	err = db.QueryRow("SELECT VERSION()").Scan(&version)
	if err != nil && err != sql.ErrNoRows {
		return err
	}

	return nil
}
