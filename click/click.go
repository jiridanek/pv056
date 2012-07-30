package click

import (
	"strconv"
    "time"
)

type Click struct {
	Id 				int
	Fake_uco           int
	Typ_aplikace       string
	Datum_operace      time.Time
	Nazev_dne_operace  string
	Adresa_prislusnost string
	Adresa_checksum    string
}

func NewClickFromList(id int, list []string) *Click {
	u, err := strconv.ParseInt(list[0], 10, 32)
	if err != nil {
		panic(err)
	}
	//To define your own format, write down what the standard time would look like formatted your way;
	t, err := time.Parse("200601021504", list[2])
	if err != nil {
		panic(err)
	}
	return &Click{
		Id:					id,
		Fake_uco:           int(u),
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

type ByIpFucoTimeId struct{ Clicks }

func (s ByIpFucoTimeId) Less(i, j int) bool {
	if s.Clicks[i].Adresa_checksum == s.Clicks[j].Adresa_checksum {
		if s.Clicks[i].Fake_uco == s.Clicks[j].Fake_uco {
			if s.Clicks[i].Datum_operace == s.Clicks[j].Datum_operace {
				return s.Clicks[i].Id < s.Clicks[i].Id
			}
			return s.Clicks[i].Datum_operace.Before(s.Clicks[j].Datum_operace)
		}
		return s.Clicks[i].Fake_uco < s.Clicks[j].Fake_uco
	}
	return s.Clicks[i].Adresa_checksum < s.Clicks[j].Adresa_checksum
}
