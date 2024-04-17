CREATE TABLE IF NOT EXISTS Animals(
    ID serial,
    Kind_Of_Animal varchar(255) NOT NULL,
    Kind_Of_Breed varchar(255) NOT NULL,
	Name varchar(255) NOT NULL,
    Age int NOT NULL,
	Description text NOT NULL
    );

CREATE TABLE IF NOT EXISTS Users(
    UserID serial,
    User_Email text NOT NULL,
	Username varchar(30) NOT NULL,
	Password text NOT NULL,
    Number_of_phone_user varchar(50) NOT NULL,
    Profile_Picture_User text NOT NULL
    );

CREATE TABLE IF NOT EXISTS Roles(
	Role_name  varchar(50) PRIMARY KEY,
	Permissions text
);
	
CREATE TABLE IF NOT EXISTS Admins(
	AdminID serial,
    Admin_Email text,
	Adminame varchar(30),
	Password text,
    Number_of_phone_Admin varchar(50),
    Profile_Picture_Admin text,
	Role varchar(50),
	foreign key(Role) references Roles(Role_name)
);
