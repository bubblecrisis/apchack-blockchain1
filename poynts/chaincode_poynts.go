package main

import (
	"errors"
	"fmt"
	"strconv"
	"github.com/hyperledger/fabric/core/chaincode/shim"
)

// SimpleChaincode example simple Chaincode implementation
type SimpleChaincode struct {
}

// ============================================================================================================================
// Main
// ============================================================================================================================
func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}

// Init resets all the things
func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	for i := 1; i < len(args); i ++ {
		stub.PutState(args[i], []byte(strconv.Itoa(0)))
		fmt.Printf("Added %d\n", args[i])
	}
	return nil, nil
}

// Invoke is our entry point to invoke a chaincode function
func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("invoke is running " + function)

	// Handle different functions
	if function == "init" {													//initialize the chaincode state, used as reset
		return t.Init(stub, "init", args)
	} else if function == "add" {
		return t.add(stub, args)
	} else if function == "deduct" {
		return t.deduct(stub, args)
	} else {
		fmt.Println("invoke did not find func: " + function)					//error
		return nil, errors.New("Received unknown function invocation: " + function)
	}
}

// Query is our entry point for queries
func (t *SimpleChaincode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("query is running " + function)

	// Handle different functions
	if function == "balance" {											//read a variable
		return t.balance(stub, args)						
	}
	fmt.Println("query did not find func: " + function)						//error

	return nil, errors.New("Received unknown function query: " + function)
}

//================================================================================================
func (t *SimpleChaincode) balance(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {	
	var err error
	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting 1")
	}

	balanceBytes, err := stub.GetState(args[0])
	if (err != nil) {
		return nil, errors.New("Error retrieving balance")
	} else {
		return balanceBytes, nil
	}
}

// args[0] - bucket
// args[1] - value to add
func (t *SimpleChaincode) add(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var value, delta int 
	if len(args) != 2 {
		return nil, errors.New("Incorrect number of arguments. Expecting 2")
	}
	valueBytes, err := stub.GetState(args[0])
	if err != nil {
		return nil, errors.New("Failed to get state: " + args[0])
	}
	if valueBytes == nil {
		return nil, errors.New("Entity not found " + args[0])
	}
	value, _ = strconv.Atoi(string(valueBytes))
	delta, err = strconv.Atoi(args[1])
	value = value + delta
	err = stub.PutState(args[0], []byte(strconv.Itoa(value)))
	if err != nil {
		return nil, err
	}
	return []byte(strconv.Itoa(value)), nil
}

// args[0] - bucket
// args[1] - value to add
func (t *SimpleChaincode) deduct(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var value, delta int 
	if len(args) != 2 {
		return nil, errors.New("Incorrect number of arguments. Expecting 2")
	}
	valueBytes, err := stub.GetState(args[0])
	if err != nil {
		return nil, errors.New("Failed to get state: " + args[0])
	}
	if valueBytes == nil {
		return nil, errors.New("Entity not found " + args[0])
	}
	value, _ = strconv.Atoi(string(valueBytes))
	delta, err = strconv.Atoi(args[1])
	value = value - delta
	err = stub.PutState(args[0], []byte(strconv.Itoa(value)))
	if err != nil {
		return nil, err
	}
	return []byte(strconv.Itoa(value)), nil
}