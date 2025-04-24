package discovery

import (
	"bytes"
	"context"
	"net/url"
	"regexp"
	"strconv"
	"strings"

	ec "github.com/bsv-blockchain/go-sdk/primitives/ec"
	"github.com/bsv-blockchain/go-sdk/wallet"
)

// IsAdvertisableURI checks if the provided URI is advertisable, with a recognized URI prefix.
// Applies scheme-specific validation rules as defined by the BRC-101 overlay advertisement spec.
//
// - For HTTPS-based schemes (https://, https+bsvauth+smf://, https+bsvauth+scrypt-offchain://, https+rtt://)
//   - Uses the URL parser (after substituting the custom scheme with "https:" where needed)
//   - Disallows "localhost" as hostname
//
// - For wss:// URIs (for real-time lookup streaming)
//   - Ensures valid URL with protocol "wss:" and non-"localhost" hostname
//
// - For JS8 Call–based URIs (js8c+bsvauth+smf:)
//   - Requires a query string with parameters: lat, long, freq, and radius.
//   - Validates that lat is between -90 and 90 and long between -180 and 180.
//   - Validates that freq and radius each include a positive number.
func IsAdvertisableURI(uri string) bool {
	if uri == "" {
		return false
	}

	// Helper function: validate a URL by substituting its scheme if needed.
	validateCustomHttpsURI := func(uri string, prefix string) bool {
		modifiedURI := strings.Replace(uri, prefix, "https://", 1)
		parsed, err := url.Parse(modifiedURI)
		if err != nil {
			return false
		}
		if strings.ToLower(parsed.Hostname()) == "localhost" {
			return false
		}
		if parsed.Path != "" && parsed.Path != "/" {
			return false
		}
		return true
	}

	// HTTPS-based schemes – disallow localhost.
	if strings.HasPrefix(uri, "https://") {
		return validateCustomHttpsURI(uri, "https://")
	} else if strings.HasPrefix(uri, "https+bsvauth://") {
		// Plain auth over HTTPS, but no payment can be collected
		return validateCustomHttpsURI(uri, "https+bsvauth://")
	} else if strings.HasPrefix(uri, "https+bsvauth+smf://") {
		// Auth and payment over HTTPS
		return validateCustomHttpsURI(uri, "https+bsvauth+smf://")
	} else if strings.HasPrefix(uri, "https+bsvauth+scrypt-offchain://") {
		// A protocol allowing you to also supply sCrypt off-chain values to the topical admissibility checking context
		return validateCustomHttpsURI(uri, "https+bsvauth+scrypt-offchain://")
	} else if strings.HasPrefix(uri, "https+rtt://") {
		// A protocol allowing overlays that deal with real-time transactions (non-finals)
		return validateCustomHttpsURI(uri, "https+rtt://")
	} else if strings.HasPrefix(uri, "wss://") {
		// WSS for real-time event-listening lookups.
		parsed, err := url.Parse(uri)
		if err != nil {
			return false
		}
		if parsed.Scheme != "wss" {
			return false
		}
		if strings.ToLower(parsed.Hostname()) == "localhost" {
			return false
		}
		return true
	} else if strings.HasPrefix(uri, "js8c+bsvauth+smf:") {
		// JS8 Call–based advertisement.
		// Expect a query string with parameters.
		queryIndex := strings.Index(uri, "?")
		if queryIndex == -1 {
			return false
		}

		queryStr := uri[queryIndex:]
		values, err := url.ParseQuery(queryStr)
		if err != nil {
			return false
		}

		// Required parameters: lat, long, freq, and radius.
		latStr := values.Get("lat")
		longStr := values.Get("long")
		freqStr := values.Get("freq")
		radiusStr := values.Get("radius")

		if latStr == "" || longStr == "" || freqStr == "" || radiusStr == "" {
			return false
		}

		// Validate latitude and longitude ranges.
		lat, err := strconv.ParseFloat(latStr, 64)
		if err != nil || lat < -90 || lat > 90 {
			return false
		}

		lon, err := strconv.ParseFloat(longStr, 64)
		if err != nil || lon < -180 || lon > 180 {
			return false
		}

		// Validate frequency: extract the first number from the freq string.
		freqRegex := regexp.MustCompile(`(\d+(\.\d+)?)`)
		freqMatch := freqRegex.FindStringSubmatch(freqStr)
		if freqMatch == nil {
			return false
		}
		freqVal, err := strconv.ParseFloat(freqMatch[1], 64)
		if err != nil || freqVal <= 0 {
			return false
		}

		// Validate radius: extract the first number from the radius string.
		radiusMatch := freqRegex.FindStringSubmatch(radiusStr)
		if radiusMatch == nil {
			return false
		}
		radiusVal, err := strconv.ParseFloat(radiusMatch[1], 64)
		if err != nil || radiusVal <= 0 {
			return false
		}

		// JS8 is more of a "demo" / "example". We include it to demonstrate that
		// overlays can be advertised in many, many ways.
		return true
	}

	// If none of the known prefixes match, the URI is not advertisable.
	return false
}

