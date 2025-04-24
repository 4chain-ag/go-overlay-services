package ship

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/4chain-ag/go-overlay-services/pkg/discovery"
	"github.com/bsv-blockchain/go-sdk/overlay"
	"github.com/bsv-blockchain/go-sdk/transaction"
	"github.com/bsv-blockchain/go-sdk/transaction/template/pushdrop"
)

// SHIPTopicManager implements the TopicManager interface for SHIP tokens.
// SHIP (Service Host Interconnect Protocol) tokens facilitate the advertisement
// of nodes hosting specific topics within the overlay network.
type SHIPTopicManager struct{}

// NewSHIPTopicManager creates a new SHIPTopicManager instance
func NewSHIPTopicManager() *SHIPTopicManager {
	return &SHIPTopicManager{}
}

// IdentifyAdmissibleOutputs analyzes a transaction for SHIP outputs.
// It identifies which outputs are valid SHIP tokens based on protocol requirements.
// SHIP tokens must have 5 fields in a PushDrop script with specific formatting.
func (m *SHIPTopicManager) IdentifyAdmissibleOutputs(
	ctx context.Context,
	beefBytes []byte,
	previousCoins map[uint32][]byte,
) (admit overlay.AdmittanceInstructions, err error) {
	_, tx, _, err := transaction.ParseBeef(beefBytes)
	if err != nil {
		return admit, fmt.Errorf("failed to parse transaction: %w", err)
	}

	for vout, output := range tx.Outputs {
		decoded := pushdrop.Decode(output.LockingScript)
		if decoded == nil || len(decoded.Fields) != 4 {
			continue
		}

		shipIdentifier := string(decoded.Fields[0])
		if shipIdentifier != "SHIP" {
			continue
		}

		// Field 3: Check for advertisable URI
		advertisedURI := string(decoded.Fields[2])
		if !discovery.IsAdvertisableURI(advertisedURI) {
			continue
		}

		// Field 4: Check for valid topic name
		topic := string(decoded.Fields[3])
		if !discovery.IsValidTopicOrServiceName(topic) {
			continue
		}
		// SHIP only accepts "tm_" (topic manager) advertisements
		if !strings.HasPrefix(topic, "tm_") {
			continue
		}

		// Verify signature linking
		if isValid, err := discovery.IsTokenSignatureCorrectlyLinked(
			ctx,
			decoded.LockingPublicKey,
			decoded.Fields,
		); err != nil || !isValid {
			continue
		}

		// Output is valid, add to the list of outputs to admit
		admit.OutputsToAdmit = append(admit.OutputsToAdmit, uint32(vout))
	}

	// Log information about admitted outputs
	if len(admit.OutputsToAdmit) > 0 {
		plural := "output"
		if len(admit.OutputsToAdmit) > 1 {
			plural = "outputs"
		}
		log.Printf("🛳️ Ahoy! Admitted %d SHIP %s!", len(admit.OutputsToAdmit), plural)
	}

	// Log information about consumed previous coins
	previousCoinsCount := len(previousCoins)
	if previousCoinsCount > 0 {
		plural := "coin"
		if previousCoinsCount > 1 {
			plural = "coins"
		}
		log.Printf("🚢 Consumed %d previous SHIP %s!", previousCoinsCount, plural)
	}

	// Log a warning if no outputs were admitted and no previous coins consumed
	if len(admit.OutputsToAdmit) == 0 && previousCoinsCount == 0 {
		log.Printf("⚓ No SHIP outputs admitted and no previous SHIP coins consumed.")
	}

	return
}

// GetDocumentation returns documentation about the SHIP topic manager
func (m *SHIPTopicManager) GetDocumentation() string {
	return `# SHIP Topic Manager

**Protocol Name**: SHIP (Service Host Interconnect Protocol)  
**Manager Name**: SHIPTopicManager  

---

## Overview

The SHIP Topic Manager is responsible for identifying _admissible outputs_ in transactions that declare themselves as part of the SHIP protocol. In other words, it looks at transaction outputs (UTXOs) that embed certain metadata via a PushDrop locking script. This metadata must meet SHIP-specific requirements so that your node or application can recognize valid SHIP advertisements.

A **SHIP token** (in the context of BRC-101 overlays) is a UTXO containing information that advertises a node or host providing some topic-based service to the network. That topic must be prefixed with tm_ — short for "topic manager."

---

## Purpose

- **Announce**: The SHIP token is used to signal that "this identity key is hosting a certain topic (prefixed with tm_)."
- **Connect**: By publishing a SHIP output, a node indicates it offers some service or is a participant in a specific overlay "topic."
- **Authorize**: The SHIP token includes a signature which binds it to an identity key, ensuring authenticity and preventing impersonation.

This allows other nodes to discover hosts by querying the lookup service for valid SHIP tokens.

---

## Requirements for a Valid SHIP Output

1. **PushDrop Fields**: Exactly five fields must be present:
   1. "SHIP" — The protocol identifier string.
   2. identityKey — The 33-byte compressed DER secp256k1 public key that claims to own this UTXO.
   3. advertisedURI — A URI string describing how or where to connect (see BRC-101).
   4. topic — A string that identifies the topic. Must:
      - Start with tm_
      - Pass the BRC-87 checks
   5. signature — A valid signature (in DER) proving that identityKey is authorizing this output, in conjunction with the PushDrop locking key.

2. **Signature Verification**:  
   - The signature in the last field must be valid for the data in the first 4 fields.
   - It must match the identity key, which in turn must match the locking public key used in the output script.  
   - See the code in isTokenSignatureCorrectlyLinked for the implementation details.

3. **Advertised URI**:  
   - Must align with what is contemplated in BRC-101, which enforces certain URI formats (e.g., https://, wss://, or custom prefixed https+bsvauth... URIs).
   - No localhost or invalid URIs allowed.

If any of these checks fail, the SHIP token output is _not_ admitted by the topic manager.

---

## Gotchas and Tips

- **Field Ordering**: The fields **must** appear in the exact order specified above (SHIP -> identityKey -> advertisedURI -> topic -> signature).
- **Exact Five Fields**: More or fewer fields will cause the manager to skip the output.
- **Proper Locking Script**: Ensure the output is locked with a valid PushDrop format. If the lockingScript can't be decoded by PushDrop, the output is invalid.
- **Signature Data**: The signature is a raw ECDSA signature over the raw bytes of the preceding fields. The manager expects that the identity key and signature match up with the logic in isTokenSignatureCorrectlyLinked.
- **Funding**: Remember to fund your SHIP output with at least one satoshi so it remains unspent if you want your advertisement to be valid.`
}

// GetMetaData returns metadata about the SHIP topic manager
func (m *SHIPTopicManager) GetMetaData() *overlay.MetaData {
	return &overlay.MetaData{
		Name:        "SHIP Topic Manager",
		Description: "Manages SHIP tokens for service host interconnect.",
	}
}
