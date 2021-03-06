package schain

import (
	"encoding/json"

	"github.com/mxmCherry/openrtb/v15/openrtb2"
	"github.com/aclrys/prebid-server/openrtb_ext"
)

// NewSChainWriter creates an ORTB 2.5 schain writer instance
func NewSChainWriter(reqExt *openrtb_ext.ExtRequest) (*SChainWriter, error) {
	if !extPrebidSChainExists(reqExt) {
		return &SChainWriter{}, nil
	}

	sChainsByBidder, err := BidderToPrebidSChains(reqExt.Prebid.SChains)
	if err != nil {
		return nil, err
	}

	writer := SChainWriter{
		sChainsByBidder: sChainsByBidder,
	}
	return &writer, nil
}

// SChainWriter is used to write the appropriate schain for a particular bidder defined in the ORTB 2.5 multi-schain
// location (req.ext.prebid.schain) to the ORTB 2.5 location (req.source.ext)
type SChainWriter struct {
	sChainsByBidder map[string]*openrtb_ext.ExtRequestPrebidSChainSChain
}

// Write selects an schain from the multi-schain ORTB 2.5 location (req.ext.prebid.schains) for the specified bidder
// and copies it to the ORTB 2.5 location (req.source.ext). If no schain exists for the bidder in the multi-schain
// location and no wildcard schain exists, the request is not modified.
func (w SChainWriter) Write(req *openrtb2.BidRequest, bidder string) {
	const sChainWildCard = "*"
	var selectedSChain *openrtb_ext.ExtRequestPrebidSChainSChain

	wildCardSChain := w.sChainsByBidder[sChainWildCard]
	bidderSChain := w.sChainsByBidder[bidder]

	// source should not be modified
	if bidderSChain == nil && wildCardSChain == nil {
		return
	}

	if bidderSChain != nil {
		selectedSChain = bidderSChain
	} else {
		selectedSChain = wildCardSChain
	}

	if req.Source == nil {
		req.Source = &openrtb2.Source{}
	} else {
		sourceCopy := *req.Source
		req.Source = &sourceCopy
	}
	schain := openrtb_ext.ExtRequestPrebidSChain{
		SChain: *selectedSChain,
	}
	sourceExt, err := json.Marshal(schain)
	if err == nil {
		req.Source.Ext = sourceExt
	}
}

// extPrebidSChainExists checks if an schain exists in the ORTB 2.5 req.ext.prebid.schain location
func extPrebidSChainExists(reqExt *openrtb_ext.ExtRequest) bool {
	if reqExt == nil {
		return false
	}
	if reqExt.Prebid.SChains == nil {
		return false
	}
	return true
}
