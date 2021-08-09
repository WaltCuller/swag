# swag-yapi

forked from [swag](https://github.com/swaggo/swag) and developed for yapi feature.

WARNING: This is still in beta stage, use with care!

## Features

```go
package main

type People struct {
	Name string `json:"name" mock:"@string(\"lower\",5)"`
	Age int `json:"age" mock:"@integer(0,100)"`
	Birthday string `json:"birthday" mock:"@datetime(\"yyyy-MM-dd A HH:mm:ss\")"`
}
```
- [x] upload swagger.json to yapi server.

- [x] adapt mock.js.


## Usage

### Installation

1. To install it, either download a version from the [release](https://github.com/WaltCuller/swag-yapi/releases/tag/v1.7.0-beta) page, use ```go install/go get```.
```bash
go get -u https://github.com/WaltCuller/swag-yapi
```

2. Add comments to your API source code. See [Declarative Comments Format](https://github.com/swaggo/swag#declarative-comments-format).

3. For response struct, swag-yapi use the customize tag: ```mock```. See [mock.js](http://mockjs.com/examples.html)

4. Add toml config to define the yapi server address and configuration.
```toml
[swagger]
type = "swagger" # yapi upload type
token = "xxxxx" # yapi token
merge = "normal" # merge type
server = "xxxxxx" # yapi server address
file_path = "docs/swagger/swagger.json" # swagger.json file path
```

6. After added comments and tags, run ```swag-yapi init``` like ```swag init``` in the project's root folder which contains the ```main.go``` file. This will parse your comments and generate the require files(```docs```folder and ```docs/docs.go```)

7. use ```swag-yapi upload``` with option ```-c``` to get the config and upload swagger.json to your yapi server.

### Client help
```shell
$ swag-yapi -h
NAME:
   swag-yapi - Automatically generate RESTful API documentation with Swagger 2.0 for Go.

USAGE:
   swag-yapi [global options] command [command options] [arguments...]

VERSION:
   v1.7.0_0.0.1

COMMANDS:
   init, i    Create docs.go
   upload, u  upload swagger.json to yapi
   help, h    Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h     show help (default: false)
   --version, -v  print the version (default: false)
```
