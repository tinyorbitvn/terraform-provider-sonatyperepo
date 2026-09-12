package xpprovider_test

import (
	"context"
	"testing"

	fwprovider "github.com/hashicorp/terraform-plugin-framework/provider"

	"github.com/tinyorbitvn/terraform-provider-sonatyperepo/xpprovider"
)

// The shim exists so that Crossplane/Upjet-based providers can embed this
// provider in-process; internal/provider is not importable from other modules.
func TestNewExposesFrameworkProvider(t *testing.T) {
	var p fwprovider.Provider = xpprovider.New("test")
	if p == nil {
		t.Fatal("New returned nil provider")
	}
	resp := &fwprovider.SchemaResponse{}
	p.Schema(context.Background(), fwprovider.SchemaRequest{}, resp)
	for _, attr := range []string{"url", "username", "password"} {
		if _, ok := resp.Schema.Attributes[attr]; !ok {
			t.Errorf("provider schema missing attribute %q", attr)
		}
	}
	if len(p.Resources(context.Background())) < 100 {
		t.Errorf("expected >=100 resources, got %d", len(p.Resources(context.Background())))
	}
}
