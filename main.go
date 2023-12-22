package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"strings"

	"github.com/oschwald/geoip2-golang"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		// Get the client's IP address using the basic approach (RemoteAddr)
		clientIP := getClientIPAddr(req)

		if clientIP == nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("500 - Something bad happened!"))
			return
		}

		db, err := geoip2.Open("/usr/share/GeoIP/GeoLite2-Country.mmdb")
		shouldReturn := WriteError(err, w)
		if shouldReturn {
			return
		}
		defer db.Close()

		record, err := db.Country(clientIP)

		shouldReturn = WriteError(err, w)
		if shouldReturn {
			return
		}

		ReturnIP(record, w, clientIP)
	})

	http.ListenAndServe(":8085", nil)
}

func ReturnIP(record *geoip2.Country, w http.ResponseWriter, clientIP net.IP) bool {
	if record.Country.IsoCode != "" {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "%s:%s\n", clientIP.String(), record.Country.IsoCode)
		return true
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "%s\n", clientIP.String())
	return false
}

func WriteError(err error, w http.ResponseWriter) bool {
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("500 - Something bad happened!"))
		return true
	}
	return false
}

func getClientIPAddr(req *http.Request) net.IP {
	log.Printf("%v", req.Header)
	// Check if X-Forwarded-For header exists
	clientIP := req.Header.Get("X-Real-Ip")
	if clientIP != "" {
		// Extract the first IP address from the comma-separated list
		clientIP = strings.Split(clientIP, ":")[0]
	} else {
		clientIP = req.Header.Get("X-Forwarded-For")
		if clientIP != "" {
			// If X-Forwarded-For is empty, fall back to RemoteAddr
			clientIP = strings.Split(req.RemoteAddr, ":")[0]
		} else {
			// Extract the first IP address from the comma-separated list
			clientIP = strings.Split(clientIP, ":")[0]
		}

	}
	return net.ParseIP(clientIP)
}
