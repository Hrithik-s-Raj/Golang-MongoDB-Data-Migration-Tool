

package main

import (
	"context"
	"fmt"
	"log"
	//"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const batchSize = 1000

func main() {

	// Configs
	sourceURI := "Source MongoDB URI"
	targetURI:= "Target MongoDb URI"
	sourceDBName := "Source Database Name"
	targetDBName := "Target Database Name"

	if sourceURI == "" || targetURI == "" || sourceDBName == "" || targetDBName == "" {
		log.Fatal("Please set SOURCE_URI, TARGET_URI, SOURCE_DB, and TARGET_DB environment variables")
	}

	// Create the Source and the Target Client
	sourceClient, err := mongo.NewClient(options.Client().ApplyURI(sourceURI))
	if err != nil {
		log.Fatalf("Failed to create source client: %v", err)
	}

	targetClient, err := mongo.NewClient(options.Client().ApplyURI(targetURI))
	if err != nil {
		log.Fatalf("Failed to create target client: %v", err)
	}

	// Creating the Context
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Minute)
	defer cancel()

	// Connect to source MongoDB
	err = sourceClient.Connect(ctx)
	if err != nil {
		log.Fatalf("Failed to connect to source MongoDB: %v", err)
	}
	defer func() {
		if err = sourceClient.Disconnect(ctx); err != nil {
			log.Fatalf("Error disconnecting source client: %v", err)
		}
	}()

	// Connect to target MongoDB
	err = targetClient.Connect(ctx)
	if err != nil {
		log.Fatalf("Failed to connect to target MongoDB: %v", err)
	}
	defer func() {
		if err = targetClient.Disconnect(ctx); err != nil {
			log.Fatalf("Error disconnecting target client: %v", err)
		}
	}()

	sourceDB := sourceClient.Database(sourceDBName)
	targetDB := targetClient.Database(targetDBName)

	collections, err := sourceDB.ListCollectionNames(ctx, bson.D{})
	if err != nil {
		log.Fatalf("Failed to list collections: %v", err)
	}

	for _, collectionName := range collections {
		fmt.Printf("Migrating collection: %s\n", collectionName)

		sourceCollection := sourceDB.Collection(collectionName)
		targetCollection := targetDB.Collection(collectionName)

		collectionCtx, collectionCancel := context.WithTimeout(context.Background(), 30*time.Minute)
		defer collectionCancel()

		cursor, err := sourceCollection.Find(collectionCtx, bson.D{})
		if err != nil {
			log.Fatalf("Failed to retrieve documents from collection %s: %v", collectionName, err)
		}

		defer cursor.Close(collectionCtx)

		var documents []interface{}
		count := 0

		for cursor.Next(collectionCtx) {
			var doc bson.M
			if err := cursor.Decode(&doc); err != nil {
				log.Fatalf("Failed to decode document: %v", err)
			}
			documents = append(documents, doc)
			count++

			if len(documents) == batchSize {
				if err := insertBatch(targetCollection, documents); err != nil {
					log.Fatalf("Failed to insert batch into collection %s: %v", collectionName, err)
				}
				documents = documents[:0] 
				fmt.Printf("Inserted %d documents into %s\n", count, collectionName)
			}
		}

		if len(documents) > 0 {
			if err := insertBatch(targetCollection, documents); err != nil {
				log.Fatalf("Failed to insert final batch into collection %s: %v", collectionName, err)
			}
			fmt.Printf("Inserted %d documents into %s\n", count, collectionName)
		}

		if err := cursor.Err(); err != nil {
			log.Fatalf("Cursor error in collection %s: %v", collectionName, err)
		}

		fmt.Printf("Successfully migrated collection: %s with %d documents\n", collectionName, count)
	}

	fmt.Println("Database migration completed successfully.")
}

func insertBatch(collection *mongo.Collection, documents []interface{}) error {
	_, err := collection.InsertMany(context.Background(), documents, options.InsertMany().SetOrdered(false))
	return err
}
