# myslqinstance in Go

This is a basic package to handle MySQL instances from go.

## Example of a my.cnf content:

```
[mysqld]
datadir= /Users/martin/gosandbox/test56/data
basedir= /Users/martin/mysqlversions/5.5.33
socket= /Users/martin/gosandbox/test56/mysql.sock
tmpdir= /Users/martin/gosandbox/test56/tmp
pid-file= /Users/martin/gosandbox/test56/mysql.pid

[client]
socket= /Users/martin/gosandbox/test56/mysql.sock
user= root
```


## How to use it


```go
import "github.com/martinarrieta/mysqlinstance"

m := mysqlinstance.New()

// Set the configuration file to use.
m.Configfile = "/tmp/my.cnf"

// Initialize the instance. This option will not create the data or tmp directories,
// it will return an error if the directories does not exists.
m.Initialize()

// Start the instance
m.Start()

// Status of the instance
if m.Status() {
    fmt.Println("Instance is running!")
}

// Stopping the instance
m.Stop()
```
