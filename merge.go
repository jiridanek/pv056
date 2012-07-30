package main

// 
//
// also sanitizes the "aplikace" field, as Weka has trouble with "%" chars in strings

import (
	"encoding/csv"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
)

func main() {
	log.Println("reading zaznamy_export_lide")
	lf, err := os.Open("zaznamy_export_lide")
	if err != nil {
		panic(err)
	}
	defer lf.Close()
	lr := csv.NewReader(lf)
	lr.FieldsPerRecord = 9

	l := make(map[int64][]string)

	lr.Read() // skip headerlr.Read()
	for {
		lline, err := lr.Read()
		if err == io.EOF {
			break
		}

		if err != nil {
			panic(err)
		}
		fuco, err := strconv.ParseInt(lline[0], 10, 32)
		if err != nil {
			panic(err)
		}
		if _, found := l[fuco]; found {
			panic("duplicite fake uco in lide")
		}

		l[fuco] = lline
	}

	df, err := os.Open("zaznamy_export_data")
	if err != nil {
		panic(err)
	}
	defer df.Close()
	dr := csv.NewReader(df)
	dr.FieldsPerRecord = 6

	rf, err := os.Create("zaznamy_merged")
	if err != nil {
		panic(err)
	}
	defer rf.Close()
	rw := csv.NewWriter(rf)

	log.Println("writing zaznamy_merged")
	dr.Read() // skip header
	for {
		dline, err := dr.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			panic(err)
		}
		fuco, err := strconv.ParseInt(dline[0], 10, 32)
		if err != nil {
			panic(err)
		}
		lline, found := l[fuco]
		if !found {
			panic("no person info for a fuco: " + string(fuco))
		}
		line := sanitizeLine(append(lline, dline[1:]...))
		//log.Println(line)
		rw.Write(line)

	}

	rw.Flush()
	log.Println("done")
}

func sanitizeLine(line []string) []string {
	aplikace := 9
	ret := line[:]
	ret[aplikace] = strings.Replace(ret[aplikace], "%", "$", -1)
	return ret
}
