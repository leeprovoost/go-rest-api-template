# go-rest-template-bluemix

*WORK IN PROGRESS*

Reusable template for building REST Web Services in Golang (deployed on IBM Bluemix / Cloud Foundry)

Uses gorilla/mux as a router/dispatcher and Negroni as a middleware handler.

## Overview

### Code Structure

TO DO

### Live Code Reloading

TO DO

### Main Function

TO DO

### Data structure

We are going to use a travel Passport for our example. I've chosen Id as the unique key for the passport because (in the UK), passport book numbers these days have a unique 9 character field length (e.g. 012345678). A passport belongs to a user and a user can have one or more passports.

```golang
type User struct {
  Id              int    `json:"id"`
  FirstName       string `json:"first_name"`
  LastName        string `json:"last_name"`
  DateOfBirth     string `json:"date_of_birth"`
  LocationOfBirth string `json:"location_of_birth"`
}

type Passport struct {
  Id           string `json:"id"`
  DateOfIssue  string `json:"date_of_issue"`
  DateOfExpiry string `json:"date_of_expiry"`
  Authority    string `json:"authority"`
  CustomerId   int    `json:"customer_id"`
}
```

The first time you create a struct, you may not be aware that uppercasing and lowercasing your field names have a meaning in Go. It's similar to public and private members in Java. Uppercase = public, lowercase = private. There are some good discussions on Stackoverflow about [this](http://stackoverflow.com/questions/21825322/why-golang-cannot-generate-json-from-struct-with-front-lowercase-character). The gist is that if field names with a lowercase won't be visible to json.Marshal.

You may not want to expose your data to the consumer of your web service in this format, so you can override the way your fields are marshalled by adding ``json:"first_name"`` to each field with the desired name.

### API Routes

Now that we have defined the data model, we need to translate that to a REST interface:

* Retrieve a list of all users: `GET /users` -> The `GET` just refers to the HTTP action you would use. If you want to test this in the command line, then you can use curl: `curl -X GET http://localhost:3009/users` or `curl -X POST http://localhost:3009/users`
* Retrieve the details of an individual user: `GET /users/{uid}` -> {uid} allows us to create a variable, named uid, that we can use in our code. An example of this url would be `GET /users/1`
* Create a new user: `POST /users`
* Update a user: `PUT /users/{uid}`
* Delete a user: `DELETE /users/{uid}`

We now need to do the same for handling passports. Don't forget that a passport belongs to a user, so to retrieve a list of all passports for a given user, we would use `GET /users/{uid}/passports`.

When we want to retrieve an specific passport, we don't need to prefix the route with `/users/{uid}` anymore because we know exactly which passport we want to retrieve. So, instead of `GET /users/{uid}/passports/{pid}`, we can just use `GET /passports/{pid}`.

Once you have the API design sorted, it's just a matter of creating the code that gets called when a specific route is hit. We implement those with Handlers.

```golang
  router.HandleFunc("/c", UsersHandler).Methods("GET")
  router.HandleFunc("/users/{uid}", UsersHandler).Methods("GET")
  router.HandleFunc("/users", UsersHandler).Methods("POST")
  router.HandleFunc("/users/{uid}", UsersHandler).Methods("PUT")
  router.HandleFunc("/users/{uid}", UsersHandler).Methods("DELETE")

  router.HandleFunc("/users/{uid}/passports", PassportsHandler).Methods("GET")
  router.HandleFunc("/passports/{pid}", PassportsHandler).Methods("GET")
  router.HandleFunc("/users/{uid}/passports", PassportsHandler).Methods("POST")
  router.HandleFunc("/passports/{pid}", PassportsHandler).Methods("PUT")
  router.HandleFunc("/passports/{pid}", PassportsHandler).Methods("DELETE")
```

Last but not least, we want to handle two special cases:

```golang
  router.HandleFunc("/", HomeHandler)
  router.HandleFunc("/healthcheck", HealthcheckHandler).Methods("GET")
```

When someone hits our API, without a specified route, then we can handle that with either a standard 404 (not found), or any other type of feedback.

We also want to set up a health check that monitoring tools like [Sensu](https://sensuapp.org/) can call: `GET /healthcheck`. The health check route can return a 200 OK when the serivce is up and running, including some extra stats. Your DevOps colleagues will be very grateful for this.

### Route Handlers

TO DO

* Use render for generating JSON, see https://github.com/unrolled/render/issues/7 for use of global variable

### Mock Data

TO DO

### Testing

TO DO

### Security

TO DO

### Bluemix / Cloud Foundry

TO DO

### Other

TO DO

## Useful references

* https://github.com/codegangsta/negroni
* http://vluxe.io/golang-web-api.html
* https://github.com/msanterre/canoe/blob/master/main.go
* http://alpacalunchbox.com/building-lightweight-apis-with-go/
* https://gist.github.com/danesparza/eb3a63ab55a7cd33923e
* http://stackoverflow.com/questions/21825322/why-golang-cannot-generate-json-from-struct-with-front-lowercase-character
