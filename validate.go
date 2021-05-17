package main

import (
	"github.com/ganluo960214/validations"
	"github.com/go-playground/validator/v10"
	"log"
)

var (
	validate = validator.New()
)

func init() {
	if err := validations.Register_CheckFileExists_Validation(validate); err != nil {
		log.Fatalln(err)
	}
}
