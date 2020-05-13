package main

import "cms-api/models"

var plugin = models.Plugin{
	ID:     "roles",
	Name:   "Roles",
	Enable: true,
}

func Init() models.Plugin {
	InitSchema(&plugin)
	return plugin
}
