package gamoshi

import (
	"testing"

	"github.com/aclrys/prebid-server/adapters/adapterstest"
	"github.com/aclrys/prebid-server/config"
	"github.com/aclrys/prebid-server/openrtb_ext"
)

func TestJsonSamplesWithConfiguredURI(t *testing.T) {
	bidder, buildErr := Builder(openrtb_ext.BidderGamoshi, config.Adapter{
		Endpoint: "https://rtb.gamoshi.io"})

	if buildErr != nil {
		t.Fatalf("Builder returned unexpected error %v", buildErr)
	}

	adapterstest.RunJSONBidderTest(t, "gamoshitest", bidder)
}

func TestJsonSamplesWithHardcodedURI(t *testing.T) {
	bidder, buildErr := Builder(openrtb_ext.BidderGamoshi, config.Adapter{})

	if buildErr != nil {
		t.Fatalf("Builder returned unexpected error %v", buildErr)
	}

	adapterstest.RunJSONBidderTest(t, "gamoshitest", bidder)
}
