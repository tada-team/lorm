# lorm

Use Django models from go code.

1. install django app: https://github.com/tada-team/lorm_exporter
2. generate `models.go` from your django models
3. setup database:

```go
import "github.com/tada-team/lorm"

func init() {
    connStr := "dbname=... user=... password=... host=... port=... sslmode=..."
    conn, err := sql.Open("postgres", connStr)
    if err != nil {
        panic(err)
    }
    lorm.SetConn(conn)
}
```

## Abstractions

### Low-level: raw queries

```go
import (
    "github.com/tada-team/lorm"
    "github.com/tada-team/lorm/op"
)

func ExampleQuery(i, j int) (int, error) {
    args := op.NewArgs()
    query := op.RawQuery("SELECT", args.Next(i), "*", args.Next(j))
    
    var res int    
    err := lorm.TxScan(nil, query, args, &res)

    return res, err
}
```

`args` and `query` in example above equals to:
```go
args := op.Args{i, j}
query := "SELECT $1 * $2"
```

`RawQuery()` works like `fmt.Splintln()`.

Supported 3 base functions for database access:
* `lorm.TxExec(tx *lorm.Tx, q op.Query, args op.Args) (sql.Result, error)`
* `lorm.TxQuery(tx *lorm.Tx, q op.Query, args op.Args, each func(*sql.Rows) error) error`
* `lorm.TxScan(tx *lorm.Tx, q op.Query, args op.Args, dest ...interface{}) error`

#### Transactions

```go
import (
    "github.com/tada-team/lorm"
    "github.com/tada-team/lorm/op"
)


func ExampleTransaction() (int, error) {
    args := op.NewArgs()
    query := op.RawQuery("SELECT", args.Next(i), "*", args.Next(j))
    
    var res int    
    if err := lorm.Atomic(func(tx *lorm.Tx) error {
        return lorm.TxScan(tx, query, args, &res)    	
    }); err != nil {
        return res, err
    }

    return res, err
} 
```
