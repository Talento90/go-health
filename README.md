# health

[![Build Status](https://travis-ci.org/Talento90/health.svg?branch=master)](https://travis-ci.org/Talento90/health) [![Go Report Card](https://goreportcard.com/badge/github.com/Talento90/health)](https://goreportcard.com/report/github.com/Talento90/health)

Health package simplifies the way you add health check to your applications.

## Supported Features

- Get aplication status (Memory, Uptime, etc...)
- Health check external dependencies
- HTTP Handler out of the box that gives you the health of your application

## How to use

- Create a new instance of health

- Register external dependencies that your application depends on like external API's, databases and so on...

- Get explicity the status of your application

- If your application is an HTTP server you can use the default handler to expose the application health


## HTTP Handler

