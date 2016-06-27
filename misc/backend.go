package misc

import (
	"strings"
	"time"

	"gopkg.in/mgo.v2"
)

type backend struct {
	Db *mgo.Session
}

var Backend *backend

func createDb(mgoInfo *mgo.DialInfo) (*mgo.Session, error) {
	mgoSess, err := mgo.DialWithInfo(mgoInfo)
	if err != nil {
		return nil, err
	}
	mgoSess.SetMode(mgo.Monotonic, true)
	return mgoSess, nil
}

func InitBackend() error {
	mgoSess, err := createDb(&mgo.DialInfo{
		Addrs:    strings.Split(Conf.Mongo.Hosts, ","),
		Timeout:  Conf.Mongo.Timeout * time.Second,
		Database: Conf.Mongo.Db,
		Username: Conf.Mongo.User,
		Password: Conf.Mongo.Passwd,
	})
	if err != nil {
		return err
	}

	Backend = &backend{
		Db: mgoSess,
	}
	return nil
}
