package main

import "cms-api/models"

var plugin = models.Plugin{
	ID:     "auth",
	Name:   "Authentication",
	Enable: true,
}

func Init() models.Plugin {
	InitSchema(&plugin)
	return plugin
}
