package controllers

import (
	"Config/interfaces"
	"Config/models"
	pro "Config/proto"
	"Config/services"
	"context"
	"fmt"
	"log"

	"google.golang.org/protobuf/encoding/protojson"
)

type RPCServer struct {
	pro.UnimplementedApplicationServiceServer
}

var (
	// ctx           gin.Context
	Configservice interfaces.Config
)

// -------->

func (s *RPCServer) AddApplication(ctx context.Context, req *pro.Application) (*pro.Application, error) {
	configdata := models.Config{
		Key:   req.Config.Key,
		Value: req.Config.Value,
	}
	data := models.Application{
		Id:     services.GenerateUniqueCustomerID(),
		Name:   req.Name,
		Config: []models.Config{configdata},
	}
	result, err := Configservice.AddApplication(&data)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	config := &pro.Config{
		Key:   result.Config[0].Key,
		Value: result.Config[0].Value,
	}
	response := &pro.Application{
		Name:   result.Name,
		Id:     result.Id,
		Config: config,
	}
	return response, nil
}

// --------->

func (s *RPCServer) UpdateConfig(ctx context.Context, req *pro.UpdateConfigRequest) (*pro.Config, error) {
	dbdata := models.Config1{
		Id:    req.ApplicationId,
		Key:   req.Key,
		Value: req.Value,
	}
	result, err := Configservice.UpdateConfig(&dbdata)
	if err != nil {
		return nil, err
	}
	response := &pro.Config{
		Key:   result.Key,
		Value: result.Value,
	}

	return response, nil

}

// --------->

func (s *RPCServer) AddConfig(ctx context.Context, req *pro.AddConfigRequest) (*pro.Empty, error) {
	dbdata := models.Config1{
		Id:    req.Id,
		Key:   req.Key,
		Value: req.Value,
	}
	err := Configservice.AddConfig(&dbdata)
	if err != nil {
		return nil, err
	}
	return &pro.Empty{}, nil
}

// ------->

func (s *RPCServer) GetConfigValue(ctx context.Context, req *pro.GetConfigValueRequest) (*pro.GetConfigValueResponse, error) {
	dbdata := models.Config1{
		Id:  req.GetId(),
		Key: req.GetKey(),
	}
	fmt.Println("000")
	fmt.Println(dbdata)
	responses, err := Configservice.GetConfigValue(&dbdata)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	var results []*pro.Config1
	for _, response := range responses {
		result := &pro.Config1{
			Id:    response.Id,
			Key:   response.Key,
			Value: response.Value,
		}
		results = append(results, result)
	}
	response := &pro.GetConfigValueResponse{Configs: results}

	marshaler := protojson.MarshalOptions{EmitUnpopulated: true}
	jsonStr, err := marshaler.Marshal(response)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	fmt.Println(response)
	fmt.Println("111")
	fmt.Println(string(jsonStr))

	return response, nil
}
