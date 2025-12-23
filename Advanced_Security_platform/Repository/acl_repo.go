package repository

import (
	domain "security/domain"
	"security/config"
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ACLRepository struct {
	coll   *mongo.Collection
	ctx    context.Context
	client *mongo.Client
}

func NewACLRepository() domain.IACLRepository {
	conn, err := Connect()
	if err != nil {
		log.Fatal("can't init ACL repository")
	}
	db := config.USER_DB
	collection := conn.Client.Database(db).Collection("acls")
	return &ACLRepository{
		coll:   collection,
		ctx:    context.TODO(),
		client: conn.Client,
	}
}

func (r *ACLRepository) Grant(entry domain.AccessControlEntry) error {
	if entry.GrantedAt.IsZero() {
		entry.GrantedAt = time.Now()
	}
	filter := bson.M{"resourceid": entry.ResourceID, "subject": entry.Subject}
	update := bson.M{"$set": entry}
	opts := options.Update().SetUpsert(true)
	_, err := r.coll.UpdateOne(r.ctx, filter, update, opts)
	return err
}

func (r *ACLRepository) Revoke(resourceID, subject string) error {
	_, err := r.coll.DeleteOne(r.ctx, bson.M{"resourceid": resourceID, "subject": subject})
	return err
}

func (r *ACLRepository) Get(resourceID, subject string) (*domain.AccessControlEntry, error) {
	var entry domain.AccessControlEntry
	err := r.coll.FindOne(r.ctx, bson.M{"resourceid": resourceID, "subject": subject}).Decode(&entry)
	return &entry, err
}

