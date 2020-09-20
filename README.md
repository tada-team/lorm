# lorm

## Install

```bash
go get -u github.com/tada-team/lorm
```

## Setup database

```go
package main

import (
    "database/sql"
    "github.com/tada-team/lorm"
)

func init() {
    connStr := "dbname=... user=... password=... host=... port=... sslmode=..."
    conn, err := sql.Open("postgres", connStr)
    if err != nil {
        panic(err)
    }
    lorm.SetConn(conn)
}
```

## Django models

1. install django app: https://github.com/tada-team/lorm_exporter
2. generate `models.go` from your django models

### Generated methods from .Filter() object

| Django  | lorm |
| ------------- | ------------- |
| ```user = User.objects.get(id=42)```               | ```user := UserTable.Filter().Id(42).MustGet()```  |
| ```users = User.objects.filter(is_admin=False, created__gt=dt)```   | ```users := UserTable.Filter().IsAdmin(false).CreatedGt(dt).MustGet()```  |
| ```user.Save()``` | ```err := user.Save()``` |
| ```user.Delete()``` | ```err := user.Delete()``` |

## Low-level queries

```go
package main

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

### API
Supported 3 base functions for database access:
* `lorm.TxExec(tx *lorm.Tx, q op.Query, args op.Args) (sql.Result, error)`
* `lorm.TxQuery(tx *lorm.Tx, q op.Query, args op.Args, each func(*sql.Rows) error) error`
* `lorm.TxScan(tx *lorm.Tx, q op.Query, args op.Args, dest ...interface{}) error`

### Query constructor

`github.com/tada-team/lorm/op` has many helpers. For example:
 * `op.Select()`
 * `op.Delete()`
 * `op.Insert()`
 * `op.Or()`
 * `op.And()`
 * `op.Not()`
 * etc. Try autocomplete.

```go
args := op.NewArgs()
query := op.Select(
    UserTable.Id(),
    UserTable.Created(),
).From(
    UserTable,
).OrderBy(
    UserTable.Created(),
).Where( // where arguments concatenated with "AND" statement
    UserTable.Created().Gt(args.Next(time.Now().Add(-time.Hour))),
    op.Or(
        UserTable.IsArchive(),        
        op.Not(UserTable.IsAdmin()),
        UserTable.GroupId().InSubquery( // wow, subqueries!
            op.Select(
                GroupTable.Id(),
            ).From(
                GroupTable,
            ).Where(
                GroupTable.Title().ILike(f.Args.Next(op.HasPrefix("aaaaa"))),
            ),
        ),
    ),
).Limit(
    args.Next(10),
)
```

### Transactions

```go
package main

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

    return res, nil
} 
```
