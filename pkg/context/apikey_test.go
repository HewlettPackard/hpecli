package context

import (
	"errors"
	"testing"

	"github.com/HewlettPackard/hpecli/pkg/db"
)

const (
	apiKeyPrefix  = "somePrefix"
	contextKey    = "someContext"
	host          = "someHost"
	key           = "someKey"
	errTempl      = "got: %s, wanted: %s"
	fail          = "fail"
	errExpected   = "error was expected"
	errUnexpected = "unexpected error: %v"
)

var ErrorDBOpen = errors.New("expected failure to open db")

func failOpenDB() (db.Store, error) {
	return nil, ErrorDBOpen
}

func TestNew(t *testing.T) {
	c := New(contextKey)

	if c.(*APIContext).contextKey != contextKey {
		t.Fatal("ContextKey value not set as expected")
	}

	if c.(*APIContext).dbOpenFunc == nil {
		t.Fatal("DB open func not set")
	}
}
func TestNewWithDB(t *testing.T) {
	c := NewWithDB(contextKey, MockOpen)

	if c.(*APIContext).contextKey != contextKey {
		t.Fatal("ContextKey value not set as expected")
	}

	if c.(*APIContext).dbOpenFunc == nil {
		t.Fatal("DB open func not set")
	}
}

func TestSetModuleContextDBOpenError(t *testing.T) {
	c := NewWithDB(contextKey, failOpenDB)

	err := c.SetModuleContext(host)
	if err != ErrorDBOpen {
		t.Fatal(errExpected)
	}
}

func TestSetModuleContextPutError(t *testing.T) {
	c := NewWithDB(contextKey, MockOpen)

	err := c.SetModuleContext("fail-host")
	if err == nil {
		t.Fatal(errExpected)
	}
}

func TestSetModuleContext(t *testing.T) {
	c := NewWithDB(contextKey, MockOpen)

	if err := c.SetModuleContext(host); err != nil {
		t.Fatalf(errUnexpected, err)
	}
}

func TestModuleContext(t *testing.T) {
	c := NewWithDB(contextKey, MockOpen)

	if err := c.SetModuleContext(host); err != nil {
		t.Fatalf(errUnexpected, err)
	}

	got, err := c.ModuleContext()
	if err != nil {
		t.Fatal(err)
	}

	if got != host {
		t.Fatal("wrong host returned")
	}
}

func TestHostData(t *testing.T) {
	cases := []struct {
		name        string
		dbOpen      DBOpenFunc
		key         string
		value       interface{}
		wantErr     bool
		expectedErr error
	}{
		{
			name:        "Invalid Key",
			dbOpen:      MockOpen,
			key:         "",
			value:       "",
			expectedErr: ErrorInvalidKey,
		},
		{
			name:        "DB Open Failure",
			dbOpen:      failOpenDB,
			key:         key,
			value:       "",
			expectedErr: ErrorDBOpen,
		},
		{
			name:        "Get Failure",
			dbOpen:      MockOpen,
			key:         "fail-host-key",
			value:       "",
			expectedErr: ErrorKeyNotFound,
		},
		{
			name:        "Get Works",
			dbOpen:      MockOpen,
			key:         key,
			value:       "some specific value",
			expectedErr: nil,
		},
	}

	// save some data so we have something to delete
	c := NewWithDB(contextKey, MockOpen)
	c.SetHostData(key, "some specific value")

	for _, test := range cases {
		test := test
		t.Run(test.name, func(t *testing.T) {
			con := NewWithDB(contextKey, test.dbOpen)
			var value string

			got := con.HostData(test.key, &value)
			if !errors.Is(got, test.expectedErr) {
				t.Fatalf("Didn't get the expected result.  got=%s, want=%s", got, test.expectedErr)
			}

			if value != test.value {
				t.Fatalf("Didn't get the expected result.  got=%s, want=%s", value, test.value)
			}
		})
	}
}

func TestSetHostData(t *testing.T) {
	cases := []struct {
		name        string
		dbOpen      DBOpenFunc
		key         string
		value       interface{}
		wantErr     bool
		expectedErr error
	}{
		{
			name:        "Invalid Key",
			dbOpen:      MockOpen,
			key:         "",
			expectedErr: ErrorInvalidKey,
		},
		{
			name:        "Invalid Value",
			dbOpen:      MockOpen,
			key:         key,
			value:       nil,
			expectedErr: ErrorInvalidValue,
		},
		{
			name:        "DB Open Failure",
			dbOpen:      failOpenDB,
			key:         key,
			value:       "value doesn't matter",
			expectedErr: ErrorDBOpen,
		},
		{
			name:        "Set Failure",
			dbOpen:      MockOpen,
			key:         "fail-host-key",
			value:       "value doesn't matter",
			expectedErr: ErrorExpected,
		},
		{
			name:        "Set Works",
			dbOpen:      MockOpen,
			key:         key,
			value:       "value stored",
			expectedErr: nil,
		},
	}

	for _, test := range cases {
		test := test
		t.Run(test.name, func(t *testing.T) {
			con := NewWithDB(contextKey, test.dbOpen)
			got := con.SetHostData(test.key, test.value)
			if !errors.Is(got, test.expectedErr) {
				t.Fatalf("Didn't get the expected result.  got=%s, want=%s", got, test.expectedErr)
			}
		})
	}
}

func TestDeleteHostData(t *testing.T) {
	cases := []struct {
		name        string
		dbOpen      DBOpenFunc
		key         string
		wantErr     bool
		expectedErr error
	}{
		{
			name:        "Invalid Key",
			dbOpen:      MockOpen,
			key:         "",
			expectedErr: ErrorInvalidKey,
		},
		{
			name:        "DB Open Failure",
			dbOpen:      failOpenDB,
			key:         "somekey",
			expectedErr: ErrorDBOpen,
		},
		{
			name:        "Delete Failure",
			dbOpen:      MockOpen,
			key:         "fail-host-key",
			expectedErr: ErrorExpected,
		},
		{
			name:        "Delete Works",
			dbOpen:      MockOpen,
			key:         key,
			expectedErr: nil,
		},
	}

	// save some data so we have something to delete
	c := NewWithDB(contextKey, MockOpen)
	c.SetHostData(key, "data doesn't matter")

	for _, test := range cases {
		test := test
		t.Run(test.name, func(t *testing.T) {
			con := NewWithDB(contextKey, test.dbOpen)
			got := con.DeleteHostData(test.key)
			if !errors.Is(got, test.expectedErr) {
				t.Fatalf("Didn't get the expected result.  got=%s, want=%s", got, test.expectedErr)
			}
		})
	}
}
