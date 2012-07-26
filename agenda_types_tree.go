package main

// 
//
// also sanitizes the "aplikace" field, as Weka has trouble with "%" chars in strings

import (
	"encoding/csv"
	"io"
	"os"
	"log"
	"fmt"
	"sort"
	"strings"
	"path"
//	"strconv"
)

func main() {
  
	log.Println("reading zaznamy_export_data")
	df, err := os.Open("zaznamy_export_data")
	if err != nil {
		panic(err)
	}
	defer df.Close()
	dr := csv.NewReader(df)
	//dr.FieldsPerRecord = 9

	d := make(map[string] int)

	dr.Read() // skip header
	for {
		dline, err := dr.Read()
		if err == io.EOF {
			break
		}

		if err != nil {
			panic(err)
		}
// 		fuco, err := strconv.ParseInt(lline[0], 10, 32)
// 		if err != nil {
// 			panic(err)
// 		}
		key := dline[1]
		d[key] += 1
		
		//debug
//  		if key == "/lide/" {
//  		  break
//  }
	}
	
	log.Println("making list")
	list := make(Sipairs, 0, len(d))
        for k, v := range d {
                list = append(list, &Sipair{k,v})
        }
        
        log.Println("processing")
        //print_counts(list)
	//print_grouped(list)
	tree := grouped_tree(list)
	sort_tree(tree)
	print_tree(tree)
	
  
	
	log.Println("finished")
}

func print_counts(list Sipairs) {
  sort.Sort(ByV{list})
	
	fmt.Println("#","\tTYP_APLIKACE")
	for _,i := range list {
	  fmt.Println(i.V, "\t" + i.K)
	  //fmt.Println(i)
	}
}

// stare, neaktualni
func print_grouped(list Sipairs) {
  sort.Sort(ByK{list})
  stack := make([]int,0)
  
  stack = append(stack, 0)
  fmt.Println(list[0].V, tabs(len(stack)) + list[0].K)
  for i := 1; i < len(list); i++ {
      // find prefix
      for len(stack) > 0 &&
	!strings.HasPrefix(list[i].K, list[stack[len(stack)-1]].K) { 
	stack = stack[:len(stack)-1]
      }
      
      stack = append(stack, i)
      fmt.Println(list[i].V, tabs(len(stack)) + list[i].K)
  }
      
}

func sort_tree(node *Node) {
  sort.Sort(Reverse{ByCumV{node.Children}})
  for _,ch := range node.Children {
    sort_tree(ch)
  }
}

type Nodes []*Node

func (s Nodes) Len() int      { return len(s) }
func (s Nodes) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

type ByCumV struct{ Nodes }
func (s ByCumV) Less(i, j int) bool { return s.Nodes[i].CumV < s.Nodes[j].CumV }



func print_tree(tree *Node) {
  tabs := ""
  print_nodes(tree, tabs) 
}

func print_nodes(tree *Node, tabs string) {
  //fmt.Println("tisknu uzel")
  for _, n := range tree.Children {
    fmt.Printf("%6v (%6v) %v%v\n", n.V, n.CumV, tabs, n.K)
    print_nodes(n, tabs+"\t")
  }
}

type Path struct {
  Elems []string
  Is_file bool
  Original_string string
}

// if they are equal, returns also true
func (p *Path) SubDirOf (q *Path) bool {
  l1 := len(p.Elems)
  l2 := len(q.Elems)
  if l1 > l2 {
    return false
  }
  for i,elem := range p.Elems {
    if elem != q.Elems[i] {
      return false
    }
  }
  return true
}

func (p *Path) String() string {
  return strings.Join(p.Elems, "/")
}

func NewPath(p string) *Path {
  p = strings.TrimSpace(p)
  original_path := p
  is_file := strings.HasSuffix(p, "/")
  p = path.Clean(p)
  return &Path{strings.Split(p, "/"), is_file, original_path}
}

type Node struct {
  K *Path
  V int
  CumV int
  Children Nodes
}

func NewNode (k string, v,cumv int) *Node {
  return &Node{NewPath(k), v, cumv, make(Nodes,0)}
}

func (n *Node) AppendChild(ch *Node) {
  n.Children = append(n.Children, ch)
}

func grouped_tree(list Sipairs) *Node {
  sort.Sort(ByK{list})
  stack := make([]*Node,0)
  tree := NewNode("",0,0)
  
  stack = append(stack, tree)
  
  for _,row := range list {
      // find prefix
      for len(stack) > 1 &&
	!stack[len(stack)-1].K.SubDirOf(NewPath(row.K)) { 
	stack[len(stack)-2].CumV += stack[len(stack)-1].CumV
	stack = stack[:len(stack)-1]
      }
      
      n := NewNode(row.K, row.V, row.V)
      //add it to the tree
      stack[len(stack)-1].AppendChild(n)
      stack = append(stack, n)
  }
  return tree
}

func test_print_grouped() {
  list := make(Sipairs,0)
  for _, k := range []string{"a", "aa", "aaa", "aaaa", "b", "bb", "bbb", "aac"} {
    list = append(list, &Sipair{k, 1})
  }
  print_grouped(list)
}

// use str buffer?
func tabs(n int) string {
  t := ""
  for i := 0; i < n; i++ {
    t += "\t"
  }
  return t
}

type Sipairs []*Sipair
func (s Sipairs) Len() int      { return len(s) }
func (s Sipairs) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
type ByV struct{ Sipairs }
func (s ByV) Less(i, j int) bool { return s.Sipairs[i].V < s.Sipairs[j].V }

type ByK struct{ Sipairs }
func (s ByK) Less(i, j int) bool { return s.Sipairs[i].K < s.Sipairs[j].K }

type Reverse struct {
    // This embedded Interface permits Reverse to use the methods of
    // another Interface implementation.
    sort.Interface
}

func (r Reverse) Less(i, j int) bool {
    return r.Interface.Less(j, i)
}

type Sipair struct {
  K string
  V int
}