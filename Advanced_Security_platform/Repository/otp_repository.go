package repository

import (
	domain "security/domain"
	"security/config"
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type UserOTPRepository struct {
	collection *mongo.Collection
	Contxt     context.Context
	Client     *mongo.Client
}

func NewUserOTPRepository() domain.IUserOTPRepository {
	userDB := config.USER_DB
	connection, err := Connect()
	if err != nil {
		log.Fatal("can't initailize ", "OTP", " Database")
	}
	coll := connection.Client.Database(userDB).Collection(config.USER_OTP_COLLECTION_NAME)
	return &UserOTPRepository{
		collection: coll,
		Contxt:     context.TODO(),
		Client:     connection.Client,
	}
}

func (r *UserOTPRepository) StoreOTP(entry domain.UserUnverified) error {
	filter := bson.M{"email": entry.Email}
	update := bson.M{"$set": entry}
	opts := options.Update().SetUpsert(true)
	_, err := r.collection.UpdateOne(context.Background(), filter, update, opts)
	return err
}

func (r *UserOTPRepository) FindOTP(email string) (*domain.UserUnverified, error) {
	var entry domain.UserUnverified
	err := r.collection.FindOne(context.Background(), bson.M{"email": email}).Decode(&entry)
	if err != nil {
		return nil, err
	}
	return &entry, nil
}

func (r *UserOTPRepository) DeleteOTP(email string) error {
	_, err := r.collection.DeleteOne(context.Background(), bson.M{"email": email})
	return err
}

func (r *UserOTPRepository) CloseDataBase() error {
	if r.Client == nil {
		return nil // Nothing to close
	}
	if err := r.Client.Disconnect(r.Contxt); err != nil {
		return fmt.Errorf("error disconnecting from MongoDB: %w", err)
	}
	log.Println("Disconnected from Blog MongoDB.")
	return nil
}
