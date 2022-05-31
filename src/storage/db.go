package storage

import (
	"context"
	"fmt"
	"github.com/volvinbur1/docs-chain/src/common"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"log"
	"time"
)

const mongoServerUri = "mongodb://localhost:27017"

const (
	dbName                   = "docsChain"
	papersCollection         = "papersList"
	papersShinglesCollection = "papersShingles"
	papersNftCollection      = "papersNft"
)

type DatabaseInterface interface {
	AddNewPaper(newPaper common.PaperMetadata) error
	AddPaperShingles(paperShingles common.PaperShingles) error
	GetAllPapersShingles() ([]common.PaperShingles, error)
	AddPaperNft(paperNft common.NftResponse) error
	GetPaperNftById(paperId string) (string, error)
}

type DatabaseManager struct {
	client *mongo.Client
}

func NewDatabaseManager() *DatabaseManager {
	m := &DatabaseManager{}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	var err error
	m.client, err = mongo.Connect(ctx, options.Client().ApplyURI(mongoServerUri))
	if err != nil {
		log.Fatalln("Database connection establishing failed:", err)
		return nil
	}

	return m
}

func (d *DatabaseManager) Disconnect() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := d.client.Disconnect(ctx); err != nil {
		log.Fatalln("Database disconnect failed. Error:", err)
	}
}

func (d *DatabaseManager) AddNewPaper(newPaper common.PaperMetadata) error {
	if err := d.pingServer(); err != nil {
		return fmt.Errorf("mongo db ping error: %s", err)
	}

	collection := d.client.Database(dbName).Collection(papersCollection)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := collection.InsertOne(ctx, newPaper)
	if err != nil {
		return fmt.Errorf("new paper metadata insertion to database failed: %s", err)
	}
	return nil
}

func (d *DatabaseManager) AddPaperShingles(paperShingles common.PaperShingles) error {
	if err := d.pingServer(); err != nil {
		return fmt.Errorf("mongo db ping error: %s", err)
	}

	collection := d.client.Database(dbName).Collection(papersShinglesCollection)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := collection.InsertOne(ctx, paperShingles)
	if err != nil {
		return fmt.Errorf("paper shingles insertion to database failed: %s", err)
	}
	return nil
}

func (d *DatabaseManager) GetAllPapersShingles() ([]common.PaperShingles, error) {
	if err := d.pingServer(); err != nil {
		return nil, fmt.Errorf("mongo db ping error: %s", err)
	}

	collection := d.client.Database(dbName).Collection(papersShinglesCollection)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cursor, err := collection.Find(ctx, bson.D{})
	if err != nil {
		return nil, fmt.Errorf("all paper shingles retriviation from database failed: %s", err)
	}

	paperShingles := make([]common.PaperShingles, 0)
	if err = cursor.All(ctx, &paperShingles); err != nil {
		return nil, fmt.Errorf("all paper shingles decode failed: %s", err)
	}
	return paperShingles, nil
}

func (d *DatabaseManager) AddPaperNft(paperNft common.NftResponse) error {
	if err := d.pingServer(); err != nil {
		return fmt.Errorf("mongo db ping error: %s", err)
	}

	collection := d.client.Database(dbName).Collection(papersNftCollection)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := collection.InsertOne(ctx, paperNft)
	if err != nil {
		return fmt.Errorf("paper nft insertion to database failed: %s", err)
	}
	return nil
}

func (d *DatabaseManager) GetPaperNftById(paperId string) (string, error) {
	if err := d.pingServer(); err != nil {
		return "", fmt.Errorf("mongo db ping error: %s", err)
	}

	collection := d.client.Database(dbName).Collection(papersShinglesCollection)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cursor, err := collection.Find(ctx, bson.D{{"id", paperId}})
	if err != nil {
		return "", fmt.Errorf("paper nft retriviation from database failed: %s", err)
	}

	var nft string
	if err = cursor.Decode(&nft); err != nil {
		return "", fmt.Errorf("paper nft retriviation from database failed: %s", err)
	}
	return nft, nil
}

func (d *DatabaseManager) pingServer() error {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	return d.client.Ping(ctx, readpref.Primary())
}
