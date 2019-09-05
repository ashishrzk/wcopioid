// This is an older version of the chaincode, included in this repo for 
// reference.

package main

import (
  "bytes"
  "encoding/json"
  "fmt"
  "strconv"
  "strings"
  "time"
  "github.com/hyperledger/fabric/core/chaincode/shim"
  pb "github.com/hyperledger/fabric/protos/peer"
)

// Empty struct that is used for future chaincode functions/methods.
type PtntVstTraceChaincode struct {
}

// Struct for patient information that should not change often.
type patientMainData struct {
  ObjectType        string  `json:docType` // On all structs, so when it is written to the chain we know what type of record we're seeing.
  PatientId         string  `json:patientId` // Unique ID for patient, for accessing other related records.
  PatientFirstName  string  `json:patientFirstName`
  PatientLastName   string  `json:patientLastName`
  PatientGender     string  `json:patientGender`
  PatientBirthYear  int     `json:patientBirthYear` // Integer for patient year of birth.
  PatientBirthMonth int     `json:patientBirthMonth` // Integer for the month (1-12) when patient was born.
  PatientBirthDay   int     `json:patientBirthDate` // Integer for the day of the month (1-31) when patient was born.
}

// Struct for health data received from FitBit IOT. This is from the object.payload.data.
type fitBitData struct {
  ObjectType            string  `json:docType` // On all structs, so when it is written to the chain we know what type of record we're seeing.
  PatientId             string  `json:patientId` // Unique ID for patient, for accessing other related records.
  ActivityCalories      string  `json:data_activityCalories`
  CaloriesBMR           string  `json:data_caloriesBMR`
  CaloriesOut           string  `json:data_caloriesOut`
  Elevation             string  `json:data_elevation`
  FairlyActiveMinutes   string  `json:data_fairlyActiveMinutes`
  Floors                string  `json:data_floors`
  RestingHeartRate      string  `json:data_restingHeartRate`
  Steps                 string  `json:data_steps`
  VeryActiveMinutes     string  `json:data_veryActiveMinutes`
  MarginalCalories      string  `json:data_marginalCalories`
  SedendaryMinutes      string  `json:data_sedentaryMinutes`
}

// Struct for health data received from Blood Pressure monitor IOT. This is from the object.payload.data.
type bloodPressureMonitorData struct {
  ObjectType            string  `json:docType` // On all structs, so when it is written to the chain we know what type of record we're seeing.
  PatientId             string  `json:patientId` // Unique ID for patient, for accessing other related records.
  SystolicPressure      int     `json:data_systolicpressure`
  DiastolicPressure     int     `json:data_diastolicpressure`
  Heartrate             int     `json:data_heartrate`
}

type prescriptionIssue struct {
  // Struct for when a perscription is issued.
  ObjectType          string  `json:docType` // On all structs, so when it is written to the chain we know what type of record we're seeing.
  PatientId           string  `json:patientId` // Unique ID for patient, for accessing other related records.
  VisitId             string  `json:"visitId"` // Unique ID for doctor visit where perscription was issued.
  RxId                string  `json:rxId` // Unique ID for this perscription.
  RxCode              string  `json:rxCode` // Drug code, for insurance purposes
  RxDrugs             string  `json:rxDrugs` // Name of the drug perscribed.
  RxInstructions      string  `json:rxInstructions` // Instructions (if any)
  DoctorId            string  `json:doctorId` // Unique ID of doctor who issued perscription.
  DoctorName          string  `json:doctorName` // Name of doctor who issued perscription.
}

type prescriptionFulfillment struct {
  // Struct for when a perscription is fulfilled
  ObjectType          string  `json:docType` // On all structs, so when it is written to the chain we know what type of record we're seeing.
  PatientId           string  `json:patientId` // Unique ID for patient, for accessing other related records.
  RxId                string  `json:rxId` // ID for the prescriptionIssue.
  RxTotBill           int     `json:rxTotBill` // Integer bill amount (amount of bill in cents).
  RxInsPay            int     `json:rxInsBill` // Integer insurance payout amount (amount insurance will pay in cents).
  RxCoPay             int     `json:rxCoPay` // Integer copay amount (amount of copay in cents).
  InsId               string  `json:insId` // Patient's insurance ID.
  InsName             string  `json:insName` // Name of Patient's insurer.
  PrescriptionFilled  bool    `json:prescriptionFilled` // Boolean to confirm that the perscription was fulfilled.
}

// Struct for patient visit with doctor.
type medicalVisit struct {
  ObjectType        string  `json:"docType"` // On all structs, so when it is written to the chain we know what type of record we're seeing.
  PatientId         string  `json:"patientId"` // Unique ID for patient, for accessing other related records.
  VisitId           string  `json:"visitId"` // Unique ID for this visit.
  DoctorId          string  `json:"doctorId"` // Unique ID of doctor being visited.
  DoctorName        string  `json:"doctorName"` // Doctor's name
  DoctorTotBill     int     `json:"doctorTotBill"` // Total bill amount, in cents.
  DoctorInsPay      int     `json:"doctorInsPay"` // Amount of bill paid by patient's insurance, in cents.
  DoctorCoPay       int     `json:"doctorCoPay"` // Integer copay amount (amount of copay in cents).
  VisitNotes        string  `json:"visitNotes"` // Field for miscellaneous notes from doctor about visit.
  InsId             string  `json:"insId"` // Patient's insurance ID.
  InsName           string  `json:"insName"` // Name of Patient's insurer.
  VisitDate         int     `json:"visitDate"` // UTC timestamp this one.
}

// =============================================================================
// The main function. We start the needed method from here.
// =============================================================================
func main() {
  err := shim.Start(new(PtntVstTraceChaincode))
  if err != nil {
    fmt.Printf("Error starting Parts Trace chaincode: %s", err)
  }
}

// =============================================================================
// Init initializes chaincode
// =============================================================================
func (t *PtntVstTraceChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
  return shim.Success(nil)
}

