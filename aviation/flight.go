package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
)

const (
	PEER1          = "PHI"
	PEER2          = "SHELL"
	PEER3          = "CVXG"
	PEER4          = "BP"
	CONTRACT_KEY   = "_Contract"
	INDEX_NAME_FLT = "company~year~month~day"
)

type FlightShrContract struct {
	OwnerCompany  string `json:"ownerCompany"`
	PercSeatAlloc uint8  `json:"percSeatAlloc"`
}

type Flights struct {
	FlightList []Flight `json:"flightList"`
}

type Flight struct {
	FlightKey    string      `json:"flightKey"`
	FlightName   string      `json:"flightName"`
	OwnerCompany string      `json:"ownerCompany"`
	FlightType   string      `json:"flightType"`
	SlNo         string      `json:"slNo"`
	Origin       string      `json:"origin"`
	Destination  string      `json:"destination"`
	DeptDate     string      `json:"deptDate"`
	DeptTime     string      `json:"deptTime"`
	ArrDate      string      `json:"arrDate"`
	ArrTime      string      `json:"arrTime"`
	NoOfSeats    uint8       `json:"noOfSeats"`
	NoOfStops    uint8       `json:"noOfStops"`
	LegDetails   []FlightLeg `json:"legDetails"`
}

type FlightLeg struct {
	Origin      string `json:"origin"`
	Destination string `json:"destination"`
	DeptDate    string `json:"deptDate"`
	DeptTime    string `json:"deptTime"`
	ArrDate     string `json:"arrDate"`
	ArrTime     string `json:"arrTime"`
	TravelMode  string `json:"travelMode"`
	LegNo       uint8  `json:"legNo"`
	AvailSeats  uint8  `json:"availSeats"`
}

// FlightSmartContract implements a simple chaincode to manage an asset
type FlightSmartContract struct {
}

// Init is called during chaincode instantiation to initialize any
// data. Note that chaincode upgrade also calls this function to reset
// or to migrate data.
func (t *FlightSmartContract) Init(stub shim.ChaincodeStubInterface) peer.Response {
	// Get the args from the transaction proposal
	fmt.Println("FlightSmartContract is initializing...")
	args := stub.GetStringArgs()
	if len(args) != 2 {
		return shim.Error("Incorrect arguments. Expecting a key and a value")
	}

	// Set up any variables or assets here by calling stub.PutState()

	// We store the key and the value on the ledger
	err := stub.PutState(args[0], []byte(args[1]))
	if err != nil {
		return shim.Error(fmt.Sprintf("Failed to create asset: %s", args[0]))
	}
	return shim.Success(nil)
}

// Invoke is called per transaction on the chaincode. Each transaction is
// either a 'get' or a 'set' on the asset created by Init function. The Set
// method may create a new asset by specifying a new key-value pair.
func (t *FlightSmartContract) Invoke(stub shim.ChaincodeStubInterface) peer.Response {
	// Extract the function and args from the transaction proposals
	fmt.Println("FlightSmartContract is invoking###...")
	fn, args := stub.GetFunctionAndParameters()

	var result string
	var err error
	if fn == "createFlight" {
		return t.createFlight(stub, args)
	}else if fn == "queryFlight" {
		return t.queryFlight(stub, args)
	}else if fn == "queryAllFlights" {
		return t.queryAllFlights(stub, args);
	}else if fn == "set" {
		return set(stub, args)
	}else { // assume 'get' even if fn is nil
		return get(stub, args)
	}
	if err != nil {
		return shim.Error(err.Error())
	}

	// Return the result as success payload
	return shim.Success([]byte(result))
}

// Set stores the asset (both key and value) on the ledger. If the key exists,
// it will override the value with the new one
func set(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if len(args) != 2 {
		return shim.Error("Incorrect arguments. Expecting a key and a value")
	}

	err := stub.PutState(args[0], []byte(args[1]))
	if err != nil {
		return shim.Error(fmt.Errorf("Failed to set asset: %s", args[0]).Error())
	}
	return shim.Success([]byte(args[1]))
}

