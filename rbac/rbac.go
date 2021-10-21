/*
SPDX-License-Identifier: Apache-2.0
*/

package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

type ServerConfig struct {
	CCID    string
	Address string
}

// SmartContract provides functions for managing a rbac system
type SmartContract struct {
	contractapi.Contract
}

type User struct {
	UID      string `json:"uid"`
	Active   string `json:"active"`
	UserType string `json:"usertype"`
	PubKey   string `json:"pubkey"`
}

// QueryResult structure used for handling result of query
type UserQueryResult struct {
	Key    string `json:"Key"`
	Record *User
}

type Project struct {
	PID        string `json:"pid"`
	CreatorUID string `json:"creatoruid"`
	CreateDate string `json:"createdate"`
	DueDate    string `json:"duedate"`
	Active     string `json:"active"`
}

// QueryResult structure used for handling result of query
type ProjectQueryResult struct {
	Key    string `json:"Key"`
	Record *Project
}

type Task struct {
	TID           string `json:"tid"`
	PID           string `json:"pid"`
	CreatorUID    string `json:"creatoruid"`
	AssigneeUID   string `json:"assigneeuid"`
	CreateDate    string `json:"createdate"`
	CompletedDate string `json:"completeddate"`
	DueDate       string `json:"duedate"`
	Responsible   string `json:"responsible"`
	Active        string `json:"active"`
}

type Permission struct {
	PermID string `json:"permid"`
	AbleTo string `json:"ableto"`
}

// QueryResult structure used for handling result of query
type PermQueryResult struct {
	Key    string `json:"Key"`
	Record *Permission
}

type Role struct {
	RID      string        `json:"rid"`
	RoleName string        `json:"rolename"`
	Perms    *[]Permission `json:"permissions"`
}

// QueryResult structure used for handling result of query
type RoleQueryResult struct {
	Key    string `json:"Key"`
	Record *Role
}

type RbacModel struct {
	RMID string `json:"rmid"`
	PID  string `json:"pid"`
	Role *Role  `json:"role"`
	UID  string `json:"uid"`
}

type RbacMatrix struct {
	MID  string       `json:"mid"`
	Rbac *[]RbacModel `json:"rbac"`
}

// QueryResult structure used for handling result of query
type MatrixQueryResult struct {
	Key    string `json:"Key"`
	Record *RbacMatrix
}

// QueryResult structure used for handling result of query
type RbacQueryResult struct {
	Key    string `json:"Key"`
	Record *RbacModel
}

// QueryResult structure used for handling result of query
type TaskQueryResult struct {
	Key    string `json:"Key"`
	Record *Task
}

// CreateUser adds a new user to the world state with given details
func (s *SmartContract) CreateUser(ctx contractapi.TransactionContextInterface, UID string, Active string, UserType string, PubKey string) error {
	user := User{
		UID:      UID,
		Active:   Active,
		UserType: UserType,
		PubKey:   PubKey,
	}

	userAsBytes, _ := json.Marshal(user)

	return ctx.GetStub().PutState(UID, userAsBytes)
}

// CreateProject adds a new project to the world state with given details
func (s *SmartContract) CreateProject(ctx contractapi.TransactionContextInterface, PID string, CreatorUID string, CreateDate string, DueDate string, Active string) error {
	project := Project{
		PID:        PID,
		CreatorUID: CreatorUID,
		CreateDate: CreateDate,
		DueDate:    DueDate,
		Active:     Active,
	}

	projectAsBytes, _ := json.Marshal(project)

	return ctx.GetStub().PutState(PID, projectAsBytes)
}

