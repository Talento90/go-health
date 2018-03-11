# health

[![Build Status](https://travis-ci.org/Talento90/health.svg?branch=master)](https://travis-ci.org/Talento90/health) [![Go Report Card](https://goreportcard.com/badge/github.com/Talento90/health)](https://goreportcard.com/report/github.com/Talento90/health) [![codecov](https://codecov.io/gh/Talento90/health/branch/master/graph/badge.svg)](https://codecov.io/gh/Talento90/health)
[![GoDoc](https://godoc.org/github.com/Talento90/health?status.svg)](https://godoc.org/github.com/Talento90/health)

Health package simplifies the way you add health check to your service.

For a real application using health please check [ImgArt](https://github.com/Talento90/imgart)

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

## Response Example

```
{  
    "service":"imgart",
    "up_time":"14m5.788341028s",
    "start_time":"2018-03-11T17:02:33Z",
    "memory":{  
        "current":{  
            "total_alloc":8359984,
            "heap_alloc":2285896,
            "rss":5767168
        },
        "initial":{  
            "total_alloc":7784792,
            "heap_alloc":1754064,
            "rss":5701632
        },
        "diff":{  
            "total_alloc":575192,
            "heap_alloc":531832,
            "rss":65536
        }
    },
    "go_routines":21,
    "is_shutting_down":false,
    "health_checkers":{  
        "mongo":{  
            "status":"CHECKED",
            "response_time":"573.813µs"
        },
        "redis":{  
            "status":"CHECKED",
            "error":{  
                "Op":"dial",
                "Net":"tcp",
                "Source":null,
                "Addr":{  
                    "IP":"172.17.0.16",
                    "Port":6379,
                    "Zone":""
                },
                "Err":{  
                    "Syscall":"getsockopt",
                    "Err":113
                }
            },
            "response_time":"93.526014ms"
        },
        "external_api": {
            "status":"TIMEOUT",
            "response_time":"1.2s"
        }
    }
}
```