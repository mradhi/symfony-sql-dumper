package main

import (
	"database/sql"
	"fmt"
	"github.com/JamesStewy/go-mysqldump"
	_ "github.com/go-sql-driver/mysql"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
)

type ConfigParamsYAML struct {
	Parameters struct {
		DatabaseHost     string `yaml:"database_host"`
		DatabasePort     string `yaml:"database_port"`
		DatabaseName     string `yaml:"database_name"`
		DatabaseUser     string `yaml:"database_user"`
		DatabasePassword string `yaml:"database_password"`
	}
}

func main() {
	// Fetch args
	// parameters.yml file from symfony project
	yamlFile := os.Args[1]
	// Directory contains the backups
	dumpDir := os.Args[2]

	yamlParams := ConfigParamsYAML{}
	// Read YAML configuration file to get db connection.
	b, err := ioutil.ReadFile(yamlFile)
	if err != nil {
		fmt.Print(err)
	}

	// Parse the YAML file into the ConfigParamsYAML object.
	err = yaml.Unmarshal(b, &yamlParams)

	if err != nil {
		panic(err)
	}

	parameters := yamlParams.Parameters

	// Check if the port is empty
	if len(parameters.DatabasePort) == 0 {
		parameters.DatabasePort = "3306"
	}

	// Accepts time layout string and add .sql at the end of file
	dumpFilenameFormat := fmt.Sprintf("%s-20060102T150405", parameters.DatabaseName)

	conStr := fmt.Sprintf("%s:%s@%s:%s/%s",
		parameters.DatabaseUser,
		parameters.DatabasePassword,
		parameters.DatabaseHost,
		parameters.DatabasePort,
		parameters.DatabaseName)

	fmt.Printf("Opening connection for: '%s' ...\n", conStr)

	db, err := sql.Open("mysql", conStr)

	if err != nil {
		panic(err)
	}

	// Create the directory
	if _, err := os.Stat(dumpDir); os.IsNotExist(err) {
		_ = os.Mkdir(dumpDir, os.ModePerm)
	}

	// Register database with mysqldump command
	dumper, err := mysqldump.Register(db, dumpDir, dumpFilenameFormat)
	if err != nil {
		panic(err)
	}

	// Dump database to file
	resultFilename, err := dumper.Dump()
	if err != nil {
		panic(err)
	}

	fmt.Printf("File is saved to '%s'.\n", resultFilename)

	// Close dumper and connected database
	_ = dumper.Close()
}