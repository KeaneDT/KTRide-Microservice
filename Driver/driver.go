package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

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
	driverValue string
	driver_list = map[string]Driver{}
)

// The Main Function Handles the API Endpoints along with which functions are associated to each endpoint
func main() {
	router := mux.NewRouter()
	// View Map of Drivers, function returns map of all drivers.
	router.HandleFunc("/api/v1/drivers", alldrivers)

	// Get/Add/Update Specific Driver, function either retrieves, creates or updates a driver based on method.
	router.HandleFunc("/api/v1/drivers/{driverid}", driver).Methods("GET", "POST", "PUT")

	// Assign Random Available Driver & Set Their Status To Busy in the Database
	router.HandleFunc("/api/v1/ride", driverbooking)

	// Update Specified Driver Status To Available in the Database
	router.HandleFunc("/api/v1/ride/{driverid}", drivercomplete)

	fmt.Println("Listening at port 3042")
	log.Fatal(http.ListenAndServe(":3042", router)) //Port 3042 will be used for the Passengers API
}

// Used in "/api/v1/drivers"
func alldrivers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json")
	db, err := sql.Open("mysql", "user:password@tcp(127.0.0.1:3306)/ETIASG1_Drivers")
	if err != nil {
		db.Close()
		panic(err.Error())
	}
	defer db.Close()

	// Retrieve All Drivers from Database and Populate Map
	RetreiveDrivers(db)

	// Prints out the Driver List If No Errors Occur
	data, err := json.Marshal(map[string]map[string]Driver{"Drivers": driver_list})
	if err != nil {
		log.Fatal(err)
	}
	if len(driver_list) != 0 {
		fmt.Fprintf(w, "%s\n", data)
	}
}

// Used in "/api/v1/drivers/{driverid}"
func driver(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json")

	// Connect to the ETIASG1 Database
	db, err := sql.Open("mysql", "user:password@tcp(127.0.0.1:3306)/ETIASG1_Drivers")
	if err != nil {
		db.Close()
		panic(err.Error())
	}
	defer db.Close()

	// Retrieve the List of Drivers in DB & Update Map
	RetreiveDrivers(db)

	// Get Params & Check if Driver Exists
	params := mux.Vars(r)
	// Retrieves the driverid specified in the API Endpoint which will be used to either Get, Post or Put.
	driverValue = params["driverid"]

	// Check if Driver ID in passengerValue Exists in Map
	driverDetails, exists := driver_list[driverValue]

	inputDriver := Driver{}
	reqBody, _ := ioutil.ReadAll(r.Body)
	json.Unmarshal(reqBody, &inputDriver)

	fmt.Printf(driverValue, inputDriver.FirstName, inputDriver.LastName, inputDriver.MobileNumber, inputDriver.Email, inputDriver.IdentificationNumber, inputDriver.CarLicense)

	if !exists { // Check If Driver Details Does Not Exist in Map
		if r.Method == "POST" { // If New Driver is Being Created
			// Insert Into DB (DriverValue Will Be Randomly Generated in the Console Application)
			_, err := db.Exec("INSERT INTO Drivers (DriverID, FirstName, LastName, MobileNumber, Email, IdentificationNumber, CarLicense, DriverStatus) values(?, ?, ?, ?, ?, ?, ?, 'Available')", driverValue, inputDriver.FirstName, inputDriver.LastName, inputDriver.MobileNumber, inputDriver.Email, inputDriver.IdentificationNumber, inputDriver.CarLicense)
			if err != nil {
				panic(err.Error())
			}

			// Insert Into Map
			driver_list[driverValue] = inputDriver

			// Write StatusAccepted Header
			w.WriteHeader(http.StatusAccepted)
			fmt.Fprintf(w, "The Driver Profile "+inputDriver.FirstName+" "+inputDriver.LastName+" ("+inputDriver.MobileNumber+") Has Been Created!")

		} else if r.Method == "PUT" { // If Driver is Being Updated but Not Found in DB
			// Missing Profile ID Error (Cannot Update a Missing Profile)
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprintf(w, "Error! Your Profile is Missing and Can Not Be Updated. Please Contact Our Staff at 91234567.")

		}

	} else if exists { // Check If Driver Details Exists in DB
		if r.Method == "GET" { // If Driver Data is Being Called
			// Print the specified driver details retrieved from the map
			data, _ := json.Marshal(driverDetails)
			fmt.Fprintf(w, "%s\n", data)

		} else if r.Method == "POST" { // If New Driver is Being Created but ID already exists in the DB
			// Existing Profile Error
			w.WriteHeader(http.StatusConflict)
			fmt.Fprintf(w, "Error! Your Mobile Number or Identification Number Is Already in Use. Please Contact Our Staff at 91234567.")

		} else if r.Method == "PUT" { // If Driver is Being Updated
			// Update the DB
			_, err := db.Exec("UPDATE Drivers SET FirstName=?, LastName=?, MobileNumber=?, Email=?, CarLicense=? WHERE DriverID=?", inputDriver.FirstName, inputDriver.LastName, inputDriver.MobileNumber, inputDriver.Email, inputDriver.CarLicense, driverValue)
			if err != nil {
				panic(err.Error())
			}
			// Update Map
			driver_list[driverValue] = inputDriver

			w.WriteHeader(http.StatusAccepted)
			fmt.Fprintf(w, "Profile Update Successful!")
		}

	} else { // Error Validation
		w.WriteHeader(http.StatusNotFound) // Error 404
		// Driver ID Does Not Exist
		fmt.Fprintf(w, "Oops! You seem to have encountered an error. Please Contact Our Staff at 91234567.")
	}
}

