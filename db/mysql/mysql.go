package mysql

import (
	"fmt"
	"github.com/cy18cn/zlog"
	"github.com/jinzhu/gorm"
	"github.com/spf13/viper"
)

type SQLDataSource struct {
	db *gorm.DB
}

type DBOptions struct {
	Address   string `json:"address" required:"true"`
	User      string `json:"user" required:"true"`
	Password  string `json:"password" required:"true"`
	Schema    string `json:"schema" required:"true"`
	Charset   string `json:"charset"`
	ParseTime bool   `json:"parse_time"`
	Loc       string `json:"loc"`
	LogMode   bool   `json:"log_mode"`
	IdleConns int    `json:"idle_conns"`
	MaxConns  int    `json:"max_conns"`
	Dialect   string `json:"dialect"`
}

func NewSQLDataSource(opts *DBOptions) (*SQLDataSource, error) {
	conn := fmt.Sprintf("%s:%s@(%s)/%s?charset=%s&parseTime=%s&loc=Local", opts.User,
		opts.Password, opts.Address, opts.Schema,
		opts.Charset, opts.ParseTime)
	db, err := gorm.Open(opts.Dialect, conn)
	if err != nil {
		zlog.Errorf("initialize repository database failed, err: %v", err)
		// logrus.WithError(err).Fatalln("initialize repository database failed")
		return nil, err
	}
	//enable Gorm repository log
	if opts.LogMode {
		db.LogMode(true)
	}

	db.DB().SetMaxIdleConns(viper.GetInt("repository.idleConns"))
	db.DB().SetMaxOpenConns(viper.GetInt("repository.maxConns"))

	return &SQLDataSource{db}, nil
}

func (self *SQLDataSource) DB() *gorm.DB {
	return self.db
}

func (self *SQLDataSource) Close() error {
	if self.db != nil {
		return self.db.Close()
	}

	return nil
}
