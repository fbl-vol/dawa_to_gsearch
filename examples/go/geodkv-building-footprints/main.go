package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"time"
)

const endpoint = "https://graphql.datafordeler.dk/GEODKV/v2"

const query = `query GetBygninger($first: Int, $where: GEODKV_BygningFilterInput, $registreringstid: DafDateTime, $virkningstid: DafDateTime) {
  GEODKV_Bygning(first: $first, where: $where, registreringstid: $registreringstid, virkningstid: $virkningstid) {
    nodes {
      geometri { crs wkt }
      bygningstype
      status
      id_lokalId
      BBRUUID
    }
  }
}`

var apiKeyParamPattern = regexp.MustCompile(`([?&]apiKey=)[^&\s]+`)

func main() {
	easting := flag.Int("easting", 689255, "EPSG:25832 easting")
	northing := flag.Int("northing", 6051787, "EPSG:25832 northing")
	bboxSize := flag.Int("bbox-size-meters", 20, "query box width and height")
	limit := flag.Int("limit", 10, "maximum building candidates")
	timeout := flag.Duration("timeout", 10*time.Second, "request timeout")
	flag.Parse()

	ctx, cancel := context.WithTimeout(context.Background(), *timeout)
	defer cancel()

	payload, err := queryBuildings(ctx, *easting, *northing, *bboxSize, *limit)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	encoder.SetEscapeHTML(false)
	if err := encoder.Encode(payload); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func queryBuildings(ctx context.Context, easting int, northing int, bboxSize int, limit int) (map[string]any, error) {
	apiKey := os.Getenv("DATAFORDELEREN_API_KEY")
	if apiKey == "" {
		apiKey = os.Getenv("DATAFORDELER_GRAPHQL_API_KEY")
	}
	if apiKey == "" {
		return nil, errors.New("DATAFORDELEREN_API_KEY or DATAFORDELER_GRAPHQL_API_KEY is required")
	}

	half := bboxSize / 2
	minE := easting - half
	minN := northing - half
	maxE := easting + half
	maxN := northing + half
	wkt := fmt.Sprintf("POLYGON((%d %d, %d %d, %d %d, %d %d, %d %d))", minE, minN, maxE, minN, maxE, maxN, minE, maxN, minE, minN)
	now := time.Now().UTC().Truncate(time.Second).Format(time.RFC3339)

	body := map[string]any{
		"query": query,
		"variables": map[string]any{
			"first":            limit,
			"registreringstid": now,
			"virkningstid":     now,
			"where": map[string]any{
				"geometri": map[string]any{
					"intersects": map[string]any{
						"crs": 25832,
						"wkt": wkt,
					},
				},
			},
		},
	}

	encodedBody, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	requestURL := endpoint + "?apiKey=" + url.QueryEscape(apiKey)
	request, err := http.NewRequestWithContext(ctx, http.MethodPost, requestURL, bytes.NewReader(encodedBody))
	if err != nil {
		return nil, err
	}
	request.Header.Set("accept", "application/json")
	request.Header.Set("content-type", "application/json")

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return nil, fmt.Errorf("GEODKV request failed: %s", redactAPIKey(err.Error()))
	}
	defer response.Body.Close()

	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	if response.StatusCode < 200 || response.StatusCode >= 300 {
		return nil, fmt.Errorf("GEODKV query failed: HTTP %d %s", response.StatusCode, string(responseBody))
	}

	var payload map[string]any
	if err := json.Unmarshal(responseBody, &payload); err != nil {
		return nil, err
	}
	return payload, nil
}

func redactAPIKey(message string) string {
	return apiKeyParamPattern.ReplaceAllString(message, "${1}<redacted>")
}
