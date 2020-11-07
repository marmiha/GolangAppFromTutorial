## GolangAppFromTutorial

### Summary
This app was made by following the tutorial of youtube channel [EQuimper](https://www.youtube.com/channel/UC7R7bcH9-KEBDiGNP1mZnmw) 
for purpose of learning Go language. The tutorial playlist can be found on [this](https://www.youtube.com/watch?v=Uuy9J33iG0E&list=PLzQWIQOqeUSPFPVfticl-CsmUv82Gb5W-&index=1) 
link. Albeit I have followed the tutorial I've tried to implement a few quirks of my own, just to experiment with the 
yet unfamiliar programming language.

### Project structure

The application made of three different components - database specific logic, business logic and API endpoints
specific configurations. Through the project our main goal is to separate our API and database choice from the business 
logic of the application itself which is mostly achieved by interfaces and packages. This way we ensure that our business
logic is independent of the implementation and choice of database and API architecture as the communication between all is
done over interfaces. I have chosen `REST` design for endpoints combined with a `PostgreSQL` database but if one chose to
change it for example to `GraphQL` support or a different database the code base of the business logic would
stay the same.

Each of the following packages can be found in the folder with the same name as the package.

#### `package domain`
Contains the business logic of the application and is used by API endpoints. This is where we defined custom 
business logic errors and the interfaces `type <Model>Repository interface {}` where model stands for the name of the 
corresponding database model. These are then implemented in the database logic. All functions return Golang errors, which
are then properly interpreted in the API endpoints `package handlers`. This is where we also define our payload structs
which are used by our API endpoints to pass data.

#### `package postgres`

Contains implementations of the business logic interfaces and the database connection specific setup functions. Mostly 
code for our CRUD operations.

#### `package handlers`

Contains our `REST` API endpoints implementation. These accept data over HTTP and respond accordingly to the data validation
or business errors that occur with the right status codes and payload values.
