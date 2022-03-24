package admixer

import (
	"testing"

	"github.com/aclrys/prebid-server/adapters/adapterstest"
	"github.com/aclrys/prebid-server/config"
	"github.com/aclrys/prebid-server/openrtb_ext"
)

func TestJsonSamples(t *testing.T) {
	bidder, buildErr := Builder(openrtb_ext.BidderAdmixer, config.Adapter{
		Endpoint: "http://inv-nets.admixer.net/pbs.aspx"})

	if buildErr != nil {
		t.Fatalf("Builder returned unexpected error %v", buildErr)
	}

	adapterstest.RunJSONBidderTest(t, "admixertest", bidder)
}
