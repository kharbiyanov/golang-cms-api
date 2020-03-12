package main

import "cms-api/models"

var plugin = models.Plugin{
	ID:     "translate",
	Name:   "Translate",
	Enable: true,
}

func Init() models.Plugin {
	InitSchema(&plugin)
	return plugin
}
