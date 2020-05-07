package main

import "cms-api/models"

var plugin = models.Plugin{
	ID:     "users",
	Name:   "Users",
	Enable: true,
}

func Init() models.Plugin {
	InitSchema(&plugin)
	return plugin
}
