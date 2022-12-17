# KTRide-Microservices
3 Microservices, Passenger, Driver, and Interface created for a Ride-Sharing Platform. All the files are created using GoLang.

## User Guide for KTRide Service Console Application

This console application allows users to login as either passengers or drivers and use the ride service. Passengers can book rides, view their ride history, and cancel rides. Drivers can view their pending rides and accept or decline them.

The first step of setting up the Ride Sharing Application is to either download or clone the files in the repository to get the 3 microservice files which are the Passenger API, the Driver API, and the Interface. From there, run the 'ETIASG1_Sript.sql' Database set-up script to set up the necessary databases and tables for the APIs to use for persistent storage. 

In order to properly run this application, both the Passenger API (passenger.go) and Driver API (driver.go) must be up and running on their respective ports so that he application can call the relevant endpoints. This can be done through the use of the 'go run x.go' where x is the API file name (driver/passenger). To start the application, run the main function with the code 'go run interface.go' in cmd. You will be presented with a menu of options:

### Application Options
- Login: Allows you to login as either a passenger or driver.
- Create an Account: Allows you to create a new account as either a passenger or driver.
- Quit: Exits the application.

### Login Options
Once you have selected the "Login" option, you will be presented with a menu of options:

- Passenger: Allows you to login as a passenger.
- Driver: Allows you to login as a driver.
- Quit: Returns you to the previous menu.

#### Passenger Login Options
If you select the "Passenger" option, a list of existing passengers will be showed with an index number before it. To log in as a specific passenger, simply input the index number. After selection, you will be logged in as the passenger and presented with the following menu of options:

- Book a Ride: Allows you to book a new ride.
- View Ride History: Shows a list of all the rides you have taken.
- Change Account Details: Allows you to change parts of your account like Name, Phone Number, and Email.
- Logout: Logs you out and returns you to the previous menu.

#### Driver Login Options
If you select the "Driver" option, a list of existing drivers will be showed with an index number before it. To log in as a specific driver, simply input the index number. After selection, you will be logged in as the driver and presented with the following menu of options:

- Start a Trip: Checks to see if there is any pending ride assigned to the Driver that they can start.
- End a Trip: Checks to see if there is any ongoing ride assigned to the Driver that they can end.
- Change Account Details: Allows you to change parts of your account like Name, Phone Number, Email, and Car License.
- Logout: Logs you out and returns you to the previous menu.

### Account Creation Options
If you select the "Create an Account" option, you will be presented with a menu of options:

- Passenger: Allows you to create a new passenger account.
- Driver: Allows you to create a new driver account.
- Quit: Returns you to the previous menu.

If you select the "Passenger" option, you will be prompted to enter your first name, last name, mobile number, and email address. A new passenger account will then be created with the provided details.

If you select the "Driver" option, you will be prompted to enter your first name, last name, mobile number, email address, identification number, and car license. A new driver account will then be created with the provided details.

<br />

## API Documentation 
### Passengers API Microservice

This microservice provides a set of endpoints for managing passenger and ride information. It allows clients to view all passengers, get and update specific passenger details, create a ride booking for a passenger, view a passenger's ride history, view driver rides with pending status, and update a driver ride status.

##### Endpoints:

- GET /api/v1/passengers: Returns a list of all passengers in the database. 

- GET /api/v1/passengers/{passengerid}: Returns details for the specified passenger.

- POST /api/v1/passengers/{passengerid}: Adds a new passenger with the specified passenger ID. The request body should contain a JSON object with the following fields: FirstName, LastName, MobileNumber, and Email. 

- PUT /api/v1/passengers/{passengerid}: Updates the details for the specified passenger ID. The request body should contain a JSON object with the following fields: FirstName, LastName, MobileNumber, and Email. 

- POST /api/v1/passengers/{passengerid}/ride/{pickup}/{dropoff}: Creates a new ride booking for the specified passenger. The pickup and dropoff parameters should contain the pickup and dropoff codes, respectively. A random available driver is assigned to the ride created and their status is set to 'Busy'. 

- GET /api/v1/passengers/{passengerid}/history: Returns a list of all rides taken by the specified passenger. 

- GET /api/v1/ride/{driverid}/driver/{ridestatus}: Returns a list of rides with the specified ride status for the specified driver. 

- PUT /api/v1/ride/{rideid}/{driverid}/{response}: Updates the ride status for the specified ride. The rideid parameter should contain the ID of the ride to update, the driverid parameter should contain the ID of the driver, and the response parameter should contain the new ride status (Pending, Ongoing or Completed).

###### Examples: 
- To get a list of all passengers: "GET /api/v1/passengers" 

- To get details for passenger with ID 123456: "GET /api/v1/passengers/123456" 

- To add a new passenger with ID 123456: "POST /api/v1/passengers/123456 {"FirstName": "John", "LastName": "Doe", "MobileNumber": "12345678", "Email": "john.doe@mail.com"}" 

- To update the details for passenger with ID 123456: "PUT /api/v1/passengers/123456 {"FirstName": "Jane", "LastName": "Doe", "MobileNumber": "98765432", "Email": "jane.doe@mail.com"}" 

- To create a ride booking for passenger with ID 123456 from pickup code ABCDEF to dropoff code QWERTY: "POST /api/v1/passengers/123456/ride/ABCDEF/QWERTY"

<br />

### Drivers API Microservice

This microservice provides a set of endpoints for managing driver information. It allows clients to view all drivers, get and update specific driver details, assign random driver to be busy, and set specific driver to be available.

##### Endpoints:

- GET /api/v1/drivers: Returns a list of all drivers in the database. 

- GET /api/v1/drivers/{driverid}: Returns details for the specified driver.

- POST /api/v1/drivers/{driverid}: Adds a new driver with the specified driver ID. The request body should contain a JSON object with the following fields: FirstName, LastName, MobileNumber, Email, IdentificationNumber, CarLicense, and DriverStatus. 

- PUT /api/v1/drivers/{driverid}: Updates the details for the specified driver ID. The request body should contain a JSON object with the following fields: FirstName, LastName, MobileNumber, Email, CarLicense, and DriverStatus. Here, the IdentificationNumber cannot be changed.

- GET /api/v1/ride: Retrieves a random driver's details from the database with an 'Available' DriverStatus. The Driver's status will be set to 'Busy'. This Endpoint will be called in the Passengers API.

- GET /api/v1/ride/{driverid}: Updates the specified driver's status to be 'Available'. This Endpoint will be called in the Passengers API.

###### Examples: 
- To get a list of all drivers: "GET /api/v1/drivers" 

- To get details for driver with ID 123456: "GET /api/v1/drivers/123456" 

- To add a new driver with ID 123456: "POST /api/v1/drivers/123456 {"FirstName": "John", "LastName": "Doe", "MobileNumber": "12345678", "Email": "john.doe@mail.com", "IdentificationNumber":"S12345678", "CarLicense": "ABCDEFG", "DriverStatus": "Available"}" 

- To update the details for driver with ID 123456: "PUT /api/v1/passengers/123456 {"FirstName": "Jane", "LastName": "Doe", "MobileNumber": "98765432", "Email": "jane.doe@mail.com", "CarLicense": "QWERTYU", "DriverStatus": "Busy"}" 

- To get a random 'Available' driver and set their status to 'Busy': "GET /api/v1/ride"

- To set driver's status with ID 123456 to 'Available': "GET /api/v1/ride/123456"