// =============================================================================
// Invoke - Our entry point for Invocations
// =============================================================================
func (t *PtntVstTraceChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
  // Get the function that was sent (from the "method" field in the body of message sent to REST endpoint) amd the parameters (from the "args" field in same message body).
  function, args := stub.GetFunctionAndParameters()
  fmt.Println("invoke is running " + function)

  // Handle different functions
  if function == "initPatient" { // Create a new patient
    return t.initPatient(stub, args)
  } else if function == "queryRecord" {
    return t.queryRecord(stub, args)
  } else if function == "medicalVisit" { // Create a new visit by patient to doctor
    return t.medicalVisit(stub, args)
  } else if function == "getAllRecordKeys" { // Getting all record keys, for debugging
    return t.getAllRecordKeys(stub, args)
  } else if function == "initPrescription" { // Create a prescriptionIssue
    return t.initPrescription(stub, args)
  } else if function == "fulfillPrescription" { // Create a prescriptionF
    return t.fulfillPrescription(stub, args)
  } else if function == "queryPrescription" { // Create a prescriptionIssue
    return t.queryPrescription(stub, args)
  } else if function == "getHistoryForRecord" { // Get the history for a record
    return t.getHistoryForRecord(stub, args)
  } else if function == "getCompostiteKey" { // Getting composite keys, links between patient/rx/doctor or patient/visit/doctor.
    return t.getCompostiteKey(stub, args)
  } else if function == "addBloodPressureMonitorData" { // Adding data to patient record from IOT device.
    return t.addBloodPressureMonitorData(stub, args)
  } else if function == "addFitBitData" { // Adding data to patient record from another IOT device.
    return t.addFitBitData(stub, args)
  } else if function == "getResearchData" { // Querying for data for researchers.
    return t.getResearchData(stub, args)
  }
  // If the method provided doesn't match any of our functions, print out debug and return error.
  fmt.Println("Invoke did not find function: " + function)
  return shim.Error("Received unknown function invocation")
}


// =============================================================================
// initPatient - Create a new patient , store into chaincode state
// =============================================================================
func (t *PtntVstTraceChaincode) initPatient(stub shim.ChaincodeStubInterface, args []string) pb.Response {
  var err error

  // Make sure we are getting 4 arguments
  if len(args) < 7 {
    return shim.Error("Incorrect number of arguments. Expecting 7: Patient ID, Patient First Name, Patient Last Name, Patient Gender, Patient Birth Year, Patient Birth Month, Patient Birth Day")
  }

  // ==== Input sanitation ====
  fmt.Println("- start init patient")

  if len(args[0]) < 1 {
    return shim.Error("Patient ID must be a non-empty string")
  }
  if len(args[1]) < 1 {
    return shim.Error("Patient first name must be a non-empty string")
  }
  if len(args[2]) < 1 {
    return shim.Error("Patient last name must be a non-empty string")
  }
  if len(args[3]) < 1 {
    return shim.Error("Patient gender must be a non-empty string")
  }
  if len(args[4]) < 1 {
    return shim.Error("Patient birth year must be a non-empty string")
  }
  if len(args[5]) < 1 {
    return shim.Error("Patient birth month must be a non-empty string")
  }
  if len(args[6]) < 1 {
    return shim.Error("Patient birth day must be a non-empty string")
  }

  // Putting the inputs into variables, and doing some further input sanitation along the way.
  patientId := args[0]
  patientFirstName := strings.ToLower(args[1])
  patientLastName := strings.ToLower(args[2])
  patientGender := strings.ToLower(args[3])
  patientBirthYear, err := strconv.Atoi(args[4])
  if err != nil {
    return shim.Error("Patient birth year must be a numeric string")
  }
  patientBirthMonth, err := strconv.Atoi(args[5])
  if err != nil {
    return shim.Error("Patient birth month must be a numeric string")
  }
  if patientBirthMonth < 1 || patientBirthMonth > 12 {
    return shim.Error("Patient birth month must be a between 1 and 12")
  }
  patientBirthDay, err := strconv.Atoi(args[6])
  if err != nil {
    return shim.Error("Patient birth day must be a numeric string")
  }
  if patientBirthMonth < 1 || patientBirthMonth > 31 {
    return shim.Error("Patient birth day must be a between 1 and 31")
  }


  //Check if patient already exists, trying to pull records for the patient ID
  patientAsBytes, err := stub.GetState(patientId)
  if err != nil {
    return shim.Error("Error received with patient: " + err.Error())
  } else if patientAsBytes != nil {
    return shim.Error("This patient ID already exists: " + patientId)
  }


  // ==== Create patient and marshal to JSON ====
  objectType := "patientMainData" // Give the objectType the name of the struct we use to save the record.

  patient := &patientMainData{
    objectType,
    patientId,
    patientFirstName,
    patientLastName,
    patientGender,
    patientBirthYear,
    patientBirthMonth,
    patientBirthDay }
  patientJSONasBytes, err := json.Marshal(patient)
  if err != nil {
    return shim.Error("Error marshalling data to JSON: " + err.Error())
  }

  // Create a record key that lets us search these separately.
  recordKey := objectType + patientId

  // === Save visit to state ===
  err = stub.PutState(recordKey, patientJSONasBytes)
  if err != nil {
    return shim.Error("Error committing record to state: " + err.Error())
  }


////////////////////////////////////////////////////////////////////////////////
  // LANCE NOTE: Not gonna lie, I still don't understand composite keys just yet. Leaving this code in for now, though. Just in case.

  //  ==== Index the mobile based on the owner
  //  An 'index' is a normal key/value entry in state.
  //  The key is a composite key, with the elements that you want to range query on listed first.
  //  In our case, the composite key is based on patient~details
  //  This will enable very efficient state range queries based on composite keys matching patient~details~*
  patientIndex := "patient~details"
  patientIndexKey, err := stub.CreateCompositeKey(patientIndex, []string{ patient.PatientLastName, patient.PatientFirstName, patient.PatientId})
  if err != nil {
    return shim.Error(err.Error())
  }
  //  Save index entry to state. Only the key name is needed, no need to store a duplicate copy of the marble.
  //  Note - passing a 'nil' value will effectively delete the key from state, therefore we pass null character as value
  value := []byte{0x00}
  stub.PutState(patientIndexKey, value)
////////////////////////////////////////////////////////////////////////////////


  // ==== We created our patient record, return success ====
  fmt.Println("- end init patient")
  return shim.Success(nil)
}


// =============================================================================
// queryRecord - Grabbing a record by key, this pulls raw record and isn't
// de-identified.
// =============================================================================
func (t *PtntVstTraceChaincode) queryRecord(stub shim.ChaincodeStubInterface, args []string) pb.Response {

  if len(args) != 1 {
    return shim.Error("Incorrect number of arguments. Need the record ID.")
  }

  key := args[0]

  recordAsBytes, _ := stub.GetState(key)

  if recordAsBytes == nil {
    return shim.Error("Could not locate record")
  }

  return shim.Success(recordAsBytes)
}



