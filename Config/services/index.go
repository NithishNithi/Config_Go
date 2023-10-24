package services

import (
	"Config/interfaces"
	"Config/models"
	"context"
	"errors"

	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type CustomerService struct {
	ConfigCollection *mongo.Collection
	ctx              context.Context
}

func InitCustomerService(collection *mongo.Collection, ctx context.Context) interfaces.Config {
	return &CustomerService{collection, ctx}
}

func (p *CustomerService) AddApplication(request *models.Application) (*models.Application, error) {
	result, err := p.ConfigCollection.InsertOne(p.ctx, request)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	filter := bson.M{"_id": result.InsertedID}

	err = p.ConfigCollection.FindOne(p.ctx, filter).Decode(&request)
	if err != nil {
		return nil, err
	}
	return request, nil
}

func (p *CustomerService) UpdateConfig(request *models.Config1) (*models.Config, error) {
	fmt.Println("1")
	filter := bson.M{
		"id":         request.Id,
		"config.key": request.Key,
	}
	update := bson.M{
		"$set": bson.M{
			"config.$.value": request.Value,
		},
	}
	_, err := p.ConfigCollection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return nil, err
	}
	var updatedConfig *models.Config
	err = p.ConfigCollection.FindOne(context.TODO(), filter).Decode(&updatedConfig)
	if err != nil {
		return nil, err
	}
	return updatedConfig, nil
}

// -------->

func (p *CustomerService) AddConfig(request *models.Config1) error {
	filter := bson.M{"id": request.Id}
	var existingmodels *models.Application
	err := p.ConfigCollection.FindOne(p.ctx, filter).Decode(&existingmodels)
	if err != nil {
		log.Fatal(err)
	}
	newconfig := models.Config{
		Key:   request.Key,
		Value: request.Value,
	}
	existingmodels.Config = append(existingmodels.Config, newconfig)
	update := bson.M{"$set": bson.M{"config": existingmodels.Config}}
	options := options.Update()
	_, err = p.ConfigCollection.UpdateOne(p.ctx, filter, update, options)
	if err != nil {
		log.Fatal(err)
		return err
	}
	return nil
}

// ------->
func (p *CustomerService) GetConfigValue(request *models.Config1) ([]*models.Config1, error) {
	filter := bson.M{"id": request.Id}
	if request.Key != "" {
		filter["config.key"] = request.Key
	}
	var result models.Application
	err := p.ConfigCollection.FindOne(p.ctx, filter).Decode(&result)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	if len(result.Config) > 0 {
		var resultConfigs []*models.Config1
		for _, foundConfig := range result.Config {
			resultConfig := &models.Config1{
				Id:    request.Id,
				Key:   foundConfig.Key,
				Value: foundConfig.Value,
			}
			resultConfigs = append(resultConfigs, resultConfig)
		}
		return resultConfigs, nil
	}
	return nil, errors.New("configuration not found")
}
