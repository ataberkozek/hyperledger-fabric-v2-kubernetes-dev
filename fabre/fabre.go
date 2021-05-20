/*
SPDX-License-Identifier: Apache-2.0
*/

package main

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// SmartContract provides functions for managing a real estate
type SmartContract struct {
	contractapi.Contract
}

// RealEstate describes basic details of what makes up a real estate
type RealEstate struct {
	Location    string `json:"location"`
	Rooms       string `json:"rooms"`
	Baths       string `json:"baths"`
	Price       string `json:"price"`
	LivingSpace string `json:"livingSpace"`
	Owner       string `json:"owner"`
}

// QueryResult structure used for handling result of query
type QueryResult struct {
	Key    string `json:"Key"`
	Record *RealEstate
}

// InitLedger adds a base set of Real Estates to the ledger
func (s *SmartContract) InitLedger(ctx contractapi.TransactionContextInterface) error {
	realEstates := []RealEstate{
		RealEstate{Location: "Istanbul", Rooms: "5", Baths: "2", Price: "$140,000", LivingSpace: "120m2", Owner: "Agency"},
		RealEstate{Location: "Izmir", Rooms: "3", Baths: "1", Price: "$75,000", LivingSpace: "90m2", Owner: "Agency"},
		RealEstate{Location: "Ankara", Rooms: "4", Baths: "2", Price: "$135,000", LivingSpace: "140m2", Owner: "Agency"},
		RealEstate{Location: "Istanbul", Rooms: "1", Baths: "1", Price: "$300,000", LivingSpace: "40m2", Owner: "Agency"},
		RealEstate{Location: "Bursa", Rooms: "5", Baths: "2", Price: "$100,000", LivingSpace: "200m2", Owner: "Agency"},
		RealEstate{Location: "Istanbul", Rooms: "2", Baths: "1", Price: "$55,000", LivingSpace: "80m2", Owner: "Agency"},
		RealEstate{Location: "Ankara", Rooms: "3", Baths: "1", Price: "$90,000", LivingSpace: "120m2", Owner: "Agency"},
		RealEstate{Location: "Istanbul", Rooms: "7", Baths: "3", Price: "$1,135,500", LivingSpace: "370m2", Owner: "Agency"},
		RealEstate{Location: "Izmir", Rooms: "2", Baths: "1", Price: "$55,000", LivingSpace: "80m2", Owner: "Agency"},
	}

	for i, re := range realEstates {
		reAsBytes, _ := json.Marshal(re)
		err := ctx.GetStub().PutState("RE"+strconv.Itoa(i), reAsBytes)

		if err != nil {
			return fmt.Errorf("Failed to put to world state. %s", err.Error())
		}
	}

	return nil
}

// AddRe adds a new Real Estate to the world state with given details
func (s *SmartContract) AddRe(ctx contractapi.TransactionContextInterface, reNumber string, location string, rooms string, baths string, price string, livingSpace string) error {
	re := RealEstate{
		Location:    location,
		Rooms:       rooms,
		Baths:       baths,
		Price:       price,
		LivingSpace: livingSpace,
		Owner:       "Agency",
	}

	reAsBytes, _ := json.Marshal(re)

	return ctx.GetStub().PutState(reNumber, reAsBytes)
}

// QueryRe returns the Real Estate stored in the world state with given id
func (s *SmartContract) QueryRe(ctx contractapi.TransactionContextInterface, reNumber string) (*RealEstate, error) {
	reAsBytes, err := ctx.GetStub().GetState(reNumber)

	if err != nil {
		return nil, fmt.Errorf("Failed to read from world state. %s", err.Error())
	}

	if reAsBytes == nil {
		return nil, fmt.Errorf("%s does not exist", reNumber)
	}

	re := new(RealEstate)
	_ = json.Unmarshal(reAsBytes, re)

	return re, nil
}

// QueryAllRes returns all Real Estates found in world state
func (s *SmartContract) QueryAllRes(ctx contractapi.TransactionContextInterface) ([]QueryResult, error) {
	startKey := ""
	endKey := ""

	resultsIterator, err := ctx.GetStub().GetStateByRange(startKey, endKey)

	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	results := []QueryResult{}

	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()

		if err != nil {
			return nil, err
		}

		re := new(RealEstate)
		_ = json.Unmarshal(queryResponse.Value, re)

		queryResult := QueryResult{Key: queryResponse.Key, Record: re}
		results = append(results, queryResult)
	}

	return results, nil
}

// ChangeReOwner updates the owner field of Real Estate with given id in world state
func (s *SmartContract) ChangeReOwner(ctx contractapi.TransactionContextInterface, reNumber string, newOwner string) error {
	re, err := s.QueryRe(ctx, reNumber)

	if err != nil {
		return err
	}

	re.Owner = newOwner

	reAsBytes, _ := json.Marshal(re)

	return ctx.GetStub().PutState(reNumber, reAsBytes)
}

// ChangeRePrice updates the price field of Real Estate with given id in world state
func (s *SmartContract) ChangeRePrice(ctx contractapi.TransactionContextInterface, reNumber string, newPrice string) error {
	re, err := s.QueryRe(ctx, reNumber)

	if err != nil {
		return err
	}

	re.Price = newPrice

	reAsBytes, _ := json.Marshal(re)

	return ctx.GetStub().PutState(reNumber, reAsBytes)
}

func main() {

	chaincode, err := contractapi.NewChaincode(new(SmartContract))

	if err != nil {
		fmt.Printf("Error create fabre chaincode: %s", err.Error())
		return
	}

	if err := chaincode.Start(); err != nil {
		fmt.Printf("Error starting fabre chaincode: %s", err.Error())
	}
}
