# gowith

![Tests](https://github.com/gnikyt/Basic-Shopify-API/workflows/Package%20Test/badge.svg?branch=master)
[![Coverage Status](https://coveralls.io/repos/github/gnikyt/gowith/badge.svg?branch=master)](https://coveralls.io/github/gnikyt/gowith?branch=master)

This function simply mocks Python's [with statement](http://docs.python.org/release/2.5.3/ref/with.html) for Go.

## Usage

```go
import (
  "githib.com/gnikyt/gowith"
)
```

The with statement is used to wrap the execution of code with methods defined by an object. This allows common tasks to be encapsulated for convenient reuse.

A with statement is defined as followed: `gowith.New[T]([EnterExiter], [fn])` (API: `New[T](ee EnterExiter, cf func(T) error) error`).

The executation of a with statement is done as followed:

1. The `[EnterExiter]`'s `Enter()` method is invoked
2. The return value from `Enter()` is assigned to the first argument of `[fn]`
3. The `[fn]` is executed (if `Enter()` returned no error)
4. The `[EnterExiter]`'s `Exit()` method is invoked (always). If an error was caused in `[fn]`, the return value from `Enter()` and `[fn]`'s error are to be passed as arguments to `Exit()`.

Here is a sample object and a with-statement:

```go
import (
  "context"
  "database/sql"
)

// Example struct
type Db struct {
  // ...
}

// Implement EnterExiter's Enter method.
// Returns transaction and error.
func (db Db) Enter() (*sql.Tx, error) {
  // Open database
  dbi, err := sql.Open("./example.db")
  if err != nil {
    // Oops, pass error to Exit().
    return nil, err
  }

  // Create a transaction
  ctx := context.Background()
	tx, err := dbi.BeginTx(ctx, nil)
  if err != nil {
    // Oops, pass error to Exit().
    return nil, err
  }

  // All good, pass the database to fn.
  return tx, nil
}

// Implement EnterExiter's Exit method.
// Accepts the Enter's return and Enter's error.
// Returns error.
func (db Db) Exit(tx *sql.Tx, err error) error {
  if err != nil {
    // Enter() or fn had an error, return it.
    tx.Rollback()
  } else {
    // All good.
    err = tx.Commit()
  }
  return err
}

// Example: The type for New is the type that will be returned from Enter
// and passed to the annonymous function.
err := gowith.New[*sql.Tx](new(Db), func(tx *sql.Tx) error {
  _, err := tx.Exec("INSERT INTO xyz (firstname, lastname) VALUES (?, ?)", "John", "Doe")
  return err
})
if err != nil {
  // Oops something went wrong in Enter(), fn, or Exit().
  fmt.Errorf(err.Error())
}
```

The above example is processed as follows:

+ `With` will call `Db.Enter`
+ `Db.Enter` will setup the database, and return a transaction.
+ `With` will now pass the transaction to the anonymous function.
+ `With` now executes fn (if `Enter` returned no error).
+ `With` now calls `Db.Exit` and passes the transaction and the error (if any).
+ `Db.Exit` now checks for an error and rollsback the changes if so, or commits them.

You could even go a step further and build a wrapper to make it easier:

```go
func (db Db) WithTransaction(cf func(tx *sql.Tx) error) error {
  return gowith.New[*sql.Tx](db, cf)
}

db := new(Db)
err := db.WithTransaction(func(tx *sql.Tx) error {
  _, err := tx.Exec("INSERT INTO xyz (firstname, lastname) VALUES (?, ?)", "John", "Doe")
  return err
})
if err != nil {
  // Oops something went wrong in Enter(), fn, or Exit().
  fmt.Errorf(err.Error())
}
```
