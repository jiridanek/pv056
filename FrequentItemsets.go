package main

import (
  "fmt"
  "log"
  //"sort"
  "encoding/gob"
  "os"
  //"flag"
  //"runtime/pprof"
) 

type BagsAgends [][]string

type TransactionDB struct {
  nexttid Id
  nextiid Id
  toIid map[string]Id
  fromIid map[Id]string
  tids TidLists
}

func NewTransactionDB () *TransactionDB{
  return &TransactionDB{nexttid: 0,
  nextiid: 0, toIid: make(map[string]Id),
			  fromIid: make(map[Id]string),
					tids: make(TidLists,0)}
}

func (t *TransactionDB) Tids(id Id) TidList {
  return t.tids[id]
}

func (t *TransactionDB) MaxTid() Id {
  return Id(t.nexttid - 1)
}

func (t *TransactionDB) MaxIid() Id {
  return Id(t.nextiid - 1)
}

func (t *TransactionDB) Item(iid Id) string {
  return t.fromIid[iid]
}

// all strings in transaction have to be unique
func (t *TransactionDB) Append (transaction []string) {
  tid := t.nextTid()
  for _,item := range transaction {
    iid := t.Iid(item)
    // uniquenes check
    if len(t.tids[iid]) != 0 {
      last := t.tids[iid][len(t.tids[iid])-1]
      if last == tid {
	panic("an element is more than once in a transaction")
      }
    }
    t.tids[iid] = append(t.tids[iid], tid)
  }
}

func (t *TransactionDB) Iid (item string) Id {
  iid, found := t.toIid[item]
  if !found {
    iid = t.nextiid
    t.nextiid++
    t.toIid[item] = iid
    t.fromIid[iid] = item
    
    t.tids = append(t.tids, make(TidList,0))
  }
  return Id(iid)
}

func (t *TransactionDB) nextTid () Id {
  tid := t.nexttid
  t.nexttid++
  return tid
}

type Id int
type TidLists []TidList
type TidList []Id
func (ci TidList) IntersectedWith(cj TidList) TidList{
  p := 0
  q := 0
  
  intersection := make(TidList,0)
  for p < len(ci) && q < len(cj) {
    if ci[p] < cj[q] {
	p++
    } else if cj[q] < ci[p] {
      q++
    } else {
    	intersection = append(intersection, ci[p])
	p++
	q++
      }
  }
  return intersection
}


type Root struct {
  Children Nodes
}

type Node struct {
  V Id
  TidList TidList
  Children Nodes
}

func NewRoot() *Root{
  return &Root{make(Nodes,0)}
}

func NewNode(id Id, tidlist TidList) *Node{
  return &Node{id, tidlist, make(Nodes,0)}
}

func Run (tdb *TransactionDB, root *Root, min_support int) {  
  for i := Id(0); i <= tdb.MaxIid(); i++ {
    tidlist := tdb.Tids(i)
    if len(tidlist) < min_support {
      root.Children = append(root.Children, NewNode(i, make(TidList,0)))
      continue
    }
    root.Children = append(root.Children, NewNode(i, tidlist))
  }
  
    // from right to left
    for i := tdb.MaxIid(); i >= 0; i-- {
      run_recursive(tdb, root, root.Children[i], min_support)
    }
}

func run_recursive (tdb *TransactionDB, root *Root, node *Node, min_support int) {
  offset := node.V+1
  for i := node.V+1; i <= tdb.MaxIid(); i++ {
    tidlist := root.Children[i].TidList.IntersectedWith(node.TidList)
    
    if len(tidlist) < min_support {
      node.Children = append(node.Children, NewNode(i, make(TidList,0)))
      continue
    }
    node.Children = append(node.Children, NewNode(i, tidlist))
    run_recursive(tdb, root, node.Children[i-offset], min_support)
  }
}

type Nodes []*Node
func (s Nodes) Len() int      { return len(s) }
func (s Nodes) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
type ByValue struct{ Nodes }
func (s ByValue) Less(i, j int) bool { return s.Nodes[i].V < s.Nodes[j].V }

func PrintTree (tdb *TransactionDB, root *Root) {
  for _,node := range root.Children {
    PrintNode_recursive(tdb, node, 1)
  }
}

func tabs(i int) string {
  t := ""
  for j := 0; j < i; j++ {
    t +="\t"
  }
  return t
}

func PrintNode_recursive (tdb *TransactionDB, node *Node, level int) {
  fmt.Printf("%v%v:%v\n", tabs(level), tdb.Item(node.V), len(node.TidList))
  for _,n := range node.Children {
    PrintNode_recursive(tdb, n, level+1)
  }
}

func ComputeSets (tdb *TransactionDB, root *Root) []Itemsets {
  sets := make([]Itemsets,0)
  for _,n := range root.Children {
   	sets = ComputeSets_recursive(tdb, n, sets, []string{}, 0)
  }
  return sets
}
  
func ComputeSets_recursive(tdb *TransactionDB, node *Node, setssofar []Itemsets, before []string, level int) []Itemsets{
  // are we in a leef?
  if len(node.TidList) == 0 {
    return setssofar
  }
  
  if level > len(setssofar)-1 {
      //always the difference is only 0 or 1
      setssofar = append(setssofar, make(Itemsets,0))

  }
  newbefore := make([]string,len(before), len(before)+1)
  copy(newbefore, before)
  newbefore = append(newbefore, tdb.Item(node.V))
  
//   if level == 6 {
//     fmt.Println(tdb.Item(node.V))
//     fmt.Println(before)
//     fmt.Println(newbefore)
//   }

  setssofar[level] = append(setssofar[level], Itemset{Support: len(node.TidList), Items: newbefore})
  for _,n := range node.Children {
    setssofar = ComputeSets_recursive(tdb, n, setssofar, newbefore, level+1)
  }
//   if len(setssofar) >= 7 {
//     fmt.Println(setssofar[6])
//   }
  return setssofar
}

func test () []Itemsets {
       tdb := NewTransactionDB()
    tdb.Append([]string{"a", "b", "c"})
    tdb.Append([]string{"a", "b", "d"})
    tdb.Append([]string{"a", "b", "e"})
    tdb.Append([]string{"b", "c", "d"})
    
    root := NewRoot()
    Run(tdb, root, 2)
    sets := ComputeSets(tdb, root)
    PrintTree(tdb, root)
    return sets
}

func naostro () []Itemsets {
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
  
  root := NewRoot()
   Run(tdb, root, 1000)
   //PrintTree(tdb, root)
   sets := ComputeSets(tdb, root)
   return sets
}

func main () {
    log.Println("Started")

//     naostro()
    sets := naostro()
    //sets := test()
      
    
    for i,r := range sets {
      fmt.Println(i, len(r))
      fmt.Println(r)
    }
     
    
    log.Println("Finished")
}

type FrequentItemsets struct {
  Min_support int
  Data []Itemsets
}

type Itemsets []Itemset
type Itemset struct {
  Support int
  Items []string
}