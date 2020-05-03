# GO CRUD Spawner

## Description 

this program is meant to create a simple CRUD (Create, Read, Update, Delete) application with a MySQL backend database.

## Prerequisites
1. install GoLang the proper files and instructions for that can be found here https://golang.org/doc/install 
2. Install MySQL This page can help you with this selecting the proper version for your OS 
    * Install instructions: https://dev.mysql.com/doc/mysql-installation-excerpt/5.7/en/windows-installation.html
    * Community downloads: https://dev.mysql.com/downloads/installer/
3. Go SQL driver which can be found here or by typing the following into the command line 
```sh
   go get -u github.com/go-sql-driver/mysql 
```
4. download (and extract if neccesary) the master git repo

## Launch Instructions
1. MAKE SURE YOU HAVE ALL THE PREREQUISITS
2. navigate or open your command shell to the root file directory where you downloaded the git repo to
3. type the following into the command line this should start the server you may have to give permissions to the program to use your network
```sh
   go run server.go 
``` 
4. the first time the program is ran it will prompt the user for credentials to access your database and the name of the database. if youve done this already the credentials and database name are saved in a JSON file in configs which is pulled upon server start so if you need to change this in the future you can just change the JSON file.

## How it Works
### Initialization (first time run)
1. it starts by creating the database and tables it if you run this subsequent times it will log an error but not stop stating that the database and table already exist.
2. it creates a few template files that have some HTML formatting in them upon startup the server checks to make sure these files exist and if they do it does nothing so editing these files for your personal preference of styling can be done without consequence. if they do not exist or you delete any of them it will recreate the file to its original state
### Server start
after user inputs the server starts listening on port 8000 by default it can be changed by hard coding the first string value on line 291 but it would be suggested to also change the log on line 277 so that it prints the correct string to reflect.
the URL http://localhost:8000 by default will take you to your server instance 


## Future Plans
##### NOTE: crossouts mean that this functionality has been implemented but other branches may still have this apply to them
1. ~~integrate a prompt to the user through the terminal to input the dbUser and dbPass variables on line 120 and line 121 which will create a json file where its stored.~~
    * ~~loop back to asking for credentials if access is denied~~
    * store the password in a secure way
2. allow for user input into the terminal to tell the system the table name, how many fields the table has and customize the column names and types with an error check that loops back to the user if the type is not valid (just to keep it all nice and happy)
3. allow for multiple tables to be created with association
4. make it runnable in either web page, API, or hybrid mode where hybrid has both an api server and the web page references the API and you can also use the seperate API server for raw JSON streams 
    * might also add support for other serialization types like XML 
5. add integration for other database types
    * SQL
        * ~~MySQL~~
        * Postgres 
    * NoSQL
        * MongoDB

