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

Troubleshooting curl commands
curl -X GET "http://localhost:8080/query_history?database=crime"
curl -X GET "http://localhost:8080/databases"
curl -X GET "http://localhost:8080/static/styles.css"

List all Databases
curl -X GET "http://localhost:8181/api/v3/configure/database?format=json" -H "Authorization: Bearer $TOKEN"
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


