package gowith

import (
	"errors"
	"testing"
)

// Struct for full working test.
type Db struct{}

// Implement EnterExiter's Enter method.
func (db Db) Enter() (bool, error) {
	return true, nil
}

// Implement EnterExiter's Exit method.
func (db Db) Exit(ent bool, err error) error {
	return err
}

// Struct for error on enter.
type DbErrOnEnt struct{}

// Implement EnterExiter's Enter method.
func (db DbErrOnEnt) Enter() (*string, error) {
	return nil, errors.New("error")
}

// Implement EnterExiter's Exit method.
func (db DbErrOnEnt) Exit(ent *string, err error) error {
	return err
}

// Test full implementation
func TestFull(t *testing.T) {
	var ran bool
	var erv bool

	err := New[bool](Db{}, func(ent bool) error {
		ran = true
		erv = ent
		return nil
	})

	if err != nil {
		t.Errorf("error happened during enter, action, or exit: %v", err)
	}

	if !ran {
		t.Errorf("expected action to run but it did not")
	}

	if erv != true {
		t.Errorf("expected enter value to be %d, but got %v", 3, erv)
	}
}

// Test when an error happens on enter, which means action should not run.
// And Exit should recieve the error.
func TestErrorOnEnter(t *testing.T) {
	var ran bool
	var erv *string

	err := New[*string](DbErrOnEnt{}, func(ent *string) error {
		ran = true
		erv = ent
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
	err := New[bool](Db{}, func(er bool) error {
		return errors.New("error")
	})

	if err.Error() != "error" {
		t.Errorf("expected error to match but did not, got \"%v\" expected \"error\"", err)
	}
}
