package engine

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/4chain-ag/go-overlay-services/pkg/core/gasp/core"
	"github.com/bsv-blockchain/go-sdk/overlay"
	"github.com/bsv-blockchain/go-sdk/util"
)

type OverlayGASPRemote struct {
	EndpointUrl string
	Topic       string
	HttpClient  util.HTTPClient
}

func (r *OverlayGASPRemote) GetInitialResponse(ctx context.Context, request *core.GASPInitialRequest) (*core.GASPInitialResponse, error) {
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(request); err != nil {
		return nil, err
	} else if req, err := http.NewRequest("POST", r.EndpointUrl+"/requestSyncResponse", io.NopCloser(&buf)); err != nil {
		return nil, err
	} else {
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("X-BSV-Topic", r.Topic)
		if resp, err := r.HttpClient.Do(req); err != nil {
			return nil, err
		} else {
			defer resp.Body.Close()
			if resp.StatusCode != http.StatusOK {
				return nil, &util.HTTPError{
					StatusCode: resp.StatusCode,
					Err:        err,
				}
			}
			result := &core.GASPInitialResponse{}
			if err := json.NewDecoder(resp.Body).Decode(result); err != nil {
				return nil, err
			}
			return result, nil
		}
	}
}

func (r *OverlayGASPRemote) RequestNode(ctx context.Context, graphID *overlay.Outpoint, outpoint *overlay.Outpoint, metadata bool) (*core.GASPNode, error) {
	if j, err := json.Marshal(&core.GASPNodeRequest{
		GraphID:     graphID,
		Txid:        &outpoint.Txid,
		OutputIndex: outpoint.OutputIndex,
		Metadata:    metadata,
	}); err != nil {
		return nil, err
	} else if req, err := http.NewRequest("POST", r.EndpointUrl+"/requestForeignGASPNode", bytes.NewReader(j)); err != nil {
		return nil, err
	} else {
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("X-BSV-Topic", r.Topic)
		if resp, err := r.HttpClient.Do(req); err != nil {
			return nil, err
		} else {
			defer resp.Body.Close()
			if resp.StatusCode != http.StatusOK {
				return nil, &util.HTTPError{
					StatusCode: resp.StatusCode,
					Err:        err,
				}
			}
			result := &core.GASPNode{}
			if err := json.NewDecoder(resp.Body).Decode(result); err != nil {
				return nil, err
			}
			return result, nil
		}
	}
}

func (r *OverlayGASPRemote) GetInitialReplay(ctx context.Context, response *core.GASPInitialResponse) (*core.GASPInitialReply, error) {
	return nil, errors.New("not-implemented")
}

func (r *OverlayGASPRemote) SubmitNode(ctx context.Context, node *core.GASPNode) (*core.GASPNodeResponse, error) {
	return nil, errors.New("not-implemented")
}
