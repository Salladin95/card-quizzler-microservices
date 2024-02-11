package models

import (
	"context"
	"github.com/Salladin95/card-quizzler-microservices/logging-service/cmd/api/entities"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
)

var client *mongo.Client

func NewModels(mongo *mongo.Client) Models {
	client = mongo

	return Models{
		LogEntry: LogEntry{},
	}
}

type Models struct {
	LogEntry LogEntry
}

type LogEntry struct {
	ID          uuid.UUID `bson:"_id" json:"id"`
	Message     string    `bson:"message" json:"message"`
	Level       string    `bson:"level" json:"level"`
	FromService string    `bson:"fromService" json:"fromService"`
	CreatedAt   time.Time `bson:"createdAt" json:"createdAt"`
	Method      string    `bson:"method,omitempty" json:"method,omitempty"`
	Name        string    `bson:"name,omitempty" json:"name,omitempty"`
}

func (l *LogEntry) getCollection() *mongo.Collection {
	return client.Database("logs").Collection("logs")
}

func (l *LogEntry) Insert(ctx context.Context, entry entities.LogMessage) error {
	collection := l.getCollection()

	_, err := collection.InsertOne(ctx, LogEntry{
		Name:        entry.Name,
		Message:     entry.Message,
		FromService: entry.FromService,
		Level:       entry.Level,
		Method:      entry.Method,
		CreatedAt:   time.Now(),
		ID:          uuid.New(),
	})
	if err != nil {
		log.Println("Error inserting into logs:", err)
		return err
	}

	return nil
}

func (l *LogEntry) All() ([]*LogEntry, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	collection := client.Database("logs").Collection("logs")

	opts := options.Find()
	opts.SetSort(bson.D{{"created_at", -1}})

	cursor, err := collection.Find(context.TODO(), bson.D{}, opts)
	if err != nil {
		log.Println("Finding all docs error:", err)
		return nil, err
	}
	defer cursor.Close(ctx)

	var logs []*LogEntry

	for cursor.Next(ctx) {
		var item LogEntry

		err := cursor.Decode(&item)
		if err != nil {
			log.Print("Error decoding log into slice:", err)
			return nil, err
		} else {
			logs = append(logs, &item)
		}
	}

	return logs, nil
}

func (l *LogEntry) GetOne(id string) (*LogEntry, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	collection := client.Database("logs").Collection("logs")

	docID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var entry LogEntry
	err = collection.FindOne(ctx, bson.M{"_id": docID}).Decode(&entry)
	if err != nil {
		return nil, err
	}

	return &entry, nil
}

func (l *LogEntry) DropCollection() error {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	collection := client.Database("logs").Collection("logs")

	if err := collection.Drop(ctx); err != nil {
		return err
	}

	return nil
}
