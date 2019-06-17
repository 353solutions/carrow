// The following directive is necessary to make the package coherent:

// +build ignore

// This program generates contributors.go. It can be invoked by running
// go generate
package main

import (	
	"log"
	"os"
	"text/template"
	"time"
)

func main() {
	const url = "https://github.com/golang/go/raw/master/CONTRIBUTORS"

	// rsp, err := http.Get(url)
	// die(err)
	// defer rsp.Body.Close()

	// sc := bufio.NewScanner(rsp.Body)
	carls := []string{}

	// for sc.Scan() {
	// 	if strings.Contains(sc.Text(), "Carl") {
	// 		carls = append(carls, sc.Text())
	// 	}
	// }

	// die(sc.Err())

	f, err := os.Create("gen/contributors.go")
	die(err)
	defer f.Close()

	packageTemplate.Execute(f, struct {
		Timestamp time.Time
		URL       string
		Carls     []string
	}{
		Timestamp: time.Now(),
		URL:       url,
		Carls:     carls,
	})
}

func die(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

var packageTemplate = template.Must(template.New("").Parse(`// Code generated by go generate; DO NOT EDIT.
// This file was generated by robots at
// {{ .Timestamp }}
// using data from
// {{ .URL }}
package project

hello
`))