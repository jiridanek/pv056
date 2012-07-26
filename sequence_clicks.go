package main

import (
	"encoding/csv"
	"io"
	"os"
	"fmt"
	"log"
	"sort"
	"time"
	"encoding/gob"
//	"strings"
//	"strconv"
	. "./click"
)

const (
	FAKE_UCO = iota
	TYP_APLIKACE
	DATUM_OPERACE
	NAZEV_DNE_OPERACE
	ADRESA_PRISLUSNOST
	ADRESA_CHECKSUM
)

func main() {

		log.Println("reading zaznamy_export_data")
		df, err := os.Open("zaznamy_export_data")
		if err != nil {
			panic(err)
		}
		defer df.Close()
		dr := csv.NewReader(df)
	
		list := make(Clicks,0)
	
		dr.Read() // skip header
		for {
			dline, err := dr.Read()
			if err == io.EOF {
				break
			}
	
			if err != nil {
				panic(err)
			}
	
			list = append(list, NewClickFromList(dline))
			
			//debug
	//  		if key == "/lide/" {
	//  		  break
	//  }
		}
		
        log.Println("processing")
        sort.Sort(ByIpFucoTime{list})
	sessions := split_sessions(list)
	

	fmt.Println("sessions:", len(sessions))
	log.Println("serializing sessions")
	serialize_sessions("sessions.gob", sessions)
	log.Println("finished")
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
	list = append(list, &Click{"1", "a", time.Date(2012, time.November, 10, 23, 0, 1, 0, time.UTC), "so", "MU", "xxx"})
	list = append(list, &Click{"1", "a", time.Date(2012, time.November, 10, 23, 0, 2, 0, time.UTC), "so", "MU", "xxx"})
	list = append(list, &Click{"1", "a", time.Date(2012, time.November, 10, 23, 0, 3, 0, time.UTC), "so", "MU", "xxx"})
	list = append(list, &Click{"1", "a", time.Date(2012, time.November, 10, 23, 0, 4, 0, time.UTC), "so", "MU", "yyy"})
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