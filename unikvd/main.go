// Command unikvd is a simple daemon for unikv
package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/apiles/unikv"
	_ "github.com/apiles/unikv/drivers"
)

var namespaces map[string]*unikv.Namespace

func main() {
	if len(os.Args) != 2 {
		fmt.Println(`USAGE: unikvd [LISTEN_ADDR]`)
		os.Exit(0)
	}
	reinit()
	defer cleanup()
	http.HandleFunc("/v1/", handle)
	fmt.Println("Listening on: ", os.Args[1])
	err := http.ListenAndServe(os.Args[1], nil)
	if err != nil {
		panic(err)
	}
}

func cleanup() {
	for k := range namespaces {
		namespaces[k].Close()
	}
}

func reinit() {
	stats.totalRequests = 0
	stats.startTime = time.Now()
	namespaces = make(map[string]*unikv.Namespace)
	for _, nsName := range unikv.GetConfigure().GetNamespaceList() {
		namespaces[nsName] = unikv.NewNamespace(nsName)
		nsConfig, _ := unikv.GetConfigure().GetNamespace(nsName)
		for _, bucketName := range nsConfig.GetBucketList() {
			_, err := namespaces[nsName].NewBucket(bucketName)
			if err != nil {
				panic(err)
			}
		}
	}
}

func resp404(rw http.ResponseWriter) {
	rw.WriteHeader(404)
	rw.Write([]byte(`{"status":"error","result":"Not Found"}`))
}

func resp400(rw http.ResponseWriter) {
	rw.WriteHeader(400)
	rw.Write([]byte(`{"status":"error","result":"Bad Request"}`))
}

func resp500(rw http.ResponseWriter, msg string) {
	rw.WriteHeader(500)
	rw.Write([]byte(fmt.Sprintf(`{"status":"error","result":"%s"}`, msg)))
}

// Result is result
type Result struct {
	Status string      `json:"status"`
	Result interface{} `json:"result"`
}

func resp200(rw http.ResponseWriter, rslt interface{}) {
	r := &Result{
		Status: "success",
		Result: rslt,
	}
	b, _ := json.Marshal(r)
	rw.Write(b)
}

var stats struct {
	totalRequests int
	startTime     time.Time
}

func handle(rw http.ResponseWriter, r *http.Request) {
	stats.totalRequests++
	log.Printf("Request %s %s\r\n", r.Method, r.RequestURI)
	rw.Header().Add("Server", "UniKVd")
	uri := strings.Split(r.URL.RequestURI(), "?")[0]
	sp := strings.Split(strings.Trim(uri, "/"), "/")
	if len(sp) != 3 {
		resp404(rw)
		return
	}
	nsName, bucketName, key := sp[0], sp[1], sp[2]
	if nsName == "_system" {
		switch bucketName {
		case "lifecycle":
			switch key {
			case "shutdown":
				cleanup()
				resp200(rw, "OK")
				log.Println("Service shutdown")
				go func() {
					time.Sleep(200 * time.Millisecond)
					os.Exit(0)
				}()
			case "restart":
				cleanup()
				resp200(rw, "OK")
				log.Println("Service restarted")
				unikv.ReloadConfigure()
				reinit()
			default:
				resp404(rw)
			}
		case "stats":
			switch key {
			case "total_requests":
				resp200(rw, stats.totalRequests)
			case "configure":
				resp200(rw, unikv.GetConfigure())
			case "start_time":
				resp200(rw, stats.startTime.String())
			case "uptime":
				resp200(rw, time.Now().Sub(stats.startTime).String())
			case "all":
				var all struct {
					TotalRequests int              `json:"total_requests"`
					StartTime     string           `json:"start_time"`
					UpTime        string           `json:"uptime"`
					Configure     *unikv.Configure `json:"configure"`
				}
				all.TotalRequests = stats.totalRequests
				all.Configure = unikv.GetConfigure()
				all.StartTime = stats.startTime.String()
				all.UpTime = time.Now().Sub(stats.startTime).String()
				resp200(rw, all)
			default:
				resp404(rw)
			}
		default:
			resp404(rw)
		}
		return
	}
	ns, ok := namespaces[nsName]
	if !ok {
		resp404(rw)
		return
	}
	bucket, err := ns.NewBucket(bucketName)
	if err != nil {
		resp404(rw)
		return
	}
	switch r.Method {
	case "GET":
		rslt, err := bucket.GetString(key)
		if err == unikv.ErrNotFound {
			resp404(rw)
			return
		}
		if err != nil {
			resp500(rw, err.Error())
			return
		}
		resp200(rw, rslt)
	case "POST":
		d, err := ioutil.ReadAll(r.Body)
		if err != nil {
			resp500(rw, err.Error())
			return
		}
		err = bucket.PutString(key, string(d))
		if err != nil {
			resp500(rw, err.Error())
			return
		}
		resp200(rw, "OK")
	case "DELETE":
		err = bucket.Unset(key)
		if err != nil {
			resp500(rw, err.Error())
			return
		}
		resp200(rw, "OK")
	default:
		resp400(rw)
	}
}
