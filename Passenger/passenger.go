package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

// Define Passenger Struct which contains all the Passenger Details as seen in the ETIASG1_Passengers Database, Passengers Table
type Passenger struct { //Uppecase to make Struct Public
	FirstName    string `json:"FirstName"`
	LastName     string `json:"LastName"`
	MobileNumber string `json:"MobileNumber"`
	Email        string `json:"Email"`
}

// Define Ride Struct which contains all the Ride Details as seen in the ETIASG1_Passengers Database, PassengerRide Table
type Ride struct {
	PassengerID string `json:"PassengerID"`
	RideDate    string `json:"RideDate"`
	PickupCode  string `json:"PickupCode"`
	DropoffCode string `json:"DropoffCode"`
	DriverID    string `json:"DriverID"`
	CarLicense  string `json:"CarLicense"`
	RideStatus  string `json:"RideStatus"`
}

// Define Diver Struct which contains all the Driver Details as seen in the ETIASG1_Drivers Database, Drivers Table
type Driver struct { // Uppecase to make Struct Public
	FirstName            string `json:"FirstName"`
	LastName             string `json:"LastName"`
	MobileNumber         string `json:"MobileNumber"`
	Email                string `json:"Email"`
	IdentificationNumber string `json:"IdentificationNumber"`
	CarLicense           string `json:"CarLicense"`
	DriverStatus         string `json:"DriverStatus"`
}

// Define the various Global Variables which will be used throughout the code
var (
	passengerValue   string
	passenger_list   = map[string]Passenger{}
	ride_list        = map[string]Ride{}
	driver_ride_list = map[string]Ride{}
)

// The Main Function Handles the API Endpoints along with which functions are associated to each endpoint
func main() {
	router := mux.NewRouter()
	// View Map of Passengers
	router.HandleFunc("/api/v1/passengers", allpassengers)

	// Get/Add/Update Specific Passengers
	router.HandleFunc("/api/v1/passengers/{passengerid}", passenger).Methods("GET", "POST", "PUT")

	// Create a Ride Booking For a Passenger
	router.HandleFunc("/api/v1/passengers/{passengerid}/ride/{pickup}/{dropoff}", passengerbooking)

	// View Passenger Ride History
	router.HandleFunc("/api/v1/passengers/{passengerid}/history", passengerhistory)

	// View Driver Rides With Pending Status
	router.HandleFunc("/api/v1/ride/{driverid}/driver/{ridestatus}", driverridehistory)

	// Update Driver Ride Status
	router.HandleFunc("/api/v1/ride/{rideid}/{driverid}/{response}", driverridestatus)

	fmt.Println("Listening at port 6703")
	log.Fatal(http.ListenAndServe(":6703", router)) //Port 6703 will be used for the Passengers API
}

// Used in "/api/v1/passengers"
func allpassengers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json")
	// Connect to ETIASG1_Passengers DB
	db, err := sql.Open("mysql", "user:password@tcp(127.0.0.1:3306)/ETIASG1_Passengers")
	if err != nil {
		db.Close()
		panic(err.Error())
	}
	defer db.Close()

	// Retrieve All Passengers from Database and Populate Map
	RetreivePassengers(db)

	// Prints out the Passenger List If No Errors Occur
	data, err := json.Marshal(map[string]map[string]Passenger{"Passengers": passenger_list})
	if err != nil {
		log.Fatal(err)
	}
	if len(passenger_list) != 0 {
		fmt.Fprintf(w, "%s\n", data)
	}
}

