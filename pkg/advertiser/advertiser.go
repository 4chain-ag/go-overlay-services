package advertiser

import (
	"github.com/bsv-blockchain/go-sdk/overlay"
	"github.com/bsv-blockchain/go-sdk/script"
)

type Advertisement struct {
	Protocol       overlay.Protocol
	IdentityKey    string
	Domain         string
	TopicOrService string
	Beef           []byte
	OutputIndex    uint32
}

type AdvertisementData struct{}

type Advertiser interface {
	CreateAdvertisements(adsData []AdvertisementData) (overlay.TaggedBEEF, error)
	FindAllAdvertisements(protocol overlay.Protocol) ([]AdvertisementData, error)
	RevokeAdvertisements(advertisements []Advertisement) (overlay.TaggedBEEF, error)
	ParseAdvertisement(outputScript *script.Script) (Advertisement, error)
}