// CreateUser adds a new user to the world state with given details
func (s *SmartContract) CreateTask(ctx contractapi.TransactionContextInterface, TID string, PID string, CreatorUID string, AssigneeUID string, CreateDate string,
	CompletedDate string, DueDate string, Responsible string, Active string) error {
	task := Task{
		TID:           TID,
		PID:           PID,
		CreatorUID:    CreatorUID,
		AssigneeUID:   AssigneeUID,
		CreateDate:    CreateDate,
		CompletedDate: CompletedDate,
		DueDate:       DueDate,
		Responsible:   Responsible,
		Active:        Active,
	}

	taskAsBytes, _ := json.Marshal(task)

	return ctx.GetStub().PutState(TID, taskAsBytes)
}

// CreateUser adds a new user to the world state with given details
func (s *SmartContract) CreatePerm(ctx contractapi.TransactionContextInterface, PermID string, AbleTo string) error {
	permission := Permission{
		PermID: PermID,
		AbleTo: AbleTo,
	}

	permissionAsBytes, _ := json.Marshal(permission)

	return ctx.GetStub().PutState(PermID, permissionAsBytes)
}

// CreateUser adds a new user to the world state with given details
func (s *SmartContract) CreateRole(ctx contractapi.TransactionContextInterface, RID string, RoleName string, Perms *[]Permission) error {
	role := Role{
		RID:      RID,
		RoleName: RoleName,
		Perms:    Perms,
	}

	roleAsBytes, _ := json.Marshal(role)

	return ctx.GetStub().PutState(RID, roleAsBytes)
}

// CreateUser adds a new user to the world state with given details
func (s *SmartContract) CreateRbacModel(ctx contractapi.TransactionContextInterface, RMID string, PID string, Role *Role, UID string) error {
	rbacmodel := RbacModel{
		RMID: RMID,
		PID:  PID,
		Role: Role,
		UID:  UID,
	}

	rbacmodelAsBytes, _ := json.Marshal(rbacmodel)

	return ctx.GetStub().PutState(RMID, rbacmodelAsBytes)
}

// CreateUser adds a new user to the world state with given details
func (s *SmartContract) CreateRbacMatrix(ctx contractapi.TransactionContextInterface, MID string, Rbac *[]RbacModel) error {
	rbacmatrix := RbacMatrix{
		MID:  MID,
		Rbac: Rbac,
	}

	rbacmatrixAsBytes, _ := json.Marshal(rbacmatrix)

	return ctx.GetStub().PutState(MID, rbacmatrixAsBytes)
}

// QueryUser returns the user stored in the world state with given uid
func (s *SmartContract) QueryUser(ctx contractapi.TransactionContextInterface, UID string) (*User, error) {
	userAsBytes, err := ctx.GetStub().GetState(UID)

	if err != nil {
		return nil, fmt.Errorf("Failed to read from world state. %s", err.Error())
	}

	if userAsBytes == nil {
		return nil, fmt.Errorf("%s does not exist", UID)
	}

	user := new(User)
	_ = json.Unmarshal(userAsBytes, user)

	return user, nil
}

// QueryProject returns the project stored in the world state with given pid
func (s *SmartContract) QueryProject(ctx contractapi.TransactionContextInterface, PID string) (*Project, error) {
	projectAsBytes, err := ctx.GetStub().GetState(PID)

	if err != nil {
		return nil, fmt.Errorf("Failed to read from world state. %s", err.Error())
	}

	if projectAsBytes == nil {
		return nil, fmt.Errorf("%s does not exist", PID)
	}

	project := new(Project)
	_ = json.Unmarshal(projectAsBytes, project)

	return project, nil
}

// QueryUser returns the user stored in the world state with given uid
func (s *SmartContract) QueryTask(ctx contractapi.TransactionContextInterface, TID string) (*Task, error) {
	userAsBytes, err := ctx.GetStub().GetState(TID)

	if err != nil {
		return nil, fmt.Errorf("Failed to read from world state. %s", err.Error())
	}

	if userAsBytes == nil {
		return nil, fmt.Errorf("%s does not exist", TID)
	}

	task := new(Task)
	_ = json.Unmarshal(userAsBytes, task)

	return task, nil
}

