package handler

import (
	"static/misc"

	log "code.google.com/p/log4go"
	"gopkg.in/mgo.v2/bson"
)

type Bucket struct{}

func (b *Bucket) CreateBucket(bucketName string) error {
	db := misc.Backend.Db.Copy()
	c := db.DB("static").C("bucket")

	res := bson.M{"bucketname": bucketName}
	err := c.Insert(res)
	if err != nil {
		log.Warn("failed to create bucket:[%s], Err:[%s]", bucketName, err)
		return err
	}

	db.Close()

	return nil
}

func (b *Bucket) DeleteBucket(bucketName string) error {
	db := misc.Backend.Db.Copy()
	c := db.DB("static").C("bucket")

	res := bson.M{"bucketname": bucketName}
	err := c.Remove(res)
	if err != nil {
		log.Warn("failed to delete bucket:[%s], Err:[%s]", bucketName, err)
		return err
	}

	db.Close()

	return nil
}
