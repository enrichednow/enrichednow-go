package enrichednow

import (
	"context"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

// DataService reads the datasets. The account must have been given Datasets; without it every call answers
// [ErrPermission]. A data request is charged one data unit for each record it returns, so an iterator that
// fetches page after page is charged for every page it reads.
type DataService struct {
	// Companies are the companies of the dataset.
	Companies *CompanyService
	// Postings are the job postings of the dataset.
	Postings *PostingService
}

// CompanyService reads the companies of the dataset.
type CompanyService struct{ client *Client }

// CompanyListParams filter and page [CompanyService.List].
type CompanyListParams struct {
	// Query is text that must appear in the domain or in a name the company was seen under.
	Query string
	// ATS keeps the companies with this applicant tracking system, e.g. "greenhouse".
	ATS string
	// Limit is the companies per page: 1 to 200. Zero leaves it to the API (50).
	Limit int
	// Cursor is the NextCursor of the page before.
	Cursor string
}

// List returns one page of companies.
func (s *CompanyService) List(ctx context.Context, params *CompanyListParams) (*CompanyPage, error) {
	q := url.Values{}
	if params != nil {
		setString(q, "q", params.Query)
		setString(q, "ats", params.ATS)
		setInt(q, "limit", params.Limit)
		setString(q, "cursor", params.Cursor)
	}
	var out CompanyPage
	if err := s.client.do(ctx, request{method: http.MethodGet, path: "/v1/data/companies", query: q, retry: true}, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// ListAll iterates over the companies, fetching page after page.
func (s *CompanyService) ListAll(ctx context.Context, params *CompanyListParams) *Iter[Company] {
	p := CompanyListParams{}
	if params != nil {
		p = *params
	}
	return newIter(ctx, p.Cursor, func(ctx context.Context, cursor string) ([]Company, string, error) {
		p.Cursor = cursor
		page, err := s.List(ctx, &p)
		if err != nil {
			return nil, "", err
		}
		return page.Companies, page.NextCursor, nil
	})
}

// Get returns one company by its domain, e.g. "figma.com".
func (s *CompanyService) Get(ctx context.Context, domain string) (*Company, error) {
	d, err := pathEscape("company domain", domain)
	if err != nil {
		return nil, err
	}
	var out Company
	if err := s.client.do(ctx, request{method: http.MethodGet, path: "/v1/data/companies/" + d, retry: true}, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// Postings returns one page of a company's postings. params.CompanyDomain is ignored: the company is the one
// named by domain.
func (s *CompanyService) Postings(ctx context.Context, domain string, params *PostingListParams) (*PostingPage, error) {
	d, err := pathEscape("company domain", domain)
	if err != nil {
		return nil, err
	}
	q := params.query()
	q.Del("company_domain")
	var out PostingPage
	if err := s.client.do(ctx, request{method: http.MethodGet, path: "/v1/data/companies/" + d + "/postings", query: q, retry: true}, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// PostingsAll iterates over a company's postings, fetching page after page.
func (s *CompanyService) PostingsAll(ctx context.Context, domain string, params *PostingListParams) *Iter[Posting] {
	p := PostingListParams{}
	if params != nil {
		p = *params
	}
	return newIter(ctx, p.Cursor, func(ctx context.Context, cursor string) ([]Posting, string, error) {
		p.Cursor = cursor
		page, err := s.Postings(ctx, domain, &p)
		if err != nil {
			return nil, "", err
		}
		return page.Postings, page.NextCursor, nil
	})
}

// PostingService reads the job postings of the dataset.
type PostingService struct{ client *Client }

// PostingListParams filter, sort and page a list of postings. The zero value of a field leaves that filter out.
type PostingListParams struct {
	// CompanyDomain is the company's domain; any host or URL of it will do.
	CompanyDomain string
	// Title is words that must all appear in the cleaned title, in any order and any case.
	Title string
	// Location is text that must appear in the location (raw text, city, region or country).
	Location string
	// Country is a country code (ISO 3166-1 alpha-2) or name.
	Country string
	// Remote, when set: true keeps postings marked remote, false those marked not remote. Postings that do
	// not say are in neither.
	Remote *bool
	// EmploymentType is one of full_time, part_time, contract, temporary, internship, volunteer.
	EmploymentType string
	// Seniority is one of c_level, vp, director, manager, lead, senior, junior, intern.
	Seniority string
	// Function is one of data, sales, customer_support, product, design, marketing, finance, legal, people,
	// it, healthcare, education, operations, engineering, skilled_trades, retail_hospitality.
	Function string
	// Status is "open", "closed" or "all" (the default).
	Status string
	// Source is the source's name: an ATS such as "greenhouse", a job board, or "careers_page".
	Source string
	// Provenance is which provenances to return, comma-separated: feed, page, board; or "all". Only feed is
	// returned unless asked.
	Provenance string
	// Unmatched also returns postings not matched to a company.
	Unmatched bool
	// PostedFrom and PostedTo bound the posting date, both included (YYYY-MM-DD).
	PostedFrom string
	PostedTo   string
	// FirstSeenFrom (included) and FirstSeenTo (excluded) bound when the posting was first seen.
	FirstSeenFrom time.Time
	FirstSeenTo   time.Time
	// LastSeenFrom (included) and LastSeenTo (excluded) bound when the posting was last seen.
	LastSeenFrom time.Time
	LastSeenTo   time.Time
	// Sort is "first_seen_at" (the default) or "last_seen_at"; newest first by that time.
	Sort string
	// Limit is the postings per page: 1 to 200. Zero leaves it to the API (50).
	Limit int
	// Cursor is the NextCursor of the page before, with the same other parameters.
	Cursor string
}

func (p *PostingListParams) query() url.Values {
	q := url.Values{}
	if p == nil {
		return q
	}
	setString(q, "company_domain", p.CompanyDomain)
	setString(q, "title", p.Title)
	setString(q, "location", p.Location)
	setString(q, "country", p.Country)
	if p.Remote != nil {
		q.Set("remote", strconv.FormatBool(*p.Remote))
	}
	setString(q, "employment_type", p.EmploymentType)
	setString(q, "seniority", p.Seniority)
	setString(q, "function", p.Function)
	setString(q, "status", p.Status)
	setString(q, "source", p.Source)
	setString(q, "provenance", p.Provenance)
	if p.Unmatched {
		q.Set("unmatched", "true")
	}
	setString(q, "posted_from", p.PostedFrom)
	setString(q, "posted_to", p.PostedTo)
	setTime(q, "first_seen_from", p.FirstSeenFrom)
	setTime(q, "first_seen_to", p.FirstSeenTo)
	setTime(q, "last_seen_from", p.LastSeenFrom)
	setTime(q, "last_seen_to", p.LastSeenTo)
	setString(q, "sort", p.Sort)
	setInt(q, "limit", p.Limit)
	setString(q, "cursor", p.Cursor)
	return q
}

// List returns one page of the postings that match.
func (s *PostingService) List(ctx context.Context, params *PostingListParams) (*PostingPage, error) {
	var out PostingPage
	if err := s.client.do(ctx, request{method: http.MethodGet, path: "/v1/data/postings", query: params.query(), retry: true}, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// ListAll iterates over the postings that match, fetching page after page.
func (s *PostingService) ListAll(ctx context.Context, params *PostingListParams) *Iter[Posting] {
	p := PostingListParams{}
	if params != nil {
		p = *params
	}
	return newIter(ctx, p.Cursor, func(ctx context.Context, cursor string) ([]Posting, string, error) {
		p.Cursor = cursor
		page, err := s.List(ctx, &p)
		if err != nil {
			return nil, "", err
		}
		return page.Postings, page.NextCursor, nil
	})
}

// Get returns one posting with its description and its history.
func (s *PostingService) Get(ctx context.Context, postingID string) (*PostingDetail, error) {
	id, err := pathEscape("posting id", postingID)
	if err != nil {
		return nil, err
	}
	var out PostingDetail
	if err := s.client.do(ctx, request{method: http.MethodGet, path: "/v1/data/postings/" + id, retry: true}, &out); err != nil {
		return nil, err
	}
	return &out, nil
}
