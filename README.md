# go-rest-template-bluemix

Reusable template for building REST Web Services in Golang (deployed on IBM Bluemix / Cloud Foundry)

Uses gorilla/mux as a router/dispatcher and Negroni as a middleware handler.

## Overview of the code

TODO

### Passport struct

We are going to use a travel Passport for our example. I've chosen Id as the unique key for the passport because (in the UK), passport book numbers these days have a unique 9 character field length (e.g. 012345678).

```
type Passport struct {
  Id              string `json:"id"`
  FirstName       string `json:"first_name"`
  LastName        string `json:"last_name"`
  DateOfBirth     string `json:"date_of_birth"`
  LocationOfBirth string `json:"location_of_birth"`
}
```

The first time you create a struct, you may not be aware that uppercasing and lowercasing your field names have a meaning in Go. It's similar to public and private members in Java. Uppercase = public, lowercase = private. There are some good discussions on Stackoverflow about [this](http://stackoverflow.com/questions/21825322/why-golang-cannot-generate-json-from-struct-with-front-lowercase-character). The gist is that if field names with a lowercase won't be visible to json.Marshal.

You may not want to expose your data to the consumer of your web service in this format, so you can override the way your fields are marshalled by adding ``json:"first_name"`` to each field with the desired name.

## Inspired by

* http://vluxe.io/golang-web-api.html
* https://github.com/msanterre/canoe/blob/master/main.go
* http://alpacalunchbox.com/building-lightweight-apis-with-go/
* https://gist.github.com/leeprovoost/879151cf696b52d6cecb