// QueryUser returns the user stored in the world state with given uid
func (s *SmartContract) QueryPerm(ctx contractapi.TransactionContextInterface, PermID string) (*Permission, error) {
	permissionAsBytes, err := ctx.GetStub().GetState(PermID)

	if err != nil {
		return nil, fmt.Errorf("Failed to read from world state. %s", err.Error())
	}

	if permissionAsBytes == nil {
		return nil, fmt.Errorf("%s does not exist", PermID)
	}

	permission := new(Permission)
	_ = json.Unmarshal(permissionAsBytes, permission)

	return permission, nil
}

// QueryUser returns the user stored in the world state with given uid
func (s *SmartContract) QueryRole(ctx contractapi.TransactionContextInterface, RID string) (*Role, error) {
	roleAsBytes, err := ctx.GetStub().GetState(RID)

	if err != nil {
		return nil, fmt.Errorf("Failed to read from world state. %s", err.Error())
	}

	if roleAsBytes == nil {
		return nil, fmt.Errorf("%s does not exist", RID)
	}

	role := new(Role)
	_ = json.Unmarshal(roleAsBytes, role)

	return role, nil
}

// QueryUser returns the user stored in the world state with given uid
func (s *SmartContract) QueryRbacModel(ctx contractapi.TransactionContextInterface, RMID string) (*RbacModel, error) {
	rbacmodelAsBytes, err := ctx.GetStub().GetState(RMID)

	if err != nil {
		return nil, fmt.Errorf("Failed to read from world state. %s", err.Error())
	}

	if rbacmodelAsBytes == nil {
		return nil, fmt.Errorf("%s does not exist", RMID)
	}

	rbacmodel := new(RbacModel)
	_ = json.Unmarshal(rbacmodelAsBytes, rbacmodel)

	return rbacmodel, nil
}

// QueryUser returns the user stored in the world state with given uid
func (s *SmartContract) QueryRbacMatrix(ctx contractapi.TransactionContextInterface, MID string) (*RbacMatrix, error) {
	rbacmatrixAsBytes, err := ctx.GetStub().GetState(MID)

	if err != nil {
		return nil, fmt.Errorf("Failed to read from world state. %s", err.Error())
	}

	if rbacmatrixAsBytes == nil {
		return nil, fmt.Errorf("%s does not exist", MID)
	}

	rbacmatrix := new(RbacMatrix)
	_ = json.Unmarshal(rbacmatrixAsBytes, rbacmatrix)

	return rbacmatrix, nil
}

// QueryAllUsers returns all users found in world state
func (s *SmartContract) QueryAllUsers(ctx contractapi.TransactionContextInterface) ([]UserQueryResult, error) {
	startKey := ""
	endKey := ""

	resultsIterator, err := ctx.GetStub().GetStateByRange(startKey, endKey)

	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	results := []UserQueryResult{}

	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()

		if err != nil {
			return nil, err
		}

		user := new(User)
		_ = json.Unmarshal(queryResponse.Value, user)

		queryResult := UserQueryResult{Key: queryResponse.Key, Record: user}
		results = append(results, queryResult)
	}

	return results, nil
}

// QueryAllProjects returns all projects found in world state
func (s *SmartContract) QueryAllProjects(ctx contractapi.TransactionContextInterface) ([]ProjectQueryResult, error) {
	startKey := ""
	endKey := ""

	resultsIterator, err := ctx.GetStub().GetStateByRange(startKey, endKey)

	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	results := []ProjectQueryResult{}

	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()

		if err != nil {
			return nil, err
		}

		project := new(Project)
		_ = json.Unmarshal(queryResponse.Value, project)

		queryResult := ProjectQueryResult{Key: queryResponse.Key, Record: project}
		results = append(results, queryResult)
	}

	return results, nil
}

