package gdpr

import (
	"context"
	"errors"
	"fmt"

	"github.com/aclrys/go-gdpr/api"
	"github.com/aclrys/go-gdpr/consentconstants"
	tcf2ConsentConstants "github.com/aclrys/go-gdpr/consentconstants/tcf2"
	"github.com/aclrys/go-gdpr/vendorconsent"
	tcf2 "github.com/aclrys/go-gdpr/vendorconsent/tcf2"
	"github.com/aclrys/go-gdpr/vendorlist"
	"github.com/aclrys/prebid-server/config"
	"github.com/aclrys/prebid-server/openrtb_ext"
)

type permissionsImpl struct {
	cfg              config.GDPR
	gdprDefaultValue Signal
	purposeConfigs   map[consentconstants.Purpose]config.TCF2Purpose
	vendorIDs        map[openrtb_ext.BidderName]uint16
	fetchVendorList  map[uint8]func(ctx context.Context, id uint16) (vendorlist.VendorList, error)
}

func (p *permissionsImpl) HostCookiesAllowed(ctx context.Context, gdprSignal Signal, consent string) (bool, error) {
	gdprSignal = SignalNormalize(gdprSignal, p.cfg)

	if gdprSignal == SignalNo {
		return true, nil
	}

	return p.allowSync(ctx, uint16(p.cfg.HostVendorID), consent, false)
}

func (p *permissionsImpl) BidderSyncAllowed(ctx context.Context, bidder openrtb_ext.BidderName, gdprSignal Signal, consent string) (bool, error) {
	gdprSignal = SignalNormalize(gdprSignal, p.cfg)

	if gdprSignal == SignalNo {
		return true, nil
	}

	id, ok := p.vendorIDs[bidder]
	if ok {
		vendorException := p.isVendorException(consentconstants.Purpose(1), bidder)
		return p.allowSync(ctx, id, consent, vendorException)
	}

	return false, nil
}

func (p *permissionsImpl) AuctionActivitiesAllowed(ctx context.Context, bidderCoreName openrtb_ext.BidderName, bidder openrtb_ext.BidderName, PublisherID string, gdprSignal Signal, consent string, weakVendorEnforcement bool, aliasGVLIDs map[string]uint16) (allowBidReq bool, passGeo bool, passID bool, err error) {
	if _, ok := p.cfg.NonStandardPublisherMap[PublisherID]; ok {
		return true, true, true, nil
	}

	gdprSignal = SignalNormalize(gdprSignal, p.cfg)

	if gdprSignal == SignalNo {
		return true, true, true, nil
	}

	if consent == "" && gdprSignal == SignalYes {
		return false, false, false, nil
	}

	if id, ok := p.resolveVendorId(bidderCoreName, bidder, aliasGVLIDs); ok {
		return p.allowActivities(ctx, id, bidderCoreName, consent, weakVendorEnforcement)
	} else if weakVendorEnforcement {
		return p.allowActivities(ctx, 0, bidderCoreName, consent, weakVendorEnforcement)
	}

	return p.defaultVendorPermissions()
}

func (p *permissionsImpl) defaultVendorPermissions() (allowBidRequest bool, passGeo bool, passID bool, err error) {
	return false, false, false, nil
}

func (p *permissionsImpl) resolveVendorId(bidderCoreName openrtb_ext.BidderName, bidder openrtb_ext.BidderName, aliasGVLIDs map[string]uint16) (id uint16, ok bool) {
	if id, ok = aliasGVLIDs[string(bidder)]; ok {
		return id, ok
	}

	id, ok = p.vendorIDs[bidderCoreName]

	return id, ok
}

