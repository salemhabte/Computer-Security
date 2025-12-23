package repository

import (
	domain "security/domain"
	"security/config"
	"context"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type RoleRepository struct {
	coll   *mongo.Collection
	ctx    context.Context
	client *mongo.Client
}

func NewRoleRepository() domain.IRoleRepository {
	conn, err := Connect()
	if err != nil {
		log.Fatal("can't init Role repository")
	}
	db := config.USER_DB
	collection := conn.Client.Database(db).Collection("roles")
	return &RoleRepository{
		coll:   collection,
		ctx:    context.TODO(),
		client: conn.Client,
	}
}

func (r *RoleRepository) Upsert(role domain.Role) error {
	filter := bson.M{"name": role.Name}
	update := bson.M{"$set": role}
	opts := options.Update().SetUpsert(true)
	_, err := r.coll.UpdateOne(r.ctx, filter, update, opts)
	return err
}

func (r *RoleRepository) Get(name string) (*domain.Role, error) {
	var role domain.Role
	err := r.coll.FindOne(r.ctx, bson.M{"name": name}).Decode(&role)
	return &role, err
}

