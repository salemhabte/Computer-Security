package repository

import (
	domain "security/domain"
	"security/config"
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type RefreshTokenRepository struct {
	Collection *mongo.Collection
	Contxt     context.Context
	Client     *mongo.Client
}

func NewRefreshTokenRepository() domain.IAuthRepo {
	userDB := config.USER_DB
	connection, err := Connect()
	if err != nil {
		log.Fatal("can't initailize ", "Refresh Token Repo", " Database")
	}
	collection := connection.Client.Database(userDB).Collection(config.USER_REFRESH_TOKEN_COLLECTION_NAME)
	return &RefreshTokenRepository{
		Collection: collection,
		Contxt:     context.TODO(),
		Client:     connection.Client,
	}
}

func (r *RefreshTokenRepository) Save(token *domain.RefreshToken) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{"user_id": token.UserID}
	update := bson.M{"$set": token}
	opts := options.Update().SetUpsert(true)

	_, err := r.Collection.UpdateOne(ctx, filter, update, opts)
	return err
}
func (r *RefreshTokenRepository) GetByToken(token string) (*domain.RefreshToken, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var rt domain.RefreshToken
	err := r.Collection.FindOne(ctx, bson.M{"token": token}).Decode(&rt)
	return &rt, err
}
func (r *RefreshTokenRepository) Delete(token string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := r.Collection.DeleteOne(ctx, bson.M{"token": token})
	return err
}

func (r *RefreshTokenRepository) CloseDataBase() error {
	if r.Client == nil {
		return nil // Nothing to close
	}
	if err := r.Client.Disconnect(r.Contxt); err != nil {
		return fmt.Errorf("error disconnecting from MongoDB: %w", err)
	}
	log.Println("Disconnected from Blog MongoDB.")
	return nil
}
