package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"os"
	"strings"
	"syscall"
	"trie"
)

func main() {
	tree := trie.InitRoot()

	fi, err := os.Open("/Users/user/Downloads/keyword.csv")
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return
	}
	defer fi.Close()

	br := bufio.NewReader(fi)
	for {
		line, _, end := br.ReadLine()
		if end == io.EOF {
			break
		}
		info := strings.Split(string(line), ",")
		trie.AddKeyword(tree, info[1])
	}
	socket, err := net.Listen("unix", "/tmp/keyword_match.sock")
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return
	}
	defer syscall.Unlink("/tmp/keyword_match.sock")
	for {
		client, err := socket.Accept()
		if err != nil {
			fmt.Printf("Error: %s\n", err)
			return
		}
		buf := make([]byte, 1024)
		data_len, err := client.Read(buf)
		if err != nil {
			fmt.Printf("Error: %s\n", err)
			return
		}

		data := buf[0:data_len]
		msg := string(data)
		matched := trie.Match(tree, msg)
		response := []byte("[]") // 给响应一个默认值
		if len(matched) > 0 {
			json_str, _ := json.Marshal(matched)
			response = []byte(string(json_str))
		}
		_, err = client.Write(response)
		if err != nil {
			fmt.Printf("Error: %s\n", err)
			return
		}
	}
}
