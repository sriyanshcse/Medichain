/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

/*
 * The sample smart contract for documentation topic:
 * Writing Your First Blockchain Application
 */

package main

/* Imports
 * 4 utility libraries for formatting, handling bytes, reading and writing JSON, and string manipulation
 * 2 specific Hyperledger Fabric specific libraries for Smart Contracts
 */
import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
	"time"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	sc "github.com/hyperledger/fabric/protos/peer"
)

// Define the Smart Contract structure
type SmartContract struct {
}

// Define the Med structure, with 4 properties.  Structure tags are used by encoding/json library
type Med struct {
	Make   string `json:"make"`
	Name  string `json:"name"`
	BatchNo string `json:"batchno"`
	ExpiryDate string `json:"expirydate"`
	Owner  string `json:"owner"`
	Location string `json:"location"`
	Timestamp time.Time `json:"timestamp"`
	Temperature string `json:"temperature"`
}

/*
 * The Init method is called when the Smart Contract "fabMed" is instantiated by the blockchain network
 * Best practice is to have any Ledger initialization in separate function -- see initLedger()
 */
func (s *SmartContract) Init(APIstub shim.ChaincodeStubInterface) sc.Response {
	return shim.Success(nil)
}

/*
 * The Invoke method is called as a result of an application request to run the Smart Contract "fabMed"
 * The calling application program has also specified the particular smart contract function to be called, with arguments
 */
func (s *SmartContract) Invoke(APIstub shim.ChaincodeStubInterface) sc.Response {

	// Retrieve the requested Smart Contract function and arguments
	function, args := APIstub.GetFunctionAndParameters()
	// Route to the appropriate handler function to interact with the ledger appropriately
	if function == "queryMed" {
		return s.queryMed(APIstub, args)
	} else if function == "initLedger" {
		return s.initLedger(APIstub)
	} else if function == "createMed" {
		return s.createMed(APIstub, args)
	} else if function == "queryAllMeds" {
		return s.queryAllMeds(APIstub)
	} else if function == "changeMedOwner" {
		return s.changeMedOwner(APIstub, args)
	}

	return shim.Error("Invalid Smart Contract function name.")
}

func (s *SmartContract) queryMed(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	medAsBytes, _ := APIstub.GetState(args[0])
	return shim.Success(medAsBytes)
}

func (s *SmartContract) initLedger(APIstub shim.ChaincodeStubInterface) sc.Response {
	meds := []Med{
		Med{Make: "Cipla", Name: "Azid", BatchNo: "123", ExpiryDate: "02/2020",Owner: "Tomoko", Location: "NA", Timestamp: time.Now(), Temperature: "23"},
		Med{Make: "GSK", Name: "Crocin", BatchNo: "101", ExpiryDate: "03/2020",Owner: "Thomas", Location: "NA", Timestamp: time.Now(), Temperature: "23"},
	}

	i := 0
	for i < len(meds) {
		fmt.Println("i is ", i)
		medAsBytes, _ := json.Marshal(meds[i])
		APIstub.PutState("Med"+strconv.Itoa(i), medAsBytes)
		fmt.Println("Added", meds[i])
		i = i + 1
	}

	return shim.Success(nil)
}

func (s *SmartContract) createMed(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	fmt.Println(len(args))
	if len(args) != 8 {
		return shim.Error("Incorrect number of arguments. Expecting 9")
	}

	var med = getUpdatedObject(args[1:])

	medAsBytes, _ := json.Marshal(med)
	APIstub.PutState(args[0], medAsBytes)

	return shim.Success(nil)
}

func getUpdatedObject (objArgs []string) Med {
	med := Med{
		Make: objArgs[0],
		Name : objArgs[1],
		BatchNo: objArgs[2],
		ExpiryDate: objArgs[3],
		Owner: objArgs[4],
		Location: objArgs[5],
		Timestamp: time.Now(),
		Temperature: objArgs[6],
	}
	return med
}

func (s *SmartContract) queryAllMeds(APIstub shim.ChaincodeStubInterface) sc.Response {

	startKey := "Med0"
	endKey := "Med999"

	resultsIterator, err := APIstub.GetStateByRange(startKey, endKey)
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	// buffer is a JSON array containing QueryResults
	var buffer bytes.Buffer
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}
		// Add a comma before array members, suppress it for the first array member
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}
		buffer.WriteString("{\"Key\":")
		buffer.WriteString("\"")
		buffer.WriteString(queryResponse.Key)
		buffer.WriteString("\"")

		buffer.WriteString(", \"Record\":")
		// Record is a JSON object, so we write as-is
		buffer.WriteString(string(queryResponse.Value))
		buffer.WriteString("}")
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")

	fmt.Printf("- queryAllMeds:\n%s\n", buffer.String())

	return shim.Success(buffer.Bytes())
}

func (s *SmartContract) changeMedOwner(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}

	medAsBytes, _ := APIstub.GetState(args[0])
	med := Med{}

	json.Unmarshal(medAsBytes, &med)
	med.Owner = args[1]

	medAsBytes, _ = json.Marshal(med)
	APIstub.PutState(args[0], medAsBytes)

	return shim.Success(nil)
}

func (s *SmartContract) changeLocTemp(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	if len (args) != 3 {
		return shim.Error("Incorrect number of arguments. Expecting 3")
	}

	medAsBytes, _ := APIstub.GetState(args[0])
	med := Med{}

	json.Unmarshal(medAsBytes, &med)
	med.Location = args[1]
	med.Temperature = args[2]
	medAsBytes, _ = json.Marshal(med)
	APIstub.PutState(args[0], medAsBytes)
	return shim.Success(nil)
}

// The main function is only relevant in unit test mode. Only included here for completeness.
func main() {

	// Create a new Smart Contract
	err := shim.Start(new(SmartContract))
	if err != nil {
		fmt.Printf("Error creating new Smart Contract: %s", err)
	}
}
