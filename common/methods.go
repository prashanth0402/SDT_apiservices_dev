package common

import (
	"errors"
	"log"

	"github.com/BurntSushi/toml"
)

// --------------------------------------------------------------------
// function reads the constants from the config.toml file
// --------------------------------------------------------------------
func ReadTomlConfig(filename string) interface{} {
	var f interface{}
	if _, err := toml.DecodeFile(filename, &f); err != nil {
		log.Println(err)
	}
	return f
}

//----------------------------------------------------------------
// Creating custom error
// ----------------------------------------------------------------

func CustomError(pErrorMsg string) error {
	err := errors.New(pErrorMsg)
	return err
}