// =============================================================================
// medicalVisit - create a new doctor's office visit, store into chaincode state
// =============================================================================
func (t *PtntVstTraceChaincode) medicalVisit(stub shim.ChaincodeStubInterface, args []string) pb.Response {

  // We need 11 arguments provided here, in this order:
  // args[0] - PatientId,
  // args[1] - VisitId,
  // args[2] - DoctorId,
  // args[3] - DoctorName,
  // args[4] - DoctorTotBill,
  // args[5] - DoctorInsPay,
  // args[6] - DoctorCoPay,
  // args[7] - VisitNotes,
  // args[8] - InsId,
  // args[9] - InsName,
  // args[10] - VisitDate
  // All these should be strings (we'll convert the variable types here).
  // DoctorInsPay/VisitNotes/InsId/InsName can be empty.

  var err error

  // Make sure we have the right number of arguments.
  if len(args) < 11 {
    return shim.Error("Incorrect number of arguments. Expecting 11")
  }

  // ==== Input sanitation ====
  fmt.Println("- start init visit")

  if len(args[0]) < 1 {
    return shim.Error("Patient ID must be a non-empty string")
  }
  if len(args[1]) < 1 {
    return shim.Error("Visit ID must be a non-empty string")
  }
  if len(args[2]) < 1 {
    return shim.Error("Doctor ID must be a non-empty string")
  }
  if len(args[3]) < 1 {
    return shim.Error("Doctor name must be a non-empty string")
  }
  if len(args[4]) < 1 {
    return shim.Error("Doctor total bill must be a non-empty string. If there is no charge for visit, enter '0'")
  }
  if len(args[6]) <= 0 {
    return shim.Error("Copay must be a non-empty string. If there is no copayment for visit, enter '0'")
  }
  if len(args[10]) <= 0 {
    return shim.Error("Visit date name be a non-empty string")
  }

  // Add patient ID to variable.
  patientId := args[0]
  // Add visit ID to variable
  visitId := args[1]
  // Add dictor ID to variable
  doctorId := args[2]
  // Add doctor's name to variable
  doctorName := strings.ToLower(args[3])
  // Add the total bill amount to variable, and sanitize input.
  doctorTotBill, err := strconv.Atoi(args[4])
  if err != nil {
    return shim.Error("Doctor total bill must be a numeric string representing total cents. Make sure there is no period.")
  }
  doctorInsPay := 0
  // Put the insurance amount into variable, and sanitize input
  if len(args[5]) > 0 { // If doctorInsPay isn't empty, convert to integer.
    doctorInsPay, err = strconv.Atoi(args[5])
    if err != nil {
      return shim.Error("Insurance payment amount must be a numeric string representing total cents. Make sure there is no period.")
    }
  }
  // Put the copay amount into variable, and sanitize input
  doctorCoPay, err := strconv.Atoi(args[6])
  if err != nil {
    return shim.Error("Insurance payment amount must be a numeric string representing total cents. Make sure there is no period.")
  }
  // Put visit notes into variable. (This may be an empty string, but that's okay.)
  visitNotes := strings.ToLower(args[7])
  // Put insurer ID into variable. (This may be an empty string, but that's okay.)
  insId :=strings.ToLower(args[8])
  // Put insurer name into variable. (This may be an empty string, but that's okay.)
  insName := strings.ToLower(args[9])
  // Put visit date timestamp into variable.
  visitDate, err := strconv.Atoi(args[10])
  if err != nil {
    return shim.Error("Visit date of birth must be a UTC timestamp.")
  }

  // ==== Check if patient already exists ====
  patientIDCheck := "patientMainData" + patientId
  patientAsBytes, err := stub.GetState(patientIDCheck)
  if err != nil {
    return shim.Error("Failed to get patient: " + err.Error())
  } else if patientAsBytes == nil {
    return shim.Error("This patient does not exist: " + patientId)
  }

  // ==== Check if visit already exists ====
  visitAsBytes, err := stub.GetState(visitId)
  if err != nil {
    return shim.Error("Failed to get Visit: " + err.Error())
  } else if visitAsBytes != nil {
    return shim.Error("This visit already exists: " + visitId)
  }

  // ==== Create vist  and marshal to JSON ====
  objectType := "medicalVisit"

  visit := &medicalVisit{
    objectType,
    patientId,
    visitId,
    doctorId,
    doctorName,
    doctorTotBill,
    doctorInsPay,
    doctorCoPay,
    visitNotes,
    insId,
    insName,
    visitDate }
  visitJSONasBytes, err := json.Marshal(visit)
  if err != nil {
    return shim.Error(err.Error())
  }

  recordKey := objectType + patientId + visitId

  // === Save mobile to state ===
  err = stub.PutState(recordKey, visitJSONasBytes)
  if err != nil {
    return shim.Error(err.Error())
  }

////////////////////////////////////////////////////////////////////////////////
  // LANCE NOTE: Not gonna lie, I still don't understand composite keys just yet. Leaving this code in for now, though. Just in case.

  //  ==== Index the visit based on the owner
  //  An 'index' is a normal key/value entry in state.
  //  The key is a composite key, with the elements that you want to range query on listed first.
  //  In our case, the composite key is based on patient~details
  //  This will enable very efficient state range queries based on composite keys matching patient~details~*

  visitIndex := "patient~visit~doctor"
  visitIndexKey, err := stub.CreateCompositeKey(visitIndex, []string{visit.PatientId, visit.VisitId, visit.DoctorId})
  if err != nil {
    return shim.Error(err.Error())
  }
  //  Save index entry to state. Only the key name is needed, no need to store a duplicate copy of the marble.
  //  Note - passing a 'nil' value will effectively delete the key from state, therefore we pass null character as value
  value := []byte{0x00}
  stub.PutState(visitIndexKey, value)
////////////////////////////////////////////////////////////////////////////////

  // ==== Visit part saved and indexed. Return success ===
  fmt.Println("- end init Visit")
  return shim.Success(nil)
}



