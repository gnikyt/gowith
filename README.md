# gowith

[![Build Status](https://secure.travis-ci.org/ohmybrew/gowith.png?branch=master)](http://travis-ci.org/ohmybrew/gowith)
[![Coverage Status](https://coveralls.io/repos/github/ohmybrew/gowith/badge.svg?branch=master)](https://coveralls.io/github/ohmybrew/gowith?branch=master)

This function simply mocks Python's [with statement](http://docs.python.org/release/2.5.3/ref/with.html) for Go.

## Usage

```go
import (
  gw "githib.com/ohmybrew/gowith"
)
```

The with statement is used to wrap the execution of code with methods defined by an object. This allows common tasks to be encapsulated for convenient reuse.

A with statement is defined as followed: `gw.With([EnterExiter], [fn]);`.

An API of: `With(ee EnterExiter, act func(er *EnterReturn) error) error`.

The executation of a with statement is done as followed:

1. The `[EnterExiter]`'s `Enter()` method is invoked
2. The return value from `Enter()` is assigned to the first argument of `[fn]`
3. The `[fn]` is executed (if `Enter()` returned no error)
4. The `[EnterExiter]`'s `Exit()` method is invoked (always). If an error was caused in `[fn]`, the return value from `Enter()` and `[fn]`'s error are to be passed as arguments to `Exit()`.

Here is a sample object and a with-statement:

```php
type Db struct{
  // ...
}

func (db Db) Enter() (*EnterReturn, error) {
  db, err := db.Open("./example.db")

  if err != nil {
    // Oops, pass error to Exit().
    return nil, err
  }

  // All good, pass the database to fn.
	return &EnterReturn{Value: db}, nil
}

func (db Db) Exit(er *EnterReturn, err error) error {
  db := er.Value

	if err != nil {
    // Enter() or fn had an error, return it.
    db.Rollback()
    return err
  }

  // All good.
  db.Commit()
  return nil
}

// fn
err := gw.With(new(Db), func(er *EnterReturn) error {
  db := er.Value

  stmt, err := db.Prepare("INSERT INTO xyz (firstname, lastname) VALUES (?, ?)")
  stmt.Exec("John", "Doe")

  return err
})

if err != nil {
  // Oops something went wrong in Enter(), fn, or Exit().
  fmt.Errorf(err.Error())
}
```

The above example is processed as follows:

+ `With` will call `Db.Enter`
+ `Db.Enter` will setup the database, and return the database as a value for `EnterReturn`
+ `With` will now pass the `EnterReturn` to the fn for use within the anonymous function
+ `With` now executes fn (if `Enter` returned no error)
+ `With` now calls `Db.Exit` and passes the `EnterReturn` and the error (if any)
+ `Db.Exit` now checks for an error and rollsback the changes if so or commits them.