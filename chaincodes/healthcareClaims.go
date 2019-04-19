package main

import (
    "bytes"
    "encoding/json"
    "fmt"
    "strconv"
    "time"
    "strings"

    "github.com/hyperledger/fabric/core/chaincode/shim"
    pb "github.com/hyperledger/fabric/protos/peer"
)

type ActionType int

const (
    batch              ActionType = 0
    readyToInsurance   ActionType = 1
    readyToThirdParty  ActionType = 2
    paid               ActionType = 3
    other              ActionType = 4
)

// Empty struct that is used for future chaincode functions/methods.
type healthCareChainCode struct {
}

type StatementBatch struct {
    BatchID     string     `json:BatchID`     // The Unique ID that will be used for the batch.
    StatementID string     `json:StatementID` // The ID for the statement that this batch contains.
    Description string     `json:description` // description of batch.
    Date        string     `json:Date`        // Date and time of Action.
    Action      ActionType `json:action`      //action for this transaction
}


// Struct for Statement information.
type Statement struct {
    Id               string     `json:Id`                  // the ID that will be used for the statement.
    Provider         string     `json:Provider`            // name of provider.
    NetworkInOut     string     `json:NetworkInOut`        // status of network.
    PatientName      string     `json:PatientName`         // name of patient.
    DiagnosisID      []string   `json:DiagnosisID`         // the diagnosisID that will be used for this statement.
    Description      string     `json:description`         // description of this statement
    CurrentPhase     string     `json:currentPhase`        // current status of the statement
    AttachementLink  string     `json:attachementLink`
    Action           ActionType `json:action`              // action for this transaction
    Date             string     `json:Date`                // data of this action happened
}

type Diagnosis struct {
    Id                   string   `json:Id`                      // the ID of diagnosis
    CPTCode              string   `json:CPTCode`                 // the Code of CPT
    CPTName              string   `json:CPTName`                 // the name of CPT
    ICD9                 string   `json:ICD9`                    // the amount of ICD9
    Description          string   `json:Description`             // the descriptio of the diagnosis
    TotalCharge          string   `json:TotalCharge`             // the total amount of the diagnosis
    PatientAmount        string   `json:PatientAmount`           // the amount of patient should pay
    BilledToInsurance    string   `json:BilledToInsurance`       // the amount of insurance 
}

// =============================================================================
// The main function. We start the needed method from here.
// =============================================================================
func main() {
    err := shim.Start(new(healthCareChainCode))
    if err != nil {
        fmt.Printf("Error starting Parts Trace chaincode: %s", err)
    }
}

// =============================================================================
// Init initializes chaincode
// =============================================================================
func (t *healthCareChainCode) Init(stub shim.ChaincodeStubInterface) pb.Response {
    return shim.Success(nil)
}