// =============================================================================
// getAllRecordKeyss - get all records in the chaincode/channel state.
// This is for debug purposes only.
// =============================================================================
func (t *PtntVstTraceChaincode) getAllRecordKeys(stub shim.ChaincodeStubInterface, args []string) pb.Response {

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


//==============================================================================
// initPrescription - create a prescriptionIssue, store into chaincode state
//==============================================================================
func (t *PtntVstTraceChaincode) initPrescription(stub shim.ChaincodeStubInterface, args []string) pb.Response {

  // ObjectType
  // args[0] - PatientId
  // args[1] - VisitId
  // args[2] - RxId
  // args[3] - RxCode
  // args[4] - RxDrugs
  // args[5] - RxInstructions
  // args[6] - DoctorId
  // args[7] - DoctorName

  var err error

  // @MODIFY_HERE extend to expect 8 arguements
  if len(args) < 8 {
    return shim.Error("Incorrect number of arguments. Expecting 8")
  }

  // ==== Input sanitation and variable assignment ====
  fmt.Println("- start init prescription")

  if len(args[0]) < 1 {
    return shim.Error("Patient ID must be a non-empty string")
  }
  patientId := args[0]

  if len(args[1]) < 1 {
    return shim.Error("Visit ID must be a non-empty string")
  }
  visitId := args[1]

  if len(args[2]) < 1 {
    return shim.Error("RX ID must be a non-empty string")
  }
  rxId := args[2]

  if len(args[3]) < 1 {
    return shim.Error("Drug code must be a non-empty string")
  }
  rxCode := args[3]

  if len(args[4]) < 1 {
    return shim.Error("Drug name must be a non-empty string")
  }
  // rxInstructions can be an empty string.
  rxDrugs := strings.ToLower(args[4])

  rxInstructions := args[5]

  if len(args[6]) <= 0 {
    return shim.Error("Doctor ID be a non-empty string")
  }
  doctorId := args[6]

  if len(args[7]) <= 0 {
    return shim.Error("Doctor name date be a non-empty string")
  }
  doctorName := strings.ToLower(args[7])


  // ==== Check if patient already exists ====
  patientIDCheck := "patientMainData" + patientId
  patientAsBytes, err := stub.GetState(patientIDCheck)
  if err != nil {
    return shim.Error("Failed to get patient: " + err.Error())
  } else if patientAsBytes == nil {
    return shim.Error("This patient does not exist: " + patientId)
  }

  // ==== Check if prescription already exists ====
  prescriptionAsBytes, err := stub.GetState(rxId)
  if err != nil {
    return shim.Error("Failed to get perscription: " + err.Error())
  } else if prescriptionAsBytes != nil {
    return shim.Error("This prescription already exists: " + rxId)
  }

  // ==== Create perscription and marshal to JSON ====
  objectType := "prescriptionIssue"

  prescription := &prescriptionIssue{
    objectType,
    patientId,
    visitId,
    rxId,
    rxCode,
    rxDrugs,
    rxInstructions,
    doctorId,
    doctorName }
  prescriptionJSONasBytes, err := json.Marshal(prescription)
  if err != nil {
    return shim.Error(err.Error())
  }

  // Key labels will be the RX ID and the object type.
  keyLabel := objectType + patientId + rxId

  // === Save mobile to state ===
  err = stub.PutState(keyLabel, prescriptionJSONasBytes)
  if err != nil {
    return shim.Error(err.Error())
  }

////////////////////////////////////////////////////////////////////////////////
  // LANCE NOTE: Not gonna lie, I still don't understand composite keys just yet. Leaving this code in for now, though. Just in case.

  //  ==== Index the mobile based on the owner
  //  An 'index' is a normal key/value entry in state.
  //  The key is a composite key, with the elements that you want to range query on listed first.
  //  In our case, the composite key is based on patient~details
  //  This will enable very efficient state range queries based on composite keys matching patient~details~*

  presriptionIndex := "patient~rx~doctor"
  presriptionIndexKey, err := stub.CreateCompositeKey(presriptionIndex, []string{prescription.PatientId, prescription.RxId, prescription.DoctorId})
  if err != nil {
    return shim.Error(err.Error())
  }
  //  Save index entry to state. Only the key name is needed, no need to store a duplicate copy of the marble.
  //  Note - passing a 'nil' value will effectively delete the key from state, therefore we pass null character as value
  value := []byte{0x00}
  stub.PutState(presriptionIndexKey, value)
////////////////////////////////////////////////////////////////////////////////

  // ==== Visit part saved and indexed. Return success ====
  fmt.Println("- end init prescription")
  return shim.Success(nil)
}

//==============================================================================
// fulfillPrescription - If pharmacy and insurance agree to fulfill a
// prescription, store a prescriptionFulfillment record in chaincode state
//==============================================================================
func (t *PtntVstTraceChaincode) fulfillPrescription(stub shim.ChaincodeStubInterface, args []string) pb.Response {

  // ObjectType
  // args[0] - PatientId
  // args[1] - RxId
  // args[2] - RxTotBill
  // args[3] - RxInsPay
  // args[4] - RxCoPay
  // args[5] - InsId
  // args[6] - InsName
  // args[7] - PrescriptionFilled

  var err error

  // Make sure the right number of arguments are provided.
  if len(args) < 8 {
    return shim.Error("Incorrect number of arguments. Expecting 8")
  }

  // ==== Input sanitation and variable assignment ====
  fmt.Println("- start init prescription fulfillment")

  if len(args[0]) < 1 {
    return shim.Error("Patient ID must be a non-empty string")
  }
  patientId := args[0]
  patientIDCheck := "patientMainData" + patientId
  // ==== Check if patient already exists ====
  patientAsBytes, err := stub.GetState(patientIDCheck)
  if err != nil {
    return shim.Error("Failed to get patient: " + err.Error())
  } else if patientAsBytes == nil {
    return shim.Error("This patient does not exist: " + patientId)
  }


  if len(args[1]) < 1 {
    return shim.Error("RX ID must be a non-empty string")
  }
  rxId := args[1]
  // ==== Check if RX was actually issued ====
  rxIssue := "prescriptionIssue" + patientId + rxId
  _, err = stub.GetState(rxIssue)
  if err != nil {
    return shim.Error("Failed to get RX: " + err.Error())
  } else if patientAsBytes == nil {
    return shim.Error("This RX does not exist: " + rxId)
  }


  if len(args[2]) < 1 {
    return shim.Error("RX total bill must be a non-empty string. If there is no charge for RX, send '0'.")
  }
  rxTotBill, err := strconv.Atoi(args[2])
  if err != nil {
    return shim.Error("RX total bill amount must be a numeric string representing total cents. Make sure there is no period.")
  }


  if len(args[3]) < 1 {
    return shim.Error("RX insurance payment amount must be a non-empty string. If there is no insurance payment for RX, send '0'.")
  }
  rxInsPay, err := strconv.Atoi(args[3])
  if err != nil {
    return shim.Error("RX insurance payment amount must be a numeric string representing payment in cents. Make sure there is no period.")
  }


  if len(args[4]) < 1 {
    return shim.Error("RX copayment amount must be a non-empty string. If there is no copayment for RX, send '0'.")
  }
  rxCoPay, err := strconv.Atoi(args[3])
  if err != nil {
    return shim.Error("RX copayment amount must be a numeric string representing copayment amount in cents. Make sure there is no period.")
  }


  insId := args[5]
  insName := args[6]


  if len(args[7]) < 1 {
    return shim.Error("Fulfillment of RX needs to be the string 'true' or 'false'.")
  }
  prescriptionFilled, err := strconv.ParseBool(args[7])
  if err != nil {
    return shim.Error("Fulfillment of RX needs to be the string 'true' or 'false'.")
  }

  // ==== Create fulfillment and marshal to JSON ====
  objectType := "prescriptionFulfillment"
  // Key labels will be the RX ID and the object type.
  keyLabel := objectType + patientId + rxId
  fulfillment := &prescriptionFulfillment{
    objectType,
    patientId,
    rxId,
    rxTotBill,
    rxInsPay,
    rxCoPay,
    insId,
    insName,
    prescriptionFilled }
  fulfillmentJSONasBytes, err := json.Marshal(fulfillment)
  if err != nil {
    return shim.Error(err.Error())
  }

  // ==== Save filfullment to state ====
  err = stub.PutState(keyLabel, fulfillmentJSONasBytes)
  if err != nil {
    return shim.Error(err.Error())
  }

  // ==== Prescription fulfillment record saved and indexed. Return success ====
  fmt.Println("- end fulfill prescription")
  return shim.Success(nil)
}



//==============================================================================
// queryPrescription - Get the record and history for a prescription
//==============================================================================
func (t *PtntVstTraceChaincode) queryPrescription(stub shim.ChaincodeStubInterface, args []string) pb.Response {

  if len(args) != 2 {
    return shim.Error("Incorrect number of arguments. Expecting 2: the patient ID and prescription ID.")
  }

  rxId := args[0]
  patientID := args[1]
  issueKey := "prescriptionIssue" + patientID + rxId
  fulfillmentKey := "prescriptionFulfillment" + patientID + rxId

  fmt.Printf("- start queryPrescription: %s\n", rxId)

  issueAsBytes, err := stub.GetState(issueKey)
  if err != nil {
    return shim.Error("Failed to get RX: " + err.Error())
  } else if issueAsBytes == nil {
    return shim.Error("No RX with that ID for that patient.")
  }



  // Buffer is where we're writing our result.
  var buffer bytes.Buffer

  // Opening the JSON we'll be sending in our reply.
  buffer.WriteString("{")
  // Writing the issue to the buffer.
    buffer.WriteString("\"Issue\" : ")
      buffer.WriteString(string(issueAsBytes))
    // Close out the item
    buffer.WriteString(",")
    buffer.WriteString("\"Fulfillment\" : [")
      // Get all history for fulfillments attempts for this RX.
      fulfillmentIterator, err := stub.GetHistoryForKey(fulfillmentKey)
      if err != nil {
        return shim.Error("Error checking for fulfillment: " + err.Error())
      }
      // Boolean to check if we're on the first item or not.
      bArrayMemberAlreadyWritten := false
      // Loop through all the fulfillment records.
      for fulfillmentIterator.HasNext() {
        fulfillment, err := fulfillmentIterator.Next()
        // If there is an error, throw it.
        if err != nil {
          return shim.Error(err.Error())
        }
        // Write a comman if we aren't on the first item.
        if bArrayMemberAlreadyWritten == true {
          buffer.WriteString(",")
        }
        // Write the transaction ID.
        buffer.WriteString("{\"TxId\":")
        buffer.WriteString("\"")
        buffer.WriteString(fulfillment.TxId)
        buffer.WriteString("\"")
        // Write the value for the string.
        buffer.WriteString(", \"Value\":")
        buffer.WriteString(string(fulfillment.Value))
        // Write the timestamp when the trip was committed.
        buffer.WriteString(", \"Timestamp\":")
        buffer.WriteString("\"")
        buffer.WriteString(time.Unix(fulfillment.Timestamp.Seconds, int64(fulfillment.Timestamp.Nanos)).String())
        buffer.WriteString("\"")
        // Close out the item
        buffer.WriteString("}")
        bArrayMemberAlreadyWritten = true
      }
    buffer.WriteString("]")
    // Closing out the array of history records.
  buffer.WriteString("}")
  // Closing out the JSON we'll be sending in our reply.

  fmt.Printf("- queryPrescription:\n%s\n", buffer.String())
  return shim.Success(buffer.Bytes())
}



//==============================================================================
// getHistoryForRecord returns the histotical state transitions for a given key
// of a record. This is full info and needs cleanup.
//==============================================================================
func (t *PtntVstTraceChaincode) getHistoryForRecord(stub shim.ChaincodeStubInterface, args []string) pb.Response {

  if len(args) < 1 {
    return shim.Error("Incorrect number of arguments. Expecting 1")
  }

  recordKey := args[0]

  fmt.Printf("- start getHistoryForRecord: %s\n", recordKey)

  resultsIterator, err := stub.GetHistoryForKey(recordKey)
  if err != nil {
    return shim.Error(err.Error())
  }
  defer resultsIterator.Close()

  // buffer is a JSON array containing historic values for the key/value pair
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
    // if it was a delete operation on given key, then we need to set the
    //corresponding value null. Else, we will write the response.Value
    if response.IsDelete {
      buffer.WriteString("null")
    } else {
      buffer.WriteString(string(response.Value))
    }
    buffer.WriteString(", \"Timestamp\":")
    buffer.WriteString("\"")
    buffer.WriteString(time.Unix(response.Timestamp.Seconds, int64(response.Timestamp.Nanos)).String())
    buffer.WriteString("\"")

    buffer.WriteString("}")
    bArrayMemberAlreadyWritten = true
  }
  buffer.WriteString("]")

  fmt.Printf("- getHistoryForRecord returning:\n%s\n", buffer.String())

  return shim.Success(buffer.Bytes())
}


