package bmclib

import (
	"context"
	"testing"
	"time"

	"github.com/bmc-toolbox/bmclib/v2/logging"
	"gopkg.in/go-playground/assert.v1"
)

func TestBMC(t *testing.T) {
	t.Skip("needs ipmitool and real ipmi server")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	host := "127.0.0.1"
	port := "623"
	user := "ADMIN"
	pass := "ADMIN"

	log := logging.DefaultLogger()
	cl := NewClient(host, port, user, pass, WithLogger(log))
	cl.Registry.Drivers = cl.Registry.FilterForCompatible(ctx)
	var err error
	err = cl.Open(ctx)
	if err != nil {
		t.Fatal(err)
	}
	defer cl.Close(ctx)
	t.Logf("metadata: %+v", cl.GetMetadata())

	cl.Registry.Drivers = cl.Registry.PreferDriver("dummy")
	state, err := cl.GetPowerState(ctx)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(state)
	t.Logf("metadata %+v", cl.GetMetadata())

	cl.Registry.Drivers = cl.Registry.PreferDriver("ipmitool")
	state, err = cl.GetPowerState(ctx)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(state)
	t.Logf("metadata: %+v", cl.GetMetadata())

	users, err := cl.ReadUsers(ctx)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(users)
	t.Logf("metadata: %+v", cl.GetMetadata())

	t.Fatal()
}

func TestWithRedfishVersionsNotCompatible(t *testing.T) {
	host := "127.0.0.1"
	port := "623"
	user := "ADMIN"
	pass := "ADMIN"

	tests := []struct {
		name     string
		versions []string
	}{
		{
			"no versions",
			[]string{},
		},
		{
			"with versions",
			[]string{"1.2.3", "4.5.6"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cl := NewClient(host, port, user, pass, WithRedfishVersionsNotCompatible(tt.versions))
			assert.Equal(t, tt.versions, cl.redfishVersionsNotCompatible)
		})
	}
}

func TestWithConnectionTimeout(t *testing.T) {
	host := "127.0.0.1"
	port := "623"
	user := "ADMIN"
	pass := "ADMIN"

	tests := []struct {
		name    string
		timeout time.Duration
	}{
		{
			"no connection timeout",
			0,
		},
		{
			"with timeout",
			5 * time.Second,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cl := NewClient(host, port, user, pass, WithConnectTimeout(tt.timeout))
			assert.Equal(t, tt.timeout, cl.connectTimeout)
		})
	}
}
