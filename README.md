[![Build Status](https://travis-ci.org/jbsmith7741/go-tools.svg?branch=master)](https://travis-ci.org/jbsmith7741/go-tools)
[![Go Report Card](https://goreportcard.com/badge/github.com/jbsmith7741/go-tools)](https://goreportcard.com/report/github.com/jbsmith7741/go-tools)

# go-tools
a collection of useful go libraries

## uri
a convenient and easy way to unmarshal a uri to a struct. 

This has been migrated to its own repo see [https://github.com/jbsmith7741/uri]


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
