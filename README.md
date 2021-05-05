# go-docker-initiator

[![Build Status](https://travis-ci.com/Storytel/go-docker-initiator.svg?branch=master)](https://travis-ci.com/Storytel/go-docker-initiator)
[![codecov](https://codecov.io/gh/Storytel/go-docker-initiator/branch/master/graph/badge.svg)](https://codecov.io/gh/Storytel/go-docker-initiator)
[![Go Report Card](https://goreportcard.com/badge/github.com/Storytel/go-docker-initiator)](https://goreportcard.com/report/github.com/Storytel/go-docker-initiator)

Utility for starting docker containers from Go code.
Useful for testing.

## Image Support

This library currently supports the following services out of the box:

- [pubsub](pubsub/pubsub.go)
- [mysql](mysql/mysql.go)
- [firestore](firestore/firestore.go)

## Installation

Install the package with:

```
go get github.com/Storytel/go-docker-initiator
```

Once installed import it to your code:

```
import dockerinitiator github.com/Storytel/go-docker-initiator
```

## Examples

This package is especially useful for testing. With `go-docker-initiator` the configuration for your external integrations lives in the code and not in 3rd party configuration files.

Below is a typical and simple example using `go-docker-initiator` in a test.

```go
package example_test

import (
	"log"
	"testing"

	dockerinitiator "github.com/Storytel/go-docker-initiator"
	mysqlinitiator "github.com/Storytel/go-docker-initiator/mysql"
	"github.com/stretchr/testify/assert"
)

// WithMySQL will clear obsolete containers an spin up a mysql container for use
func WithMySQL() *mysqlinitiator.MysqlInstance {
	if err := dockerinitiator.ClearObsolete(); err != nil {
		log.Panic(err)
	}

	instance, err := mysqlinitiator.Mysql(mysqlinitiator.MysqlConfig{
		Password: "",
		DbName:   "testdb",
	})
	if err != nil {
		log.Panic(err)
	}

	// Set the needed environment variables
	// MYSQL_SERVER, MYSQL_USER, MYSQL_PASSWORD, MYSQL_DATABASE
	if err = instance.Setenv(); err != nil {
		log.Panic(err)
	}

	return instance
}

func TestDatabaseIntegration(t *testing.T) {
	mysqlInstance := WithMySQL()
	defer mysqlInstance.Stop()

	// Establish a database connection to the exposed environment variables
	db, err := InitAndCreateDatabase()
	assert.NoError(t, err)
	defer db.Close()

	// Run any DB seeds here

	// Setup your service and inject the database
	exampleService := ExampleService{
		Db: db,
	}

	// Test your integration
	_, err = exmapleService.Create()
	assert.NoError(t, err)
}
```

## Notes

A single image can be easily shared between tests using techniques such as [TableDrivenTests](https://github.com/golang/go/wiki/TableDrivenTests).

As is often the case, if you have tests set to automatically run on each file-save you might do best to tag the tests running with `go-docker-initiator` to avoid running them on-save.

Tag a file like so:

```
// +build <TAG>
```

Since the tests with the tag wont run automatically you have to manually invoke it with:

```
go test -tags=<TAG> ./...
```

## Storytel Go

https://github.com/Storytel/go-mysql-seed - Simple MySQL seeding package
