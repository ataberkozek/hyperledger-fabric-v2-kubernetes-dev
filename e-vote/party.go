package main

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// SmartContract provides functions for managing an Party
type SmartContract struct {
	contractapi.Contract
}

// Party describes basic details of what makes up a simple Party
type Party struct {
	PartyName      string `json:"Partyname"`
	VoteCount      int `json:"VoteCount"`
}

// InitLedger adds a base set of Partys to the ledger
func (s *SmartContract) InitLedger(ctx contractapi.TransactionContextInterface) error {
	Partys := []Party{
		{PartyName: "Republicans", VoteCount: 0},
		{PartyName: "Democrats", VoteCount: 0},
	}

	for _, Party := range Partys {
		PartyJSON, err := json.Marshal(Party)
		if err != nil {
			return err
		}

		err = ctx.GetStub().PutState(Party.PartyName, PartyJSON)
		if err != nil {
			return fmt.Errorf("failed to put to world state. %v", err)
		}
	}

	return nil
}

// CreateParty issues a new Party to the world state with given details.
func (s *SmartContract) CreateParty(ctx contractapi.TransactionContextInterface, PartyName string) error {
	exists, err := s.PartyExists(ctx, PartyName)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("the Party %s already exists", PartyName)
	}

	Party := Party{
		PartyName:      PartyName,
		VoteCount:      0,
	}
	PartyJSON, err := json.Marshal(Party)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(PartyName, PartyJSON)
}

// ReadParty returns the Party stored in the world state with given id.
func (s *SmartContract) ReadParty(ctx contractapi.TransactionContextInterface, PartyName string) (*Party, error) {
	PartyJSON, err := ctx.GetStub().GetState(PartyName)
	if err != nil {
		return nil, fmt.Errorf("failed to read from world state: %v", err)
	}
	if PartyJSON == nil {
		return nil, fmt.Errorf("the Party %s does not exist", PartyName)
	}

	var Party Party
	err = json.Unmarshal(PartyJSON, &Party)
	if err != nil {
		return nil, err
	}

	return &Party, nil
}

// UpdateParty updates an existing Party in the world state with provided parameters.
func (s *SmartContract) UpdateParty(ctx contractapi.TransactionContextInterface, PartyName string, VoteCount int) error {
	exists, err := s.PartyExists(ctx, PartyName)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("the Party %s does not exist", PartyName)
	}

	// overwriting original Party with new Party
	Party := Party{
		PartyName:      PartyName,
		VoteCount:      VoteCount,
	}
	PartyJSON, err := json.Marshal(Party)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(PartyName, PartyJSON)
}

// DeleteParty deletes an given Party from the world state.
func (s *SmartContract) DeleteParty(ctx contractapi.TransactionContextInterface, PartyName string) error {
	exists, err := s.PartyExists(ctx, PartyName)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("the Party %s does not exist", PartyName)
	}

	return ctx.GetStub().DelState(PartyName)
}

// PartyExists returns true when Party with given ID exists in world state
func (s *SmartContract) PartyExists(ctx contractapi.TransactionContextInterface, PartyName string) (bool, error) {
	PartyJSON, err := ctx.GetStub().GetState(PartyName)
	if err != nil {
		return false, fmt.Errorf("failed to read from world state: %v", err)
	}

	return PartyJSON != nil, nil
}

// TransferParty updates the owner field of Party with given id in world state.
func (s *SmartContract) TransferParty(ctx contractapi.TransactionContextInterface, PartyName string, VoteCount int) error {
	Party, err := s.ReadParty(ctx, PartyName)
	if err != nil {
		return err
	}

	Party.VoteCount = VoteCount
	PartyJSON, err := json.Marshal(Party)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(PartyName, PartyJSON)
}

// CastVote updates the owner field of Party with given id in world state.
func (s *SmartContract) CastVote(ctx contractapi.TransactionContextInterface, PartyName string) error {
	Party, err := s.ReadParty(ctx, PartyName)
	if err != nil {
		return err
	}
	
	Party.VoteCount = Party.VoteCount + 1
	PartyJSON, err := json.Marshal(Party)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(PartyName, PartyJSON)
}

// GetAllPartys returns all Partys found in world state
func (s *SmartContract) GetAllPartys(ctx contractapi.TransactionContextInterface) ([]*Party, error) {
	// range query with empty string for startKey and endKey does an
	// open-ended query of all Partys in the chaincode namespace.
	resultsIterator, err := ctx.GetStub().GetStateByRange("", "")
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	var Partys []*Party
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		var Party Party
		err = json.Unmarshal(queryResponse.Value, &Party)
		if err != nil {
			return nil, err
		}
		Partys = append(Partys, &Party)
	}

	return Partys, nil
}

func main() {

	chaincode, err := contractapi.NewChaincode(new(SmartContract))

	if err != nil {
		fmt.Printf("Error create fabcar chaincode: %s", err.Error())
		return
	}

	if err := chaincode.Start(); err != nil {
		fmt.Printf("Error starting fabcar chaincode: %s", err.Error())
	}
}
