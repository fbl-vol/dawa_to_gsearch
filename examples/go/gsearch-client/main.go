package main

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"
	"sync"
	"time"
)

const baseURL = "https://api.dataforsyningen.dk/rest/gsearch/v2.0"

type result map[string]any

var httpClient = &http.Client{Transport: gsearchTransport()}
var tokenParamPattern = regexp.MustCompile(`([?&](?:token|TOKEN)=)[^&\s]+`)

type selection struct {
	Provider    string   `json:"provider"`
	Resource    string   `json:"resource"`
	ID          string   `json:"id,omitempty"`
	HusnummerID string   `json:"husnummerId,omitempty"`
	Label       string   `json:"label,omitempty"`
	Kommunekode string   `json:"kommunekode,omitempty"`
	Vejkode     string   `json:"vejkode,omitempty"`
	Postnummer  string   `json:"postnummer,omitempty"`
	Longitude   *float64 `json:"longitude,omitempty"`
	Latitude    *float64 `json:"latitude,omitempty"`
	Raw         result   `json:"raw"`
}

func main() {
	limit := flag.Int("limit", 5, "maximum results per resource")
	srid := flag.Int("srid", 4326, "response geometry SRID")
	timeout := flag.Duration("timeout", 10*time.Second, "request timeout")
	flag.Parse()

	query := "Søbakkevej 8, Tilst"
	if flag.NArg() > 0 {
		query = strings.Join(flag.Args(), " ")
	}

	ctx, cancel := context.WithTimeout(context.Background(), *timeout)
	defer cancel()

	suggestions, err := searchAddressSuggestions(ctx, query, *limit, *srid)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	encoder.SetEscapeHTML(false)
	if err := encoder.Encode(suggestions); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func searchAddressSuggestions(ctx context.Context, query string, limit int, srid int) ([]selection, error) {
	var wg sync.WaitGroup
	var husnummer []result
	var adresse []result
	var husnummerErr error
	var adresseErr error

	wg.Add(2)
	go func() {
		defer wg.Done()
		husnummer, husnummerErr = searchGSearch(ctx, "husnummer", query, limit, srid)
	}()
	go func() {
		defer wg.Done()
		adresse, adresseErr = searchGSearch(ctx, "adresse", query, limit, srid)
	}()
	wg.Wait()

	if husnummerErr != nil && adresseErr != nil {
		return nil, fmt.Errorf("%v; %v", husnummerErr, adresseErr)
	}
	if husnummerErr != nil {
		return nil, husnummerErr
	}
	if adresseErr != nil {
		return nil, adresseErr
	}
	return mergeAddressResults(husnummer, adresse), nil
}

func searchGSearch(ctx context.Context, resource string, query string, limit int, srid int) ([]result, error) {
	token := os.Getenv("GSEARCH_TOKEN")
	if token == "" {
		return nil, errors.New("GSEARCH_TOKEN is required")
	}

	endpoint, err := url.Parse(baseURL + "/" + resource)
	if err != nil {
		return nil, err
	}
	params := endpoint.Query()
	params.Set("token", token)
	params.Set("q", query)
	params.Set("limit", fmt.Sprint(limit))
	params.Set("srid", fmt.Sprint(srid))
	endpoint.RawQuery = params.Encode()

	request, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint.String(), nil)
	if err != nil {
		return nil, err
	}
	request.Header.Set("accept", "application/json")

	response, err := httpClient.Do(request)
	if err != nil {
		return nil, fmt.Errorf("GSearch %s request failed: %s", resource, redactToken(err.Error()))
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	if response.StatusCode < 200 || response.StatusCode >= 300 {
		return nil, fmt.Errorf("GSearch %s failed: HTTP %d %s", resource, response.StatusCode, string(body))
	}

	var payload []result
	if err := json.Unmarshal(body, &payload); err != nil {
		return nil, err
	}
	return payload, nil
}

func gsearchTransport() *http.Transport {
	transport := http.DefaultTransport.(*http.Transport).Clone()
	// The Dataforsyningen gateway can close Go HTTP/2 requests with unexpected EOF.
	transport.ForceAttemptHTTP2 = false
	transport.TLSClientConfig = &tls.Config{NextProtos: []string{"http/1.1"}}
	transport.TLSNextProto = map[string]func(string, *tls.Conn) http.RoundTripper{}
	return transport
}

func redactToken(message string) string {
	return tokenParamPattern.ReplaceAllString(message, "${1}<redacted>")
}

func mergeAddressResults(husnummer []result, adresse []result) []selection {
	husnummerIDs := map[string]string{}
	for _, item := range husnummer {
		key := addressKey(item, "husnummertekst")
		if key != "" {
			husnummerIDs[key] = stringField(item, "id")
		}
	}

	seen := map[string]bool{}
	merged := []selection{}

	for _, item := range adresse {
		key := addressKey(item, "husnummer")
		normalized := normalizeResult("adresse", item, husnummerIDs[key])
		if normalized.Label != "" && !seen[normalized.Label] {
			seen[normalized.Label] = true
			merged = append(merged, normalized)
		}
	}

	for _, item := range husnummer {
		id := stringField(item, "id")
		normalized := normalizeResult("husnummer", item, id)
		if normalized.Label != "" && !seen[normalized.Label] {
			seen[normalized.Label] = true
			merged = append(merged, normalized)
		}
	}

	return merged
}

func normalizeResult(resource string, item result, husnummerID string) selection {
	longitude, latitude := firstLonLat(item["geometri"])
	return selection{
		Provider:    "dataforsyningen-gsearch",
		Resource:    resource,
		ID:          stringField(item, "id"),
		HusnummerID: husnummerID,
		Label:       stringField(item, "visningstekst"),
		Kommunekode: stringField(item, "kommunekode"),
		Vejkode:     stringField(item, "vejkode"),
		Postnummer:  stringField(item, "postnummer"),
		Longitude:   longitude,
		Latitude:    latitude,
		Raw:         item,
	}
}

func addressKey(item result, houseNumberField string) string {
	parts := []string{
		stringField(item, "kommunekode"),
		stringField(item, "vejkode"),
		stringField(item, houseNumberField),
	}
	for _, part := range parts {
		if part == "" {
			return ""
		}
	}
	return strings.Join(parts, ":")
}

func stringField(item result, key string) string {
	value, ok := item[key].(string)
	if !ok {
		return ""
	}
	return value
}

func firstLonLat(geometry any) (*float64, *float64) {
	geometryMap, ok := geometry.(map[string]any)
	if !ok {
		return nil, nil
	}
	position, ok := firstPosition(geometryMap["coordinates"])
	if !ok {
		return nil, nil
	}
	return &position[0], &position[1]
}

func firstPosition(value any) ([2]float64, bool) {
	values, ok := value.([]any)
	if !ok {
		return [2]float64{}, false
	}
	if len(values) >= 2 {
		x, xOK := values[0].(float64)
		y, yOK := values[1].(float64)
		if xOK && yOK {
			return [2]float64{x, y}, true
		}
	}
	for _, child := range values {
		position, ok := firstPosition(child)
		if ok {
			return position, true
		}
	}
	return [2]float64{}, false
}
