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
	"strings"
//	"strconv"
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

	l := make(map[string] int)

	lr.Read() // skip headerlr.Read()
	for {
		lline, err := lr.Read()
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
		key := make_key(lline[2:9])	
		l[key] += 1
		
	}
	fmt.Println("STUDIUM_NA_FI,AKTIVNI_STUDIUM_NA_FI,USPESNE_STUDIUM_NA_FI,STUDIUM_NA_MU,AKTIVNI_STUDIUM_NA_MU,USPESNE_STUDIUM_NA_MU,UCITEL")
	for k,v := range l {
	  fmt.Println(k, v)
	}
}

func make_key(flags []string) string{
 return strings.Join(flags,"")
}