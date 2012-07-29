package main

import (
  "fmt"
  "log"
  //"sort"
  "encoding/gob"
  "os"
  . "./assoc"
  //"flag"
  //"runtime/pprof"
)

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

func naostro (min_support int) (*TransactionDB, *Root, []Itemsets)  {
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
   Run(tdb, root, min_support)
   //PrintTree(tdb, root)
   sets := ComputeSets(tdb, root)
   return tdb, root, sets
}

func main () {
    log.Println("Started")

//     naostro()
    tdb, tree, sets := naostro(1000)
    //sets := test()
      
    
//     for i,r := range sets {
//       fmt.Println(i, len(r))
//       fmt.Println(r)
//     }
     
  log.Println("Writing frequentitemsets.gob")
    f,err := os.Create("frequentitemsets.gob")
    if err != nil {
      panic(err)
    }
    defer f.Close()
    enc := gob.NewEncoder(f)
    err = enc.Encode(FrequentItemsets{1000, tdb, tree, sets})
	if err != nil {
		panic(err)
	}
    
    log.Println("Finished")
}