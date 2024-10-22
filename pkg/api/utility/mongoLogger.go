package utility

import (
	"context"
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap/zapcore"
)

type MongoDBCore struct {
	collection *mongo.Collection
}

type MongoDBWriteSyncer struct {
	core *MongoDBCore
}

func NewMongoDBCore(collection *mongo.Collection) (zapcore.Core, error) {
	level := zapcore.InfoLevel

	core := &MongoDBWriteSyncer{&MongoDBCore{collection: collection}}
	return zapcore.NewCore(zapcore.NewJSONEncoder(zapcore.EncoderConfig{
		TimeKey:        "timestamp",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}), core, level), nil
}

func (m MongoDBCore) Write(entry zapcore.Entry, fields []zapcore.Field) error {
	logEntry := bson.M{
		"level":     entry.Level.String(),
		"message":   entry.Message,
		"timestamp": entry.Time,
		"caller":    entry.Caller.TrimmedPath(),
		"stack":     entry.Stack,
	}

	for _, field := range fields {
		logEntry[field.Key] = field.Interface
	}

	_, err := m.collection.InsertOne(context.TODO(), logEntry)
	return err
}

func (m MongoDBCore) Sync() error {
	return nil
}

func (m *MongoDBWriteSyncer) Write(p []byte) (n int, err error) {
	var entry zapcore.Entry
	err = json.Unmarshal(p, &entry)
	if err != nil {
		return 0, err
	}
	return len(p), m.core.Write(entry, nil)
}

func (m *MongoDBWriteSyncer) Sync() error {
	return nil
}
