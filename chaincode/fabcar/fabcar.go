
// This code is for giving all the instances where temparture violated specific
// threshold, the query call be seen from queryTemp.js

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
	//"bytes"
	"encoding/json"
	"fmt"
	"strconv"
	"time"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	sc "github.com/hyperledger/fabric/protos/peer"
)

var OpenTradeKey = "_opentrades"

// Define the Smart Contract structure
type SmartContract struct {
}

// Define the Med structure, with 4 properties.  Structure tags are used by encoding/json library
type Med struct {
	ID 			string 		`json:"id"`
	Description string 		`json:"description"`
	Owner  		string 		`json:"owner"`
	Location 	string 		`json:"location"`
	Timestamp 	time.Time 	`json:"timestamp"`
	Temperature float64 	`json:"temperature"`
	Status 		bool 		`json:"status"`
	Lower_Limit float64 	`json:"lower_limit"`
	Upper_Limit float64		`json:"upper_limit"`
}

type OpenTrade struct {
	Owner 	  string   	`json:"owner"`
	Timestamp time.Time `json:"timestamp"`   
	MedList   []string 	`json:medlist"`
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
	} else if function == "transferOwnership" {
		return s.transferOwnership(APIstub, args)
	} else if function == "changeLocTemp" {
		return s.changeLocTemp(APIstub, args)
	} else if function == "queryMedHistory" {
		return s.queryMedHistory(APIstub, args)
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

func (s *SmartContract) queryMedHistory(APIstub shim.ChaincodeStubInterface, args[]string) sc.Response {
	
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}


	versionsIterator, err := APIstub.GetHistoryForKey(args[0])

	if err != nil {
		return shim.Error(err.Error())
	}

	defer versionsIterator.Close()

	MedHistoryList := []Med{}

	for versionsIterator.HasNext() {
		queryResponse, err := versionsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}
		med := Med{}
		json.Unmarshal(queryResponse.Value, &med)
		med.ID = args[0]
		MedHistoryList = append(MedHistoryList, med)
	}

	fmt.Printf("- queryMedHistory:\n%+v\n", MedHistoryList)
	medAsBytes, _ := json.Marshal(MedHistoryList)
	return shim.Success(medAsBytes)
 }

func (s *SmartContract) initLedger(APIstub shim.ChaincodeStubInterface) sc.Response {
	meds := []Med{
		Med{ID: "4501", Description: "", Owner: "Tomoko", Location: "NA", Timestamp: time.Now(), Temperature: 23, Status : true, Lower_Limit : 20, Upper_Limit: 25},
		Med{ID: "4570", Description: "", Owner: "Thomas", Location: "NA", Timestamp: time.Now(), Temperature: 23, Status : true, Lower_Limit: 18, Upper_Limit: 25},
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
	if len(args) != 7 {
		return shim.Error("Incorrect number of arguments. Expecting 9")
	}

	var med = getUpdatedObject(args[1:])
	med.Status = validateTemperature(med.Lower_Limit, med.Upper_Limit, med.Temperature)
	medAsBytes, _ := json.Marshal(med)
	APIstub.PutState(args[0], medAsBytes)
	if (!med.Status) {
		return shim.Error("Temperature Conditions have been violated")
	}
	return shim.Success(nil)
}

func getUpdatedObject (objArgs []string) Med {
	if n, err := strconv.ParseFloat(objArgs[4], 64); err == nil {
		LL, _ := strconv.ParseFloat(objArgs[5], 64)
		UL, _ := strconv.ParseFloat(objArgs[6], 64)
        med := Med{
			ID : objArgs[0],
			Description : objArgs[1],
			Owner : objArgs[2],
			Location: objArgs[3],
			Timestamp: time.Now(),
			Temperature: n,
			Lower_Limit: LL,
			Upper_Limit: UL,
		}
		return med
    }
    med := Med{}
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
	AllMedList := []Med{}
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}
		med := Med{}
		json.Unmarshal(queryResponse.Value, &med)
		med.ID = queryResponse.Key
		AllMedList = append(AllMedList, med)
	}

	fmt.Printf("- queryAllMeds:\n%+v\n", AllMedList)
	medAsBytes, err := json.Marshal(AllMedList)
	return shim.Success(medAsBytes)
}

func (s *SmartContract) transferOwnership(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

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

func validateTemperature(lower_limit float64, upper_limit float64, given float64) bool {
    if given >= lower_limit && given <= upper_limit {
        return true
    }
    return false
}

func (s *SmartContract) changeLocTemp(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	
	if len (args) != 3 {
		return shim.Error("Incorrect number of arguments. Expecting 3")
	}

    if n, err := strconv.ParseFloat(args[2], 64) ; err == nil{
		medAsBytes, _ := APIstub.GetState(args[0])
		med := Med{}
		json.Unmarshal(medAsBytes, &med)
		if(!validateTemperature(med.Lower_Limit, med.Upper_Limit, n)){
			return shim.Error("Temparture conditions have been violated")
		}
		med.Location = args[1]
		med.Temperature = n
		medAsBytes, _ = json.Marshal(med)
		APIstub.PutState(args[0], medAsBytes)
		return shim.Success(nil)
    }
	return shim.Success(nil)
}



func (s *SmartContract) open_trade(APIstub shim.ChaincodeStubInterface, args []string) {

	//if len(args) != 
}


/*func (s *SmartContract) tempartureViolations(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	fmt.Printf("- tempartureViolations:")
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
		medAsBytes := queryResponse.Value
		med := Med{}
		json.Unmarshal(medAsBytes, &med)
		
		var lower_limit = 18.0
	    var upper_limit = 26.0

		if(!validateTemperature(lower_limit, upper_limit, med.Temperature)){
			return shim.Error("Temparture conditions have been violated somewhere")
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

	fmt.Printf("- tempartureViolations:\n%s\n", buffer.String())

	return shim.Success(buffer.Bytes())
}
*/


// The main function is only relevant in unit test mode. Only included here for completeness.
func main() {

	// Create a new Smart Contract
	err := shim.Start(new(SmartContract))
	if err != nil {
		fmt.Printf("Error creating new Smart Contract: %s", err)
	}
}



