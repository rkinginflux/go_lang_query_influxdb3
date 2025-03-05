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
    http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

    // Serve index.html at root "/"
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        data, err := content.ReadFile("index.html")
        if err != nil {
            http.Error(w, "File not found", http.StatusNotFound)
            return
        }
        w.Header().Set("Content-Type", "text/html")
        w.Write(data)
    })

    // Fetch all databases from InfluxDB
    http.HandleFunc("/databases", func(w http.ResponseWriter, r *http.Request) {
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

    // Since InfluxDB returns an ARRAY, use []map[string]interface{} to parse it correctly
    var rawResponse []map[string]interface{}
    if err := json.NewDecoder(resp.Body).Decode(&rawResponse); err != nil {
        log.Println("Error decoding JSON response:", err)
        http.Error(w, "Failed to parse response", http.StatusInternalServerError)
        return
    }

    log.Println("Raw response from InfluxDB:", rawResponse)

    // Extract database names from "iox::database" instead of "name"
    var dbList []string
    for _, db := range rawResponse {
        if name, exists := db["iox::database"].(string); exists {
            dbList = append(dbList, name)
        }
    }

    log.Println("Extracted databases:", dbList)

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]interface{}{"databases": dbList})
})

    // Handle user queries with selected database
    http.HandleFunc("/query", func(w http.ResponseWriter, r *http.Request) {
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
            http.Error(w, "Failed to connect to database", http.StatusInternalServerError)
            return
        }
        defer client.Close()

        startTime := time.Now()
        iterator, err := client.Query(context.Background(), reqBody.Query)
        if err != nil {
            w.Header().Set("Content-Type", "application/json")
            w.WriteHeader(http.StatusBadRequest)
            json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
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

        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(response)
    })

    log.Println("Server is running at http://0.0.0.0:8080 ...")
    log.Fatal(http.ListenAndServe(":8080", nil))
}
