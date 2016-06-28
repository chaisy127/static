package handler

import (
	"static/misc"

	log "code.google.com/p/log4go"
	"gopkg.in/mgo.v2/bson"
)

type Storage struct{}

func (s *Storage) UploadFile(fname, fid, bucketName string, b []byte) error {
	db := misc.Backend.Db.Copy()
	fs := db.DB("static").GridFS("fs")
	c := db.DB("static").C("bucket")

	cond := bson.M{"bucketname": bucketName}
	res := bson.M{"bucketname": bson.M{"$set": bucketName}}
	_, err := c.Upsert(cond, res)
	if err != nil {
		log.Warn("failed to create bucket:[%s], Err:[%s]", bucketName, err)
		return err
	}

	file, err := fs.Create(fid)
	if err != nil {
		log.Warn("failed to create file:[%s:%s:%s], Err:[%s]", bucketName, fname, fid, err)
		return err
	}

	_, err = file.Write(b)
	if err != nil {
		log.Warn("failed to write file:[%s:%s:%s], Err:[%s]", bucketName, fname, fid, err)
		return err
	}

	meta := bson.M{"fid": fid, "fname": fname, "bucketname": bucketName}
	file.SetMeta(meta)

	file.Close()
	db.Close()

	return nil
}

func (s *Storage) DownloadFile(fid, bucketName string) (interface{}, error) {
	db := misc.Backend.Db.Copy()
	fs := db.DB("static").GridFS("fs")

	var res map[string]int = nil
	cond := bson.M{"metadata.bucketname": bucketName, "filename": fid}
	filter := bson.M{"length": 1}
	err := fs.Find(cond).Select(filter).One(&res)
	if err != nil || res == nil {
		log.Warn("failed to find file:[%s:%s], Err:[%s]", bucketName, fid, err)
		return nil, err
	}

	file, err := fs.Open(fid)
	if err != nil {
		log.Warn("failed to open file:[%s:%s], Err:[%s]", bucketName, fid, err)
		return nil, err
	}

	b := make([]byte, res["length"])
	_, err = file.Read(b)
	if err != nil {
		log.Warn("failed to read file:[%s:%s], Err:[%s]", bucketName, fid, err)
		return nil, err
	}

	file.Close()
	db.Close()

	return string(b), nil
}

func (s *Storage) DeleteFile(fid, bucketName string) error {
	db := misc.Backend.Db.Copy()
	fs := db.DB("static").GridFS("fs")

	var res interface{} = nil
	cond := bson.M{"filename": fid, "metadata.bucketname": bucketName}
	err := fs.Find(cond).One(&res)
	if err != nil || res == nil {
		log.Warn("failed to find file:[%s:%s], Err:[%s]", bucketName, fid, err)
		return err
	}

	err = fs.Remove(fid)
	if err != nil {
		log.Warn("failed to remove file:[%s:%s], Err:[%s]", bucketName, fid, err)
		return err
	}

	db.Close()

	return nil
}

func (s *Storage) InitUploadUrl(fid string) bool {
	db := misc.Backend.Db.Copy()
	fs := db.DB("static").GridFS("fs")

	var res interface{} = nil
	cond := bson.M{"filename": fid}
	err := fs.Find(cond).One(&res)
	if err != nil {
		log.Warn("failed to find file:[%s], Err:[%s]", fid, err)
		return false
	}

	db.Close()

	return res != nil
}

func (s *Storage) GetMeta(fid, bucketName string) (interface{}, error) {
	db := misc.Backend.Db.Copy()
	fs := db.DB("static").GridFS("fs")

	var res interface{} = nil
	cond := bson.M{"filename": fid}
	err := fs.Find(cond).One(&res)
	if err != nil {
		log.Warn("failed to find file:[%s], Err:[%s]", fid, err)
		return nil, err
	}

	db.Close()

	return res, nil
}
