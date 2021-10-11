package main

import (
	"encoding/json"
	"fmt"
	"time"

	"errors"

	"github.com/rs/xid"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// SmartContract provides functions for managing an Party
type SmartContract struct {
	contractapi.Contract
}

type Project struct {
	ProjectID   string `json:"projectid"`
	ProjectName string `json:"projectname"`
	CreatorID   string `json:"creatorid"`
	CreateDate  string `json:"createdate"`
	DueDate     string `json:"duedate"`
}

// Party describes basic details of what makes up a simple Party
type Role struct {
	RoleID   string `json:"roleid"`
	RoleName string `json:"rolename"`
}

type Permission struct {
	PermissionID   string `json:"permid"`
	PermissionName string `json:"permname"`
}

type UserRole struct {
	ID     string `json:"urid"`
	UserID string `json:"userid"`
	RoleID string `json:"roleid"`
}

type RolePermission struct {
	RPID   string       `json:"rpid"`
	RoleID string       `json:"roleid"`
	Perms  []Permission `json:"perms"`
}

var (
	ErrPermissionInUse     = errors.New("cannot delete assigned permission")
	ErrPermissionNotFound  = errors.New("permission not found")
	ErrRoleAlreadyAssigned = errors.New("this role is already assigned to the user")
	ErrRoleInUse           = errors.New("cannot delete assigned role")
	ErrRoleNotFound        = errors.New("role not found")
)

//RoleExists returns true when Party with given ID exists in world state
func (s *SmartContract) ProjectExists(ctx contractapi.TransactionContextInterface, ProjectName string) (bool, error) {
	ProjectJSON, err := ctx.GetStub().GetState(ProjectName)
	if err != nil {
		return false, fmt.Errorf("failed to read from world state: %v", err)
	}

	return ProjectJSON != nil, nil
}

//RoleExists returns true when Party with given ID exists in world state
func (s *SmartContract) RoleExists(ctx contractapi.TransactionContextInterface, RoleName string) (bool, error) {
	RoleJSON, err := ctx.GetStub().GetState(RoleName)
	if err != nil {
		return false, fmt.Errorf("failed to read from world state: %v", err)
	}

	return RoleJSON != nil, nil
}

//RoleExists returns true when Party with given ID exists in world state
func (s *SmartContract) PermissionExists(ctx contractapi.TransactionContextInterface, PermissionName string) (bool, error) {
	PermJSON, err := ctx.GetStub().GetState(PermissionName)
	if err != nil {
		return false, fmt.Errorf("failed to read from world state: %v", err)
	}

	return PermJSON != nil, nil
}

//RoleExists returns true when Party with given ID exists in world state
func (s *SmartContract) RolePermExists(ctx contractapi.TransactionContextInterface, RoleID string, PermID string) (bool, error) {
	RPJSON, err := ctx.GetStub().GetState(RoleID)
	if err != nil {
		return false, fmt.Errorf("failed to read from world state: %v", err)
	}

	return RPJSON != nil, nil
}

//CreateRole issues a new Party to the world state with given details.
func (s *SmartContract) CreateProject(ctx contractapi.TransactionContextInterface, ProjectName string, CreatorID string) error {
	exists, err := s.ProjectExists(ctx, ProjectName)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("the Role %s already exists", ProjectName)
	}

	Project := Project{
		ProjectName: ProjectName,
		ProjectID:   xid.New().String(),
		CreatorID:   CreatorID,
		CreateDate:  time.Now().Format("2006.01.02 15:04:05"),
		DueDate:     time.Now().Format("2006.01.02 15:04:05"),
	}
	ProjectJSON, err := json.Marshal(Project)
	if err != nil {
		return err
	}

	s.CreateRole(ctx, ProjectName+"_user")
	s.CreateRole(ctx, ProjectName+"_PM")

	return ctx.GetStub().PutState(ProjectName, ProjectJSON)
}

//CreateRole issues a new Party to the world state with given details.
func (s *SmartContract) CreateRole(ctx contractapi.TransactionContextInterface, RoleName string) error {
	exists, err := s.RoleExists(ctx, RoleName)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("the Role %s already exists", RoleName)
	}

	Role := Role{
		RoleName: RoleName,
		RoleID:   xid.New().String(),
	}
	RoleJSON, err := json.Marshal(Role)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(RoleName, RoleJSON)
}

