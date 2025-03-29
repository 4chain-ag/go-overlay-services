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

type AdvertisementData struct {
	Protocol           overlay.Protocol
	TopicOrServiceName string
}

type Advertiser interface {
	CreateAdvertisements(adsData []*AdvertisementData) (overlay.TaggedBEEF, error)
	FindAllAdvertisements(protocol overlay.Protocol) ([]*Advertisement, error)
	RevokeAdvertisements(advertisements []*Advertisement) (overlay.TaggedBEEF, error)
	ParseAdvertisement(outputScript *script.Script) (*Advertisement, error)
}
