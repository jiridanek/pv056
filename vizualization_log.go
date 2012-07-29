package main

import (
	"fmt"
	"time"
	"os"
	"encoding/csv"
	"io"
	"bufio"
	"strings"
)

func main () {
	//input := make(chan []string,100000)
	input := read_csv("zaznamy_merged")
	done := make(chan bool,0)
	go write_log("log.log", input, done)
	// wait for it
	<-done
}

// from sequence_clicks.go
func select_click(record []string, mask string) bool {
	id := strings.Join(record[1:9],"")
	for i,b := range mask {
		if rune(id[i]) != b && b != '.' {
			return false
		}
	}
	return true
}

func read_csv(filename string) [][]string {
	lines := make([][]string,0)
	f, err := os.Open(filename)
		if err != nil {
			panic(err)
		}
		buf := bufio.NewReader(f)
		defer f.Close()
		fr := csv.NewReader(buf)
	
		// no header line
		for {
			dline, err := fr.Read()
			if err == io.EOF {
				break
			}
	
			if err != nil {
				panic(err)
			}
			//					    12345678
			if select_click(dline, ".......1") {
				lines = append(lines, dline)
			}
		}
		return lines
}

func write_log(filename string, input [][]string, done chan bool) {
	f, err := os.Create(filename)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	buf := bufio.NewWriter(f)
	initialcommit := make(map[string]bool)
	for _,v := range input {
		// more processing may be needed
		t, err := time.Parse("200601021504", v[10])
		if err != nil {
			panic(err)
		}
		timestamp := t.Unix()
		username := v[0]
		//file := "/agendy" + strings.TrimRight(v[9],"/")
		file := strings.TrimRight(v[9],"/") + "/spread/out/now"
		//color := 
		
		action := "A"
 		if initialcommit[file] {
			fmt.Fprintf(buf, "%v|%s|%s|%s\n", timestamp, username, action, strings.TrimRight(v[9],"/") + "/spreat/out/now")
 			action = "M"
 		}
		initialcommit[file] = true
		fmt.Fprintf(buf, "%v|%s|%s|%s\n", timestamp, username, action, file)
	}
	done <- true
}