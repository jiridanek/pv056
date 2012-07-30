package main

import (
	"encoding/gob"
	"flag"
	"fmt"
	"log"
	"os"
	//	"io"
	"sort"
	"strings"
	"math"

	. "./click"
)

var lidigobnamep = flag.String("lidigob", "", "Path to gobbed associative array from fake uco to user record")
var relationgobnamep = flag.String("relationgob", "", "Path to a gob with clicks split into relations")
var outarffbnamep = flag.String("outarff", "", "Path for the reslut arff")

type Dvec map[string]int
type Markov map[string]Dvec

func (m Markov) TransitionMap(s string) Dvec {
	tm, found := m[s]
	if !found {
		tm = make(Dvec)
		m[s] = tm
	}
	return tm
}

type ArffCols map[string]bool

func (a ArffCols) Add(s string) {
	a[s] = true
}

func main() {

	flag.Parse()
	if *lidigobnamep == "" {
		fmt.Println("required: -lidigob, path to assoc array of people")
		return
	}
	if *relationgobnamep == "" {
		fmt.Println("required: -relationgob, path to relations gob")
		return
	}
	if *outarffbnamep == "" {
		fmt.Println("required: -outarff, path to write the resluts to")
		return
	}

	users := load_dblidi(*lidigobnamep)
	rels := load_clicks(*relationgobnamep)

	markovs := make(map[int]Markov)
	arffcols := make(ArffCols)
	arffcols.Add("@start")
	arffcols.Add("@end")

	log.Println("Processing")
	for _, rel := range rels {
		user, found := markovs[rel[0].Fake_uco]
		if !found {
			user = make(Markov)
			markovs[rel[0].Fake_uco] = user
		}
		for i := range rel {
			currapl := apl_name(rel[i].Typ_aplikace)
			arffcols.Add(currapl)
			var transition Dvec
			if i == 0 {
				transition = user.TransitionMap("@start")
			} else {
				prevapl := apl_name(rel[i-1].Typ_aplikace)
				transition = user.TransitionMap(prevapl)
			}

			transition[currapl]++

			if i == len(rel)-1 {
				transition = user.TransitionMap(currapl)
				transition["@end"]++
			}
		}
	}

	log.Println("Writing ARFF")
	if false {
		print_arff(arffcols, markovs, users)
	}
	print_distmatrix(arffcols, markovs)
}

func print_distmatrix(arffcols ArffCols, markovs map[int]Markov) {
	nousers := len(markovs)
	listusers := make([]int,0,nousers)
	listmarkovs := make([]Markov,0,nousers)
	for fuco,markov := range markovs {
		listusers = append(listusers, fuco)
		listmarkovs = append(listmarkovs, markov)
	}
	fmt.Println(nousers)
	for i,fuco1 := range listusers {
		valuesonline := 0
		fmt.Printf("%v", fuco1)
		for j,_ := range listusers {
			d := dist(arffcols, listmarkovs[i], listmarkovs[j])
			
			if valuesonline > 40 {
				fmt.Printf(" %.4f\n", d)
				valuesonline = 0
			} else {
				fmt.Printf(" %.4f", d)
				valuesonline++
			}
		}
		fmt.Printf("\n")
	}
}

func dist(arffcols ArffCols, m1, m2 Markov) float64{
	sum := 0.0
	set := make(map[string]bool)
	for a := range m1 {
		set[a] = true
	}
	for b := range m2 {
		set[b] = true
	}
	for agenda := range set {
		v1, f1 := m1[agenda]
		v2, f2 := m2[agenda]
		if !f1 && !f2 {
			continue
		}
		sum += distdvec(v1, v2)
	}
// 	for agenda,_ := range arffcols {
// 		v1, f1 := m1[agenda]
// 		v2, f2 := m2[agenda]
// 		if !f1 && !f2 {
// 			continue
// 		}
// 		sum += distdvec(v1, v2)
// 	}
	return sum
}

