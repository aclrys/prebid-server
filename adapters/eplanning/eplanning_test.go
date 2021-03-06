package eplanning

import (
	"testing"

	"github.com/aclrys/prebid-server/adapters"
	"github.com/aclrys/prebid-server/adapters/adapterstest"
	"github.com/aclrys/prebid-server/config"
	"github.com/aclrys/prebid-server/openrtb_ext"
)

func TestJsonSamples(t *testing.T) {
	bidder, buildErr := Builder(openrtb_ext.BidderEPlanning, config.Adapter{
		Endpoint: "http://rtb.e-planning.net/pbs/1"})

	if buildErr != nil {
		t.Fatalf("Builder returned unexpected error %v", buildErr)
	}

	setTesting(bidder)
	adapterstest.RunJSONBidderTest(t, "eplanningtest", bidder)
}

func setTesting(bidder adapters.Bidder) {
	bidderEplanning := bidder.(*EPlanningAdapter)
	bidderEplanning.testing = true
}
