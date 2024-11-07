package main

type fileAlias struct {
	Htmx             string
	StyleIndex       string
	FilterIndexTable string
	StyleSublicense  string
}

var ALIAS fileAlias = fileAlias{
	Htmx:             "cdn/htmx-2.0.2.min.js",
	StyleIndex:       "cdn/style-index-0.1.css",
	FilterIndexTable: "cdn/filter-0.1.js",
	StyleSublicense:  "cdn/style-sublicense.css",
}
