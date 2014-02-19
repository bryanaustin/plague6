
package main

import (
	"flag"
	"encoding/xml"
	"io/ioutil"
	"net/http"
)

type AppConfig struct {
	Quiet bool
	Listen string
	Concurrent int
	Requests uint64
	Time int
	Walks []Walk
}

type Walk struct {
	Steps []IStep
}

type IStep interface {
	Compile() error
	Run() *StepResponse
}

type TrackFile struct {
	HttpSteps []HttpStep	`xml:"httpstep"`
	Walks []TrackWalk			`xml:"walk"`
}

type TrackWalk struct {
	Name string 					`xml:"name,attr"`
	HttpSteps []HttpStep	`xml:"httpstep"`
}

type HttpStep struct {
	Name string 			`xml:"name,attr"`
	Method string 		`xml:"method,attr"`
	Url string				`xml:"url,attr"`
	Headers []Header	`xml:"header"`
	Body string				`xml:"body"`
	Request *http.Request
}

type Header struct {
	Name string				`xml:"name,attr"`
	Value string			`xml:"value,attr"`
}

func ParseArguments() []string {
	ako.AppConfig = new(AppConfig)
	flag.BoolVar(&ako.AppConfig.Quiet, "q", false, "supress messages")
	flag.StringVar(&ako.AppConfig.Listen, "l", "", "listen on this addres for instructions")
	flag.IntVar(&ako.AppConfig.Concurrent, "c", 1, "maximum concurrency")
	flag.Uint64Var(&ako.AppConfig.Requests, "n", 0, "number of requests to make")
	flag.IntVar(&ako.AppConfig.Time, "t", -1, "duration of time to send out requests")
	flag.Parse()
	return flag.Args()
}

func CompileWalkList(filepaths ... string) ([]Walk, error) {
	walks := make([]Walk, 0, 1)
	for _, path := range filepaths {
		newwalks, err := ParseInstructions(path)
		if err != nil { return nil, err }
		walks = append(walks, newwalks...)
	}
	return walks, nil
}

func ParseInstructions(path string) ([]Walk, error) {
	nw := make([]Walk, 0, 1)
	tf := new(TrackFile)
	data, ioerr := ioutil.ReadFile(path)
	if ioerr != nil { return nil, ioerr }
	xmlerr := xml.Unmarshal([]byte(data), tf)
	if xmlerr != nil { return nil, xmlerr }

	if len(tf.Walks) > 0 {
		for _, walk := range tf.Walks {
			moresteps := make([]IStep, len(walk.HttpSteps))
			for wi, wstep := range walk.HttpSteps {
				wstep.Compile()
				moresteps[wi] = &wstep
			}
			nw = append(nw, Walk{ moresteps })
		}
	} else {
		for s := range tf.HttpSteps {
			tf.HttpSteps[s].Compile()
			nw = append(nw, Walk{ []IStep{ &tf.HttpSteps[s] } })
		}
	}

	return nw, nil
}
