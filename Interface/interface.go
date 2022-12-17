package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

// Create Loading Bar String Array which will be used in Booking Ride
var loadingBar = []string{
	"00%: [                    ]",
	"20%: [####                ]",
	"40%: [########            ]",
	"60%: [############        ]",
	"80%: [################    ]",
	"100%:[####################]\n",
}

// Define the various Global Variables which will be used throughout the code
var (
	passenger_list map[string]map[string]Passenger
	driver_list    map[string]map[string]Driver
)

// Define Passenger Struct which contains all the Passenger Details as seen in the ETIASG1_Passengers Database, Passengers Table
type Passenger struct { // Uppecase to make Struct Public
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
type Driver struct {
	FirstName            string `json:"FirstName"`
	LastName             string `json:"LastName"`
	MobileNumber         string `json:"MobileNumber"`
	Email                string `json:"Email"`
	IdentificationNumber string `json:"IdentificationNumber"`
	CarLicense           string `json:"CarLicense"`
}

func main() {
outer:
	for {
		// Print Console Options
		fmt.Println(strings.Repeat("=", 10))
		fmt.Println("KTRide Platform")
		fmt.Println(strings.Repeat("=", 10),
			"\n [1] Login\n",
			"[2] Create an Account\n",
			"[9] Quit")

		// Read User Input into outerOption
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("\nEnter an option: ")
		userInput, _ := reader.ReadString('\n')
		outerOption, _ := strconv.Atoi(strings.TrimSpace(userInput))

		// Option Actions
		switch outerOption {
		case 1:
			Login()

		case 2:
			Create()

		case 9:
			break outer

		default: // If Any Other Key Was Pressed
			fmt.Println("### Invalid Input ###")
		}
	}
}

// Login Function For Passengers and Drivers
func Login() {
login:
	for {
		// Print Login Console Options
		fmt.Println(strings.Repeat("=", 10))
		fmt.Println("How Would You Like to Use The Service?\n",
			"[1] Passenger\n",
			"[2] Driver\n",
			"[9] Quit")

		// Read User Input into innerOption
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("\nEnter an option: ")
		userInput, _ := reader.ReadString('\n')
		innerOption, _ := strconv.Atoi(strings.TrimSpace(userInput))

		// Option Actions
		switch innerOption {
		case 1:
			// Retrieve All Passengers from DB
			resp, err := http.Get("http://localhost:6703/api/v1/passengers")
			if err != nil {
				fmt.Println("Error! The Program Was Unable to Retrieve the Passenger Data. Do Try Again!")
				continue
			}
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				log.Fatalln(err)
			}
			json.Unmarshal([]byte(body), &passenger_list)

			// Check If Program Retrieved Passenger List with No Accounts
			if len(passenger_list) != 0 {
				// Prompt User for Account Selection
				fmt.Println(strings.Repeat("=", 10))
				fmt.Println("Which Passenger Would You Like to Login As?")

				// Create Numeric ID Map
				idMap := map[string]string{}
				count := 1

				// Print Login Options
				for key, element := range passenger_list["Passengers"] {
					countString := strconv.Itoa(count)
					idMap[countString] = key
					fmt.Println(" " + countString + ": " + element.FirstName + " " + element.LastName + " [" + element.MobileNumber + "/" + element.Email + "]")
					count += 1
				}

				// Read User Input & Retrieve User Details
				reader := bufio.NewReader(os.Stdin)
				fmt.Print("\nEnter Login Number: ")
				userNumber, _ := reader.ReadString('\n')
				userNumber = strings.TrimSpace(userNumber)
				userID := idMap[userNumber]

				// Check if User is in DB
				loggedInPassenger, exists := passenger_list["Passengers"][userID]
				if !exists {
					fmt.Println("\nSorry, the UserID Inputted Does Not Exist. Press Any Key to Try Again...")
					fmt.Scanln() // Wait for Keypress
					continue
				} else {
					fmt.Print("\nSuccessfully Logged In As Passenger [" + loggedInPassenger.FirstName + " " + loggedInPassenger.LastName + "]!\n")
					PassengerActions(userID, loggedInPassenger)
				}

			} else { // If Map is Empty (Should Not Happen as Sample Accounts Inserted during Creation)
				fmt.Println("There are Currently No Passenger Accounts! Do Create an Account or Login as a Driver.")
				continue
			}

		case 2:
			// Retrieve All Drivers from DB
			resp, err := http.Get("http://localhost:3042/api/v1/drivers")
			if err != nil {
				fmt.Println("Error! The Program Was Unable to Retrieve the Driver Data. Do Try Again!")
				continue
			}
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				log.Fatalln(err)
			}
			json.Unmarshal([]byte(body), &driver_list)

			// Check If Program Retrieved Driver List with No Accounts
			if len(driver_list) != 0 {
				// Prompt User for Account Selection
				fmt.Println(strings.Repeat("=", 10))
				fmt.Println("Which Driver Would You Like to Login As?")

				// Create Numeric ID Map
				idMap := map[string]string{}
				count := 1

				// Print Login Options
				for key, element := range driver_list["Drivers"] {
					countString := strconv.Itoa(count)
					idMap[countString] = key
					fmt.Println(" " + countString + ": " + element.FirstName + " " + element.LastName + " - " + element.CarLicense + " [" + element.MobileNumber + "/" + element.Email + "]")
					count += 1
				}

				// Read User Input & Retrieve User Details
				reader := bufio.NewReader(os.Stdin)
				fmt.Print("\nEnter Login Number: ")
				userNumber, _ := reader.ReadString('\n')
				userNumber = strings.TrimSpace(userNumber)
				userID := idMap[userNumber]

				// Check if User is in DB
				loggedInDriver, exists := driver_list["Drivers"][userID]
				if !exists {
					fmt.Println("\nSorry, the UserID Inputted Does Not Exist. Press Any Key to Try Again...")
					fmt.Scanln() // Wait for Keypress
					continue
				} else {
					fmt.Print("\nSuccessfully Logged In As Driver [" + loggedInDriver.FirstName + " " + loggedInDriver.LastName + "]!\n")
					DriverActions(userID, loggedInDriver)
				}

			} else { // If Map is Empty (Should Not Happen as Sample Accounts Inserted during Creation)
				fmt.Println("There are Currently No Driver Accounts! Do Create an Account or Login as a Passenger.")
				continue
			}

		case 9:
			break login

		default:
			fmt.Println("### Invalid Input ###")
		}
	}
}

