package adocean

import (
	"testing"

	"github.com/aclrys/prebid-server/adapters/adapterstest"
	"github.com/aclrys/prebid-server/config"
	"github.com/aclrys/prebid-server/openrtb_ext"
	"github.com/stretchr/testify/assert"
)

func TestJsonSamples(t *testing.T) {
	bidder, buildErr := Builder(openrtb_ext.BidderAdOcean, config.Adapter{
		Endpoint: "https://{{.Host}}"})

	if buildErr != nil {
		t.Fatalf("Builder returned unexpected error %v", buildErr)
	}

	adapterstest.RunJSONBidderTest(t, "adoceantest", bidder)
}

func TestEndpointTemplateMalformed(t *testing.T) {
	_, buildErr := Builder(openrtb_ext.BidderAdOcean, config.Adapter{
		Endpoint: "{{Malformed}}"})

	assert.Error(t, buildErr)
}
