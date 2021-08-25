# lorm
Django ispired ORM.

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

## Examples

### Single object

Django:
```python
try:
    user = User.objects.get(id=42)
except User.DoesNotExist:
    raise Exception("user not found")

user.is_archive = False
user.save()

user.delete()
```

lorm:
```go
user := UserTable.MustGet(42) // same as UserTable.Filter().Id(42).MustGet() 
if user == nil {
    panic("user not found")
}

user.IsArchive = false
if err := user.Save(); err != nil {
    panic(err)
}

if err := user.Delete(); err != nil {
    panic(err)
}

```

### Filter

Django:
```python
for user in User.objects.filter(is_admin=False, created__gt=dt):
    print(user.login)
```

lorm:
```go
for _, user := range UserTable.Filter().IsAdmin(false).CreatedGt(dt).MustList() {
    fmt.Println(user.login)
}
```

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
).Where( // .Where() arguments concatenated with "AND" statement
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
                // HasPrefix is shortcut for "...%"
                GroupTable.Title().ILike(args.Next(op.HasPrefix("aaaaa"))),
            ),
        ),
    ),
).Limit(
    args.Next(10),
)
```

```go
args := op.NewArgs()
query := op.Insert(t, op.Set{
    op.Column("date"): op.Raw("CURRENT_DATE"),
    op.Column("value"): args.Next(1),
}, op.Set{
    op.Column("date"): op.Raw("CURRENT_DATE"),
    op.Column("value"): args.Next(2),
})
```

### Transactions

```go
package main

import (
    "github.com/tada-team/lorm"
    "github.com/tada-team/lorm/op"
)

func ExampleTransaction(i, j int) (int, error) {
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

### Internal naming
* `r` — record (row in database)
* `l` – list of records
* `t` — table
* `f` – filter
* `q` — query
