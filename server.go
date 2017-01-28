package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/jung-kurt/gofpdf"

	"appengine"
	"appengine/datastore"
	"appengine/user"
)

func init() {
	http.HandleFunc("/", index)
	http.HandleFunc("/getlist", getlist)
	http.HandleFunc("/createpdf", createPDF)
}

func index(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "onlypage.html")
}

type Data struct {
	NoCM      string
	NamaPts   string
	Diagnosis string
	Tanggal   string
	IKI       bool
}
type Tabel struct {
	Tanggal string
	DataTab []Data
}
type KunjunganPasien struct {
	Diagnosis, LinkID      string
	GolIKI, ATS, ShiftJaga string
	JamDatang              time.Time
	Dokter                 string
	Hide                   bool
	JamDatangRiil          time.Time
}

type DataPasien struct {
	NamaPasien              string
	NomorCM, JenKel, Alamat string
	TglDaftar, Umur         time.Time
}

func pKey(ctx appengine.Context, par, idpar, chld, idchld string) (*datastore.Key, *datastore.Key) {

	gpKey := datastore.NewKey(ctx, "IGD", "fasttrack", 0, nil)
	parKey := datastore.NewKey(ctx, par, idpar, 0, gpKey)
	chldKey := datastore.NewKey(ctx, chld, idchld, 0, parKey)

	return parKey, chldKey
}
func UbahTanggal(tgl time.Time, shift string) time.Time {

	ubah := tgl
	jam := ubah.Hour()

	if jam < 12 && shift == "3" {
		ubah = tgl.AddDate(0, 0, -1)
	}
	return ubah
}
func getlist(w http.ResponseWriter, r *http.Request) {
	awal := time.Date(2016, time.December, 1, 0, 0, 0, 0, time.UTC)
	ctx := appengine.NewContext(r)
	u := user.Current(ctx)
	email := u.Email
	q := datastore.NewQuery("KunjunganPasien").Filter("Dokter =", email).Filter("Hide =", false).Order("-JamDatang")
	t := q.Run(ctx)
	var daf KunjunganPasien
	var pts DataPasien
	var dat Data
	var list []Data

	for {
		k, err := t.Next(&daf)
		if err == datastore.Done {
			break
		}
		if err != nil {
			fmt.Fprintln(w, err)
		}
		jam := UbahTanggal(daf.JamDatang, daf.ShiftJaga)
		if jam.Before(awal) == true {
			break
		}
		if daf.Hide == true {
			continue
		}
		dat.Tanggal = jam.Format("2-01-2006")
		dat.Diagnosis = daf.Diagnosis
		if daf.GolIKI == "1" {
			dat.IKI = true
		} else {
			dat.IKI = false
		}
		nocm := k.Parent()
		dat.NoCM = nocm.StringID()
		err = datastore.Get(ctx, nocm, &pts)
		dat.NamaPts = pts.NamaPasien
		list = append(list, dat)
	}
	jlist, err := json.Marshal(&list)
	if err != nil {
		fmt.Fprintln(w, err)
	}
	fmt.Fprintln(w, string(jlist))
}
func getall(w http.ResponseWriter, r *http.Request) []Data {
	awal := time.Date(2016, time.December, 1, 0, 0, 0, 0, time.UTC)
	ctx := appengine.NewContext(r)
	u := user.Current(ctx)
	email := u.Email
	q := datastore.NewQuery("KunjunganPasien").Filter("Dokter =", email).Filter("Hide =", false).Order("-JamDatang")
	t := q.Run(ctx)
	var daf KunjunganPasien
	var pts DataPasien
	var dat Data
	var list []Data

	for {
		k, err := t.Next(&daf)
		if err == datastore.Done {
			break
		}
		if err != nil {
			fmt.Fprintln(w, err)
		}
		jam := UbahTanggal(daf.JamDatang, daf.ShiftJaga)
		if jam.Before(awal) == true {
			break
		}
		if daf.Hide == true {
			continue
		}
		dat.Tanggal = jam.Format("2-01-2006")
		dat.Diagnosis = daf.Diagnosis
		if daf.GolIKI == "1" {
			dat.IKI = true
		} else {
			dat.IKI = false
		}
		nocm := k.Parent()
		dat.NoCM = nocm.StringID()
		err = datastore.Get(ctx, nocm, &pts)
		dat.NamaPts = pts.NamaPasien
		list = append(list, dat)
	}
	return list
}
func createPDF(w http.ResponseWriter, r *http.Request) {
	//list := getall(w, r)
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 16)
	pdf.Cell(40, 10, "Hello, world!")
	err := pdf.OutputFileAndClose("hello.pdf")
	if err != nil {
		fmt.Fprintln(w, err)
	}

}
