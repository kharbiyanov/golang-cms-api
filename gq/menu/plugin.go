package main

import "cms-api/models"

var plugin = models.Plugin{
	ID:     "menu",
	Name:   "Menu",
	Enable: true,
}

func Init() models.Plugin {
	InitSchema(&plugin)
	return plugin
}
