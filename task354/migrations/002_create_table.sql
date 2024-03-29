-- +goose Up

CREATE TABLE IF NOT EXISTS Users (
ID SERIAL PRIMARY KEY,
Name VARCHAR(255) NOT NULL
);

CREATE TABLE IF NOT EXISTS Books (
ID SERIAL PRIMARY KEY,
Title VARCHAR(255) NOT NULL,
AuthorID INT NOT NULL,
IsRented BOOLEAN DEFAULT FALSE,
FOREIGN KEY (AuthorID) REFERENCES Authors(ID)
);

CREATE TABLE IF NOT EXISTS Authors (
ID SERIAL PRIMARY KEY,
Name VARCHAR(255) NOT NULL
);

CREATE TABLE IF NOT EXISTS Rentals (
RentalID SERIAL PRIMARY KEY,
UserID INT,
BookID INT,
RentDate TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
ReturnDate TIMESTAMP,
FOREIGN KEY (UserID) REFERENCES Users(ID),
FOREIGN KEY (BookID) REFERENCES Books(ID)
);
