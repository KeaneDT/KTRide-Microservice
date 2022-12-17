DROP DATABASE IF EXISTS ETIASG1_Passengers;
CREATE DATABASE ETIASG1_Passengers;
USE ETIASG1_Passengers;
DROP TABLE IF EXISTS PassengerRides;
DROP TABLE IF EXISTS Passengers;
CREATE TABLE Passengers (
    PassengerID varchar(255) NOT NULL,
    FirstName varchar(255) NOT NULL,
    LastName varchar(255),
    MobileNumber varchar(255) UNIQUE,
    Email varchar(255) UNIQUE,
    PRIMARY KEY (PassengerID)
);
INSERT INTO Passengers (PassengerID, FirstName, LastName, MobileNumber, Email) values("P02183701823210", "Hugo", "Boss", "727289112", "hboss@gmail.com");
INSERT INTO Passengers (PassengerID, FirstName, LastName, MobileNumber, Email) values("P03122022125349", "Tim", "Limothy", "2329392839", "tlimo@gmail.com");
INSERT INTO Passengers (PassengerID, FirstName, LastName, MobileNumber, Email) values("P08731239139281", "Jane", "Foster", "12614523232", "ladythor@gmail.com");

CREATE TABLE PassengerRides (
    RideID varchar(255) NOT NULL,
    RideDate DateTime,
    PassengerID varchar(255) NOT NULL,
    PickupCode varchar(255),
    DropoffCode varchar(255),
    DriverID varchar(255) NOT NULL,
    CarLicense varchar(255),
    RideStatus varchar(255),
    primary key (RideID),
	foreign key (PassengerID) references Passengers (PassengerID)
);	

DROP DATABASE IF EXISTS ETIASG1_Drivers;
CREATE DATABASE ETIASG1_Drivers;
USE ETIASG1_Drivers;
DROP TABLE IF EXISTS Drivers;
CREATE TABLE Drivers (
    DriverID varchar(255) NOT NULL,
    FirstName varchar(255) NOT NULL,
    LastName varchar(255),
    MobileNumber varchar(255) UNIQUE,
    Email varchar(255),
    IdentificationNumber varchar(255) NOT NULL UNIQUE,
    CarLicense varchar(255),
    DriverStatus varchar(255),
    PRIMARY KEY (DriverID)
);
INSERT INTO Drivers (DriverID, FirstName, LastName, MobileNumber, Email, IdentificationNumber, CarLicense, DriverStatus) values("D03912291184839", "Keane", "Travasso", "90607457", "keanet@gmail.com", "S192837739", "L1R48902A", 'Available');
INSERT INTO Drivers (DriverID, FirstName, LastName, MobileNumber, Email, IdentificationNumber, CarLicense, DriverStatus) values("D89239820193840", "Dylan", "Tan", "278738273", "dTan@gmail.com", "S382167871", "L8298HA1", 'Available');
INSERT INTO Drivers (DriverID, FirstName, LastName, MobileNumber, Email, IdentificationNumber, CarLicense, DriverStatus) values("D98123012381023", "Ben", "Thia", "983727372", "bThia@gmail.com", "S993762392", "L28782H29", 'Available');