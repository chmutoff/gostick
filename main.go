// Visit: http://127.0.0.1:8080
package main

import (
	"html/template"
	"log"
	"net/http"
	"path"

	"github.com/tarm/serial"
)

var serialPort *serial.Port

//JoystickValues saves the axis values of 2 joysticks
type JoystickValues struct {
	X1 int `json:"x1"`
	Y1 int `json:"y1"`
	X2 int `json:"x2"`
	Y2 int `json:"y2"`
}

// Default Request Handler
func defaultHandler(w http.ResponseWriter, r *http.Request) {
	fp := path.Join("html", "index.html")
	tmpl, err := template.ParseFiles(fp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := tmpl.Execute(w, nil); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// AJAX Request Handler
func ajaxHandler(w http.ResponseWriter, r *http.Request) {

	var keys = r.URL.Query()
	//log.Println(keys)

	x1 := keys["x1"][0]
	y1 := keys["y1"][0]
	x2 := keys["x2"][0]
	y2 := keys["y2"][0]

	var controlString string = "x1:" + x1 + ",y1:" + y1 + ",x2:" + x2 + ",y2:" + y2
	log.Print(controlString)

	n, err := serialPort.Write([]byte(controlString))
	if err != nil {
		log.Fatal(err)
	} else if n != len(controlString) {
		log.Printf("Lenght of control string is different from written bytes: %d != %d: ", len(controlString), n)
	}

}

func main() {
	c := &serial.Config{Name: "COM4", Baud: 115200}
	var err error
	serialPort, err = serial.OpenPort(c)
	if err != nil {
		log.Fatal(err)
	}

	fs := http.FileServer(http.Dir("./html/js"))
	http.Handle("/js/", http.StripPrefix("/js/", fs))
	http.HandleFunc("/", defaultHandler)
	http.HandleFunc("/ajax", ajaxHandler)
	http.ListenAndServe(":8080", nil)
}
