///////////////////////////////////////////////////////////////////////////////
//
// The MIT License (MIT)
// Copyright (c) 2018 Tom Kralidis
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to
// deal in the Software without restriction, including without limitation the
// rights to use, copy, modify, merge, publish, distribute, sublicense, and/or
// sell copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
// EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF
// MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT.
// IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM,
// DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR
// OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE
// USE OR OTHER DEALINGS IN THE SOFTWARE.
//
///////////////////////////////////////////////////////////////////////////////

// Package web - simple HTTP Wrapper
package web

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/go-spatial/geocatalogo"
	"github.com/go-spatial/geocatalogo/metadata"
	"github.com/go-spatial/geocatalogo/search"
	"github.com/gorilla/mux"
)

type properties struct {
	start    *time.Time `json:"start"`
	end      *time.Time `json:"end"`
	provider string     `json:"provider"`
	license  string     `json:"license"`
}

type links struct {
	rel  string `json:"rel"`
	href string `json:"href"`
}

type assets struct {
	name string `json:"name"`
	href string `json:"href"`
}

type STACItem struct {
	Type       string     `json:"type"`
	id         string     `json:"id"`
	bbox       [4]float64 `json:"bbox"`
	geometry   metadata.Geometry
	properties properties `json:"properties"`
	links      []links    `json:"links"`
	assets     []assets   `json:"assets"`
}

// STACAPIDescription provides the API description
func STACAPIDescription(w http.ResponseWriter, r *http.Request, cat *geocatalogo.GeoCatalogue) {
	fmt.Fprintf(w, "hi there")
}

// STACItems provides STAC compliant Items matching filters
func STACItems(w http.ResponseWriter, r *http.Request, cat *geocatalogo.GeoCatalogue) {
	var value []string
	var filter string
	var bbox []float64
	var timeVal []time.Time
	var limit = 10
	var next int
	var identifiers []string
	var results search.Results
	var tmp string

	kvp := make(map[string][]string)

	for k, v := range r.URL.Query() {
		kvp[strings.ToLower(k)] = v
	}

	vars := mux.Vars(r)
	tmp, ok := vars["id"]

	if ok == true {
		identifiers = append(identifiers, tmp)
	}

	value, _ = kvp["bbox"]
	if len(value) > 0 {
		bboxTokens := strings.Split(value[0], ",")
		if len(bboxTokens) != 4 {
			exception := search.Exception{
				Code:        20002,
				Description: "bbox format error (should be minx,miny,maxx,maxy)"}
			EmitResponseNotOK(w, cat.Config.Server.MimeType, cat.Config.Server.PrettyPrint, &exception)
			return
		}
		for _, bt := range bboxTokens {
			bt, _ := strconv.ParseFloat(bt, 64)
			bbox = append(bbox, bt)
		}
	}
	value, _ = kvp["time"]
	if len(value) > 0 {
		for _, t := range strings.Split(value[0], ",") {
			timestep, err := time.Parse(time.RFC3339, t)
			if err != nil {
				exception := search.Exception{
					Code:        20002,
					Description: "time format error (should be ISO 8601/RFC3339)"}
				EmitResponseNotOK(w, cat.Config.Server.MimeType, cat.Config.Server.PrettyPrint, &exception)
				return
			}
			timeVal = append(timeVal, timestep)
		}
	}

	value, _ = kvp["filter"]
	if len(value) > 0 {
		filter = value[0]
	}

	value, _ = kvp["limit"]
	if len(value) > 0 {
		limit, _ = strconv.Atoi(value[0])
	}

	value, _ = kvp["next"]
	if len(value) > 0 {
		next, _ = strconv.Atoi(value[0])
	}

	if len(identifiers) > 0 {
		results = cat.Get(identifiers)
	} else {
		results = cat.Search(filter, bbox, timeVal, next, limit)
	}

	EmitResponseOK(w, cat.Config.Server.MimeType, cat.Config.Server.PrettyPrint, &results)
	return
}

// STACRouter provides STAC API Routing
func STACRouter(cat *geocatalogo.GeoCatalogue) *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/api", 301)
	}).Methods("GET")

	router.HandleFunc("/api", func(w http.ResponseWriter, r *http.Request) {
		source, _ := ioutil.ReadFile(cat.Config.Server.OpenAPIDef)
		w.Header().Set("Content-Type", cat.Config.Server.MimeType)
		fmt.Fprintf(w, "%s", source)
	}).Methods("GET")

	router.HandleFunc("/items", func(w http.ResponseWriter, r *http.Request) {
		STACItems(w, r, cat)
	}).Methods("GET", "POST")

	router.HandleFunc("/items/{id}", func(w http.ResponseWriter, r *http.Request) {
		STACItems(w, r, cat)
	}).Methods("GET")

	return router
}

// EmitResponseOK provides HTTP response for successful requests
func IEmitResponseOK(w http.ResponseWriter, contentType string, prettyPrint bool, results *search.Results) {
	var jsonBytes []byte

	w.Header().Set("Content-Type", contentType)
	w.WriteHeader(200)
	if prettyPrint == true {
		jsonBytes, _ = json.MarshalIndent(results, "", "    ")
	} else {
		jsonBytes, _ = json.Marshal(results)
	}
	fmt.Fprintf(w, "%s", jsonBytes)
	return
}

// EmitResponseNotOK provides HTTP response for unsuccessful requests
func IEmitResponseNotOK(w http.ResponseWriter, contentType string, prettyPrint bool, exception *search.Exception) {
	var jsonBytes []byte

	w.Header().Set("Content-Type", contentType)
	w.WriteHeader(400)

	if prettyPrint == true {
		jsonBytes, _ = json.MarshalIndent(exception, "", "    ")
	} else {
		jsonBytes, _ = json.Marshal(exception)
	}
	fmt.Fprintf(w, "%s", jsonBytes)
	return
}