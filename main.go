package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"static/handler"
	"static/misc"

	log "code.google.com/p/log4go"
)

type Response struct {
	ErrNo  int         `json:"errno"`
	ErrMsg string      `json:"errmsg"`
	Data   interface{} `json:"data"`
}

func responseError(w http.ResponseWriter, errno int, errmsg string) {
	r := &Response{
		ErrNo:  errno,
		ErrMsg: errmsg,
	}
	b, _ := json.Marshal(r)
	fmt.Fprintf(w, string(b))
}

func staticHandler(w http.ResponseWriter, r *http.Request) {
	resp := &Response{
		ErrNo:  10000,
		ErrMsg: "",
	}

	s := &handler.Storage{}
	var data interface{} = nil
	switch r.Method {
	case "GET":
		fid := r.URL.Query().Get("fid")
		bucketName := r.URL.Query().Get("bucketname")
		if fid == "" || bucketName == "" {
			http.Error(w, "Bad Request", 400)
			return
		}

		d, err := s.DownloadFile(fid, bucketName)
		if err != nil {
			resp.ErrNo = 10001
			resp.ErrMsg = "failed to download"
			b, _ := json.Marshal(resp)
			fmt.Fprintf(w, "%s", string(b))
			return
		}

		data = d
	case "POST":
		fname := r.URL.Query().Get("fname")
		fid := r.URL.Query().Get("fid")
		bucketName := r.URL.Query().Get("bucketname")

		if fname == "" || fid == "" || bucketName == "" {
			http.Error(w, "Bad Request", 400)
			return
		}

		b, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Bad Request", 400)
			return
		}

		err = s.UploadFile(fname, fid, bucketName, b)
		if err != nil {
			resp.ErrNo = 10002
			resp.ErrMsg = "failed to upload"
			b, _ := json.Marshal(resp)
			fmt.Fprintf(w, "%s", string(b))
			return
		}

		data = fid
	case "DELETE":
		fid := r.URL.Query().Get("fid")
		bucketName := r.URL.Query().Get("bucketname")
		if fid == "" || bucketName == "" {
			http.Error(w, "Bad Request", 400)
			return
		}

		err := s.DeleteFile(fid, bucketName)
		if err != nil {
			resp.ErrNo = 10003
			resp.ErrMsg = "failed to delete"
			b, _ := json.Marshal(resp)
			fmt.Fprintf(w, "%s", string(b))
			return
		}
		data = "done"
	case "HEAD":
		fid := r.URL.Query().Get("fid")
		bucketName := r.URL.Query().Get("bucketname")
		if fid == "" || bucketName == "" {
			http.Error(w, "Bad Request", 400)
			return
		}

		d, err := s.GetMeta(fid, bucketName)
		if err != nil {
			resp.ErrNo = 10004
			resp.ErrMsg = "failed to head"
			b, _ := json.Marshal(resp)
			fmt.Fprintf(w, "%s", string(b))
			return
		}
		data = d
	default:
		http.Error(w, "Method not allowed", 405)
		return
	}

	resp.Data = data
	b, _ := json.Marshal(resp)
	fmt.Fprintf(w, "%s", string(b))
}

func initHandler(w http.ResponseWriter, r *http.Request) {
	resp := &Response{
		ErrNo:  10000,
		ErrMsg: "",
	}

	s := &handler.Storage{}
	var data interface{} = nil
	switch r.Method {
	case "POST":
		fid := r.URL.Query().Get("fid")

		if fid == "" {
			http.Error(w, "Bad Request", 400)
			return
		}

		data = s.InitUploadUrl(fid)
	default:
		http.Error(w, "Method not allowed", 405)
		return
	}

	resp.Data = data
	b, _ := json.Marshal(resp)
	fmt.Fprintf(w, "%s", string(b))
}

func bucketHandler(w http.ResponseWriter, r *http.Request) {
	resp := &Response{
		ErrNo:  10000,
		ErrMsg: "",
	}

	bucket := &handler.Bucket{}
	var data interface{} = nil
	switch r.Method {
	case "GET":
		bucketName := r.URL.Query().Get("bucketname")
		if bucketName == "" {
			http.Error(w, "Bad Request", 400)
			return
		}

		d, err := bucket.ListBucket(bucketName)
		if err != nil {
			resp.ErrNo = 10005
			resp.ErrMsg = "failed to create bucket"
			b, _ := json.Marshal(resp)
			fmt.Fprintf(w, "%s", string(b))
			return
		}
		data = d
	case "POST":
		bucketName := r.URL.Query().Get("bucketname")
		if bucketName == "" {
			http.Error(w, "Bad Request", 400)
			return
		}

		err := bucket.CreateBucket(bucketName)
		if err != nil {
			resp.ErrNo = 10006
			resp.ErrMsg = "failed to create bucket"
			b, _ := json.Marshal(resp)
			fmt.Fprintf(w, "%s", string(b))
			return
		}
		data = "done"
	case "DELETE":
		bucketName := r.URL.Query().Get("bucketname")
		if bucketName == "" {
			http.Error(w, "Bad Request", 400)
			return
		}

		err := bucket.DeleteBucket(bucketName)
		if err != nil {
			resp.ErrNo = 10007
			resp.ErrMsg = "failed to delete bucket"
			b, _ := json.Marshal(resp)
			fmt.Fprintf(w, "%s", string(b))
			return
		}
		data = "done"
	default:
		http.Error(w, "Method not allowed", 405)
		return
	}

	resp.Data = data
	b, _ := json.Marshal(resp)
	fmt.Fprintf(w, "%s", string(b))
}

func main() {
	logConfigFile := flag.String("l", "./conf/log4go.xml", "Log config file")
	configFile := flag.String("c", "./conf/conf.json", "Config file")

	flag.Parse()

	log.LoadConfiguration(*logConfigFile)

	if err := misc.LoadConf(*configFile); err != nil {
		fmt.Printf("failed to load configure, Err:[%s]", err)
		os.Exit(1)
	}

	if err := misc.InitBackend(); err != nil {
		fmt.Printf("failed to init database, Err:[%s]", err)
	}

	http.HandleFunc("/api/v1/static", staticHandler)
	http.HandleFunc("/api/v1/init", initHandler)
	http.HandleFunc("/api/v1/bucket", bucketHandler)

	err := http.ListenAndServe(misc.Conf.Addr, nil)
	if err != nil {
		fmt.Printf("failed to ListenAndServe, Err:[%s]", err)
		os.Exit(1)
	}
}
