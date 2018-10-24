# IGCInfoViewer

# About
Online service that allows users to browse information about IGC files. IGC is an international file format for soaring track files that are used by paragliders and gliders.

# Can be reached at
https://igcinfoviewer.herokuapp.com

# GET igcinfo/api
What: meta information about the API
Response type: application/json
Response code: 200
Body template
`
{
  "uptime": <uptime>
  "info": "Service for IGC tracks."
  "version": "v1"
}
`
where: <uptime> is the current uptime of the service formatted according to Duration format as specified by ISO 8601. 

# POST igcinfo/api/igc
What: track registration
Response type: application/json
Response code: 200 if everything is OK, appropriate error code otherwise, eg. when provided body content, is malformed or URL does not point to a proper IGC file, etc. Handle all errors gracefully. 
Request body template
`
{
  "url": "<url>"
}`
Response body template
`
{
  "id": "<id>"
}
`
where: <url> represents a normal URL, that would work in a browser, eg: http://skypolaris.org/wp-content/uploads/IGS%20Files/Madrid%20to%20Jerez.igc

# GET igcinfo/api/igc
What: returns the array of all tracks ids
Response type: application/json
Response code: 200 if everything is OK, appropriate error code otherwise. 
Response: the array of IDs, or an empty array if no tracks have been stored yet.
`
[<id1>, <id2>, ...]
`
# GET igcinfo/api/igc/<id>
What: returns the meta information about a given track with the provided <id>, or NOT FOUND response code with an empty body.
Response type: application/json
Response code: 200 if everything is OK, appropriate error code otherwise. 
Response: 
  `
{
"H_date": <date from File Header, H-record>,
"pilot": <pilot>,
"glider": <glider>,
"glider_id": <glider_id>,
"track_length": <calculated total track length>
}
`
# GET igcinfo/api/igc/<id>/<field>
What: returns the single detailed meta information about a given track with the provided <id>, or NOT FOUND response code with an empty body. The response should always be a string, with the exception of the calculated track length, that should be a number.
Response type: text/plain
Response code: 200 if everything is OK, appropriate error code otherwise. 
Response
  `
<pilot> for pilot
<glider> for glider
<glider_id> for glider_id
<calculated total track length> for track_length
<H_date> for H_date
`