//==============================================================================
// Searches based on the composite key and returns all the data requested for
// based on search string and the key like, Search String
// (patient~rx~doctor, patient~visit~doctor) and
// Key could be Patient ID or Doctor ID or Visit ID
// Work around in case, couchDB cannot be configured
// LANCE NOTE: Keeping this as well though it may not see use.
//==============================================================================

func (t *PtntVstTraceChaincode) getCompostiteKey(stub shim.ChaincodeStubInterface, args []string) pb.Response {

  if len(args) < 2 {
    return shim.Error("Incorrect number of arguments. Expecting 2")
  }

  keyString := args[0]
  keyToSearch := args[1]

  resultsIterator, err := stub.GetStateByPartialCompositeKey(keyString, []string{keyToSearch})
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
    objectType, compositeKeyParts, err := stub.SplitCompositeKey(queryResponse.Key)
    if err != nil {
      return shim.Error(err.Error())
    }
    key1 := compositeKeyParts[0]
    key2 := compositeKeyParts[1]
    key3 := compositeKeyParts[2]
    // Add a comma before array members, suppress it for the first array member
    if bArrayMemberAlreadyWritten == true {
      buffer.WriteString(",")
    }
    buffer.WriteString("{\"key1\":")
    buffer.WriteString("\"")
    buffer.WriteString(string(key1))
    buffer.WriteString("\"")

    buffer.WriteString(",\"key2\":")
    buffer.WriteString("\"")
    buffer.WriteString(string(key2))
    buffer.WriteString("\"")

    buffer.WriteString(",\"key3\":")
    buffer.WriteString("\"")
    buffer.WriteString(string(key3))
    buffer.WriteString("\"")

    buffer.WriteString(",\"objectType\":")
    buffer.WriteString("\"")
    buffer.WriteString(string(objectType))
    buffer.WriteString("\"")

    // buffer.WriteString(", \"Record\":")
    // // Record is a JSON object, so we write as-is
    // buffer.WriteString(string(queryResponse.Value))
    buffer.WriteString("}")
    bArrayMemberAlreadyWritten = true
  }
  buffer.WriteString("]")
  fmt.Printf("- getCompositeKey queryResult:\n%s\n", buffer.String())
  return shim.Success(buffer.Bytes())
}