// Get returns the value of the specified asset key
func get(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if len(args) != 1 {
		return shim.Error("Incorrect arguments. Expecting a key")
	}

	value, err := stub.GetState(args[0])
	if err != nil {
		return shim.Error(fmt.Errorf("Failed to get asset: %s with error: %s", args[0], err).Error())
	}
	if value == nil {
		return shim.Error(fmt.Errorf("Asset not found: %s", args[0]).Error())
	}
	return shim.Success(value)
}

func (t *FlightSmartContract) createFlight(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	//legDetails1 := []FlightLeg{FlightLeg{Origin: "LOC1", Destination: "LOC2", DeptDate:"19-07-2017", DeptTime:"10:00", ArrDate:"19-07-2017", ArrTime:"11:00", TravelMode: "Fixed Wing", LegNo: 1, AvailSeats: 100},FlightLeg{Origin: "LOC2", Destination: "LOC3", DeptDate:"19-07-2017", DeptTime:"11:10", ArrDate:"19-07-2017", ArrTime:"12:30", TravelMode: "Fixed Wing", LegNo: 1, AvailSeats: 100}}

	//flight := Flight{FlightKey: "Flight#", FlightName: "TEST_FLT", OwnerCompany: "PHI", FlightType: "FTYPE1", SlNo: "SL01", Origin: "LOC1", Destination: "LOC3", DeptDate: "19-07-2017", DeptTime: "10:00", ArrDate: "19-07-2017", ArrTime: "12:30", NoOfSeats: 100, NoOfStops: 1, LegDetails: legDetails1}
	fmt.Println("createFlight is running ")
	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}
	fmt.Println("company ", args[0])
	flight := Flight{}
	flight_json := args[1]
	fmt.Println("flight_json ", flight_json)
	flightByteArray := []byte(flight_json)
	fmt.Println("flightByteArray created...")
	err := json.Unmarshal(flightByteArray, &flight)
	if err != nil {
		fmt.Println("unmarshalling failed... err=", err)
		fmt.Println("Error while parsing file")
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}
	fmt.Println("unmarshalling completed... flight=", flight)
	return createSharedFlights(stub, flight)
}

func createSharedFlights(stub shim.ChaincodeStubInterface, flight Flight) peer.Response {
	fmt.Println("createSharedFlights is running ")
	fltShrContracts := []FlightShrContract{FlightShrContract{OwnerCompany: PEER2, PercSeatAlloc: 20}, FlightShrContract{OwnerCompany: PEER3, PercSeatAlloc: 30}, FlightShrContract{OwnerCompany: PEER4, PercSeatAlloc: 10}}
	totalSeats := flight.NoOfSeats
	availSeat := totalSeats
	var noOfSeats uint8

	i := 0
	fmt.Println("fltShrContracts= ", fltShrContracts)

	compositKey, err := createFlightCompositKey(stub, flight)
	if err != nil {
		return shim.Error(err.Error())
	}
	addFlightToLedger(stub, compositKey, makeFlightsData([]Flight{flight}))

	for i < len(fltShrContracts) {
		fltShrContract := fltShrContracts[i]
		if fltShrContract.PercSeatAlloc > 0 {
			noOfSeats = totalSeats * (fltShrContract.PercSeatAlloc / 100)
			fmt.Println("noOfSeats= ", noOfSeats)
			fmt.Println("availSeat= ", availSeat)
			if availSeat >= noOfSeats {
				newFlight := prepareFlight(flight, noOfSeats, &availSeat, fltShrContract.OwnerCompany)
				compositKey, err := createFlightCompositKey(stub, newFlight)
				if err != nil {
					return shim.Error(err.Error())
				}
				addFlightToLedger(stub, compositKey, makeFlightsData([]Flight{newFlight}))
			}
		}
		i++
	}

	newFlight := prepareFlight(flight, availSeat, &availSeat, flight.OwnerCompany)
	compositKey, err = createFlightCompositKey(stub, newFlight)
	if err != nil {
		return shim.Error(err.Error())
	}
	return addFlightToLedger(stub, compositKey, makeFlightsData([]Flight{newFlight}))
}

func makeFlightsData(flightArr []Flight) Flights {
	flights := Flights{}
	flights.FlightList = flightArr
	return flights
}

