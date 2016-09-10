# BC-FEATURE-REQUESTS

[![Build Status](https://travis-ci.org/kujtimiihoxha/bc-feature-requests.svg?branch=master)](https://travis-ci.org/kujtimiihoxha/bc-feature-requests)
[![Build status](https://ci.appveyor.com/api/projects/status/kaetcdk7u9j9tiji?svg=true)](https://ci.appveyor.com/project/kujtimiihoxha/bc-feature-requests)

This is the API for [bc-feature-requests-client](https://github.com/kujtimiihoxha/bc-feature-requests-client).

## Technologies
- Based on GOLANG a language designed for the CLOUD.
- Beego framework for web apps.
- RethinkDB a modern NoSQL database.
- Goconvey a very powerful coverage report testing framework.

## Platform independent
- This rest API can be run in Linux, Windows, Mac.

## Install
 To build this project you need to have GO installed on your machine.
 Setup the Go environment variables and ``go get`` the project. 

```
go get github.com/kujtimiihoxha/bc-feature-requests
```
After you get the project cd to the project folder and install dependencies.

```
 go get -v -t ./...
```

Now you can simply run.
```
go run main.go
```

If you want to run coverage report then install ``goconvey``.
```
 go get github.com/smartystreets/goconvey
```

Then run:
```
 goconvey -port 
```
Navigate to ```http://localhost:8080``` and you will get the coverage report in a very user friendly web app.

## Features
- Very fast, golang is a very fast programming language and it is designed with web apps in mind.
- Completely decoupled from the frontend (besides the user verification link)
- JWT for user authentication.
- Test suit (not complete)
- CI for Windows,Linux,Mac
- Highly scalable.
- Easily configurable using one configuration file. 
- Very easy deployment, only one executable contains the whole api.

To deploy the app you need to build the app.
```
go build
```
Than copy the config folder, logs folder(if you want to keep the default logs output) and the created binary to the server.
```
|-conf/
|   |- app.conf
|-logs/
|   |- bc.log
|created-binary.
```

## Configuration file
**app.conf**
```
# app name
appname = "bc-feature-requests"
# the server url (ex. api.britecore.com)
server-url = "bc.kujtimhoxha.com"
# port to listen to
httpport = 8084
# enviorment mode
runmode = "prod"
#bee configurations
autorender = false
copyrequestbody = true
EnableDocs = false
graceful = true

#server name, this will bee the value of the Server header.
server = "bc.kujtimhoxha.com"

# Log configurations
# enabled : if true server will log to the specific file
# path : path to the log file
[log]
enabled = true
path = "./logs/bc.log"

# Database configurations
# host: database server url with port
# databse: the database name
[db]
host = "kujtimhoxha.com:28015"
database = "bc_feature_requests"

#Cors configurations
[cors]
allow-origins = *
allow-methods = GET;POST;PUT;DELETE;PATCH
allow-credentials = true
expose-headers = Content-Length

#JWT configurations
[jwt]
key = "snC17VHM5DcpPuWpJcCl6f78j3A9AB8L"
hours = 24

#JWT configurations
[mail]
gmail_account =
gmail_account_password =
mail_host = smtp.gmail.com
mail_host_port = 587
```

Most of the important configurations can be set on this configuration file so there is no need to recompile the app if used in different environments