//==============================================================================
// addBloodPressureMonitorData - Method to write data from an IOT device
// and put it on the blockchain.
//==============================================================================
func (t *PtntVstTraceChaincode) addBloodPressureMonitorData(stub shim.ChaincodeStubInterface, args []string) pb.Response {

  // ObjectType
  // args[0] - PatientId
  // args[1] - SystolicPressure
  // args[2] - DiastolicPressure
  // args[3] - Heartrate

  if len(args) != 4 {
    return shim.Error("Incorrect number of arguments. Expecting 4")
  }

  objectType := "bloodPressureMonitorData"

  // ========= Assigning variables, sanitizing inputs. =========
  if len(args[0]) < 1 {
    return shim.Error("Patient ID cannot be blank.")
  }
  patientId := args[0]
  patientKey := "patientMainData" + patientId
  patientAsBytes, err := stub.GetState(patientKey)
  if err != nil {
    return shim.Error("Error received with patient: " + err.Error())
  } else if patientAsBytes == nil {
    return shim.Error("No patient exists with this ID: " + patientId)
  }

  if len(args[1]) < 1 {
    return shim.Error("Systolic pressure cannot be blank.")
  }
  systolicPressure, err := strconv.Atoi(args[1])
  if err != nil {
    return shim.Error("Systolic pressure must be an integer string.")
  }

  if len(args[2]) < 1 {
    return shim.Error("Diastolic pressure cannot be blank.")
  }
  diastolicPressure, err := strconv.Atoi(args[2])
  if err != nil {
    return shim.Error("Diastolic pressure must be an integer string.")
  }

  if len(args[3]) < 1 {
    return shim.Error("Heart rate cannot be blank.")
  }
  heartRate, err := strconv.Atoi(args[3])
  if err != nil {
    return shim.Error("Heart rate must be an integer string.")
  }

  recordKey := objectType + patientId
  iotRecord := &bloodPressureMonitorData{
    objectType,
    patientId,
    systolicPressure,
    diastolicPressure,
    heartRate }

  iotRecordJSONasBytes, err := json.Marshal(iotRecord)
  if err != nil {
    return shim.Error("Error formatting data for use.")
  }

  // === Save record to state ===
  err = stub.PutState(recordKey, iotRecordJSONasBytes)
  if err != nil {
    return shim.Error(err.Error())
  }

  return shim.Success(nil)
}



//==============================================================================
// addFitBitData - Adding data from IOT services gathered from a device
// that we can use
//==============================================================================
func (t *PtntVstTraceChaincode) addFitBitData(stub shim.ChaincodeStubInterface, args []string) pb.Response {

  // ObjectType
  // args[0] - PatientId
  // args[1] - ActivityCalories
  // args[2] - CaloriesBMR
  // args[3] - CaloriesOut
  // args[4] - Elevation
  // args[5] - FairlyActiveMinutes
  // args[6] - Floors
  // args[7] - RestingHeartRate
  // args[8] - Steps
  // args[9] - VeryActiveMinutes
  // args[10] - MarginalCalories
  // args[11] - SedendaryMinutes

  if len(args) != 12 {
    return shim.Error("Incorrect number of arguments. Expecting 12. Send empty strings for any missing pieces of data.")
  }

  objectType := "fitBitData"

  // ========= Assigning variables, sanitizing inputs. =========
  if len(args[0]) < 1 {
    return shim.Error("Patient ID cannot be blank.")
  }
  patientId := args[0]
  patientKey := "patientMainData" + patientId
  patientAsBytes, err := stub.GetState(patientKey)
  if err != nil {
    return shim.Error("Error received with patient: " + err.Error())
  } else if patientAsBytes == nil {
    return shim.Error("No patient exists with this ID: " + patientId)
  }
  // ===== These can be empty and we're storing them all as strings
  // ===== so no sanitizing going on. ======
  activityCalories := args[1]
  caloriesBMR := args[2]
  caloriesOut := args[3]
  elevation := args[4]
  fairlyActiveMinutes := args[5]
  floors := args[6]
  restingHeartRate := args[7]
  steps := args[8]
  veryActiveMinutes := args[9]
  marginalCalories := args[10]
  sedendaryMinutes := args[11]

  recordKey := objectType + patientId
  iotRecord := &fitBitData{
    objectType,
    patientId,
    activityCalories,
    caloriesBMR,
    caloriesOut,
    elevation,
    fairlyActiveMinutes,
    floors,
    restingHeartRate,
    steps,
    veryActiveMinutes,
    marginalCalories,
    sedendaryMinutes }

  iotRecordJSONasBytes, err := json.Marshal(iotRecord)
  if err != nil {
    return shim.Error("Error formatting data for use.")
  }

  // === Save record to state ===
  err = stub.PutState(recordKey, iotRecordJSONasBytes)
  if err != nil {
    return shim.Error(err.Error())
  }

  return shim.Success(nil)
}