//CreateRole issues a new Party to the world state with given details.
func (s *SmartContract) CreatePerm(ctx contractapi.TransactionContextInterface, PermName string) error {
	exists, err := s.RoleExists(ctx, PermName)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("the Permission %s already exists", PermName)
	}

	Permission := Permission{
		PermissionName: PermName,
		PermissionID:   xid.New().String(),
	}
	PermJSON, err := json.Marshal(Permission)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(PermName, PermJSON)
}

//CreateRole issues a new Party to the world state with given details.
func (s *SmartContract) CreateRolePerm(ctx contractapi.TransactionContextInterface, RoleID string, Perms []Permission) error {
	RolePermission := RolePermission{
		RoleID: RoleID,
		Perms:  Perms,
		RPID:   xid.New().String(),
	}
	RPJSON, err := json.Marshal(RolePermission)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(RolePermission.RPID, RPJSON)
}

//ReadRole returns the Role stored in the world state with given id.
func (s *SmartContract) ReadRole(ctx contractapi.TransactionContextInterface, RoleName string) (*Role, error) {
	RoleJSON, err := ctx.GetStub().GetState(RoleName)
	if err != nil {
		return nil, fmt.Errorf("failed to read from world state: %v", err)
	}
	if RoleJSON == nil {
		return nil, fmt.Errorf("the Party %s does not exist", RoleName)
	}

	var Role Role
	err = json.Unmarshal(RoleJSON, &Role)
	if err != nil {
		return nil, err
	}

	return &Role, nil
}

func (s *SmartContract) ReadPerm(ctx contractapi.TransactionContextInterface, PermName string) (*Permission, error) {
	PermJSON, err := ctx.GetStub().GetState(PermName)
	if err != nil {
		return nil, fmt.Errorf("failed to read from world state: %v", err)
	}
	if PermJSON == nil {
		return nil, fmt.Errorf("the Party %s does not exist", PermName)
	}

	var Permission Permission
	err = json.Unmarshal(PermJSON, &Permission)
	if err != nil {
		return nil, err
	}

	return &Permission, nil
}

func (s *SmartContract) ReadRolePerm(ctx contractapi.TransactionContextInterface, RPID string, RoleID string, PermID string) (*RolePermission, error) {
	RPJSON, err := ctx.GetStub().GetState(RPID)
	if err != nil {
		return nil, fmt.Errorf("failed to read from world state: %v", err)
	}
	if RPJSON == nil {
		return nil, fmt.Errorf("the Party %s does not exist", RPID)
	}

	var RolePermission RolePermission
	err = json.Unmarshal(RPJSON, &RolePermission)
	if err != nil {
		return nil, err
	}

	return &RolePermission, nil
}

func (s *SmartContract) AssignPermissions(ctx contractapi.TransactionContextInterface, RoleName string, Perms []Permission) error {

	role, err := s.ReadRole(ctx, RoleName)

	if err != nil {
		return ErrRoleNotFound
	}

	// var perms []Permission

	// for _, permName := range PermNames {

	// 	perm, err1 := s.ReadPerm(ctx, permName)
	// 	if err1 != nil {
	// 		return ErrPermissionNotFound
	// 	}
	// 	perms = append(perms, *perm)
	// }

	// for _, perm := range perms {
	// 	s.CreateRolePerm(ctx, role.RoleID, perm.PermissionID)
	// }
	s.CreateRolePerm(ctx, role.RoleID, Perms)
	return nil
}

//GetAllRoles returns all Roles found in world state
func (s *SmartContract) GetAllRoles(ctx contractapi.TransactionContextInterface) ([]*Role, error) {
	//range query with empty string for startKey and endKey does an
	//open-ended query of all Partys in the chaincode namespace.
	resultsIterator, err := ctx.GetStub().GetStateByRange("", "")
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	var Roles []*Role
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		var Role Role
		err = json.Unmarshal(queryResponse.Value, &Role)
		if err != nil {
			return nil, err
		}
		Roles = append(Roles, &Role)
	}

	return Roles, nil
}

//GetAllRoles returns all Roles found in world state
func (s *SmartContract) GetAllProjects(ctx contractapi.TransactionContextInterface) ([]*Project, error) {
	//range query with empty string for startKey and endKey does an
	//open-ended query of all Partys in the chaincode namespace.
	resultsIterator, err := ctx.GetStub().GetStateByRange("", "")
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	var Projects []*Project
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		var Project Project
		err = json.Unmarshal(queryResponse.Value, &Project)
		if err != nil {
			return nil, err
		}
		Projects = append(Projects, &Project)
	}

	return Projects, nil
}

