package main

// duplicate code with merge.go 

import (
	"os"
	"encoding/csv"
	"encoding/gob"
	"io"
	"strconv"
)

func main () {
	ppl := read_people("zaznamy_export_lide")
	gob_people("pple.gob", ppl)
}

func read_people(filename string) map[int][]string {
	lf, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer lf.Close()
	lr := csv.NewReader(lf)

	l := make(map[int][]string)

	lr.Read() // skip headerlr.Read()
	for {
		lline, err := lr.Read()
		if err == io.EOF {
			break
		}

		if err != nil {
			panic(err)
		}
		converted, err := strconv.ParseInt(lline[0], 10, 32)
		fuco := int(converted)
		if err != nil {
			panic(err)
		}
		if _, found := l[fuco]; found {
			panic("duplicite fake uco in lide")
		}

		l[fuco] = lline
	}
	return l
}

func gob_people (filename string, pple map[int][]string){
	f, err := os.Create(filename)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	enc := gob.NewEncoder(f)
	enc.Encode(pple)
}