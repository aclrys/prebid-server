package apacdex

import (
	"testing"

	"github.com/aclrys/prebid-server/adapters/adapterstest"
	"github.com/aclrys/prebid-server/config"
	"github.com/aclrys/prebid-server/openrtb_ext"
)

func TestJsonSamples(t *testing.T) {
	bidder, buildErr := Builder(openrtb_ext.BidderApacdex, config.Adapter{
		Endpoint: "//host"})

	if buildErr != nil {
		t.Fatalf("Builder returned unexpected error %v", buildErr)
	}

	adapterstest.RunJSONBidderTest(t, "apacdextest", bidder)
}
