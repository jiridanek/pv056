package main

import (
	. "./click"
	"encoding/gob"
	"flag"
	"fmt"
	"log"
	"os"
	"sort"
	"strings"
	"io"
)

var infnamep = flag.String("ingob", "", "input .gob file name. File is produced by a seqence click script.")
var outgobfnamep = flag.String("outgob", "", "output .gob file name.")
var outarfffnamep = flag.String("outarff", "", "output .arff file name.")

type BagsAgends [][]string

func main() {
	flag.Parse()
	if *infnamep == "" {
		log.Println("No input filename given")
		return
	}
	if *outgobfnamep == "" && *outarfffnamep == "" {
		log.Println("No output filename given")
		return
	}
	ses, err := os.Open(*infnamep)
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
	bags := make(BagsAgends, 0)
	setofallitems := make(map[string]bool, 0)
	for _, session := range sessions {
		set := make(map[string]bool, 0)
		for _, click := range session {
			var agenda string
			for _, seg := range strings.Split(click.Typ_aplikace, "/") {
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
		bag := make([]string, 0)
		for k, _ := range set {
			item := "TYP_AGENDY:"+k
			bag = append(bag, item)
			setofallitems[item] = true
		}
		sort.Strings(bag)
		bags = append(bags, bag)
	}

	//   fmt.Println(len(bags))
	//   fmt.Println(len(bags[0]))
	//   fmt.Println(len(bags[1]))
	//debugging output
	for _, row := range bags {
		fmt.Println(row)
	}

	if *outgobfnamep != "" {
		log.Println("Encoding gob")
		bagf, err := os.Create(*outgobfnamep)
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
	if *outarfffnamep != "" {
		log.Println("Encoding arff")
		bagf, err := os.Create(*outarfffnamep)
		defer bagf.Close()
		if err != nil {
			panic(err)
		}
		
		listofallitems := make([]string,0)
		for k,_ := range setofallitems {
			listofallitems = append(listofallitems, k)
		}
		sort.Strings(listofallitems)
		at2atnum := make(map[string]int)
		for i, item := range listofallitems {
			at2atnum[item] = i
		}

		io.WriteString(bagf, "@RELATION bagofagends\n")
		for _,item := range listofallitems {
			io.WriteString(bagf, "@ATTRIBUTE " + item + " {0,1}\n")
		}
		io.WriteString(bagf, "@DATA\n")
		for _,bag := range bags {
			//fmt.Fprintln(bagf, bag)
// 			j := 0
// 			for i, item := range listofallitems {
// 				if j < len(bag) && item == bag[j] {
// 					io.WriteString(bagf, "1")
// 					j++
// 				} else {
// 					io.WriteString(bagf, "0")
// 				}
// 				// last item
// 				if i == len(listofallitems) - 1 {
// 					io.WriteString(bagf, "\n")
// 				} else {
// 					io.WriteString(bagf, ",")
// 				}
// 			}
			// we want sparse ARFF for this
			io.WriteString(bagf, "{")
			for i,item := range bag {
				atnum, found := at2atnum[item]
				if !found {
					panic("")
				}
				fmt.Fprintf(bagf, "%d 1", atnum)
				if i != len(bag)-1 {
					io.WriteString(bagf, ",")
				}
			}
			io.WriteString(bagf, "}\n")
		}
		log.Println("Encoding done")
	}
}
