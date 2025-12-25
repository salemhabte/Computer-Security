package repository

import (
	"context"
	"errors"
	"log"
	"security/config"
	domain "security/domain"

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

func (r *ResourceRepository) Create(resource *domain.Resource) error {
	_, err := r.coll.InsertOne(r.ctx, resource)
	return err
}

func (r *ResourceRepository) GetByOwner(ownerEmail string) ([]*domain.Resource, error) {
	cursor, err := r.coll.Find(r.ctx, bson.M{"owneremail": ownerEmail})
	if err != nil {
		return nil, err
	}
	var resources []*domain.Resource
	if err = cursor.All(r.ctx, &resources); err != nil {
		return nil, err
	}
	return resources, nil
}

func (r *ResourceRepository) GetByIDs(ids []string) ([]*domain.Resource, error) {
	if len(ids) == 0 {
		return []*domain.Resource{}, nil
	}
	filter := bson.M{"id": bson.M{"$in": ids}}
	cursor, err := r.coll.Find(r.ctx, filter)
	if err != nil {
		return nil, err
	}
	var resources []*domain.Resource
	if err = cursor.All(r.ctx, &resources); err != nil {
		return nil, err
	}
	return resources, nil
}
