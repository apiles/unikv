# unikv

(sub project of apiles)

Simple KV memory middleware library (for Golang).

## Installation

```bash
go get github.com/apiles/unikv
```

## Usage

First import packages:

```go
import (
    "github.com/apiles/unikv"
    _ "github.com/apiles/unikv/drivers"
)
```

Note that `github.com/apiles/unikv/drivers` contains all the available drivers.

You can specify that with importing `github.com/apiles/unikv/memory` or something.

Available drivers:

| name   | address                          |
| ------ | -------------------------------- |
| memory | `github.com/apiles/unikv/memory` |
| redis  | `github.com/apiles/unikv/redis`  |

Then you need to specify an `unikv.yml`. An example is included in [docs/unikv.yml](docs/unikv.yml).

Note that unikv will look it up in the current work directory. You can also specify it by setting the environment varaible `UNIKV_CONFIGURE`.

Then you can see this example:

```go
func main() {
    ns := unikv.NewNamespace("hello")
    bucket, err := ns.NewBucket("new")
    if err != nil {
        panic(err)
    }
    bucket.PutString("hello", "world")
    bucket.PutInt("new", 1234)
    bucket.Put("c", true)
    var c bool
    bucket.Get("c", &c)
    value, _ := bucket.GetString("hello")
    in, _ := bucket.GetInt("new")
    fmt.Println(value, in, c)
}
```

## API Docs

Go to <https://pkg.go.dev/github.com/apiles/unikv>.
