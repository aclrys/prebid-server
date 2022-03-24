package videobyte

import (
	"testing"

	"github.com/aclrys/prebid-server/config"

	"github.com/aclrys/prebid-server/adapters/adapterstest"
)

func TestJsonSamples(t *testing.T) {
	bidder, buildErr := Builder("videobyte", config.Adapter{Endpoint: "https://mock.videobyte.com"})
	if buildErr != nil {
		t.Fatalf("Builder returned unexpected error %v", buildErr)
	}
	adapterstest.RunJSONBidderTest(t, "videobytetest", bidder)
}
