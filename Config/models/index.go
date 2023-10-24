package models

import pb "Config/proto"

type Application struct {
	Id     string   `json:"id" bson:"id"`
	Name   string   `json:"name" bson:"name"`
	Config []Config `json:"config" bson:"config"`
}

type Config struct {
	Key   string `json:"key" bson:"key"`
	Value string `json:"value" bson:"value"`
}

type Config1 struct {
	Id    string `json:"id" bson:"id"`
	Key   string `json:"key" bson:"key"`
	Value string `json:"value" bson:"value"`
}

type ConfigRequest struct {
	Id    string     `json:"id"`
	Key   string     `json:"key"`
	Value *pb.Config `json:"value"`
}
