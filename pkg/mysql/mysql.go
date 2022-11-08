package mysql

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/ninehills/go-webapp-template/pkg/logger"
	"github.com/qustavo/sqlhooks/v2"
)

const (
	defaultMaxOpenConns    = 10
	defaultMaxIdleConns    = 10
	defaultConnMaxLifetime = time.Second * 120
	connAttempts           = 3
	connAttemptPeriod      = time.Second * 5
)

// MySQL -.
type MySQL struct {
	maxOpenConns    int
	maxIdleConns    int
	connMaxLifetime time.Duration

	DB *sql.DB
}

// New -.
func New(l logger.Logger, dsn string, opts ...Option) (*MySQL, error) {
	ms := &MySQL{
		maxOpenConns:    defaultMaxOpenConns,
		maxIdleConns:    defaultMaxIdleConns,
		connMaxLifetime: defaultConnMaxLifetime,
	}

	// Custom options
	for _, opt := range opts {
		opt(ms)
	}

	// Register hook
	sql.Register("mysqllog", sqlhooks.Wrap(&mysql.MySQLDriver{}, NewHook(l)))

	// Open database
	var err error

	ms.DB, err = sql.Open("mysqllog", dsn)
	if err != nil {
		return nil, fmt.Errorf("mysql - NewMySQL - Open mysql database failed: %w", err)
	}

	ms.DB.SetConnMaxLifetime(ms.connMaxLifetime)
	ms.DB.SetMaxOpenConns(ms.maxOpenConns)
	ms.DB.SetMaxIdleConns(ms.maxIdleConns)

	i := connAttempts
	for i > 0 {
		// Ping database
		err = ms.DB.Ping()
		if err == nil {
			break
		}

		l.Infof("mysql - NewMySQL - Trying to connect, attempts left: %d", i)
		time.Sleep(connAttemptPeriod)
		i--
	}

	if err != nil {
		return nil, fmt.Errorf("mysql - NewMySQL - Connect attempts %d times failed: %w", connAttempts, err)
	}

	return ms, nil
}

// Close -.
func (m *MySQL) Close() {
	if m.DB != nil {
		m.DB.Close()
	}
}
