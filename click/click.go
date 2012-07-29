package click

import (
    "time"
)

type Click struct {
	Fake_uco           string // int
	Typ_aplikace       string
	Datum_operace      time.Time
	Nazev_dne_operace  string
	Adresa_prislusnost string
	Adresa_checksum    string
}

func NewClickFromList(list []string) *Click {
	//To define your own format, write down what the standard time would look like formatted your way;
	t, err := time.Parse("200601021504", list[2])
	if err != nil {
		panic(err)
	}
	return &Click{
		Fake_uco:           list[0],
		Typ_aplikace:       list[1],
		Datum_operace:      t, //list[2],
		Nazev_dne_operace:  list[3],
		Adresa_prislusnost: list[4],
		Adresa_checksum:    list[5],
	}
}

type Clicks []*Click

func (c Clicks) Len() int      { return len(c) }
func (c Clicks) Swap(i, j int) { c[i], c[j] = c[j], c[i] }

type ByIpFucoTime struct{ Clicks }

func (s ByIpFucoTime) Less(i, j int) bool {
	if s.Clicks[i].Adresa_checksum == s.Clicks[j].Adresa_checksum {
		if s.Clicks[i].Fake_uco == s.Clicks[j].Fake_uco {
			return s.Clicks[i].Datum_operace.Before(s.Clicks[j].Datum_operace)
		} else {
			return s.Clicks[i].Fake_uco < s.Clicks[j].Fake_uco
		}
	}
	return s.Clicks[i].Adresa_checksum < s.Clicks[j].Adresa_checksum
}
