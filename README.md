# About
This project is a parallel rewrite of my Apprenticeship final project
From the Java Fullstack with an inhouse developed Java Framework (not written by me) to open source technologies
As that is more interesting, I think
# Dependencies
- htmx @ 2.0.2 in the folder static as "htmx.min.js" (not vendor'd into the repo)
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
- Have a counter that shows total reserved seats out of total seats in the OEM License (if the oem license isn't infinite) then just a running counter
- Check before exporting if the file was manually changed and make a warning
# The Gameplan
- Write a Golang application that has 2 subcommands:
    - Run the backend for the htmx frontend
    - Run an export of the database and write it to disk (used for the cronjob)
- have a bash script (cronjob) that runs a query on the database to set sublicenses as deactivated when they expire and after that query finishes re-export the database state to the config
(maybe have the query also as part of the 2. subcommand and just have cron run that secound sub command)
- Use SQLite as the database because this is would be internal software that has barely any users
