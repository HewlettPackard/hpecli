package oneview

import (
	"testing"
)

const testHost = "host1"
const testUser = "someUser"
const testPass = "somePass"

func TestUserPassConstructor(t *testing.T) {
	ovc := NewOVClient(testHost, testUser, testPass)
	if ovc == nil {
		t.Fatal("constructor returned nil")
	}

	if ovc.Client.Endpoint != testHost {
		t.Fatal("user wasn't populated in constructor")
	}

	if ovc.Client.User != testUser {
		t.Fatal("user wasn't populated in constructor")
	}

	if ovc.Client.Password != testPass {
		t.Fatal("user wasn't populated in constructor")
	}
}

func TestAPIKeyConstructor(t *testing.T) {
	const testAPIKey = "some key"

	ovc := NewOVClientFromAPIKey(testHost, testAPIKey)
	if ovc == nil {
		t.Fatal("constructor returned nil")
	}

	if ovc.Client.APIKey != testAPIKey {
		t.Fatal("user wasn't populated in constructor")
	}
}