func (p *permissionsImpl) allowSync(ctx context.Context, vendorID uint16, consent string, vendorException bool) (bool, error) {
	if consent == "" {
		return false, nil
	}

	parsedConsent, vendor, err := p.parseVendor(ctx, vendorID, consent)
	if err != nil {
		return false, err
	}

	if vendor == nil {
		return false, nil
	}

	if p.cfg.TCF2.Purpose1.EnforcePurpose == config.TCF2NoEnforcement {
		return true, nil
	}
	consentMeta, ok := parsedConsent.(tcf2.ConsentMetadata)
	if !ok {
		err := errors.New("Unable to access TCF2 parsed consent")
		return false, err
	}
	return p.checkPurpose(consentMeta, vendor, vendorID, tcf2ConsentConstants.InfoStorageAccess, vendorException, false), nil
}

func (p *permissionsImpl) allowActivities(ctx context.Context, vendorID uint16, bidder openrtb_ext.BidderName, consent string, weakVendorEnforcement bool) (allowBidRequest bool, passGeo bool, passID bool, err error) {
	parsedConsent, vendor, err := p.parseVendor(ctx, vendorID, consent)
	if err != nil {
		return false, false, false, err
	}

	// vendor will be nil if not a valid TCF2 consent string
	if vendor == nil {
		if weakVendorEnforcement && parsedConsent.Version() == 2 {
			vendor = vendorTrue{}
		} else {
			return false, false, false, nil
		}
	}

	if !p.cfg.TCF2.Enabled {
		return true, false, false, nil
	}

	consentMeta, ok := parsedConsent.(tcf2.ConsentMetadata)
	if !ok {
		err = fmt.Errorf("Unable to access TCF2 parsed consent")
		return
	}

	if p.cfg.TCF2.SpecialFeature1.Enforce {
		vendorException := p.isSpecialFeatureVendorException(bidder)
		passGeo = vendorException || (consentMeta.SpecialFeatureOptIn(1) && (vendor.SpecialFeature(1) || weakVendorEnforcement))
	} else {
		passGeo = true
	}
	if p.cfg.TCF2.Purpose2.EnforcePurpose == config.TCF2FullEnforcement {
		vendorException := p.isVendorException(consentconstants.Purpose(2), bidder)
		allowBidRequest = p.checkPurpose(consentMeta, vendor, vendorID, consentconstants.Purpose(2), vendorException, weakVendorEnforcement)
	} else {
		allowBidRequest = true
	}
	for i := 2; i <= 10; i++ {
		vendorException := p.isVendorException(consentconstants.Purpose(i), bidder)
		if p.checkPurpose(consentMeta, vendor, vendorID, consentconstants.Purpose(i), vendorException, weakVendorEnforcement) {
			passID = true
			break
		}
	}

	return
}

func (p *permissionsImpl) isVendorException(purpose consentconstants.Purpose, bidder openrtb_ext.BidderName) (vendorException bool) {
	if _, ok := p.purposeConfigs[purpose].VendorExceptionMap[bidder]; ok {
		vendorException = true
	}
	return
}

func (p *permissionsImpl) isSpecialFeatureVendorException(bidder openrtb_ext.BidderName) (vendorException bool) {
	if _, ok := p.cfg.TCF2.SpecialFeature1.VendorExceptionMap[bidder]; ok {
		vendorException = true
	}
	return
}

const pubRestrictNotAllowed = 0
const pubRestrictRequireConsent = 1
const pubRestrictRequireLegitInterest = 2

func (p *permissionsImpl) checkPurpose(consent tcf2.ConsentMetadata, vendor api.Vendor, vendorID uint16, purpose consentconstants.Purpose, vendorException, weakVendorEnforcement bool) bool {
	if purpose == tcf2ConsentConstants.InfoStorageAccess && p.cfg.TCF2.PurposeOneTreatment.Enabled && consent.PurposeOneTreatment() {
		return p.cfg.TCF2.PurposeOneTreatment.AccessAllowed
	}
	if consent.CheckPubRestriction(uint8(purpose), pubRestrictNotAllowed, vendorID) {
		return false
	}

	if vendorException {
		return true
	}

	purposeAllowed := p.consentEstablished(consent, vendor, vendorID, purpose, weakVendorEnforcement)
	legitInterest := p.legitInterestEstablished(consent, vendor, vendorID, purpose, weakVendorEnforcement)

	if consent.CheckPubRestriction(uint8(purpose), pubRestrictRequireConsent, vendorID) {
		return purposeAllowed
	}
	if consent.CheckPubRestriction(uint8(purpose), pubRestrictRequireLegitInterest, vendorID) {
		// Need LITransparency here
		return legitInterest
	}

	return purposeAllowed || legitInterest
}