func Create() {
create:
	for {
		// Print Console Options
		fmt.Println(strings.Repeat("=", 10))
		fmt.Println("Which Account Would You Like To Create?\n",
			"[1] Passenger\n",
			"[2] Driver\n",
			"[9] Quit")

		// Read User Input into innerOption
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("\nEnter an option: ")
		userInput, _ := reader.ReadString('\n')
		innerOption, _ := strconv.Atoi(strings.TrimSpace(userInput))

		// Option Actions
		switch innerOption {
		case 1:
			newPassenger := Passenger{}

			// Passenger First Name
			fmt.Println(strings.Repeat("=", 10))
			reader = bufio.NewReader(os.Stdin)
			fmt.Println("\nEnter First Name: ")
			inputString, _ := reader.ReadString('\n')
			inputString = strings.TrimSpace(inputString)
			newPassenger.FirstName = inputString

			// Passenger Last Name
			reader = bufio.NewReader(os.Stdin)
			fmt.Println("\nEnter Last Name: ")
			inputString, _ = reader.ReadString('\n')
			inputString = strings.TrimSpace(inputString)
			newPassenger.LastName = inputString

			// Passenger Mobile Number
			reader = bufio.NewReader(os.Stdin)
			fmt.Println("\nEnter Mobile Number: ")
			inputString, _ = reader.ReadString('\n')
			inputString = strings.TrimSpace(inputString)
			newPassenger.MobileNumber = inputString

			// Passenger Email
			reader = bufio.NewReader(os.Stdin)
			fmt.Println("\nEnter Email: ")
			inputString, _ = reader.ReadString('\n')
			inputString = strings.TrimSpace(inputString)
			newPassenger.Email = inputString

			// Convert the Passenger Details into JSON
			jsonBody, _ := json.Marshal(newPassenger)

			// Generate Passenger ID
			userID := "P" + time.Now().Format("02012006150405")

			// Post Passenger Data to passengers.go API to create new Passenger in the Database
			client := &http.Client{}
			if req, err := http.NewRequest("POST", "http://localhost:6703/api/v1/passengers/"+userID, bytes.NewBuffer(jsonBody)); err == nil {
				if _, err := client.Do(req); err == nil {
					fmt.Print("\nPassenger ", newPassenger.FirstName+" "+newPassenger.LastName, " Created Successfully!\n")
					break create
				} else {
					fmt.Println("Error! Your Passenger Account Could Not Be Created. Please Try Again.")
					continue
				}
			}

		case 2:
			newDriver := Driver{}

			// Driver First Name
			fmt.Println(strings.Repeat("=", 10))
			reader = bufio.NewReader(os.Stdin)
			fmt.Println("Enter First Name: ")
			inputString, _ := reader.ReadString('\n')
			inputString = strings.TrimSpace(inputString)
			newDriver.FirstName = inputString

			// Driver Last Name
			reader = bufio.NewReader(os.Stdin)
			fmt.Println("\nEnter Last Name: ")
			inputString, _ = reader.ReadString('\n')
			inputString = strings.TrimSpace(inputString)
			newDriver.LastName = inputString

			// Driver Mobile Number
			reader = bufio.NewReader(os.Stdin)
			fmt.Println("\nEnter Mobile Number: ")
			inputString, _ = reader.ReadString('\n')
			inputString = strings.TrimSpace(inputString)
			newDriver.MobileNumber = inputString

			// Driver Email
			reader = bufio.NewReader(os.Stdin)
			fmt.Println("\nEnter Email: ")
			inputString, _ = reader.ReadString('\n')
			inputString = strings.TrimSpace(inputString)
			newDriver.Email = inputString

			// Driver Car License
			reader = bufio.NewReader(os.Stdin)
			fmt.Println("\nEnter Car License: ")
			inputString, _ = reader.ReadString('\n')
			inputString = strings.TrimSpace(inputString)
			newDriver.CarLicense = inputString

			// Convert the Driver Details into JSON
			jsonBody, _ := json.Marshal(newDriver)

			// Generate Driver ID
			userID := "D" + time.Now().Format("02012006150405")

			// Post Driver Data to drivers.go API to create new Driver in the Database
			client := &http.Client{}
			if req, err := http.NewRequest("POST", "http://localhost:3042/api/v1/drivers/"+userID, bytes.NewBuffer(jsonBody)); err == nil {
				if _, err := client.Do(req); err == nil {
					fmt.Print("\nDriver ", newDriver.FirstName+" "+newDriver.LastName, " Created Successfully!\n")
					break create
				} else {
					fmt.Println("Error! Your Driver Account Could Not Be Created. Please Try Again.")
					continue
				}
			}

		case 9:
			break create

		default:
			fmt.Println("### Invalid Input ###")
		}
	}
}

