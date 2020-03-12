package main

import "cms-api/models"

var plugin = models.Plugin{
	ID:     "posts",
	Name:   "Posts",
	Enable: true,
}

func Init() models.Plugin {
	InitSchema(&plugin)
	return plugin
}
