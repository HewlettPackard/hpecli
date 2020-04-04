package cfm

import (
	"testing"

	"github.com/HewlettPackard/hpecli/internal/platform/context"
)

const testHost = "host1"
const testUser = "someUser"
const testPass = "somePass"

func init() {
	context.DefaultDBOpenFunc = context.MockOpen
}

func TestUserPassConstructor(t *testing.T) {
	cfmClient := newCFMClient(testHost, testUser, testPass)
	if cfmClient == nil {
		t.Fatal("constructor returned nil")
	}

	if cfmClient.Host != testHost {
		t.Fatal("host wasn't populated in constructor")
	}

	if cfmClient.Username != testUser {
		t.Fatal("user wasn't populated in constructor")
	}

	if cfmClient.Password != testPass {
		t.Fatal("password wasn't populated in constructor")
	}
}

func TestAPIKeyConstructor(t *testing.T) {
	const testAPIKey = "some key"

	cfmClient := newCFMClientFromAPIKey(testHost, testAPIKey)
	if cfmClient == nil {
		t.Fatal("constructor returned nil")
	}

	if cfmClient.Token != testAPIKey {
		t.Fatal("user wasn't populated in constructor")
	}
}
