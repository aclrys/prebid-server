package sharethrough

import (
	"github.com/aclrys/prebid-server/adapters/adapterstest"
	"github.com/aclrys/prebid-server/config"
	"github.com/aclrys/prebid-server/openrtb_ext"
	"testing"
)

func TestJsonSamples(t *testing.T) {
	adapterVersion = "10.0"

	bidder, buildErr := Builder(openrtb_ext.BidderSharethrough, config.Adapter{
		Endpoint: "http://whatever.url",
	})

	if buildErr != nil {
		t.Fatalf("Builder returned unexpected error %v", buildErr)
	}

	adapterstest.RunJSONBidderTest(t, "sharethroughtest", bidder)
}
