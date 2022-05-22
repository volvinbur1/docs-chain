package storage

import (
	"context"
	"fmt"
	"github.com/volvinbur1/docs-chain/src/common"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"log"
	"time"
)

const mongoServerUri = "mongodb://localhost:27017"

const (
	dbName               = "docsChain"
	papersCollectionName = "papersList"
)

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

func (d *DatabaseManager) AddNewPaper(newPaper common.UploadedPaper) error {
	if err := d.pingServer(); err != nil {
		return err
	}

	collection := d.client.Database(dbName).Collection(papersCollectionName)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := collection.InsertOne(ctx, newPaper)
	if err != nil {
		return fmt.Errorf("newly uploaded paper insertion to database failed: %s", err)
	}
	return nil
}

func (d *DatabaseManager) pingServer() error {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	return d.client.Ping(ctx, readpref.Primary())
}
