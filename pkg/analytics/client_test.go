// (C) Copyright 2020 Hewlett Packard Enterprise Development LP.

package analytics

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/HewlettPackard/hpecli/internal/platform/db"
	"github.com/sirupsen/logrus"
)

const (
	version            = "1"
	eventHitType       = "eventHitType"
	eventCategory      = "eventCategory"
	eventAction        = "eventAction"
	eventLabel         = "eventLabel"
	eventValue         = "eventValue"
	userAgent          = "hpecli/0.0.1"
	applicationVersion = "0.0.1"
	applicationName    = "hpecli"
	errTempl           = "got=%s, want=%s"
)

const (
	ClientIDKey = "someClientIDKey"
	dbOpenErr   = "Unable to open DB to get client ID"
)

func TestNewAnalyticsClient(t *testing.T) {
	c := newAnalyticsClient(version, eventHitType, eventCategory, eventAction,
		eventValue, eventLabel, userAgent, applicationVersion, applicationName)
	if c == nil {
		t.Fatal("expected AnalyticsClient to not be nil")
	}

	if version != c.Version {
		t.Fatal("version doesn't match")
	}

	if eventHitType != c.EventHitType {
		t.Fatal("eventHitType doesn't match")
	}

	if eventCategory != c.Eventcategory {
		t.Fatal("eventCategory doesn't match")
	}

	if eventAction != c.EventAction {
		t.Fatal("eventAction doesn't match")
	}

	if eventLabel != c.EventLabel {
		t.Fatal("eventLabel doesn't match")
	}

	if eventValue != c.EventValue {
		t.Fatal("eventValue doesn't match")
	}

	if userAgent != c.UserAgent {
		t.Fatal("userAgent doesn't match")
	}

	if applicationVersion != c.ApplicationVersion {
		t.Fatal("applicationVersion doesn't match")
	}

	if applicationName != c.ApplicationName {
		t.Fatal("applicationName doesn't match")
	}
}

