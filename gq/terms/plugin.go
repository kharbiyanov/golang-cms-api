package main

import "cms-api/models"

var plugin = models.Plugin{
	ID:     "terms",
	Name:   "Terms",
	Enable: true,
}

func Init() models.Plugin {
	InitSchema(&plugin)
	return plugin
}
