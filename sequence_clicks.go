package main

// Splits a csv into sessions, serialized into a .gob file as a []Clicks structure.
// Sessions from the same user are separated by 30 minutes or more of inactivity.

import (
	"encoding/csv"
	"encoding/gob"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"sort"
	"strings"
	"time"
	//	"strconv"

	. "./click"
)

// const (
// 	FAKE_UCO = iota
// 	TYP_APLIKACE
// 	DATUM_OPERACE
// 	NAZEV_DNE_OPERACE
// 	ADRESA_PRISLUSNOST
// 	ADRESA_CHECKSUM
// )

var fnamep = flag.String("incsv", "", "csv file to read, output of the merging script")

//																						  12345678 
var patternp = flag.String("pattern", "........", "patern to filter people by. Example: \"m1..1..0\". Defaults to no filtering.")
var outfnamep = flag.String("outgob", "", "a .gob file to write results into")

func main() {
	flag.Parse()
	if *fnamep == "" {
		log.Println("-incsv: No input filename given")
		return
	}
	if *outfnamep == "" {
		log.Println("-outgob: No output filename given")
		return
	}
	if len(*patternp) != 8 {
		log.Println("-pattern: Pattern is of wrong length, must be 8")
		return
	}

	log.Println("Reading " + *fnamep)

	list := read_csv(*fnamep, *patternp)
	log.Println("Processing")
	sort.Sort(ByIpFucoTimeId{list})
	sessions := split_sessions(list)

	fmt.Println("no of sessions:", len(sessions))
	log.Println("Serializing sessions to " + *outfnamep)
	serialize_sessions(*outfnamep, sessions)
	log.Println("Done")
}

func read_csv(fname, pattern string) Clicks {
	df, err := os.Open(fname)
	if err != nil {
		panic(err)
	}
	defer df.Close()
	dr := csv.NewReader(df)

	nextid := 0
	list := make(Clicks, 0)

	// no header line
	for {
		dline, err := dr.Read()
		if err == io.EOF {
			break
		}

		if err != nil {
			panic(err)
		}
		//                      
		if select_click(dline, pattern) {
			// we need to have fuco at [0] :(
			data_line := make([]string, 6)
			data_line[0] = dline[0] // FUCO
			copy(data_line[1:6], dline[9:])
			list = append(list, NewClickFromList(nextid, data_line))
		}

		//debug
		//  	  		if dline[9] == "/lide/" {
		//  	  		  break
		//  			}
	}
	return list
}

func select_click(record []string, mask string) bool {
	id := strings.Join(record[1:9], "")
	for i, b := range mask {
		if rune(id[i]) != b && b != '.' {
			return false
		}
	}
	return true
}

func serialize_sessions(filename string, sessions []Clicks) {
	f, err := os.Create(filename)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	e := gob.NewEncoder(f)
	e.Encode(sessions)
}

func test_split_sessions() []Clicks {
	list := make(Clicks, 0)
	list = append(list, &Click{0, 1, "a", time.Date(2012, time.November, 10, 23, 0, 1, 0, time.UTC), "so", "MU", "xxx"})
	list = append(list, &Click{1, 1, "a", time.Date(2012, time.November, 10, 23, 0, 2, 0, time.UTC), "so", "MU", "xxx"})
	list = append(list, &Click{2, 1, "a", time.Date(2012, time.November, 10, 23, 0, 3, 0, time.UTC), "so", "MU", "xxx"})
	list = append(list, &Click{3, 1, "a", time.Date(2012, time.November, 10, 23, 0, 4, 0, time.UTC), "so", "MU", "yyy"})
	sessions := split_sessions(list)
	return sessions
}

// split clickstream into sessions
// 30 minute interval
func split_sessions(list Clicks) []Clicks {
	sessions := make([]Clicks, 0)

	var session Clicks
	for i, v := range list {
		if i != 0 { // there is a previous record
			previous := list[i-1]
			if previous.Fake_uco != v.Fake_uco ||
				previous.Adresa_checksum != v.Adresa_checksum ||
				v.Datum_operace.Sub(previous.Datum_operace).Minutes() >= 30.0 {
				if session != nil {
					sessions = append(sessions, session)
				}
				session = make(Clicks, 0)
			}
		}
		session = append(session, v)
	}
	sessions = append(sessions, session)
	return sessions
}
