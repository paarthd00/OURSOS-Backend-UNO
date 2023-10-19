# Backend for Oursos

## Introduction

This guide outlines the steps to set up and run the backend for Oursos, an application that presumably relies on a PostgreSQL database. It provides instructions for creating a `.env` file, downloading Go modules, building the application, and organizing the project's folder structure.

## Installing Go onto your device

To install Go language onto your local device follow these steps:

1. Please Visit this link to install go on your system [Go Installation](https://go.dev/doc/install)

2. Click on the Download button and download the file for your relevant system

3. On Windows open up the powershell command line and type the command `go version`

4. Confirm that the command line returns the version of go installed

## Environment Setup

To configure your local respository to access the database and API keys follow these steps:

1. Create a `.env` file in your project's root directory to store the local environment variables

2. On the OurSoS server download the file `envVariables.txt` and paste the text into the `.env` file

## Building and Running the Application

To build and run the Oursos backend, follow these steps:

1. Open your terminal and navigate to the project's root directory.

2. Download Go modules using the following command:

    ```bash
    go mod download
    ```

3. Build the application with the following command:

    ```bash
    go build -o out
    ```

4. Run the application:

    ```bash
    ./out
    ```

## Folder Structure

### Modules

In Go the entire project is referred to as a module. If you navigate to the `go.mod` file you will see:
```
module oursos.com/packages
```
The module declaration gives our project a name and informs Go that any folders, files or APIs being called
are associated to this 'module' or project.

### Packages

In Go packages are folders and the files contained within the folder. The packages encapsulate the logic inside them
and create a seperation between different pieces of logic. In the case of this project these are our packages:
```
Root
├── alerts
├── api
├── db
├── users
└── util
```
If you open one of these packages and look at the files within. They will always start with an definition that references the package that they are a part of. For example if you open the API package directory, all files within start with:

```
package api
```

This informs Go that the file within the API package folder, contains logic that is relevant to the API package.

Each package has its own logic and the files within it are directly related to achieving whatever functionality the package
was made for. That being said functions and logic within one package can be referenced in another without importing. Below are some guidelines for organizing your project's folder structure:

- **Package Naming:** Conventionally, all files within a folder should have the same package name, which is typically the folder name.

- **Exported Functions:** Functions that need to be accessible outside the package should be capitalized. Functions that start with a lowercase letter are treated as private and won't be accessible outside the package, except for the `main` function.

## Servers + Routing

`Echo` is the web framework that Go uses in order to host the backend server. In many ways this is very similar to an Express.js server. The setup of the Go `Echo` server can be found in `out.go` as listed below:

```
func main() {
	// go db.SeedDatabase()
	err := godotenv.Load()

    if err != nil {
        fmt.println("Failed to to load the environment variables")
    }
	
    e := echo.New()
	e.GET("/", homeHandler)
	e.DELETE("/deleteuser/:id", users.DeleteUser)
	e.PUT("/updateuser/:id", users.UpdateUser)

	e.Logger.Fatal(e.Start("0.0.0.0:" + os.Getenv("PORT")))
}
```
* `err := godotenv.Load()` - informs app there is a `.env` file
    - `err` - err variable will hold the error if `.env` cannot be located
    - The only time `err` holds a value is if there is an issue in retrieving the `.env` file, otherwise
    the `godotenv.Load()` function performs and resolves itself without the variable definition updating
    - The `err` variable exists so the `error trace` can be accessed

* `e := echo.New()` - creates a variable `e` that is defined as an instance of the `Echo` server

* `e.GET('/route', getRouteHandler)` - runs a GET process at the defined route using the defined handler

* `e.DELETE('/route/:param', deleteRouteHandler)` - runs a DELETE process at the defined route through the defined handler
    - `:param` - allows the route to inherit a URL parameter for its backend logic

* `e.PUT('/route', updateRouteHandler)` - runs a UPDATE process at the defined route through the defined handler

* `e.logger.Fatal(e.Start("0.0.0.0:" + os.getEnv("PORT)))` - starts the configured `Echo` server
    - `e.logger.Fatal` - inbuilt Go feature that will log errors in server startup
    - `e.start('IP + PORT')` - function starts the server on the given IP + `PORT`
    - `0.0.0.0` - informs the server it can be run on any `PORT`

When a server is configured to listen on IP address "0.0.0.0," it means that it will accept incoming connections from any available network interface on the host. In other words, it binds to all available network interfaces.

* `os.Getenv("PORT")` - retrieves env variables from `.env` file using `godotenv.Load()`
    - `PORT` - `PORT` is an env variable that tells app what `PORT` to host the app on

## Understanding Functions

Within Go every function has two general possibilities for an output. Either the function returns something or the function
returns an error.

### Function Structure

To understand the function structure let's go through an example function as listed below:

```
func ExampleFunction(c echo.Context) error {

    err,output := // Function logic output

	if err != nil{
        return new Error(err)
    }

	return c.JSON(http.StatusOK, output)

}
```

* `func ExampleFunction(c echo.Context) error {`
    - `c`` - variable name, representative of the word 'context'
    - `echo.Context` - variable data type, Context is a variable provided by `Echo` containing:
        - `Response Data`
        - `Resquest Data`
        - `Params`
        - Handles the `Req, Res` functionality in `HTTP` communication and networking
    - `error` - data type for the return, we have specified `error` but if the return is `void` it is left empty
        - In this case the function either returns `error` or `void` so we only list `error` 

* `err, output := // Function logic output`
    - `err, output` - two potential variables that the output could be assigned to
    - `:=` - variable assignment syntax in Go
    - Function returns some kind of output, if the output is successful it is assigned to `output`, if
    it is erroneous it is assigned to `err`

* `if err != nil { return new Error(err)}`
    - `if != nil` - if check that checks wheter `err` variable holds a definition
    - `return new Error(err)` - creates an instance of an `Error` class that holds your error information
    - Function structure built so if there is an `err` then it is returned

* `return c.JSON(http.StatusOK, output)`
    - `c.JSON` - data type of the output being returned via route response (JSON)
    - `http.StatusOK` - return `OK` status code `200`
    - `output` - JSON data to return that is defined in the `err, output` variable definition
    - Functions structure built so if there is an `output` then it is returned.

As can be seen, all functions in Go require a structure that takes into account errors. In most cases you must have error as a potential output strictly defined in the function output type. We can see this in the fact that we never defined the data type for the JSON output in the above function, but it still works because we've dealt with the possibility of an error in strictly typed syntax.

## Managing Externally Integrated Packages in Go (Similar to package.json)

### Installing an External Package

If you want to add an external dependency to the project, using the command line input the following command:
```
go get "packageName"
```

`go get "github.com/joho/godotenv"` - this command would import the `godotenv` library into the project

Functionally this is the same as `npm install`

### Creating a Dependency Reference

Tracking dependencies is automated in Go. When you run the above command the `go.mod` file will update to reflect the newly added dependency. In this section we will briefly discuss the structure of the `go.mod` file:

```
module oursos.com/packages

go 1.20

require (
	github.com/go-sql-driver/mysql v1.7.1
)

require (
	github.com/go-resty/resty/v2 v2.9.1 // indirect
)
```
* `module oursos.com/packages` - references the specific project that all these dependencies are for

* `go 1.20` - version of Go being used

* `require()` - two functions that hold the direct and indirect dependencies the projhect needs to run
    - **Direct Dependency** - packages that your Go module is directly importing and uses in the project
    - **Indirect Dependency** - packages that your direct dependencies need in order to work

### Cleaning Up Dependencies

If you end up importing a number of dependencies that you don't actually use, you can use the following command to clean up `go.mod`:

```
go mod tidy
```

## Managing Internally Integrated Packages in Go (Similar to Exporting Functions)

Internally refernced packages have their import routes defined in the beginning of the file importing them. For example if we wanted to import the `util` package inside `api/earthquake.go` we would need the following at the top of the `api/earthquake.go` file:

```
import (
    "oursos.com/packages/util"
)
```

We would need a similar import on the top of any file that we were importing the `util` functionality into. Functions inside the imported package can be accessed through dot notation. An example of this is:

```
util.CheckErrors(err)
```

## Conclusion

By the end of this guide you should have a basic understanding of the backend structure that Go uses, and be able to setup the environment and basic architecture for a project.

To gain more understanding of the syntax and logic patterns that Go uses [reference their documentation](https://go.dev/doc/).
