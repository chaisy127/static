package handler

import (
	"static/misc"

	log "code.google.com/p/log4go"
	"gopkg.in/mgo.v2/bson"
)

type Bucket struct{}

func (b *Bucket) ListBucket(bucketName string) ([]string, error) {
	db := misc.Backend.Db.Copy()
	fs := db.DB("statis").GridFS("fs")

	res := make([]map[string]map[string]string, 0)
	cond := bson.M{"metadata.bucketname": bucketName}
	filter := bson.M{"metadata.fname": 1}
	err := fs.Find(cond).Select(filter).All(&res)
	if err != nil {
		log.Warn("failed to find file:[%s], Err:[%s]", bucketName, err)
		return nil, err
	}

	r := make([]string, 0)
	for _, i := range res {
		fname := i["metadata"]["fname"]
		r = append(r, fname)
	}

	db.Close()

	return r, nil
}

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