//==============================================================================
// getResearchData - Pulling records for use removing identifiers so it
// isn't PHI
//==============================================================================
func (t *PtntVstTraceChaincode) getResearchData(stub shim.ChaincodeStubInterface, args []string) pb.Response {

  currentTime := time.Now()

  // First we grab the patients...
  patientsResultsIterator, err := stub.GetStateByRange("patientMainData", "patientMainDatazzzz")
  if err != nil {
    return shim.Error("Error committing initial search for patient records.")
  }
  // Defer closing out the iterator we've opened up.
  defer patientsResultsIterator.Close()

  // Create the buffer that will be our output.
  var buffer bytes.Buffer
  buffer.WriteString("[")

    // Boolean so we know if we're on the first item...
    patientArrayMemberAlreadyWritten := false
    // Iterate through the array of patients...
    for patientsResultsIterator.HasNext() {
      patient, err := patientsResultsIterator.Next()
      if err != nil {
        return shim.Error("Error getting patient - " + err.Error())
      }
      var patientData patientMainData
      err = json.Unmarshal(patient.Value, &patientData)
      if err != nil {
        return shim.Error(err.Error())
      }
      // Add a comma before array members, suppress it for the first array member
      if patientArrayMemberAlreadyWritten == true {
        buffer.WriteString(",")
      }
      buffer.WriteString("{")
        buffer.WriteString("\"Birth Year\":")
        // Find out if patient age is over 89, as that will require extra de-identification.
        ageCheck := currentTime.Year() - patientData.PatientBirthYear
        if ageCheck < 89 {
          buffer.WriteString(strconv.Itoa(patientData.PatientBirthYear))
        } else {
          buffer.WriteString(">")
          buffer.WriteString(strconv.Itoa(currentTime.Year()-88))
        }
        buffer.WriteString(",")
        buffer.WriteString("\"Gender\":\"")
        buffer.WriteString(patientData.PatientGender)
        buffer.WriteString("\",")

        buffer.WriteString("\"MedicalVisits\":")
        buffer.WriteString("[")
        // Begin writing the medical visits for patient.
        medicalVisitStartKey := "medicalVisit" + patientData.PatientId
        medicalVisitEndKey := "medicalVisit" + patientData.PatientId + "zzzzz"
        medicalVisitsIterator, err := stub.GetStateByRange(medicalVisitStartKey, medicalVisitEndKey)
        medicalVisitsIteratorWrittenBoolean := false
        for medicalVisitsIterator.HasNext() {
          if medicalVisitsIteratorWrittenBoolean == true {
            buffer.WriteString(",")
          }
          visit, err := medicalVisitsIterator.Next()
          if err != nil {
            return shim.Error("Error attempting to get medical visit information.")
          }
          var visitData medicalVisit
          err = json.Unmarshal(visit.Value, &visitData)
          if err != nil {
            return shim.Error("Error parsing visit data from JSON.")
          }
          buffer.WriteString("{")
            buffer.WriteString("\"VisitDate\":\"")
            buffer.WriteString(strconv.Itoa(visitData.VisitDate))
            buffer.WriteString("\",")
            buffer.WriteString("\"DoctorName\":\"")
            buffer.WriteString(visitData.DoctorName)
            buffer.WriteString("\",")
            buffer.WriteString("\"VisitNotes\":\"")
            buffer.WriteString(visitData.VisitNotes)
            buffer.WriteString("\"")
          buffer.WriteString("}")
          medicalVisitsIteratorWrittenBoolean = true
        }
        buffer.WriteString("],")

        // Begin writing all the prescriptions.

        buffer.WriteString("\"Prescriptions\":")
        buffer.WriteString("[")
          // Now we grab the prescriptions...
          prescriptionIssueStartKey := "prescriptionIssue" + patientData.PatientId
          prescriptionIssueEndKey := "prescriptionIssue" + patientData.PatientId + "zzzzz"
          prescriptionIssuesIterator, err := stub.GetStateByRange(prescriptionIssueStartKey, prescriptionIssueEndKey)
          if err != nil {
            return shim.Error("Error checking for prescriptions.")
          }
          prescriptionMemberAlreadyWritten := false
          for prescriptionIssuesIterator.HasNext() {
            prescription, err := prescriptionIssuesIterator.Next()
            if err != nil {
              return shim.Error("Error attempting to get prescription issue item.")
            }
            if prescriptionMemberAlreadyWritten == true {
              buffer.WriteString(",")
            }
            var thePrescription prescriptionIssue
            err = json.Unmarshal(prescription.Value, &thePrescription)
            if err != nil {
              return shim.Error("Error parsing prescription issue data from JSON.")
            }
            buffer.WriteString("{")
              // Write the elements of prescription...
              buffer.WriteString("\"RxDrugs\":\"")
              buffer.WriteString(thePrescription.RxDrugs)
              buffer.WriteString("\",")
              buffer.WriteString("\"RxCode\":\"")
              buffer.WriteString(thePrescription.RxCode)
              buffer.WriteString("\",")
              buffer.WriteString("\"RxInstructions\":\"")
              buffer.WriteString(thePrescription.RxInstructions)
              buffer.WriteString("\",")
              buffer.WriteString("\"FulfillmentAttempts\":")
              buffer.WriteString("[")
                // Here is where we get the history where the prescription fullfillment attempts are made.
                fulfillmentKey := "prescriptionFulfillment" + patientData.PatientId + thePrescription.RxId
                fulfillmentIterator, err := stub.GetHistoryForKey(fulfillmentKey)
                if err != nil {
                  return shim.Error("Error checking for fulfillment: " + err.Error())
                }
                // Boolean to check if we're on the first item or not.
                fulfillmentMemberAlreadyWritten := false
                // Loop through all the fulfillment records.
                for fulfillmentIterator.HasNext() {
                  fulfillment, err := fulfillmentIterator.Next()
                  // If there is an error, throw it.
                  if err != nil {
                    return shim.Error(err.Error())
                  }
                  // Write a comman if we aren't on the first item.
                  if fulfillmentMemberAlreadyWritten == true {
                    buffer.WriteString(",")
                  }
                  var fulfillmentAttempt prescriptionFulfillment
                  err = json.Unmarshal(fulfillment.Value, &fulfillmentAttempt)
                  // Write the value for the string.
                  buffer.WriteString("{\"Fulfilled\":")
                  buffer.WriteString(strconv.FormatBool(fulfillmentAttempt.PrescriptionFilled))
                  // Write the timestamp when the trip was committed.
                  buffer.WriteString(", \"Timestamp\":")
                  buffer.WriteString("\"")
                  buffer.WriteString(time.Unix(fulfillment.Timestamp.Seconds, int64(fulfillment.Timestamp.Nanos)).String())
                  buffer.WriteString("\"")
                  // Close out the item
                  buffer.WriteString("}")
                  fulfillmentMemberAlreadyWritten = true
                }
              buffer.WriteString("]")
            buffer.WriteString("}")
            prescriptionMemberAlreadyWritten = true
          }
          prescriptionIssuesIterator.Close()
        buffer.WriteString("],")
        buffer.WriteString("\"Wearables\": {")
          buffer.WriteString("\"Fitbit\": [")
            fitBitKey := "fitBitData" + patientData.PatientId
            fitBitRecords, err := stub.GetHistoryForKey(fitBitKey)
            if err != nil {
              return shim.Error(err.Error())
            }
            fitBitRecordWrittenBoolean := false
            for fitBitRecords.HasNext() {
              thisFitBitRecord, err := fitBitRecords.Next()
              if err != nil {
                return shim.Error("Error getting FitBit records. " + err.Error())
              }
              // Write a comma if needed.
              if fitBitRecordWrittenBoolean == true {
                buffer.WriteString(",")
              }
              buffer.WriteString("{")
                // Write the timestamp when the trip was committed.
                buffer.WriteString("\"Timestamp\":")
                buffer.WriteString("\"")
                buffer.WriteString(time.Unix(thisFitBitRecord.Timestamp.Seconds, int64(thisFitBitRecord.Timestamp.Nanos)).String())
                buffer.WriteString("\",")
                var thisFitBitRecordJSON fitBitData
                err = json.Unmarshal(thisFitBitRecord.Value, &thisFitBitRecordJSON)
                buffer.WriteString("\"ActivityCalories\":\"")
                buffer.WriteString(thisFitBitRecordJSON.ActivityCalories)
                buffer.WriteString("\",")
                buffer.WriteString("\"CaloriesBMR\":\"")
                buffer.WriteString(thisFitBitRecordJSON.CaloriesBMR)
                buffer.WriteString("\",")
                buffer.WriteString("\"CaloriesOut\":\"")
                buffer.WriteString(thisFitBitRecordJSON.CaloriesOut)
                buffer.WriteString("\",")
                buffer.WriteString("\"Elevation\":\"")
                buffer.WriteString(thisFitBitRecordJSON.Elevation)
                buffer.WriteString("\",")
                buffer.WriteString("\"FairlyActiveMinutes\":\"")
                buffer.WriteString(thisFitBitRecordJSON.FairlyActiveMinutes)
                buffer.WriteString("\",")
                buffer.WriteString("\"Floors\":\"")
                buffer.WriteString(thisFitBitRecordJSON.Floors)
                buffer.WriteString("\",")
                buffer.WriteString("\"RestingHeartRate\":\"")
                buffer.WriteString(thisFitBitRecordJSON.RestingHeartRate)
                buffer.WriteString("\",")
                buffer.WriteString("\"Steps\":\"")
                buffer.WriteString(thisFitBitRecordJSON.Steps)
                buffer.WriteString("\",")
                buffer.WriteString("\"VeryActiveMinutes\":\"")
                buffer.WriteString(thisFitBitRecordJSON.VeryActiveMinutes)
                buffer.WriteString("\",")
                buffer.WriteString("\"MarginalCalories\":\"")
                buffer.WriteString(thisFitBitRecordJSON.MarginalCalories)
                buffer.WriteString("\",")
                buffer.WriteString("\"SedendaryMinutes\":\"")
                buffer.WriteString(thisFitBitRecordJSON.SedendaryMinutes)
                buffer.WriteString("\"")
              buffer.WriteString("}")
              fitBitRecordWrittenBoolean = true
            }
            fitBitRecords.Close()
          buffer.WriteString("],")
          buffer.WriteString("\"BloodPressureMonitor\": [")
            bpmKey := "bloodPressureMonitorData" + patientData.PatientId
            bpmRecords, err := stub.GetHistoryForKey(bpmKey)
            if err != nil {
              return shim.Error(err.Error())
            }
            bpmRecordWrittenBoolean := false
            for bpmRecords.HasNext() {
              bpmRecord, err := bpmRecords.Next()
              if err != nil {
                return shim.Error("Error getting BPM records. " + err.Error())
              }
              // Write a comma if needed.
              if bpmRecordWrittenBoolean == true {
                buffer.WriteString(",")
              }
              buffer.WriteString("{")
                // Write the timestamp when the trip was committed.
                buffer.WriteString("\"Timestamp\":")
                buffer.WriteString("\"")
                buffer.WriteString(time.Unix(bpmRecord.Timestamp.Seconds, int64(bpmRecord.Timestamp.Nanos)).String())
                buffer.WriteString("\",")
                var bpmRecordBytes bloodPressureMonitorData
                err = json.Unmarshal(bpmRecord.Value, &bpmRecordBytes)
                buffer.WriteString("\"SystolicPressure\":\"")
                buffer.WriteString(strconv.Itoa(bpmRecordBytes.SystolicPressure))
                buffer.WriteString("\",")
                buffer.WriteString("\"DiastolicPressure\":\"")
                buffer.WriteString(strconv.Itoa(bpmRecordBytes.DiastolicPressure))
                buffer.WriteString("\",")
                buffer.WriteString("\"Heartrate\":\"")
                buffer.WriteString(strconv.Itoa(bpmRecordBytes.Heartrate))
                buffer.WriteString("\"")
              buffer.WriteString("}")
              bpmRecordWrittenBoolean = true
            }
          buffer.WriteString("]")
        buffer.WriteString("}")
      buffer.WriteString("}")
      patientArrayMemberAlreadyWritten = true
    }

  // Close out the array in the buffer.
  buffer.WriteString("]")

  // We've got the full record, now let's output it.
  return shim.Success(buffer.Bytes())

}
