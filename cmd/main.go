package main

import (
	"encoding/json"
	"fmt"
	"github.com/rakyll/statik/fs"
	"html/template"
	apiv1 "k8s.io/api/core/v1"
	"kubedemo/pkg/podinfo"
	_ "kubedemo/statik"
	"log"
	"net"
	"net/http"
	"os"
	u "os/user"
	"runtime"
	"sort"
	"strings"
	"time"
)

const Version string = "1.0"
const defaultPort = "8080"

type Data struct {
	Version   string
	GoVersion string
	StartTime time.Time
	PodName   string
	Port      string
	PodIP     string
	ClientIP  string
	NodeName  string
	NodeIP    string
	Request   []string
	Env       map[string]string
}

var (
	username string
	err      error
	podInfo  *apiv1.Pod
	data     Data
)

func main() {

	statikFS, err := fs.New()
	if err != nil {
		log.Fatal(err)
	}

	htmlindex, err := fs.ReadFile(statikFS, "/index.html")
	if err != nil {
		panic(err)
	}

	htmlenv, err := fs.ReadFile(statikFS, "/env.html")
	if err != nil {
		panic(err)
	}

	templ, err := template.New("index.html").Parse(string(htmlindex))
	templates := template.Must(templ.New("env.html").Parse(string(htmlenv)))

	data := new(Data)
	data.Version = Version
	data.StartTime = time.Now()
	data.Port = defaultPort
	// $PORT environment variable is provided in the Kubernetes deployment.
	if p := os.Getenv("PORT"); p != "" {
		data.Port = p
	}

	data.GoVersion = runtime.Version()
	data.PodName, _ = os.Hostname()
	podIPL, _ := getLocalIP()
	data.PodIP = strings.Join(podIPL, ", ") // Kubernetes dışında çalıştığı durumlad için.

	podInfo, err = podinfo.GetPodInfo(data.PodName)
	if err == nil {
		data.NodeIP = podInfo.Status.HostIP
		data.StartTime = podInfo.Status.StartTime.Time
		data.NodeName = podInfo.Spec.NodeName
		data.PodIP = podInfo.Status.PodIP

	} else {
		log.Printf("Pod Info Err.: %s\n", err)
	}

	user, err := u.Current()
	if err == nil {
		username = user.Username
	}

	data.Env = make(map[string]string)
	data.Env = getenvironment(os.Environ(), func(item string) (key, val string) {
		splits := strings.Split(item, "=")
		key = splits[0]
		val = splits[1]
		return
	})

	// http://patorjk.com/software/taag/#p=display&f=Doom&t=KubeDemo

	fmt.Println(" _   __      _         ______                     ")
	fmt.Println("| | / /     | |        |  _  \\                    ")
	fmt.Println("| |/ / _   _| |__   ___| | | |___ _ __ ___   ___  ")
	fmt.Println("|    \\| | | | '_ \\ / _ \\ | | / _ \\ '_ ` _ \\ / _ \\ ")
	fmt.Println("| |\\  \\ |_| | |_) |  __/ |/ /  __/ | | | | | (_) |")
	fmt.Printf("\\_| \\_/\\__,_|_.__/ \\___|___/ \\___|_| |_| |_|\\___/  v%s\n", Version)
	fmt.Println()
	fmt.Printf("Go Version: %s\n", data.GoVersion)
	fmt.Printf("Start Time: %v\n", data.StartTime)
	fmt.Printf("User      : %s\n", username)
	fmt.Printf("Port      : %s\n", data.Port)
	fmt.Printf("Pod Name  : %s\n", data.PodName)
	fmt.Printf("Pod Ip    : %s\n", data.PodIP)
	fmt.Printf("Node Name : %s\n", data.NodeName)
	fmt.Printf("Node Ip   : %s\n", data.NodeIP)

	var keys []string
	for k := range data.Env {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	/*
		fmt.Printf("Env. Var. : \n")
		for _, k := range keys {
			log.Printf("  %s =  %s \n", k, env[k])
		}
	*/

	http.HandleFunc("/api", func(w http.ResponseWriter, r *http.Request) {

		ip, _, _ := net.SplitHostPort(r.RemoteAddr)
		data.ClientIP = ip
		data.Request = fmtRequest(r)
		sort.Strings(data.Request)

		json.NewEncoder(w).Encode(data)
		return
	})

	http.HandleFunc("/env", func(w http.ResponseWriter, r *http.Request) {
		ip, _, _ := net.SplitHostPort(r.RemoteAddr)
		data.ClientIP = ip
		data.Request = fmtRequest(r)
		sort.Strings(data.Request)
		err := templates.ExecuteTemplate(w, "env.html", data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		ip, _, _ := net.SplitHostPort(r.RemoteAddr)

		data.ClientIP = ip
		data.Request = fmtRequest(r)
		sort.Strings(data.Request)

		err := templates.ExecuteTemplate(w, "index.html", data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	})

	if err := http.ListenAndServe(":"+data.Port, nil); err != nil {
		panic(err)
	}
}
