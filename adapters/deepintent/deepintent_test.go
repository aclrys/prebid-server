package deepintent

import (
	"testing"

	"github.com/aclrys/prebid-server/config"
	"github.com/aclrys/prebid-server/openrtb_ext"

	"github.com/aclrys/prebid-server/adapters/adapterstest"
)

func TestJsonSamples(t *testing.T) {
	bidder, buildErr := Builder(openrtb_ext.BidderDeepintent, config.Adapter{
		Endpoint: "https://prebid.deepintent.com/prebid"})

	if buildErr != nil {
		t.Fatalf("Builder returned unexpected error %v", buildErr)
	}

	adapterstest.RunJSONBidderTest(t, "deepintenttest", bidder)
}
