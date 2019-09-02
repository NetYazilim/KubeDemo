package main

import (
	"fmt"
	"net"
	"net/http"
)

func getLocalIP() ([]string, error) {
	var ipList []string
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ipList, err
	}
	for _, address := range addrs {
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				ipList = append(ipList, ipnet.IP.String())
			}
		}
	}

	if len(ipList) > 0 {
		return ipList, err
	}
	return ipList, fmt.Errorf("non loopback ip address not found")
}

func fmtRequest(r *http.Request) []string {
	var request []string
	url := fmt.Sprintf("%v %v %v", r.Method, r.URL, r.Proto)
	request = append(request, url)
	request = append(request, fmt.Sprintf("Host: %v", r.Host))
	for name, headers := range r.Header {
		for _, h := range headers {
			request = append(request, fmt.Sprintf("%v: %v", name, h))
		}
	}
	if r.Method == "POST" {
		r.ParseForm()
		request = append(request, r.Form.Encode())
	}
	return request
}
func getenvironment(data []string, getkeyval func(item string) (key, val string)) map[string]string {
	items := make(map[string]string)
	for _, item := range data {
		key, val := getkeyval(item)
		items[key] = val
	}
	return items
}