func PassengerActions(userID string, p Passenger) {
passengerActions:
	for {
		// Print Console Options
		fmt.Println(strings.Repeat("=", 10))
		fmt.Println("What Would You Like To Do?\n",
			"[1] Book a Ride\n",
			"[2] View Ride History\n",
			"[3] Change Account Details\n",
			"[9] Logout")

		// Read User Input into innerOption
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("\nEnter an option: ")
		userInput, _ := reader.ReadString('\n')
		innerOption, _ := strconv.Atoi(strings.TrimSpace(userInput))
		fmt.Println(strings.Repeat("=", 10))

		// Option Actions
		switch innerOption {
		case 1:
			// Ride Pick-Up Postal Code
			reader := bufio.NewReader(os.Stdin)
			fmt.Println("Input the Pick-Up Postal Code: ")
			inputPickUp, _ := reader.ReadString('\n')
			inputPickUp = strings.TrimSpace(inputPickUp)

			// Ride Drop-Off Postal Code
			reader = bufio.NewReader(os.Stdin)
			fmt.Println("\nInput the Drop-Off Postal Code: ")
			inputDropOff, _ := reader.ReadString('\n')
			inputDropOff = strings.TrimSpace(inputDropOff)

			// Get Random Available Driver Data from drivers.go API
			var assignedDriver map[string]Driver
			resp, err := http.Get("http://localhost:6703/api/v1/passengers/" + userID + "/ride/" + inputPickUp + "/" + inputDropOff)
			if err != nil {
				fmt.Println("Error! The Program Was Unable to Find a Driver. Do Try Again!")
				continue
			}
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				log.Fatalln(err)
			}
			json.Unmarshal([]byte(body), &assignedDriver)

			// Print Driver Search Progress Bar
			fmt.Println(strings.Repeat("=", 10))
			fmt.Println("Searching For a Driver...\n ")
			for _, progress := range loadingBar {
				fmt.Printf("\r \a%s", progress)
				time.Sleep(1 * time.Second)
			}

			// Check That Random Driver Has Been Assigned
			if len(assignedDriver) != 0 {
				// Print Assigned Driver Details
				fmt.Println("\nSearch Complete!\n\n==========\nRide Details\n==========")
				for _, element := range assignedDriver {
					fmt.Println("Driver Name: " + element.FirstName + " " + element.LastName + "\nMobile Number: " + element.MobileNumber + "\nEmail: " + element.Email + "\nCar License: " + element.CarLicense + "\n\nPick-Up Postal Code: " + inputPickUp + "\nDrop-Off Postal Code: " + inputDropOff)
				}
				fmt.Println("\nYour Driver will Arrive Shortly. Press any Key to Continue...")
				fmt.Scanln() // Wait for Keypress

			} else { // If All Drivers have their Status Set to Busy
				fmt.Println("\nSorry, All The Drivers In Your Area Are Busy. Please Try Again Shortly.")
				continue
			}

		case 2:
			// Get Passenger Ride Data from passengers.go API
			resp, err := http.Get("http://localhost:6703/api/v1/passengers/" + userID + "/history")
			if err != nil {
				fmt.Println("Error! The Program Was Unable to Retrieve Your History. Do Try Again!")
				continue
			}
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				log.Fatalln(err)
			}
			var ride_history map[string]map[string]Ride
			json.Unmarshal([]byte(body), &ride_history)

			// Check If Passenger Has Rides In Their History
			if len(ride_history) != 0 {
				fmt.Println("Ride History: " + p.FirstName + " " + p.LastName)
				fmt.Println(strings.Repeat("=", 10))

				// Print Driver Details
				for _, element := range ride_history["Rides"] {
					fmt.Println("Ride Date: " + element.RideDate + "\nPick-Up Postal Code: " + element.PickupCode + "\nDrop-Off Postal Code: " + element.DropoffCode + "\nCar License: " + element.CarLicense + "\n==========")
				}
				fmt.Println("\nPress any Key to Continue...")
				fmt.Scanln() // Wait for Keypress

			} else { // If Passenger Has No Rides
				fmt.Println("\nYou Have No Rides In Your Ride History!")
			}

		case 3:
			updatedPassenger := Passenger{}

			// Passenger First Name
			reader = bufio.NewReader(os.Stdin)
			fmt.Println("Enter First Name: ")
			inputString, _ := reader.ReadString('\n')
			inputString = strings.TrimSpace(inputString)
			updatedPassenger.FirstName = inputString

			// Passenger Last Name
			reader = bufio.NewReader(os.Stdin)
			fmt.Println("Enter Last Name: ")
			inputString, _ = reader.ReadString('\n')
			inputString = strings.TrimSpace(inputString)
			updatedPassenger.LastName = inputString

			// Passenger Mobile Number
			reader = bufio.NewReader(os.Stdin)
			fmt.Println("Enter Mobile Number: ")
			inputString, _ = reader.ReadString('\n')
			inputString = strings.TrimSpace(inputString)
			updatedPassenger.MobileNumber = inputString

			// Passenger Email
			reader = bufio.NewReader(os.Stdin)
			fmt.Println("Enter Email: ")
			inputString, _ = reader.ReadString('\n')
			inputString = strings.TrimSpace(inputString)
			updatedPassenger.Email = inputString

			jsonBody, _ := json.Marshal(updatedPassenger)

			// Put Passenger Data to passengers.go API
			client := &http.Client{}
			if req, err := http.NewRequest("PUT", "http://localhost:6703/api/v1/passengers/"+userID, bytes.NewBuffer(jsonBody)); err == nil {
				if _, err := client.Do(req); err == nil {
					fmt.Print("\nPassenger ", updatedPassenger.FirstName+" "+updatedPassenger.LastName, " Updated Successfully!\n")
				}
			}

		case 9:
			break passengerActions

		default:
			fmt.Println("### Invalid Input ###")
		}

	}
}

