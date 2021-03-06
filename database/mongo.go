package database

import (
	"context"
	"fmt"
	"log"
	"regexp"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

//This repository implementation is not really reliable if the user wants to use multiple drivers (postgresql/etc) or have different types of structs moving around. Made for simplicity

type (
	MongoString string //to prevent accidental string input, those should be pre defined or explicitly typecasted.

	Repository interface {
		UpdateOrInsertOne(database MongoString, collecton MongoString, input Document) error
		RemoveAllDocumentsInCollection(database MongoString, collecton MongoString) error
		FindOneByDate(database MongoString, collecton MongoString, date string) (MongoDocument, error)
	}

	database struct { //internal struct
		client *mongo.Client
		ctx    context.Context
	}
)

var (
	DaatabaseName    MongoString    = "blockchain"
	CollectionName   MongoString    = "binance"
	TestDatabaseName MongoString    = "test_blockchain"
	DateFormat       *regexp.Regexp = regexp.MustCompile(`\d{1,2}.\d{1,2}.\d{4}`) //DD.MM.YYYY
	ctx                             = context.TODO()                              //should replace with a proper context at some point
)

//New Initializes mongodb connection and returns a `Repository` interface
func New(connURL string) (Repository, error) {

	clientOptions := options.Client().ApplyURI(connURL)
	clientOptions = clientOptions.SetMaxPoolSize(50)
	connctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	dbClient, err := mongo.Connect(connctx, clientOptions)
	if err != nil {
		return &database{}, fmt.Errorf("failed to connect to mongodb: %e", err)
	}
	// Check the connection
	err = dbClient.Ping(connctx, nil)
	if err != nil {
		return &database{}, fmt.Errorf("failed to ping mongodb: %e", err)
	}

	log.Println("Connected to mongodb database") //TODO: replace with custom logger
	return &database{
		client: dbClient,
	}, nil
}

//RemoveAllDocumentsInCollection deletes all documents in a given collection, used for tests
func (db *database) RemoveAllDocumentsInCollection(database MongoString, collecton MongoString) error {
	collection := db.client.Database(string(database)).Collection(string(collecton))
	opts := options.Delete()
	_, err := collection.DeleteMany(ctx, bson.D{}, opts)
	if err != nil {
		return fmt.Errorf("Failed to drop collection: %e", err)
	}
	return nil
}

//UpdateOrInsertOne inserts or updates a bson-complaint object into the database
func (db *database) UpdateOrInsertOne(database MongoString, collecton MongoString, input Document) error {
	//check the date

	if !DateFormat.MatchString(input.Date) {
		return fmt.Errorf("Unexpected date format, GOT: %s", input.Date)
	}

	if input.Price == "" || input.Time == "" {
		return fmt.Errorf("%s", "Insufficient data was provided")
	}

	collection := db.client.Database(string(database)).Collection(string(collecton))

	//check if the current day is already present in the db (could be optimized with use of something like collection.FindOneAndUpdate(), but should be fine since this is not performance hungry)
	opts := options.FindOne()
	filter := bson.D{primitive.E{
		Key: "date",
		Value: bson.D{primitive.E{
			Key:   "$in",
			Value: bson.A{input.Date},
		}}}}
	res := collection.FindOne(db.ctx, filter, opts)
	if res.Err() != nil {
		if res.Err().Error() == "mongo: no documents in result" {
			//new day
			var stamps []MongoStamp
			stamps = append(stamps, MongoStamp{
				Time: input.Time, Price: input.Price,
			})
			doc := MongoDocument{
				Date:   input.Date,
				Stamps: stamps,
			}

			if _, err := collection.InsertOne(db.ctx, doc); err != nil {
				return fmt.Errorf("Failed to insert new document: %e", err)
			}
			return nil
		}
		//unexpected err
		return fmt.Errorf("Failed to FindOne: %e", res.Err())
	}
	//doc for current day exist, update the timefield
	out := new(Document)
	if err := res.Decode(out); err != nil {
		return fmt.Errorf("Failed to decode doc: %e", err)
	}
	updateOpts := options.FindOneAndUpdate()

	stamp := MongoStamp{
		Time:  input.Time,
		Price: input.Price,
	}
	//payload
	update := bson.M{"$push": bson.M{"stamps": stamp}}

	updRes := collection.FindOneAndUpdate(ctx, filter, update, updateOpts)
	if updRes.Err() != nil {
		return fmt.Errorf("Failed to FindOneAndUpdate: %e", updRes.Err())
	}

	return nil
}

//FindOneByDate retrieves a document by given date (format: DD.MM.YYYY)
func (db *database) FindOneByDate(database MongoString, collecton MongoString, date string) (MongoDocument, error) {
	collection := db.client.Database(string(database)).Collection(string(collecton))
	opts := options.FindOne()
	filter := bson.D{primitive.E{
		Key: "date",
		Value: bson.D{primitive.E{
			Key:   "$in",
			Value: bson.A{date},
		}}}}
	res := collection.FindOne(db.ctx, filter, opts)
	if res.Err() != nil {
		if res.Err().Error() == "mongo: no documents in result" {
			return MongoDocument{}, fmt.Errorf("Failed to find the document with desired date")
		}
		return MongoDocument{}, fmt.Errorf("Unexpected error on finding the document: %e", res.Err())
	}
	out := MongoDocument{}

	if err := res.Decode(&out); err != nil {
		return MongoDocument{}, fmt.Errorf("Failed to decode document: %e", err)
	}

	return out, nil
}
