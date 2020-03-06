package oneview

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
	ovc := newOVClient(testHost, testUser, testPass)
	if ovc == nil {
		t.Fatal("constructor returned nil")
	}

	if ovc.Client.Endpoint != testHost {
		t.Fatal("host wasn't populated in constructor")
	}

	if ovc.Client.User != testUser {
		t.Fatal("user wasn't populated in constructor")
	}

	if ovc.Client.Password != testPass {
		t.Fatal("password wasn't populated in constructor")
	}
}

func TestAPIKeyConstructor(t *testing.T) {
	const testAPIKey = "some key"

	ovc := newOVClientFromAPIKey(testHost, testAPIKey)
	if ovc == nil {
		t.Fatal("constructor returned nil")
	}

	if ovc.Client.APIKey != testAPIKey {
		t.Fatal("user wasn't populated in constructor")
	}
}
