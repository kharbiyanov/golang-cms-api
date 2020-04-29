package main

import "cms-api/models"

var plugin = models.Plugin{
	ID:     "files",
	Name:   "Files",
	Enable: true,
}

func Init() models.Plugin {
	InitSchema(&plugin)
	return plugin
}
