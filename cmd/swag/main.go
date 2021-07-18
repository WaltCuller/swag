package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/spf13/viper"
	"github.com/swaggo/swag"
	"github.com/swaggo/swag/gen"
	"github.com/urfave/cli/v2"
)

const (
	searchDirFlag        = "dir"
	excludeFlag          = "exclude"
	generalInfoFlag      = "generalInfo"
	propertyStrategyFlag = "propertyStrategy"
	outputFlag           = "output"
	parseVendorFlag      = "parseVendor"
	parseDependencyFlag  = "parseDependency"
	markdownFilesFlag    = "markdownFiles"
	codeExampleFilesFlag = "codeExampleFiles"
	parseInternalFlag    = "parseInternal"
	generatedTimeFlag    = "generatedTime"
	parseDepthFlag       = "parseDepth"
	configFlag           = "config"
)

var initFlags = []cli.Flag{
	&cli.StringFlag{
		Name:    generalInfoFlag,
		Aliases: []string{"g"},
		Value:   "main.go",
		Usage:   "Go file path in which 'swagger general API Info' is written",
	},
	&cli.StringFlag{
		Name:    searchDirFlag,
		Aliases: []string{"d"},
		Value:   "./",
		Usage:   "Directory you want to parse",
	},
	&cli.StringFlag{
		Name:  excludeFlag,
		Usage: "Exclude directories and files when searching, comma separated",
	},
	&cli.StringFlag{
		Name:    propertyStrategyFlag,
		Aliases: []string{"p"},
		Value:   "camelcase",
		Usage:   "Property Naming Strategy like snakecase,camelcase,pascalcase",
	},
	&cli.StringFlag{
		Name:    outputFlag,
		Aliases: []string{"o"},
		Value:   "./docs",
		Usage:   "Output directory for all the generated files(swagger.json, swagger.yaml and doc.go)",
	},
	&cli.BoolFlag{
		Name:  parseVendorFlag,
		Usage: "Parse go files in 'vendor' folder, disabled by default",
	},
	&cli.BoolFlag{
		Name:  parseDependencyFlag,
		Usage: "Parse go files in outside dependency folder, disabled by default",
	},
	&cli.StringFlag{
		Name:    markdownFilesFlag,
		Aliases: []string{"md"},
		Value:   "",
		Usage:   "Parse folder containing markdown files to use as description, disabled by default",
	},
	&cli.StringFlag{
		Name:    codeExampleFilesFlag,
		Aliases: []string{"cef"},
		Value:   "",
		Usage:   "Parse folder containing code example files to use for the x-codeSamples extension, disabled by default",
	},
	&cli.BoolFlag{
		Name:  parseInternalFlag,
		Usage: "Parse go files in internal packages, disabled by default",
	},
	&cli.BoolFlag{
		Name:  generatedTimeFlag,
		Usage: "Generate timestamp at the top of docs.go, disabled by default",
	},
	&cli.IntFlag{
		Name:  parseDepthFlag,
		Value: 100,
		Usage: "Dependency parse depth",
	},
}

var uploadFlags = []cli.Flag{
	&cli.StringFlag{
		Name:    configFlag,
		Aliases: []string{"c"},
		Value:   "./config/swag.toml",
		Usage:   "swag upload config",
	},
	&cli.StringFlag{
		Name:    outputFlag,
		Aliases: []string{"o"},
		Value:   "./docs",
		Usage:   "Output directory for all the generated files(swagger.json, swagger.yaml and doc.go)",
	},
}

func initAction(c *cli.Context) error {
	strategy := c.String(propertyStrategyFlag)

	switch strategy {
	case swag.CamelCase, swag.SnakeCase, swag.PascalCase:
	default:
		return fmt.Errorf("not supported %s propertyStrategy", strategy)
	}

	return gen.New().Build(&gen.Config{
		SearchDir:           c.String(searchDirFlag),
		Excludes:            c.String(excludeFlag),
		MainAPIFile:         c.String(generalInfoFlag),
		PropNamingStrategy:  strategy,
		OutputDir:           c.String(outputFlag),
		ParseVendor:         c.Bool(parseVendorFlag),
		ParseDependency:     c.Bool(parseDependencyFlag),
		MarkdownFilesDir:    c.String(markdownFilesFlag),
		ParseInternal:       c.Bool(parseInternalFlag),
		GeneratedTime:       c.Bool(generatedTimeFlag),
		CodeExampleFilesDir: c.String(codeExampleFilesFlag),
		ParseDepth:          c.Int(parseDepthFlag),
	})
}

func main() {
	app := cli.NewApp()
	app.Version = swag.Version
	app.Usage = "Automatically generate RESTful API documentation with Swagger 2.0 for Go."
	app.Commands = []*cli.Command{
		{
			Name:    "init",
			Aliases: []string{"i"},
			Usage:   "Create docs.go",
			Action:  initAction,
			Flags:   initFlags,
		},
		{
			Name:    "upload",
			Aliases: []string{"u"},
			Usage:   "upload swagger.json to yapi",
			Action:  Post,
			Flags:   uploadFlags,
		},
	}
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func initConfig(c *cli.Context) {
	if c.String(configFlag) != "" {
		// Use config file from the flag.
		viper.SetConfigFile(c.String(configFlag))
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		if err != nil {
			panic(err)
		}

		// Search config in home directory with name ".cobra" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".cobra")
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}
	fmt.Println("Using config file:", viper.ConfigFileUsed())
}

type config struct {
	Type   string `json:"type"`
	Token  string `json:"token"`
	File   string `json:"file"`
	Merge  string `json:"merge"`
	Server string `json:"server"`
}

// Post post swagger.json
func Post(c *cli.Context) error {
	initConfig(c)
	config := config{
		Type: viper.Get("swagger.type").(string),
	}

	path := viper.Get("swagger.server").(string)
	token := viper.Get("swagger.token").(string)
	filePath := viper.Get("swagger.file_path").(string)

	fmt.Println("type: ", config.Type)
	fmt.Println("path: ", path)
	fmt.Println("token: ", token)
	fmt.Println("filePath: ", filePath)

	bytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		fmt.Printf("read file[%s] failed err=%v", filePath, err)
		return err
	}
	actualToken := config.Token
	if config.Token == "" {
		actualToken = token
	}
	values := url.Values{
		"type":  {config.Type},
		"json":  {string(bytes)},
		"merge": {config.Merge},
		"token": {actualToken},
	}
	fmt.Println(string(bytes))
	client := &http.Client{}
	uri := path + "/api/open/import_data"
	req, err := http.NewRequest("POST", uri, strings.NewReader(values.Encode()))
	if err != nil {
		fmt.Printf("http.NewRequest failed err=%v", err)
		return err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Close = true
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("client.Do failed err=%v", err)
		return err
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			panic(err)
		}
	}(resp.Body)

	result, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("read responde failed err=%v", err)
		return err
	}
	fmt.Println(string(result))
	return nil
}
