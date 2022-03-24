package onetag

import (
	"testing"

	"github.com/aclrys/prebid-server/adapters/adapterstest"
	"github.com/aclrys/prebid-server/config"
	"github.com/aclrys/prebid-server/openrtb_ext"
	"github.com/stretchr/testify/assert"
)

func TestEndpointTemplateMalformed(t *testing.T) {
	_, buildErr := Builder(openrtb_ext.BidderOneTag, config.Adapter{
		Endpoint: "{{Malformed}}"})

	assert.Error(t, buildErr)
}

func TestJsonSamples(t *testing.T) {
	bidder, buildErr := Builder(openrtb_ext.BidderOneTag, config.Adapter{
		Endpoint: "https://example.com/prebid-server/{{.PublisherID}}"})

	assert.NoError(t, buildErr, "Builder returned unexpected error %v", buildErr)

	adapterstest.RunJSONBidderTest(t, "onetagtest", bidder)
}
