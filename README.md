# go-rest-api-template

A template for building REST Web Services in Go. Uses gorilla/mux as a router/dispatcher and Negroni as a middleware handler. Tested against `go1.15.6 darwin/amd64`.

## Introduction

### Why?

After writing many REST APIs with Java Dropwizard, Node.js/Express and Go, I wanted to distill my lessons learned into a reusable template for writing REST APIs, in the Go language.

It's mainly for myself. I don't want to keep reinventing the wheel and just want to get the foundation of my REST API 'ready to go' so I can focus on the business logic and integration with other systems and data stores.

Just to be clear: this is not a framework, library, package or anything like that. This tries to use a couple of very good Go packages and libraries that I like and cobbled them together.

The main ones are:

* [gorilla/mux](http://www.gorillatoolkit.org/pkg/mux) for routing
* [urfave/negroni](https://github.com/urfave/negroni) as a middleware handler
* [strechr/testify](https://github.com/stretchr/testify) for writing easier test assertions
* [unrolled/render](https://github.com/unrolled/render) for HTTP response rendering
* [palantir/stacktrace](https://github.com/palantir/stacktrace) to provide more context to error messages
* [unrolled/secure](https://github.com/unrolled/secure) to improve API security
* [sirupsen/logrus](https://github.com/sirupsen/logrus) for structured logging

Whilst working on this, I've tried to write up my thought process as much as possible. Everything from the design of the API and routes, some details of the Go code like JSON formatting in structs and my thoughts on testing. However, if you feel that there is something missing, send a PR or raise an issue.

Last but not least: I'm occasionally writing about my lessons learned developing Go applications and Amazon Web Services systems on [Medium](https://medium.com/@leeprovoost).

### Knowledge of Go

If you're new to programming in Go, I would highly recommend you to read the following resources:

* [A Tour of Go](https://tour.golang.org/welcome/1)
* [Effective Go](https://golang.org/doc/effective_go.html)
* [50 Shades of Go: Traps, Gotchas, and Common Mistakes for New Golang Devs](http://devs.cloudimmunity.com/gotchas-and-common-mistakes-in-go-golang/)
* [Learn X in Y Minutes](https://learnxinyminutes.com/docs/go/)

You can work your way through those in two to three days. I did try out a couple of books but didn't find them that useful. If you really want to get some books, there is a [github](https://github.com/dariubs/GoBooks) repo that tries to list the main ones.

A good project to keep on eye on to discover new Go packages and software is [awesome-go](https://github.com/avelino/awesome-go). The maintainers enforce strict standards around documentation, tests, structure, etc. so they are doing a brilliant job improving the overall quality of code in the Go community. It's not easy to keep track of the high rate of change on that project, so my suggestion would be to track new additions using the [changedetection.com](https://www.changedetection.com/) website.

### Development tools

The choice of an editor/IDE is a personal one, so the only advice I can give you is to try different ones. It looks like currently the following seem to be the most mature:

* [Atom](https://atom.io/) with [go-plus](https://github.com/joefitzgerald/go-plus) 
* vim with [vim-go](https://github.com/fatih/vim-go)
* Microsoft did a pretty good job with their [Visual Studio Code](https://code.visualstudio.com/) and then install the [Go plugin](https://code.visualstudio.com/docs/languages/go)
* IntelliJ [Goland](https://www.jetbrains.com/go/)

### How to run

`go run main.go` works fine if you have a single file or a script you're working on, but once you have a more complex project with lots of files then you'll have to start using the proper go build tool and run the compiled executable:

```
go build && ./api-service
```

Which is faster than first typing `go build` to generate an executable called `api-service` and then run that executable by typing `./api-service`.

The `api-service` app will bind itself to the port that you have defined in your environment variables or your Makefile (in our case `3001`). You can copy paste the following in your terminal to run the application (ensure you are in the `cmd/api-service` folder):

```
export ENV=LOCAL
export VERSION=VERSION
export PORT=3001
export FIXTURES=./fixtures.json
go build && ./api-service
```

One comment from a security point of view: you will commonly read that port 80 is the default port for HTTP (and 443 for HTTPS), it's usually not recommended to actually bind your application to that port when it runs on a server instance like AWS EC2. That's because ports below 1024 are privileged ports and require elevated privileges for the app to bind itself to such port. If someone compromises your app, it could potentially take advantage of the application user with the elevated privileges. Instead, run it on a higher port like 8001 and run it with a newly created restricted user (instead of using the root user). You can then use a load balancer or a proxy (e.g. AWS ELB/ALB, HAProxy, nginx) to accept your incoming request from the internet on port 80 or 443 and then forward that to your application port.

### Live Code Reloading

Manually stopping and restarting your server during development can get quite annoying after a while, so let's set up a task runner that automatically restarts the server when it detects changes. People from the JavaScript world will be quite familiar with tools like gulp. I've chosen a Go native task runner called fresh.

Install [fresh](https://github.com/pilu/fresh):

```
go get github.com/pilu/fresh
```

And during development, you can now just type in the following command in your project root directory:

```
fresh
```

You should see following output if all goes well:

```
Loading settings from ./runner.conf
11:10:42 runner      | InitFolders
11:10:42 runner      | mkdir ./tmp
11:10:42 runner      | mkdir ./tmp: file exists
11:10:42 watcher     | Watching .
11:10:42 watcher     | Watching vendor
11:10:42 watcher     | Watching vendor/github.com
11:10:42 watcher     | Watching vendor/github.com/codegangsta
11:10:42 watcher     | Watching vendor/github.com/codegangsta/negroni
11:10:42 watcher     | Watching vendor/github.com/davecgh
11:10:42 watcher     | Watching vendor/github.com/davecgh/go-spew
11:10:42 watcher     | Watching vendor/github.com/davecgh/go-spew/spew
11:10:42 watcher     | Watching vendor/github.com/gorilla
11:10:42 watcher     | Watching vendor/github.com/gorilla/context
11:10:42 watcher     | Watching vendor/github.com/gorilla/mux
11:10:42 watcher     | Watching vendor/github.com/palantir
11:10:42 watcher     | Watching vendor/github.com/palantir/stacktrace
11:10:42 watcher     | Watching vendor/github.com/palantir/stacktrace/cleanpath
11:10:42 watcher     | Watching vendor/github.com/pmezard
11:10:42 watcher     | Watching vendor/github.com/pmezard/go-difflib
11:10:42 watcher     | Watching vendor/github.com/pmezard/go-difflib/difflib
11:10:42 watcher     | Watching vendor/github.com/stretchr
11:10:42 watcher     | Watching vendor/github.com/stretchr/testify
11:10:42 watcher     | Watching vendor/github.com/stretchr/testify/assert
11:10:42 watcher     | Watching vendor/github.com/unrolled
11:10:42 watcher     | Watching vendor/github.com/unrolled/render
11:10:42 watcher     | Watching vendor/github.com/unrolled/secure
11:10:42 main        | Waiting (loop 1)...
11:10:42 main        | receiving first event /
11:10:42 main        | sleeping for 600 milliseconds
11:10:42 main        | flushing events
11:10:42 main        | Started! (75 Goroutines)
11:10:42 main        | remove tmp/go-rest-api-template.log: no such file or directory
11:10:42 build       | Building...
11:10:48 runner      | Running...
11:10:48 main        | --------------------
11:10:48 main        | Waiting (loop 2)...
11:10:48 app         | [negroni] listening on localhost:3001
11:10:48 app         | 2016/09/28 11:10:48 ===> Starting app (v0.2.0) on port 3001 in LOCAL mode.
```

Fresh should work without any configuration, but to make it more explicit you can add a `runner.conf` file in your project root:

```
root:              .
tmp_path:          ./tmp
build_name:        api-service
build_log:         api-service.log
valid_ext:         .go, .tpl, .tmpl, .html, .mod, .sum, modules.txt
build_delay:       600
colors:            1
log_color_main:    cyan
log_color_build:   yellow
log_color_runner:  green
log_color_watcher: magenta
```

As you can see, it creates a `tmp` directory in your project root and a log file. You can tell `.gitignore` to stop checking it into your git repository by adding the following lines in your `.gitignore` file:

```
# fresh
tmp
api-service.log
```

I've added a simple task to our `Makefile` that allows you to run the app with all the correct environment variables. Just run:

```
make fresh
```

Where I work, we've hit the limitations of fresh and have switched to Gulp as a task runner. Our Gulp setup takes care of running tests, using OSX notifications to warn you of build failures, etc. 

### High-level Code Structure

The following structure is something that has been evolved over the years. Back when I started writing Go applications, we didn't have best practices like this project ```https://github.com/golang-standards/project-layout```. I will sound very old now, but "back in the olden days", Go didn't even have a proper package/module management system, so we've come a long way. Let me first show the project structure and explain what the key folders do, immediately followed by some explanation what each file does:

```
- cmd --> folder that has one or more application entry points
  |- api-service
     |- tmp
     |- fixtures.json
     |- main.go
     |- Makefile
     |- runner.conf
     |- VERSION
- internal --> folder where we store application code that should not be reused by other applications
  |- passport
     |- models
        |- passport.go
        |- user.go
     |- appenv.go
     |- db_user_test.go
     |- db_user.go
     |- handlers_test.go
     |- handlers.go
     |- main.go
     |- routes.go
- pkg --> folder where we store application code that can be reused by other applications
  |- health
     |- check.go
  |- status
     |- response.go
  |- version
     |- parser.go
- vendor --> directory where we store the vendored modules
- .gitignore
- go.mod
- go.sum
- LICENSE
- README.md
```

And now with some annotations:

```
- cmd
  |- api-service --> folder for our passport api
     |- tmp --> temporary folder for our fresh task runner
     |- fixtures.json --> test data in JSON format
     |- main.go --> this ties all the various application logic packages together and loads the correct environment variables
     |- Makefile --> Makefile that describes some tasks:
     |- runner.conf --> configuration file for [fresh](go get github.com/pilu/fresh)
     |- VERSION --> defines application version using semantic versioning
- internal 
  |- passport
     |- models
        |- passport.go --> passport data structure definition
        |- user.go --> user data structure definition
     |- appenv.go --> defines the application context object
     |- db_user_test.go --> tests our mock/fake database
     |- db_user.go --> our mock/fake database implementation
     |- handlers_test.go --> tests for our route handlers
     |- handlers.go --> defines the actual logic that gets executed when you visit a route and takes care of the response to the client
     |- main.go --> actual server initialisation is happening here: mainly loading of routes and middleware
     |- routes.go --> server routes, e.g. GET /healthcheck, and binding of those routes to the correct route handlers
- pkg 
  |- health
     |- check.go --> healthcheck data structure definition
  |- status
     |- response.go --> status data structure definition
  |- version
     |- parser.go --> parses VERSION file
- vendor
- .gitignore --> define what files git should ignore and not check in
- go.mod --> Go modules definition file 
- go.sum --> Go modules definition file 
- LICENSE --> (optional) License file
- README.md --> this document
```

### Starting the application

The main entry point to our app is `cmd/api-service/main.go` where we take care of the following:

Configuring our logrus structured logging package. Instead of using the standard log package, we use logrus that allows us to print information in a JSON format. However, during development, we just want to see a human-readable format:

```
func init() {
	if "LOCAL" == strings.ToUpper(os.Getenv("ENV")) {
		log.SetFormatter(&log.TextFormatter{})
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetFormatter(&log.JSONFormatter{})
		log.SetLevel(log.InfoLevel)
	}
}
```

Next up is to load our environment variables. The default values that we are using during development are defined in our `Makefile`.

```
var (
  // try to read environment variables
  env      = os.Getenv("ENV")      // LOCAL, DEV, STG, PRD
  port     = os.Getenv("PORT")     // server traffic on this port
  version  = os.Getenv("VERSION")  // path to VERSION file
  fixtures = os.Getenv("FIXTURES") // path to fixtures file
)
```

I'm using the environment variables technique to tell the app in what environment it lives:

* `LOCAL`: local development environment
* `DEV`: development or integration server
* `STG`: staging servers
* `PRD`: production servers

We then load our version file. I have added a helper function in `pkg/version/parse.go` that can parse a `VERSION` file using semantic versioning.

```
// reading version from file
version, err := vparse.ParseVersionFile(version)
if err != nil {
  log.WithFields(log.Fields{
    "env":  env,
    "err":  err,
    "path": os.Getenv("VERSION"),
  }).Fatal("Can't find a VERSION file")
  return
}
log.WithFields(log.Fields{
  "env":     env,
  "path":    os.Getenv("VERSION"),
  "version": version,
}).Info("Loaded VERSION file")
```

When we can't find a version file, then we throw an error. When we do find a version file, then we're just writing a log entry.

The check for error (`if err != nil`) is a common Go way of handling return values that contain errors. the `main` function is one of the few places where I use log.Fatal in my application. `log.Fatal` logs your file to the console, but also exits your application. This is for me the correct behaviour here because without a correct `VERSION` file, the app shouldn't work. It's too risky that you may have deployed the incorrect application version.

This API doesn't talk to a real database but to a very simple in-memory mock database. The data is hardcoded into the `CreateMockDataSet()` function and loaded at application start time. I will change this with an actual database at some point. 

```
userStore := passport.NewUserService(passport.CreateMockDataSet())
```

Note that `func NewUserService(list map[int]models.User, count int) models.UserStorage` expects two parameters, but it looks like we are only passing one? That's because the `func CreateMockDataSet() (map[int]models.User, int) ` function returns two results. The first one is the list of users, the second one is the count. Hence why this satisfies the type signature of `NewUserService`.

Once all the data is read from our environment variables and the mock data has been loaded, we can then initialise our application context. `AppEnv` is a struct that holds some application info that we can pass around (e,g, provide to our handlers). It's a neat way to avoid using global variables and avoiding that you have application variables all over the place.

I have defined the `AppEnv` struct in `internal/passport/appenv.go`:

```
// AppEnv holds application configuration data
type AppEnv struct {
	Render  *render.Render
	Version string
	Env     string
	Port    string
	DB      DataStorer
}
```

`Render` helps us with rendering the correct response to the client, e.g. JSON, XML, etc. Version holds the path to our `VERSION` file. `Env` holds information about our platform environment (e.g. STG, DEV). `Port` is the server port that our application binds to. `DB` is our database struct that provides a data abstraction layer.

Once all of that is set up properly, we can now start the server. The `StartServer` function takes our `AppEnv` struct and initialises all our routes and starts our `negroni` server. This is all defined in `internal/passport/main.go`.

Most `negroni` code examples ([like the one on the official GitHub repo](https://github.com/urfave/negroni)) will start the application using:

```
n := negroni.Classic()
```

When you look at our `internal/passport/main.go` file, then you'll see that we do it slightly differently:

```
// start now
n := negroni.New()
```

`negroni.Classic()` adds a little bit of magic in the background and for instance loads the `/public` folder as the location of your static files. I personally don't like that and prefer to make things more explicit. It's up to you to decide which one you prefer, as long as you are aware what exactly happens when you play with magic.

One important thing worth mentioning is the way we run the application and bind the app to a given port. It's common to see Go apps being started as:

```
n.Run(":3001")
```

However, when you're running this on OSX, you will get constant warnings about accepting incoming connections. I've written a piece on [Medium](https://medium.com/@leeprovoost/suppressing-accept-incoming-network-connections-warnings-on-osx-7665b33927ca#.crake4bm9) about this but the tl;dr solution is to explicitly bind your application to `localhost`:

```
n.Run("localhost:3001")
```

Last but not least, a big thank you to [edoardo849](https://github.com/edoardo849) for providing some great feedback on structuring the API and reducing `main.go` complexity.

## Data

### Data model

We are going to use a travel Passport for our example. I've chosen Id as the unique key for the passport because (in the UK), passport book numbers these days have a unique 9 character field length (e.g. 012345678). A passport belongs to a user and a user can have one or more passports. We'll define this in `internal/passport/models/passport.go` and `internal/passport/models/user.go`.

```
type User struct {
  ID              int       `json:"id"`
  FirstName       string    `json:"firstName"`
  LastName        string    `json:"lastName"`
  DateOfBirth     time.Time `json:"dateOfBirth"`
  LocationOfBirth string    `json:"locationOfBirth"`
}

type Passport struct {
  ID           string    `json:"id"`
  DateOfIssue  time.Time `json:"dateOfIssue"`
  DateOfExpiry time.Time `json:"dateOfExpiry"`
  Authority    string    `json:"authority"`
  UserID       int       `json:"userId"`
}
```

The first time you create a struct, you may not be aware that uppercasing and lowercasing your field names have a meaning in Go. It's similar to public and private members in Java. Uppercase = public, lowercase = private. There are some good discussions on Stackoverflow about [this](http://stackoverflow.com/questions/21825322/why-golang-cannot-generate-json-from-struct-with-front-lowercase-character). The gist is that field names that start with a lowercase letter will not be visible to json.Marshal.

You may not want to expose your data to the consumer of your web service in this format, so you can override the way your fields are marshalled by adding ``json:"firstName"`` to each field with the desired name. I admit that in the past I had the habit of using underscores for my json field names, e.g. `first_name`. However after reading [this](http://www.slideshare.net/stormpath/rest-jsonapis) excellent presentation on API design, I got reminded that the JS in JSON stands for JavaScript and in the JavaScript world, it's common to use camelCasing so the preferred way of writing the same field name would be: `firstName`.

Note the use of `time.Time` for the dates instead of using a standard `string` type. We'll discuss the pain of marshalling and unmarshalling of JSON dates a bit later.

If you want to prevent a certain struct field to be marshalled/unmarshalled then use `json:"-"`.

During development, it can be useful to print the contents of some variables to your console. We've added a little helper so Go knows how to properly display the contents of the structs. See this example in `models/user.go`:

```
// User holds personal user information
type User struct {
	ID              int       `json:"id"`
	FirstName       string    `json:"firstName"`
	LastName        string    `json:"lastName"`
	DateOfBirth     time.Time `json:"dateOfBirth"`
	LocationOfBirth string    `json:"locationOfBirth"`
}

// GoString implements the GoStringer interface so we can display the full struct during debugging
// usage: fmt.Printf("%#v", i)
// ensure that i is a pointer, so might need to do &i in some cases
func (u *User) GoString() string {
	return fmt.Sprintf(`
{
	ID: %d,
	FirstName: %s,
	LastName: %s,
	DateOfBirth: %s,
	LocationOfBirth: %s,
}`,
		u.ID,
		u.FirstName,
		u.LastName,
		u.DateOfBirth,
		u.LocationOfBirth,
	)
}
```

### Operations on our (mock) database

I wanted to create a template REST API that didn't depend on a database, so started with creating a simple in-memory database that we can work with. I admit that this was mainly to satisfy my own curiosity. I might change that in the future to using an off-the-shelf in-memory Go database.

The good thing is that this will be the start of a so-called data access layer that abstracts away the underlying data store. We can achieve that by starting with creating an interface (which is a good practice in Go anyway). 

Open the `internal/passport/models/user.go` file and you will see:

```
// UserStorage defines all the database operations
type UserStorage interface {
	ListUsers() ([]User, error)
	GetUser(i int) (User, error)
	AddUser(u User) (User, error)
	UpdateUser(u User) (User, error)
	DeleteUser(i int) error
}
```

This allows us to define a set of operations on the data as a contract, without people having to worry about the actual implementation of how the data is stored and accessed. I've added the basic operations to list, retrieve, create, update and delete data, so the standard CRUD-style operations (accepting that CRUD has some subtle differences with REST).

Let's now have a look at the type signature of the `Get` operation:

```
GetUser(i int) (User, error)
```

What this tells us is that it is expecting an integer as an argument (which will be the User id in our case), and returns a pair of values: a User object and an error object. Returning pairs of values is a nice Go feature and is often used to return information about errors.

An example of how this could be used is in `internal/passport/handlers.go:

```
// GetUserHandler returns a user object
func GetUserHandler(w http.ResponseWriter, req *http.Request, appEnv AppEnv) {
	vars := mux.Vars(req)
	uid, _ := strconv.Atoi(vars["uid"])
	user, err := appEnv.UserStore.GetUser(uid)
	if err != nil {
		response := status.Response{
			Status:  strconv.Itoa(http.StatusNotFound),
			Message: "can't find user",
		}
		log.WithFields(log.Fields{
			"env":    appEnv.Env,
			"status": http.StatusNotFound,
		}).Error("Can't find user")
		appEnv.Render.JSON(w, http.StatusNotFound, response)
		return
	}
	appEnv.Render.JSON(w, http.StatusOK, user)
}
```

We check whether the error object is nil. If it is, then we return a HTTP 200 OK, if not then we return HTTP 404 Not Found. Let's go into more detail when we talk about our API handlers.

BTW in other languages you might be used to write the above as:

```
user, err := appEnv.UserStore.GetUser(uid)
if err == nil {
  ctx.render.JSON(w, http.StatusOK, user)
} else {
  ctx.render.JSON(w, http.StatusNotFound, err)
}
```

But the Go way is to check if an error exists and handle that immediately. If there is no error, just continue with the normal application flow.

Let's have a look at the actual mock in-memory database. We need to create a struct that will hold the data:

```
// UserService will hold the connection and key db info
type UserService struct {
	UserList  map[int]models.User
	MaxUserID int
}
```
The UserList will hold a list of User structs and the MaxUserID holds the latest used integer. MaxUserID mimics the behaviour of an auto-generated ID in conventional databases. In this case MaxUserID represents the highest used database ID.

Let's have a closer at the `GetUser` function:

```
// GetUser returns a single JSON document
func (service *UserService) GetUser(i int) (models.User, error) {
	user, ok := service.UserList[i]
	if !ok {
		return models.User{}, stacktrace.NewError("Failure trying to retrieve user")
	}
	return user, nil
}
```

* `func (service *UserService)`: This might be something new for you and I admit that it took me a little while to appreciate this. It's called a pointer receiver and is part of [A Tour of Go](https://tour.golang.org/methods/4). This means that the `GetUser` function can work on a `UserService` struct. For people coming from languages that use objects and methods, it's similar. 
* `GetUser(i int) (models.User, error)`: This is the function signature and says that we have a function named `GetUser` that accepts one argument of the type `int` (e.g. 0, 1, 2, 3) and will return two variables. The first one is a `User` struct, the second one is an `error` struct.
* `return models.User{}, stacktrace.NewError("Failure trying to retrieve user")`: In case there is an error, we return an empty `User` object and an `error`. We have to return an `User` object because that's what the function signature enforces, but this will be empty and frankly doesn't really matter. It doesn't matter because we will first check the `error` struct. The basic `error` struct in Go doesn't really say much, which is why we use here Palantir's `stacktrace` package. It adds some more contextual information about where the issue occured.
* `return user, nil`: If everything goes well, then we just return the user data as well as an empty error (nil).

Before we move on to loading the mock data, I wanted to briefly touch on the `UserStorage` interface. One of the advantages of using an interface is that you can swap out the implementation. This can be useful for testing. So instead of `GetUser` calling the actual datastore, you could write a different implementation that just returns a hardcoded `User` struct. However, Go has a implicit way of linking the interface to the struct. Unlike in Java where you have to explicitly state that a Class implements an Interface (and thus you benefit from the IDE checking whether you have implemented all the methods), Go just tries to figure out whether there are structs that match the definition of your interface. In order to make this more explicit, I always add a small test, see at the top of `internal/passport/db_user.go`:

```
// Compile-time proof of interface implementation
var _ models.UserStorage = (*UserService)(nil)
```

This tries to typecast the `UserService` pointer to the `UserStorage` interface. If you haven't implemented your interface functions properly, then your Go tests will complain.

### Mock data

In order to make it a bit more useful, we will initialise it with some user objects. I've created a helper function in `internal/passport/db_user.go` called `CreateMockDataSet` that just loads some test data. The structure of that data looks as follows:

```
dt, _ := time.Parse(time.RFC3339, "1985-12-31T00:00:00Z")
list[0] = models.User{
	ID:              0,
	FirstName:       "John",
	LastName:        "Doe",
	DateOfBirth:     dt,
	LocationOfBirth: "London",
}
```

The date string looks a bit odd. Why not just use `31-12-1985` or `1985-12-31`? The first is discouraged altogether because that's a European way of writing dates and will cause confusion around the world. Imagine you have 3-4-2015. Is it April the third or March the fourth? Unfortunately there isn't an "enforced standard" for dates in JSON, so I've tried to use one that is commonly used and also understood by Go's `json.Marshaler` and `json.Unmarshaler` to avoid that we have to write our own custom marshaler/unmarshaler.

If you have a look at Go's `time/format` [code](http://golang.org/src/time/format.go) then you'll see on line 54:

```
54    RFC3339     = "2006-01-02T15:04:05Z07:00"
```

That's the one we need. It has the following format:

```
year-month-day
T (for time)
hour:minutes:seconds
Z (for time zone)
offset from UTC
```

## Defining the API

### API Routes

Now that we have defined the data access layer, we need to translate that to a REST interface:

* Retrieve a list of all users: `GET /users` -> The `GET` just refers to the HTTP action you would use. If you want to test this in the command line, then you can use curl: `curl -X GET http://localhost:3009/users`
* Retrieve the details of an individual user: `GET /users/{uid}` -> {uid} allows us to create a variable, named uid, that we can use in our code. An example of this url would be `GET /users/1`
* Create a new user: `POST /users`
* Update a user: `PUT /users/{uid}`
* Delete a user: `DELETE /users/{uid}`

We now need to do the same for handling passports. Don't forget that a passport belongs to a user, so to retrieve a list of all passports for a given user, we would use `GET /users/{uid}/passports`.

When we want to retrieve an specific passport, we don't need to prefix the route with `/users/{uid}` anymore because we know exactly which passport we want to retrieve. So, instead of `GET /users/{uid}/passports/{pid}`, we can just use `GET /passports/{pid}`.

Once you have the API design sorted, it's just a matter of creating the code that gets called when a specific route is hit. We implement those with Handlers and those routes are defined in `routes.go`:

```
Route{"Healthcheck", "GET", "/healthcheck", HealthcheckHandler},
//=== USERS ===
Route{"ListUsers", "GET", "/users", ListUsersHandler},
Route{"GetUser", "GET", "/users/{uid:[0-9]+}", GetUserHandler},
Route{"CreateUser", "POST", "/users", CreateUserHandler},
Route{"UpdateUser", "PUT", "/users/{uid:[0-9]+}", UpdateUserHandler},
Route{"DeleteUser", "DELETE", "/users/{uid:[0-9]+}", DeleteUserHandler},
```

In order to make our code more robust, I've added pattern matching in the routes.  This `[0-9]+` pattern says that we only accept digits from 0 to 9 and we can have one or more digits. Everything else will most likely trigger an HTTP 404 Not Found status code being returned to the client.

### API Handlers

Most Go code that show HandleFunc examples, will show something slightly different, something more like this:

```
router.HandleFunc("/users", ListUsersHandler).Methods("GET")
```

So Go's HandleFunc function takes two arguments, one string that defines the route and a second argument of the type `http.HandleFunc(w http.ResponseWriter, r *http.Request)`. This works very well, but doesn't allow you to pass any extra data to your handlers. So for instance when you want to use the `unrolled/render` package, then you'd have to define in your `main.go` file a global variable like this: `var Render  *render.Render`, then initialise that in your `func main()` so that your handlers can access this global variable later on. Using global variables in Go is not a good practice (there are exceptions like certain database drivers, but that's a different discussion).

So our initial handler function for returning a list of users was:

```
func ListUsersHandler(w http.ResponseWriter, req *http.Request) {
  // return user data
}
```

Where the `render` variable is a global variable. We'd prefer to pass the Render variable to that function so will rewrite it to:

```
func ListUsersHandler(w http.ResponseWriter, req *http.Request, render Render) {
  // return user data
}
```

Perfect. Or, is it? What if we want to pass more variables to the handler function? Like Metrics? Or some environment variables? We'd continuously have to change ALL our handlers and the type signature will become quite long and hard to maintain. An alternative is to create a server or an environment struct and use that as a container for our variables we want to pass around the system.

In our `internal/passport/appenv.go` we'll add:

```
// AppEnv holds application configuration data
type AppEnv struct {
	Render    *render.Render
	Version   string
	Env       string
	Port      string
	UserStore models.UserStorage
}
```

And in our `func main()` function (in `cmd/api-service/main.go`) we initialise that struct:

```
// initialse application context
appEnv := passport.AppEnv{
	Render:    render.New(),
	Version:   version,
	Env:       env,
	Port:      port,
	UserStore: userStore,
}
```

Our handler looks like this now:

```
func ListUsersHandler(w http.ResponseWriter, req *http.Request, appEnv AppEnv) {
  // return user data
}
```

The only problem is that this handler's type signature is not `http.ResponseWriter, *http.Request` but `http.ResponseWriter, *http.Request, appContext` so Go's HandleFunc function will complain about this. That's why we are introducing a helper function `MakeHandler` that takes our appContext struct and our handlers with the special type signature and converts it to `func(w http.ResponseWriter, r *http.Request)`:

```
func MakeHandler(appEnv AppEnv, fn func(http.ResponseWriter, *http.Request, AppEnv)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fn(w, r, appEnv)
	}
}
```

We can use this `MakeHandler` function now in our `handlers.go` when we initialise all the routes and the corresponding route handlers.

In the previous examples we've ignored the implementation of the `ListUsersHandler`, but let's have a closer look at that now:

```
// ListUsersHandler returns a list of users
func ListUsersHandler(w http.ResponseWriter, req *http.Request, appEnv AppEnv) {
	list, err := appEnv.UserStore.ListUsers()
	if err != nil {
		response := status.Response{
			Status:  strconv.Itoa(http.StatusNotFound),
			Message: "can't find any users",
		}
		log.WithFields(log.Fields{
			"env":    appEnv.Env,
			"status": http.StatusNotFound,
		}).Error("Can't find any users")
		appEnv.Render.JSON(w, http.StatusNotFound, response)
		return
	}
	responseObject := make(map[string]interface{})
	responseObject["users"] = list
	responseObject["count"] = len(list)
	appEnv.Render.JSON(w, http.StatusOK, responseObject)
}
```

After we retrieve the data from our database, we check whether there are errors in the response. Rather than passing on the response to the user, we create a new custom `Status` struct with a new bespoke error message. Why? That's to avoid that we're leaking sensitive information from our systems to the user. It's the equivalent of sending back a Java stacktrace to a user. It can disclose information to people that are trying to hack into your system by exposing some of the weaknesses.

So this code:

```
if err != nil {
	response := status.Response{
		Status:  strconv.Itoa(http.StatusNotFound),
		Message: "can't find any users",
	}
	log.WithFields(log.Fields{
		"env":    appEnv.Env,
		"status": http.StatusNotFound,
	}).Error("Can't find any users")
	appEnv.Render.JSON(w, http.StatusNotFound, response)
	return
}
```

constructs a new `Status` response struct. Adds a 404 status so that computers can understand it, but also adds a human-readable message. We need to convert the `http.StatusNotFound` from an integer to a string using `strconv.Itoa` because the status codes in the `Status` object are strings. It then prints the error to the server log. Last but not least you send the `Status` struct as a JSON response back to the client and add the 404 Not Found response to the HTTP headers (using `http.StatusNotFound`).


### Special Route: Health-check

When you look at the overview of the handlers in `routes.go`, you will notice a special route:

```
Route{"Healthcheck", "GET", "/healthcheck", HealthcheckHandler},
```

Monitoring tools like [Sensu](https://sensuapp.org/) can call: `GET /healthcheck`. The health check route can return a 200 OK when the service is up and running and will also say what the application name is and the version number.

```
// HealthcheckHandler returns useful info about the app
func HealthcheckHandler(w http.ResponseWriter, req *http.Request, appEnv AppEnv) {
	check := health.Check{
		AppName: "go-rest-api-template",
		Version: appEnv.Version,
	}
	appEnv.Render.JSON(w, http.StatusOK, check)
}
```

This health check is very simple. It just checks whether the service is up and running, which can be useful in a build and deployment pipelines where you can check whether your newly deployed API is running (as part of a smoke test). More advanced health checks will also check whether it can reach the database, message queue or anything else you'd like to check. Trust me, your platform/infra colleagues will be very grateful for this. 

### Returning Data

Let's have a look at interacting with our data. Returning a list of users is quite easy, it's just showing the UserList:

```
func ListUsersHandler(w http.ResponseWriter, req *http.Request, appEnv AppEnv) {
	list, err := appEnv.UserStore.ListUsers()
	if err != nil {
		response := status.Response{
			Status:  strconv.Itoa(http.StatusNotFound),
			Message: "can't find any users",
		}
		log.WithFields(log.Fields{
			"env":    appEnv.Env,
			"status": http.StatusNotFound,
		}).Error("Can't find any users")
		appEnv.Render.JSON(w, http.StatusNotFound, response)
		return
	}
	responseObject := make(map[string]interface{})
	responseObject["users"] = list
	responseObject["count"] = len(list)
	appEnv.Render.JSON(w, http.StatusOK, responseObject)
}
```

BTW, did you notice the `appEnv.render.JSON`? That's part of `"github.com/unrolled/render"` and allows us to render JSON output when we send data back to the client.

So, this will return the following to the client:

```
{
    "users": [
        {
            "dateOfBirth": "1992-01-01T00:00:00Z",
            "firstName": "Jane",
            "id": 1,
            "lastName": "Doe",
            "locationOfBirth": "Milton Keynes"
        },
        {
            "dateOfBirth": "1985-12-31T00:00:00Z",
            "firstName": "John",
            "id": 0,
            "lastName": "Doe",
            "locationOfBirth": "London"
        }
    ]
}
```

It may surprise you that we are returning a JSON object that holds an array with multiple JSON objects, rather than an array with multiple JSON objects, as seen in the example below:

```
{
    [
        {
            "dateOfBirth": "1992-01-01T00:00:00Z",
            "firstName": "Jane",
            "id": 1,
            "lastName": "Doe",
            "locationOfBirth": "Milton Keynes"
        },
        {
            "dateOfBirth": "1985-12-31T00:00:00Z",
            "firstName": "John",
            "id": 0,
            "lastName": "Doe",
            "locationOfBirth": "London"
        }
    ]
}
```

They're both valid and there are lots of views and opinions (as always in developer / architecture communities!), but the reason why I prefer to wrap the array in a JSON object is because later on we can easily add more data without causing significant changes to the client. What if we want to add the concept of pagination to our API?

Example:

```
{
    "offset": 0,
    "limit":  25,
    "users":  [
        {
            "dateOfBirth": "1992-01-01T00:00:00Z",
            "firstName": "Jane",
            "id": 1,
            "lastName": "Doe",
            "locationOfBirth": "Milton Keynes"
        },
        {
            "dateOfBirth": "1985-12-31T00:00:00Z",
            "firstName": "John",
            "id": 0,
            "lastName": "Doe",
            "locationOfBirth": "London"
        }
    ]
}
```

In this project, we're adding an extra field named `count` so we tell the client how many items they can expect

```
{
    "count": 2,
    "users": [
        {
            "dateOfBirth": "1985-12-31T00:00:00Z",
            "firstName": "John",
            "id": 0,
            "lastName": "Doe",
            "locationOfBirth": "London"
        },
        {
            "dateOfBirth": "1992-01-01T00:00:00Z",
            "firstName": "Jane",
            "id": 1,
            "lastName": "Doe",
            "locationOfBirth": "Milton Keynes"
        }
    ]
}
```

It's important that we use the `interface{}` type in the map, rather than a specific `[]User` type. Otherwise, we won't be able to assign whatever data we want to the map:

```
responseObject := make(map[string]interface{})
responseObject["users"] = list
responseObject["count"] = len(list)
appEnv.Render.JSON(w, http.StatusOK, responseObject)
```

Another example is the retrieval of a specific object:

```
// GetUserHandler returns a user object
func GetUserHandler(w http.ResponseWriter, req *http.Request, appEnv AppEnv) {
	vars := mux.Vars(req)
	uid, _ := strconv.Atoi(vars["uid"])
	user, err := appEnv.UserStore.GetUser(uid)
	if err != nil {
		response := status.Response{
			Status:  strconv.Itoa(http.StatusNotFound),
			Message: "can't find user",
		}
		log.WithFields(log.Fields{
			"env":    appEnv.Env,
			"status": http.StatusNotFound,
		}).Error("Can't find user")
		appEnv.Render.JSON(w, http.StatusNotFound, response)
		return
	}
	appEnv.Render.JSON(w, http.StatusOK, user)
}
```

This reads the uid (User ID) variable from the route (`/users/{uid}`), converts the string to an integer and then looks up the user in our UserList by ID. If the user does not exit, we return a 404 and an error object. If the user exists, we return a 200 and a JSON object with the user.

Example:

```
{
    "dateOfBirth": "1992-01-01T00:00:00Z",
    "firstName": "Jane",
    "id": 1,
    "lastName": "Doe",
    "locationOfBirth": "Milton Keynes"
}
```

We add the attributes of the user object to the root of the JSON response, rather than wrapping it up in an explicit JSON object. I can quite easily add extra data to the response, without breaking the existing data.

## Testing

### Manually testing your API routes with curl commands

Let's start with some simple curl tests. Open your terminal and try the following curl commands.

Retrieve a list of users:

```
curl -X GET http://localhost:3001/users | python -mjson.tool
```

That should result in the following result:

```
  % Total    % Received % Xferd  Average Speed   Time    Time     Time  Current
                                 Dload  Upload   Total   Spent    Left  Speed
100   246  100   246    0     0  14470      0 --:--:-- --:--:-- --:--:-- 15375
{
    "count": 2,
    "users": [
        {
            "id": 0,
            "firstName": "John",
            "lastName": "Doe",
            "dateOfBirth": "1985-12-31T00:00:00Z",
            "locationOfBirth": "London"
        },
        {
            "id": 1,
            "firstName": "Jane",
            "lastName": "Doe",
            "dateOfBirth": "1992-01-01T00:00:00Z",
            "locationOfBirth": "Milton Keynes"
        }
    ]
}
```

The `| python -mjson.tool` at the end is for pretty printing (formatting). It essentially tells to pipe the output of the curl command to the SJON formatting tool. If we only typed `curl -X GET http://localhost:3009/users` then we'd have something like this:

```
{"users":[{"id":0,"firstName":"John","lastName":"Doe","dateOfBirth":"1985-12-31T00:00:00Z","locationOfBirth":"London"},{"id":1,"firstName":"Jane","lastName":"Doe","dateOfBirth":"1992-01-01T00:00:00Z","locationOfBirth":"Milton Keynes"}]}
```

So not that easy to read as the earlier nicely formatted example.

Get a specific user:

```
curl -X GET http://localhost:3001/users/0 | python -mjson.tool
```

Results in:

```
  % Total    % Received % Xferd  Average Speed   Time    Time     Time  Current
                                 Dload  Upload   Total   Spent    Left  Speed
100   104  100   104    0     0   6625      0 --:--:-- --:--:-- --:--:--  6933
{
    "dateOfBirth": "1992-01-01T00:00:00Z",
    "firstName": "Jane",
    "id": 1,
    "lastName": "Doe",
    "locationOfBirth": "Milton Keynes"
}
```

### Testing the Database

There are lots of opinions on testing, how much you should be testing, which layers of your applications, etc. When I'm working with API services, I tend to focus on two types of tests to start with: testing the data access layer and testing the actual HTTP service.

In this example, we want to test the List, Add, Get, Update and Delete operations on our in-memory document database. The data access code is stored in the `internal/passport/db_user.go` file, so following Go convention we will create a new file called `internal/passport/db_user_test.go`.

Before we can run the tests in `db_user_test.go` we need to make sure that we have a valid application context with a mock database initialised. I've added a helper function in `internal/passport/appenv.go` named `CreateContextForTestSetup` and that takes care of both the database setup as well as the AppEnv.

```
// CreateContextForTestSetup initialises an application context struct
// for testing purposes
func CreateContextForTestSetup() AppEnv {
	testVersion := "0.0.0"
	appEnv := AppEnv{
		Render:    render.New(),
		Version:   testVersion,
		Env:       "LOCAL",
		Port:      "3001",
		UserStore: NewUserService(CreateMockDataSet()),
	}
	return appEnv
}
```

We then have a helper function in `internal/passport/db_user.go` that sets up a mock data set:

```
// CreateMockDataSet initialises a database for test purposes
func CreateMockDataSet() map[int]models.User {
	list := make(map[int]models.User)
	dt, _ := time.Parse(time.RFC3339, "1985-12-31T00:00:00Z")
	list[0] = models.User{
		ID:              0,
		FirstName:       "John",
		LastName:        "Doe",
		DateOfBirth:     dt,
		LocationOfBirth: "London",
	}
	dt, _ = time.Parse(time.RFC3339, "1992-01-01T00:00:00Z")
	list[1] = models.User{
		ID: 1, FirstName: "Jane",
		LastName:        "Doe",
		DateOfBirth:     dt,
		LocationOfBirth: "Milton Keynes",
	}
	return list
}
```

Let's have a look at the easiest test where we list the elements in our database:

```
func TestListUsers(t *testing.T) {
	appEnv := CreateContextForTestSetup()
	list, _ := appEnv.UserStore.ListUsers()
	count := len(list)
	assert.Equal(t, 2, count, "There should be 2 items in the list.")
}
```

This first calls `CreateContextForTestSetup()` the create the application context and mock database, then the `appEnv.UserStore.ListUsers()` function, which returns a list of users. We then count the number of elements and last but not least we then check whether that count equals 2.

In standard Go, you would actually write something like:

```
if 2 != count {
  t.Errorf("Expected 2 elements in the list, instead got %v", count)
}
```

However there is a neat Go package called [testify](https://github.com/stretchr/testify) that gives you assertions like Java and that's why we can write cleaner test code like:

```
assert.Equal(t, 2, count, "There should be 2 items in the list.")
```

The `TestList` is only testing for a positive result, but we really need to test for failures as well.

This is our test code for the Delete functionality:

```
func TestDeleteUserSuccess(t *testing.T) {
	appEnv := CreateContextForTestSetup()
	err := appEnv.UserStore.DeleteUser(1)
	assert.Nil(t, err)
}

func TestDeleteUserFail(t *testing.T) {
	appEnv := CreateContextForTestSetup()
	err := appEnv.UserStore.DeleteUser(10)
	assert.NotNil(t, err)
}
```

The first test function `TestDeleteSuccess` tries to delete a known existing user, with id 1. We're expecting that the error object is Nil. The second test function `TestDeleteFail` tries to look up a non-existing user with id 10, and as expected, this should return an actual Error object.

How do we run the tests?

Simple. Open your terminal, go to your project root:

```
go test ./...
```

We need to add `./...` so that it also runs all the tests in the subpackages. If you want it more verbose, then:

```
go test ./... -v
```

Which will give you:

```
?   	github.com/leeprovoost/go-rest-api-template/cmd/api-service	[no test files]
=== RUN   TestListUsers
--- PASS: TestListUsers (0.00s)
=== RUN   TestGetUserSuccess
--- PASS: TestGetUserSuccess (0.00s)
=== RUN   TestGetUserFail
--- PASS: TestGetUserFail (0.00s)
=== RUN   TestAddUser
--- PASS: TestAddUser (0.00s)
=== RUN   TestUpdateUserSuccess
--- PASS: TestUpdateUserSuccess (0.00s)
=== RUN   TestUpdateUserFail
--- PASS: TestUpdateUserFail (0.00s)
=== RUN   TestDeleteUserSuccess
--- PASS: TestDeleteUserSuccess (0.00s)
=== RUN   TestDeleteUserFail
--- PASS: TestDeleteUserFail (0.00s)
=== RUN   TestHealthcheckHandler
--- PASS: TestHealthcheckHandler (0.00s)
=== RUN   TestListUsersHandler
--- PASS: TestListUsersHandler (0.00s)
PASS
ok  	github.com/leeprovoost/go-rest-api-template/internal/passport	0.222s
?   	github.com/leeprovoost/go-rest-api-template/internal/passport/models	[no test files]
?   	github.com/leeprovoost/go-rest-api-template/pkg/health	[no test files]
?   	github.com/leeprovoost/go-rest-api-template/pkg/status	[no test files]
?   	github.com/leeprovoost/go-rest-api-template/pkg/version	[no test files]
```

Do you want to get some more info on your code coverage? No worries, Go has you covered (no pun intended):

```
go test ./... -cover
```

This will give you:

```
?   	github.com/leeprovoost/go-rest-api-template/cmd/api-service	[no test files]
ok  	github.com/leeprovoost/go-rest-api-template/internal/passport	0.165s	coverage: 38.3% of statements
?   	github.com/leeprovoost/go-rest-api-template/internal/passport/models	[no test files]
?   	github.com/leeprovoost/go-rest-api-template/pkg/health	[no test files]
?   	github.com/leeprovoost/go-rest-api-template/pkg/status	[no test files]
?   	github.com/leeprovoost/go-rest-api-template/pkg/version	[no test files]
```

## Vendoring of dependencies

Before we wrap up, I'd like to briefly talk about vendoring. Most software relies on some third-party (open source) software. In the Node.js world they are called npm packages, in the Java world they are Maven packages and in Go they are just called a Go package.

What you commonly see in the Node.js world is that people push their own code to a server and then run `npm install` on that server. There have been a few instances where some npm packages disappeared and that caused a lot of problems for people who all of a sudden couldn't install their software anymore.

One way to solve it is to simply store a copy of the third-party dependency into your own project and treat it as you would treat the code you've written yourself. They're part of the git project and checked in together. Since Go 1.5, the Go world has agreed to use the `/vendor` folder in your project for this, as well as introduced a common go modules approach. 

If you want to add a package to your `/vendor` folder, then just run `go get your/desired/package`

This project has already all its dependencies vendored into the `/vendor` folder, so no need to do it again.

When you want to upgrade a module:

```
go get -u github.com/unnrolled/render
go mod tidy
go mod vendor
```

Side note: if you are upgrading from a pre-go mod project (in my case govendor) then that's quite easy. Enter the following commands in the root of your project step by step:

```
go mod init github.com/leeprovoost/go-rest-api-template
go mod tidy
rm ./vendor/vendor.json
go mod vendor
```

This is documented [here](https://blog.golang.org/migrating-to-go-modules) and the last `go mod vendor` is mentioned [here](https://github.com/golang/go/issues/37734#issuecomment-596695647).

## Starting the app on a production server

This is how you could run your app on a server:

First, you copy the binary and the `VERSION` file into a directory, e.g. `/opt/go-rest-api-template`.

Then start the app as a service. Store the app's PID in a text file so we can kill it later.

```
#!/bin/bash
export ENV=DEV
export PORT=8080
export VERSION=/opt/go-rest-api-template/VERSION
sudo nohup /opt/go-rest-api-template/api-service >> /var/log/go-rest-api-template.log 2>&1&
echo $! > /opt/go-rest-api-template/api-service-pid.txt
```

When you want to kill your app later during a redeployment or a server shutdown, then you can kill the app by looking up the previously stored PID:

```
#!/bin/bash
if [ -f /opt/go-rest-api-template/api-service-pid.txt ]; then
  kill -9 `cat /opt/go-rest-api-template/api-service-pid.txt`
  rm -f /opt/go-rest-api-template/api-service-pid.txt
fi
```

## Useful references

### General

* [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
* [Structuring applications in Go](https://medium.com/@benbjohnson/structuring-applications-in-go-3b04be4ff091)
* [Writing modular GO REST APIs](http://thenewstack.io/make-a-restful-json-api-go/)

### HTTP, REST and JSON

* [Structs and JSON formatting](http://stackoverflow.com/questions/21825322/why-golang-cannot-generate-json-from-struct-with-front-lowercase-character)
* [JSON and Go](http://blog.golang.org/json-and-go)
* [Design beautiful REST + JSON APIs](http://www.slideshare.net/stormpath/rest-jsonapis)
* [Use render for generating JSON](https://github.com/unrolled/render/issues/7) for use of global variable
* [Read JSON POST body](http://stackoverflow.com/questions/15672556/handling-json-post-request-in-go)
* [How to pass a parameter to a Http handler function](https://groups.google.com/forum/#!topic/golang-nuts/SGn1gd290zI)
* Go and datetime parsing/formatting: [ISO 8601, the International Standard for the representation of dates and times](http://www.w3.org/TR/NOTE-datetime), [Go by Example: Time Formatting / Parsing](https://gobyexample.com/time-formatting-parsing), [JSON datetime formatting](http://stackoverflow.com/a/15952652), [src/time/format.go](http://golang.org/src/time/format.go)

### Testing

* [Testing techniques](https://talks.golang.org/2014/testing.slide#1)
* [Testing Go HTTP API](http://dennissuratna.com/testing-in-go/)
* [Great overview of HTTP response codes](http://stackoverflow.com/a/2342631)

### Go core language concepts

* [Understanding method receivers and pointers](http://nathanleclaire.com/blog/2014/08/09/dont-get-bitten-by-pointer-vs-non-pointer-method-receivers-in-golang/)
* [HTTP Closures gist](https://gist.github.com/tsenart/5fc18c659814c078378d)
* [Introducing Function Literals and Closures](https://golang.org/doc/articles/wiki/)
* [Custom Handlers and Avoiding Globals in Go Web Applications](https://elithrar.github.io/article/custom-handlers-avoiding-globals/)
