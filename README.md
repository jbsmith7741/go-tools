[![Build Status](https://travis-ci.org/jbsmith7741/go-tools.svg?branch=master)](https://travis-ci.org/jbsmith7741/go-tools)
[![Go Report Card](https://goreportcard.com/badge/github.com/jbsmith7741/go-tools)](https://goreportcard.com/report/github.com/jbsmith7741/go-tools)

# go-tools
a collection of useful go libraries

## uri
a convenient and easy way to unmarshal a uri to a struct.
 
### keywords
- schema
- host
- path
- authority (schema:host)
- origin (schema:host/path)


### example
If we have the uri "http://example.com/path/to/page?name=ferret&color=purple" we can unmarshal this to a predefined struct as follows
``` go 
type Example struct {
    Schema `uri:"schema"`
    Host   `uri:"Host"`
    Path   `uri:"path"`
    Name   `uri:"name"`
    Color  `uri:"color"`
}

func() {
e := Example{}

err := uri.Unmarshal("http://example.com/path/to/page?name=ferret&color=purple", &e)
 
}
```
this would become the following struct 
``` go
e := Example{
    Schema: "www",
    Host:   "example.com",
    Path:   "path/to/page",
    Name:   "ferret",
    Color:  "purple",
    }
 
```

### example 

``` golang 
uri = http://example.org/wiki/Main_Page?Option1=10&Text=hello 

type MyStruct struct {
    Schema `uri:"scheme"`
    Host `uri:"host"`
    Path `uri:"path"`
    Option1 int
    Text string 
}

func Parse() {
    var s *MyStruct
    uri.Unmarshal(s, uri)
}
```

## appenderr
A lot of times we have functions that have multiple error checks in them. Sometimes its helpful to be able to lot at full set of errors rather than the first occurring error. AppendErr is an easy way to add multiple errors together and return the whole set in a single error interface. appenderr is thread safe and can be used to collect errors that occur in different go routines. Each error is counted (based on the string value) and displayed on it's own line.

``` go
// Checkline verifies that the line is a list of comma seperated integers.
// It returns a list of invalid fields if any exist
func CheckLine(line string)  error {
    errs := appenderr.New() 
    for _, v := range strings.Split(line, ",") {
        _, err := strconv.Atoi(v)
        errs.Add(err)
    }
    return errs.ErrOrNil()
}

```

## sqlh
a sql helper class that simplifies configs and connections to databases. sqlh removes the need for anonymous imports and allows basic testing of database calls without having to worry about connecting to a real database.

sqlh should be used to test around database calls and not to test the calls themselves. If you need to test the database logic use [sqlmock](https://github.com/DATA-DOG/go-sqlmock)

## trial
trial is a helper to write Table Driven Tests.