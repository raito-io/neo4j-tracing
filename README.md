# RAITO - Neo4J Tracing

![Version](https://img.shields.io/github/v/tag/raito-io/neo4j-tracing?sort=semver&label=version&color=651FFF)
[![Build](https://img.shields.io/github/actions/workflow/status/raito-io/go-dynamo-utils/build.yml?branch=main)](https://github.com/raito-io/go-dynamo-utils/actions/workflows/build.yml)
[![Contribute](https://img.shields.io/badge/Contribute-🙌-green.svg)](/CONTRIBUTING.md)
[![Go version](https://img.shields.io/github/go-mod/go-version/raito-io/neo4j-tracing?color=7fd5ea)](https://golang.org/)
[![Software License](https://img.shields.io/badge/license-Apache%202-brightgreen.svg?label=license)](/LICENSE)
[![Go Reference](https://pkg.go.dev/badge/github.com/raito-io/neo4j-tracing.svg)](https://pkg.go.dev/github.com/raito-io/neo4j-tracing)

## Introduction
`neo4jtracing` is a go library that enables otel distribute tracing for neo4j driver v5. 

## Getting Started
Add this library as a dependency via `go get github.com/raito-io/neo4j-tracing`

## Enable tracing
Tracing can be enabled by using the `neo4j_tracing.Neo4jTracer` object. 
The `Neo4jTracer` a factory that creates `neo4j.DriverWithContext` objects that are wrapped so distributed tracing can be applied.

Start using tracing is very easy. A regular neo4j driver will be created as follows:
```go
package main

import (
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

func main() {
    dbUri := "neo4j://localhost" // scheme://host(:port) (default port is 7687)
    driver, err := neo4j.NewDriverWithContext(dbUri, neo4j.BasicAuth("neo4j", "letmein!", ""))
    if err != nil {
        panic(err)
    }
    // Do something useful
}
```

To enable tracing you need to create your driver by using the `Neo4jTracer` object.
```go
package main

import (
    "github.com/neo4j/neo4j-go-driver/v5/neo4j"
    neo4j_tracing "github.com/raito-io/neo4j-tracing"
)

func main() {
    driverFactory := neo4j_tracing.NewNeo4jTracer()
	
    dbUri := "neo4j://localhost" // scheme://host(:port) (default port is 7687)
    driver, err := driverFactory.NewDriverWithContext(dbUri, neo4j.BasicAuth("neo4j", "letmein!", ""))
    if err != nil {
        panic(err)
    }
    // Do something useful
}
```

### Options
The following options could be used to customize the tracing behavior:
- `WithTracerProvider(provider)`: Specifies a custom tracer provider. By default, the global OpenTelemetry tracer provider is used.

Those options are passed as argument to the `neo4j_tracing.NewNeo4jTracer()` function.