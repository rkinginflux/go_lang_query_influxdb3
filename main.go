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

const (
    influxHost  = "http://192.168.0.63:8181"
    influxToken = "apiv3_j864z0VmbPEdJIKyeLRLdJI5uagYAHZFgZC2BKuy_WsKxLo8PZ9R-GLWskSCVBp7jTzb16z1uLMijdHnc9MdTQ"
    influxDB    = "support_ear_data"
)

var content embed.FS

func main() {
    client, err := influxdb3.New(influxdb3.ClientConfig{
        Host:     influxHost,
        Token:    influxToken,
        Database: influxDB,
    })
    if err != nil {
        log.Fatalf("Failed to create InfluxDB client: %v", err)
    }
    defer client.Close()

    http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        http.ServeFile(w, r, "index.html")
    })

    http.HandleFunc("/query", func(w http.ResponseWriter, r *http.Request) {
        if r.Method != http.MethodPost {
            http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
            return
        }

        var reqBody struct {
            Query string `json:"query"`
        }
        if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil || reqBody.Query == "" {
            http.Error(w, "Invalid request. Please provide a SQL query.", http.StatusBadRequest)
            return
        }

        startTime := time.Now()
        iterator, err := client.Query(context.Background(), reqBody.Query)
        if err != nil {
            w.Header().Set("Content-Type", "application/json")
            w.WriteHeader(http.StatusBadRequest)
            json.NewEncoder(w).Encode(map[string]string{
                "error": err.Error(),
            })
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

    http.HandleFunc("/query_history", func(w http.ResponseWriter, r *http.Request) {
        query := `SELECT query_text AS "Query" FROM system.queries`
        iterator, err := client.Query(context.Background(), query)
        if err != nil {
            w.Header().Set("Content-Type", "application/json")
            w.WriteHeader(http.StatusBadRequest)
            json.NewEncoder(w).Encode(map[string]string{
                "error": err.Error(),
            })
            return
        }

        var queries []string
        for iterator.Next() {
            row := iterator.Value()
            if queryText, ok := row["Query"].(string); ok {
                queries = append(queries, queryText)
            }
        }

        response := map[string]interface{}{
            "queries": queries,
        }

        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(response)
    })

    log.Println("Server is running at http://0.0.0.0:8080 ...")
    if err := http.ListenAndServe(":8080", nil); err != nil {
        log.Fatalf("Server failed: %v", err)
    }
}