func TestTrackEvent(t *testing.T) {
	ts := newTestServer("/collect", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	defer ts.Close()

	c := newAnalyticsClient(version, eventHitType, eventCategory,
		eventAction, eventValue, eventLabel, userAgent, applicationVersion, applicationName)

	err := c.trackEvent()
	if err != nil {
		t.Fatalf("unexpected error in sending GA data")
	}
}

func newTestServer(path string, h func(w http.ResponseWriter, r *http.Request)) *httptest.Server {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	mux.HandleFunc(path, h)

	return server
}

func TestNewClientID(t *testing.T) {
	cases := []struct {
		name string
		want string
	}{
		{
			name: "generate new client id",
			want: "someRandonID",
		},
	}
	for _, c := range cases {
		c := c
		t.Run(c.name, func(t *testing.T) {
			if got := newClientID(); got == "" {
				t.Errorf("NewClientID() = %v, want %v", got, c.want)
			}
		})
	}
}

func TestClientID(t *testing.T) {
	cases := []struct {
		delete bool
		value  bool
		key    string
		name   string
		want   string
	}{
		{
			name:   "delete the key",
			delete: true,
			key:    "GA_CLIENT_ID",
			value:  true,
			want:   "someRandonID",
		},
		{
			name: "check client id not present",
			want: "someRandonID",
		},
	}
	for _, c := range cases {
		c := c
		t.Run(c.name, func(t *testing.T) {
			if c.delete {
				DBCheck(false, c.key, c.value)
			}
			if got, err := clientID(); got == "" || err != nil {
				t.Errorf("clientID() = %v, want %v", got, c.want)
			}
		})
	}
}

func TestEnableGoogleAnalyticsDBError(t *testing.T) {
	cases := []struct {
		want     bool
		name     string
		funcName string
	}{
		{
			name:     "error opening DB for EnableGoogleAnalytics",
			funcName: "enableGA",
			want:     false,
		},
		{
			name:     "error opening DB DisableGoogleAnalytics",
			funcName: "disableGA",
			want:     false,
		},
		{
			name:     "error opening DB analyticsEnabled",
			funcName: "checkGA",
			want:     false,
		},
		{
			name:     "error opening DB CLientID",
			funcName: "clientID",
			want:     false,
		},
	}
	for _, c := range cases {
		c := c
		t.Run(c.name, func(t *testing.T) {
			var got bool
			var err error
			d, err := db.Open()
			if err != nil {
				logrus.Debug(dbOpenErr)
			}
			defer d.Close()
			switch c.funcName {
			case "enableGA":
				err = enableAnalytics()
				if err == nil {
					t.Error("expected error")
				}
			case "disableGA":
				err = disableAnalytics()
				if err == nil {
					t.Error("expected error")
				}
			default:
				got = analyticsEnabled()
				if got != c.want {
					t.Errorf("got %v, want %v", got, c.want)
				}
			}
		})
	}
}

func TestEnableGA(t *testing.T) { //nolint:funlen // long test methods ok
	cases := []struct {
		put   bool
		flag  bool
		value bool
		want  bool
		name  string
		key   string
	}{
		{
			name: "enable GA, if key not exist, then generate",
			flag: false,
			want: true,
		},
		{
			name:  "put the key value as true GA",
			put:   true,
			key:   "GA_DISABLE",
			value: true,
			flag:  true,
			want:  true,
		},
		{
			name: "run enable GA again",
			flag: false,
			want: true,
		},
		{
			name:  "put the key value as false GA",
			put:   true,
			key:   "GA_DISABLE",
			value: false,
			flag:  true,
			want:  true,
		},
		{
			name: "verify enable GA again",
			flag: false,
			want: true,
		},
		{
			name:  "delete the key GA_DISABLE GA",
			put:   false,
			key:   "GA_DISABLE",
			value: true,
			flag:  true,
			want:  true,
		},
		{
			name: "Key not found GA",
			flag: false,
			want: true,
		},
	}
	for _, c := range cases {
		c := c
		t.Run(c.name, func(t *testing.T) {
			if c.flag == true {
				DBCheck(c.put, c.key, c.value)
			}
			err := enableAnalytics()
			if err != nil {
				t.Fatalf("unexpected error in enabling google analytics")
			}
		})
	}
}

func TestClientIDDBError(t *testing.T) {
	cases := []struct {
		wantErr bool
		name    string
	}{
		{
			name:    "DB open error",
			wantErr: true,
		},
	}
	for _, c := range cases {
		c := c
		t.Run(c.name, func(t *testing.T) {
			d, err := db.Open()
			if err != nil {
				logrus.Debug(dbOpenErr)
			}
			defer d.Close()
			got, err := clientID()
			if err == nil {
				t.Errorf("got %v, want error as %v", got, c.wantErr)
			}
		})
	}
}

func TestDisableGA(t *testing.T) { //nolint:funlen // long test methods ok
	cases := []struct {
		put   bool
		flag  bool
		value bool
		want  bool
		name  string
		key   string
	}{
		{
			name: "disable GA, if key not exist, then generate",
			flag: false,
			want: true,
		},
		{
			name:  "put the key value as true for disable GA",
			put:   true,
			key:   "GA_DISABLE",
			value: true,
			flag:  true,
			want:  true,
		},
		{
			name: "run disable GA again",
			flag: false,
			want: true,
		},
		{
			name:  "put the key value as false for disable GA",
			put:   true,
			key:   "GA_DISABLE",
			value: false,
			flag:  true,
			want:  true,
		},
		{
			name: "verify disable GA again",
			flag: false,
			want: true,
		},
		{
			name:  "delete key GA_DISABLE GA",
			put:   false,
			key:   "GA_DISABLE",
			value: true,
			flag:  true,
			want:  true,
		},
		{
			name: "Key not found for disable GA",
			flag: false,
			want: true,
		},
	}
	for _, c := range cases {
		c := c
		t.Run(c.name, func(t *testing.T) {
			if c.flag == true {
				DBCheck(c.put, c.key, c.value)
			}
			err := disableAnalytics()
			if err != nil {
				t.Fatalf("unexpected error in disabling google analytics")
			}
		})
	}
}

func DBCheck(put bool, key string, value bool) {
	d, err := db.Open()
	if err != nil {
		logrus.Debug(dbOpenErr)
	}

	if put {
		if err := d.Put(key, value); err != nil {
			logrus.Debugf("Unable to put the key %s in to DB for TestDisableGA test case", key)
		}
	} else {
		if err := d.Delete(key); err != nil {
			logrus.Debugf("Unable to delete the key %s in DB for TestDisableGA test case", key)
		}
	}

	d.Close()
}

func TestCheckGA(t *testing.T) {
	cases := []struct {
		enable bool
		want   bool
		name   string
	}{
		{
			name:   "check google analytics enabled",
			enable: true,
			want:   true,
		},
		{
			name:   "check google analytics disabled",
			enable: false,
			want:   false,
		},
	}
	for _, c := range cases {
		c := c
		t.Run(c.name, func(t *testing.T) {
			if c.enable {
				if err := enableAnalytics(); err != nil {
					t.Fatalf("unexpected error in enabling google analytics")
				}
			} else {
				if err := disableAnalytics(); err != nil {
					t.Fatalf("unexpected error in disabling google analytics")
				}
			}
			got := analyticsEnabled()
			if got != c.want {
				t.Errorf("TestCheckGA() = %v, want %v", got, c.want)
			}
		})
	}
}
