### Status:
[![Maintainability](https://api.codeclimate.com/v1/badges/80fad80606ae11afb8c1/maintainability)](https://codeclimate.com/github/RomanUtolin/RESTful-CRUD/maintainability)
### Description of project:
This is an example Golang project.
### Install
Here is the steps to run it with `docker-compose`
### move to directory
```
cd workspace
```
### Clone into your workspace
```
git clone git@github.com:RomanUtolin/RESTful-CRUD.git
```
### move to project
```
cd RESTful-CRUD
```
### Run the application
```
make start
```
### Run the test(test db in Docker)
```
make testDb

make test
```
### Routes
```
# Return all person
GET /person

Query param:
    email=
    phone=
    first_name=
    page=
    limit=

    example:
    GET /person?email=test@test.test&phone=1234&first_name=test$page=1&limit=5


# Return one person
GET /person/id

# Create person
POST /person

# Update person
PUT /person/id

# Delete person
DELETE /person/id
```