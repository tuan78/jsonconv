# jsonconv

[![Go Reference](https://pkg.go.dev/badge/github.com/tuan78/jsonconv.svg)](https://pkg.go.dev/github.com/tuan78/jsonconv)
[![CI](https://github.com/tuan78/jsonconv/actions/workflows/codecov.yml/badge.svg)](https://github.com/tuan78/jsonconv/actions/workflows/codecov.yml)
[![codecov](https://codecov.io/gh/tuan78/jsonconv/branch/main/graph/badge.svg?token=4NLTWR5XNW)](https://codecov.io/gh/tuan78/jsonconv)
[![Go Report Card](https://goreportcard.com/badge/github.com/tuan78/jsonconv)](https://goreportcard.com/report/github.com/tuan78/jsonconv)
![GitHub](https://img.shields.io/github/license/tuan78/jsonconv)

# Description

A Golang library and cmd for flattening JSON and converting JSON to CSV.

With jsonconv, you can:
- Flatten a JSON object which contains deeply nested JSON object and JSON array.
- Convert a JSON object or JSON array to CSV data in a flexible way.
- Use jsonconv cmd (built with [cobra](https://github.com/spf13/cobra)) as convenient tool.

# Installation
First, use `go get` to install the latest version of the library:

```
go get -u github.com/tuan78/jsonconv@latest
```

Next, include jsonconv in your application:
```go
import "github.com/tuan78/jsonconv" 
```

# Usage

## Flatten JSON Object

```go
obj := jsonconv.JsonObject{ // equivalent to map[string]interface{}
    "a": 1,
    "b": 2,
    "c": jsonconv.JsonObject{
        "d": jsonconv.JsonObject{
            "e": 3,
        },
    },
    "f": []int{4, 5, 6},
    "g": jsonconv.JsonObject{
        "h": "A",
        "i": true,
        "j": 1,
        "k": 1.5,
        "l": nil,
    },
}
jsonconv.FlattenJsonObject(obj, nil) // The 'obj' will be modified as result
```

Result:

```go
jsonconv.JsonObject{
    "a": 1,
    "b": 2,
    "c__d__e": 3,
    "f[0]": 4,
    "f[1]": 5,
    "f[2]": 6,
    "g__h": "A",
    "g__i": true,
    "g__j": 1,
    "g__k": 1.5,
    "g__l": nil,
}
```

To customize the flattened data, you can create `FlattenOption` and put it in 2nd parameter. For example:

```go
obj := jsonconv.JsonObject{
    "a": 1,
    "b": 2,
    "c": jsonconv.JsonObject{
        "d": jsonconv.JsonObject{
            "e": 3,
        },
    },
    "f": []int{4, 5, 6},
    "g": jsonconv.JsonObject{
        "h": "A",
        "i": true,
        "j": 1,
        "k": 1.5,
        "l": nil,
    },
}
opt := &jsonconv.FlattenOption{
    Level: 1, // Only flatten the 1st nested JSON object and JSON array
    Gap: "|", // Change the gap style between parent object and nested object
    SkipArray: true, // Skip JSON array flattening
    SkipObject: false, // Skip JSON object flattening. For example, leave it as 'false'
}
jsonconv.FlattenJsonObject(obj, opt) // The 'obj' will be modified as result
```

Result:

```go
jsonconv.JsonObject{
    "a": 1,
    "b": 2,
    "c|d": jsonconv.JsonObject{
      "e": 3,
    },
    "f": []int{4, 5, 6},
    "g|h": "A",
    "g|i": true,
    "g|j": 1,
    "g|k": 1.5,
    "g|l": nil,
}
```

## Convert JSON Object or JSON Array to CSV Data

```go
arr := jsonconv.JsonArray{ // equivalent to []map[string]interface{}
    {
        "id": "b042ab5c-ca73-4460-b739-96410ea9d3a6",
        "user": "Jon Doe",
        "score": -100,
        "is active": false,
    },
    {
    	"id": "ce06f5b1-5721-42c0-91e1-9f72a09c250a",
    	"user": "Tuấn",
    	"score": 1.5,
    	"is active": true,
    	"nested": JsonObject{
    		"a": 1,
    		"b": 2,
    	},
    },
    {
    	"id": "4e01b638-44e5-4079-8043-baabbff21cc8",
    	"user": "高橋",
    	"score": 100000000000000000,
    	"is active": true,
    },
}
result := jsonconv.ToCsv(arr, nil)
```

Result:

```go
jsonconv.CsvData{
    {
        "id", "is active", "nested", "score", "user"
    },
    {
        "b042ab5c-ca73-4460-b739-96410ea9d3a6", "false", "", "-100", "Jon Doe"
    },
    {
        "ce06f5b1-5721-42c0-91e1-9f72a09c250a", "true", "map[a:1 b:2]", "1.5", "Tuấn"
    },
    {
        "4e01b638-44e5-4079-8043-baabbff21cc8", "true", "", "100000000000000000", "高橋"
    },
}
```

From the result above, we can see that the header is sorted and JSON object is not flattened. 

So to customize the csv data, you can create `ToCsvOption` and put it in 2nd parameter. For example:

```go
arr := jsonconv.JsonArray{ // equivalent to []map[string]interface{}
    {
        "id": "b042ab5c-ca73-4460-b739-96410ea9d3a6",
        "user": "Jon Doe",
        "score": -100,
        "is active": false,
    },
    {
    	"id": "ce06f5b1-5721-42c0-91e1-9f72a09c250a",
    	"user": "Tuấn",
    	"score": 1.5,
    	"is active": true,
    	"nested": JsonObject{
    		"a": 1,
    		"b": 2,
    	},
    },
    {
    	"id": "4e01b638-44e5-4079-8043-baabbff21cc8",
    	"user": "高橋",
    	"score": 100000000000000000,
    	"is active": true,
    },
}
opt := &jsonconv.ToCsvOption{
    FlattenOption: jsonconv.DefaultFlattenOption, // Let's try default flatten option this time
    BaseHeaders: []string{"id", "user"}, // To keep 'id' column and 'user' column before the rest
}
result := jsonconv.ToCsv(arr, opt)
```

Result:

```go
jsonconv.CsvData{
    {
        "id", "user", "is active", "nested_a", "nested_b", "score"
    },
    {
        "b042ab5c-ca73-4460-b739-96410ea9d3a6", "Jon Doe", "false", "", "", "-100"
    },
    {
        "ce06f5b1-5721-42c0-91e1-9f72a09c250a", "Tuấn", "true", "1", "2", "1.5"
    },
    {
        "4e01b638-44e5-4079-8043-baabbff21cc8", "高橋", "true", "", "", "100000000000000000"
    },
}
```

# Cmd

To install the latest version of jsonconv cmd, you can use `go install` command:

```
go install github.com/tuan78/jsonconv/cmd/jsonconv@latest
```

Next, to see what jsonconv cmd can do, you can run `help` for any command. For example:

```
jsonconv help
jsonconv csv --help
jsonconv flatten --help
```

## Flatten JSON Object or JSON Array

To flatten JSON from JSON file and output fattened JSON file, you just simply run:

```
jsonconv flatten -i sample.json -o fattened.json
```

Alternatively, without `-i` (or `--in`) and without `-o` (or `--out`) in your command, it reads data from `Stdin` and prints result to `Stdout`. So you can execute command in your convenient way:

```
cat sample.json | jsonconv flatten
```

## Convert JSON Object or JSON Array to CSV Data

To convert JSON from JSON file to CSV file, you can run:

```
jsonconv csv -i sample.json -o converted.csv
```

Or you want to read from `Stdin` and print to `Stdout`. For example: 

```
cat sample.json | jsonconv csv
```

Notes that the `csv` command will flatten the json by default. If you don't want JSON flattening feature, you can skip it by using `--noft` (no flattening) flag. For example:

```
cat sample.json | jsonconv csv --noft
```

# License
jsonconv is released under the MIT license. See [LICENSE](https://github.com/tuan78/jsonconv/blob/main/LICENSE)
