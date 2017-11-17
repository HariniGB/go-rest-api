# go-rest-api-mgo
## Jenkins Build Status   ![Success](https://raw.github.com/gigablah/jenkins-status/master/images/success.png)
-------------------------


# [Go (Golang)](https://golang.org)


## Quick Start

To install: Download the .pkg [here](https://golang.org/dl/)

### To set up environment
```bash
$ mkdir golang
$ export GOPATH=[full_path_to_dir]/golang
```

### Recommended folder structure
**bin/** and **src/**<br>
**bin/** will containe the compiled files<br>
**src/** should have one package file per folder

Go files has file ending **.go**<br>
To compile the files:<br>
```bash
$ go install folder_path*
```
*folder_path is the path from src/, do not include src/ in the folder path


To run the compiled file:
```bash
$ bin/folder_path
```
or
```bash
$ go run (file-name).go
```

## What is Golang

Go, also commonly referred to as golang, is a programming language developed at Google in 2007 by Robert Griesemer, Rob Pike, and Ken Thompson. It is a statically typed language with syntax loosely derived from that of C, adding garbage collection, type safety, some structural typing capabilities, additional built-in types such as variable-length arrays & key-value maps, and a large standard library.
Android support was added in version 1.4, which has since been ported to also run on iOS. - [Wikipedia](https://en.wikipedia.org/wiki/Go_(programming_language))<br>

## Dependencies:
1. Install mongo db 
#### [Install Mongo in Mac details](https://treehouse.github.io/installation-guides/mac/mongo-mac.html)
2. mongo db path/bin/mongod to start mongodb server
3. mongo db path/bin/mongo to start mongodb shell
4. To view the data inside database collections 
```bash
$ mongo show
$ show dbs
$ use go_rest
$ show colections
$ db.users.find()
```
####  For futher reference: [Install Mongo in Mac details](https://docs.mongodb.com/manual/reference/mongo-shell/)

## Resources
### [Go Doc](https://godoc.org/github.com/)
### [Interactive tutorial](https://tour.golang.org/welcome/1)
### [Repl.it like editor](http://play.golang.org/)
