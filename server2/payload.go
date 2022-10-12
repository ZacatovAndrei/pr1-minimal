package main

type Payload struct {
	Id          int    `json:"id"`
	GeneratorId int    `json:"generator_id"`
	ConsumerId  int    `json:"consumer_id"`
	Payload     string `json:"payload"`
}