//GetAllRoles returns all Roles found in world state
func (s *SmartContract) GetAllPerms(ctx contractapi.TransactionContextInterface) ([]*Permission, error) {
	//range query with empty string for startKey and endKey does an
	//open-ended query of all Partys in the chaincode namespace.
	resultsIterator, err := ctx.GetStub().GetStateByRange("", "")
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	var Perms []*Permission
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		var Perm Permission
		err = json.Unmarshal(queryResponse.Value, &Perm)
		if err != nil {
			return nil, err
		}
		Perms = append(Perms, &Perm)
	}

	return Perms, nil
}

//GetAllRoles returns all Roles found in world state
func (s *SmartContract) GetAllRolePerms(ctx contractapi.TransactionContextInterface) ([]*RolePermission, error) {
	//range query with empty string for startKey and endKey does an
	//open-ended query of all Partys in the chaincode namespace.
	resultsIterator, err := ctx.GetStub().GetStateByRange("", "")
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	var RolePerms []*RolePermission
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		var RolePerm RolePermission
		err = json.Unmarshal(queryResponse.Value, &RolePerm)
		if err != nil {
			return nil, err
		}
		RolePerms = append(RolePerms, &RolePerm)
	}

	return RolePerms, nil
}

// InitLedger adds a base set of Partys to the ledger
// func (s *SmartContract) InitLedger(ctx contractapi.TransactionContextInterface) error {
// 	Partys := []Party{
// 		{PartyName: "Republicans", VoteCount: 0},
// 		{PartyName: "Democrats", VoteCount: 0},
// 	}

// 	for _, Party := range Partys {
// 		PartyJSON, err := json.Marshal(Party)
// 		if err != nil {
// 			return err
// 		}

// 		err = ctx.GetStub().PutState(Party.PartyName, PartyJSON)
// 		if err != nil {
// 			return fmt.Errorf("failed to put to world state. %v", err)
// 		}
// 	}

// 	return nil
// }

// UpdateParty updates an existing Party in the world state with provided parameters.
// func (s *SmartContract) UpdateParty(ctx contractapi.TransactionContextInterface, PartyName string, VoteCount int) error {
// 	exists, err := s.PartyExists(ctx, PartyName)
// 	if err != nil {
// 		return err
// 	}
// 	if !exists {
// 		return fmt.Errorf("the Party %s does not exist", PartyName)
// 	}

// 	overwriting original Party with new Party
// 	Party := Party{
// 		PartyName:      PartyName,
// 		VoteCount:      VoteCount,
// 	}
// 	PartyJSON, err := json.Marshal(Party)
// 	if err != nil {
// 		return err
// 	}

// 	return ctx.GetStub().PutState(PartyName, PartyJSON)
// }

// DeleteParty deletes an given Party from the world state.
// func (s *SmartContract) DeleteParty(ctx contractapi.TransactionContextInterface, PartyName string) error {
// 	exists, err := s.PartyExists(ctx, PartyName)
// 	if err != nil {
// 		return err
// 	}
// 	if !exists {
// 		return fmt.Errorf("the Party %s does not exist", PartyName)
// 	}

// 	return ctx.GetStub().DelState(PartyName)
// }

// TransferParty updates the owner field of Party with given id in world state.
// func (s *SmartContract) TransferParty(ctx contractapi.TransactionContextInterface, PartyName string, VoteCount int) error {
// 	Party, err := s.ReadParty(ctx, PartyName)
// 	if err != nil {
// 		return err
// 	}

// 	Party.VoteCount = VoteCount
// 	PartyJSON, err := json.Marshal(Party)
// 	if err != nil {
// 		return err
// 	}

// 	return ctx.GetStub().PutState(PartyName, PartyJSON)
// }

// CastVote updates the owner field of Party with given id in world state.
// func (s *SmartContract) CastVote(ctx contractapi.TransactionContextInterface, PartyName string) error {
// 	Party, err := s.ReadParty(ctx, PartyName)
// 	if err != nil {
// 		return err
// 	}

// 	Party.VoteCount = Party.VoteCount + 1
// 	PartyJSON, err := json.Marshal(Party)
// 	if err != nil {
// 		return err
// 	}

// 	return ctx.GetStub().PutState(PartyName, PartyJSON)
// }

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
