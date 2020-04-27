package main

import "cms-api/models"

var plugin = models.Plugin{
	ID:     "media",
	Name:   "Media",
	Enable: true,
}

func Init() models.Plugin {
	InitSchema(&plugin)
	return plugin
}