// QueryAllUsers returns all users found in world state
func (s *SmartContract) QueryAllTasks(ctx contractapi.TransactionContextInterface) ([]TaskQueryResult, error) {
	startKey := ""
	endKey := ""

	resultsIterator, err := ctx.GetStub().GetStateByRange(startKey, endKey)

	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	results := []TaskQueryResult{}

	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()

		if err != nil {
			return nil, err
		}

		task := new(Task)
		_ = json.Unmarshal(queryResponse.Value, task)

		queryResult := TaskQueryResult{Key: queryResponse.Key, Record: task}
		results = append(results, queryResult)
	}

	return results, nil
}

// QueryAllUsers returns all users found in world state
func (s *SmartContract) QueryAllPerms(ctx contractapi.TransactionContextInterface) ([]PermQueryResult, error) {
	startKey := ""
	endKey := ""

	resultsIterator, err := ctx.GetStub().GetStateByRange(startKey, endKey)

	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	results := []PermQueryResult{}

	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()

		if err != nil {
			return nil, err
		}

		permission := new(Permission)
		_ = json.Unmarshal(queryResponse.Value, permission)

		queryResult := PermQueryResult{Key: queryResponse.Key, Record: permission}
		results = append(results, queryResult)
	}

	return results, nil
}

// QueryAllUsers returns all users found in world state
func (s *SmartContract) QueryAllRoles(ctx contractapi.TransactionContextInterface) ([]RoleQueryResult, error) {
	startKey := ""
	endKey := ""

	resultsIterator, err := ctx.GetStub().GetStateByRange(startKey, endKey)

	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	results := []RoleQueryResult{}

	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()

		if err != nil {
			return nil, err
		}

		role := new(Role)
		_ = json.Unmarshal(queryResponse.Value, role)

		queryResult := RoleQueryResult{Key: queryResponse.Key, Record: role}
		results = append(results, queryResult)
	}

	return results, nil
}

// QueryAllUsers returns all users found in world state
func (s *SmartContract) QueryAllRbacModels(ctx contractapi.TransactionContextInterface) ([]RbacQueryResult, error) {
	startKey := ""
	endKey := ""

	resultsIterator, err := ctx.GetStub().GetStateByRange(startKey, endKey)

	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	results := []RbacQueryResult{}

	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()

		if err != nil {
			return nil, err
		}

		rbacmodel := new(RbacModel)
		_ = json.Unmarshal(queryResponse.Value, rbacmodel)

		queryResult := RbacQueryResult{Key: queryResponse.Key, Record: rbacmodel}
		results = append(results, queryResult)
	}

	return results, nil
}

// QueryAllUsers returns all users found in world state
func (s *SmartContract) QueryAllRbacMatrices(ctx contractapi.TransactionContextInterface) ([]MatrixQueryResult, error) {
	startKey := ""
	endKey := ""

	resultsIterator, err := ctx.GetStub().GetStateByRange(startKey, endKey)

	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	results := []MatrixQueryResult{}

	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()

		if err != nil {
			return nil, err
		}

		rbacmatrix := new(RbacMatrix)
		_ = json.Unmarshal(queryResponse.Value, rbacmatrix)

		queryResult := MatrixQueryResult{Key: queryResponse.Key, Record: rbacmatrix}
		results = append(results, queryResult)
	}

	return results, nil
}

// ChangeUserActiveStatus updates the active field of an user with given uid in world state
func (s *SmartContract) ChangeUserActiveStatus(ctx contractapi.TransactionContextInterface, UID string, Active string) error {
	user, err := s.QueryUser(ctx, UID)

	if err != nil {
		return err
	}

	user.Active = Active

	userAsBytes, _ := json.Marshal(user)

	return ctx.GetStub().PutState(UID, userAsBytes)
}

// ChangeProjectActiveStatus updates the active field of a project with given pid in world state
func (s *SmartContract) ChangeProjectActiveStatus(ctx contractapi.TransactionContextInterface, PID string, Active string) error {
	project, err := s.QueryProject(ctx, PID)

	if err != nil {
		return err
	}

	project.Active = Active

	projectAsBytes, _ := json.Marshal(project)

	return ctx.GetStub().PutState(PID, projectAsBytes)
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
