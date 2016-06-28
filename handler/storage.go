package handler

import (
	"static/misc"
	"strings"

	log "code.google.com/p/log4go"
	"gopkg.in/mgo.v2/bson"
)

type Storage struct{}

func (s *Storage) UploadFile(uid, fname, fid string, b []byte) error {
	db := misc.Backend.Db.Copy()
	fs := db.DB("static").GridFS("fs")

	file, err := fs.Create(fname)
	if err != nil {
		log.Warn("failed to create file:[%s:%s:%s:%s], Err:[%s]", uid, fname, err)
		return err
	}

	_, err = file.Write(b)
	if err != nil {
		log.Warn("failed to write file:[%s:%s:%s:%s], Err:[%s]", uid, fname, err)
		return err
	}

	meta := bson.M{"uid": uid, "fid": fid}
	file.SetMeta(meta)

	file.Close()
	db.Close()

	return nil
}

func (s *Storage) DownloadFile(uid, fname string) (interface{}, error) {
	db := misc.Backend.Db.Copy()
	fs := db.DB("static").GridFS("fs")

	file, err := fs.Open(fname)
	if err != nil {
		log.Warn("failed to open file:[%s:%s], Err:[%s]", uid, fname, err)
		return nil, err
	}

	b := make([]byte, 0)
	_, err = file.Read(b)
	if err != nil {
		log.Warn("failed to read file:[%s:%s], Err:[%s]", uid, fname, err)
		return nil, err
	}

	file.Close()
	db.Close()

	return b, nil
}

func (s *Storage) DeleteFile(uid, fname string) error {
	db := misc.Backend.Db.Copy()
	fs := db.DB("static").GridFS("fs")

	err := fs.Remove(fname)
	if err != nil {
		log.Warn("failed to remove file:[%s:%s], Err:[%s]", uid, fname, err)
		return err
	}

	db.Close()

	return nil
}

func (s *Storage) InitUploadUrl(fid string) bool {
	db := misc.Backend.Db.Copy()
	fs := db.DB("static").GridFS("fs")

	var res interface{} = nil
	cond := bson.M{"fid": fid}
	err := fs.Find(cond).One(&res)
	if err != nil {
		log.Warn("failed to find file:[%s], Err:[%s]", fid, err)
		return false
	}

	db.Close()

	return res != nil
}

func (s *Storage) GetMeta(uid, fid string) (interface{}, error) {
	db := misc.Backend.Db.Copy()
	fs := db.DB("static").GridFS("fs")

	res := make(map[string]interface{})
	fids := strings.Split(fid, ",")
	cond := bson.M{"uid": uid, "fid": bson.M{"$in": fids}}
	filter := bson.M{"_id": 1}
	err := fs.Find(cond).Select(filter).One(&res)
	if err != nil {
		log.Warn("failed to find file:[%s:%s], Err:[%s]", uid, fid, err)
		return nil, err
	}

	file, err := fs.OpenId(res["_id"])
	if err != nil {
		log.Warn("failed to open file:[%s:%s:%s], Err:[%s]", uid, fid, res["_id"], err)
		return nil, err
	}

	err = file.GetMeta(&res)
	if err != nil {
		log.Warn("failed to get meta:[%s:%s], Err:[%s]", uid, fid, err)
		return nil, err
	}

	file.Close()
	db.Close()

	return res, nil
}
