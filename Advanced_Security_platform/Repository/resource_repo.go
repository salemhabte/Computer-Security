package repository

import (
	domain "security/domain"
	"security/config"
	"context"
	"log"
	"errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type ResourceRepository struct {
	coll   *mongo.Collection
	ctx    context.Context
	client *mongo.Client
}

func NewResourceRepository() domain.IResourceRepository {
	conn, err := Connect()
	if err != nil {
		log.Fatal("can't init Resource repository")
	}
	db := config.USER_DB
	collection := conn.Client.Database(db).Collection("resources")
	return &ResourceRepository{
		coll:   collection,
		ctx:    context.TODO(),
		client: conn.Client,
	}
}

func (r *ResourceRepository) GetByID(resourceID string) (*domain.Resource, error) {
	var resource domain.Resource
	err := r.coll.FindOne(r.ctx, bson.M{"id": resourceID}).Decode(&resource)
	if err == mongo.ErrNoDocuments {
		return nil, errors.New("resource not found")
	}
	if err != nil {
		return nil, err
	}
	return &resource, nil
}