func createFlightCompositKey(stub shim.ChaincodeStubInterface, flight Flight) (string, error) {
	tObj, err := time.Parse("02-01-2006", flight.DeptDate) // dd-MM-yyyy
	var key string
	if err != nil {
		return key, errors.New("Invalid date " + flight.DeptDate)
	}
	year, month, day := tObj.Date()
	key, err = stub.CreateCompositeKey(INDEX_NAME_FLT, []string{flight.OwnerCompany, strconv.Itoa(year), strconv.Itoa(int(month)), strconv.Itoa(day)})
	return key, nil
}

func prepareFlight(flight Flight, noOfSeats uint8, availSeat *uint8, ownerCompany string) Flight {
	fmt.Println("prepareFlight is running...")
	newFlight := Flight{}
	*(&newFlight) = *(&flight)
	newFlight.LegDetails = copyLegDetails(flight.LegDetails, noOfSeats)
	newFlight.OwnerCompany = ownerCompany
	newFlight.NoOfSeats = noOfSeats
	*availSeat = *availSeat - noOfSeats
	fmt.Printf("Flight prepared... %+v\n", newFlight)
	return newFlight
}

func copyLegDetails(flightLegs []FlightLeg, noOfSeats uint8) []FlightLeg {
	fmt.Println("copyLegDetails is running...")
	var newFlightLegs []FlightLeg
	var flightLeg FlightLeg
	i := 0
	for i < len(flightLegs) {
		flightLeg = FlightLeg{}
		*(&flightLeg) = *(&flightLegs[i])
		flightLeg.AvailSeats = noOfSeats
		newFlightLegs = append(newFlightLegs, flightLeg)
		i++
	}
	fmt.Printf("created flight legs... %v", newFlightLegs)
	return newFlightLegs
}

func addFlightToLedger(stub shim.ChaincodeStubInterface, key string, flights Flights) peer.Response {
	fmt.Println(">> start writing flight to ledger - Key:", key)
	value, err := stub.GetState(key)
	if err != nil {
		return shim.Error(err.Error())
	}
	if value != nil {
		oldFlights := Flights{}
		err = json.Unmarshal(value, &flights)
		if err != nil {
			fmt.Println("unmarshalling failed... err=", err)
			fmt.Println("Error while parsing file")
			return shim.Error(err.Error())
		}
		flights.FlightList = append(Flights(oldFlights).FlightList, flights.FlightList...)
	}
	flightsAsBytes, _ := json.Marshal(flights)
	stub.PutState(key, flightsAsBytes)
	fmt.Println(">> writing flight to ledger completed- Key:", key)
	return shim.Success(nil)
}

func (t *FlightSmartContract) queryFlight(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	return get(stub, args)
}

func (t *FlightSmartContract) queryAllFlights(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}
	ownerCompany := args[1]
	flightDate := args[2]
	tObj, err := time.Parse("02-01-2006", flightDate) // dd-MM-yyyy
	if err != nil {
		return shim.Error("Invalid date " + flightDate)
	}
	year, month, day := tObj.Date()
	fmt.Printf("%s%d", "...day=", day)
	fltResultsIterator, err := stub.GetStateByPartialCompositeKey(INDEX_NAME_FLT, []string{ownerCompany, strconv.Itoa(year), strconv.Itoa(int(month))})
	if err != nil {
		return shim.Error(err.Error())
	}
	defer fltResultsIterator.Close()
	var i int
	allFlights := Flights{}
	for i = 0; fltResultsIterator.HasNext(); i++ {
		responseRange, err := fltResultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}
		//objectType, compositeKeyParts, err := stub.SplitCompositeKey(responseRange.Key)
		flightsByteArray := responseRange.Value
		flights := Flights{}
		fmt.Println("flightByteArray created...")
		err = json.Unmarshal(flightsByteArray, &flights)
		if err != nil {
			fmt.Println("unmarshalling failed... err=", err)
			fmt.Println("Error while parsing file")
			return shim.Error("Incorrect number of arguments. Expecting 2")
		}
		if i == 0 {
			allFlights = flights
		} else {
			allFlights.FlightList = append(allFlights.FlightList, flights.FlightList...)
		}
	}
	allFlightsAsBytes, _ := json.Marshal(allFlights)
	return shim.Success(allFlightsAsBytes)
}

// main function starts up the chaincode in the container during instantiate
func main() {
	if err := shim.Start(new(FlightSmartContract)); err != nil {
		fmt.Printf("Error starting FlightSmartContract chaincode: %s", err)
	}
}
