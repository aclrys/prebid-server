package richaudience

import (
	"encoding/json"
	"testing"

	"github.com/aclrys/prebid-server/openrtb_ext"
)

func TestValidParams(t *testing.T) {
	validator, err := openrtb_ext.NewBidderParamsValidator("../../static/bidder-params")
	if err != nil {
		t.Fatalf("Failed to fetch the json schema. %v", err)
	}

	for _, p := range validParams {
		if err := validator.Validate(openrtb_ext.BidderRichaudience, json.RawMessage(p)); err != nil {
			t.Errorf("Schema rejected valid params: %s", p)
		}
	}
}

func TestInvalidParams(t *testing.T) {
	validator, err := openrtb_ext.NewBidderParamsValidator("../../static/bidder-params")
	if err != nil {
		t.Fatalf("Failed to fetch the json schema. %v", err)
	}

	for _, p := range invalidParams {
		if err := validator.Validate(openrtb_ext.BidderRichaudience, json.RawMessage(p)); err == nil {
			t.Errorf("Schema allowed invalid params: %s", p)
		}
	}
}

var validParams = []string{
	`{"pid":"hash", "supplyType":"site"}`,
	`{"pid":"hash", "supplyType":"site", "test": true}`,
}

var invalidParams = []string{
	`{"pid": 42}`,
	`{"pid": "", "supplyType":0}`,
	`{"pid": 11, "supplyType":"site"}`,
	`{"pid": "hash", "supplyType":11}`,
	`{"pid": "hash"}`,
	`{"supplyType":"site"}`,
	`{}`,
}