func distdvec(v1, v2 Dvec) float64 {
	set := make( map[string]int)
	for ag1,val1 := range v1 {
		set[ag1] += val1
	}
	for ag2, val2 := range v2 {
		set[ag2] -= val2
	}
	sum := 0.0
	for _,val := range set {
		sum += math.Sqrt(float64(val)*float64(val))
	}
	return sum
}

func print_arff(arffcols ArffCols, markovs map[int]Markov, users map[int][]string) {
	fmt.Println("@RELATION markov")
	fmt.Println("@ATTRIBUTE fuco STRING")
	fmt.Println("@ATTRIBUTE type {1011010,1101100,1001001,1111110,1011110,0001110,0001101,0001000,1101111,1001110,1001111,1011111,0000000,0001001,1001011,1001101,0001100,1001000,1001010,1111111,1101110,1101101,1011011,1001100,0001011,0000001,0001111,0001010}")
	fmt.Println("@ATTRIBUTE ucitel {0,1}")
	// to ensure we iterate it always in the same order
	collist := make([]string, 0)
	arffmapping := make(map[string]int)
	nextagend := 0
	for r, _ := range arffcols {
		collist = append(collist, r)
		arffmapping[r] = nextagend
		nextagend++
	}

	reducedarffmapping := make(map[string]int)
	nextdoubleagend := 0
	for _, markov := range markovs {
		for agenda, dvec := range markov {
			for agenda2, _ := range dvec {
				doubleagend := agenda + "|" + agenda2
				_, found := reducedarffmapping[doubleagend]
				if !found {
					fmt.Printf("@ATTRIBUTE %s NUMERIC\n", doubleagend)
					reducedarffmapping[doubleagend] = nextdoubleagend
					nextdoubleagend++
				}
			}
		}
	}

	log.Println("Attributes:", nextdoubleagend)

	// awfully lot of attributes here, hope Weka's EM can cope

	fmt.Println("@DATA")
	offset := 3
	for fuco, markov := range markovs {
		pairs := make(IntFloatPairs, 0)
		for agenda, dvec := range markov {
			normalizer := 0
			for _, value := range dvec {
				normalizer += value
			}
			for agenda2, value := range dvec {
				doubleagend := agenda + "|" + agenda2
				num := offset + reducedarffmapping[doubleagend]

				//log.Println("\n", offset, arffmapping[agenda], nextagend, arffmapping[agenda2], num)
				normal := float64(value) / float64(normalizer)
				pairs = append(pairs, &IntFloatPair{num, normal})
			}
		}
		sort.Sort(ByInt{pairs})
		fmt.Printf("{0 %v, 1 %v, 2 %v", fuco, strings.Join(users[fuco][2:9],""), users[fuco][8])
		for _, pair := range pairs {

			fmt.Printf(",%v %0.3f", pair.Int, pair.Float)
		}
		fmt.Printf("}\n")
	}
	//fmt.Println(markovs)
}

type IntFloatPairs []*IntFloatPair
type IntFloatPair struct {
	Int   int
	Float float64
}

func (s IntFloatPairs) Len() int      { return len(s) }
func (s IntFloatPairs) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

type ByInt struct{ IntFloatPairs }

func (s ByInt) Less(i, j int) bool { return s.IntFloatPairs[i].Int < s.IntFloatPairs[j].Int }

func apl_name(name string) string {
	name = strings.TrimSpace(name)
	name = strings.Trim(name, "/")
	subs := strings.SplitN(name, "/", 2)
	if len(subs) == 2 {
		subs = subs[:1]
	}
	return strings.Join(subs, "/")
}

func load_dblidi(filename string) map[int][]string{
	log.Println("Loading " + filename)
	f, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	dec := gob.NewDecoder(f)
	var users map[int][]string
	err = dec.Decode(&users)
	if err != nil {
		panic(err)
	}
	return users
}
func load_clicks(filename string) []Clicks {
	log.Println("Loading " + filename)
	f, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	dec := gob.NewDecoder(f)
	var clicks []Clicks
	err = dec.Decode(&clicks)
	if err != nil {
		panic(err)
	}
	return clicks
}
