package repository

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type AdminRepository struct {
	collection *mongo.Collection
}

func NewAdminRepository(client *mongo.Client, dbName, collectionName string) *AdminRepository {
	db := client.Database(dbName)
	collection := db.Collection(collectionName)
	return &AdminRepository{
		collection: collection,
	}
}

func (repo *AdminRepository) GetUsersID() ([]string, error) {
	filter := bson.M{}

	projection := bson.M{"_id": 1}

	cursor, err := repo.collection.Find(context.Background(), filter, options.Find().SetProjection(projection))
	if err != nil {
		fmt.Println("error in repository", err.Error())
		return nil, err
	}
	defer cursor.Close(context.Background())

	var ids []primitive.ObjectID

	for cursor.Next(context.Background()) {
		var result bson.M
		if err := cursor.Decode(&result); err != nil {
			fmt.Println("error in repository2", err.Error())
			return nil, err
		}
		id, ok := result["_id"].(primitive.ObjectID)
		if ok {
			ids = append(ids, id)
		}
	}

	if err := cursor.Err(); err != nil {
		fmt.Println("error in repository3", err.Error())
		return nil, err
	}

	response := make([]string, 0, len(ids))
	for _, id := range ids {
		response = append(response, id.Hex())
	}

	return response, err

}

func (repo *AdminRepository) AssignAdmin(ID primitive.ObjectID) error {
	filter := bson.M{"_id": ID}

	update := bson.M{
		"$set": bson.M{"role": []string{"admin"}},
	}
	_, err := repo.collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		fmt.Println("cant find", err)
		return err
	}

	return nil
}
