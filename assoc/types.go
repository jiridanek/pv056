package assoc

import (
  "sort"
)

type BagsAgends [][]string

type FrequentItemsets struct {
  Min_support int
  TDB *TransactionDB
  Tree *Root
  Data []Itemsets
}

type Itemsets []Itemset
type Itemset struct {
  Support int
  Items []string
} 

func (fisets *FrequentItemsets) Sort() {
  for _,isets := range fisets.Data {
    for _,iset := range isets {
      sort.Strings(iset.Items)
    }
  }
}

///////////////////////////////////////////////

type TransactionDB struct {
  Nexttid Id
  nextiid Id
  ToIid map[string]Id
  FromIid map[Id]string
  tids TidLists
}

func NewTransactionDB () *TransactionDB{
  return &TransactionDB{Nexttid: 0,
  nextiid: 0,
  ToIid: make(map[string]Id),
			  FromIid: make(map[Id]string),
					tids: make(TidLists,0)}
}

func (t *TransactionDB) Tids(id Id) TidList {
  return t.tids[id]
}

func (t *TransactionDB) MaxTid() Id {
  return Id(t.Nexttid - 1)
}

func (t *TransactionDB) MaxIid() Id {
  return Id(t.nextiid - 1)
}

func (t *TransactionDB) Item(iid Id) string {
  return t.FromIid[iid]
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
  iid, found := t.ToIid[item]
  if !found {
    iid = t.nextiid
    t.nextiid++
    t.ToIid[item] = iid
    t.FromIid[iid] = item
    
    t.tids = append(t.tids, make(TidList,0))
  }
  return Id(iid)
}

func (t *TransactionDB) nextTid () Id {
  tid := t.Nexttid
  t.Nexttid++
  return tid
}

type Id int
type TidLists []TidList
type TidList []Id
func (s TidList) Len() int      { return len(s) }
func (s TidList) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
func (s TidList) Less(i, j int) bool { return s[i] < s[j] }
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

///////////////////////////////////////////////

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

type Nodes []*Node
func (s Nodes) Len() int      { return len(s) }
func (s Nodes) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
type ByValue struct{ Nodes }
func (s ByValue) Less(i, j int) bool { return s.Nodes[i].V < s.Nodes[j].V }