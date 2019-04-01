package gowith

import (
	"errors"
	"testing"
)

// Struct for full working test.
type Db struct{}

func (db Db) Enter() (*EnterReturn, error) {
	return &EnterReturn{Value: 3}, nil
}

func (db Db) Exit(er *EnterReturn, err error) error {
	return err
}

// Struct for error on enter.
type DbErrOnEnt struct{}

func (db DbErrOnEnt) Enter() (*EnterReturn, error) {
	return nil, errors.New("error")
}

func (db DbErrOnEnt) Exit(er *EnterReturn, err error) error {
	return err
}

// Test full implementation
func TestFull(t *testing.T) {
	var ran bool
	var erv interface{}

	err := With(new(Db), func(er *EnterReturn) error {
		ran = true
		erv = er.Value

		return nil
	})

	if err != nil {
		t.Errorf("error happened during enter, action, or exit: %v", err)
	}

	if !ran {
		t.Errorf("expected action to run but it did not")
	}

	if erv != 3 {
		t.Errorf("expected enter value to be %d, but got %v", 3, erv)
	}
}

// Test when an error happens on enter, which means action should not run.
// And Exit should recieve the error.
func TestErrorOnEnter(t *testing.T) {
	var ran bool
	var erv interface{}

	err := With(new(DbErrOnEnt), func(er *EnterReturn) error {
		ran = true
		erv = er.Value

		return nil
	})

	if err.Error() != "error" {
		t.Errorf("expected error to match but did not, got \"%v\" expected \"error\"", err)
	}

	if ran || erv != nil {
		t.Errorf("action should not have ran, but it did")
	}
}

// Test when an error happens on action.
// And Exit should recieve the error.
func TestErrorOnAction(t *testing.T) {
	err := With(new(Db), func(er *EnterReturn) error {
		return errors.New("error")
	})

	if err.Error() != "error" {
		t.Errorf("expected error to match but did not, got \"%v\" expected \"error\"", err)
	}
}