// IsValidTopicOrServiceName checks if the provided service name is valid based on BRC-87 guidelines.
// A valid topic or service name must:
// - Be between 1 and 50 characters
// - Start with "tm_" or "ls_"
// - Contain only lowercase letters and underscores
// - Follow the pattern of words separated by underscores
func IsValidTopicOrServiceName(service string) bool {
	serviceRegex := regexp.MustCompile(`^(?:tm_|ls_)[a-z]+(?:_[a-z]+)*$`)
	if !serviceRegex.MatchString(service) {
		return false
	}

	// Check length separately (1-50 characters)
	if len(service) < 1 || len(service) > 50 {
		return false
	}

	return true
}

// IsTokenSignatureCorrectlyLinked checks that the BRC-48 locking key and the signature are valid
// and linked to the claimed identity key.
// Parameters:
// - lockingPublicKey: The public key used in the output's locking script
// - fields: The fields of the PushDrop token for the SHIP or SLAP advertisement
// Returns:
// - true if the token's signature is properly linked to the claimed identity key
// - false otherwise
func IsTokenSignatureCorrectlyLinked(
	ctx context.Context,
	lockingPublicKey *ec.PublicKey,
	fields [][]byte,
) (bool, error) {
	if len(fields) < 3 {
		return false, nil
	}

	// The signature is the last field, which needs to be removed for verification
	fieldsWithoutSignature := fields[:len(fields)-1]

	// The protocol is in the first field
	protocolName := string(fieldsWithoutSignature[0])
	protocolID := &wallet.Protocol{
		SecurityLevel: 2,
	}
	if protocolName == "SHIP" {
		protocolID.Protocol = "service host interconnect"
	} else {
		protocolID.Protocol = "service lookup availability"
	}

	// The identity key is in the second field
	identityKey, err := ec.PublicKeyFromBytes(fieldsWithoutSignature[1])
	if err != nil {
		return false, err
	}

	// First, we ensure that the signature over the data is valid for the claimed identity key.
	// Concatenate all fields to create the message for verification
	data := make([]byte, 0, 128)
	for _, field := range fieldsWithoutSignature {
		data = append(data, field...)
	}

	// Verify the signature
	// messageHash := hash.Sha256d(data)
	signature, err := ec.ParseDERSignature(fields[len(fields)-1])
	if err != nil {
		return false, nil
	}

	anyoneWallet, err := wallet.NewProtoWallet(wallet.ProtoWalletArgs{
		Type: wallet.ProtoWalletArgsTypeAnyone,
	})
	if err != nil {
		return false, err
	}
	counterparty := wallet.Counterparty{
		Type:         wallet.CounterpartyTypeOther,
		Counterparty: identityKey,
	}
	if result, err := anyoneWallet.VerifySignature(ctx, wallet.VerifySignatureArgs{
		EncryptionArgs: wallet.EncryptionArgs{
			ProtocolID:   *protocolID,
			Counterparty: counterparty,
			KeyID:        "1",
		},
		Data:      data,
		Signature: *signature,
	}, ""); err != nil {
		return false, err
	} else if !result.Valid {
		return false, nil
	}

	// Then, we ensure that the locking public key matches the correct derived child.
	if expectedLockingPublickKey, err := anyoneWallet.GetPublicKey(ctx, wallet.GetPublicKeyArgs{
		EncryptionArgs: wallet.EncryptionArgs{
			ProtocolID:   *protocolID,
			Counterparty: counterparty,
			KeyID:        "1",
		},
	}, ""); err != nil {
		return false, err
	} else {
		return bytes.Equal(expectedLockingPublickKey.PublicKey.Compressed(), lockingPublicKey.Compressed()), nil
	}
}
