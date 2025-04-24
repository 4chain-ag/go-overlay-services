package ship

import (
	"context"
	"time"

	// "github.com/bsv-blockchain/go-overlay-services/pkg/discovery"
	"github.com/4chain-ag/go-overlay-services/pkg/discovery"
	"github.com/bsv-blockchain/go-sdk/chainhash"
	"github.com/bsv-blockchain/go-sdk/overlay"
	ec "github.com/bsv-blockchain/go-sdk/primitives/ec"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// SHIPStorage implements a storage engine for SHIP protocol
type SHIPStorage struct {
	shipRecords *mongo.Collection
}

// NewSHIPStorage constructs a new SHIPStorage instance
func NewSHIPStorage(db *mongo.Database) *SHIPStorage {
	return &SHIPStorage{
		shipRecords: db.Collection("shipRecords"),
	}
}

// EnsureIndexes creates the necessary indexes for the collections
func (s *SHIPStorage) EnsureIndexes(ctx context.Context) error {
	// Create compound index on domain and topic
	_, err := s.shipRecords.Indexes().CreateOne(
		ctx,
		mongo.IndexModel{
			Keys: bson.M{
				"domain": 1,
				"topic":  1,
			},
		},
	)
	return err
}

// StoreSHIPRecord stores a SHIP record in the database
func (s *SHIPStorage) StoreSHIPRecord(
	ctx context.Context,
	outpoint *overlay.Outpoint,
	identityKey *ec.PublicKey,
	domain string,
	topic string,
) error {
	record := discovery.SHIPRecord{
		Txid:        &outpoint.Txid,
		OutputIndex: outpoint.OutputIndex,
		IdentityKey: identityKey,
		Domain:      domain,
		Topic:       topic,
		CreatedAt:   time.Now(),
	}

	_, err := s.shipRecords.InsertOne(ctx, record)
	return err
}

// DeleteSHIPRecord deletes a SHIP record from the database
func (s *SHIPStorage) DeleteSHIPRecord(
	ctx context.Context,
	outpoint *overlay.Outpoint,
) error {
	filter := bson.M{
		"txid":        outpoint.Txid.String(),
		"outputIndex": outpoint.OutputIndex,
	}

	_, err := s.shipRecords.DeleteOne(ctx, filter)
	return err
}

// FindRecord finds SHIP records based on a given query object
func (s *SHIPStorage) FindRecord(
	ctx context.Context,
	query discovery.SHIPQuery,
) ([]*overlay.Outpoint, error) {
	mongoQuery := bson.M{}

	// Add domain to the query if provided
	if query.Domain != "" {
		mongoQuery["domain"] = query.Domain
	}

	// Add topics to the query if provided
	if len(query.Topics) > 0 {
		mongoQuery["topic"] = bson.M{"$in": query.Topics}
	}

	// Add identityKey to the query if provided
	if query.IdentityKey != nil {
		mongoQuery["identityKey"] = query.IdentityKey
	}

	// Create a projection to only return txid and outputIndex
	projection := bson.M{
		"txid":        1,
		"outputIndex": 1,
		"_id":         0,
	}

	opts := options.Find().SetProjection(projection)
	cursor, err := s.shipRecords.Find(ctx, mongoQuery, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []*overlay.Outpoint
	if err := cursor.All(ctx, &results); err != nil {
		return nil, err
	}

	return results, nil
}

// FindAll returns all SHIP records tracked by the overlay
func (s *SHIPStorage) FindAll(ctx context.Context) ([]*overlay.Outpoint, error) {
	// Create a projection to only return txid and outputIndex
	projection := bson.M{
		"txid":        1,
		"outputIndex": 1,
		"_id":         0,
	}

	opts := options.Find().SetProjection(projection)
	cursor, err := s.shipRecords.Find(ctx, bson.M{}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []*overlay.Outpoint
	if err := cursor.All(ctx, &results); err != nil {
		return nil, err
	}

	return results, nil
}

// DeleteByDomain deletes all SHIP records for a specific domain
func (s *SHIPStorage) DeleteByDomain(ctx context.Context, domain string) error {
	filter := bson.M{"domain": domain}
	_, err := s.shipRecords.DeleteMany(ctx, filter)
	return err
}

// GetUTXOByReference retrieves the full SHIP record for a given UTXO reference
func (s *SHIPStorage) GetRecordByReference(
	ctx context.Context,
	txid *chainhash.Hash,
	outputIndex uint32,
) (*discovery.SHIPRecord, error) {
	filter := bson.M{
		"txid":        txid,
		"outputIndex": outputIndex,
	}

	var result discovery.SHIPRecord
	err := s.shipRecords.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}
