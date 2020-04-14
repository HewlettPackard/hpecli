// (C) Copyright 2019 Hewlett Packard Enterprise Development LP.

package greenlake

import (
	"errors"
	"testing"

	"github.com/HewlettPackard/hpecli/internal/platform/context"
)

const logoutURI = "/rest/login-sessions"
const expectedErrMsg = "expected to see an error here but didn't"

func init() {
	context.DefaultDBOpenFunc = context.MockOpen
}

func TestHostPrefixAddedForLogout(t *testing.T) {
	const testHost = "https://testhost.net"
	// // clear everything from the mock store
	context.MockClear()

	saveContextAndSessionData(&sessionData{testHost, "token", "tenantID"})

	cmd := newLogoutCommand()
	cmd.SetArgs([]string{"--host", testHost})
	_ = cmd.Execute()

	// check the db to make sure it was persisted
	_, err := hostData(testHost)
	if !errors.Is(err, context.ErrorKeyNotFound) {
		t.Fatal("logout should delete the context")
	}
}

func TestLogoutRemovesAPIKeyFromContext(t *testing.T) {
	// set context to the test server host
	_ = saveContextAndSessionData(&sessionData{greenlakeDefaultHost, "token", "tenantID"})

	_ = runLogout("")

	// verify the data is gone
	var token string

	c := context.New(glContextKey)

	err := c.HostData(dataKey(greenlakeDefaultHost), &token)
	if !errors.Is(err, context.ErrorKeyNotFound) {
		t.Fatal("expected ErrorKeyNotFound error")
	}
}

func TestLogoutRemovesAPIKeyFromParameter(t *testing.T) {
	const testHost = "https://yetanothertesthost.net"
	
	// set context to the test server host
	_ = saveContextAndSessionData(&sessionData{testHost, "token", "tenantID"})

	// erase the context .. just to make sure that it doesn't pickup the host from the context
	_ = setContext("")

	// specify the param like it was passed on the command line
	_ = runLogout(testHost)

	// verify the data is gone
	var token string

	c := context.New(glContextKey)

	err := c.HostData(dataKey(testHost), &token)
	if !errors.Is(err, context.ErrorKeyNotFound) {
		t.Fatal("expected ErrorKeyNotFound error")
	}
}

func TestLogoutFailsWhenItCantGetContext(t *testing.T) {
	// erase context for runnig the command
	_ = setContext("")

	// run the command that will fail because of the missing context
	if err := runLogout(""); err == nil {
		t.Fatal(expectedErrMsg)
	}
}