// Used in "/api/v1/passengers/{passengerid}"
func passenger(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json")
	params := mux.Vars(r)
	// Retrieves the passengerid specified in the API Endpoint which will be used to either Get, Post or Put.
	passengerValue = params["passengerid"]

	// Connect to ETIASG1_Passengers DB
	db, err := sql.Open("mysql", "user:password@tcp(127.0.0.1:3306)/ETIASG1_Passengers")
	if err != nil {
		db.Close()
		panic(err.Error())
	}
	defer db.Close()

	// Retrieve All Passengers from Database and Populate Map
	RetreivePassengers(db)

	// Check if Passenger ID in passengerValue Exists in Map
	passengerDetails, exists := passenger_list[passengerValue]

	if !exists { // Check that Passenger Does Not Exists in DB
		if r.Method == "POST" { // If New Passenger is Being Created
			// Create New Passenger Struct
			newPassenger := Passenger{}
			// Read JSON Data from API Request into Passenger Struct
			reqBody, _ := ioutil.ReadAll(r.Body)
			json.Unmarshal(reqBody, &newPassenger)
			//Insert Into DB (passengerValue Will Be Randomly Generated in the Console Application)
			_, err := db.Exec("INSERT INTO Passengers (PassengerID, FirstName, LastName, MobileNumber, Email) values(?, ?, ?, ?, ?)", passengerValue, newPassenger.FirstName, newPassenger.LastName, newPassenger.MobileNumber, newPassenger.Email)
			if err != nil {
				panic(err.Error())
			}
			//Insert Into Map
			passenger_list[passengerValue] = newPassenger
			w.WriteHeader(http.StatusAccepted)

		} else if r.Method == "PUT" { // If Passenger is Being Updated but Not Found in DB
			// Missing Profile ID Error
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprintf(w, "Error! Your Profile is Missing. Please Contact Our Staff at 91234567.")

		}
	} else if exists { // Check that Passenger Exists in DB
		if r.Method == "GET" { // If Passenger Data is Being Called
			// Print the Passenger Data retrieved from the Map
			data, _ := json.Marshal(passengerDetails)
			fmt.Fprintf(w, "%s\n", data)

		} else if r.Method == "POST" { // If New Passenger is Being Created but ID already exists in the DB
			// Existing Profile Error
			w.WriteHeader(http.StatusConflict)
			fmt.Fprintf(w, "Error! Your Profile already Exists. Please Contact Our Staff at 91234567.")

		} else if r.Method == "PUT" { // If Passenger is Being Updated
			// Create Updated Passenger Struct
			updatePassenger := Passenger{}
			// Read JSON Data from API Request into Passenger Struct
			reqBody, _ := ioutil.ReadAll(r.Body)
			json.Unmarshal(reqBody, &updatePassenger)
			// Update the DB
			_, err := db.Exec("UPDATE Passengers SET FirstName=?, LastName=?, MobileNumber=?, Email=? WHERE PassengerID=?", updatePassenger.FirstName, updatePassenger.LastName, updatePassenger.MobileNumber, updatePassenger.Email, passengerValue)
			if err != nil {
				panic(err.Error())
			}
			// Update Map
			passenger_list[passengerValue] = updatePassenger
			w.WriteHeader(http.StatusAccepted)
		}

	} else { // Error Validation
		w.WriteHeader(http.StatusNotFound) // Error 404
		// Course ID Does Not Exist
		fmt.Fprintf(w, "Oops! You seem to have encountered an error. Please Contact Our Staff at 91234567.")
	}
}

