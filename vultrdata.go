package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/kuangchanglang/graceful"
)

var (
	httpClient = &http.Client{
		Timeout: time.Second * 10,
	}
)

var (
	addr     = flag.String("addr", "0.0.0.0", "Listen address")
	port     = flag.Int("port", 8888, "Listen port")
	userdata = flag.Bool("userdata", false, "Also return 'userdata'")
	apiKey   = ""
)

func instanceHandler(w http.ResponseWriter, r *http.Request) {
	body, err := doApiRequest(w, r, "GET", "https://api.vultr.com/v1/server/list", nil)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Header().Set("content-type", "text/plain")
		w.Write([]byte(err.Error() + "\n"))
		return
	}

	var subidData map[string]map[string]interface{}
	err = json.Unmarshal(body, &subidData)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Header().Set("content-type", "text/plain")
		w.Write([]byte(err.Error() + "\n"))
		return
	}

	remoteAddr := strings.Split(r.RemoteAddr, ":")[0]

	found := findByIpAddress(subidData, remoteAddr)
	if found == nil {
		w.WriteHeader(http.StatusNotFound)
		w.Header().Set("content-type", "text/plain")
		w.Write([]byte(""))
		return
	}

	delete(found, "default_password")
	delete(found, "kvm_url")

	if *userdata {
		url := fmt.Sprintf("https://api.vultr.com/v1/server/get_user_data?SUBID=%v", found["SUBID"])
		body, err = doApiRequest(w, r, "GET", url, nil)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Header().Set("content-type", "text/plain")
			w.Write([]byte(err.Error() + "\n"))
			return
		}
		var kv map[string]string
		json.Unmarshal(body, &kv)
		if kv != nil {
			if v, ok := kv["userdata"]; ok {
				found["userdata"] = v
			}
		}
	}

	data, err := json.MarshalIndent(found, "", "  ")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Header().Set("content-type", "text/plain")
		w.Write([]byte(err.Error() + "\n"))
		return
	}

	w.Header().Set("content-type", "application/json")
	w.Write(data)
}

func doApiRequest(w http.ResponseWriter, r *http.Request, method, url string, body io.Reader) ([]byte, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("API-Key", apiKey)
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func findByIpAddress(subidData map[string]map[string]interface{}, remoteAddr string) map[string]interface{} {
	for _, data := range subidData {
		for _, name := range []string{"internal_ip", "main_ip", "v6_main_ip"} {
			value := data[name]
			if s, ok := value.(string); ok {
				if s == remoteAddr {
					return data
				}
			}
		}
	}
	return nil
}

func main() {
	flag.Parse()

	apiKey = os.Getenv("API_KEY")
	if apiKey == "" {
		fmt.Printf("Missing API_KEY environment variable.")
		os.Exit(1)
	}

	http.HandleFunc("/", instanceHandler)

	server := graceful.NewServer()

	server.Register(*addr+":"+strconv.Itoa(*port), http.DefaultServeMux)
	if graceful.IsWorker() {
		log.Printf("Listening on %v:%d\n", *addr, *port)
	}
	if err := server.Run(); err != nil {
		fmt.Printf("graceful.Server.Run: %v", err)
	}
}
