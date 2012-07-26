package main

import (
  "fmt"
  "log"
  "sort"
  "encoding/gob"
  "os"
  "flag"
  "runtime/pprof"
)

type BagsAgends [][]string

var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")

func main () {
  
  flag.Parse()
  if *cpuprofile != "" {
    cpuprofile, err := os.Create(*cpuprofile)
    if err != nil {
        log.Fatal(err)
    }
   pprof.StartCPUProfile(cpuprofile)
   defer pprof.StopCPUProfile()
  }
  
  f, err := os.Open("bag.gob")
  if err != nil {
    panic(err)
  }
  log.Println("Decoding")
  dec := gob.NewDecoder(f)
  var bags BagsAgends
  err = dec.Decode(&bags)
  if err != nil {
    panic(err)
  }
  log.Println("Adding to DB")
  tdb := NewTransactionDB()
  for _,bag := range bags {
      tdb.Append(bag)
  }
  log.Println("Computing C1")
  c1 := tdb.generateC1(1)
  fmt.Println("c1",len(c1))
  
  cprev := c1
  for k := 2; k < 5; k++ {
    log.Printf("Computing C%v\n", k)
    ck := tdb.generateC(k, cprev, 1000)
    
  fmt.Println("ck",len(ck))
  fmt.Println("ck",ck)
  cprev = ck
  }
  
//   tdb := NewTransactionDB()
//   tdb.Append([]string{"A", "B", "C"})
//   tdb.Append([]string{"A", "B", "D"})
//   fmt.Println(tdb)
//   c1 := tdb.generateC1(2)
//   fmt.Println(c1)
//   c2 := tdb.generateC(2, c1, 3)
//   fmt.Println(c2)
log.Println("Done")
}

type TransactionDB struct {
//  next_id int
  next_row int
//  string_to_id map[string]int
  tid_lists map[string][]int
}

func NewTransactionDB() *TransactionDB {
  tdb := &TransactionDB{
//    0,
    0,
//    make(map[string]int, 0),
    make(map[string][]int,0),
  }
  return tdb
}

// func (t *TransactionDB) nextId () int {
//   return t.next_id++
// }

func (t *TransactionDB) nextRow () int {
  row := t.next_row
  t.next_row++
  return row
}

func (t *TransactionDB) Append (items []string) {
  row := t.nextRow()
  for _, item := range items {
//     id, found := t.string_to_id[item]
//     if !found {
//       id = t.nextId()
//       t.string_to_id[item] = id
//     }
    list, found := t.tid_lists[item]
    if !found {
      list = make([]int,0)
    }
    list = append(list, row)
    t.tid_lists[item] = list
  } 
}

// not necesarry, rows are numbered in ascending order
func (t *TransactionDB) Sort () {
  for _,itemlist := range t.tid_lists {
    sort.Ints(itemlist)
  }
}

func (t *TransactionDB) generateC1(min_support int) [][]string{
  c1 := make([][]string,0)
  for k,v := range t.tid_lists {
    if len(v) >= min_support {
      c1 = append(c1, []string{k})
    }
  }
  return c1
}

func (t *TransactionDB) gen_check_candidate(iv, jv []string, min_support int) (candidate []string, admissable bool) {
  set := make(map[string]int,0)
  for _, item := range iv {
    set[item]++
  }
  for _, item := range jv {
    set[item]--
  }
  
  ivplus := ""
  jvplus := ""
  intersection := make([]string, 0)
  
  for key,val := range set {
    if val == -1 {
      if jvplus != "" {
	return make([]string,0), false
      }
      jvplus = key
      continue
    }
    if val == 1 {
      if ivplus != "" {
	return make([]string,0), false
      }
      ivplus = key
      continue
    }
    if val == 0 {
      intersection = append(intersection, key)
      continue
    }
    return make([]string,0), false
  }
  
  if ivplus == "" || jvplus == "" {
    return make([]string,0), false
  }
  
  double := make([]string,0,2)
  double = append(double, ivplus, jvplus)
  if t.support(double) < min_support {
    return make([]string,0), false
  }
  candidate = append(intersection, double...)
  return candidate, true
}

func (t *TransactionDB) generateC(k int, prev_c [][]string, min_support int) [][]string{
  ck := make(map[string][]string,0) // to remove duplicities
  if k < 2 {
    panic("k must be >= 2")
  }
  for i,iv := range prev_c {
    for _,jv := range prev_c[:i]{
      candidate, admissable := t.gen_check_candidate(iv, jv, min_support)
      if !admissable {
	continue
      }

      //check support
      if t.support(candidate) < min_support {
	continue
      }
      //now we are good
      sort.Strings(candidate)
      str := fmt.Sprintf("%v", candidate)
      ck[str] = candidate
    }
  }
  list := make([][]string, 0)
  for _,v := range ck{
    list = append(list, v)
  }
  return list
}

func (t *TransactionDB) union(seta, setb []string) []string{
  set := make(map[string]bool)
  for _, item := range seta {
    set[item] = true
  }
  for _, item := range setb {
    set[item] = true
  }
  list := make([]string,0,len(set))
  for key,_ := range set {
    list = append(list, key)
  }
  return list
}

func (t *TransactionDB) support(set []string) int{
  cnt := make(map[int]int)
  for _, item := range set {
    for _,tid := range t.tid_lists[item] {
      cnt[tid]++
    }
  }
  max := len(set)
  support := 0
  for _,v := range cnt {
    if v == max {
      support++
    }
  }
  return support
}