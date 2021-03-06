package medianet

import (
	"github.com/stretchr/testify/assert"
	"testing"

	"github.com/aclrys/prebid-server/adapters/adapterstest"
	"github.com/aclrys/prebid-server/config"
	"github.com/aclrys/prebid-server/openrtb_ext"
)

func TestJsonSamples(t *testing.T) {
	bidder, buildErr := Builder(openrtb_ext.BidderMedianet, config.Adapter{
		Endpoint:         "https://example.media.net/rtb/prebid",
		ExtraAdapterInfo: "http://localhost:8080/extrnal_url",
	})

	if buildErr != nil {
		t.Fatalf("Builder returned unexpected error %v", buildErr)
	}

	adapterstest.RunJSONBidderTest(t, "medianettest", bidder)
}

func TestEndpointTemplateMalformed(t *testing.T) {
	_, buildErr := Builder(openrtb_ext.BidderMedianet, config.Adapter{
		Endpoint: "{{Malformed}}"})

	assert.Nil(t, buildErr)
}
