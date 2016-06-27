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
		uid := r.URL.Query().Get("uid")
		fname := r.URL.Query().Get("fname")
		if uid == "" || fname == "" {
			http.Error(w, "Bad Request", 400)
			return
		}

		d, err := s.DownloadFile(uid, fname)
		if err != nil {
			resp.ErrNo = 10001
			resp.ErrMsg = "failed to download"
			b, _ := json.Marshal(resp)
			fmt.Fprintf(w, "%s", string(b))
			return
		}

		data = d
	case "POST":
		uid := r.URL.Query().Get("uid")
		fname := r.URL.Query().Get("fname")
		fid := r.URL.Query().Get("fid")

		if uid == "" || fname == "" || fid == "" {
			http.Error(w, "Bad Request", 400)
			return
		}

		b, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Bad Request", 400)
			return
		}

		err = s.UploadFile(uid, fname, fid, b)
		if err != nil {
			resp.ErrNo = 10002
			resp.ErrMsg = "failed to upload"
			b, _ := json.Marshal(resp)
			fmt.Fprintf(w, "%s", string(b))
			return
		}

		data = "done"
	case "DELETE":
		uid := r.URL.Query().Get("uid")
		fname := r.URL.Query().Get("fname")
		if uid == "" || fname == "" {
			http.Error(w, "Bad Request", 400)
			return
		}

		err := s.DeleteFile(uid, fname)
		if err != nil {
			resp.ErrNo = 10003
			resp.ErrMsg = "failed to delete"
			b, _ := json.Marshal(resp)
			fmt.Fprintf(w, "%s", string(b))
			return
		}
		data = "done"
	case "HEAD":
		uid := r.URL.Query().Get("uid")
		fid := r.URL.Query().Get("fid")
		if uid == "" || fid == "" {
			http.Error(w, "Bad Request", 400)
			return
		}

		d, err := s.GetMeta(uid, fid)
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
		uid := r.URL.Query().Get("uid")
		fsign := r.URL.Query().Get("fsign")

		if uid == "" || fsign == "" {
			http.Error(w, "Bad Request", 400)
			return
		}

		err = s.InitUploadUrl(uid, fsign)
		if err != nil {
			resp.ErrNo = 10004
			resp.ErrMsg = "failed to upload"
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
		fmt.Printf("failed to ListenAndServe: (%s)", err)
		os.Exit(1)
	}

	misc.InitBackend()

	http.HandleFunc("/api/v1/static", staticHandler)
	http.HandleFunc("/api/v1/init", initHandler)

	err := http.ListenAndServe(misc.Conf.Addr, nil)
	if err != nil {
		fmt.Printf("failed to ListenAndServe: (%s)", err)
		os.Exit(1)
	}
}