// Used in "/api/v1/passengers/{passengerid}/ride/{pickup}/{dropoff}" to create a ride
func passengerbooking(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json")
	params := mux.Vars(r)
	// Retrieves the passengerid, pickupCode, and dropoffCode specified in the API Endpoint
	pId := params["passengerid"]
	pickupCode := params["pickup"]
	dropoffCode := params["dropoff"]

	// Connect to ETIASG1_Passengers DB
	db, err := sql.Open("mysql", "user:password@tcp(127.0.0.1:3306)/ETIASG1_Passengers")
	if err != nil {
		db.Close()
		panic(err.Error())
	}
	defer db.Close()

	// Update Map with Passenger Data from DB
	RetreivePassengers(db)

	// Check Whether Passenger Account Exists in Map
	_, exists := passenger_list[pId]
	if exists {
		// Get Assigned Driver From Driver API
		assignedDriverMap := map[string]Driver{}
		resp, err := http.Get("http://localhost:3042/api/v1/ride")
		if err != nil {
			log.Fatalln(err)
		}
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatalln(err)
		}
		// Add JSON to Driver Struct
		json.Unmarshal([]byte(body), &assignedDriverMap)

		// Extract Assigned Driver Data
		var assignedDriverID string
		var assignedDriver Driver

		for key, element := range assignedDriverMap {
			assignedDriverID = key
			assignedDriver = element
		}

		// Generate RideID using ddMMyyyy
		generatedRideID := "R" + time.Now().Format("02012006150405")

		// Create New PassengerRide Record
		_, err = db.Exec("INSERT INTO PassengerRides (RideID, RideDate, PassengerID, PickupCode, DropoffCode, DriverID, CarLicense, RideStatus) values(?, ?, ?, ?, ?, ?, ?, ?)", generatedRideID, time.Now(), pId, pickupCode, dropoffCode, assignedDriverID, assignedDriver.CarLicense, "Pending")
		if err != nil {
			panic(err.Error())
		}

		// Print Driver Data to API (Will Just be the Driver Struct)
		data, err := json.Marshal(map[string]Driver{assignedDriverID: assignedDriver})
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(assignedDriver.FirstName)
		fmt.Fprintf(w, "%s\n", data)

	} else {
		w.WriteHeader(http.StatusNotFound) // Error 404
		// Passenger ID Does Not Exist
		fmt.Fprintf(w, "Oops! The Passenger Account is Missing. Please Contact Our Staff at 91234567.")
	}
}

// Used in "/api/v1/passengers/{passengerid}/history"
func passengerhistory(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json")
	params := mux.Vars(r)
	// Retrieves the passengerid which will be used to retrieve the passenger rides
	pId := params["passengerid"]

	// Connect to ETIASG1_Passengers Database
	db, err := sql.Open("mysql", "user:password@tcp(127.0.0.1:3306)/ETIASG1_Passengers")
	if err != nil {
		db.Close()
		panic(err.Error())
	}
	defer db.Close()

	// Empty Ride List to remove previously stored data
	ride_list = map[string]Ride{}

	// Retrieve All Passengers Rides using PassengerID from Database and Populate Map
	RetreivePassengerRides(db, pId)

	// Print Passenger Ride Details in JSON
	data, err := json.Marshal(map[string]map[string]Ride{"Rides": ride_list})
	if err != nil {
		log.Fatal(err)
	}
	if len(ride_list) != 0 {
		fmt.Fprintf(w, "%s\n", data)
	}

}

// Used in "/api/v1/ride/{driverid}/driver/{ridestatus}"
func driverridehistory(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json")
	params := mux.Vars(r)
	// Retrieves the driverid which will be used to retrieve the driver rides
	dId := params["driverid"]
	// Retrieves the ridestatus which will be used to retrieve rides of the specified status
	ridestatus := params["ridestatus"]

	// Connect to ETIASG1_Passengers Database
	db, err := sql.Open("mysql", "user:password@tcp(127.0.0.1:3306)/ETIASG1_Passengers")
	if err != nil {
		db.Close()
		panic(err.Error())
	}
	defer db.Close()

	// Empty Driver Ride List to remove previously stored data
	driver_ride_list = map[string]Ride{}

	// Retrieve All Driver Rides using DriverID and RideStatus from Database and Populate Map
	RetreiveDriverRides(db, dId, ridestatus)

	// Print Driver Ride Details with specified status in JSON
	data, err := json.Marshal(map[string]map[string]Ride{"Rides": driver_ride_list})
	if err != nil {
		log.Fatal(err)
	}
	if len(driver_ride_list) != 0 {
		fmt.Fprintf(w, "%s\n", data)
	}
}

