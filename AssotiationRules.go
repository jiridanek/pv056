package main

import (
	. "./assoc"
	"encoding/gob"
	"fmt"
	"os"
	"sort"
)

func LoadData(filename string) *FrequentItemsets {
	f, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	dec := gob.NewDecoder(f)
	var fiset FrequentItemsets
	err = dec.Decode(&fiset)
	if err != nil {
		panic(err)
	}
	return &fiset
}

func test_partition() {
	set := make([]string, 0)
	set = append(set, "TYP_AGENDY:titulka", "TYP_AGENDY:diskuse", "TYP_AGENDY:gdf", "TYP_AGENDY:cd")
	for _, r := range partition(set) {
		fmt.Println(*r)
	}
}

func main() {
	
//	test_partition()
	
	min_confidence := 0.5
	var fiset *FrequentItemsets
	fiset = LoadData("frequentitemsets.gob")
	fiset.Sort()

	rules := make([]*AssociationRule,0)
	for _, itemsets := range fiset.Data {
		for _, itemset := range itemsets {
			partitioned := partition(itemset.Items)
			for _, rule := range partitioned {
					if rule.confidence(fiset) >= min_confidence {
					  rules = append(rules, NewAssociationRule(rule, fiset))
					}
			}
			
		}
	}
	
	sort.Sort(Reverse{ByConfidence{rules}})
// 	  for _,rule := range rules {
// 		fmt.Println(*rule)
// 	  }
	print_association_rules(fiset, rules)
	
}

type AssociationRules []*AssociationRule

func (s AssociationRules) Len() int      { return len(s) }
func (s AssociationRules) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

type ByConfidence struct{ AssociationRules }

func (s ByConfidence) Less(i, j int) bool {
	return s.AssociationRules[i].confidence < s.AssociationRules[j].confidence
}

type Reverse struct{ sort.Interface }

func (r Reverse) Less(i, j int) bool { return r.Interface.Less(j, i) }

type AssociationRule struct {
	Rule
	confidence float64
	lift       float64
}

func NewAssociationRule(r *Rule, f *FrequentItemsets) *AssociationRule {
	return &AssociationRule{*r, r.confidence(f), r.lift(f)}
}

type Rule struct {
	Left  []string
	Right []string
}

func (r *Rule) Sort() {
	sort.Strings(r.Left)
	sort.Strings(r.Right)
}

func CopyRule(r *Rule) *Rule {
	left := make([]string, len(r.Left))
	right := make([]string, len(r.Right))
	copy(left, r.Left)
	copy(right, r.Right)
	return &Rule{left, right}
}

func (r *Rule) AppendLeft(s string) *Rule {
	r.Left = append(r.Left, s)
	return r
}

func (r *Rule) AppendRight(s string) *Rule {
	r.Right = append(r.Right, s)
	return r
}

// Pr(Right|Left)
func (r *Rule) confidence(fset *FrequentItemsets) float64 {
	union := make([]string, 0)
	union = append(union, r.Left...)
	union = append(union, r.Right...)
	//sort.Strings(union)
	return float64(support(fset, union)) / float64(support(fset, r.Left))
}

// Pr(Right,Left) / (Pr(Left) * Pr(Right))
func (r *Rule) lift(fset *FrequentItemsets) float64 {
	union := make([]string, 0)
	union = append(union, r.Left...)
	union = append(union, r.Right...)
	//sort.Strings(union)
	
	//Tids are zero based
	tidno := fset.TDB.MaxTid()+1
	//fmt.Println(tidno)
	return (float64(support(fset, union)) * float64(tidno)) / ( float64(support(fset, r.Left))*float64(support(fset, r.Right)) )
}

func partition(set []string) []*Rule {
	rules := make([]*Rule, 0)
	rules = append(rules, &Rule{[]string{}, []string{}})
	for _, v := range set {
		for _, rule := range rules {
			rules = append(rules, CopyRule(rule).AppendLeft(v))
			rule.AppendRight(v)
		}
	}
	// first rule is []|[…] and last […]|[]
	return rules[1:len(rules)-1]
}

func support(fset *FrequentItemsets, set []string) int {
	idset := make(TidList,0,len(set))
	for _,v := range set {
		idset = append(idset, fset.TDB.Iid(v))
	}
	sort.Sort(idset)
	return faster_support(fset, idset)
/*	
	if len(fset.Data)-1 < len(set)-1 {
		panic("set too big")
	}
	sort.Strings(set)
nextset:
	for _, v := range fset.Data[len(set)-1] {
		// both lists will be ordered in the same way
		for i := range set {
			if v.Items[i] != set[i] {
				continue nextset
			}
		}
		return v.Support
	}
	panic("set not found")
*/
}

func faster_support(fiset *FrequentItemsets, set []Id) int {
	node := fiset.Tree.Children[set[0]]
	for _,v := range set[1:] {
		offset := node.V + 1
		// the way this function is used, this will never be out of bounds
		node = node.Children[v-offset]
	}
	return len(node.TidList)
}

func print_association_rules(fiset *FrequentItemsets, rules []*AssociationRule) {
	for _,rule := range rules {
		fmt.Printf("%v %v => %v %v", rule.Left, support(fiset, rule.Left), rule.Right, support(fiset, rule.Right))
		fmt.Printf(" conf: %.2f, lift: %0.2f", rule.confidence, rule.lift)
		fmt.Printf("\n")
	}
	
}