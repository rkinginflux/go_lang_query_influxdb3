package main

import (
    "context"
    "embed"
    "encoding/json"
    "log"
    "net/http"
    "time"

    influxdb3 "github.com/InfluxCommunity/influxdb3-go/v2/influxdb3"
)

// Constants for InfluxDB
const (
    influxHost  = "http://192.168.0.63:8181"
    influxToken = "apiv3_j864z0VmbPEdJIKyeLRLdJI5uagYAHZFgZC2BKuy_WsKxLo8PZ9R-GLWskSCVBp7jTzb16z1uLMijdHnc9MdTQ"
)

// Embed static files (HTML, CSS, etc.)
//go:embed index.html
var content embed.FS

func main() {
    log.Println("Starting InfluxDB Query Client...")

    // Serve Static Files (CSS, Images, JS)
    http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

    // API Endpoints (Must be registered BEFORE serving index.html)
    http.HandleFunc("/databases", fetchDatabasesHandler)
    http.HandleFunc("/query_history", fetchQueryHistoryHandler)
    http.HandleFunc("/query", executeQueryHandler)

    // Serve index.html at root "/"
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        http.ServeFile(w, r, "index.html")
    })

    log.Println("Server is running at http://0.0.0.0:8080 ...")
    log.Fatal(http.ListenAndServe(":8080", nil))
}

// Fetch available databases
func fetchDatabasesHandler(w http.ResponseWriter, r *http.Request) {
    log.Println("Fetching databases from InfluxDB...")

    req, err := http.NewRequest("GET", influxHost+"/api/v3/configure/database?format=json", nil)
    if err != nil {
        log.Println("Error creating request:", err)
        http.Error(w, "Failed to create request", http.StatusInternalServerError)
        return
    }
    req.Header.Set("Authorization", "Bearer "+influxToken)

    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        log.Println("Error making request to InfluxDB:", err)
        http.Error(w, "Failed to fetch databases", http.StatusInternalServerError)
        return
    }
    defer resp.Body.Close()

    var rawResponse []map[string]interface{}
    if err := json.NewDecoder(resp.Body).Decode(&rawResponse); err != nil {
        log.Println("Error decoding JSON response:", err)
        http.Error(w, "Failed to parse response", http.StatusInternalServerError)
        return
    }

    var dbList []string
    for _, db := range rawResponse {
        if name, exists := db["iox::database"].(string); exists {
            dbList = append(dbList, name)
        }
    }

    log.Println("Extracted databases:", dbList)

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]interface{}{"databases": dbList})
}

// Fetch query history
func fetchQueryHistoryHandler(w http.ResponseWriter, r *http.Request) {
    log.Println("Fetching query history from InfluxDB...")

    client, err := influxdb3.New(influxdb3.ClientConfig{
        Host:  influxHost,
        Token: influxToken,
    })
    if err != nil {
        log.Println("Failed to create InfluxDB client:", err)
        http.Error(w, "Failed to connect to database", http.StatusInternalServerError)
        return
    }
    defer client.Close()

    query := `SELECT DISTINCT(query_text) AS "Query" FROM system.queries`
    iterator, err := client.Query(context.Background(), query)
    if err != nil {
        log.Println("Error executing query:", err)
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    var queries []string
    for iterator.Next() {
        row := iterator.Value()
        if queryText, ok := row["Query"].(string); ok {
            queries = append(queries, queryText)
        }
    }

    log.Println("Extracted queries:", queries)

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]interface{}{"queries": queries})
}

// Execute query on selected database
func executeQueryHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

    var reqBody struct {
        Query    string `json:"query"`
        Database string `json:"database"`
    }
    if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil || reqBody.Query == "" || reqBody.Database == "" {
        http.Error(w, "Invalid request. Provide a query and a database.", http.StatusBadRequest)
        return
    }

    log.Println("Executing query on database:", reqBody.Database)

    client, err := influxdb3.New(influxdb3.ClientConfig{
        Host:     influxHost,
        Token:    influxToken,
        Database: reqBody.Database,
    })
    if err != nil {
        log.Println("Failed to connect to database:", err)
        http.Error(w, "Failed to connect to database", http.StatusInternalServerError)
        return
    }
    defer client.Close()

    startTime := time.Now()
    iterator, err := client.Query(context.Background(), reqBody.Query)
    if err != nil {
        log.Println("Error executing query:", err)
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    var results []map[string]interface{}
    for iterator.Next() {
        row := iterator.Value()
        results = append(results, row)
    }

    queryDuration := time.Since(startTime).Seconds()
    response := map[string]interface{}{
        "duration": queryDuration,
        "results":  results,
    }

    log.Println("Query executed successfully in", queryDuration, "seconds")

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(response)
}
