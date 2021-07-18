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
-[x] upload swagger.json to yapi server.
-[x] adapt mock.js.

## Usage

### Installation

To install it, either download a version from the [release](https://github.com/WaltCuller/swag-yapi/releases/tag/v1.7.0-beta) page, use ```go install``` or pull the source code from ```master```.


