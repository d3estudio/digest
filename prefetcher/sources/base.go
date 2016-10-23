package sources

// SourceData represents a set of data returned by a Source after processing a
// given URL.
type SourceData map[string]interface{}

// SourceResult represents a SourceResult structure that is stored into the
// database
type SourceResult struct {
	Type string     `json:"type"`
	Data SourceData `json:"data"`
}

// Source defines a set of methods that every Source processor must implement
// in order to be used by the Prefetcher
type Source interface {
	CanHandle(url string) bool
	Process(url string) *SourceResult
}
