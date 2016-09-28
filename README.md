# go-rest-api-template [ ![Codeship Status for leeprovoost/go-rest-api-template](https://codeship.com/projects/89ed2300-9d8f-0133-5815-1a74f7994c2d/status?branch=master)](https://codeship.com/projects/127524)

*WORK IN PROGRESS: documentation is not yet on par with a recent major refactoring exercise (28 September 2016)*

Template for building REST Web Services in Golang. Uses gorilla/mux as a router/dispatcher and Negroni as a middleware handler. Tested against Go 1.7.1.

## Introduction

### Why?

After writing many REST APIs with Java Dropwizard, Node.js/Express and Go, I wanted to distill my lessons learned into a reusable template for writing REST APIs, in the Go language.

It's mainly for myself. I don't want to keep reinventing the wheel and just want to get the foundation of my REST API 'ready to go' so I can focus on the business logic and integration with other systems and data stores.

Just to be clear: this is not a framework, library, package or anything like that. This tries to use a couple of very good Go packages and libraries that I like and cobbled them together.

The main ones are:

* [gorilla/mux](http://www.gorillatoolkit.org/pkg/mux) for routing
* [codegangsta/negroni](https://github.com/codegangsta/negroni) as a middleware handler
* [strechr/testify](https://github.com/stretchr/testify) for writing easier test assertions
* [unrolled/render](https://github.com/unrolled/render) for HTTP response rendering
* [palantir/stacktrace](https://github.com/palantir/stacktrace) to provide more context to error messages
* [unrolled/secure](https://github.com/unrolled/secure) to improve API security

Whilst working on this, I've tried to write up my thought process as much as possible. Everything from the design of the API and routes, some details of the Go code like JSON formatting in structs and my thoughts on testing. However, if you feel that there is something missing, send a PR or raise an issue.

Last but not least: I'm occasionally writing about my lessons learned developing Go applications and Amazon Web Services systems on [Medium](https://medium.com/@leeprovoost).

### Knowledge of Go

If you're new to programming in Go, I would highly recommend you to read the following three resources:

* [A Tour of Go](https://tour.golang.org/welcome/1)
* [Effective Go](https://golang.org/doc/effective_go.html)
* [50 Shades of Go: Traps, Gotchas, and Common Mistakes for New Golang Devs](http://devs.cloudimmunity.com/gotchas-and-common-mistakes-in-go-golang/)
* [Learn X in Y Minutes](https://learnxinyminutes.com/docs/go/)

You can work your way through those in two to three days. I did try out a couple of books but didn't find them that useful. If you really want to get some books, there is a [github](https://github.com/dariubs/GoBooks) repo that tries to list the main ones.

A good project to keep on eye on to discover new Go packages and software is [awesome-go](https://github.com/avelino/awesome-go). The maintainers enforce strict standards around documentation, tests, structure, etc. so they are doing a brilliant job improving the overall quality of code in the Go community. It's not easy to keep track of the high rate of change on that project, so my suggestion would be to track new additions using the [changedetection.com](https://www.changedetection.com/) website.

### Development tools

I've tried many different editors and it seems like it's a common Go newbie [frustration](https://groups.google.com/forum/#!topic/golang-nuts/6ZgrZsPzHr0). Coming from Java, I've tried "proper" IDEs like [LiteIDE](https://github.com/visualfc/liteide) and IntelliJ with the [golang plugin](https://github.com/go-lang-plugin-org/go-lang-idea-plugin) but never really fell in love with them.

I've used Sublime Text with the [GoSublime](https://github.com/DisposaBoy/GoSublime) plugin for a long time, but have now settled on Atom with the [go-plus](https://github.com/joefitzgerald/go-plus) plugin and that is working very well. It has code completion, linting, formatting and even shows you what parts of the code are covered by tests.

The choice of an editor/IDE is a personal one, so the only advice I can give you is to try different ones. It looks like currently the following seem to be the most mature:

* Sublime Text 3 with [GoSublime](https://github.com/DisposaBoy/GoSublime). I've been made aware that there is another Go plugin for Sublime Text called [GoTools](https://packagecontrol.io/packages/GoTools)
* Atom with [go-plus](https://github.com/joefitzgerald/go-plus). There is a good article describing a full Go setup [here](http://marcio.io/2015/07/supercharging-atom-editor-for-go-development/).
* vim with [vim-go](https://github.com/fatih/vim-go)
* [LiteIDE](https://github.com/visualfc/liteide)
* Microsoft did a pretty good job with their [Visual Studio Code](https://code.visualstudio.com/) and Luke Hoban created a brilliant Go binding for it: [vscode-go](https://github.com/Microsoft/vscode-go)
* IntelliJ with the [Go plugin](https://github.com/go-lang-plugin-org/go-lang-idea-plugin)

Last but not least, I can highly recommend trying out [GitKraken](https://www.gitkraken.com/). I've never been a huge fan of desktop UI clients for git, but this one sets a new bar. Slick UI and very useful information about your local and remote branches and code history.

### How to run

`go run main.go` works fine if you have a single file or a script you're working on, but once you have a more complex project with lots of files then you'll have to start using the proper go build tool and run the compiled executable.

```
go build && ./go-rest-api-template
```

Which is faster than first typing `go build` to generate an executable called `go-rest-api-template` and then run that executable by typing `./go-rest-api-template`.

The `go-rest-api-template` app will bind itself to port 3001 (defined in `main.go`). If you want to change it (e.g. bind it to the default http port 80), then use the following environment variable `PORT` (see `main.go`). Same for the location of the `fixtures.json` model.

```
export PORT=80
export FIXTURES=/tmp/fixtures.json
go build && ./go-rest-api-template
```

Alternatively, you can just change line 23 of `main.go` and change the hardcoded port to whichever one you prefer.

One comment from a security point of view: whilst 80 is the default port for HTTP (and 443 for HTTPS), it's usually not recommended to actually bind your application to that port when it runs on a server instance like AWS EC2. That's because ports below 1024 are privileged ports and require elevated privileges for the app to bind itself to such port. If someone compromises your app, it could potentially take advantage of the application user with the elevated privileges. Instead, run it on a higher port like 8001 and run it with a newly created restricted user (instead of using the root user). You can then use a load balancer or a proxy (e.g. AWS ELB, HAProxy, nginx) to accept your incoming request from the internet on port 80 or 443 and then forward that to your application port.

### Live Code Reloading

Manually stopping and restarting your server can get quite annoying after a while, so let's set up a task runner that automatically restarts the server when it detects changes.

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
build_name:        go-rest-api-template
build_log:         go-rest-api-template.log
valid_ext:         .go, .tpl, .tmpl, .html
build_delay:       600
colors:            1
log_color_main:    cyan
log_color_build:   yellow
log_color_runner:  green
log_color_watcher: magenta
log_color_app:
```

As you can see, it creates a `tmp` directory in your project root and a log file. You can tell `.gitignore` to stop checking it into your git repository by adding the following lines in your `.gitignore` file:

```
# fresh
tmp
go-rest-api-template.log
```

Where I work, we've hit the limitations of fresh and have switched to Gulp as a task runner. Our Gulp setup takes care of running tests, using OSX notifications to warn you of build failures and provides an overview of your `govendor` package status. I'll open source our Gulp for Go configuration at some point.

### Code Structure

Main server files (bootstrapping of http server):

```
main.go
server.go
```

Route definitions and handlers:

```
router.go
handlers.go
handlers_test.go
```

Data model descriptions and operations on the data:

```
models.go
database.go
database_test.go
```

Test data:

```
fixtures.json
```

Configuration file for [fresh](go get github.com/pilu/fresh):

```
runner.conf
```

Test that checks whether structs comply with the DataStorer interface:

```
interface_test.go
```

Helper structs and functions:

```
helpers.go
```

Defines application version using semantic versioning:

```
VERSION
```

Vendored go packages using govendor:

```
vendor/
```

### Starting the application: main.go and server.go

The main entry point to our app is `main.go` where we take care of the following:

Loading our environment variables, but also provide some default values that are useful when you run it on your local development machine:

```
var (
  // environment variables
  env      = os.Getenv("ENV")      // LOCAL, DEV, STG, PRD
  port     = os.Getenv("PORT")     // server traffic on this port
  version  = os.Getenv("VERSION")  // path to VERSION file
  fixtures = os.Getenv("FIXTURES") // path to fixtures file
)
if env == "" || env == local {
  // running from localhost, so set some default values
  env = local
  port = "3001"
  version = "VERSION"
  fixtures = "fixtures.json"
}
```

I'm using the environment variables technique to tell the app in what environment it lives:

* `LOCAL`: local development machine
* `DEV`: development or integration server
* `STG`: staging servers
* `PRD`: production servers

When it can't detect an environment variable, it assumes that it's running on your local development workstation.

We then load our version file. I have added a helper function in `helpers.go` that can parse a `VERSION` file using semantic versioning.

```
// reading version from file
version, err := ParseVersionFile(version)
if err != nil {
  log.Fatal(err)
}
```

The check for error (`if err != nil`) is a common Go way of handling return values that contain errors. the `main` function is one of the few places where I use log.Fatal in my application. `log.Fatal` logs your file to the console, but also exits your application. This is for me the correct behaviour here because without a correct `VERSION` file, the app shouldn't work. It's too risky that you may have deployed the incorrect application version.

This API doesn't talk to a real database but to a very simple in-memory mock database. That's why I have a `fixtures.json` file with some data. I have added a helper function in `helpers.go` that can read and parse that file into our in-memory mock database. Have a look at the next Data section to find out more about the in-memory database. I may change this in the future into an integration with an actual database.

```
// load fixtures data into mock database
db, err := LoadFixturesIntoMockDatabase(fixtures)
```

Once all the data is read from our environment variables and `fixtures.json` file, we can then initialise our application context. `AppContext` is a struct that holds some application info that we can pass around (e,g, provide to our handlers). It's a neat way to avoid using global variables and avoiding that you have application variables all over the place.

I have defined the `AppContext` struct in `helpers.go`:

```
// AppContext holds application configuration data
type AppContext struct {
	Render  *render.Render
	Version string
	Env     string
	Port    string
	DB      DataStorer
}
```

`Render` helps us with rendering the correct response to the client, e.g. JSON, XML, etc. Version holds the path to our `VERSION` file. `Env` holds information about our platform environment (e.g. STG, DEV). `Port` is the server port that our application binds to. `DB` is our database struct that provides a data abstraction layer.

Once all of that is set up properly, we can now start the server. The `StartServer` function takes our `AppContext` struct and initialises all our routes and starts our Negroni server. This is all defined in `server.go` and we'll discuss the details over the next few sections.

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

We are going to use a travel Passport for our example. I've chosen Id as the unique key for the passport because (in the UK), passport book numbers these days have a unique 9 character field length (e.g. 012345678). A passport belongs to a user and a user can have one or more passports. We'll define this in `models.go`.

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

If you want to prevent a certain struct field to be marshalled/unmarshalled then add `json:"-"`.

### Operations on our (mock) database

I wanted to create a template REST API that didn't depend on a database, so started with a simple in-memory database that we can work with. The good thing is that this will be the start of a so-called data access layer that abstracts away the underlying data store. We can achieve that by starting with creating an interface (which is a good practice in Go anyway). Note the use of the -er at the end of the interface name, as per Go convention.

Open the `database.go` file and you will see:

```
type DataStorer interface {
	ListUsers() ([]User, error)
	GetUser(i int) (User, error)
	AddUser(u User) (User, error)
	UpdateUser(u User) (User, error)
	DeleteUser(i int) error
}
```

This allows us to define a set of operations on the data as a contract, without people having to worry about the actual implementation of how the data is stored and accessed. I've added the basic operations to list, retrieve, create, update and delete data, so the standard CRUD-style operations (accepting that CRUD has some subtle differences with REST).

Let's have a look at the type signature of the `Get` operation:

```
GetUser(i int) (User, error)
```

What this tells us is that it is expecting an integer as an argument (which will be the User id in our case), and returns a pair of values: a user object and an error object. Returning pairs of values is a nice Go feature and is often used to return information about errors.

An example of how this could be used is the following:

```
user, err := ctx.db.GetUser(uid)
if err != nil {
  ctx.render.JSON(w, http.StatusNotFound, err)
  return
}
ctx.render.JSON(w, http.StatusOK, user)
```

We check whether the error object is nil. If it is, then we return a HTTP 200 OK, if not then we return HTTP 404 NOT FOUND. Let's go into more detail when we talk about our API handlers.

BTW in other languages you might be used to write the above as:

```
user, err := ctx.db.GetUser(uid)
if err == nil {
  ctx.render.JSON(w, http.StatusOK, user)
} else {
  ctx.render.JSON(w, http.StatusNotFound, err)
}
```

But the Go way is to check if an error exists and handle that immediately. If there is no error, just continue with the normal application flow.

Let's have a look at the actual mock in-memory database. We need to create a Database struct that will hold the data:

```
type MockDB struct {
	UserList  map[int]User
	MaxUserID int
}
```
The UserList will hold a list of User structs and the MaxUserID holds the latest used integer. MaxUserID mimics the behaviour of an auto-generated ID in conventional databases. In this case MaxUserID represents the highest used database ID.

Let's have a closer at the `GetUser` function:

```
// GetUser returns a single JSON document
func (db *MockDB) GetUser(i int) (User, error) {
	user, ok := db.UserList[i]
	if !ok {
		return User{}, stacktrace.NewError("Failure trying to retrieve user")
	}
	return user, nil
}
```

* `func (db *MockDB)`: TODO
* `GetUser(i int) (User, error)`: TODO
* `return User{}, stacktrace.NewError("Failure trying to retrieve user")`: TODO
* `return user, nil`: TODO

When I discussed the `DataStorer` interface, I didn't touch on a very important use case: testing.

TODO

### Fixtures

In order to make it a bit more useful, we will initialise it with some user objects. Luckily, we can make use of the `init` function that gets automatically called when you start the application. This init() function will be in our `main.go` file when you start up the server:

```
func init() {
  // read JSON fixtures file
  var jsonObject map[string][]User
  file, err := ioutil.ReadFile("./fixtures.json")
  if err != nil {
    log.Fatalf("File error: %v\n", err)
  }
  err = json.Unmarshal(file, &jsonObject)
  if err != nil {
    log.Fatal(err)
  }
  // load data in database
  list := make(map[int]User)
  list[0] = jsonObject["users"][0]
  list[1] = jsonObject["users"][1]
  db = &Database{list, 1}
}
```

We are first going to load the data from a `fixtures.json` file:

```
{
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

When we can't load the file, we will stop the bootstrapping of the application. This is taken care of by Go's log handler, which fires off a fatal error:

```
file, err := ioutil.ReadFile("./fixtures.json")
if err != nil {
  log.Fatalf("File error: %v\n", err)
}
```

The `fixtures.json` file contains a JSON representation of a Go map where the key is a string (i.e. `"users"`) and the map value is a string of User objects. In Go, this would be represented as: `map[string][]User`. We load the fixtures file, marshal it into the type we just defined and then load it into our database.

The date string looks a bit odd. Why not just use `31-12-1985` or `1985-12-31`? The first is discouraged altogether because that's an European way of writing dates and will cause confusion around the world. Not this particular example, but imagine you have 3-4-2015. Is it third of April or fourth of March? Unfortunately there isn't an "enforced standard" for dates in JSON, so I've tried to use one that is commonly used and also understood by Go's `json.Marshaler` and `json.Unmarshaler` to avoid that we have to write our own custom marshaler/unarshaler.

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

We now need to implement the various methods from our DataStore interface.

## Defining the API

### API Routes

Now that we have defined the data access layer, we need to translate that to a REST interface:

* Retrieve a list of all users: `GET /users` -> The `GET` just refers to the HTTP action you would use. If you want to test this in the command line, then you can use curl: `curl -X GET http://localhost:3009/users` or `curl -X POST http://localhost:3009/users`
* Retrieve the details of an individual user: `GET /users/{uid}` -> {uid} allows us to create a variable, named uid, that we can use in our code. An example of this url would be `GET /users/1`
* Create a new user: `POST /users`
* Update a user: `PUT /users/{uid}`
* Delete a user: `DELETE /users/{uid}`

We now need to do the same for handling passports. Don't forget that a passport belongs to a user, so to retrieve a list of all passports for a given user, we would use `GET /users/{uid}/passports`.

When we want to retrieve an specific passport, we don't need to prefix the route with `/users/{uid}` anymore because we know exactly which passport we want to retrieve. So, instead of `GET /users/{uid}/passports/{pid}`, we can just use `GET /passports/{pid}`.

Once you have the API design sorted, it's just a matter of creating the code that gets called when a specific route is hit. We implement those with Handlers.

```
router.HandleFunc("/users", makeHandler(env, ListUsersHandler)).Methods("GET")
router.HandleFunc("/users/{uid:[0-9]+}", makeHandler(env, GetUserHandler)).Methods("GET")
router.HandleFunc("/users", makeHandler(env, CreateUserHandler)).Methods("POST")
router.HandleFunc("/users/{uid:[0-9]+}", makeHandler(env, UpdateUserHandler)).Methods("PUT")
router.HandleFunc("/users/{uid:[0-9]+}", makeHandler(env, DeleteUserHandler)).Methods("DELETE")

router.HandleFunc("/users/{uid}/passports", makeHandler(env, PassportsHandler)).Methods("GET")
router.HandleFunc("/passports/{pid:[0-9]+}", makeHandler(env, PassportsHandler)).Methods("GET")
router.HandleFunc("/users/{uid}/passports", makeHandler(env, PassportsHandler)).Methods("POST")
router.HandleFunc("/passports/{pid:[0-9]+}", makeHandler(env, PassportsHandler)).Methods("PUT")
router.HandleFunc("/passports/{pid:[0-9]+}", makeHandler(env, PassportsHandler)).Methods("DELETE")
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
  render.JSON(w, http.StatusOK, db.List())
}
```

Where the `render` variable is a global variable. We'd prefer to pass the Render variable to that function so will rewrite it to:

```
func ListUsersHandler(w http.ResponseWriter, req *http.Request, render Render) {
  render.JSON(w, http.StatusOK, db.List())
}
```

Perfect. Or, is it? What if we want to pass more variables to the handler function? Like Metrics? Or some environment variables? We'd continuously have to change ALL our handlers and the type signature will become quite long and hard to maintain. An alternative is to create a server or an environment struct and use that as a container for our variables we want to pass around the system.

In our `main.go` we'll add:

```
type appContext struct {
  metrics *stats.Stats
  render  *render.Render
  version string
	env     string
	port    string
	db      DataStorer
}
```

And in our `func main()` function we initialise that struct:

```
func main() {
  ctx := appContext{
    metrics: stats.New(),
    render:  render.New(),
    version: version,
		env:     env,
		port:    port,
		db:      db,
  }
  // ...
}
```

Our handler looks now like this:

```
func ListUsersHandler(w http.ResponseWriter, req *http.Request, ctx appContext) {
  // code that retrieves users from database
  ctx.render.JSON(w, http.StatusOK, responseObject
}
```

The only problem is that this handler's type signature is not `http.ResponseWriter, *http.Request` but `http.ResponseWriter, *http.Request, appContext` so Go's HandleFunc function will complain about this. That's why we are introducing a helper function `makeHandler` that takes our appContext struct and our handlers with the special type signature and converts it to `func(w http.ResponseWriter, r *http.Request)`:

```
func makeHandler(ctx appContext, fn func(http.ResponseWriter, *http.Request, appContext)) http.HandlerFunc {
  return func(w http.ResponseWriter, r *http.Request) {
    fn(w, r, ctx)
  }
}
```

### Special Route: Health check

When you look at the overview of the handlers in `routes.go`, you will notice a special route:

```
Route{"Healthcheck", "GET", "/healthcheck", HealthcheckHandler},
```

Monitoring tools like [Sensu](https://sensuapp.org/) can call: `GET /healthcheck`. The health check route can return a 200 OK when the service is up and running and will also say what the application name is and the version number.

```
func HealthcheckHandler(w http.ResponseWriter, req *http.Request, ctx appContext) {
	check := Healthcheck{
		AppName: "go-rest-api-template",
		Version: ctx.version,
	}
	ctx.render.JSON(w, http.StatusOK, check)
}
```

This health check is very simple. It just checks whether the service is up and running, which can be useful in a build and deployment pipelines where you can check whether your newly deployed API is running (as part of a smoke test). More advanced health checks will also check whether it can reach the database, message queue or anything else you'd like to check. Trust me, your DevOps colleagues will be very grateful for this. (Don't forget to change your HTTP status code to 200 if you want to report on the various components that your health check is checking.)

### Returning Data

Let's have a look at interacting with our data. Returning a list of users is quite easy, it's just showing the UserList:

```
func ListUsersHandler(w http.ResponseWriter, req *http.Request, ctx appContext) {
  ctx.render.JSON(w, http.StatusOK, db.List())
}
```

BTW, notice the `ctx.render.JSON`? That's part of `"github.com/unrolled/render"` and allows us to render JSON output when we send data back to the client.

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

Another example is the retrieval of a specific object:

```
func GetUserHandler(w http.ResponseWriter, req *http.Request, ctx appContext) {
  vars := mux.Vars(req)
  uid, _ := strconv.Atoi(vars["uid"])
  user, err := ctx.db.Get(uid)
  if err == nil {
    ctx.render.JSON(w, http.StatusOK, user)
  } else {
    ctx.render.JSON(w, http.StatusNotFound, err)
  }
}
```

This reads the uid variable from the route (`/users/{uid}`), converts the string to an integer and then looks up the user in our UserList by ID. If the user does not exit, we return a 404 and an error object. If the user exists, we return a 200 and a JSON object with the user.

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
curl -X GET http://localhost:3009/users | python -mjson.tool
```

That should result in the following result:

```
  % Total    % Received % Xferd  Average Speed   Time    Time     Time  Current
                                 Dload  Upload   Total   Spent    Left  Speed
100   216  100   216    0     0  14466      0 --:--:-- --:--:-- --:--:-- 15428
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

The `| python -mjson.tool` at the end is for pretty printing (formatting). It essentially tells to pipe the output of the curl command to the SJON formatting tool. If we only typed `curl -X GET http://localhost:3009/users` then we'd have something like this:

```
{"users":[{"id":0,"firstName":"John","lastName":"Doe","dateOfBirth":"1985-12-31T00:00:00Z","locationOfBirth":"London"},{"id":1,"firstName":"Jane","lastName":"Doe","dateOfBirth":"1992-01-01T00:00:00Z","locationOfBirth":"Milton Keynes"}]}
```

So not that easy to read as the earlier nicely formatted example.

Get a specific user:

```
curl -X GET http://localhost:3009/users/0 | python -mjson.tool
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

Adding a user:

```
TO DO
```

Deleting a user:

```
TO DO
```

Updating an existing user:

```
TO DO
```

### Testing the Handlers

TO DO

### Testing the Database

There are lots of opinions on testing, how much you should be testing, which layers of your applications, etc. When I'm working with micro services, I tend to focus on two types of tests to start with: testing the data access layer and testing the actual HTTP service.

In this example, we want to test the List, Add, Get, Update and Delete operations on our in-memory document database. The data access code is stored in the `database.go` file, so following Go convention we will create a new file called `database_test.go`.

In the `database_test.go` file, we have two sections:

First, we're going to create a common initialiser, this code sets up our database, inserts a couple of test records and then fires off the tests:

```
func TestMain(m *testing.M) {
  list := make(map[int]User)
  list[0] = User{0, "John", "Doe", "1985-12-31T00:00:00Z", "London"}
  list[1] = User{1, "Jane", "Doe", "1992-01-01T00:00:00Z", "Milton Keynes"}
  db = &Database{list, 1}
  retCode := m.Run()
  os.Exit(retCode)
}
```
Once this is ready, we can start writing tests. Let's have a look at the easiest one where we list the elements in our database:

```
func TestList(t *testing.T) {
  list := db.List()
  count := len(list["users"])
  assert.Equal(t, 2, count, "There should be 2 items in the list.")
}
```

This first calls the `db.List()` function, which returns a list of users. We then count the number of elements and last but not least we then check whether that count equals 2.

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
func TestDeleteSuccess(t *testing.T) {
  ok, err := db.Delete(1)
  assert.Equal(t, true, ok, "they should be equal")
  assert.Nil(t, err)
}

func TestDeleteFail(t *testing.T) {
  ok, err := db.Delete(10)
  assert.Equal(t, false, ok, "they should be equal")
  assert.NotNil(t, err)
}
```

The first test function `TestDeleteSuccess` tries to delete a known existing user, with Id 1. We're expecting that the error object is Nil. The second test function `TestDeleteFail` tries to look up a non-existing user with Id 10, and as expected, this should return an actual Error object.

How do we run the tests?

Simple:

```
go test
```

If you want it more verbose, then:

```
go test -v
```

Which will give you:

```
=== RUN TestList
--- PASS: TestList (0.00s)
=== RUN TestGetSuccess
--- PASS: TestGetSuccess (0.00s)
=== RUN TestGetFail
--- PASS: TestGetFail (0.00s)
=== RUN TestAdd
--- PASS: TestAdd (0.00s)
=== RUN TestUpdateSuccess
--- PASS: TestUpdateSuccess (0.00s)
=== RUN TestUpdateFail
--- PASS: TestUpdateFail (0.00s)
=== RUN TestDeleteSuccess
--- PASS: TestDeleteSuccess (0.00s)
=== RUN TestDeleteFail
--- PASS: TestDeleteFail (0.00s)
PASS
ok    github.com/leeprovoost/go-rest-api-template 0.008s
```

Do you want to get some more info on your code coverage? No worries, Go has you covered (no pun intended):

```
go test -cover
```

This will give you:

```
PASS
coverage: 34.9% of statements
ok    github.com/leeprovoost/go-rest-api-template 0.009s
```

## Starting the app on a production server

This is how you could run your app on a server:

First, you copy the binary and the `fixtures.json` + `VERSION` files into a directory, e.g. `/opt/go-rest-api-template`.

Then start the app as a service. Store the app's PID in a text file so we can kill it later.

```
#!/bin/bash
export ENV=DEV
export PORT=8080
export VERSION=/opt/go-rest-api-template/VERSION
export FIXTURES=/opt/go-rest-api-template/fixtures.json
sudo nohup /opt/go-rest-api-template/go-rest-api-template >> /var/log/go-rest-api-template.log 2>&1&
echo $! > /opt/go-rest-api-template/go-rest-api-template-pid.txt
```

When you want to kill your app later during a redeployment or a server shutdown, then you can kill the app by looking up the previously stored PID:

```
#!/bin/bash
if [ -f /opt/go-rest-api-template/go-rest-api-template-pid.txt ]; then
  kill -9 `cat /opt/go-rest-api-template/go-rest-api-template-pid.txt`
  rm -f /opt/go-rest-api-template/go-rest-api-template-pid.txt
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
