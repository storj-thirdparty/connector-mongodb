package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// ConfigMongoDB defines the variables and types.
type ConfigMongoDB struct {
	Hostname   string `json:"hostname"`
	Portnumber string `json:"port"`
	Username   string `json:"username"`
	Password   string `json:"password"`
	Database   string `json:"database"`
}

// MongoReader implements an io.Reader interface
type MongoReader struct {
	database          *mongo.Database
	collectionNames   []string
	lastDocumentIndex int
}

// Read reads collections' data as per the buffer capacity.
// It returns number of bytes (int) read and any error, if occurred.
// EOF error is returned after complete read.
func (mongoReader *MongoReader) Read(buf []byte) (int, error) {

	var err error
	ctx := context.TODO()

	bufferCapacity := cap(buf)
	filterBSON := bson.M{}

	if mongoReader.collectionNames == nil {
		// Retrieve ALL collections in the database.
		mongoReader.collectionNames, err = mongoReader.database.ListCollectionNames(ctx, filterBSON)
		if err != nil {
			log.Fatalf("Failed to retrieve collection names: %s\n", err)
			return 0, err
		}
	}

	var numOfBytesRead = 0
	// Go through ALL collections.
	for ij := range mongoReader.collectionNames {
		collection := mongoReader.database.Collection(mongoReader.collectionNames[ij])
		cursor, err := collection.Find(ctx, filterBSON)
		if err != nil {
			log.Printf("Failed to retrieve data about %s collection: %s\n", mongoReader.collectionNames[ij], err)
			return numOfBytesRead, err
		}
		defer cursor.Close(ctx)

		var documentCount, lastIndex int
		// Retrieve each document of the selected collection.
		for cursor.Next(ctx) {
			// Start reading from the last document that was under process.
			if documentCount >= mongoReader.lastDocumentIndex {
				// Convert JSON to BSON.
				rawDocumentBSON, _ := bson.Marshal(cursor.Current)
				documentSize := len(rawDocumentBSON)
				// Ensure required space is available in buf.
				if (numOfBytesRead + documentSize) < bufferCapacity {
					lastIndex = numOfBytesRead + len(rawDocumentBSON)
					copy(buf[numOfBytesRead:lastIndex], rawDocumentBSON)
					numOfBytesRead += documentSize
				} else {
					// Insufficient space in the buffer.
					err = io.ErrShortBuffer
					// Next time, start with the unprocessed collections.
					mongoReader.collectionNames = mongoReader.collectionNames[ij:]
					// Also, operate from the current document instance.
					mongoReader.lastDocumentIndex = documentCount

					return numOfBytesRead, err
				}
			}
			documentCount++
		}

		if err := cursor.Err(); err != nil {
			// Unexpected error occurred while processing cursors.
			// Next time, start with the unprocessed collections.
			mongoReader.collectionNames = mongoReader.collectionNames[ij:]
			// Also, operate from the current document instance.
			mongoReader.lastDocumentIndex = documentCount + 1
			// The caller needs to recall the Reader to fetch left-over data.
			return numOfBytesRead, err
		}

		fmt.Printf("\nAll data of the collections %s are uploaded!", mongoReader.collectionNames[ij])
		// All documents of the selected collection have been read.
		if mongoReader.lastDocumentIndex > 0 {
			// Reset the document index to be read from.
			mongoReader.lastDocumentIndex = 0
		}
	}
	// All collections have been read and processed.
	mongoReader.collectionNames = nil

	return numOfBytesRead, io.EOF
}

// LoadMongoProperty reads and parses the JSON file
// that contain a MongoDB instance's credentials.
// It returns all the properties embedded in a configuration object.
func LoadMongoProperty(fullFileName string) ConfigMongoDB {

	var configMongoDB ConfigMongoDB
	// Open and read the file
	fileHandle, err := os.Open(filepath.Clean(fullFileName))
	if err != nil {
		log.Fatal(err)
	}

	jsonParser := json.NewDecoder(fileHandle)
	err = jsonParser.Decode(&configMongoDB)
	if err != nil {
		log.Fatal(err)
	}

	if err = fileHandle.Close(); err != nil {
		log.Fatal(err)
	}

	// Display read information.
	fmt.Println("\nRead MongoDB configuration from the ", fullFileName, " file")
	fmt.Println("Hostname\t", configMongoDB.Hostname)
	fmt.Println("Portnumber\t", configMongoDB.Portnumber)
	fmt.Println("Username \t", configMongoDB.Username)
	fmt.Println("Password \t", configMongoDB.Password)
	fmt.Println("Database \t", configMongoDB.Database)

	return configMongoDB
}

// ConnectToDB will connect to a MongoDB instance based on the specified credentials.
// It returns a reference to an io.Reader with MongoDB instance information
func ConnectToDB(configMongoDB ConfigMongoDB) *MongoReader {

	fmt.Println("Connecting to MongoDB...")

	mongoURL := fmt.Sprintf("mongodb://%s:%s@%s:%s/%s?authSource="+configMongoDB.Database, configMongoDB.Username, configMongoDB.Password, configMongoDB.Hostname, configMongoDB.Portnumber, configMongoDB.Database)
	clientOptions := options.Client().ApplyURI(mongoURL)
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	// Check the connection with MongoDB.
	if err = client.Ping(context.TODO(), nil); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Successfully connected to MongoDB!")

	return &MongoReader{database: client.Database(configMongoDB.Database)}
}