// Used in "/api/v1/ride"
func driverbooking(w http.ResponseWriter, r *http.Request) {
	// Set Up Function
	w.Header().Set("Content-type", "application/json")

	// Connect to ETIASG1_Drivers Database
	db, err := sql.Open("mysql", "user:password@tcp(127.0.0.1:3306)/ETIASG1_Drivers")
	if err != nil {
		db.Close()
		panic(err.Error())
	}
	defer db.Close()

	// Retrieve 1 random available driver from the database
	results, err := db.Query("SELECT * FROM Drivers WHERE DriverStatus = 'Available' ORDER BY RAND() LIMIT 1")
	if err != nil {
		db.Close()
		panic(err.Error())
	}
	for results.Next() { // Add the Random Driver Details to a Driver Structc
		var randDriver Driver
		var driverID string
		err = results.Scan(&driverID, &randDriver.FirstName, &randDriver.LastName, &randDriver.MobileNumber, &randDriver.Email, &randDriver.IdentificationNumber, &randDriver.CarLicense, &randDriver.DriverStatus)
		if err != nil {
			db.Close()
			panic(err.Error())
		}

		// Set Driver Status to Busy
		_, err = db.Exec("UPDATE Drivers SET DriverStatus='Busy' WHERE DriverID=?", driverID)
		if err != nil {
			panic(err.Error())
		}
		randDriver.DriverStatus = "Busy"

		// Print Driver Data to API (Will Just be the Driver Struct)
		data, err := json.Marshal(map[string]Driver{driverID: randDriver})
		if err != nil {
			log.Fatal(err)
		}
		fmt.Fprintf(w, "%s\n", data)
	}
}

// Used in "/api/v1/ride/{driverid}"
func drivercomplete(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json")
	params := mux.Vars(r)
	// Retrieves the driverid specified in the API Endpoint which will have its status set to Available
	driverValue = params["driverid"]

	// Connect to the ETIASG1_Drivers Database
	db, err := sql.Open("mysql", "user:password@tcp(127.0.0.1:3306)/ETIASG1_Drivers")
	if err != nil {
		db.Close()
		panic(err.Error())
	}
	defer db.Close()

	// Set Driver Status to Busy
	_, err = db.Exec("UPDATE Drivers SET DriverStatus='Available' WHERE DriverID=?", driverValue)
	if err != nil {
		panic(err.Error())
	}
}

// This Function Updates the Display Map (driver_list) with the Driver Data in the DB (Used in Letting user choose which driver to login as)
func RetreiveDrivers(db *sql.DB) {
	// SQL Query to Retrieve all Driver Data
	results, err := db.Query("SELECT * FROM Drivers")
	if err != nil {
		db.Close()
		panic(err.Error())
	}
	for results.Next() { // Scan through each Driver Record from the DB Query to add to a Driver Struct
		var d Driver
		var driverID string
		err = results.Scan(&driverID, &d.FirstName, &d.LastName, &d.MobileNumber, &d.Email, &d.IdentificationNumber, &d.CarLicense, &d.DriverStatus)
		if err != nil {
			db.Close()
			panic(err.Error())
		}

		// Update Driver List Map with Drivers
		driver_list[driverID] = d
	}
}
