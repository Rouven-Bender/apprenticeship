# About
This project is a parallel rewrite of my Apprenticeship final project from the Java Fullstack with an inhouse developed Java Framework 
(not written by me) to open source technologies as that is more interesting
# Usage
This program uses subcommands for the features:
- "serve" - runs the server
- "export" - runs the export
- "addUser" - this will take a -u {username} and ask for a password and add those credentials
- "deactivateExpired" - this will set all licenses with an expiry date past the current time as inactiv
# Dependencies
- htmx @ 2.0.2 in the folder static as "htmx-2.0.2.min.js" (not vendored into the repo)
- github.com/mattn/go-sqlite3
    - needs gcc
- back button icon as svg as ./static/arrow-left.svg
- an "JWT_SECRET" file in the project folder with an string that is at least 32 chars long
# The Project
The License Server(of my apprenticeship company) has a feature to register sublicenses (with limited number of seats) of a bigger License(maybe limited in total number of seats).
But to configure that sublicense to be accepted by the License Server it needs to be written to that servers config
The servers config accepts to import json as part of the config so I need to export a JSON file containing all sublicenses that need to be activ
# JSON Schema
(not how its actually on that license server but how I'm gonna export it in this rewrite)
```
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
```
# Features
- removing expired licenses from file automaticly
- Search for Customer or License Key
- Adjust Number of Seats
- Create and Deactivate licenses manually
# The Gameplan
- Write a Golang application that does these things:
    - Run the backend for the htmx frontend
    - Run an export of the database and write it to disk (used for a cronjob)
    - A command for the sysadmins to add a User for the webapp
    - a command to check for expired licenses and deactivate them
- Use SQLite as the database because this is would be internal software that has barely any concurrent users
