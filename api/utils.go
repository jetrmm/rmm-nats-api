package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	nats "github.com/nats-io/nats.go"
)

func setupNatsOptions(key string) []nats.Option {
	opts := []nats.Option{
		nats.Name("JetRMM"),
		nats.UserInfo("jetrmm", key),
		nats.ReconnectWait(time.Second * 2),
		nats.RetryOnFailedConnect(true),
		nats.MaxReconnects(-1),
		nats.ReconnectBufSize(-1),
	}
	return opts
}

func GetConfig(cfg string) (db *sqlx.DB, r WebConfig, err error) {
	if cfg == "" {
		cfg = "nats-api.conf"
		if !FileExists(cfg) {
			err = errors.New("unable to find config file")
			return
		}
	}

	jret, _ := ioutil.ReadFile(cfg)
	err = json.Unmarshal(jret, &r)
	if err != nil {
		return
	}

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		r.Host, r.Port, r.User, r.Pass, r.DBName)

	db, err = sqlx.Connect("postgres", psqlInfo)
	if err != nil {
		return
	}
	db.SetMaxOpenConns(20)
	return
}

func FileExists(path string) bool {
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}
