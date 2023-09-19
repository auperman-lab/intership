package repository

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"intership/internal/models"
)

type UserRepository struct {
	collection *mongo.Collection
}

func NewUserRepository(client *mongo.Client, dbName, collectionName string) *UserRepository {
	db := client.Database(dbName)
	collection := db.Collection(collectionName)
	return &UserRepository{
		collection: collection,
	}
}

func (repo *UserRepository) CheckUserId(userID primitive.ObjectID) (models.UserModel, error) {
	var user models.UserModel
	filter := bson.M{"_id": userID}

	if err := repo.collection.FindOne(context.Background(), filter).Decode(&user); err != nil {
		if err == mongo.ErrNoDocuments {
			return user, nil
		}
		return user, err
	}
	return user, nil
}

func (repo *UserRepository) Delete(userID primitive.ObjectID) error {
	filter := bson.M{"_id": userID}
	result, err := repo.collection.DeleteOne(context.Background(), filter)
	if err != nil {
		return err
	}
	if result.DeletedCount == 0 {
		return mongo.ErrNoDocuments
	}
	return nil
}
