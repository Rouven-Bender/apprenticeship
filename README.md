# About
This project is a parallel rewrite of my Apprenticeship final project
From the Java Fullstack with an inhouse developed Java Framework (not written by me) to open source technologies
As that is more interesting, I think
# Usage
This program uses subcommands for the features:
    - "serve" - runs the server
    - "export" - runs the export
# Dependencies
- htmx @ 2.0.2 in the folder static as "htmx-2.0.2.min.js" (not vendor'd into the repo)
- github.com/mattn/go-sqlite3
    - needs gcc
- back button icon as svg as ./static/arrow-left.svg
# The Project
The License Server has a feature to register sublicenses (with limited number of seats) of a bigger License(maybe limited in total number of seats).
But to configure that sublicense to be accepted by the big central License Server it needs to be written to that servers config
The servers config accepts to import json as part of the config so I need to export a JSON file containing all sublicenses that need to be activ
# JSON (not how its actually on that server but how I'm gonna export it here)
´´´json
{
    "sublicense": [
        {
            "name": "Customer 1",
            "number_of_seats": 5,
            "license_key": "KEY-FOR-SOFT-WARE"
        },
        {
            "name": "Customer 1",
            "number_of_seats": 5,
            "license_key": "KEY-ROF-SOFT-WARE"
        }
    ]
}
´´´
# Features
- Expiry date
- removing expired licenses from file automaticly
- Search for Customer or License Key
- Adjust Number of Seats
- Create and Deactivate licenses manually
# The Gameplan
- Write a Golang server that does these three things:
    - Run the backend for the htmx frontend
    - Run an export of the database and write it to disk (used for the cronjob)
- Use SQLite as the database because this is would be internal software that has barely any concurrent users
# Things to add that are part of the feature set but haven't done yet
- Login and Auth - plan: wrap the endpoints that need JWT Key in a middleware that checks for valid key, have Login give the browser an jwt key into the cookies that has an expiry date