func (p *permissionsImpl) consentEstablished(consent tcf2.ConsentMetadata, vendor api.Vendor, vendorID uint16, purpose consentconstants.Purpose, weakVendorEnforcement bool) bool {
	if !consent.PurposeAllowed(purpose) {
		return false
	}
	if weakVendorEnforcement {
		return true
	}
	if !p.purposeConfigs[purpose].EnforceVendors {
		return true
	}
	if vendor.Purpose(purpose) && consent.VendorConsent(vendorID) {
		return true
	}
	return false
}

func (p *permissionsImpl) legitInterestEstablished(consent tcf2.ConsentMetadata, vendor api.Vendor, vendorID uint16, purpose consentconstants.Purpose, weakVendorEnforcement bool) bool {
	if !consent.PurposeLITransparency(purpose) {
		return false
	}
	if weakVendorEnforcement {
		return true
	}
	if !p.purposeConfigs[purpose].EnforceVendors {
		return true
	}
	if vendor.LegitimateInterest(purpose) && consent.VendorLegitInterest(vendorID) {
		return true
	}
	return false
}

func (p *permissionsImpl) parseVendor(ctx context.Context, vendorID uint16, consent string) (parsedConsent api.VendorConsents, vendor api.Vendor, err error) {
	parsedConsent, err = vendorconsent.ParseString(consent)
	if err != nil {
		err = &ErrorMalformedConsent{
			Consent: consent,
			Cause:   err,
		}
		return
	}

	version := parsedConsent.Version()
	if version != 2 {
		return
	}

	vendorList, err := p.fetchVendorList[version](ctx, parsedConsent.VendorListVersion())
	if err != nil {
		return
	}

	vendor = vendorList.Vendor(vendorID)
	return
}

// AllowHostCookies represents a GDPR permissions policy with host cookie syncing always allowed
type AllowHostCookies struct {
	*permissionsImpl
}

// HostCookiesAllowed always returns true
func (p *AllowHostCookies) HostCookiesAllowed(ctx context.Context, gdprSignal Signal, consent string) (bool, error) {
	return true, nil
}

// Exporting to allow for easy test setups
type AlwaysAllow struct{}

func (a AlwaysAllow) HostCookiesAllowed(ctx context.Context, gdprSignal Signal, consent string) (bool, error) {
	return true, nil
}

func (a AlwaysAllow) BidderSyncAllowed(ctx context.Context, bidder openrtb_ext.BidderName, gdprSignal Signal, consent string) (bool, error) {
	return true, nil
}

func (a AlwaysAllow) AuctionActivitiesAllowed(ctx context.Context, bidderCoreName openrtb_ext.BidderName, bidder openrtb_ext.BidderName, PublisherID string, gdprSignal Signal, consent string, weakVendorEnforcement bool, aliasGVLIDs map[string]uint16) (allowBidReq bool, passGeo bool, passID bool, err error) {
	return true, true, true, nil
}

// vendorTrue claims everything.
type vendorTrue struct{}

func (v vendorTrue) Purpose(purposeID consentconstants.Purpose) bool {
	return true
}
func (v vendorTrue) PurposeStrict(purposeID consentconstants.Purpose) bool {
	return true
}
func (v vendorTrue) LegitimateInterest(purposeID consentconstants.Purpose) bool {
	return true
}
func (v vendorTrue) LegitimateInterestStrict(purposeID consentconstants.Purpose) bool {
	return true
}
func (v vendorTrue) SpecialFeature(featureID consentconstants.SpecialFeature) (hasSpecialFeature bool) {
	return true
}
func (v vendorTrue) SpecialPurpose(purposeID consentconstants.Purpose) (hasSpecialPurpose bool) {
	return true
}
