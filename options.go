package main

type Options struct  {
	Port int `short:"p" long:"port" desctiption:"Port number"`
	ConsumerKey string `short:"k" long:"key" description:"Twitter consumer key"`
	ConsumerSecret string `short:"s" long:"secret" description:"Twitter consumer secret"`
	Database string `short:"d" long:"db" description:"Kind of database supported mysql, postgres, and sqlite"`
	DBOptions string `short:"o" long:"dboptions" description:"Database options"`
}