func DriverActions(userID string, d Driver) {
driverActions:
	for {
		// Print Console Options
		fmt.Println(strings.Repeat("=", 10))
		fmt.Println("What Would You Like To Do?\n",
			"[1] Start a Trip\n",
			"[2] End a Trip\n",
			"[3] Change Account Details\n",
			"[9] Logout")

		// Read User Input into innerOption
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("\nEnter an option: ")
		userInput, _ := reader.ReadString('\n')
		innerOption, _ := strconv.Atoi(strings.TrimSpace(userInput))

		// Option Actions
		switch innerOption {
		case 1:
			// Get Driver Ride Data from passengers.go API (passengers.go API is will call drivers.go API)
			resp, err := http.Get("http://localhost:6703/api/v1/ride/" + userID + "/driver/Pending")
			if err != nil {
				fmt.Println("Error! The Program Was Unable to Retrieve the Driver Ride Data. Do Try Again!")
			}
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				log.Fatalln(err)
			}
			driver_ride_list := map[string]map[string]Ride{}
			json.Unmarshal([]byte(body), &driver_ride_list)

			// Check If Driver Has Any Rides To Start Available
			if len(driver_ride_list["Rides"]) != 0 {
				fmt.Println(strings.Repeat("=", 10))
				fmt.Println("Ride Pending")
				fmt.Println(strings.Repeat("=", 10))

				// Create Numeric ID Map
				var pendingRideID string

				// Print Out Pending Ride Details
				for key, element := range driver_ride_list["Rides"] {
					pendingRideID = key
					fmt.Println("\nRide Request Date: " + element.RideDate + "\nPick-Up Postal Code: " + element.PickupCode + "\nDrop-Off Postal Code: " + element.DropoffCode + "\nCar License: " + element.CarLicense + "\n")
				}

				reader = bufio.NewReader(os.Stdin)
				fmt.Println("Would You Like To Start The Trip? (Yes/No): ")
				userInput, _ := reader.ReadString('\n')
				userInput = strings.ToUpper(strings.TrimSpace(userInput))

				// Check If Driver Started Trip
				if userInput == "YES" || userInput == "Y" {
					// Start Ride
					resp, err = http.Get("http://localhost:6703/api/v1/ride/" + pendingRideID + "/" + userID + "/Ongoing")
					if err != nil {
						fmt.Println("Error! The Program Was Unable to Start the Ride. Do Try Again!")
					}

					body, err = ioutil.ReadAll(resp.Body)
					if err != nil {
						log.Fatalln(err)
					}

					driver_ride_list = map[string]map[string]Ride{}
					json.Unmarshal([]byte(body), &driver_ride_list)

					if len(driver_ride_list["Rides"]) != 0 {
						for key, element := range driver_ride_list["Rides"] {
							fmt.Println("\nRide " + key + " From " + element.PickupCode + " To " + element.DropoffCode + " Started!")
						}
					}
				} else { // If Driver Does Not Start Trip
					fmt.Println("\nDo Start the Trip Soon or Contact 91234567 to Cancel any Pending Trips.")
				}

			} else { // If Driver Has No Pending Trips
				fmt.Println("\nThere Are No Pending Trips To Start!\n ")
				continue
			}

		case 2:
			// Get Driver Ride Data from passengers.go API (passengers.go API is will call drivers.go API)
			resp, err := http.Get("http://localhost:6703/api/v1/ride/" + userID + "/driver/Ongoing")
			if err != nil {
				fmt.Println("Error! The Program Was Unable to Retrieve the Driver Ride Data. Do Try Again!")
			}
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				log.Fatalln(err)
			}
			driver_ride_list := map[string]map[string]Ride{}
			json.Unmarshal([]byte(body), &driver_ride_list)

			// Check If Driver Has Any Rides To Start Available
			if len(driver_ride_list["Rides"]) != 0 {
				fmt.Println(strings.Repeat("=", 10))
				fmt.Println("Ride Ongoing")
				fmt.Println(strings.Repeat("=", 10))

				// Create Numeric ID Map
				var pendingRideID string

				// Print Out Ongoing Ride Details
				for key, element := range driver_ride_list["Rides"] {
					pendingRideID = key
					fmt.Println("\nRide Start Date: " + element.RideDate + "\nPick-Up Postal Code: " + element.PickupCode + "\nDrop-Off Postal Code: " + element.DropoffCode + "\nCar License: " + element.CarLicense + "\n")
				}

				reader = bufio.NewReader(os.Stdin)
				fmt.Println("Would You Like To End The Trip? (Yes/No): ")
				userInput, _ := reader.ReadString('\n')
				userInput = strings.ToUpper(strings.TrimSpace(userInput))

				// Check If Driver Ended Ride
				if userInput == "YES" || userInput == "Y" {
					// End Ride and set DriverStatus to 'Available'
					resp, err = http.Get("http://localhost:6703/api/v1/ride/" + pendingRideID + "/" + userID + "/Completed")
					if err != nil {
						fmt.Println("Error! The Program Was Unable to End the Ride. Do Try Again!")
					}

					body, err = ioutil.ReadAll(resp.Body)
					if err != nil {
						log.Fatalln(err)
					}

					// Empty the driver_ride_list map so no previous data remains
					driver_ride_list = map[string]map[string]Ride{}

					// Unmarshal the Compelted Ride Details into driver_ride_list
					json.Unmarshal([]byte(body), &driver_ride_list)

					// Print Ride End Message
					if len(driver_ride_list["Rides"]) != 0 {
						for key, element := range driver_ride_list["Rides"] {
							fmt.Println("\nRide " + key + " From " + element.PickupCode + " To " + element.DropoffCode + " Ended!")
						}
					}
				} else { // If Driver Does Not Start Trip
					fmt.Println("\nDo End the Trip Soon or Contact 91234567 for Help.")
				}

			} else { // If Driver Has No Pending Trips
				fmt.Println("\nThere Are No Pending Trips To End!\n ")
				continue
			}

		case 3:
			updatedDriver := Driver{}

			// Driver First Name
			reader = bufio.NewReader(os.Stdin)
			fmt.Println("Enter First Name: ")
			inputString, _ := reader.ReadString('\n')
			inputString = strings.TrimSpace(inputString)
			updatedDriver.FirstName = inputString

			// Driver Last Name
			reader = bufio.NewReader(os.Stdin)
			fmt.Println("Enter Last Name: ")
			inputString, _ = reader.ReadString('\n')
			inputString = strings.TrimSpace(inputString)
			updatedDriver.LastName = inputString

			// Driver Mobile Number
			reader = bufio.NewReader(os.Stdin)
			fmt.Println("Enter Mobile Number: ")
			inputString, _ = reader.ReadString('\n')
			inputString = strings.TrimSpace(inputString)
			updatedDriver.MobileNumber = inputString

			// Driver Email
			reader = bufio.NewReader(os.Stdin)
			fmt.Println("Enter Email: ")
			inputString, _ = reader.ReadString('\n')
			inputString = strings.TrimSpace(inputString)
			updatedDriver.Email = inputString

			// Driver Car License
			reader = bufio.NewReader(os.Stdin)
			fmt.Println("Enter Car License: ")
			inputString, _ = reader.ReadString('\n')
			inputString = strings.TrimSpace(inputString)
			updatedDriver.CarLicense = inputString

			// Convert updatedDriver struct to JSON to be passed into the API POST
			jsonBody, _ := json.Marshal(updatedDriver)

			// Post the Driver JSON to add the driver to the database
			client := &http.Client{}
			if req, err := http.NewRequest("PUT", "http://localhost:3042/api/v1/drivers/"+userID, bytes.NewBuffer(jsonBody)); err == nil {
				if _, err := client.Do(req); err == nil {
					fmt.Print("\nDriver ", updatedDriver.FirstName+" "+updatedDriver.LastName, " Updated Successfully!\n")
				}
			}

		case 9:
			break driverActions

		default:
			fmt.Println("### Invalid Input ###")
		}
	}
}
