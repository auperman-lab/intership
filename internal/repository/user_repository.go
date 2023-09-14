package repository

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"intership/internal/models"
	"log"
)

type UserRepository struct {
	collection *mongo.Collection
	sequence   *mongo.Collection
}

func NewUserRepository(client *mongo.Client, dbName, collectionName string) *UserRepository {
	db := client.Database(dbName)
	collection := db.Collection(collectionName)
	sequence := db.Collection("sequence")
	return &UserRepository{
		collection: collection,
		sequence:   sequence,
	}
}

func (repo *UserRepository) Insert(user *models.UserModel) error {
	_, err := repo.collection.InsertOne(context.TODO(), user)
	if err != nil {
		log.Printf("failed to insert user in database: %v\n", err)
		return err
	}
	return nil
}

func (repo *UserRepository) CheckUserEmail(userEmail string) (*models.UserModel, error) {
	var user models.UserModel
	err := repo.collection.FindOne(context.TODO(), bson.M{"email": userEmail}).Decode(&user)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}

	return &user, nil
}
