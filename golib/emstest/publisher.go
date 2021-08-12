package main

import (
	"ems"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	_ "net/http/pprof"
	"reflect"
	"runtime"
	"strings"
	"time"
)

// hello world, the web server
func HelloServer(w http.ResponseWriter, req *http.Request) {
	io.WriteString(w, "hello, world!\n")
}
func bigBytes() *[]byte {
	s := make([]byte, 100000000)
	return &s
}

var connection ems.Connection
var connections map[int]ems.Connection
var r1 *rand.Rand

func main() {

	var emsServerUrl string = "tcp://emsa-perf-t01.nordstrom.net:7840"
	var userName string = "merqry"
	var password string = "merqry"
	connections = make(map[int]ems.Connection)
	for i := 0; i < 10; i++ {
		connection = ems.CreateConnection(emsServerUrl, userName, password)
		connections[i] = connection
	}
	s1 := rand.NewSource(time.Now().UnixNano())
	r1 = rand.New(s1)
	fmt.Println(reflect.TypeOf(connection))
	http.HandleFunc("/perf", task)
	http.HandleFunc("/publish/merqry.perf", publish)
	log.Fatal(http.ListenAndServe(":8090", nil))
	//var wg sync.WaitGroup

	// go func() {
	// 	log.Println(http.ListenAndServe("localhost:6060", nil))
	// }()

}
func task(w http.ResponseWriter, r *http.Request) {
	fmt.Println("I am runnning task.")
	var mem runtime.MemStats
	runtime.ReadMemStats(&mem)
	fmt.Fprintf(w, "Alloc %d\n", mem.Alloc/1024/1024)
	fmt.Fprintf(w, "TotalAlloc %d\n", mem.TotalAlloc/1024/1024)
	fmt.Fprintf(w, "Sys %d\n", mem.Sys/1024/1024)
	fmt.Fprintf(w, "Lookups %d\n", mem.Lookups/1024/1024)
	fmt.Fprintf(w, "Mallocs %d\n", mem.Mallocs/1024/1024)
	fmt.Fprintf(w, "Frees %d\n", mem.Frees/1024/1024)
	fmt.Fprintf(w, "HeapAlloc %d\n", mem.HeapAlloc/1024/1024)
	fmt.Fprintf(w, "HeapSys %d\n", mem.HeapSys/1024/1024)
	fmt.Fprintf(w, "HeapIdle %d\n", mem.HeapIdle/1024/1024)
	fmt.Fprintf(w, "HeapInuse %d\n", mem.HeapInuse/1024/1024)
	fmt.Fprintf(w, "HeapReleased %d\n", mem.HeapReleased/1024/1024)
	fmt.Fprintf(w, "HeapObjects %d\n", mem.HeapObjects/1024/1024)
	fmt.Fprintf(w, "StackInuse %d\n", mem.StackInuse/1024/1024)
	fmt.Fprintf(w, "StackSys %d\n", mem.StackSys/1024/1024)
	fmt.Fprintf(w, "MSpanInuse %d\n", mem.MSpanInuse/1024/1024)
	fmt.Fprintf(w, "MSpanSys %d\n", mem.MSpanSys/1024/1024)
	fmt.Fprintf(w, "MCacheInuse %d\n", mem.MCacheInuse/1024/1024)
	fmt.Fprintf(w, "MCacheSys %d\n", mem.MCacheSys/1024/1024)
	fmt.Fprintf(w, "BuckHashSys %d\n", mem.BuckHashSys/1024/1024)
	fmt.Fprintf(w, "GCSys %d\n", mem.GCSys/1024/1024)
	fmt.Fprintf(w, "OtherSys %d\n", mem.OtherSys/1024/1024)
	fmt.Fprintf(w, "NextGC %d\n", mem.NextGC/1024/1024)
	fmt.Fprintf(w, "LastGC %d\n", mem.LastGC/1024/1024)
	fmt.Fprintf(w, "PauseTotalNs %d\n", mem.PauseTotalNs/1024/1024)
	fmt.Fprintf(w, "NumGC %d\n", mem.NumGC/1024/1024)
	fmt.Fprintf(w, "NumForcedGC %d\n", mem.NumForcedGC/1024/1024)
}

// hello() is home handler.
func publish(w http.ResponseWriter, r *http.Request) {
	start := time.Now().UnixNano() / int64(time.Millisecond)
	switch r.Method {
	case "GET":
		fmt.Fprintf(w, "Sorry,   POST methods are supported.")
	case "POST":
		// Call ParseForm() to parse the raw query and update r.PostForm and r.Form.
		if err := r.ParseForm(); err != nil {
			fmt.Fprintf(w, "Hello, POT method. ParseForm() err: %v", err)
			return
		}
		Send(r)

		end := time.Now().UnixNano() / int64(time.Millisecond)
		fmt.Fprintln(w, "{\"status\": \"success\"} ", end-start)
		//fmt.Printf("%d", syscall.Gettid())
	//	time.Sleep(55 * time.Second)
	default:
		fmt.Fprintf(w, "Sorry, only GET and POST methods are supported.")
	}
}
func Send(req *http.Request) {
	bodyByte, err := ioutil.ReadAll(req.Body)
	defer req.Body.Close()
	queuename := strings.SplitAfter(req.RequestURI, "publish/")[1]
	env := "SIT"
	qlist := []string{queuename, env}
	queuenamesit := strings.Join(qlist, ".")

	if err != nil {
		panic(err)
	}
	body := string(bodyByte)
	index := r1.Intn(10)
	fmt.Fprintf(ioutil.Discard, "", index, body, queuenamesit)
	var session = connections[index].CreateSession(false, 0)
	var messageProducer = session.CreateQueueProducer(queuenamesit)
	var txtMessage = ems.CreateTextMessage(body)
	messageProducer.Send(txtMessage)
	messageProducer.Close()
	txtMessage.Destroy()
	session.Close()
}
