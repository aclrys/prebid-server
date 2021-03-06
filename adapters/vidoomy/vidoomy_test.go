package vidoomy

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/aclrys/prebid-server/adapters/adapterstest"
	"github.com/aclrys/prebid-server/config"
	"github.com/aclrys/prebid-server/openrtb_ext"
)

func TestVidoomyBidderEndpointConfig(t *testing.T) {
	bidder, buildErr := Builder(openrtb_ext.BidderVidoomy, config.Adapter{
		Endpoint: "http://localhost/bid",
	})

	if buildErr != nil {
		t.Fatalf("Builder returned unexpected error %v", buildErr)
	}

	bidderVidoomy := bidder.(*adapter)

	assert.Equal(t, "http://localhost/bid", bidderVidoomy.endpoint)
}

func TestJsonSamples(t *testing.T) {
	bidder, buildErr := Builder(openrtb_ext.BidderVidoomy, config.Adapter{})

	if buildErr != nil {
		t.Fatalf("Builder returned unexpected error %v", buildErr)
	}

	adapterstest.RunJSONBidderTest(t, "vidoomytest", bidder)
}
