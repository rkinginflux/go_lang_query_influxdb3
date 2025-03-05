# go_lang_query_influxdb3
Web app to spin up a site to query your Influxdb3 databases using Go Lang

Directory structure should look like this. 

```bash
<your directory>
├── go.mod
├── go.sum
├── index.html
├── main.go
└── static
    └── styles.css
-----------------------------------------------------------------------

Troubleshooting curl commands

curl -X GET "http://localhost:8080/query_history?database=ev_cars"

Should look something like...
{
  "queries": [
    "SELECT DISTINCT(query_text) AS \"Query\" FROM system.queries",
    "select DISTINCT(query_text) from system.queries"
  ]
}

-----------------------------------------------------------------------

curl -X GET "http://localhost:8080/databases"

Should look somethng like...

{
  "databases": [
    "crime",
    "ev_cars",
    "support_ear_data"
  ]
}
-----------------------------------------------------------------------

This should display the contents of the styles.css file 
curl -X GET "http://localhost:8080/static/styles.css"
-----------------------------------------------------------------------

List all Databases
curl -X GET "http://db3_server:8181/api/v3/configure/database?format=json" -H "Authorization: Bearer $TOKEN"

Should look someghing like...
[
  {
    "iox::database": "crime"
  },
  {
    "iox::database": "ev_cars"
  },
  {
    "iox::database": "support_ear_data"
  }
]

OR

curl -X GET "http://db3_server:8181/api/v3/configure/database?format=pretty" -H "Authorization: Bearer $TOKEN"

Should look something like...
+------------------+
| iox::database    |
+------------------+
| crime            |
| ev_cars          |
| support_ear_data |
+------------------+

