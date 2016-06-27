package handler

import (
	"static/misc"

	log "code.google.com/p/log4go"
	"gopkg.in/mgo.v2/bson"
)

type Storage struct{}

func (s *Storage) UploadFile(uid, fname string, b []byte) error {
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

	meta := bson.M{"uid": uid}
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
