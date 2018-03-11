# health

[![Build Status](https://travis-ci.org/Talento90/health.svg?branch=master)](https://travis-ci.org/Talento90/health) [![Go Report Card](https://goreportcard.com/badge/github.com/Talento90/health)](https://goreportcard.com/report/github.com/Talento90/health) [![codecov](https://codecov.io/gh/Talento90/health/branch/master/graph/badge.svg)](https://codecov.io/gh/Talento90/health)
[![GoDoc](https://godoc.org/github.com/Talento90/health?status.svg)](https://godoc.org/github.com/Talento90/health)

Health package simplifies the way you add health check to your service.

## Supported Features

- Service [health status](https://godoc.org/github.com/Talento90/health#Status)
- Graceful Shutdown Pattern
- Health check external dependencies
- HTTP Handler out of the box that returns the health status

## Installation

```
go get -u github.com/Talento90/health
```

## How to use

```
    // Create a new instance of Health
    h := New("service-name", Options{checkersTimeout: time.Second * 1})

    // Register external dependencies	
    h.RegisterChecker("redis", redisDb)
	h.RegisterChecker("mongo", mongoDb)
	h.RegisterChecker("external_api", api)
	
    // Get service health status
    s := h.GetStatus()

    // Listen interrupt OS signals for graceful shutdown
    var gracefulShutdown = make(chan os.Signal)

	signal.Notify(gracefulShutdown, syscall.SIGTERM)
	signal.Notify(gracefulShutdown, syscall.SIGINT)

    go func() {
	    <-gracefulShutdown
		health.Shutdown()

        // Close Databases gracefully        
        // Close HttpServer gracefully
    }


    // if you have an http server you can register the default handler
    // ServeHTTP return 503 (Service Unavailable) if service is shutting down
    http.HandleFunc("/health", h.ServeHTTP)
```
