package salunamedia

import (
	"testing"

	"github.com/aclrys/prebid-server/adapters/adapterstest"
	"github.com/aclrys/prebid-server/config"
	"github.com/aclrys/prebid-server/openrtb_ext"
	"github.com/stretchr/testify/assert"
)

func TestJsonSamples(t *testing.T) {
	bidder, buildErr := Builder(openrtb_ext.BidderSaLunaMedia, config.Adapter{
		Endpoint: "http://test.com/pserver"})

	assert.NoError(t, buildErr)
	adapterstest.RunJSONBidderTest(t, "salunamediatest", bidder)
}
