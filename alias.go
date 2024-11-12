package main

type fileAlias struct {
	Htmx             string
	StyleIndex       string
	FilterIndexTable string
	StyleSublicense  string
	ErrorMessageHook string
}

var ALIAS fileAlias = fileAlias{
	//Libs
	Htmx: "cdn/htmx-2.0.2.min.js",
	//CSS
	StyleIndex:      "cdn/style-index-0.1.css",
	StyleSublicense: "cdn/style-sublicense.css",
	//JS
	FilterIndexTable: "cdn/filter-0.1.js",
	ErrorMessageHook: "cdn/errormessage.js",
}
