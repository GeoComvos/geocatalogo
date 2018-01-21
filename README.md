# geocatalogo

[![Build Status](https://travis-ci.org/tomkralidis/geocatalogo.png)](https://travis-ci.org/tomkralidis/geocatalogo)
[![Report Card](https://goreportcard.com/badge/github.com/tomkralidis/geocatalogo)](https://goreportcard.com/report/github.com/tomkralidis/geocatalogo)
[![Godoc](http://img.shields.io/badge/godoc-reference-blue.svg?style=flat)](https://godoc.org/github.com/tomkralidis/geocatalogo)

## Overview

Geospatial Catalogue in Go

## Installation

# Requirements

geocatalogo's default backend is [Elasticsearch](https://www.elastic.co/) and
requires an Elasticsearch endpoint as defined in configuration.

```bash
# install dependencies
go get golang.org/x/text/encoding
go get github.com/sirupsen/logrus
go get gopkg.in/yaml.v2
go build ./...
config/config.go:34:2: cannot find package "gopkg.in/yaml.v2"
# install geocatalogo
go get github.com/tomkralidis/geocatalogo/...
# install utilities/helpers
# default backend
go get github.com/olivere/elastic
# or bleve backend
go get github.com/blevesearch/bleve/...
go install github.com/golang/lint/...
# set configuration
# (sample in $GOPATH/src/github.com/tomkralidis/geocatalogo/geocatalogo-config.env)
cp geocatalogo-config.env local.env
vi local.env  # update accordingly
. local.env
```

## Running

### Using the geocatalogo command line utility

```bash
# list commands
geocatalogo

# index a metadata record
geocatalogo index --file=/path/to/record.xml

# index a directory of metadata records
geocatalogo index --dir=/path/to/dir

# dedicated importers

# Landsat on AWS (https://aws.amazon.com/public-datasets/landsat/)
curl http://landsat-pds.s3.amazonaws.com/c1/L8/scene_list.gz | gunzip > /tmp/scene_list
landsat-aws-importer --file /tmp/scene_list

# OpenAerialMap Catalog (https://docs.openaerialmap.org/catalog/)
curl https://api.openaerialmap.org/meta?limit=5000 > /tmp/oam.json
oam-catalog-importer --file /tmp/scene_list

# search index
geocatalogo search --term=landsat

# search by bbox
geocatalogo search --bbox -152,42,-52,84

# search by time instant
geocatalogo search --time 2018-01-19T18:28:02Z

# search by time range
geocatalogo search --time 2007-11-11T12:43:29Z,2018-01-19T18:28:02Z

# search by any combination exclusively (term, bbox, time)
geocatalogo search --time 2007-11-11T12:43:29Z,2018-01-19T18:28:02Z --bbox -152,42,-52,84 --term landsat

# get a metadata record by id
geocatalogo get --id=12345

# get a metadata record by list of ids
geocatalogo get --id=12345,67890

# run as an HTTP server (default port 8000)
geocatalogo serve
# run as an HTTP server on a custom port
geocatalogo serve --port 8001
# run as an HTTP server honouring the STAC API
geocatalogo serve --api stac

# get version
geocatalogo version
```

### Using the API

```go
// init a Geocatalogue from environment
import (
	"encoding/json"
	"fmt"
	"github.com/tomkralidis/geocatalogo"
	"github.com/tomkralidis/geocatalogo/metadata/parsers"
)

cat, err := geocatalogo.NewFromEnv()
if err != nil {
	fmt.Println(err)
}

// index a Dublin Core metadata record
source, err := ioutil.ReadFile(file)
if err != nil {
	fmt.Printf("Could not read file: %s\n", err)
}
metadataRecord, err := parsers.ParseCSWRecord(source)
if err != nil {
	fmt.Printf("Could not parse metadata: %s\n", err)
	continue
}
result := cat.Index(metadataRecord)
if !result {
	fmt.Println("Error Indexing")
}

// search records and present records 0 - 10
results := cat.Search("birds", 0, 10)

// get record by id
results := cat.Get("record-id-123")

// process results
for _, result := range results.Records {
	b, _ := json.MarshalIndent(result, "", "    ")
	fmt.Printf("%s\n", b)
}
```

## Development

### Running Tests

## Releasing

### Bugs and Issues

All bugs, enhancements and issues are managed on [GitHub](https://github.com/tomkralidis/geocatalogo).

## Contact

* [Tom Kralidis](https://github.com/tomkralidis)
