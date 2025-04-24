package ship

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/4chain-ag/go-overlay-services/pkg/discovery"
	"github.com/bsv-blockchain/go-sdk/overlay"
	"github.com/bsv-blockchain/go-sdk/overlay/lookup"
	ec "github.com/bsv-blockchain/go-sdk/primitives/ec"
)

// SHIPLookupService provides the interface for querying SHIP records
type SHIPLookupService struct {
	storage *SHIPStorage
}

// NewSHIPLookupService creates a new SHIPLookupService instance
func NewSHIPLookupService(storage *SHIPStorage) *SHIPLookupService {
	return &SHIPLookupService{
		storage: storage,
	}
}

func (ls *SHIPLookupService) OutputAdded(
	ctx context.Context,
	outpoint *overlay.Outpoint,
	identityKey *ec.PublicKey,
	domain string,
	topic string,
) error {
	return ls.storage.StoreSHIPRecord(ctx, outpoint, identityKey, domain, topic)
}

func (ls *SHIPLookupService) OutputSpent(
	ctx context.Context,
	outpoint *overlay.Outpoint,
	topic string,
) error {
	if topic != "tm_ship" {
		return nil
	}
	return ls.storage.DeleteSHIPRecord(ctx, outpoint)
}

func (ls *SHIPLookupService) OutputDeleted(
	ctx context.Context,
	outpoint *overlay.Outpoint,
	topic string,
) error {
	if topic != "tm_ship" {
		return nil
	}
	return ls.storage.DeleteSHIPRecord(ctx, outpoint)
}

// Lookup implements the lookup functionality for SHIP records
func (ls *SHIPLookupService) Lookup(
	ctx context.Context,
	question *lookup.LookupQuestion,
) (*lookup.LookupAnswer, error) {
	if question.Query == nil {
		return nil, errors.New("a valid query must be provided")
	}

	if question.Service != "ls_ship" {
		return nil, errors.New("lookup service not supported")
	}

	if bytes.Equal(question.Query, []byte("findAll")) {
		outpoints, err := ls.storage.FindAll(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to find all records: %w", err)
		}
		return &lookup.LookupAnswer{
			Type:   lookup.AnswerTypeFreeform,
			Result: outpoints,
		}, nil
	}

	// Parse as SHIPQuery object
	var shipQuery discovery.SHIPQuery
	if err := json.Unmarshal(question.Query, &shipQuery); err != nil {
		return nil, fmt.Errorf("invalid query format: %w", err)
	}

	outpoints, err := ls.storage.FindRecord(ctx, shipQuery)
	if err != nil {
		return nil, fmt.Errorf("failed to find records: %w", err)
	}
	return &lookup.LookupAnswer{
		Type:   lookup.AnswerTypeFreeform,
		Result: outpoints,
	}, nil
}

// GetDocumentation returns documentation for the SHIP lookup service
func (ls *SHIPLookupService) GetDocumentation() string {
	return `# SHIP Lookup Service

**Protocol Name**: SHIP (Service Host Interconnect Protocol)  
**Lookup Service Name**: 'SHIPLookupService'  

## Overview

The SHIP Lookup Service is used to query the known SHIP tokens in your overlay database. 
It allows you to discover nodes that have published SHIP outputs, indicating they host 
or participate in certain topics (prefixed 'tm_').

This lookup service is typically invoked by sending a LookupQuestion with:
- 'service = ls_ship'
- 'query' containing parameters for searching.

## Query Examples

1. **Find all SHIP records**:
   {"service": "ls_ship", "query": "findAll"}

2. **Find by domain**:
   {"service": "ls_ship", "query": {"domain": "https://example.com"}}

3. **Find by topics**:
   {"service": "ls_ship", "query": {"topics": ["tm_bridge", "tm_sync"]}}

4. **Find by domain AND topics**:
   {"service": "ls_ship", "query": {"domain": "https://example.com", "topics": ["tm_bridge"]}}`
}

// GetMetaData returns metadata about the SHIP lookup service
func (ls *SHIPLookupService) GetMetadata() *overlay.MetaData {
	return &overlay.MetaData{
		Name:        "SHIP Lookup Service",
		Description: "Provides lookup capabilities for SHIP tokens.",
	}
}

// Helper function to check if a value is a string
func isString(val interface{}) bool {
	_, ok := val.(string)
	return ok
}
