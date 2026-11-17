package enrichednow

import (
	"net/http"
	"time"
)

const (
	// DefaultBaseURL is where the API is served.
	DefaultBaseURL = "https://api.enrichednow.com"
	// DefaultTimeout is the time one HTTP attempt may take. A request that asks the API to wait for a run
	// (see [RunParams.Wait]) is given that long on top.
	DefaultTimeout = 30 * time.Second
	// DefaultMaxRetries is how many times a request that is safe to repeat is sent again after a failure.
	DefaultMaxRetries = 2
	// EnvAPIKey is the environment variable [NewClient] reads the API key from.
	EnvAPIKey = "ENRICHEDNOW_API_KEY"
)

// Option configures a [Client].
type Option func(*config)

type config struct {
	apiKey     string
	anonymous  bool
	baseURL    string
	httpClient *http.Client
	timeout    time.Duration
	maxRetries int
	userAgent  string
}

// WithAPIKey sets the API key (en_live_…). Without it the key is read from ENRICHEDNOW_API_KEY.
func WithAPIKey(key string) Option { return func(c *config) { c.apiKey = key } }

// WithoutAPIKey lets [NewClient] build a client that has no key. Such a client can call only what needs none:
// [Client.Health], [ScraperService.List] and [ScraperService.Get].
func WithoutAPIKey() Option { return func(c *config) { c.anonymous = true } }

// WithBaseURL sets the address of the API. The default is [DefaultBaseURL].
func WithBaseURL(baseURL string) Option { return func(c *config) { c.baseURL = baseURL } }

// WithHTTPClient sets the HTTP client requests are sent with. Its own Timeout, if any, applies to every attempt
// in addition to [WithTimeout]; leave it at zero when you use [RunParams.Wait].
func WithHTTPClient(client *http.Client) Option { return func(c *config) { c.httpClient = client } }

// WithTimeout sets the time one HTTP attempt may take. Zero or less means no limit other than the context's.
func WithTimeout(timeout time.Duration) Option { return func(c *config) { c.timeout = timeout } }

// WithMaxRetries sets how many times a request that is safe to repeat is sent again. Zero turns retries off.
func WithMaxRetries(n int) Option {
	return func(c *config) {
		if n < 0 {
			n = 0
		}
		c.maxRetries = n
	}
}

// WithUserAgent adds a product token, such as "my-app/1.2", in front of the SDK's own User-Agent.
func WithUserAgent(product string) Option { return func(c *config) { c.userAgent = product } }