// =============================================================================
// Invoke - Our entry point for Invocations
// =============================================================================
func (t *healthCareChainCode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {

    function, args := stub.GetFunctionAndParameters()
    fmt.Println("invoke is running " + function)

    // Handle different functions
    if function == "initDiagnosis" { // Method for creating a new Diagnosis
        return t.initDiagnosis(stub, args)
    } else if function == "updateDiagnosis" { // Method for updating a diagnosis
        return t.updateDiagnosis(stub, args)
    } else if function == "deleteDiagnosis" { // Method for delete a diagnosis
        return t.deleteDiagnosis(stub, args)
    } else if function == "initBatch" {       // Method for creating a new Batch
        return t.initBatch(stub, args)
    } else if function == "updateBatch" {     // Method for updating a Batch (i.e. adding an readytoThirdparty to a batch).
        return t.updateBatch(stub, args)
    } else if function == "deleteBatch" {     // Method for delete a batch 
        return t.deleteBatch(stub, args)
    } else if function == "initStatement" {   // Method for creating a new Statement
        return t.initStatement(stub, args)
    } else if function == "updateStatement" { // Method for updating a statement
        return t.updateStatement(stub, args)
    } else if function == "deleteStatement" { // Method for delete a statement
        return t.deleteStatement(stub, args)
    } else if function == "queryBatchRecord"{  // Method for looking up a record (e.g. batchID)
        return t.queryBatchRecord(stub, args)
    } else if function == "queryStatementHistory" { // Method for looking up the history of a Statement, all of its prior states on the chain.
        return t.queryStatementHistory(stub, args)
    } else if function == "queryDiagnosisHistory" { // Method for looking up the history of a Diagnosis, all of its prior states on the chain.
        return t.queryDiagnosisHistory(stub, args)
    } else if function == "getAllRecordKeys" { // Method for getting all record keys, for debugging
        return t.getAllRecordKeys(stub, args)
    }

    // If the method provided doesn't match any of our functions, print out debug and return error.
    fmt.Println("Invoke did not find function: " + function)
    return shim.Error("Received unknown function invocation")
}


// =============================================================================
// Functions for Diagnosis
// =============================================================================
func (t *healthCareChainCode) initDiagnosis(stub shim.ChaincodeStubInterface, args []string) pb.Response {
    // We need 8 arguments provided here, in this order:
    // args[0] - DiagnosisID,
    // args[1] - CPTCode,
    // args[2] - CPTName,
    // args[3] - ICD9,
    // args[4] - Description,
    // args[5] - TotalCharge
    // args[6] - PatientAmount
    // args[7] - BilledToInsurance
    // All these should be strings (we'll convert the variable types here).

    if len(args) != 8 {
        return shim.Error("Incorrect number of arguments. Expecting 8: ID [string], CPTCode [string], CPTName [string], ICD9 [string], Description[string], TotalCharge[string], PatientAmount[string], BilledToInsurance[string]")
    }
    if len(args[0]) < 1 {
        return shim.Error("DiagnosisID must be a non-empty string.")
    }
    if len(args[1]) < 1 {
        return shim.Error("CPT Code must be a non-empty string.")
    }
    if len(args[2]) < 1 {
        return shim.Error("CPT Name must be a non-empty string.")
    }
    if len(args[3]) < 1 {
        return shim.Error("ICD9 must be a non-empty string with amount value.")
    }
    if len(args[4]) < 1 {
        return shim.Error("Description must be a non-empty string.")
    }
    if len(args[5]) < 1 {
        return shim.Error("TotalCharge must be a non-empty string with amount value.")
    }
    if len(args[6]) < 1 {
        return shim.Error("PatientAmount must be a non-empty string, if it is not calculated, put 0.0.")
    }
    if len(args[7]) < 1 {
        return shim.Error("BilledToInsurance must be a non-empty string, if it is not calculated, put 0.0.")
    }

    diagnosis := &Diagnosis{args[0], args[1], args[2], args[3], args[4], args[5], args[6], args[7]}
    diagnosisJSONNasBytes, err := json.Marshal(diagnosis)
    if err != nil {
        return shim.Error("Error marshalling data: " + err.Error())
    }
    err = stub.PutState(args[0], diagnosisJSONNasBytes)
    if err != nil {
        return shim.Error("Error putting data to state" + err.Error())
    }

    fmt.Println("- end init diagnosis")
    return shim.Success(nil)
}

func (t *healthCareChainCode) updateDiagnosis(stub shim.ChaincodeStubInterface, args []string) pb.Response{
    // We need 8 arguments provided here, in this order:
    // args[0] - DiagnosisID,
    // args[1] - CPTCode,
    // args[2] - CPTName,
    // args[3] - ICD9,
    // args[4] - Description,
    // args[5] - TotalCharge
    // args[6] - PatientAmount
    // args[7] - BilledToInsurance
    // All these should be strings (we'll convert the variable types here).

    var err error
    if len(args) != 8 {
        return shim.Error("Incorrect number of arguments. Expecting 8: ID [string], CPTCode [string], CPTName [string], ICD9 [string], Description[string], TotalCharge[string], PatientAmount[string], BilledToInsurance[string]")
    }
    if len(args[0]) < 1 {
        return shim.Error("DiagnosisID must be a non-empty string.")
    }
    if len(args[1]) < 1 {
        return shim.Error("CPT Code must be a non-empty string.")
    }
    if len(args[2]) < 1 {
        return shim.Error("CPT Name must be a non-empty string.")
    }
    if len(args[3]) < 1 {
        return shim.Error("ICD9 must be a non-empty string with amount value.")
    }
    if len(args[4]) < 1 {
        return shim.Error("Description must be a non-empty string.")
    }
    if len(args[5]) < 1 {
        return shim.Error("TotalCharge must be a non-empty string with amount value.")
    }
    if len(args[6]) < 1 {
        return shim.Error("PatientAmount must be a non-empty string, if it is not calculated, put 0.0.")
    }
    if len(args[7]) < 1 {
        return shim.Error("BilledToInsurance must be a non-empty string, if it is not calculated, put 0.0.")
    }

    DiagnosisID := args[0]
    DiagnosisAsBytes, err := stub.GetState(DiagnosisID)

    if err != nil {
        return shim.Error("Failed to get Statement: " + err.Error())
    } else if DiagnosisAsBytes == nil {
        return shim.Error("This inspection does not exist: " + DiagnosisID)
    }


    diagnosis := &Diagnosis{DiagnosisID, args[1], args[2], args[3], args[4], args[5], args[6], args[7]}
    diagnosisJSONNasBytes, err := json.Marshal(diagnosis)
    if err != nil {
        return shim.Error("Error marshalling data: " + err.Error())
    }
    err = stub.PutState(args[0], diagnosisJSONNasBytes)
    if err != nil {
        return shim.Error("Error putting data to state" + err.Error())
    }

    fmt.Println("- end update diagnosis")
    return shim.Success(nil)
}

func (t *healthCareChainCode) deleteDiagnosis(stub shim.ChaincodeStubInterface, args []string) pb.Response{
    // args[0] - DiagnosisID,
    // args[1] - CPTCode
    // args[3] - Date
    // All these should be strings (we'll convert the variable types here).
    var err error
    var jsonResp string
    var DiagnosisJSON Diagnosis

    // Make sure we have the right number of arguments.
    if len(args) < 1 {
        return shim.Error("Incorrect number of arguments. Expecting 1")
    }
    if len(args[0]) < 1 {
        return shim.Error("Batch ID must be a non-empty string")
    }

    DiagnosisID := args[0]
    DiagnosisAsBytes, err := stub.GetState(DiagnosisID)

    if err != nil {
        return shim.Error("Failed to get Diagnosis: " + err.Error())
    } else if DiagnosisAsBytes == nil {
        return shim.Error("Failed to get Diagnosis: " + err.Error())
    }
    err = json.Unmarshal([]byte(DiagnosisAsBytes), &DiagnosisJSON)
    if err != nil {
        jsonResp = "{\"Error\":\"Failed to decode JSON of: " + DiagnosisID + "\"}"
        return shim.Error(jsonResp)
    }

    err = stub.DelState(DiagnosisID) //remove the Product from chaincode state
    if err != nil {
        return shim.Error("Failed to delete state:" + err.Error())
    }
    return shim.Success(nil)
}

// =============================================================================
// Functions for Statement
// =============================================================================
func (t *healthCareChainCode) initStatement(stub shim.ChaincodeStubInterface, args []string) pb.Response {
    // We need 8 arguments provided here, in this order:
    // args[0] - StatementID,
    // args[1] - Provider,
    // args[2] - NetworkInOut,
    // args[3] - PatientName,
    // args[4] - DiagnosisID,
    // args[5] - Description
    // args[6] - currentPhase
    // args[7] - attachementLink
    // args[8] - Action
    // args[9] - Date
    // All these should be strings (we'll convert the variable types here).

    if len(args) != 10 {
        return shim.Error("Incorrect number of arguments. Expecting 8: ID [string], Provider [string], NetworkInOut [string], PatientName [string], DiagnosisID[string], Description[string], currentPhase[string], attachementLink[string], Action[int], Date[string:ISO Date]")
    }
    if len(args[0]) < 1 {
        return shim.Error("StatementID must be a non-empty string.")
    }
    if len(args[1]) < 1 {
        return shim.Error("Provider must be a non-empty string.")
    }
    if len(args[2]) < 1 {
        return shim.Error("NetworkInOut must be a non-empty string.")
    }
    if len(args[3]) < 1 {
        return shim.Error("PatientName must be a non-empty string with amount value.")
    }
    if len(args[4]) < 1 {
        return shim.Error("DiagnosisID must be a non-empty string array.")
    }
    if len(args[5]) < 1 {
        return shim.Error("Description must be a non-empty string.")
    }
    if len(args[6]) < 1 {
        return shim.Error("currentPhase must be a non-empty string.")
    }
    if len(args[7]) < 1 {
        return shim.Error("attachementLink must be a non-empty doc link string.")
    }
    if len(args[8]) < 1 {
        return shim.Error("Action must be a non-empty integer.")
    }
    if len(args[9]) < 1 {
        return shim.Error("Date must be a non-empty ISO date string.")
    }

    tags := strings.Split(args[4], "|")

    actionNum, err := strconv.Atoi(args[8])
    action := ActionType(actionNum)
    if err != nil {
        return shim.Error("Failed converting action: " + err.Error())
    }


    statement := &Statement{args[0], args[1], args[2], args[3], tags, args[5], args[6], args[7], action, args[9]}
    // statement := &Statement{args[0]}
    // statement.provider = args[1]
    // statement.NetworkInOut = args[2]
    // statement.PatientName = args[3]
    // statement.DiagnosisID = tags
    // statement.Description = args[5]
    // statement.CurrentPhase = args[6]
    // statement.AttachementLink = args[7]
    // statement.Action = action
    // statement.Date = args[9]
    statementJSONNasBytes, err := json.Marshal(statement)
    if err != nil {
        return shim.Error("Error marshalling data: " + err.Error())
    }
    err = stub.PutState(args[0], statementJSONNasBytes)
    if err != nil {
        return shim.Error("Error putting data to state" + err.Error())
    }

    fmt.Println("- end init statement")
    return shim.Success(nil)
}

func (t *healthCareChainCode) updateStatement(stub shim.ChaincodeStubInterface, args []string) pb.Response{
    // We need 8 arguments provided here, in this order:
    // args[0] - StatementID,
    // args[1] - Provider,
    // args[2] - NetworkInOut,
    // args[3] - PatientName,
    // args[4] - DiagnosisID,
    // args[5] - Description
    // args[6] - currentPhase
    // args[7] - attachementLink
    // args[8] - Action
    // args[9] - Date
    // All these should be strings (we'll convert the variable types here).

    var err error
    if len(args) != 10 {
        return shim.Error("Incorrect number of arguments. Expecting 8: ID [string], Provider [string], NetworkInOut [string], PatientName [string], DiagnosisID[string], Description[string], currentPhase[string], attachementLink[string], Action[int], Date[string:ISO Date]")
    }
    if len(args[0]) < 1 {
        return shim.Error("StatementID must be a non-empty string.")
    }
    if len(args[1]) < 1 {
        return shim.Error("Provider must be a non-empty string.")
    }
    if len(args[2]) < 1 {
        return shim.Error("NetworkInOut must be a non-empty string.")
    }
    if len(args[3]) < 1 {
        return shim.Error("PatientName must be a non-empty string with amount value.")
    }
    if len(args[4]) < 1 {
        return shim.Error("DiagnosisID must be a non-empty string array.")
    }
    if len(args[5]) < 1 {
        return shim.Error("Description must be a non-empty string.")
    }
    if len(args[6]) < 1 {
        return shim.Error("currentPhase must be a non-empty string.")
    }
    if len(args[7]) < 1 {
        return shim.Error("attachementLink must be a non-empty doc link string.")
    }
    if len(args[8]) < 1 {
        return shim.Error("Action must be a non-empty integer.")
    }
    if len(args[9]) < 1 {
        return shim.Error("Date must be a non-empty ISO date string.")
    }

    actionNum, err := strconv.Atoi(args[8])
    if err != nil {
        return shim.Error("Failed converting action: " + err.Error())
    }

    action := ActionType(actionNum)
    StatementID := args[0]
    StatementAsBytes, err := stub.GetState(StatementID)

    if err != nil {
        return shim.Error("Failed to get Statement: " + err.Error())
    } else if StatementAsBytes == nil {
        return shim.Error("This inspection does not exist: " + StatementID)
    }

    tags := strings.Split(args[4], "|")

    
    statement := &Statement{StatementID, args[1], args[2], args[3], tags, args[5], args[6], args[7], action, args[9]}
    // statement.provider = args[1]
    // statement.NetworkInOut = args[2]
    // statement.PatientName = args[3]
    // statement.DiagnosisID = tags
    // statement.Description = args[5]
    // statement.CurrentPhase = args[6]
    // statement.AttachementLink = args[7]
    // statement.Action = args[8]
    // statement.Date = args[9]
    statementJSONNasBytes, err := json.Marshal(statement)
    if err != nil {
        return shim.Error("Error marshalling data: " + err.Error())
    }
    err = stub.PutState(args[0], statementJSONNasBytes)
    if err != nil {
        return shim.Error("Error putting data to state" + err.Error())
    }

    fmt.Println("- end update statement")
    return shim.Success(nil)
}

func (t *healthCareChainCode) deleteStatement(stub shim.ChaincodeStubInterface, args []string) pb.Response{
    // args[0] - StatementID,
    // args[1] - DiagnosisID,
    // args[2] - Provider
    // args[3] - Date
    // All these should be strings (we'll convert the variable types here).
    var err error
    var jsonResp string
    var StatementJSON Statement

    // Make sure we have the right number of arguments.
    if len(args) < 1 {
        return shim.Error("Incorrect number of arguments. Expecting 1")
    }
    if len(args[0]) < 1 {
        return shim.Error("Batch ID must be a non-empty string")
    }

    StatementID := args[0]
    StatementAsBytes, err := stub.GetState(StatementID)

    if err != nil {
        return shim.Error("Failed to get Statement: " + err.Error())
    } else if StatementAsBytes == nil {
        return shim.Error("Failed to get Statement: " + err.Error())
    }
    err = json.Unmarshal([]byte(StatementAsBytes), &StatementJSON)
    if err != nil {
        jsonResp = "{\"Error\":\"Failed to decode JSON of: " + StatementID + "\"}"
        return shim.Error(jsonResp)
    }

    err = stub.DelState(StatementID) //remove the Product from chaincode state
    if err != nil {
        return shim.Error("Failed to delete state:" + err.Error())
    }
    return shim.Success(nil)
}

// =============================================================================
// Functions for Batch
// =============================================================================
func (t *healthCareChainCode) initBatch(stub shim.ChaincodeStubInterface, args []string) pb.Response {
    if len(args) != 5 {
        return shim.Error("Incorrect number of arguments. Expecting 5: BatchID [string], StatementID [string], Description[string], Date[string:ISO Date],  Action[int]")
    }
    if len(args[0]) < 1 {
        return shim.Error("Batch Id must be a non-empty string.")
    }
    if len(args[1]) < 1 {
        return shim.Error("Statement Id must be a non-empty string.")
    }
    if len(args[2]) < 1 {
        return shim.Error("Description must be a non-empty string.")
    }
    if len(args[3]) < 1 {
        return shim.Error("Date must be a non-empty ISO date string.")
    }
    if len(args[4]) < 1 {
        return shim.Error("Action must be a non-empty integer.")
    }

    actionNum, err := strconv.Atoi(args[4])
    action := ActionType(actionNum)
    if err != nil {
        return shim.Error("Failed converting action: " + err.Error())
    }

    statementBatch := &StatementBatch{args[0], args[1], args[2], args[3], action}
    statementBatchJSONNasBytes, err := json.Marshal(statementBatch)
    if err != nil {
        return shim.Error("Error marshalling data: " + err.Error())
    }
    err = stub.PutState(args[0], statementBatchJSONNasBytes)
    if err != nil {
        return shim.Error("Error putting data to state" + err.Error())
    }

    fmt.Println("- end init batch")
    return shim.Success(nil)
}

func (t *healthCareChainCode) updateBatch(stub shim.ChaincodeStubInterface, args []string) pb.Response {
    // We need 5 arguments provided here, in this order:
    // args[0] - BatchID,
    // args[1] - Statement,
    // args[2] - Description,
    // args[3] - Date,
    // args[4] - Action,
    // All these should be strings (we'll convert the variable types here).

    var err error
    if len(args) != 5 {
        return shim.Error("Incorrect number of arguments. Expecting 5: BatchID [string], StatementID [string], Description[string], Date[string:ISO Date],  Action[int]")
    }
    if len(args[0]) < 1 {
        return shim.Error("Batch Id must be a non-empty string.")
    }
    if len(args[1]) < 1 {
        return shim.Error("Statement Id must be a non-empty string.")
    }
    if len(args[2]) < 1 {
        return shim.Error("Description must be a non-empty string.")
    }
    if len(args[3]) < 1 {
        return shim.Error("Date must be a non-empty ISO date string.")
    }
    if len(args[4]) < 1 {
        return shim.Error("Action must be a non-empty integer.")
    }


    actionNum, err := strconv.Atoi(args[4])
    if err != nil {
        return shim.Error("Failed converting action: " + err.Error())
    }

    action := ActionType(actionNum)
    BatchId := args[0]
    BatchAsBytes, err := stub.GetState(BatchId)

    if err != nil {
        return shim.Error("Failed to get Statement: " + err.Error())
    } else if BatchAsBytes == nil {
        return shim.Error("This inspection does not exist: " + BatchId)
    }

    statementBatch := &StatementBatch{BatchId, args[1], args[2], args[3], action}
    statementBatchJSONNasBytes, err := json.Marshal(statementBatch)
    if err != nil {
        return shim.Error("Error marshalling data: " + err.Error())
    }
    err = stub.PutState(args[0], statementBatchJSONNasBytes)
    if err != nil {
        return shim.Error("Error putting data to state" + err.Error())
    }

    fmt.Println("- end init batch")
    return shim.Success(nil)
}

func (t *healthCareChainCode) deleteBatch(stub shim.ChaincodeStubInterface, args []string) pb.Response {
    // args[0] - BatchId,
    // args[1] - StatementID,
    // All these should be strings (we'll convert the variable types here).
    var err error
    var jsonResp string
    var BatchJSON StatementBatch

    // Make sure we have the right number of arguments.
    if len(args) < 1 {
        return shim.Error("Incorrect number of arguments. Expecting 1")
    }
    if len(args[0]) < 1 {
        return shim.Error("Batch ID must be a non-empty string")
    }

    BatchId := args[0]
    BatchAsBytes, err := stub.GetState(BatchId)

    if err != nil {
        return shim.Error("Failed to get Product: " + err.Error())
    } else if BatchAsBytes == nil {
        return shim.Error("Failed to get Product: " + err.Error())
    }
    err = json.Unmarshal([]byte(BatchAsBytes), &BatchJSON)
    if err != nil {
        jsonResp = "{\"Error\":\"Failed to decode JSON of: " + BatchId + "\"}"
        return shim.Error(jsonResp)
    }

    err = stub.DelState(BatchId) //remove the Product from chaincode state
    if err != nil {
        return shim.Error("Failed to delete state:" + err.Error())
    }
    return shim.Success(nil)
}

// =============================================================================
// Functions for Querying Data for Chain
// =============================================================================
func (t *healthCareChainCode) queryBatchRecord(stub shim.ChaincodeStubInterface, args []string) pb.Response {
    var jsonResp string

    if len(args) != 1 {
        return shim.Error("Incorrect number of arguments. Expecting 1: ID [string]")
    }

    valAsbytes, err := stub.GetState(args[0])
    if err != nil {
        jsonResp = "{\"Error\":\"Failed to get state for " + args[0] + "\"}"
        return shim.Error(jsonResp)
    } else if valAsbytes == nil {
        jsonResp = "{\"Error\":\"Batch does not exist: " + args[0] + "\"}"
        return shim.Error(jsonResp)
    }
    return shim.Success(valAsbytes)
}

func (t *healthCareChainCode) queryStatementHistory(stub shim.ChaincodeStubInterface, args []string) pb.Response {
    if len(args) < 1 {
        return shim.Error("Incorrect number of arguments. Expecting 1: ID [string]")
    }

    StatementID := args[0]

    fmt.Printf("- start queryStatementHistory: %s\n", StatementID)

    resultsIterator, err := stub.GetHistoryForKey(StatementID)
    if err != nil {
        return shim.Error(err.Error())
    }
    defer resultsIterator.Close()

    // buffer is a JSON array containing historic values for the Product
    var buffer bytes.Buffer
    buffer.WriteString("[")

    bArrayMemberAlreadyWritten := false
    for resultsIterator.HasNext() {
        response, err := resultsIterator.Next()
        if err != nil {
            return shim.Error(err.Error())
        }
        // Add a comma before array members, suppress it for the first array member
        if bArrayMemberAlreadyWritten == true {
            buffer.WriteString(",")
        }
        buffer.WriteString("{\"TxId\":")
        buffer.WriteString("\"")
        buffer.WriteString(response.TxId)
        buffer.WriteString("\"")

        buffer.WriteString(", \"Value\":")
        buffer.WriteString(string(response.Value))

        buffer.WriteString(", \"Timestamp\":")
        buffer.WriteString("\"")
        buffer.WriteString(time.Unix(response.Timestamp.Seconds, int64(response.Timestamp.Nanos)).String())
        buffer.WriteString("\"")

        buffer.WriteString("}")
        bArrayMemberAlreadyWritten = true
    }
    buffer.WriteString("]")

    fmt.Printf("- getHistoryForStatements returning:\n%s\n", buffer.String())

    return shim.Success(buffer.Bytes())
}

func (t *healthCareChainCode) queryDiagnosisHistory(stub shim.ChaincodeStubInterface, args []string) pb.Response {
    if len(args) < 1 {
        return shim.Error("Incorrect number of arguments. Expecting 1: ID [string]")
    }

    DiagnosisID := args[0]

    fmt.Printf("- start queryDiagnosisHistory: %s\n", DiagnosisID)

    resultsIterator, err := stub.GetHistoryForKey(DiagnosisID)
    if err != nil {
        return shim.Error(err.Error())
    }
    defer resultsIterator.Close()

    // buffer is a JSON array containing historic values for the Product
    var buffer bytes.Buffer
    buffer.WriteString("[")

    bArrayMemberAlreadyWritten := false
    for resultsIterator.HasNext() {
        response, err := resultsIterator.Next()
        if err != nil {
            return shim.Error(err.Error())
        }
        // Add a comma before array members, suppress it for the first array member
        if bArrayMemberAlreadyWritten == true {
            buffer.WriteString(",")
        }
        buffer.WriteString("{\"TxId\":")
        buffer.WriteString("\"")
        buffer.WriteString(response.TxId)
        buffer.WriteString("\"")

        buffer.WriteString(", \"Value\":")
        buffer.WriteString(string(response.Value))

        buffer.WriteString(", \"Timestamp\":")
        buffer.WriteString("\"")
        buffer.WriteString(time.Unix(response.Timestamp.Seconds, int64(response.Timestamp.Nanos)).String())
        buffer.WriteString("\"")

        buffer.WriteString("}")
        bArrayMemberAlreadyWritten = true
    }
    buffer.WriteString("]")

    fmt.Printf("- getHistoryDiagnosis returning:\n%s\n", buffer.String())

    return shim.Success(buffer.Bytes())
}

func (t *healthCareChainCode) getAllRecordKeys(stub shim.ChaincodeStubInterface, args []string) pb.Response {
    startKey := ""
    endKey := ""
    resultsIterator, err := stub.GetStateByRange(startKey, endKey)
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
        buffer.WriteString("}")
        bArrayMemberAlreadyWritten = true
    }
    buffer.WriteString("]")
    fmt.Printf("- queryAllTransaction:\n%s\n", buffer.String())
    return shim.Success(buffer.Bytes())
}






