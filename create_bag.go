package main

import (
  "encoding/gob"
  "os"
  "strings"
  "fmt"
  "log"
  . "./click"
)

type BagsAgends [][]string

func main () {
  ses, err := os.Open("sessions.gob")
  if err != nil {
    panic(err)
  }
  dec := gob.NewDecoder(ses)
  var sessions []Clicks
  log.Println("Decoding")
  err = dec.Decode(&sessions)
  log.Println("Decoding done")
  if err != nil {
    panic(err)
  }
  log.Println("Processing")
  bags := make(BagsAgends,0)
  for _,session := range sessions {
    set := make(map[string]bool,0)
    for _,click := range session {
      var agenda string
      for _,seg := range strings.Split(click.Typ_aplikace,"/") {
	if seg != "" {
	  agenda = seg
	  break
	}
      }
      if agenda == "" {
	continue
      }
      set[agenda] = true
    }
    bag := make([]string,0)
    for k,_ := range set {
      bag = append(bag, "TYP_AGENDY:"+k)
    }
    bags = append(bags, bag)
  }

//   fmt.Println(len(bags))
//   fmt.Println(len(bags[0]))
//   fmt.Println(len(bags[1]))
  //debugging output
   for _,row := range bags {
     fmt.Println(row)
   }
   
  log.Println("Encoding")
  bagf, err := os.Create("bag.gob")
  defer bagf.Close()
  if err != nil {
    panic(err)
  }
  enc := gob.NewEncoder(bagf)

  err = enc.Encode(bags)
  if err != nil {
    panic(err)
  }
  log.Println("Encoding done")
}
    