// Used in "/api/v1/ride/{rideid}/{driverid}/{response}"
func driverridestatus(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json")
	params := mux.Vars(r)
	// Retrieves the driverid, rideid, and the response (rideStatus) which will be used to update the ride record
	rID := params["rideid"]
	dID := params["driverid"]
	response := params["response"]

	// Connect to ETIASG1_Passengers Database
	db, err := sql.Open("mysql", "user:password@tcp(127.0.0.1:3306)/ETIASG1_Passengers")
	if err != nil {
		db.Close()
		panic(err.Error())
	}
	defer db.Close()

	// Update DB
	_, err = db.Exec("UPDATE PassengerRides SET RideStatus=? WHERE RideID=?", response, rID)
	if err != nil {
		panic(err.Error())
	}

	// Empty Driver Ride List to remove previously stored data
	driver_ride_list = map[string]Ride{}

	// Retrieve All Driver Rides using DriverID and RideStatus from Database and Populate Map
	RetreiveDriverRides(db, dID, response)

	// Check If Driver Set Ride to Be Completed
	if strings.ToUpper(strings.TrimSpace(response)) == "COMPLETED" {
		// Calls the Driver API to Set Driver Status from Busy to Available
		_, err := http.Get("http://localhost:3042/api/v1/ride/" + dID)
		if err != nil {
			log.Fatalln(err)
		}
	}

	// Print Out Driver Rides in JSON
	data, err := json.Marshal(map[string]map[string]Ride{"Rides": driver_ride_list})
	if err != nil {
		log.Fatal(err)
	}
	if len(driver_ride_list) != 0 {
		fmt.Fprintf(w, "%s\n", data)
	}
}

// This Function Updates the Display Map (passenger_list) with the Passenger Data in the DB (Used in Letting user choose which passenger to login as)
func RetreivePassengers(db *sql.DB) {
	// SQL Query to Retrieve all Passenger Data
	results, err := db.Query("SELECT * FROM Passengers")
	if err != nil {
		db.Close()
		panic(err.Error())
	}
	for results.Next() { // Scan through each Passenger Record from the DB Query to add to a Passenger Struct
		var p Passenger
		var passengerID string
		err = results.Scan(&passengerID, &p.FirstName, &p.LastName, &p.MobileNumber, &p.Email)
		if err != nil {
			db.Close()
			panic(err.Error())
		}

		// Update Passenger List Map with Passengers
		passenger_list[passengerID] = p
	}
}

// This Function Updates the Display Map (ride_list) with the Passenger Ride Data in the DB
func RetreivePassengerRides(db *sql.DB, passengerID string) {
	// SQL Query to Select Completed Rides where PassengerID is the value specified, sorted in reverse chronological order
	results, err := db.Query("SELECT * FROM PassengerRides WHERE PassengerID = ? AND RideStatus = 'Completed' ORDER BY RideDate DESC", passengerID)
	if err != nil {
		db.Close()
		panic(err.Error())
	}
	for results.Next() { // Scan through each Ride Record from the DB Query
		var r Ride
		var rideID string
		err = results.Scan(&rideID, &r.RideDate, &r.PassengerID, &r.PickupCode, &r.DropoffCode, &r.DriverID, &r.CarLicense, &r.RideStatus)
		if err != nil {
			db.Close()
			panic(err.Error())
		}

		// Update Ride List Map with Rides
		ride_list[rideID] = r
	}
}

// This Function Updates the Display Map (driver_ride_list) with the Driver Ride Data in the DB
func RetreiveDriverRides(db *sql.DB, driverID string, ridestatus string) {
	// SQL Query to Select Rides where DriverID and RideStatus are the values specified, sorted in reverse chronological order
	results, err := db.Query("SELECT * FROM PassengerRides WHERE DriverID = ? AND RideStatus = ? ORDER BY RideDate DESC", driverID, ridestatus)
	if err != nil {
		db.Close()
		panic(err.Error())
	}
	for results.Next() { // Scan through each Ride Record from the DB Query
		var r Ride
		var rideID string
		err = results.Scan(&rideID, &r.RideDate, &r.PassengerID, &r.PickupCode, &r.DropoffCode, &r.DriverID, &r.CarLicense, &r.RideStatus)
		if err != nil {
			db.Close()
			panic(err.Error())
		}

		// Update Driver Ride List Map with Rides
		driver_ride_list[rideID] = r
	}
}
