package engine_test

import (
	"bytes"
	"context"
	"io"
	"net/http"

	"github.com/bsv-blockchain/go-sdk/overlay"
	"github.com/bsv-blockchain/go-sdk/overlay/lookup"
	"github.com/bsv-blockchain/go-sdk/script"

	"github.com/4chain-ag/go-overlay-services/pkg/core/advertiser"
)

type fakeRoundTripper struct{}

func (f fakeRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	return &http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(bytes.NewReader([]byte(`{"outputs":[]}`))),
		Header:     make(http.Header),
	}, nil
}

var fakeHTTPClient = &http.Client{
	Transport: fakeRoundTripper{},
}

func init() {
    http.DefaultClient = fakeHTTPClient
}

type fakeAdvertiser struct {
	findAllAdvertisements func(protocol overlay.Protocol) ([]*advertiser.Advertisement, error)
	createAdvertisements  func(data []*advertiser.AdvertisementData) (overlay.TaggedBEEF, error)
	revokeAdvertisements  func(data []*advertiser.Advertisement) (overlay.TaggedBEEF, error)
	parseAdvertisement    func(script *script.Script) (*advertiser.Advertisement, error)
}

func (f fakeAdvertiser) FindAllAdvertisements(protocol overlay.Protocol) ([]*advertiser.Advertisement, error) {
	if f.findAllAdvertisements != nil {
		return f.findAllAdvertisements(protocol)
	}
	return nil, nil
}
func (f fakeAdvertiser) CreateAdvertisements(data []*advertiser.AdvertisementData) (overlay.TaggedBEEF, error) {
	if f.createAdvertisements != nil {
		return f.createAdvertisements(data)
	}
	return overlay.TaggedBEEF{}, nil
}
func (f fakeAdvertiser) RevokeAdvertisements(data []*advertiser.Advertisement) (overlay.TaggedBEEF, error) {
	if f.revokeAdvertisements != nil {
		return f.revokeAdvertisements(data)
	}
	return overlay.TaggedBEEF{}, nil
}
func (f fakeAdvertiser) ParseAdvertisement(script *script.Script) (*advertiser.Advertisement, error) {
	if f.parseAdvertisement != nil {
		return f.parseAdvertisement(script)
	}
	return nil, nil
}

type fakeLookupResolver struct {
	lookupAnswer *lookup.LookupAnswer
	err          error
}

func (f fakeLookupResolver) Lookup(ctx context.Context, question *lookup.LookupQuestion) (*lookup.LookupAnswer, error) {
	return f.lookupAnswer, f.err
}
func (f fakeLookupResolver) OutputAdded(ctx context.Context, outpoint *overlay.Outpoint, outputScript *script.Script, topic string, blockHeight uint32, blockIndex uint64) error {
	return nil
}
func (f fakeLookupResolver) OutputSpent(ctx context.Context, outpoint *overlay.Outpoint, topic string) error {
	return nil
}
func (f fakeLookupResolver) OutputDeleted(ctx context.Context, outpoint *overlay.Outpoint, topic string) error {
	return nil
}
func (f fakeLookupResolver) OutputBlockHeightUpdated(ctx context.Context, outpoint *overlay.Outpoint, blockHeight uint32, blockIndex uint64) error {
	return nil
}
func (f fakeLookupResolver) GetDocumentation() string {
	return ""
}
func (f fakeLookupResolver) GetMetaData() *overlay.MetaData {
	return &overlay.MetaData{}
}
