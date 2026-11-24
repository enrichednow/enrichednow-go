// Code generated from the API's formula parameter schemas. DO NOT EDIT.

package enrichednow

// FormulaParams is what the parameter types of this package have in common: the formula they are for.
// [FormulaService.Run], [FormulaService.Estimate] and [SavedFormulaService.Create] check that it is the formula named.
type FormulaParams interface {
	FormulaID() string
}

// IntRange is a range of whole numbers, both ends included.
type IntRange struct {
	Min int `json:"min"`
	Max int `json:"max"`
}

// FormulaIDs are the formulas this version of the package has a parameter type for.
var FormulaIDs = []string{
	"account-hiring-check",
	"account-research",
	"acquisition-universe",
	"ad-activity-trigger",
	"ad-creative-research",
	"agency-ad-mining",
	"agency-lists",
	"announcement-trigger",
	"app-intelligence",
	"award-teaming",
	"broken-link-prospects",
	"broker-directory",
	"business-for-sale-alerts",
	"candidate-reverse-marketing",
	"champion-tracker",
	"company-diligence",
	"competitor-price-tracker",
	"competitor-review-mining",
	"compound-signal-score",
	"contractor-license-watch",
	"crowdfunding-leads",
	"deal-site-monitor",
	"directory-to-leads",
	"ecommerce-panel",
	"enterprise-ats-sweep",
	"event-organizer-leads",
	"exhibitor-prospects",
	"footprint-prospects",
	"freelance-gig-feed",
	"funded-and-hiring",
	"hiring-signal-outbound",
	"hiring-velocity",
	"home-services-pros",
	"hotel-rate-shop",
	"influencer-discovery",
	"job-order-hunter",
	"job-seeker-targets",
	"knowledge-base-crawl",
	"landlord-agent-leads",
	"launch-batch-lists",
	"link-intersect",
	"listing-compliance",
	"literature-sweep",
	"local-business-leads",
	"long-open-roles",
	"lookalike-companies",
	"marketplace-flips",
	"menu-price-tracker",
	"out-of-stock-opportunity",
	"people-list-with-openers",
	"permit-status-watch",
	"permit-trade-leads",
	"price-seller-watch",
	"property-comps",
	"property-distress-signals",
	"rank-tracker",
	"reverse-sourcing",
	"review-mining",
	"review-velocity",
	"salary-benchmark",
	"sec-filing-monitor",
	"seller-identity",
	"shopify-niche-research",
	"social-sentiment",
	"stale-listings",
	"storefront-research",
	"str-underwriting",
	"supplier-discovery",
	"talent-sourcing",
	"tender-alerts",
	"tender-change-watch",
	"unlinked-mentions",
	"used-car-arbitrage",
	"vc-deal-sourcing",
	"web-rank-tracker",
	"wholesale-catalog-check",
}

// AccountHiringCheckParams are the parameters of Account Hiring Check (`account-hiring-check`): Every account in
// your CRM flagged hiring or not, with its open roles, departments and the roles you care about.
//
// Accounts is required. Any other field left nil or empty is not sent, and the formula's default applies; [Ptr]
// makes the pointer of a value.
type AccountHiringCheckParams struct {
	// Accounts: Each account's website domain, as in your CRM (acme.com). Upload a CSV for many. 1 to 100 items.
	// Required.
	Accounts []string `json:"accounts"`

	// Roles you care about: Openings whose title holds one of these as a whole word or phrase are counted apart,
	// such as Data Engineer. At most 20 items. Default: [].
	RolesOfInterest []string `json:"roles_of_interest,omitempty"`

	// New within (days): An opening posted in the last this many days counts as new. From 1 to 180. Default: 30.
	NewWithinDays *int `json:"new_within_days,omitempty"`

	// Hiring at this many openings: An account with at least this many openings is flagged hiring. At least 1.
	// Default: 1.
	HiringThreshold *int `json:"hiring_threshold,omitempty"`

	// Count only tracking-system openings: Counts only openings an applicant tracking system lists. What a careers
	// page's own text shows is counted apart, as unconfirmed. Default: true.
	VerifiedOnly *bool `json:"verified_only,omitempty"`

	// Most rows: One row per account; the run stops at this many, most active first. From 1 to 5000. Default:
	// 1000.
	MaxResults *int `json:"max_results,omitempty"`
}

// FormulaID is "account-hiring-check".
func (AccountHiringCheckParams) FormulaID() string { return "account-hiring-check" }

// AccountResearchParams are the parameters of Account Research (`account-research`): For each company domain, what
// its site and company page say, a line of research and a compliment built only from facts.
//
// Domains is required. Any other field left nil or empty is not sent, and the formula's default applies; [Ptr]
// makes the pointer of a value.
type AccountResearchParams struct {
	// Company websites: Each company's domain or website address, typed or from a CSV. 1 to 500 items. Required.
	Domains []string `json:"domains"`

	// Read their LinkedIn company page: When the site links to one; adds industry, size, headquarters and when it
	// was founded. Default: true.
	ReadCompanyPage *bool `json:"read_company_page,omitempty"`

	// Read the technologies of their site: The site's platform, analytics and tools, as its own pages reveal them.
	// Default: true.
	ReadTechnologies *bool `json:"read_technologies,omitempty"`

	// Most rows: The run stops at this many companies. From 1 to 500. Default: 500.
	MaxResults *int `json:"max_results,omitempty"`
}

// FormulaID is "account-research".
func (AccountResearchParams) FormulaID() string { return "account-research" }

// AcquisitionUniverseParams are the parameters of Acquisition Universe (`acquisition-universe`): Every business of
// your acquisition thesis in your metros, off the market from Google Maps and on it from sale listings.
//
// Thesis and Metros are required. Any other field left nil or empty is not sent, and the formula's default applies;
// [Ptr] makes the pointer of a value.
type AcquisitionUniverseParams struct {
	// Kinds of business: What you would buy, as you would search Google Maps for it, such as HVAC contractor or
	// landscaping company. 1 to 5 items. Required.
	Thesis []string `json:"thesis"`

	// Metros: A city or metro for each place, with its state. 1 to 10 items. Required.
	Metros []string `json:"metros"`

	// States for listings: The US states, by name, whose businesses for sale are read. Empty reads none. At most
	// 10 items. Default: [].
	States []string `json:"states,omitempty"`

	// Words a listing must hold: Whole words a for-sale listing's title or summary must hold, such as HVAC or air
	// conditioning. Empty uses the kinds of business. At most 10 items. Default: [].
	ListingWords []string `json:"listing_words,omitempty"`

	// Places read per search: How many places Google Maps is read for, for each kind of business in each metro.
	// From 10 to 300. Default: 60.
	PlacesPerSearch *int `json:"places_per_search,omitempty"`

	// Fewest reviews: Places with fewer reviews are left out; a business with none may be too small or not
	// trading. At least 0. Default: 0.
	MinReviews *int `json:"min_reviews,omitempty"`

	// Most reviews: Places with more reviews are left out; the largest are rarely owner-run. At least 0. Default:
	// 1000000.
	MaxReviews *int `json:"max_reviews,omitempty"`

	// Leave out chains: A name found five or more times in the results is taken for a chain and left out. Default:
	// true.
	ExcludeChains *bool `json:"exclude_chains,omitempty"`

	// Larger companies without their name: BusinessesForSale.com's M&A Vault describes larger companies without
	// naming them, from any country. Default: false.
	IncludeVault *bool `json:"include_vault,omitempty"`

	// Read each website for an address: The address each business's own website publishes, checked. Off saves one
	// read a business. Default: true.
	FindEmail *bool `json:"find_email,omitempty"`

	// Show a business again after (days): For a saved formula, how long before a delivered business or listing may
	// come back. From 1 to 365. Default: 180.
	ResurfaceAfterDays *int `json:"resurface_after_days,omitempty"`

	// Most rows: The run stops at this many rows, listings first, then the places with the most reviews. From 1 to
	// 5000. Default: 300.
	MaxResults *int `json:"max_results,omitempty"`
}

// FormulaID is "acquisition-universe".
func (AcquisitionUniverseParams) FormulaID() string { return "acquisition-universe" }

// AdActivityTriggerParams are the parameters of Ad-Activity Trigger (`ad-activity-trigger`): Companies running ads
// about your topic on Meta and LinkedIn now, sized up, with their marketing lead to contact.
//
// Keywords is required. Any other field left nil or empty is not sent, and the formula's default applies; [Ptr]
// makes the pointer of a value.
type AdActivityTriggerParams struct {
	// What the ads are about: Words the ads say, each searched in each library, such as "payroll software". 1 to 5
	// items. Required.
	Keywords []string `json:"keywords"`

	// Countries: Where the ads are shown, as two-letter codes such as US or GB. 1 to 5 items. Default: ["US"].
	Countries []string `json:"countries,omitempty"`

	// Ad libraries: Meta covers Facebook and Instagram. Each one of "meta", "linkedin". At least 1 item. Default:
	// ["meta", "linkedin"].
	Platforms []string `json:"platforms,omitempty"`

	// Ads started within (days): Meta says when each ad started; older ones are left out. LinkedIn's library does
	// not say, and its ads are kept. From 1 to 365. Default: 90.
	StartedWithinDays *int `json:"started_within_days,omitempty"`

	// Ads read per search: How many ads to read for each phrase, country and library. From 5 to 100. Default: 30.
	AdsPerSearch *int `json:"ads_per_search,omitempty"`

	// Least ads: Ads a company needs to be running to be listed. At least 1. Default: 1.
	MinAds *int `json:"min_ads,omitempty"`

	// Company size (employees): Companies whose size is not known are kept. Default: {"min": 2, "max": 10000}.
	CompanySize *IntRange `json:"company_size,omitempty"`

	// Companies to leave out: Names or domains, such as your existing customers. At most 5000 items. Default: [].
	ExcludeCompanies []string `json:"exclude_companies,omitempty"`

	// Who to contact: The roles a connected people provider looks for at each company. Default: ["Chief Marketing
	// Officer", "VP Marketing", "Head of Marketing", "Marketing Director", "Head of Growth"].
	BuyerTitles []string `json:"buyer_titles,omitempty"`

	// Contacts per company: 0 keeps one row per company. From 0 to 5. Default: 2.
	ContactsPerCompany *int `json:"contacts_per_company,omitempty"`

	// Only addresses that pass the check: Keeps addresses whose domain takes mail, and those a verification
	// provider confirmed. Default: true.
	RequireVerifiedEmail *bool `json:"require_verified_email,omitempty"`

	// Most rows: The run stops at this many rows, the most active advertisers first. From 1 to 2000. Default: 200.
	MaxResults *int `json:"max_results,omitempty"`
}

// FormulaID is "ad-activity-trigger".
func (AdActivityTriggerParams) FormulaID() string { return "ad-activity-trigger" }

// AdCreativeResearchParams are the parameters of Ad Creative Research (`ad-creative-research`): A competitor's ads
// on Meta, Google Ads and LinkedIn, with copy, landing page and how long each has run, longest first.
//
// Competitors is required. Any other field left nil or empty is not sent, and the formula's default applies; [Ptr]
// makes the pointer of a value.
type AdCreativeResearchParams struct {
	// Competitors: Each competitor's name as it advertises, such as Canva or HubSpot. Ads are kept where the
	// advertiser's name holds it as a whole word. 1 to 10 items. Required.
	Competitors []string `json:"competitors"`

	// Their websites: The competitors' websites, such as canva.com. Google Ads are found by the website they lead
	// to, so Google Ads is searched only when you give one. At most 10 items. Default: [].
	Domains []string `json:"domains,omitempty"`

	// Ad libraries: Meta covers Facebook, Instagram, Messenger and Threads; Google Ads covers Search, YouTube and
	// the display network. Each one of "meta", "google", "linkedin". At least 1 item. Default: ["meta", "google",
	// "linkedin"].
	Libraries []string `json:"libraries,omitempty"`

	// Country: Only ads shown in this country, as a two-letter code such as US or DE. Leave empty for every
	// country. Default: "".
	Country *string `json:"country,omitempty"`

	// Only ads running now: Leaves out ads that have stopped. Off, a stopped ad still shows how long an angle
	// lasted. Default: false.
	RunningOnly *bool `json:"running_only,omitempty"`

	// Least days running: Only ads that have run at least this many days; 30 keeps the angles an advertiser kept
	// paying for. From 0 to 3650. Default: 0.
	MinDaysRunning *int `json:"min_days_running,omitempty"`

	// Facebook Page ids: For an exact Meta search, the competitor's Page id (the number after view_all_page_id= in
	// the Ad Library's address). Every ad of the Page is read. At most 10 items. Default: [].
	FacebookPageIds []string `json:"facebook_page_ids,omitempty"`

	// Ads read per search: The most ads read for each competitor in each library. Each ad read is one unit. From 5
	// to 120. Default: 30.
	AdsPerSearch *int `json:"ads_per_search,omitempty"`

	// Read each LinkedIn ad's page: LinkedIn's search shows no dates or landing page; each ad's own page has them,
	// for one unit an ad. Default: true.
	LinkedinAdPages *bool `json:"linkedin_ad_pages,omitempty"`

	// Most rows: The run stops at this many ads, longest running first. From 1 to 2000. Default: 200.
	MaxResults *int `json:"max_results,omitempty"`
}

// FormulaID is "ad-creative-research".
func (AdCreativeResearchParams) FormulaID() string { return "ad-creative-research" }

// AgencyAdMiningParams are the parameters of Agency Ad Mining (`agency-ad-mining`): Jobs competing agencies
// advertise, what they say about the client, and the employers posting the same job directly.
//
// JobTitles and Locations are required. Any other field left nil or empty is not sent, and the formula's default
// applies; [Ptr] makes the pointer of a value.
type AgencyAdMiningParams struct {
	// Job titles: The roles you work, each searched as you type it. 1 to 10 items. Required.
	JobTitles []string `json:"job_titles"`

	// Locations: A city, region or country for each market you watch. 1 to 10 items. Required.
	Locations []string `json:"locations"`

	// Posted within (days): Only postings from the last this many days. From 1 to 30. Default: 7.
	PostedWithinDays *int `json:"posted_within_days,omitempty"`

	// Only these agencies: Names of the competing agencies to watch. Leave empty for every agency. At most 50
	// items. Default: [].
	Agencies []string `json:"agencies,omitempty"`

	// Job boards: Which boards to search. Reed and Dice say who placed a posting; Reed and Totaljobs cover the UK,
	// Dice technology roles in the US. Each one of "linkedin", "indeed", "reed", "totaljobs", "dice". At least 1
	// item. Default: ["linkedin", "reed", "dice"].
	Sources []string `json:"sources,omitempty"`

	// Jobs read per search: How many postings to read for each title, place and board, agencies' and employers'
	// alike. From 5 to 200. Default: 50.
	ResultsPerSearch *int `json:"results_per_search,omitempty"`

	// Agency postings read in full: The newest agency postings whose full text is read for what they say about the
	// client. From 5 to 300. Default: 60.
	PostingsRead *int `json:"postings_read,omitempty"`

	// Words that describe the client: The first of these a posting uses is quoted with the sentence around it.
	// Default: ["our client", "my client", "on behalf of", "client is", "client based in", "working with a", "a
	// leading", "well-established"].
	ClientPhrases []string `json:"client_phrases,omitempty"`

	// Most rows: The run stops at this many rows, best first. From 1 to 2000. Default: 100.
	MaxResults *int `json:"max_results,omitempty"`
}

// FormulaID is "agency-ad-mining".
func (AgencyAdMiningParams) FormulaID() string { return "agency-ad-mining" }

// AgencyListsParams are the parameters of Agency Lists (`agency-lists`): Agencies by service and place with size,
// age, rates, minimum project, services and reviews, for partnering or buying.
//
// Directories is required. Any other field left nil or empty is not sent, and the formula's default applies; [Ptr]
// makes the pointer of a value.
type AgencyListsParams struct {
	// Clutch directories: Each as in its address on clutch.co, the service then the place: agencies/seo,
	// web-developers/new-york, uk/agencies/digital-marketing/london, or the whole address. 1 to 10 items.
	// Required.
	Directories []string `json:"directories"`

	// Pages to read: Pages of each directory to read; a page lists about 50 to 80 agencies. From 1 to 10. Default:
	// 2.
	Pages *int `json:"pages,omitempty"`

	// What the list is for: Partner ranks agencies by rating, reviews and a published email. Acquisition ranks
	// established agencies first: ten years or more in business, 10 to 249 staff, and a body of reviews. One of
	// "partner", "acquisition". Default: "partner".
	Purpose *string `json:"purpose,omitempty"`

	// Head count: Agencies whose head count on Clutch falls outside this range are left out; an agency that does
	// not say is kept. Default: {"min": 2, "max": 1000}.
	Employees *IntRange `json:"employees,omitempty"`

	// Founded in or before: Only agencies founded in this year or earlier, such as 2015 for ten years in business.
	// Empty for any. From 1900 to 2100. Default: null.
	FoundedBefore map[string]any `json:"founded_before,omitempty"`

	// Highest hourly rate (USD): Agencies whose lowest hourly rate is above this are left out. 0 for any. From 0
	// to 1000. Default: 0.
	MaxHourlyRate *int `json:"max_hourly_rate,omitempty"`

	// Largest minimum project (USD): Agencies that take no project smaller than this are left out. 0 for any. From
	// 0 to 10000000. Default: 0.
	MaxMinProject *int `json:"max_min_project,omitempty"`

	// Sells one of these services: Words from Clutch's service names, such as SEO, Pay Per Click or Web
	// Development. Leave empty for any. At most 20 items. Default: [].
	Services []string `json:"services,omitempty"`

	// Office in these countries: Two-letter country codes (GB, US). A directory of a place also lists sponsors
	// based elsewhere; this leaves them out. At most 50 items. Default: [].
	Countries []string `json:"countries,omitempty"`

	// Lowest rating: Agencies rated below this on Clutch are left out. From 0 to 5. Default: 0.
	MinRating *float64 `json:"min_rating,omitempty"`

	// Fewest reviews: Agencies with fewer Clutch reviews are left out. At least 0. Default: 0.
	MinReviews *int `json:"min_reviews,omitempty"`

	// Only agencies with an email: Keeps the agencies whose own website publishes an email address that passes the
	// check. Default: false.
	RequireEmail *bool `json:"require_email,omitempty"`

	// Agencies to leave out: Names or domains, such as agencies you already work with. At most 5000 items.
	// Default: [].
	ExcludeCompanies []string `json:"exclude_companies,omitempty"`

	// Most rows: The run stops at this many agencies, best first. Each agency's profile and website are read, so
	// this sets the cost. From 1 to 2000. Default: 100.
	MaxResults *int `json:"max_results,omitempty"`
}

// FormulaID is "agency-lists".
func (AgencyListsParams) FormulaID() string { return "agency-lists" }

// AnnouncementTriggerParams are the parameters of Funding and Announcement Trigger (`announcement-trigger`):
// Companies that just announced funding, an expansion, a launch or a senior hire, with who to contact and the fact.
//
// EventTypes is required. Any other field left nil or empty is not sent, and the formula's default applies; [Ptr]
// makes the pointer of a value.
type AnnouncementTriggerParams struct {
	// Kinds of news: Which announcements to look for. Each one of "funding", "expansion", "product_launch",
	// "executive_hire", "acquisition". At least 1 item. Required.
	EventTypes []string `json:"event_types"`

	// Words in the news: Only releases that hold one of these words, such as fintech or healthcare. Leave empty
	// for any. At most 5 items. Default: [].
	Keywords []string `json:"keywords,omitempty"`

	// Published within (days): Only releases from the last this many days. From 1 to 30. Default: 3.
	PublishedWithinDays *int `json:"published_within_days,omitempty"`

	// News sources: The press-release wires, and Google News for stories as well as releases. Each one of
	// "prnewswire", "globenewswire", "google_news". At least 1 item. Default: ["prnewswire", "globenewswire"].
	Sources []string `json:"sources,omitempty"`

	// Releases read in full: How many of the newest matching releases are opened to read who, how much and who
	// invested. From 5 to 200. Default: 60.
	ReleasesToRead *int `json:"releases_to_read,omitempty"`

	// Places: Countries, states or cities the company or the release is from. Leave empty for anywhere. Default:
	// [].
	Regions []string `json:"regions,omitempty"`

	// Company size (employees): Companies whose size is not known are kept. Default: {"min": 1, "max": 10000}.
	CompanySize *IntRange `json:"company_size,omitempty"`

	// Only these industries: Leave empty for any. Default: [].
	IndustriesInclude []string `json:"industries_include,omitempty"`

	// Companies to leave out: Names or domains, such as customers. At most 5000 items. Default: [].
	ExcludeCompanies []string `json:"exclude_companies,omitempty"`

	// Who to contact: The roles of the people to find at each company. Default: ["Chief Executive Officer",
	// "Founder", "VP Sales"].
	BuyerTitles []string `json:"buyer_titles,omitempty"`

	// Contact the new executive: For a senior hire, look for the person in the new role rather than the roles
	// above. Default: false.
	ContactNewExecutive *bool `json:"contact_new_executive,omitempty"`

	// Contacts per company: 0 turns contacts off. From 0 to 5. Default: 1.
	ContactsPerCompany *int `json:"contacts_per_company,omitempty"`

	// Only addresses that pass the check: Keeps addresses whose domain takes mail, and those a verification
	// provider confirmed. Default: true.
	RequireVerifiedEmail *bool `json:"require_verified_email,omitempty"`

	// Most rows: The run stops at this many rows, newest first. From 1 to 1000. Default: 100.
	MaxResults *int `json:"max_results,omitempty"`
}

// FormulaID is "announcement-trigger".
func (AnnouncementTriggerParams) FormulaID() string { return "announcement-trigger" }

// AppIntelligenceParams are the parameters of App Intelligence (`app-intelligence`): Ratings, versions, release
// notes and chart ranks of apps on the App Store and Google Play, with what changed since last run.
//
// Apps is required. Any other field left nil or empty is not sent, and the formula's default applies; [Ptr] makes
// the pointer of a value.
type AppIntelligenceParams struct {
	// Apps: App Store apps by their number (the digits after id in the app's address) or their App Store link;
	// Google Play apps by their package name, such as com.spotify.music. 1 to 100 items. Required.
	Apps []string `json:"apps"`

	// Search words: Also the apps that rank in the stores' search for these words, such as habit tracker. At most
	// 10 items. Default: [].
	SearchTerms []string `json:"search_terms,omitempty"`

	// Top chart: Also the top apps of this category's chart, with their rank. An app you named gets its rank when
	// it is on the chart. One of "none", "all", "books", "business", "education", "entertainment", "finance",
	// "food_and_drink", "games", "health_and_fitness", "lifestyle", "medical", "music", "navigation", "news",
	// "photo_and_video", "productivity", "shopping", "social_networking", "sports", "travel", "utilities",
	// "weather". Default: "none".
	ChartCategory *string `json:"chart_category,omitempty"`

	// Stores: The stores searched and charted. Apps you name are read in their own store. Each one of "app_store",
	// "google_play". At least 1 item. Default: ["app_store", "google_play"].
	Stores []string `json:"stores,omitempty"`

	// Chart: Top free, top paid or top grossing. One of "free", "paid", "grossing". Default: "free".
	Chart *string `json:"chart,omitempty"`

	// Chart places read: How many of the chart's top apps are read in each store. From 5 to 100. Default: 25.
	ChartSize *int `json:"chart_size,omitempty"`

	// Apps read per search: How many of the top results of each search word are read in each store. From 1 to 50.
	// Default: 10.
	AppsPerSearch *int `json:"apps_per_search,omitempty"`

	// Store country: The country's store, as a two-letter code such as us, gb or de. Ratings and charts differ by
	// country. Default: "us".
	Country *string `json:"country,omitempty"`

	// Lowest rating: Leaves out apps rated below this. Apps you name are always kept. From 0 to 5. Default: 0.
	MinRating *float64 `json:"min_rating,omitempty"`

	// Recent reviews read per app: Reads this many of each app's newest reviews for its recent rating and the
	// share of one- and two-star reviews. Each review read is one unit; 0 turns it off. From 0 to 100. Default: 0.
	RecentReviews *int `json:"recent_reviews,omitempty"`

	// Most rows: The run stops at this many apps, those with the biggest changes first. From 1 to 1000. Default:
	// 200.
	MaxResults *int `json:"max_results,omitempty"`
}

// FormulaID is "app-intelligence".
func (AppIntelligenceParams) FormulaID() string { return "app-intelligence" }

// AwardTeamingParams are the parameters of Award Teaming Leads (`award-teaming`): Companies that just won
// government contracts in your field, with what they won, their website and the email they publish.
//
// Any other field left nil or empty is not sent, and the formula's default applies; [Ptr] makes the pointer of a
// value.
type AwardTeamingParams struct {
	// NAICS codes: The industry codes of the work, 2 to 6 digits, such as 541512 (computer systems design). Each
	// is searched on its own. At most 10 items. Default: [].
	NaicsCodes []string `json:"naics_codes,omitempty"`

	// Words: Words in the award's description, such as cybersecurity or janitorial. Each is searched on its own.
	// At most 10 items. Default: [].
	Keywords []string `json:"keywords,omitempty"`

	// Awarding departments: US departments by their full name, such as Department of Defense or Department of
	// Veterans Affairs. At most 5 items. Default: [].
	Agencies []string `json:"agencies,omitempty"`

	// Awarded within (days): Awards whose work started, or whose notice was published, in the last this many days.
	// Recently signed work that starts later is kept. From 7 to 364. Default: 30.
	AwardedWithinDays *int `json:"awarded_within_days,omitempty"`

	// Where awards are read: USAspending.gov (US federal contracts), SAM.gov award notices where offered, and TED
	// (EU contract awards). Each one of "usaspending", "sam_gov", "ted". At least 1 item. Default: ["usaspending",
	// "sam_gov"].
	Sources []string `json:"sources,omitempty"`

	// Places of performance: Only US awards for work done in these states, as two letters, such as VA or TX. At
	// most 10 items. Default: [].
	USStates []string `json:"us_states,omitempty"`

	// Buyers' countries (TED): Only EU awards made by buyers in these countries, as three letters, such as DEU or
	// FRA. At most 10 items. Default: [].
	Countries []string `json:"countries,omitempty"`

	// Smallest award: Awards below this amount, in the award's currency, are left out. An award that states no
	// amount is kept. At least 0. Default: 100000.
	MinAwardAmount *float64 `json:"min_award_amount,omitempty"`

	// Winner's size: Small businesses or large primes, as USAspending records them. A winner it has no record of
	// is kept only for any. One of "any", "small_business", "other_than_small". Default: "any".
	WinnerSize *string `json:"winner_size,omitempty"`

	// Companies to leave out: Names, such as your competitors or the partners you already have. At most 1000
	// items. Default: [].
	ExcludeCompanies []string `json:"exclude_companies,omitempty"`

	// Only companies with an email: Keeps the companies whose own website publishes an email address that passes
	// the check. Default: true.
	RequireEmail *bool `json:"require_email,omitempty"`

	// Pages read per search: Each page is the next 100 awards, newest first. Raise it for a busy code or a long
	// window. From 1 to 5. Default: 2.
	PagesPerSearch *int `json:"pages_per_search,omitempty"`

	// Most rows: The run stops at this many companies, the biggest winners first. From 1 to 1000. Default: 100.
	MaxResults *int `json:"max_results,omitempty"`
}

// FormulaID is "award-teaming".
func (AwardTeamingParams) FormulaID() string { return "award-teaming" }

// BrokenLinkProspectsParams are the parameters of Broken Link Prospects (`broken-link-prospects`): Resource pages
// on your topic that link to pages or sites that are gone, each with the dead link and the site's email.
//
// Topic is required. Any other field left nil or empty is not sent, and the formula's default applies; [Ptr] makes
// the pointer of a value.
type BrokenLinkProspectsParams struct {
	// Topic: The subject of the resource pages, as people search for it, such as genealogy or home schooling.
	// Required.
	Topic string `json:"topic"`

	// Footprints: Words searched beside the topic to find resource pages, each its own search. Quote a phrase to
	// find it exactly. 1 to 12 items. Default: ["\"useful links\"", "\"helpful resources\"", "\"recommended
	// resources\"", "\"helpful links\"", "\"resource list\"", "resources"].
	Footprints []string `json:"footprints,omitempty"`

	// Resource pages to check: How many resource pages are read and have their links checked, one page per site,
	// best placed first. From 1 to 30. Default: 10.
	PagesToRead *int `json:"pages_to_read,omitempty"`

	// Links looked at per page: The first this many links of each page are looked at, and those to other sites are
	// checked, one unit each. From 10 to 100. Default: 80.
	LinksPerPage *int `json:"links_per_page,omitempty"`

	// Results pages per search: How many pages of results to read for each footprint, about ten results a page.
	// From 1 to 3. Default: 1.
	SearchPages *int `json:"search_pages,omitempty"`

	// Sites to leave out: Domains whose pages are not read and whose links are not checked, such as social
	// networks and share buttons. At most 2000 items. Default: a list of 22 values, which `get` on the formula
	// shows.
	ExcludeDomains []string `json:"exclude_domains,omitempty"`

	// Search country: The two-letter country whose search results are read, such as us, gb or de. Default: "us".
	Country *string `json:"country,omitempty"`

	// Search language: The language of the results, such as en or de. Default: "en".
	Language *string `json:"language,omitempty"`

	// Only sites with an email: Keeps the pages whose site publishes an email address that passes the check.
	// Default: false.
	RequireEmail *bool `json:"require_email,omitempty"`

	// Show a broken link again after (days): For a saved formula, how long before a delivered broken link may come
	// back. From 1 to 365. Default: 90.
	ResurfaceAfterDays *int `json:"resurface_after_days,omitempty"`

	// Most rows: The run stops at this many broken links, gone sites first. From 1 to 2000. Default: 200.
	MaxResults *int `json:"max_results,omitempty"`
}

// FormulaID is "broken-link-prospects".
func (BrokenLinkProspectsParams) FormulaID() string { return "broken-link-prospects" }

// BrokerDirectoryParams are the parameters of Broker Directory (`broker-directory`): The business brokers of your
// states, one row for each firm, with specialties, designations and their current listings.
//
// States is required. Any other field left nil or empty is not sent, and the formula's default applies; [Ptr] makes
// the pointer of a value.
type BrokerDirectoryParams struct {
	// States: US states or Canadian provinces, by name or two-letter code. Listings are read for US states only. 1
	// to 10 items. Required.
	States []string `json:"states"`

	// Specialties: IBBA brokers who list one of these industries. Empty keeps every broker. Each one of
	// "accommodations", "agriculture", "arts", "auto-related-business", "business-services", "construction",
	// "ecommerce", "education-services", "entertainment-recreation", "finance-insurance", "food-beverage",
	// "health-care-social-assistance", "industrial-services", "information", "logistics", "manufacturing",
	// "print-and-media", "public-administration", "real-estate-rental-leasing", "retail-trade",
	// "technical-services", "transportation", "waste-management", "wholesale-trade". Default: [].
	Specialties []string `json:"specialties,omitempty"`

	// Designation: IBBA brokers holding this designation, such as CBI (Certified Business Intermediary) or M&AMI.
	// One of "any", "cbi", "mami", "master_cbi", "cbi_or_mami". Default: "any".
	Certification *string `json:"certification,omitempty"`

	// Add the firms selling businesses now: Reads the state's businesses for sale on BusinessesForSale.com, and
	// each brokered one's page for its firm. Default: true.
	IncludeListings *bool `json:"include_listings,omitempty"`

	// Listing pages read a state: BusinessesForSale.com's pages of 25, newest first. From 1 to 4. Default: 2.
	ListingPages *int `json:"listing_pages,omitempty"`

	// Fewest listings now: Keep firms with at least this many of the listings read. 0 keeps every broker. At least
	// 0. Default: 0.
	MinListings *int `json:"min_listings,omitempty"`

	// Only brokers who name a firm: Leaves out IBBA members whose profile names no firm. Default: false.
	FirmsOnly *bool `json:"firms_only,omitempty"`

	// Read each website for an address: The address each firm's own website publishes, checked. Off saves one read
	// a firm. Default: true.
	FindEmail *bool `json:"find_email,omitempty"`

	// Show a firm again after (days): For a saved formula, how long before a delivered firm may come back. From 1
	// to 365. Default: 180.
	ResurfaceAfterDays *int `json:"resurface_after_days,omitempty"`

	// Most rows: The run stops at this many firms, those with the most listings and brokers first. From 1 to 2000.
	// Default: 200.
	MaxResults *int `json:"max_results,omitempty"`
}

// FormulaID is "broker-directory".
func (BrokerDirectoryParams) FormulaID() string { return "broker-directory" }

// BusinessForSaleAlertsParams are the parameters of Business-for-Sale Alerts (`business-for-sale-alerts`):
// Businesses newly listed for sale in your states and sectors, with asking price, revenue, cash flow, sector and
// broker.
//
// States is required. Any other field left nil or empty is not sent, and the formula's default applies; [Ptr] makes
// the pointer of a value.
type BusinessForSaleAlertsParams struct {
	// States: US states, by name, such as Texas or North Carolina. 1 to 10 items. Required.
	States []string `json:"states"`

	// Kinds of business: Words a listing's title or summary must hold as whole words, such as laundromat or HVAC.
	// Empty keeps every kind. At most 10 items. Default: [].
	Keywords []string `json:"keywords,omitempty"`

	// Least cash flow ($): The yearly cash flow or net profit a listing must state. 0 keeps every listing. At
	// least 0. Default: 0.
	MinCashFlow *int `json:"min_cash_flow,omitempty"`

	// Highest asking price ($): Listings asking more are left out. 0 sets no limit. At least 0. Default: 0.
	MaxAskingPrice *int `json:"max_asking_price,omitempty"`

	// Lowest asking price ($): Listings asking less are left out. At least 0. Default: 0.
	MinAskingPrice *int `json:"min_asking_price,omitempty"`

	// Least revenue ($): The yearly revenue a listing must state. 0 keeps every listing. At least 0. Default: 0.
	MinRevenue *int `json:"min_revenue,omitempty"`

	// Keep listings that withhold a figure: A listing whose cash flow, revenue or price is on request is kept
	// rather than judged by a number it does not give. Default: true.
	KeepUndisclosed *bool `json:"keep_undisclosed,omitempty"`

	// Sold by: Every listing, only those a broker sells, or only those the owner sells. One of "any", "broker",
	// "owner". Default: "any".
	Seller *string `json:"seller,omitempty"`

	// Marketplaces: BusinessesForSale.com and BusinessBroker.net. Each one of "businessesforsale",
	// "businessbroker". At least 1 item. Default: ["businessesforsale", "businessbroker"].
	Sources []string `json:"sources,omitempty"`

	// Pages read a state: BusinessesForSale.com's pages of 25, newest first; BusinessBroker.net's first page of
	// about 30 is read either way. From 1 to 4. Default: 2.
	Pages *int `json:"pages,omitempty"`

	// Larger companies without their name: BusinessesForSale.com's M&A Vault describes larger companies without
	// naming them, from any country, among a state's listings. Default: false.
	IncludeVault *bool `json:"include_vault,omitempty"`

	// Read each listing's page: For the sector, employees, year founded, premises, reason for selling and the
	// broker. Off saves one read a listing. Default: true.
	ReadListings *bool `json:"read_listings,omitempty"`

	// Show a listing again after (days): For a saved formula, how long before a delivered listing may come back.
	// From 1 to 365. Default: 365.
	ResurfaceAfterDays *int `json:"resurface_after_days,omitempty"`

	// Most rows: The run stops at this many listings, highest cash flow first. From 1 to 1000. Default: 100.
	MaxResults *int `json:"max_results,omitempty"`
}

// FormulaID is "business-for-sale-alerts".
func (BusinessForSaleAlertsParams) FormulaID() string { return "business-for-sale-alerts" }

// CandidateReverseMarketingParams are the parameters of Candidate Reverse Marketing
// (`candidate-reverse-marketing`): Companies hiring for what your best candidate does, with the skills they ask
// for, who decides, and an opening line.
//
// CandidateTitles and Skills and Locations are required. Any other field left nil or empty is not sent, and the
// formula's default applies; [Ptr] makes the pointer of a value.
type CandidateReverseMarketingParams struct {
	// Roles the candidate fits: The job titles your candidate would be hired as, each searched as you type it. 1
	// to 5 items. Required.
	CandidateTitles []string `json:"candidate_titles"`

	// The candidate's skills: Tools, methods or certificates the candidate has. A posting is kept when it asks for
	// enough of them. 1 to 15 items. Required.
	Skills []string `json:"skills"`

	// Where they will work: A city, region or country for each place the candidate would work. 1 to 5 items.
	// Required.
	Locations []string `json:"locations"`

	// Posted within (days): Only postings from the last this many days; a fresh posting is still unfilled. From 1
	// to 30. Default: 7.
	PostedWithinDays *int `json:"posted_within_days,omitempty"`

	// Job boards: Which boards to search. Reed and Totaljobs cover the UK; Dice is technology roles in the US.
	// Each one of "linkedin", "indeed", "reed", "totaljobs", "dice". At least 1 item. Default: ["linkedin",
	// "indeed"].
	Sources []string `json:"sources,omitempty"`

	// Jobs read per search: How many postings to read for each role, place and board. From 5 to 100. Default: 25.
	ResultsPerSearch *int `json:"results_per_search,omitempty"`

	// Least skills a posting must ask for: Of the candidate's skills, how many a posting has to name to be kept.
	// From 1 to 10. Default: 2.
	MinSkillsMatched *int `json:"min_skills_matched,omitempty"`

	// Postings read in full: The newest postings whose full text is read to match skills; each costs a unit when
	// the board gave only a summary. From 10 to 300. Default: 80.
	PostingsRead *int `json:"postings_read,omitempty"`

	// Company size (employees): Companies whose size is not known are kept. Default: {"min": 11, "max": 10000}.
	CompanySize *IntRange `json:"company_size,omitempty"`

	// Companies to leave out: Names or domains, such as the candidate's current employer or companies they already
	// applied to. At most 5000 items. Default: [].
	ExcludeCompanies []string `json:"exclude_companies,omitempty"`

	// Words to leave out of titles: For example junior or intern. Default: [].
	TitleExcludeKeywords []string `json:"title_exclude_keywords,omitempty"`

	// Who to contact: Leave empty to look for whoever decides on that kind of hire, such as the VP Engineering for
	// a developer. Default: [].
	ContactTitles []string `json:"contact_titles,omitempty"`

	// Contacts per company: 0 turns contacts off. From 0 to 5. Default: 1.
	ContactsPerCompany *int `json:"contacts_per_company,omitempty"`

	// Only addresses that pass the check: Keeps addresses whose domain takes mail, and those a verification
	// provider confirmed. Default: true.
	RequireVerifiedEmail *bool `json:"require_verified_email,omitempty"`

	// Most rows: The run stops at this many rows, best first. From 1 to 1000. Default: 50.
	MaxResults *int `json:"max_results,omitempty"`
}

// FormulaID is "candidate-reverse-marketing".
func (CandidateReverseMarketingParams) FormulaID() string { return "candidate-reverse-marketing" }

// ChampionTrackerParams are the parameters of Champion Tracker (`champion-tracker`): People who held your
// champions' roles at your accounts and moved on, found at their new company, with contact details.
//
// Accounts and ChampionTitles are required. Any other field left nil or empty is not sent, and the formula's
// default applies; [Ptr] makes the pointer of a value.
type ChampionTrackerParams struct {
	// Accounts: The website domains of the accounts where your champions work or worked, such as your customers. 1
	// to 100 items. Required.
	Accounts []string `json:"accounts"`

	// Champion roles: The roles your champions held, such as VP Sales. 1 to 10 items. Required.
	ChampionTitles []string `json:"champion_titles"`

	// Who to find: moved: people who held these roles at the account and left, where they are now. successor: who
	// holds them at the account today. One of "moved", "successor". Default: "moved".
	Track *string `json:"track,omitempty"`

	// Least employees at the new company: Only follow people to companies at least this big. 1 follows everyone.
	// At least 1. Default: 1.
	MinEmployees *int `json:"min_employees,omitempty"`

	// New companies to leave out: Domains of companies not to follow people to, such as competitors or customers
	// already won. At most 5000 items. Default: [].
	ExcludeCompanies []string `json:"exclude_companies,omitempty"`

	// People per account: The most people to find for each account. From 1 to 10. Default: 3.
	PeoplePerAccount *int `json:"people_per_account,omitempty"`

	// Only addresses that pass the check: Keeps addresses whose domain takes mail, and those a verification
	// provider confirmed. Default: true.
	RequireVerifiedEmail *bool `json:"require_verified_email,omitempty"`

	// Show a person again after (days): For a saved formula, how long before a person already delivered may come
	// back. From 1 to 365. Default: 180.
	ResurfaceAfterDays *int `json:"resurface_after_days,omitempty"`

	// Most rows: The run stops at this many people, best first. From 1 to 2000. Default: 200.
	MaxResults *int `json:"max_results,omitempty"`
}

// FormulaID is "champion-tracker".
func (ChampionTrackerParams) FormulaID() string { return "champion-tracker" }

// CompanyDiligenceParams are the parameters of Company Diligence (`company-diligence`): Each company you name
// checked on its register, in adverse news, reviews and complaints, and its website, with red flags.
//
// Companies is required. Any other field left nil or empty is not sent, and the formula's default applies; [Ptr]
// makes the pointer of a value.
type CompanyDiligenceParams struct {
	// Companies: Each company by its name, its register number (a UK company number such as 09446231, a Swedish
	// organisation number such as 556703-7485) or its website's domain. A name is searched as written, in the news
	// too, so give the name the company is known by (Thames Water, not Thames Water Utilities Limited); a domain
	// also lets its reviews and website be read. 1 to 200 items. Required.
	Companies []string `json:"companies"`

	// Country: Where the companies are registered. The UK and Sweden have their register read; the US and Canada
	// have the Better Business Bureau read. The news, reviews and website checks run everywhere. One of "GB",
	// "US", "CA", "SE", "IE", "AU", "other". Default: "GB".
	Country *string `json:"country,omitempty"`

	// Place, for the Better Business Bureau: The city or state the companies trade in (Austin, TX). Only used in
	// the US and Canada, to find the right business. Default: "".
	Location *string `json:"location,omitempty"`

	// Checks: Which checks to run. Each costs about one unit a company; finding a company's website from its name
	// costs a few more. Each one of "register", "news", "press_releases", "reviews", "website". At least 1 item.
	// Default: ["register", "news", "press_releases", "reviews", "website"].
	Checks []string `json:"checks,omitempty"`

	// Adverse words: The news is searched for the company together with these words, and a headline that holds one
	// is counted as adverse. 1 to 30 items. Default: a list of 18 values, which `get` on the formula shows.
	AdverseWords []string `json:"adverse_words,omitempty"`

	// What the check is for: Counterparty weighs the register and adverse news most. Acquisition also weighs
	// customer reviews and complaints, which say most about a small business a buyer is about to take on. One of
	// "counterparty", "acquisition". Default: "counterparty".
	Focus *string `json:"focus,omitempty"`

	// News within (days): How far back the news and the press releases are read. From 30 to 1825. Default: 365.
	WindowDays *int `json:"window_days,omitempty"`

	// Most rows: The run stops at this many companies, riskiest first. From 1 to 2000. Default: 200.
	MaxResults *int `json:"max_results,omitempty"`
}

// FormulaID is "company-diligence".
func (CompanyDiligenceParams) FormulaID() string { return "company-diligence" }

// CompetitorPriceTrackerParams are the parameters of Competitor Price Tracker (`competitor-price-tracker`): Every
// product and price of competitors' Shopify stores and Amazon listings, read on a schedule, with what changed since
// last time.
//
// Stores is required. Any other field left nil or empty is not sent, and the formula's default applies; [Ptr] makes
// the pointer of a value.
type CompetitorPriceTrackerParams struct {
	// Competitor stores: Shopify stores by domain, such as allbirds.com. Each store's whole catalogue is read. At
	// most 20 items. Required.
	Stores []string `json:"stores"`

	// Amazon listings: Competitors' Amazon products by ASIN, read one by one. At most 100 items. Default: [].
	Asins []string `json:"asins,omitempty"`

	// Only products about: Keep products whose title or type holds one of these. Empty keeps every product. At
	// most 20 items. Default: [].
	Keywords []string `json:"keywords,omitempty"`

	// Store collection: Read only this collection of each store, by its handle (the name in
	// /collections/<handle>). Empty reads the whole catalogue. Default: "".
	Collection *string `json:"collection,omitempty"`

	// Products read per store: A store's catalogue is read 250 products at a time; this many are read. From 10 to
	// 250. Default: 250.
	ProductsPerStore *int `json:"products_per_store,omitempty"`

	// Amazon marketplace: Which Amazon the listings are on. One of "us", "uk", "de", "fr", "it", "es", "ca".
	// Default: "us".
	AmazonMarketplace *string `json:"amazon_marketplace,omitempty"`

	// Only what changed: Show only products that are new or whose price, compare-at price or stock changed since
	// the last run. Default: false.
	OnlyChanges *bool `json:"only_changes,omitempty"`

	// Most rows: The run stops at this many products, changes first. From 1 to 5000. Default: 1000.
	MaxResults *int `json:"max_results,omitempty"`
}

// FormulaID is "competitor-price-tracker".
func (CompetitorPriceTrackerParams) FormulaID() string { return "competitor-price-tracker" }

// CompetitorReviewMiningParams are the parameters of Competitor Review Mining (`competitor-review-mining`): Your
// competitors' newest one- and two-star reviews, each with what went wrong, to sell against and to win their
// customers.
//
// Competitors is required. Any other field left nil or empty is not sent, and the formula's default applies; [Ptr]
// makes the pointer of a value.
type CompetitorReviewMiningParams struct {
	// Competitors: Each competitor's website domain, as Trustpilot lists it, such as gusto.com. 1 to 20 items.
	// Required.
	Competitors []string `json:"competitors"`

	// Stars: Which reviews to read, by their stars. At least 1 item. Default: [1, 2].
	Stars []int `json:"stars,omitempty"`

	// Phrases to look for: Words you can answer, such as "customer service" or "hidden fees"; each match is
	// quoted. At most 30 items. Default: [].
	Phrases []string `json:"phrases,omitempty"`

	// Written within (days): Only reviews written in the last this many days. From 1 to 730. Default: 180.
	PublishedWithinDays *int `json:"published_within_days,omitempty"`

	// Pages per competitor: Twenty reviews a page, newest first. From 1 to 3. Default: 2.
	PagesPerCompetitor *int `json:"pages_per_competitor,omitempty"`

	// Only reviews with a phrase: Keeps only the reviews that contain one of your phrases. Default: false.
	OnlyMatching *bool `json:"only_matching,omitempty"`

	// Show a review again after (days): For a saved formula, a review once delivered is not delivered again before
	// this. From 1 to 365. Default: 365.
	ResurfaceAfterDays *int `json:"resurface_after_days,omitempty"`

	// Most rows: The run stops at this many reviews, the best matches first. From 1 to 2000. Default: 200.
	MaxResults *int `json:"max_results,omitempty"`
}

// FormulaID is "competitor-review-mining".
func (CompetitorReviewMiningParams) FormulaID() string { return "competitor-review-mining" }

// CompoundSignalScoreParams are the parameters of Compound Signal Score (`compound-signal-score`): Your target
// accounts scored by what they are doing now: hiring, raising, new leaders, expanding and advertising.
//
// Accounts is required. Any other field left nil or empty is not sent, and the formula's default applies; [Ptr]
// makes the pointer of a value.
type CompoundSignalScoreParams struct {
	// Accounts: The companies to score, by name, by domain, or both as Rippling (rippling.com); a domain given
	// needs no search to match. 1 to 50 items. Required.
	Accounts []string `json:"accounts"`

	// Look back (days): How recent news, postings and ads must be to count. From 7 to 180. Default: 30.
	WindowDays *int `json:"window_days,omitempty"`

	// Signals to read: What to read for each account. Each costs units for every account. Each one of "hiring",
	// "announcements", "ads". At least 1 item. Default: ["hiring", "announcements", "ads"].
	Signals []string `json:"signals,omitempty"`

	// Weight of hiring: Points for open roles on its careers page, in full at ten roles or more. From 0 to 100.
	// Default: 30.
	WeightHiring *int `json:"weight_hiring,omitempty"`

	// Weight of funding: Points for a funding announcement in the window. From 0 to 100. Default: 25.
	WeightFunding *int `json:"weight_funding,omitempty"`

	// Weight of a senior hire: Points for an announced executive hire in the window. From 0 to 100. Default: 15.
	WeightExecutive *int `json:"weight_executive,omitempty"`

	// Weight of an expansion: Points for an announced expansion or new office in the window. From 0 to 100.
	// Default: 10.
	WeightExpansion *int `json:"weight_expansion,omitempty"`

	// Weight of a launch: Points for an announced product launch in the window. From 0 to 100. Default: 5.
	WeightLaunch *int `json:"weight_launch,omitempty"`

	// Weight of advertising: Points for ads shown in the window. From 0 to 100. Default: 15.
	WeightAds *int `json:"weight_ads,omitempty"`

	// Most open roles: Leave out accounts with more open roles than this, to find those that will hire before they
	// post. At least 0. Default: 100000.
	MaxOpenRoles *int `json:"max_open_roles,omitempty"`

	// Ads read per account: How many of an account's newest ads are read; each costs a unit. From 1 to 40.
	// Default: 10.
	AdsRead *int `json:"ads_read,omitempty"`

	// Least score: Accounts scoring below this are left out. At least 0. Default: 0.
	MinScore *float64 `json:"min_score,omitempty"`

	// Most rows: The run stops at this many accounts, best first. From 1 to 50. Default: 50.
	MaxResults *int `json:"max_results,omitempty"`
}

// FormulaID is "compound-signal-score".
func (CompoundSignalScoreParams) FormulaID() string { return "compound-signal-score" }

// ContractorLicenseWatchParams are the parameters of Contractor Licence Watch (`contractor-license-watch`):
// Contractor businesses whose licence, bond or insurance expires in the next weeks, from official state and city
// registers.
//
// Register is required. Any other field left nil or empty is not sent, and the formula's default applies; [Ptr]
// makes the pointer of a value.
type ContractorLicenseWatchParams struct {
	// Register: Washington L&I, Texas TDLR, Oregon CCB (with bond and insurance dates), Chicago home repair,
	// Buffalo, or New York's mold and elevator contractors. One of "washington", "texas", "oregon", "chicago",
	// "buffalo", "new_york_mold", "new_york_elevator". Required.
	Register string `json:"register"`

	// Expiring within (days): Licences that run out between today and this many days from now. From 1 to 180.
	// Default: 30.
	ExpiringWithinDays *int `json:"expiring_within_days,omitempty"`

	// Trades: Words in the licence's type or trade, such as electrical, plumbing or roofing. Each is one search.
	// Empty reads every trade. At most 8 items. Default: [].
	Trades []string `json:"trades,omitempty"`

	// Only licences in good standing: Leave out licences the register shows as suspended, revoked, cancelled or
	// expired. Default: true.
	ActiveOnly *bool `json:"active_only,omitempty"`

	// Licences read per search: The register is read up to this many licences for each trade word. From 10 to 200.
	// Default: 200.
	LicencesToRead *int `json:"licences_to_read,omitempty"`

	// Show a licence again after (days): For a saved formula, how long before the same licence and expiry date may
	// come back. From 1 to 365. Default: 300.
	ResurfaceAfterDays *int `json:"resurface_after_days,omitempty"`

	// Most rows: The run stops at this many licences, soonest expiry first. From 1 to 2000. Default: 500.
	MaxResults *int `json:"max_results,omitempty"`
}

// FormulaID is "contractor-license-watch".
func (ContractorLicenseWatchParams) FormulaID() string { return "contractor-license-watch" }

// CrowdfundingLeadsParams are the parameters of Crowdfunding Leads (`crowdfunding-leads`): Kickstarter and
// Indiegogo projects with traction, each with funding, backers, the creator's website and the email it publishes.
//
// Keywords is required. Any other field left nil or empty is not sent, and the formula's default applies; [Ptr]
// makes the pointer of a value.
type CrowdfundingLeadsParams struct {
	// Keywords: What the projects make or are about, each searched on each platform, such as e-bike or tabletop
	// game. 1 to 10 items. Required.
	Keywords []string `json:"keywords"`

	// Campaign stage: Live campaigns, campaigns that ended funded (the ones that now have to make and ship), or
	// both. One of "live", "funded", "any". Default: "live".
	Stage *string `json:"stage,omitempty"`

	// Platforms: Where to search. Each one of "kickstarter", "indiegogo". At least 1 item. Default:
	// ["kickstarter", "indiegogo"].
	Platforms []string `json:"platforms,omitempty"`

	// Fewest backers: Projects with fewer backers are left out. At least 0. Default: 10.
	MinBackers *int `json:"min_backers,omitempty"`

	// Least funded (% of goal): 100 keeps the projects that reached their goal. A project that states no goal is
	// kept. From 0 to 100000. Default: 0.
	MinPercentFunded *int `json:"min_percent_funded,omitempty"`

	// Ended within (days): For funded projects, how recently the campaign ended. Live campaigns always pass. From
	// 1 to 3650. Default: 120.
	EndedWithinDays *int `json:"ended_within_days,omitempty"`

	// Only these categories: The platform's categories or tags, or part of them, such as 3D Printing, Product
	// Design or Home. Leave empty for any. At most 30 items. Default: [].
	Categories []string `json:"categories,omitempty"`

	// Only these countries (Kickstarter): Two-letter codes such as US or GB. Kickstarter states a project's
	// country; Indiegogo does not, so its projects are kept. At most 50 items. Default: [].
	Countries []string `json:"countries,omitempty"`

	// Pages read per keyword: Each page holds about 12 Kickstarter or 24 Indiegogo projects. From 1 to 5. Default:
	// 2.
	Pages *int `json:"pages,omitempty"`

	// Only creators with an email: Keeps the projects whose creator's own website publishes an email address that
	// passes the check. Default: true.
	RequireEmail *bool `json:"require_email,omitempty"`

	// Most projects: How many projects, best backed first, are matched to a website; with Only creators with an
	// email, fewer rows may remain. From 1 to 500. Default: 100.
	MaxResults *int `json:"max_results,omitempty"`
}

// FormulaID is "crowdfunding-leads".
func (CrowdfundingLeadsParams) FormulaID() string { return "crowdfunding-leads" }

// DealSiteMonitorParams are the parameters of Deal-Site Monitor (`deal-site-monitor`): New Slickdeals posts for
// your words, matched to Amazon, with the deal's price against Amazon's and the resale margin.
//
// Queries is required. Any other field left nil or empty is not sent, and the formula's default applies; [Ptr]
// makes the pointer of a value.
type DealSiteMonitorParams struct {
	// Words to watch: Products, brands or stores, such as lego, airpods or costco. Each is one search, newest
	// deals first. 1 to 10 items. Required.
	Queries []string `json:"queries"`

	// Only deals at: Store names as Slickdeals shows them, such as Walmart or Target. Empty keeps every store. At
	// most 30 items. Default: [].
	Stores []string `json:"stores,omitempty"`

	// Least deal score: The community's votes on the deal. From -50 to 1000. Default: 0.
	MinScore *int `json:"min_score,omitempty"`

	// Compare with Amazon: Read each deal's page and find the product on Amazon. Two reads a deal. Default: true.
	CompareOnAmazon *bool `json:"compare_on_amazon,omitempty"`

	// Amazon referral fee (%). From 0 to 50. Default: 15.
	ReferralFeePct *float64 `json:"referral_fee_pct,omitempty"`

	// Fulfilment and shipping a unit (USD). From 0 to 100. Default: 5.
	FulfilmentCost *float64 `json:"fulfilment_cost,omitempty"`

	// Show a deal again after (days). From 1 to 365. Default: 365.
	ResurfaceAfterDays *int `json:"resurface_after_days,omitempty"`

	// Most rows: The run stops at this many deals, newest first. Each is two reads when compared with Amazon. From
	// 1 to 500. Default: 40.
	MaxResults *int `json:"max_results,omitempty"`
}

// FormulaID is "deal-site-monitor".
func (DealSiteMonitorParams) FormulaID() string { return "deal-site-monitor" }

// DirectoryToLeadsParams are the parameters of Directory to Leads (`directory-to-leads`): The businesses a
// directory lists, each with its website and the email that site publishes, checked.
//
// Directory and What are required. Any other field left nil or empty is not sent, and the formula's default
// applies; [Ptr] makes the pointer of a value.
type DirectoryToLeadsParams struct {
	// Directory: Which directory to read; a web page is any directory page of your own, such as a chamber of
	// commerce's members. One of "bbb", "gelbeseiten", "checkatrade", "clutch", "web_page". Required.
	Directory string `json:"directory"`

	// Kind of business: As the directory names it: plumber on BBB or Checkatrade, Zahnarzt on Gelbe Seiten, a
	// Clutch directory with its place such as web-developers/london, or the address of the directory page.
	// Required.
	What string `json:"what"`

	// Place: A town, city or postcode. Not needed for Clutch, whose place is in its directory, or for a web page.
	// Default: "".
	Where *string `json:"where,omitempty"`

	// Pages to read: How many pages of the directory to read; for a web page, how many of its pages to follow.
	// From 1 to 10. Default: 3.
	Pages *int `json:"pages,omitempty"`

	// Only businesses with an email: Keeps the businesses whose own website publishes an email address that passes
	// the check. Default: true.
	RequireEmail *bool `json:"require_email,omitempty"`

	// Only named addresses: Leaves out shared inboxes such as info@. Default: false.
	ExcludeRoleAddresses *bool `json:"exclude_role_addresses,omitempty"`

	// Fewest reviews: Businesses with fewer reviews on the directory are left out. At least 0. Default: 0.
	MinReviews *int `json:"min_reviews,omitempty"`

	// Businesses to leave out: Names or domains, such as your existing customers. At most 5000 items. Default: [].
	ExcludeCompanies []string `json:"exclude_companies,omitempty"`

	// Show a business again after (days): For a saved formula, how long before a delivered business may come back.
	// From 1 to 365. Default: 180.
	ResurfaceAfterDays *int `json:"resurface_after_days,omitempty"`

	// Most rows: The run stops at this many businesses. From 1 to 5000. Default: 200.
	MaxResults *int `json:"max_results,omitempty"`
}

// FormulaID is "directory-to-leads".
func (DirectoryToLeadsParams) FormulaID() string { return "directory-to-leads" }

// EcommercePanelParams are the parameters of E-commerce Panel (`ecommerce-panel`): A brand's listings on Amazon,
// Walmart and Target, with price, stock, rating and ratings count, and what changed since the last run.
//
// Queries is required. Any other field left nil or empty is not sent, and the formula's default applies; [Ptr]
// makes the pointer of a value.
type EcommercePanelParams struct {
	// Brands or products: Each is searched on every retailer you choose, as you type it, such as a brand or a
	// brand and a product line. 1 to 20 items. Required.
	Queries []string `json:"queries"`

	// Only these brands: Keeps the listings whose brand or title names one of these, so that other brands a search
	// shows are left out. Leave empty to keep every listing. At most 20 items. Default: [].
	Brands []string `json:"brands,omitempty"`

	// Retailers: Where to search. Walmart and Target are their US stores. Each one of "amazon", "walmart",
	// "target". At least 1 item. Default: ["amazon", "walmart", "target"].
	Retailers []string `json:"retailers,omitempty"`

	// Listings per search: How many listings to read from the top of each retailer's results for each brand or
	// product. From 1 to 40. Default: 20.
	ResultsPerSearch *int `json:"results_per_search,omitempty"`

	// Include sponsored listings: Paid placements are left out by default, so that a position is the retailer's
	// own ranking. Default: false.
	IncludeSponsored *bool `json:"include_sponsored,omitempty"`

	// Amazon marketplace: Which Amazon store to read. One of "us", "uk", "de", "ca", "it", "es", "nl", "jp", "au",
	// "mx", "br", "in", "ae". Default: "us".
	AmazonMarketplace *string `json:"amazon_marketplace,omitempty"`

	// Amazon Best Sellers lists: Categories as their Best Sellers address names them, such as kitchen, beauty or
	// electronics/172456. A listing in their top 100 gets its rank, and a best seller of your brands is added. At
	// most 10 items. Default: [].
	AmazonBestsellerCategories []string `json:"amazon_bestseller_categories,omitempty"`

	// Check stock on Amazon and Target: Their search results do not say whether a product is in stock, so each
	// Amazon and Target listing's page is read, one unit a listing. Default: true.
	CheckStock *bool `json:"check_stock,omitempty"`

	// Most rows: The run stops at this many listings, the biggest changes first. From 1 to 2000. Default: 200.
	MaxResults *int `json:"max_results,omitempty"`
}

// FormulaID is "ecommerce-panel".
func (EcommercePanelParams) FormulaID() string { return "ecommerce-panel" }

// EnterpriseAtsSweepParams are the parameters of Enterprise ATS Sweep (`enterprise-ats-sweep`): The new openings at
// your target companies, read from their own careers pages, with the companies hiring hardest flagged P1.
//
// TargetCompanies is required. Any other field left nil or empty is not sent, and the formula's default applies;
// [Ptr] makes the pointer of a value.
type EnterpriseAtsSweepParams struct {
	// Target companies: A domain (acme.com) or a careers page address for each company. 1 to 100 items. Required.
	TargetCompanies []string `json:"target_companies"`

	// Job keywords: Openings whose title holds one of these as a whole word or phrase. Leave empty for every
	// opening. At most 20 items. Default: [].
	JobKeywords []string `json:"job_keywords,omitempty"`

	// Locations: Openings in one of these places; remote openings are kept too. Leave empty for anywhere. At most
	// 20 items. Default: [].
	Locations []string `json:"locations,omitempty"`

	// Posted within (days): Only openings posted in the last this many days. An opening with no date is kept. From
	// 1 to 365. Default: 30.
	PostedWithinDays *int `json:"posted_within_days,omitempty"`

	// Words to leave out of titles: For example intern or graduate. Default: [].
	ExcludeKeywords []string `json:"exclude_keywords,omitempty"`

	// Only openings from the tracking system: Leaves out what was read from a careers page's own text, where a
	// heading or a name can pass for an opening. Off reads those too. Default: true.
	VerifiedOnly *bool `json:"verified_only,omitempty"`

	// Least matching openings: A company with fewer matching openings is left out. At least 1. Default: 1.
	MinMatchingRoles *int `json:"min_matching_roles,omitempty"`

	// P1 at this many recent openings: A company with at least this many matching openings posted in the P1 window
	// is flagged P1. At least 1. Default: 5.
	P1MinRoles *int `json:"p1_min_roles,omitempty"`

	// P1 window (days): How recent an opening must be to count towards P1. From 1 to 90. Default: 7.
	P1WithinDays *int `json:"p1_within_days,omitempty"`

	// Show an opening again after (days): For a saved formula, how long before a delivered opening may come back.
	// From 1 to 365. Default: 30.
	ResurfaceAfterDays *int `json:"resurface_after_days,omitempty"`

	// Most rows: The run stops at this many rows, best first. From 1 to 5000. Default: 500.
	MaxResults *int `json:"max_results,omitempty"`
}

// FormulaID is "enterprise-ats-sweep".
func (EnterpriseAtsSweepParams) FormulaID() string { return "enterprise-ats-sweep" }

// EventOrganizerLeadsParams are the parameters of Event Organizer Leads (`event-organizer-leads`): The groups and
// organisers running events on your topics in a city, from Meetup and Eventbrite, busiest first.
//
// Topics and City and Region are required. Any other field left nil or empty is not sent, and the formula's default
// applies; [Ptr] makes the pointer of a value.
type EventOrganizerLeadsParams struct {
	// Topics: What the events are about, as you would search for them, such as startup, python or yoga. 1 to 8
	// items. Required.
	Topics []string `json:"topics"`

	// City: One city, by its name alone. Required.
	City string `json:"city"`

	// State or province: For a city in the United States or Canada, its state or province in two letters, which
	// both sites need there; leave it empty elsewhere. Required.
	Region string `json:"region"`

	// Country: The country, in two letters, such as US, GB or DE. Default: "US".
	Country *string `json:"country,omitempty"`

	// Where to look: Meetup's groups and Eventbrite's organisers. Each one of "meetup", "eventbrite". At least 1
	// item. Default: ["meetup", "eventbrite"].
	Sources []string `json:"sources,omitempty"`

	// Eventbrite's name for the place: Outside the United States and Canada, the place as Eventbrite's addresses
	// write it, such as united-kingdom--london. Empty there leaves Eventbrite out. Default: "".
	EventbritePlace *string `json:"eventbrite_place,omitempty"`

	// Events within (days): Only events starting in the next this many days count. From 1 to 365. Default: 90.
	WithinDays *int `json:"within_days,omitempty"`

	// Fewest group members: Meetup groups with fewer members are left out. Eventbrite does not count followers, so
	// its organisers are kept. At least 0. Default: 0.
	MinMembers *int `json:"min_members,omitempty"`

	// Fewest events found: Organisers with fewer of the events found are left out. At least 1. Default: 1.
	MinEvents *int `json:"min_events,omitempty"`

	// Show an organiser again after (days): For a saved formula, how long before a delivered organiser may come
	// back. From 1 to 365. Default: 90.
	ResurfaceAfterDays *int `json:"resurface_after_days,omitempty"`

	// Most rows: The run stops at this many organisers, best first. From 1 to 1000. Default: 200.
	MaxResults *int `json:"max_results,omitempty"`
}

// FormulaID is "event-organizer-leads".
func (EventOrganizerLeadsParams) FormulaID() string { return "event-organizer-leads" }

// ExhibitorProspectsParams are the parameters of Exhibitor Prospects (`exhibitor-prospects`): The companies
// exhibiting at the trade shows you name, with booth, website, location, categories and contact address.
//
// EventURLs is required. Any other field left nil or empty is not sent, and the formula's default applies; [Ptr]
// makes the pointer of a value.
type ExhibitorProspectsParams struct {
	// Shows: The address of each show's exhibitor list or floor plan on its a2z site (a2zinc.net), as your browser
	// shows it. 1 to 10 items. Required.
	EventURLs []string `json:"event_urls"`

	// Categories or words: Keep exhibitors whose categories, products, brands or description name one of these, as
	// whole words. Empty keeps every exhibitor. At most 20 items. Default: [].
	Keywords []string `json:"keywords,omitempty"`

	// Your show: The a2z address of your own show. Its exhibitors are left out, to find the companies that exhibit
	// elsewhere but not with you. Default: "".
	YourEventURL *string `json:"your_event_url,omitempty"`

	// Countries: Keep exhibitors from these countries, in words as the shows write them, such as United States or
	// Canada. Empty keeps all. At most 50 items. Default: [].
	Countries []string `json:"countries,omitempty"`

	// Only upgraded listings: Keep exhibitors who paid the show for an enhanced listing, a sign of a bigger
	// exhibiting budget. Default: false.
	EnhancedOnly *bool `json:"enhanced_only,omitempty"`

	// Fewest shows: With several shows, keep companies that exhibit at least at this many of them. From 1 to 10.
	// Default: 1.
	MinShows *int `json:"min_shows,omitempty"`

	// Your accounts: Company names or domains from your CRM, matched against each exhibitor's name and website. At
	// most 5000 items. Default: [].
	CrmCompanies []string `json:"crm_companies,omitempty"`

	// Your accounts are: Use your accounts to keep only those exhibiting, or to leave them out and see only new
	// companies. One of "ignore", "only_accounts", "leave_out_accounts". Default: "ignore".
	CrmMatch *string `json:"crm_match,omitempty"`

	// Read each website for an address: The address the company's own website publishes, checked. Off saves one
	// read a company. Default: true.
	IncludeSiteContact *bool `json:"include_site_contact,omitempty"`

	// Show a company again after (days): For a saved formula, how long before a delivered company may come back.
	// From 1 to 365. Default: 120.
	ResurfaceAfterDays *int `json:"resurface_after_days,omitempty"`

	// Most exhibitors read: The exhibitors whose show pages are read, most shows and upgraded listings first;
	// words and countries then keep those that match. A show page and a website are one read each. From 1 to 1000.
	// Default: 100.
	MaxResults *int `json:"max_results,omitempty"`
}

// FormulaID is "exhibitor-prospects".
func (ExhibitorProspectsParams) FormulaID() string { return "exhibitor-prospects" }

// FootprintProspectsParams are the parameters of Footprint Prospects (`footprint-prospects`): Sites on your topic
// that take guest posts or keep resource pages, found by search footprints, each with its email.
//
// Topic is required. Any other field left nil or empty is not sent, and the formula's default applies; [Ptr] makes
// the pointer of a value.
type FootprintProspectsParams struct {
	// Topic: The subject, as people search for it, such as gardening or personal finance. Required.
	Topic string `json:"topic"`

	// Footprints: Words searched beside the topic, each its own search. Quote a phrase to find it exactly. Search
	// operators such as intitle: work too. 1 to 20 items. Default: ["\"write for us\"", "\"guest post
	// guidelines\"", "\"submit a guest post\"", "\"become a contributor\"", "\"accepting guest posts\"",
	// "\"contribute an article\""].
	Footprints []string `json:"footprints,omitempty"`

	// Results pages per search: How many pages of results to read for each footprint, about ten results a page.
	// From 1 to 3. Default: 1.
	SearchPages *int `json:"search_pages,omitempty"`

	// Sites to leave out: Domains never shown, such as social networks, marketplaces and your own sites. At most
	// 2000 items. Default: a list of 21 values, which `get` on the formula shows.
	ExcludeDomains []string `json:"exclude_domains,omitempty"`

	// Search country: The two-letter country whose search results are read, such as us, gb or de. Default: "us".
	Country *string `json:"country,omitempty"`

	// Search language: The language of the results, such as en or de. Default: "en".
	Language *string `json:"language,omitempty"`

	// Read each page for its author: Reads each page found as an article, for its author and the day it was
	// published. Default: false.
	ReadArticles *bool `json:"read_articles,omitempty"`

	// Only pages with a named author: When pages are read, keeps the sites whose page names who wrote it. Default:
	// false.
	RequireAuthor *bool `json:"require_author,omitempty"`

	// Only sites with an email: Keeps the sites whose own website publishes an email address that passes the
	// check. Default: false.
	RequireEmail *bool `json:"require_email,omitempty"`

	// Show a site again after (days): For a saved formula, how long before a delivered site may come back. From 1
	// to 365. Default: 90.
	ResurfaceAfterDays *int `json:"resurface_after_days,omitempty"`

	// Most rows: The run stops at this many sites, those matching the most footprints first. From 1 to 1000.
	// Default: 100.
	MaxResults *int `json:"max_results,omitempty"`
}

// FormulaID is "footprint-prospects".
func (FootprintProspectsParams) FormulaID() string { return "footprint-prospects" }

// FreelanceGigFeedParams are the parameters of Freelance Gig Feed (`freelance-gig-feed`): New projects on Upwork
// and Freelancer that match your skills, with budget and applicants, only the ones you have not seen.
//
// Keywords is required. Any other field left nil or empty is not sent, and the formula's default applies; [Ptr]
// makes the pointer of a value.
type FreelanceGigFeedParams struct {
	// Keywords: The skills or kinds of project to look for. Each is searched, and a project must mention one of
	// them. 1 to 10 items. Required.
	Keywords []string `json:"keywords"`

	// Posted within (days): Only projects posted in the last this many days. From 1 to 30. Default: 2.
	PostedWithinDays *int `json:"posted_within_days,omitempty"`

	// Least fixed budget: A fixed-price project must offer at least this much, in its own currency. A project with
	// no budget shown is kept. At least 0. Default: 0.
	MinFixedBudget *float64 `json:"min_fixed_budget,omitempty"`

	// Least hourly rate: An hourly project must offer up to at least this much an hour, in its own currency. One
	// with no rate shown is kept. At least 0. Default: 0.
	MinHourlyRate *float64 `json:"min_hourly_rate,omitempty"`

	// Marketplaces: Where to look. Each one of "upwork", "freelancer". At least 1 item. Default: ["upwork",
	// "freelancer"].
	Platforms []string `json:"platforms,omitempty"`

	// How it pays: Fixed price, hourly, or leave empty for both. Each one of "fixed", "hourly". Default: [].
	PaymentTypes []string `json:"payment_types,omitempty"`

	// Words to leave out: A project whose title or description holds one of these is left out, such as wordpress
	// or unpaid. Default: [].
	ExcludeKeywords []string `json:"exclude_keywords,omitempty"`

	// Most bids so far: Leave out a project with more bids than this, where the count is known. 0 for any. At
	// least 0. Default: 0.
	MaxApplicants *int `json:"max_applicants,omitempty"`

	// Show a project again after (days): For a saved formula, how long before a delivered project may come back.
	// From 1 to 365. Default: 30.
	ResurfaceAfterDays *int `json:"resurface_after_days,omitempty"`

	// Most rows: The run stops at this many rows, best first. From 1 to 1000. Default: 100.
	MaxResults *int `json:"max_results,omitempty"`
}

// FormulaID is "freelance-gig-feed".
func (FreelanceGigFeedParams) FormulaID() string { return "freelance-gig-feed" }

// FundedAndHiringParams are the parameters of Funded and Hiring (`funded-and-hiring`): Companies that just raised a
// round and are hiring on their own careers page, with the founder or talent lead to contact.
//
// Any other field left nil or empty is not sent, and the formula's default applies; [Ptr] makes the pointer of a
// value.
type FundedAndHiringParams struct {
	// Industry words: Keep announcements whose headline or text holds one of these words. Leave empty for any. At
	// most 20 items. Default: [].
	Industries []string `json:"industries,omitempty"`

	// Announced within (days): Only rounds announced in the last this many days. From 1 to 60. Default: 14.
	PublishedWithinDays *int `json:"published_within_days,omitempty"`

	// Least open roles: Openings the company's careers page must list. 0 keeps companies with none. At least 0.
	// Default: 1.
	MinOpenRoles *int `json:"min_open_roles,omitempty"`

	// Count openings no tracking system confirms: Many young companies list jobs on a plain careers page. Those
	// count towards the least open roles, and are shown apart. Default: true.
	CountUnconfirmed *bool `json:"count_unconfirmed,omitempty"`

	// Searches: The words each wire is searched for. Each search reads the newest 50 releases that hold them. 1 to
	// 8 items. Default: ["raises Series A", "seed funding", "raises Series B", "funding round"].
	SearchTerms []string `json:"search_terms,omitempty"`

	// Newswires: Which newswires to read. Each one of "prnewswire", "globenewswire". At least 1 item. Default:
	// ["prnewswire", "globenewswire"].
	Wires []string `json:"wires,omitempty"`

	// Releases read: How many of the newest likely funding releases are read in full. From 5 to 100. Default: 50.
	ReleasesRead *int `json:"releases_read,omitempty"`

	// Who to contact: The roles of the people to find at each company. Default: ["Founder", "CEO", "Head of
	// Talent", "Head of People"].
	ContactTitles []string `json:"contact_titles,omitempty"`

	// Contacts per company: 0 turns contacts off. From 0 to 5. Default: 2.
	ContactsPerCompany *int `json:"contacts_per_company,omitempty"`

	// Only addresses that pass the check: Keeps addresses whose domain takes mail, and those a verification
	// provider confirmed. Default: true.
	RequireVerifiedEmail *bool `json:"require_verified_email,omitempty"`

	// Show a company again after (days): For a saved formula, how long before a delivered company may come back.
	// From 1 to 365. Default: 90.
	ResurfaceAfterDays *int `json:"resurface_after_days,omitempty"`

	// Most rows: The run stops at this many rows, best first. From 1 to 1000. Default: 100.
	MaxResults *int `json:"max_results,omitempty"`
}

// FormulaID is "funded-and-hiring".
func (FundedAndHiringParams) FormulaID() string { return "funded-and-hiring" }

// HiringSignalOutboundParams are the parameters of Hiring-Signal Outbound (`hiring-signal-outbound`): Companies
// whose job posts reveal a need for what you sell, with the decision-maker and a true line to open with.
//
// SignalTerms and Locations are required. Any other field left nil or empty is not sent, and the formula's default
// applies; [Ptr] makes the pointer of a value.
type HiringSignalOutboundParams struct {
	// Roles, tools or phrases: What the postings must name. Each is searched on the boards and looked for in the
	// posting's text. 1 to 10 items. Required.
	SignalTerms []string `json:"signal_terms"`

	// Locations: A city, region or country for each place to search. 1 to 10 items. Required.
	Locations []string `json:"locations"`

	// Signal: role: companies hiring these roles. technology: postings that name these tools. pain: postings that
	// describe these problems. growth: a department hiring several people. One of "role", "technology", "pain",
	// "growth". Default: "technology".
	SignalType *string `json:"signal_type,omitempty"`

	// Posted within (days): Only postings from the last this many days. From 1 to 30. Default: 7.
	PostedWithinDays *int `json:"posted_within_days,omitempty"`

	// Job boards: Which boards to search. Reed and Totaljobs cover the UK; Dice is technology roles. Each one of
	// "linkedin", "indeed", "dice", "reed", "totaljobs". At least 1 item. Default: ["linkedin", "indeed"].
	Sources []string `json:"sources,omitempty"`

	// Postings read per search: How many postings to read for each term, place and board. From 5 to 100. Default:
	// 25.
	ResultsPerSearch *int `json:"results_per_search,omitempty"`

	// Postings read in full: How many of the newest postings are opened to read their whole text, where the search
	// gave only a title. From 5 to 300. Default: 60.
	PostingsToRead *int `json:"postings_to_read,omitempty"`

	// Postings for growth: For the growth signal, how many postings one department of a company needs in the
	// window. From 2 to 50. Default: 3.
	GrowthMinPostings *int `json:"growth_min_postings,omitempty"`

	// Company size (employees): Companies whose size is not known are kept. Default: {"min": 11, "max": 5000}.
	CompanySize *IntRange `json:"company_size,omitempty"`

	// Only these industries: Leave empty for any. Default: [].
	IndustriesInclude []string `json:"industries_include,omitempty"`

	// Not these industries: Companies in these industries are left out. Default: [].
	IndustriesExclude []string `json:"industries_exclude,omitempty"`

	// Companies to leave out: Names or domains, such as customers and open deals. At most 5000 items. Default: [].
	ExcludeCompanies []string `json:"exclude_companies,omitempty"`

	// Who to contact: The roles of the people to find. Leave empty to look for whoever decides for the role being
	// hired, such as the VP Sales for an SDR. Default: [].
	BuyerTitles []string `json:"buyer_titles,omitempty"`

	// Contact seniority: How senior the people to find should be. Each one of "manager", "director", "vp",
	// "c_level". Default: ["director", "vp", "c_level"].
	ContactSeniority []string `json:"contact_seniority,omitempty"`

	// Contacts per company: 0 turns contacts off. From 0 to 5. Default: 1.
	ContactsPerCompany *int `json:"contacts_per_company,omitempty"`

	// Only addresses that pass the check: Keeps addresses whose domain takes mail, and those a verification
	// provider confirmed. Default: true.
	RequireVerifiedEmail *bool `json:"require_verified_email,omitempty"`

	// Show a company again after (days): For a saved formula, how long before a company already delivered may come
	// back. From 1 to 365. Default: 90.
	ResurfaceAfterDays *int `json:"resurface_after_days,omitempty"`

	// Most rows: The run stops at this many rows, best first. From 1 to 2000. Default: 100.
	MaxResults *int `json:"max_results,omitempty"`
}

// FormulaID is "hiring-signal-outbound".
func (HiringSignalOutboundParams) FormulaID() string { return "hiring-signal-outbound" }

// HiringVelocityParams are the parameters of Hiring Velocity (`hiring-velocity`): Open roles at the companies you
// watch, counted by department from their own careers sites, with the change since the last run.
//
// Companies is required. Any other field left nil or empty is not sent, and the formula's default applies; [Ptr]
// makes the pointer of a value.
type HiringVelocityParams struct {
	// Companies: Each company as a stock ticker in capitals (DUOL), a domain (janestreet.com), or the address of
	// its careers page or job board (https://boards.greenhouse.io/wehrtyou). A job board address is the most
	// exact. 1 to 100 items. Required.
	Companies []string `json:"companies"`

	// Recent means posted within (days): A role posted in the last this many days counts as a recent role. From 1
	// to 90. Default: 30.
	RecentDays *int `json:"recent_days,omitempty"`

	// Only these departments: Leave empty for every department. Each one of "engineering", "data", "research",
	// "trading", "investments", "risk", "product", "design", "sales", "marketing", "customer_support", "finance",
	// "legal", "people", "it", "operations", "healthcare", "education", "skilled_trades", "retail_hospitality",
	// "other". Default: [].
	Departments []string `json:"departments,omitempty"`

	// Fewest open roles: A department with fewer open roles than this is left out. From 1 to 1000. Default: 1.
	MinOpenRoles *int `json:"min_open_roles,omitempty"`

	// Postings to leave out: Postings whose title holds one of these words are not counted; by default, evergreen
	// talent lists and recruiting events, which are not roles. At most 50 items. Default: a list of 12 values,
	// which `get` on the formula shows.
	ExcludeTitleWords []string `json:"exclude_title_words,omitempty"`

	// Count jobs read from the careers page itself: Off: only jobs from the company's applicant tracking system
	// are counted, which is exact. On: jobs read from the text of a careers page that has none are counted too;
	// such a page can list things that are not jobs, or show only some. Default: false.
	IncludePageReads *bool `json:"include_page_reads,omitempty"`

	// Most rows: The run stops at this many departments, biggest movers first. From 1 to 5000. Default: 500.
	MaxResults *int `json:"max_results,omitempty"`
}

// FormulaID is "hiring-velocity".
func (HiringVelocityParams) FormulaID() string { return "hiring-velocity" }

// HomeServicesProsParams are the parameters of Home-Services Pros (`home-services-pros`): The plumbers,
// electricians or cleaners of a place, with ratings, licences, years in business and an email.
//
// Trade and City are required. Any other field left nil or empty is not sent, and the formula's default applies;
// [Ptr] makes the pointer of a value.
type HomeServicesProsParams struct {
	// Trade: The trade in the plural, as the sites write it, such as plumbers, electricians or house cleaning.
	// Required.
	Trade string `json:"trade"`

	// Town or city: One town or city. Required.
	City string `json:"city"`

	// State or county: For the United States, the two-letter state, which Thumbtack needs. For the United Kingdom,
	// leave empty or give the county. Default: "".
	Region *string `json:"region,omitempty"`

	// Country: The United States or the United Kingdom; the marketplaces read depend on it. One of "US", "GB".
	// Default: "US".
	Country *string `json:"country,omitempty"`

	// Where to look: Each is read only for the country it covers. Each one of "thumbtack", "bbb", "checkatrade",
	// "google_maps". At least 1 item. Default: ["thumbtack", "bbb", "checkatrade", "google_maps"].
	Sources []string `json:"sources,omitempty"`

	// Lowest rating: Out of 5 on most sites and of 10 on Checkatrade; 0 keeps those with no rating. From 0 to 10.
	// Default: 0.
	MinRating *float64 `json:"min_rating,omitempty"`

	// Fewest reviews: Businesses with fewer reviews are left out. At least 0. Default: 0.
	MinReviews *int `json:"min_reviews,omitempty"`

	// Only businesses with an email: Keeps the businesses whose own website publishes an email address that passes
	// the check. Default: false.
	RequireEmail *bool `json:"require_email,omitempty"`

	// Read each marketplace profile: Years in business, licences and background checks on Thumbtack; website and
	// accreditations on Checkatrade. Default: true.
	ReadProfiles *bool `json:"read_profiles,omitempty"`

	// Show a business again after (days): For a saved formula, how long before a delivered business may come back.
	// From 1 to 365. Default: 120.
	ResurfaceAfterDays *int `json:"resurface_after_days,omitempty"`

	// Most rows: The run stops at this many businesses, most reviewed first. From 1 to 2000. Default: 150.
	MaxResults *int `json:"max_results,omitempty"`
}

// FormulaID is "home-services-pros".
func (HomeServicesProsParams) FormulaID() string { return "home-services-pros" }

// HotelRateShopParams are the parameters of Hotel Rate Shop (`hotel-rate-shop`): What your competing hotels charge
// a night for the same stays, from Google Hotels, Booking.com and MakeMyTrip, with each change.
//
// Hotels is required. Any other field left nil or empty is not sent, and the formula's default applies; [Ptr] makes
// the pointer of a value.
type HotelRateShopParams struct {
	// Hotels: Each hotel's name as Google Hotels shows it. Its rate there is read when the first hotel Google
	// Hotels finds has that name. 1 to 20 items. Required.
	Hotels []string `json:"hotels"`

	// Town: The town the hotels are in, added to each name when it is searched, such as London. Leave empty when
	// the names say it. Default: "".
	City *string `json:"city,omitempty"`

	// Stays starting in (days): How many days after each run the stays begin. A saved formula asks the same
	// distances every time it runs. 1 to 6 items. Default: [7, 14, 30].
	DaysAhead []int `json:"days_ahead,omitempty"`

	// Booking.com pages: The hotels' pages on Booking.com, such as
	// https://www.booking.com/hotel/gb/strandpalace.html, to read their room prices there too. At most 20 items.
	// Default: [].
	BookingPages []string `json:"booking_pages,omitempty"`

	// MakeMyTrip hotels: The hotels' pages on MakeMyTrip, or their hotel ids there, to read their rates in rupees
	// too. At most 20 items. Default: [].
	MakemytripHotels []string `json:"makemytrip_hotels,omitempty"`

	// Nights: How long each stay is. Prices are shown a night. From 1 to 14. Default: 1.
	Nights *int `json:"nights,omitempty"`

	// Guests: Adults in one room. On Booking.com only the offers for at least this many guests count. From 1 to 8.
	// Default: 2.
	Adults *int `json:"adults,omitempty"`

	// Currency: The currency of the Google Hotels and Booking.com prices. MakeMyTrip prices are in rupees. One of
	// "USD", "EUR", "GBP", "AED", "INR", "AUD", "CAD", "JPY", "CHF", "SGD". Default: "USD".
	Currency *string `json:"currency,omitempty"`

	// Most rows: The run stops at this many rates, changes and sold-out hotels first. From 1 to 2000. Default:
	// 500.
	MaxResults *int `json:"max_results,omitempty"`
}

// FormulaID is "hotel-rate-shop".
func (HotelRateShopParams) FormulaID() string { return "hotel-rate-shop" }

// InfluencerDiscoveryParams are the parameters of Influencer Discovery (`influencer-discovery`): Creators who make
// videos about your topic, with followers, views, engagement, their links and a published business email.
//
// Topics is required. Any other field left nil or empty is not sent, and the formula's default applies; [Ptr] makes
// the pointer of a value.
type InfluencerDiscoveryParams struct {
	// Topics: What the creators make videos about, as a viewer would search it, such as home espresso. A hashtag
	// works too. 1 to 10 items. Required.
	Topics []string `json:"topics"`

	// Subscribers: The size of creator to keep, by YouTube subscribers. Default: {"min": 10000, "max": 1000000}.
	Subscribers *IntRange `json:"subscribers,omitempty"`

	// Countries: Only creators whose channel states one of these countries, as YouTube writes it (United States,
	// Germany). Leave empty for any; a channel that states none is then kept. At most 20 items. Default: [].
	Countries []string `json:"countries,omitempty"`

	// Least videos on the topic: How many of a creator's videos must turn up in the topic searches; 2 or more
	// keeps creators who cover it regularly. From 1 to 10. Default: 1.
	MinTopicVideos *int `json:"min_topic_videos,omitempty"`

	// Creators to leave out: Handles (@name) or channel ids of creators you already work with. At most 1000 items.
	// Default: [].
	ExcludeChannels []string `json:"exclude_channels,omitempty"`

	// Videos read per topic: How many of each topic's videos are read to find creators. Each video read is one
	// unit. From 10 to 100. Default: 40.
	VideosPerTopic *int `json:"videos_per_topic,omitempty"`

	// Creators read: How many creators' channels are read, the most on topic first; each channel, profile and
	// website read is one unit. From 5 to 200. Default: 40.
	MaxCreators *int `json:"max_creators,omitempty"`

	// Also read: The creator's Instagram and TikTok accounts, where their channel links to them. Each one of
	// "instagram", "tiktok". Default: ["instagram", "tiktok"].
	OtherPlatforms []string `json:"other_platforms,omitempty"`

	// Business email from their website: Reads the creator's own website, where their channel links to one, for
	// the business address it publishes. Default: true.
	WebsiteEmail *bool `json:"website_email,omitempty"`

	// Most rows: The run stops at this many creators, best fit first. From 1 to 1000. Default: 100.
	MaxResults *int `json:"max_results,omitempty"`
}

// FormulaID is "influencer-discovery".
func (InfluencerDiscoveryParams) FormulaID() string { return "influencer-discovery" }

// JobOrderHunterParams are the parameters of Job Order Hunter (`job-order-hunter`): Companies newly hiring for your
// roles and places, qualified by size and industry, with the person to contact.
//
// JobTitles and Locations are required. Any other field left nil or empty is not sent, and the formula's default
// applies; [Ptr] makes the pointer of a value.
type JobOrderHunterParams struct {
	// Job titles: The roles you recruit for, each searched as you type it. 1 to 10 items. Required.
	JobTitles []string `json:"job_titles"`

	// Locations: A city, region or country for each place you cover. 1 to 10 items. Required.
	Locations []string `json:"locations"`

	// Posted within (days): Only jobs posted in the last this many days. From 1 to 30. Default: 1.
	PostedWithinDays *int `json:"posted_within_days,omitempty"`

	// Job boards: Which boards to search. The careers pages of target companies are read when you choose them and
	// name the companies. Each one of "linkedin", "indeed", "google_jobs", "dice", "reed", "totaljobs",
	// "stepstone", "seek", "simplyhired", "wellfound", "xing", "irishjobs", "ycombinator", "ats_boards". At least
	// 1 item. Default: ["linkedin", "indeed", "google_jobs"].
	Sources []string `json:"sources,omitempty"`

	// Search other names of the same job: Also searches the names the same job is often advertised under, such as
	// RN for Registered Nurse. Default: false.
	TitleSynonyms *bool `json:"title_synonyms,omitempty"`

	// Target companies: Domains or careers pages whose own openings are read when Careers pages is chosen as a
	// board. At most 100 items. Default: [].
	TargetCompanies []string `json:"target_companies,omitempty"`

	// Jobs read per search: How many jobs to read for each title, place and board. From 5 to 200. Default: 50.
	ResultsPerSearch *int `json:"results_per_search,omitempty"`

	// Employment types: Leave empty for any. Each one of "full_time", "part_time", "contract", "temporary".
	// Default: [].
	EmploymentTypes []string `json:"employment_types,omitempty"`

	// Company size (employees): Companies whose size is not known are kept. Default: {"min": 11, "max": 5000}.
	CompanySize *IntRange `json:"company_size,omitempty"`

	// Only these industries: Leave empty for any. Default: [].
	IndustriesInclude []string `json:"industries_include,omitempty"`

	// Not these industries: Companies in these industries are left out. Default: [].
	IndustriesExclude []string `json:"industries_exclude,omitempty"`

	// Leave out staffing agencies: Recruiters' own postings are rarely a job order. Default: true.
	ExcludeStaffingAgencies *bool `json:"exclude_staffing_agencies,omitempty"`

	// Companies to leave out: Names or domains, such as your existing clients. At most 5000 items. Default: [].
	ExcludeCompanies []string `json:"exclude_companies,omitempty"`

	// Least open roles: Matching roles a company needs to have open. At least 1. Default: 1.
	MinOpenRoles *int `json:"min_open_roles,omitempty"`

	// Words to leave out of titles: For example intern or unpaid. Default: [].
	TitleExcludeKeywords []string `json:"title_exclude_keywords,omitempty"`

	// Check their careers pages: Reads the careers page of the companies with the most matching roles; a role
	// listed there ranks higher. Default: true.
	CheckCareersPages *bool `json:"check_careers_pages,omitempty"`

	// Careers pages to read: How many companies' careers pages to read, those with the most matching roles first.
	// From 1 to 100. Default: 25.
	CareersPagesMax *int `json:"careers_pages_max,omitempty"`

	// Who to contact: The roles of the people to find at each company. Default: ["Talent Acquisition", "HR
	// Manager", "Hiring Manager"].
	ContactTitles []string `json:"contact_titles,omitempty"`

	// Contact seniority: How senior the people to find should be. Each one of "manager", "director", "vp",
	// "c_level". Default: ["manager", "director", "vp", "c_level"].
	ContactSeniority []string `json:"contact_seniority,omitempty"`

	// Contacts per company: 0 turns contacts off. From 0 to 5. Default: 2.
	ContactsPerCompany *int `json:"contacts_per_company,omitempty"`

	// Only addresses that pass the check: Keeps addresses whose domain takes mail, and those a verification
	// provider confirmed. Default: true.
	RequireVerifiedEmail *bool `json:"require_verified_email,omitempty"`

	// Show a company again after (days): For a saved formula, how long before a delivered row may come back. From
	// 1 to 365. Default: 45.
	ResurfaceAfterDays *int `json:"resurface_after_days,omitempty"`

	// Most rows: The run stops at this many rows, best first. From 1 to 5000. Default: 200.
	MaxResults *int `json:"max_results,omitempty"`
}

// FormulaID is "job-order-hunter".
func (JobOrderHunterParams) FormulaID() string { return "job-order-hunter" }

// JobSeekerTargetsParams are the parameters of Job Seeker Targets (`job-seeker-targets`): Employers hiring now for
// the job you want, with the person who decides, their address and a first line to send.
//
// TargetTitles and Locations are required. Any other field left nil or empty is not sent, and the formula's default
// applies; [Ptr] makes the pointer of a value.
type JobSeekerTargetsParams struct {
	// Jobs you want: The job titles you are applying for, each searched as you type it. 1 to 5 items. Required.
	TargetTitles []string `json:"target_titles"`

	// Where: A city, region or country for each place you would work, or Remote. 1 to 5 items. Required.
	Locations []string `json:"locations"`

	// Posted within (days): Only postings from the last this many days, while the role is still open. From 1 to
	// 30. Default: 3.
	PostedWithinDays *int `json:"posted_within_days,omitempty"`

	// Job boards: Which boards to search. Reed and Totaljobs cover the UK; Dice is technology roles in the US.
	// Each one of "linkedin", "indeed", "reed", "totaljobs", "dice". At least 1 item. Default: ["linkedin",
	// "indeed"].
	Sources []string `json:"sources,omitempty"`

	// Jobs read per search: How many postings to read for each job, place and board. From 5 to 100. Default: 25.
	ResultsPerSearch *int `json:"results_per_search,omitempty"`

	// Company size (employees): Companies whose size is not known are kept. Default: {"min": 1, "max": 100000}.
	CompanySize *IntRange `json:"company_size,omitempty"`

	// Only these industries: Leave empty for any. Default: [].
	IndustriesInclude []string `json:"industries_include,omitempty"`

	// Employers to leave out: Names or domains, such as your current employer or places you already applied to. At
	// most 2000 items. Default: [].
	ExcludeCompanies []string `json:"exclude_companies,omitempty"`

	// Only titles that name the job: Keeps postings whose title names one of your jobs as whole words; boards also
	// return near matches. Default: true.
	TitleMustMatch *bool `json:"title_must_match,omitempty"`

	// Words to leave out of titles: For example intern, senior or contract. Default: [].
	TitleExcludeKeywords []string `json:"title_exclude_keywords,omitempty"`

	// Who to write to: Leave empty to look for whoever decides on that kind of hire, such as the CTO for a
	// developer. Default: [].
	ContactTitles []string `json:"contact_titles,omitempty"`

	// People per employer: 0 turns contacts off. From 0 to 3. Default: 1.
	ContactsPerCompany *int `json:"contacts_per_company,omitempty"`

	// Only addresses that pass the check: Keeps addresses whose domain takes mail, and those a verification
	// provider confirmed. Default: true.
	RequireVerifiedEmail *bool `json:"require_verified_email,omitempty"`

	// Most rows: The run stops at this many rows, newest postings first. From 1 to 500. Default: 50.
	MaxResults *int `json:"max_results,omitempty"`
}

// FormulaID is "job-seeker-targets".
func (JobSeekerTargetsParams) FormulaID() string { return "job-seeker-targets" }

// KnowledgeBaseCrawlParams are the parameters of Knowledge Base Crawl (`knowledge-base-crawl`): A site's or docs'
// pages as clean Markdown, one row a page; saved, it returns only the pages that are new, changed or gone.
//
// StartURL is required. Any other field left nil or empty is not sent, and the formula's default applies; [Ptr]
// makes the pointer of a value.
type KnowledgeBaseCrawlParams struct {
	// Start address: The address of the site or the documentation, such as https://docs.example.com/guide/. Only
	// pages of the same site are read. Required.
	StartURL string `json:"start_url"`

	// Most pages read: Each page read is one unit. Following links reads at most 100 pages a run; with a sitemap,
	// up to 2000. From 1 to 2000. Default: 50.
	MaxPages *int `json:"max_pages,omitempty"`

	// Only these sections: Read only pages whose path starts with one of these, such as /docs/. Leave empty for
	// the whole site below the start address. At most 20 items. Default: [].
	IncludePaths []string `json:"include_paths,omitempty"`

	// Sitemap: To read the pages a sitemap lists instead of following links, its address
	// (https://example.com/sitemap.xml) or the site's bare address. Only listed pages under the start address are
	// read. Default: "".
	SitemapURL *string `json:"sitemap_url,omitempty"`

	// Leave out these sections: Pages whose path starts with one of these are not read, such as /blog/ or
	// /changelog/. At most 20 items. Default: [].
	ExcludePaths []string `json:"exclude_paths,omitempty"`

	// Link depth: How many links away from the start address a page may be when links are followed. From 0 to 10.
	// Default: 3.
	MaxDepth *int `json:"max_depth,omitempty"`

	// Keep menus, headers and footers: Off, each page's Markdown is its main content only, which is what an index
	// wants. Default: false.
	WholePage *bool `json:"whole_page,omitempty"`

	// Least words on a page: Leaves out near-empty pages, such as redirects and section stubs. From 0 to 5000.
	// Default: 20.
	MinWords *int `json:"min_words,omitempty"`

	// Report pages that are gone: A saved formula also returns the pages the last run read that this one did not
	// find. Turn it off when the site has more pages than a run reads. Default: true.
	ReportRemoved *bool `json:"report_removed,omitempty"`

	// Most rows: The run stops at this many pages, changed ones first. From 1 to 2000. Default: 500.
	MaxResults *int `json:"max_results,omitempty"`
}

// FormulaID is "knowledge-base-crawl".
func (KnowledgeBaseCrawlParams) FormulaID() string { return "knowledge-base-crawl" }

// LandlordAgentLeadsParams are the parameters of Letting Agent Leads (UK) (`landlord-agent-leads`): The letting
// agents and letting businesses with the most rentals on Rightmove and OnTheMarket in a UK area, with their branch
// phone.
//
// Locations is required. Any other field left nil or empty is not sent, and the formula's default applies; [Ptr]
// makes the pointer of a value.
type LandlordAgentLeadsParams struct {
	// Places: UK towns, areas or postcodes, such as Manchester, Leeds or SW11. 1 to 10 items. Required.
	Locations []string `json:"locations"`

	// Least rentals on the market: Only agents advertising at least this many of the rentals read. From 1 to 100.
	// Default: 1.
	MinListings *int `json:"min_listings,omitempty"`

	// Portals: Where to read the rentals. Each one of "rightmove", "onthemarket". At least 1 item. Default:
	// ["rightmove", "onthemarket"].
	Portals []string `json:"portals,omitempty"`

	// Result pages read: Which pages of each portal's results to read for each place, about 25 rentals a page.
	// Each one of 1, 2, 3, 4, 5. At least 1 item. Default: [1, 2, 3].
	Pages []int `json:"pages,omitempty"`

	// Radius (miles): Also read rentals this far around each place. One of 0, 0.25, 0.5, 1, 3, 5, 10. Default: 0.
	RadiusMiles *float64 `json:"radius_miles,omitempty"`

	// Most rent a month (GBP): 0 for no ceiling. At least 0. Default: 0.
	MaxRent *int `json:"max_rent,omitempty"`

	// Agents to leave out: Agents you already work with, by the name the portal shows. At most 2000 items.
	// Default: [].
	ExcludeAgents []string `json:"exclude_agents,omitempty"`

	// Show an agent again after (days): For a saved formula, how long before a delivered agent may come back. From
	// 1 to 365. Default: 60.
	ResurfaceAfterDays *int `json:"resurface_after_days,omitempty"`

	// Most rows: The run stops at this many agents, most rentals first. From 1 to 1000. Default: 100.
	MaxResults *int `json:"max_results,omitempty"`
}

// FormulaID is "landlord-agent-leads".
func (LandlordAgentLeadsParams) FormulaID() string { return "landlord-agent-leads" }

// LaunchBatchListsParams are the parameters of Launch and Batch Lists (`launch-batch-lists`): The startups of a Y
// Combinator batch, newest launches first, with their founders and a way to reach each one.
//
// Batch is required. Any other field left nil or empty is not sent, and the formula's default applies; [Ptr] makes
// the pointer of a value.
type LaunchBatchListsParams struct {
	// Batch: As Y Combinator writes it, such as Fall 2026 or F26. Required.
	Batch string `json:"batch"`

	// Industry: As the directory names it, such as Fintech, Healthcare or B2B. Leave empty for every industry.
	Industry *string `json:"industry,omitempty"`

	// Only companies that are hiring: Keeps the companies the directory marks as hiring. Default: false.
	HiringOnly *bool `json:"hiring_only,omitempty"`

	// Companies to leave out: Names or domains, such as your existing customers. At most 5000 items. Default: [].
	ExcludeCompanies []string `json:"exclude_companies,omitempty"`

	// Largest team: Companies with a bigger team are left out. At least 1. Default: 500.
	MaxTeamSize *int `json:"max_team_size,omitempty"`

	// Who to contact: The roles a connected people provider looks for at each company. Default: ["Founder",
	// "Co-Founder", "CEO"].
	ContactTitles []string `json:"contact_titles,omitempty"`

	// Contacts per company: People a connected provider finds at each company. 0 keeps one row per company. From 0
	// to 5. Default: 1.
	ContactsPerCompany *int `json:"contacts_per_company,omitempty"`

	// Most companies: How many companies of the batch to read, newest launches first. From 1 to 500. Default: 100.
	MaxResults *int `json:"max_results,omitempty"`
}

// FormulaID is "launch-batch-lists".
func (LaunchBatchListsParams) FormulaID() string { return "launch-batch-lists" }

// LinkIntersectParams are the parameters of Link Intersect (`link-intersect`): Sites that link to two or more of
// your competitors but not to you, each with the page and the email the site publishes.
//
// YourDomain and Competitors are required. Any other field left nil or empty is not sent, and the formula's default
// applies; [Ptr] makes the pointer of a value.
type LinkIntersectParams struct {
	// Your domain: Your site's domain, such as example.com. Pages that already link to it are left out. Required.
	YourDomain string `json:"your_domain"`

	// Competitors' domains: Two to six domains, such as rival.com. Every pair of them is searched for, so each one
	// added costs more. 2 to 6 items. Required.
	Competitors []string `json:"competitors"`

	// Topic words: Words added to every search to keep it on your subject, such as crm. Leave empty to search the
	// domains alone. Default: "".
	Topic *string `json:"topic,omitempty"`

	// Competitors linked, at least: A page must link to at least this many of your competitors. With 1, each
	// competitor is also searched for alone. From 1 to 6. Default: 2.
	MinCompetitors *int `json:"min_competitors,omitempty"`

	// Count a competitor named without a link: A page that writes a competitor's domain in its text without
	// linking to it counts toward the least number. Default: true.
	CountNamed *bool `json:"count_named,omitempty"`

	// Results pages per search: How many pages of results to read for each search, about ten results a page. From
	// 1 to 3. Default: 1.
	SearchPages *int `json:"search_pages,omitempty"`

	// Pages to check: How many of the pages found are read to see whom they link to, those found by the most
	// searches first; never more than the most rows. From 5 to 100. Default: 30.
	PagesToRead *int `json:"pages_to_read,omitempty"`

	// Sites to leave out: Domains never shown, such as social networks and review sites where no editor can add
	// your link. At most 2000 items. Default: a list of 19 values, which `get` on the formula shows.
	ExcludeDomains []string `json:"exclude_domains,omitempty"`

	// Search country: The two-letter country whose search results are read, such as us, gb or de. Default: "us".
	Country *string `json:"country,omitempty"`

	// Search language: The language of the results, such as en or de. Default: "en".
	Language *string `json:"language,omitempty"`

	// Only sites with an email: Keeps the sites whose own website publishes an email address that passes the
	// check. Default: false.
	RequireEmail *bool `json:"require_email,omitempty"`

	// Show a site again after (days): For a saved formula, how long before a delivered site may come back. From 1
	// to 365. Default: 90.
	ResurfaceAfterDays *int `json:"resurface_after_days,omitempty"`

	// Most rows: The run stops at this many sites, those linking to the most competitors first. From 1 to 1000.
	// Default: 100.
	MaxResults *int `json:"max_results,omitempty"`
}

// FormulaID is "link-intersect".
func (LinkIntersectParams) FormulaID() string { return "link-intersect" }

// ListingComplianceParams are the parameters of Listing Compliance (`listing-compliance`): Your Amazon listings
// checked against your own content rules (brand, title words, images, bullets, Buy Box seller), with changes.
//
// Asins is required. Any other field left nil or empty is not sent, and the formula's default applies; [Ptr] makes
// the pointer of a value.
type ListingComplianceParams struct {
	// Your listings: The ASINs to check. 1 to 200 items. Required.
	Asins []string `json:"asins"`

	// Brand the listing must show: As your brand is written, such as JBL. Empty skips the check. Default: "".
	Brand *string `json:"brand,omitempty"`

	// Words every title must hold: Such as your brand or the product line. At most 20 items. Default: [].
	RequiredTitleWords []string `json:"required_title_words,omitempty"`

	// Amazon marketplace. One of "us", "uk", "de", "fr", "it", "es", "ca". Default: "us".
	Marketplace *string `json:"marketplace,omitempty"`

	// Words no title may hold. At most 30 items. Default: ["refurbished", "renewed", "used", "replica",
	// "compatible with", "generic", "unbranded"].
	ForbiddenWords []string `json:"forbidden_words,omitempty"`

	// Least images. From 0 to 20. Default: 6.
	MinImages *int `json:"min_images,omitempty"`

	// Least bullet points. From 0 to 10. Default: 5.
	MinBullets *int `json:"min_bullets,omitempty"`

	// Longest title (characters). From 50 to 500. Default: 200.
	MaxTitleLength *int `json:"max_title_length,omitempty"`

	// Sellers who may hold the Buy Box: Such as Amazon.com and your own storefront. Empty skips the check. At most
	// 100 items. Default: [].
	BuyBoxSellers []string `json:"buy_box_sellers,omitempty"`

	// Only listings with a breach or a change. Default: false.
	OnlyBreaches *bool `json:"only_breaches,omitempty"`

	// Most rows: The run stops at this many listings, most breaches first. From 1 to 500. Default: 200.
	MaxResults *int `json:"max_results,omitempty"`
}

// FormulaID is "listing-compliance".
func (ListingComplianceParams) FormulaID() string { return "listing-compliance" }

// LiteratureSweepParams are the parameters of Literature Sweep (`literature-sweep`): New papers on your subject
// from PubMed, arXiv and Crossref, each paper once, with its abstract, DOI and free full text.
//
// Query is required. Any other field left nil or empty is not sent, and the formula's default applies; [Ptr] makes
// the pointer of a value.
type LiteratureSweepParams struct {
	// Subject: Words that describe the subject; every word must appear in a paper's listing. PubMed also reads its
	// own field tags, such as review[pt], when the words are not searched as one phrase. Required.
	Query string `json:"query"`

	// Search the words as one phrase: Finds protein language model only where those words stand together, as a
	// term of the field. Crossref reads the words alone either way. Default: false.
	ExactPhrase *bool `json:"exact_phrase,omitempty"`

	// Published within (days): Only papers that appeared in the last this many days; 0 for any date. PubMed counts
	// the day a paper was added to it. From 0 to 3650. Default: 30.
	PublishedWithinDays *int `json:"published_within_days,omitempty"`

	// Where to search: PubMed for medicine and life sciences, arXiv for preprints in physics, maths and computing,
	// Crossref for journals of every field. Each one of "pubmed", "arxiv", "crossref", "openalex",
	// "semanticscholar". At least 1 item. Default: ["pubmed", "arxiv", "crossref"].
	Sources []string `json:"sources,omitempty"`

	// Only papers by these authors: A surname, or a name as papers print it (Doudna, Liu DR). Leave empty for any
	// author. At most 50 items. Default: [].
	Authors []string `json:"authors,omitempty"`

	// Only these journals: Part of a journal's or server's name; Nature also keeps Nature Methods. Leave empty for
	// any. At most 50 items. Default: [].
	Journals []string `json:"journals,omitempty"`

	// Include preprints: Papers not yet peer reviewed, from arXiv, bioRxiv and medRxiv. Default: true.
	IncludePreprints *bool `json:"include_preprints,omitempty"`

	// Only papers free to read: Keeps the papers with a link to free full text. Default: false.
	OpenAccessOnly *bool `json:"open_access_only,omitempty"`

	// Pages read per source: Each page is the 100 newest (Crossref, the 100 best matching) papers of one source;
	// raise it for a busy subject or a long window. From 1 to 3. Default: 1.
	Pages *int `json:"pages,omitempty"`

	// Read PubMed abstracts: PubMed lists papers without their abstracts; this reads each PubMed paper's own page
	// for its abstract and the day it went online. Default: true.
	ReadAbstracts *bool `json:"read_abstracts,omitempty"`

	// Most rows: The run stops at this many papers, best matching first. From 1 to 500. Default: 100.
	MaxResults *int `json:"max_results,omitempty"`
}

// FormulaID is "literature-sweep".
func (LiteratureSweepParams) FormulaID() string { return "literature-sweep" }

// LocalBusinessLeadsParams are the parameters of Local Business Leads (`local-business-leads`): Businesses of a
// kind in your places, with the email their own website publishes and facts to open with.
//
// BusinessTypes and Locations are required. Any other field left nil or empty is not sent, and the formula's
// default applies; [Ptr] makes the pointer of a value.
type LocalBusinessLeadsParams struct {
	// Kinds of business: What to search for, as you would type it into Google Maps, such as roofer or vegan
	// restaurant. 1 to 10 items. Required.
	BusinessTypes []string `json:"business_types"`

	// Places: A city, district or region for each place, with its state or country when the name is common. 1 to
	// 10 items. Required.
	Locations []string `json:"locations"`

	// Search each neighbourhood too: Searches again in the neighbourhoods and postcodes the first search found, to
	// reach past its few hundred results. Default: true.
	ExpandLocations *bool `json:"expand_locations,omitempty"`

	// Neighbourhoods per search: How many neighbourhoods or postcodes to search again, those with the most places
	// first. From 1 to 40. Default: 8.
	MaxSubAreas *int `json:"max_sub_areas,omitempty"`

	// Places read per search: How many places to read from each search. From 10 to 300. Default: 60.
	PlacesPerSearch *int `json:"places_per_search,omitempty"`

	// Website: Keep every business, only those with a website, or only those without one. One of "any", "has",
	// "missing". Default: "any".
	Website *string `json:"website,omitempty"`

	// Lowest rating: Businesses rated below this are left out. 0 keeps those with no rating. From 0 to 5. Default:
	// 0.
	MinRating *float64 `json:"min_rating,omitempty"`

	// Highest rating: Businesses rated above this are left out. From 0 to 5. Default: 5.
	MaxRating *float64 `json:"max_rating,omitempty"`

	// Fewest reviews: Businesses with fewer reviews are left out. At least 0. Default: 0.
	MinReviews *int `json:"min_reviews,omitempty"`

	// Most reviews: Businesses with more reviews are left out. At least 0. Default: 1000000.
	MaxReviews *int `json:"max_reviews,omitempty"`

	// Leave out chains: A name found five or more times in the results is taken for a chain and left out. Default:
	// true.
	ExcludeChains *bool `json:"exclude_chains,omitempty"`

	// Only businesses with an email: Keeps the businesses whose own website publishes an email address that passes
	// the check. Default: true.
	RequireEmail *bool `json:"require_email,omitempty"`

	// Only named addresses: Leaves out shared inboxes such as info@ and keeps addresses that passed every check.
	// Default: false.
	ExcludeRoleAddresses *bool `json:"exclude_role_addresses,omitempty"`

	// Only businesses with social profiles: Keeps the businesses whose website links to at least one social
	// profile. Default: false.
	RequireSocials *bool `json:"require_socials,omitempty"`

	// Site facts: What each website is built on, and its speed score on a phone. Default: false.
	IncludeSiteFacts *bool `json:"include_site_facts,omitempty"`

	// Slowest site score kept: With site facts, keeps the sites scoring this or less; 50 keeps the slow ones. From
	// 0 to 100. Default: 100.
	MaxSpeedScore *int `json:"max_speed_score,omitempty"`

	// Company facts: Industry and size from the company's LinkedIn page, when its website links to one. Default:
	// false.
	IncludeFirmographics *bool `json:"include_firmographics,omitempty"`

	// Show a business again after (days): For a saved formula, how long before a delivered business may come back.
	// From 1 to 365. Default: 180.
	ResurfaceAfterDays *int `json:"resurface_after_days,omitempty"`

	// Most rows: The run stops at this many businesses, most reviewed first. From 1 to 5000. Default: 500.
	MaxResults *int `json:"max_results,omitempty"`
}

// FormulaID is "local-business-leads".
func (LocalBusinessLeadsParams) FormulaID() string { return "local-business-leads" }

// LongOpenRolesParams are the parameters of Long-Open Roles (`long-open-roles`): Jobs open for a month or more at
// companies that are struggling to fill them, with who decides on the hire.
//
// JobTitles and Locations are required. Any other field left nil or empty is not sent, and the formula's default
// applies; [Ptr] makes the pointer of a value.
type LongOpenRolesParams struct {
	// Job titles: The roles you place, each searched as you type it. 1 to 10 items. Required.
	JobTitles []string `json:"job_titles"`

	// Locations: A city, region or country for each place you cover. 1 to 10 items. Required.
	Locations []string `json:"locations"`

	// Open at least (days): Only jobs that have been posted at least this many days. From 7 to 180. Default: 30.
	MinDaysOpen *int `json:"min_days_open,omitempty"`

	// Open at most (days): Older postings are often forgotten rather than open; they are left out. From 14 to 365.
	// Default: 120.
	MaxDaysOpen *int `json:"max_days_open,omitempty"`

	// Job boards: Which boards to search. Indeed lists jobs of the last 30 days only. Reed and Totaljobs cover the
	// UK; Dice is technology roles in the US. Each one of "linkedin", "indeed", "reed", "totaljobs", "dice". At
	// least 1 item. Default: ["linkedin", "indeed"].
	Sources []string `json:"sources,omitempty"`

	// Jobs read per search: How many postings to read for each title, place and board; most are newer than you
	// want, so read more. From 5 to 200. Default: 50.
	ResultsPerSearch *int `json:"results_per_search,omitempty"`

	// Company size (employees): Companies whose size is not known are kept. Default: {"min": 11, "max": 10000}.
	CompanySize *IntRange `json:"company_size,omitempty"`

	// Companies to leave out: Names or domains, such as your existing clients. At most 5000 items. Default: [].
	ExcludeCompanies []string `json:"exclude_companies,omitempty"`

	// Words to leave out of titles: For example intern or volunteer. Default: [].
	TitleExcludeKeywords []string `json:"title_exclude_keywords,omitempty"`

	// Who to contact: Leave empty to look for whoever decides on that kind of hire, such as the Director of
	// Nursing for a nurse. Default: [].
	ContactTitles []string `json:"contact_titles,omitempty"`

	// Contacts per company: 0 turns contacts off. From 0 to 5. Default: 1.
	ContactsPerCompany *int `json:"contacts_per_company,omitempty"`

	// Only addresses that pass the check: Keeps addresses whose domain takes mail, and those a verification
	// provider confirmed. Default: true.
	RequireVerifiedEmail *bool `json:"require_verified_email,omitempty"`

	// Most rows: The run stops at this many rows, the longest-open first. From 1 to 2000. Default: 100.
	MaxResults *int `json:"max_results,omitempty"`
}

// FormulaID is "long-open-roles".
func (LongOpenRolesParams) FormulaID() string { return "long-open-roles" }

// LookalikeCompaniesParams are the parameters of Lookalike Companies (`lookalike-companies`): Startups like the
// accounts you won, found where the Y Combinator directory files them, with the people to contact.
//
// SeedCompanies is required. Any other field left nil or empty is not sent, and the formula's default applies;
// [Ptr] makes the pointer of a value.
type LookalikeCompaniesParams struct {
	// Accounts you won: Company names as the Y Combinator directory lists them. 1 to 20 items. Required.
	SeedCompanies []string `json:"seed_companies"`

	// Places: As the directory names them, such as United States of America or Europe. Leave empty for anywhere.
	// Default: [].
	Regions []string `json:"regions,omitempty"`

	// Team size: Companies whose size is not known are kept. Default: {"min": 2, "max": 1000}.
	TeamSize *IntRange `json:"team_size,omitempty"`

	// Only active companies: Leaves out companies the directory marks as acquired, public or inactive. Default:
	// true.
	ActiveOnly *bool `json:"active_only,omitempty"`

	// Companies read per account: How many companies to read from the directory for each account you won. From 5
	// to 100. Default: 40.
	PeersPerAccount *int `json:"peers_per_account,omitempty"`

	// Companies to leave out: Names or domains, such as your existing customers. At most 5000 items. Default: [].
	ExcludeCompanies []string `json:"exclude_companies,omitempty"`

	// Who to contact: The roles a connected people provider looks for at each company. Default: ["CEO", "Founder",
	// "VP Sales", "Head of Operations"].
	BuyerTitles []string `json:"buyer_titles,omitempty"`

	// Contacts per company: 0 keeps one row per company. From 0 to 5. Default: 1.
	ContactsPerCompany *int `json:"contacts_per_company,omitempty"`

	// Most companies: The run stops at this many companies, the most alike first. From 1 to 1000. Default: 100.
	MaxResults *int `json:"max_results,omitempty"`
}

// FormulaID is "lookalike-companies".
func (LookalikeCompaniesParams) FormulaID() string { return "lookalike-companies" }

// MarketplaceFlipsParams are the parameters of Amazon Marketplace Flips (Europe) (`marketplace-flips`): Products
// cheaper on one European Amazon than another, by the same ASIN, with the price gap and margin in euros.
//
// Queries is required. Any other field left nil or empty is not sent, and the formula's default applies; [Ptr]
// makes the pointer of a value.
type MarketplaceFlipsParams struct {
	// Search words: Searched on the Amazon you buy from, such as lego technic. Not needed when ASINs are given. At
	// most 10 items. Required.
	Queries []string `json:"queries"`

	// Buy on: The Amazon you buy from (amazon.de, .it, .es or .nl). One of "de", "it", "es", "nl". Default: "de".
	BuyMarketplace *string `json:"buy_marketplace,omitempty"`

	// Sell on: The Amazons whose price is compared. The one you buy on is skipped. Each one of "de", "fr", "it",
	// "es", "nl". At least 1 item. Default: ["fr", "it", "es"].
	SellMarketplaces []string `json:"sell_marketplaces,omitempty"`

	// ASINs: Products to compare, instead of a search. At most 50 items. Default: [].
	Asins []string `json:"asins,omitempty"`

	// Least price gap (%): The selling price above the buying price, as a share of the buying price. From 0 to
	// 500. Default: 15.
	MinGapPct *float64 `json:"min_gap_pct,omitempty"`

	// Referral fee (%): The share of the sale price Amazon keeps. From 0 to 50. Default: 15.
	ReferralFeePct *float64 `json:"referral_fee_pct,omitempty"`

	// Fulfilment and shipping a unit (EUR): Your fulfilment fee and the cost of moving a unit between countries.
	// From 0 to 100. Default: 4.
	FulfilmentCost *float64 `json:"fulfilment_cost,omitempty"`

	// Products per search: How many products of each search are compared. From 1 to 48. Default: 15.
	ResultsPerSearch *int `json:"results_per_search,omitempty"`

	// Most rows: The run stops at this many pairs, widest gap first. From 1 to 1000. Default: 100.
	MaxResults *int `json:"max_results,omitempty"`
}

// FormulaID is "marketplace-flips".
func (MarketplaceFlipsParams) FormulaID() string { return "marketplace-flips" }

// MenuPriceTrackerParams are the parameters of Menu Price Tracker (`menu-price-tracker`): Every dish and price on
// the menus of the restaurants you watch, on Uber Eats, TheFork and Zomato, with each price change.
//
// Restaurants is required. Any other field left nil or empty is not sent, and the formula's default applies; [Ptr]
// makes the pointer of a value.
type MenuPriceTrackerParams struct {
	// Restaurants: The restaurants' pages on Uber Eats, TheFork or Zomato, one address each. 1 to 50 items.
	// Required.
	Restaurants []string `json:"restaurants"`

	// Only these menu sections: Words of the sections to keep, such as Mains or Pizza. Leave empty for the whole
	// menu. At most 30 items. Default: [].
	Sections []string `json:"sections,omitempty"`

	// Only dishes named with: Keeps the dishes whose name holds one of these words, such as burger or margherita.
	// Leave empty for every dish. At most 30 items. Default: [].
	DishWords []string `json:"dish_words,omitempty"`

	// Include set menus: TheFork's set menus and their prices, beside the à la carte dishes. Default: true.
	IncludeSetMenus *bool `json:"include_set_menus,omitempty"`

	// Most dishes per restaurant: The first this many dishes of each menu, in the menu's own order. From 1 to 500.
	// Default: 300.
	MaxItemsPerRestaurant *int `json:"max_items_per_restaurant,omitempty"`

	// Most rows: The run stops at this many dishes, price changes first. From 1 to 10000. Default: 2000.
	MaxResults *int `json:"max_results,omitempty"`
}

// FormulaID is "menu-price-tracker".
func (MenuPriceTrackerParams) FormulaID() string { return "menu-price-tracker" }

// OutOfStockOpportunityParams are the parameters of Out-of-Stock Opportunity (`out-of-stock-opportunity`): Amazon
// best sellers that are out of stock on Amazon, and the US retailers that have the same product in stock now.
//
// Categories is required. Any other field left nil or empty is not sent, and the formula's default applies; [Ptr]
// makes the pointer of a value.
type OutOfStockOpportunityParams struct {
	// Best Sellers categories: As the list's address names them, such as toys-and-games, electronics, kitchen, or
	// electronics/172456 for a category further down. 1 to 10 items. Required.
	Categories []string `json:"categories"`

	// Retailers: Where to look for stock. Walmart often answers automated visits with a check, so its products may
	// be missing. Each one of "target", "bestbuy", "newegg", "walmart". At least 1 item. Default: ["target",
	// "bestbuy"].
	Retailers []string `json:"retailers,omitempty"`

	// List pages: Page 1 is ranks 1 to 50, page 2 ranks 51 to 100. Each one of 1, 2. At least 1 item. Default: [1,
	// 2].
	Pages []int `json:"pages,omitempty"`

	// Results compared per retailer. From 1 to 5. Default: 2.
	CandidatesPerRetailer *int `json:"candidates_per_retailer,omitempty"`

	// Most rows: The run stops at this many product and retailer pairs, best rank first. From 1 to 500. Default:
	// 100.
	MaxResults *int `json:"max_results,omitempty"`
}

// FormulaID is "out-of-stock-opportunity".
func (OutOfStockOpportunityParams) FormulaID() string { return "out-of-stock-opportunity" }

// PeopleListWithOpenersParams are the parameters of People List with Openers (`people-list-with-openers`): The
// people in the roles you sell to at your target companies, with checked addresses and a line from their own site.
//
// Companies and Titles are required. Any other field left nil or empty is not sent, and the formula's default
// applies; [Ptr] makes the pointer of a value.
type PeopleListWithOpenersParams struct {
	// Companies: The website domains of the companies to find people at. 1 to 500 items. Required.
	Companies []string `json:"companies"`

	// Roles: The roles of the people to find. 1 to 10 items. Required.
	Titles []string `json:"titles"`

	// Seniority: How senior the people should be. Each one of "manager", "director", "vp", "c_level". Default:
	// ["director", "vp", "c_level"].
	Seniority []string `json:"seniority,omitempty"`

	// People per company: The most people to find at each company. From 1 to 10. Default: 2.
	PeoplePerCompany *int `json:"people_per_company,omitempty"`

	// Only addresses that pass the check: Keeps addresses whose domain takes mail, and those a verification
	// provider confirmed. Default: true.
	RequireVerifiedEmail *bool `json:"require_verified_email,omitempty"`

	// Most rows: The run stops at this many people. From 1 to 5000. Default: 500.
	MaxResults *int `json:"max_results,omitempty"`
}

// FormulaID is "people-list-with-openers".
func (PeopleListWithOpenersParams) FormulaID() string { return "people-list-with-openers" }

// PermitStatusWatchParams are the parameters of Permit Status Watch (`permit-status-watch`): Building permits in a
// US city read again on a schedule, with each status change and each newly issued permit.
//
// City is required. Any other field left nil or empty is not sent, and the formula's default applies; [Ptr] makes
// the pointer of a value.
type PermitStatusWatchParams struct {
	// City: The US city or county whose permits are watched. One of "austin", "chicago", "new_york",
	// "los_angeles", "san_francisco", "seattle", "dallas", "new_orleans", "cincinnati", "baton_rouge",
	// "montgomery_county_md", "prince_georges_county_md", "norfolk", "buffalo", "orlando", "mesa", "raleigh",
	// "washington_dc", "minneapolis", "detroit". Required.
	City string `json:"city"`

	// Permit numbers: The permits to watch, by their number at the city. Leave empty to watch the permits issued
	// in an area instead. At most 50 items. Default: [].
	PermitNumbers []string `json:"permit_numbers,omitempty"`

	// Address holds: Without permit numbers, only permits whose address holds this text, such as a street name.
	// Empty watches the whole city. Default: "".
	AddressContains *string `json:"address_contains,omitempty"`

	// Issued within (days): Without permit numbers, the permits issued in the last this many days are watched.
	// From 1 to 365. Default: 30.
	WithinDays *int `json:"within_days,omitempty"`

	// Work words: Without permit numbers, only permits whose type, work or description holds this, such as ROOF or
	// SOLAR. Default: "".
	Keyword *string `json:"keyword,omitempty"`

	// Permits read: Without permit numbers, the newest this many permits of the area are read on each run. From 10
	// to 200. Default: 200.
	PermitsToRead *int `json:"permits_to_read,omitempty"`

	// Only what changed: On a saved formula's later runs, show only permits that are new or whose status changed.
	// The first run shows every permit. Default: true.
	OnlyChanges *bool `json:"only_changes,omitempty"`

	// Most rows: The run stops at this many permits, changes first. From 1 to 1000. Default: 200.
	MaxResults *int `json:"max_results,omitempty"`
}

// FormulaID is "permit-status-watch".
func (PermitStatusWatchParams) FormulaID() string { return "permit-status-watch" }

// PermitTradeLeadsParams are the parameters of Permit Trade Leads (`permit-trade-leads`): Building permits just
// issued for your trades in a US city, with each job's contractor, or the busiest contractors ranked.
//
// City and Trades are required. Any other field left nil or empty is not sent, and the formula's default applies;
// [Ptr] makes the pointer of a value.
type PermitTradeLeadsParams struct {
	// City: The US city or county whose permits are read. Contractor names are published by Austin, Seattle, New
	// York, Dallas, Cincinnati, Baton Rouge, Orlando, Mesa, Raleigh and Chicago. One of "austin", "chicago",
	// "new_york", "los_angeles", "san_francisco", "seattle", "dallas", "new_orleans", "cincinnati", "baton_rouge",
	// "montgomery_county_md", "prince_georges_county_md", "norfolk", "buffalo", "orlando", "mesa", "raleigh",
	// "washington_dc", "minneapolis", "detroit". Required.
	City string `json:"city"`

	// Trades: Words in a permit's type, work or description, such as roof, solar, hvac, deck or addition. Each is
	// one search. 1 to 8 items. Required.
	Trades []string `json:"trades"`

	// Issued within (days): Only permits issued in the last this many days. From 1 to 365. Default: 30.
	IssuedWithinDays *int `json:"issued_within_days,omitempty"`

	// Rows are: One row a permit, or one row a contractor with how many permits it pulled. One of "permits",
	// "contractors". Default: "permits".
	View *string `json:"view,omitempty"`

	// Least job value (USD): Only permits whose declared value is at least this. Permits with no value are kept
	// when this is 0. At least 0. Default: 0.
	MinValue *float64 `json:"min_value,omitempty"`

	// Only permits that name a contractor: Leave out permits the city published without a contractor. Default:
	// false.
	RequireContractor *bool `json:"require_contractor,omitempty"`

	// Contractors to leave out: Names of contractors, such as your own business or current customers. At most 2000
	// items. Default: [].
	ExcludeContractors []string `json:"exclude_contractors,omitempty"`

	// Permits read per trade: The newest this many permits are read for each trade word. From 10 to 200. Default:
	// 200.
	PermitsPerTrade *int `json:"permits_per_trade,omitempty"`

	// Show a row again after (days): For a saved formula, how long before a delivered permit or contractor may
	// come back. From 1 to 365. Default: 90.
	ResurfaceAfterDays *int `json:"resurface_after_days,omitempty"`

	// Most rows: The run stops at this many rows, newest permits or busiest contractors first. From 1 to 2000.
	// Default: 300.
	MaxResults *int `json:"max_results,omitempty"`
}

// FormulaID is "permit-trade-leads".
func (PermitTradeLeadsParams) FormulaID() string { return "permit-trade-leads" }

// PriceSellerWatchParams are the parameters of Price and Seller Watch (`price-seller-watch`): Who sells each of
// your products on Amazon, Best Buy, Target and Newegg, at what price, what changed, and who is under your minimum.
//
// Products is required. Any other field left nil or empty is not sent, and the formula's default applies; [Ptr]
// makes the pointer of a value.
type PriceSellerWatchParams struct {
	// Products: One product an entry: its name as shoppers search it, and where you have them its UPC or EAN
	// (code), Amazon ASIN, brand, model and minimum advertised price (min_price). 1 to 50 items. Required.
	Products []map[string]any `json:"products"`

	// Marketplaces: Where to look. Walmart often answers automated visits with a check, so its listings may be
	// missing from a run. Each one of "amazon", "bestbuy", "target", "newegg", "walmart". At least 1 item.
	// Default: ["amazon", "bestbuy", "target"].
	Marketplaces []string `json:"marketplaces,omitempty"`

	// Authorised sellers: Seller names as the marketplaces show them, such as Amazon.com or Best Buy. Empty treats
	// every seller as authorised. At most 500 items. Default: [].
	AuthorisedSellers []string `json:"authorised_sellers,omitempty"`

	// Price change worth an alert (%): A price that moved by at least this much since the last run raises an
	// alert. From 0 to 100. Default: 5.
	PriceChangePct *float64 `json:"price_change_pct,omitempty"`

	// Amazon marketplace: Which Amazon to read. Best Buy, Target, Newegg and Walmart are read in the United
	// States. One of "us", "uk", "de", "fr", "it", "es", "ca". Default: "us".
	AmazonMarketplace *string `json:"amazon_marketplace,omitempty"`

	// Listings compared per marketplace: How many search results of each marketplace are checked against each
	// product. From 1 to 5. Default: 2.
	CandidatesPerMarketplace *int `json:"candidates_per_marketplace,omitempty"`

	// Only rows with an alert: Leave out the offers that raise no alert. Default: false.
	AlertsOnly *bool `json:"alerts_only,omitempty"`

	// Most rows: The run stops at this many offers, alerts first. From 1 to 2000. Default: 500.
	MaxResults *int `json:"max_results,omitempty"`
}

// FormulaID is "price-seller-watch".
func (PriceSellerWatchParams) FormulaID() string { return "price-seller-watch" }

// PropertyCompsParams are the parameters of Property Comps (`property-comps`): Recently sold and listed homes like
// yours in a US area, with sale price, price per square foot and the ZIP code average.
//
// Locations is required. Any other field left nil or empty is not sent, and the formula's default applies; [Ptr]
// makes the pointer of a value.
type PropertyCompsParams struct {
	// Places: ZIP codes, cities with their state ("Austin, TX") or neighbourhoods, as Zillow names them. 1 to 10
	// items. Required.
	Locations []string `json:"locations"`

	// Kinds of home: The kinds of home that count as comparable. Each one of "house", "condo", "townhouse",
	// "multi_family", "manufactured". Default: ["house"].
	HomeTypes []string `json:"home_types,omitempty"`

	// Least bedrooms: 0 for any. From 0 to 5. Default: 0.
	MinBeds *int `json:"min_beds,omitempty"`

	// Living area (sq ft): The size range of a comparable home, such as 1200 to 1800 for a 1500 sq ft home.
	// Default: {"min": 0, "max": 20000}.
	LivingArea *IntRange `json:"living_area,omitempty"`

	// Show homes for sale too: The matching homes on the market now, as the competition. Default: true.
	IncludeForSale *bool `json:"include_for_sale,omitempty"`

	// Sold within (days): Only homes sold in the last this many days. From 7 to 730. Default: 180.
	SoldWithinDays *int `json:"sold_within_days,omitempty"`

	// Homes read per place and kind: Zillow's first page of each place is read, up to this many homes, for sold
	// and for sale each. From 5 to 41. Default: 41.
	HomesPerPlace *int `json:"homes_per_place,omitempty"`

	// Most rows: The run stops at this many homes, newest sales first. Each sold home's page is read for its
	// price. From 1 to 500. Default: 60.
	MaxResults *int `json:"max_results,omitempty"`
}

// FormulaID is "property-comps".
func (PropertyCompsParams) FormulaID() string { return "property-comps" }

// PropertyDistressSignalsParams are the parameters of Property Distress Signals (`property-distress-signals`):
// Properties in a US city with recent code violations, distress permits or a stale listing, stacked by address and
// ranked.
//
// City is required. Any other field left nil or empty is not sent, and the formula's default applies; [Ptr] makes
// the pointer of a value.
type PropertyDistressSignalsParams struct {
	// City: The US city or county whose open data is read. These publish both code cases and permits. One of
	// "austin", "chicago", "new_york", "los_angeles", "seattle", "new_orleans", "cincinnati",
	// "montgomery_county_md", "norfolk", "buffalo". Required.
	City string `json:"city"`

	// Opened or issued within (days): Only cases opened, and permits issued, in the last this many days. From 1 to
	// 365. Default: 60.
	WithinDays *int `json:"within_days,omitempty"`

	// Least lists a property is on: 2 keeps only properties on two lists or more (a code case and a distress
	// permit, or either and a listing). From 1 to 3. Default: 1.
	MinLists *int `json:"min_lists,omitempty"`

	// Case words: Only cases whose type or description holds one of these, such as structure, vacant or abatement.
	// Each word is one search; empty reads the newest cases of any kind. At most 8 items. Default: [].
	CaseKeywords []string `json:"case_keywords,omitempty"`

	// Only cases still open: Leave out cases the city has closed. Default: false.
	OpenCasesOnly *bool `json:"open_cases_only,omitempty"`

	// Distress permit words: Permits whose type, work or description holds one of these are a signal. Each word is
	// one search. Empty leaves permits out. At most 8 items. Default: ["demoli", "fire damage", "vacant",
	// "unsafe"].
	PermitKeywords []string `json:"permit_keywords,omitempty"`

	// Cases read: The newest this many code cases are read. From 10 to 200. Default: 200.
	CasesToRead *int `json:"cases_to_read,omitempty"`

	// Permits read per word: The newest this many permits are read for each permit word. From 10 to 200. Default:
	// 100.
	PermitsPerWord *int `json:"permits_per_word,omitempty"`

	// Check each property on Zillow: Whether the home is listed, how long and at what price, with its size and
	// value estimate. One lookup a property. Default: true.
	CheckListings *bool `json:"check_listings,omitempty"`

	// Show a property again after (days): For a saved formula, how long before a delivered property may come back.
	// From 1 to 365. Default: 30.
	ResurfaceAfterDays *int `json:"resurface_after_days,omitempty"`

	// Most rows: The run stops at this many properties, most signals first. From 1 to 1000. Default: 100.
	MaxResults *int `json:"max_results,omitempty"`
}

// FormulaID is "property-distress-signals".
func (PropertyDistressSignalsParams) FormulaID() string { return "property-distress-signals" }

// RankTrackerParams are the parameters of Rank Tracker (`rank-tracker`): Your domain's and your competitors' Google
// Search positions for each keyword, with the top results and what moved since last run.
//
// Keywords and YourDomain are required. Any other field left nil or empty is not sent, and the formula's default
// applies; [Ptr] makes the pointer of a value.
type RankTrackerParams struct {
	// Keywords: The searches to track, each as people type it. 1 to 100 items. Required.
	Keywords []string `json:"keywords"`

	// Your domain: Your site's domain, such as example.com; its subdomains count as it. Required.
	YourDomain string `json:"your_domain"`

	// Competitors' domains: Up to ten domains, such as rival.com, each tracked for every keyword. At most 10
	// items. Default: [].
	Competitors []string `json:"competitors,omitempty"`

	// Country: The two-letter country whose results are read, such as us, gb or de. Default: "us".
	Country *string `json:"country,omitempty"`

	// Results pages read: 1 reads the top 10 results of each keyword, 2 the top 20, 3 the top 30; each page is one
	// search. From 1 to 3. Default: 1.
	Pages *int `json:"pages,omitempty"`

	// Language: The language of the results, such as en or de. Default: "en".
	Language *string `json:"language,omitempty"`

	// Top results shown: How many of the domains at the top of each keyword's results are listed on its rows. From
	// 0 to 10. Default: 3.
	TopResults *int `json:"top_results,omitempty"`

	// Most rows: The run stops at this many rows, one for each keyword and domain. From 1 to 1100. Default: 500.
	MaxResults *int `json:"max_results,omitempty"`
}

// FormulaID is "rank-tracker".
func (RankTrackerParams) FormulaID() string { return "rank-tracker" }

// ReverseSourcingParams are the parameters of Reverse Sourcing (`reverse-sourcing`): Products that sell on Amazon
// for more than Target, Best Buy, Newegg or Walmart charge, with profit and ROI after your fees.
//
// Queries is required. Any other field left nil or empty is not sent, and the formula's default applies; [Ptr]
// makes the pointer of a value.
type ReverseSourcingParams struct {
	// Search words: Searched on Amazon (reverse, without ASINs) or on each retailer (forward), such as lego duplo
	// or air fryer. At most 10 items. Required.
	Queries []string `json:"queries"`

	// Start from: reverse starts from Amazon products and looks for them at the retailers; forward starts from the
	// retailers' search and looks for each product on Amazon. One of "reverse", "forward". Default: "reverse".
	Mode *string `json:"mode,omitempty"`

	// Amazon products: For reverse, the ASINs to source. Leave empty to search Amazon for the words below instead.
	// At most 50 items. Default: [].
	Asins []string `json:"asins,omitempty"`

	// Retailers: Where to buy. Walmart often answers automated visits with a check, so its products may be missing
	// from a run. Each one of "target", "bestbuy", "newegg", "walmart". At least 1 item. Default: ["target",
	// "bestbuy"].
	Retailers []string `json:"retailers,omitempty"`

	// Amazon referral fee (%): The share of the sale price Amazon keeps; 15% in most categories, 8% in some. From
	// 0 to 50. Default: 15.
	ReferralFeePct *float64 `json:"referral_fee_pct,omitempty"`

	// Fulfilment and shipping a unit (USD): Your FBA fee, prep and inbound shipping for one unit. From 0 to 100.
	// Default: 5.
	FulfilmentCost *float64 `json:"fulfilment_cost,omitempty"`

	// Least return on cost (%): Profit as a share of the retailer's price. From -100 to 1000. Default: 15.
	MinRoiPct *float64 `json:"min_roi_pct,omitempty"`

	// Least profit a unit (USD): Profit after the referral fee and fulfilment cost. At least -1000. Default: 2.
	MinProfit *float64 `json:"min_profit,omitempty"`

	// Best Buy order: For forward, the order of Best Buy's results; best_discount reads its deepest reductions
	// first. One of "best_match", "best_discount", "best_selling". Default: "best_match".
	BestbuySort *string `json:"bestbuy_sort,omitempty"`

	// Results compared per search: How many search results are checked against each product. From 1 to 5. Default:
	// 2.
	CandidatesPerRetailer *int `json:"candidates_per_retailer,omitempty"`

	// Amazon products per search: For reverse without ASINs, how many products of each Amazon search are sourced.
	// From 1 to 48. Default: 10.
	AmazonResults *int `json:"amazon_results,omitempty"`

	// Retailer products per search: For forward, how many products of each retailer search are looked up on
	// Amazon. From 1 to 48. Default: 12.
	RetailerResults *int `json:"retailer_results,omitempty"`

	// Most rows: The run stops at this many pairs, best return first. From 1 to 1000. Default: 100.
	MaxResults *int `json:"max_results,omitempty"`
}

// FormulaID is "reverse-sourcing".
func (ReverseSourcingParams) FormulaID() string { return "reverse-sourcing" }

// ReviewMiningParams are the parameters of Review Mining (`review-mining`): The complaints in your products' or
// competitors' reviews on Target, Walmart and Amazon, by topic, with the buyer's own words.
//
// Products is required. Any other field left nil or empty is not sent, and the formula's default applies; [Ptr]
// makes the pointer of a value.
type ReviewMiningParams struct {
	// Product pages: The address of each product's page on Target, Walmart, or an Amazon outside the US
	// (amazon.co.uk, .de, .ca ...). 1 to 30 items. Required.
	Products []string `json:"products"`

	// Topics: Phrases to look for, matched as words and their forms ("broke" also finds "broken"). 1 to 40 items.
	// Default: a list of 20 values, which `get` on the formula shows.
	Topics []string `json:"topics,omitempty"`

	// Most stars: Only reviews of this many stars or fewer. 5 reads every review. From 1 to 5. Default: 3.
	MaxStars *int `json:"max_stars,omitempty"`

	// Reviews read first: Target and Walmart read their reviews in this order; Amazon reads the page's own. One of
	// "lowest_rating", "newest", "most_helpful". Default: "lowest_rating".
	Order *string `json:"order,omitempty"`

	// Only reviews that mention a topic. Default: true.
	OnlyWithTopics *bool `json:"only_with_topics,omitempty"`

	// Most rows: The run stops at this many reviews, the most topics first. From 1 to 2000. Default: 300.
	MaxResults *int `json:"max_results,omitempty"`
}

// FormulaID is "review-mining".
func (ReviewMiningParams) FormulaID() string { return "review-mining" }

// ReviewVelocityParams are the parameters of Review Velocity (`review-velocity`): Ratings and reviews of the apps,
// companies and games you follow, with how many were added since the last run and per day.
//
// Apps is required. Any other field left nil or empty is not sent, and the formula's default applies; [Ptr] makes
// the pointer of a value.
type ReviewVelocityParams struct {
	// Apps: An App Store app's address or its number (368677368 for Uber), or a Google Play app's package name, as
	// its address shows it after id= (com.ubercab). An app on both stores is given twice, once for each. At most
	// 50 items. Required.
	Apps []string `json:"apps"`

	// Companies on Trustpilot: The company's website domain, such as chewy.com, for its Trustpilot reviews. At
	// most 50 items. Default: [].
	Companies []string `json:"companies,omitempty"`

	// Games on Steam: A Steam game's number or its store address, such as 1091500 for Cyberpunk 2077. At most 50
	// items. Default: [].
	Games []string `json:"games,omitempty"`

	// Store country: The country whose App Store, Google Play and Steam store are read, as a two-letter code.
	// Their counts and ratings are per country. Default: "us".
	Country *string `json:"country,omitempty"`

	// Most rows: The run stops at this many rows, the biggest moves first. From 1 to 200. Default: 200.
	MaxResults *int `json:"max_results,omitempty"`
}

// FormulaID is "review-velocity".
func (ReviewVelocityParams) FormulaID() string { return "review-velocity" }

// SalaryBenchmarkParams are the parameters of Salary Benchmark (`salary-benchmark`): The pay employers advertise
// for a job in each place, as yearly ranges per place, with every posting behind them.
//
// JobTitles and Locations are required. Any other field left nil or empty is not sent, and the formula's default
// applies; [Ptr] makes the pointer of a value.
type SalaryBenchmarkParams struct {
	// Job: The job to benchmark, and up to two other names it is advertised under. 1 to 3 items. Required.
	JobTitles []string `json:"job_titles"`

	// Places: Each city, region or country to compare. 1 to 10 items. Required.
	Locations []string `json:"locations"`

	// Posted within (days): Only postings from the last this many days, so the figures are today's. From 1 to 60.
	// Default: 30.
	PostedWithinDays *int `json:"posted_within_days,omitempty"`

	// Job boards: Which boards to read. Reed and Totaljobs cover the UK (Reed states pay on most postings), Dice
	// technology roles in the US, Seek Australia, StepStone Germany. Each one of "indeed", "linkedin", "reed",
	// "totaljobs", "dice", "seek", "stepstone". At least 1 item. Default: ["indeed", "linkedin"].
	Sources []string `json:"sources,omitempty"`

	// Jobs read per search: How many postings to read for each job name, place and board; many do not state pay.
	// From 5 to 200. Default: 50.
	ResultsPerSearch *int `json:"results_per_search,omitempty"`

	// Only titles that name the job: Keeps postings whose title names the job as whole words; boards also return
	// near matches, which would blur the figures. Default: true.
	TitleMustMatch *bool `json:"title_must_match,omitempty"`

	// Leave out agencies' postings: Agencies often advertise a wide range; leave them out to see only what
	// employers state. Default: false.
	ExcludeAgencies *bool `json:"exclude_agencies,omitempty"`

	// Companies to add from Levels.fyi: Company names as Levels.fyi writes them; their pay for the job is added by
	// level, from what people reported. At most 20 items. Default: [].
	LevelsCompanies []string `json:"levels_companies,omitempty"`

	// The job as Levels.fyi names it: Such as Software Engineer or Data Scientist. Empty uses the first job name.
	// Default: "".
	LevelsRole *string `json:"levels_role,omitempty"`

	// Most rows: The run stops at this many rows, best paid first. From 1 to 2000. Default: 200.
	MaxResults *int `json:"max_results,omitempty"`
}

// FormulaID is "salary-benchmark".
func (SalaryBenchmarkParams) FormulaID() string { return "salary-benchmark" }

// SecFilingMonitorParams are the parameters of SEC Filing Monitor (`sec-filing-monitor`): New SEC filings of the
// companies you follow, insider trades, 8-Ks and activist stakes first, each linked to the filing.
//
// Companies is required. Any other field left nil or empty is not sent, and the formula's default applies; [Ptr]
// makes the pointer of a value.
type SecFilingMonitorParams struct {
	// Companies: Tickers (AAPL, BRK-B) or EDGAR's CIK numbers; a fund or a person who files is named by CIK. 1 to
	// 100 items. Required.
	Companies []string `json:"companies"`

	// Form types: The forms to watch, as EDGAR writes them. Leave empty for every form. At most 20 items. Default:
	// a list of 10 values, which `get` on the formula shows.
	Forms []string `json:"forms,omitempty"`

	// Filed within (days): Only filings made in the last this many days. From 1 to 90. Default: 7.
	FiledWithinDays *int `json:"filed_within_days,omitempty"`

	// Only 8-Ks reporting these items: For example 5.02 (officers leaving or joining) or 1.01 (a material
	// agreement). Leave empty for every 8-K. Default: [].
	EightKItems []string `json:"eight_k_items,omitempty"`

	// Include amendments: Amended filings, whose form ends in /A. Default: true.
	IncludeAmendments *bool `json:"include_amendments,omitempty"`

	// Filings read per company: The newest filings read for each company; raise it for a long window on a company
	// that files often. From 1 to 100. Default: 40.
	PerCompany *int `json:"per_company,omitempty"`

	// Most rows: The run stops at this many filings, most significant first. From 1 to 2000. Default: 200.
	MaxResults *int `json:"max_results,omitempty"`
}

// FormulaID is "sec-filing-monitor".
func (SecFilingMonitorParams) FormulaID() string { return "sec-filing-monitor" }

// SellerIdentityParams are the parameters of Seller Identity (`seller-identity`): The legal business behind each
// marketplace seller of your products, with its address, website, domain record and register entry.
//
// Asins is required. Any other field left nil or empty is not sent, and the formula's default applies; [Ptr] makes
// the pointer of a value.
type SellerIdentityParams struct {
	// Your products on Amazon: ASINs whose sellers to identify. Every seller with an offer on them is read. At
	// most 30 items. Required.
	Asins []string `json:"asins"`

	// Amazon seller ids: Sellers to identify directly, by the id in their storefront's address (seller=...). At
	// most 100 items. Default: [].
	SellerIds []string `json:"seller_ids,omitempty"`

	// Walmart sellers: Walmart Marketplace sellers by id or storefront address. At most 100 items. Default: [].
	WalmartSellers []string `json:"walmart_sellers,omitempty"`

	// Amazon marketplace: Which Amazon the products and sellers are on. One of "us", "uk", "de", "fr", "it", "es",
	// "ca". Default: "us".
	Marketplace *string `json:"marketplace,omitempty"`

	// Authorised sellers: Storefront names to leave out, such as your own and your distributors'. At most 500
	// items. Default: ["Amazon.com", "Amazon"].
	AuthorisedSellers []string `json:"authorised_sellers,omitempty"`

	// Find each business's website and its domain record: Looks the business up by name, and reads the
	// registration of the domain it finds. Default: true.
	FindWebsite *bool `json:"find_website,omitempty"`

	// Look UK businesses up at Companies House: For a seller whose address is in the United Kingdom. Default:
	// true.
	CheckCompanyRegister *bool `json:"check_company_register,omitempty"`

	// Most rows: The run stops at this many sellers, the most identified first. From 1 to 500. Default: 100.
	MaxResults *int `json:"max_results,omitempty"`
}

// FormulaID is "seller-identity".
func (SellerIdentityParams) FormulaID() string { return "seller-identity" }

// ShopifyNicheResearchParams are the parameters of Shopify Niche Research (`shopify-niche-research`): The products
// of a niche's Shopify stores, with price, newness and reductions, and similar products on AliExpress with the
// margin.
//
// Stores is required. Any other field left nil or empty is not sent, and the formula's default applies; [Ptr] makes
// the pointer of a value.
type ShopifyNicheResearchParams struct {
	// Shopify stores: The stores of the niche, by domain. 1 to 20 items. Required.
	Stores []string `json:"stores"`

	// Niche words: Keep products whose title, type or tags hold one of these. Empty keeps every product. At most
	// 20 items. Default: [].
	NicheWords []string `json:"niche_words,omitempty"`

	// Collection: A collection to read in each store, by its handle, such as best-sellers. Empty reads the
	// catalogue. Default: "".
	Collection *string `json:"collection,omitempty"`

	// Look for similar products on AliExpress: One search a product. Default: true.
	FindSuppliers *bool `json:"find_suppliers,omitempty"`

	// Products read per store. From 10 to 250. Default: 250.
	ProductsPerStore *int `json:"products_per_store,omitempty"`

	// Most rows: The run stops at this many products, newest and reduced first. Each is one AliExpress search.
	// From 1 to 500. Default: 30.
	MaxResults *int `json:"max_results,omitempty"`
}

// FormulaID is "shopify-niche-research".
func (ShopifyNicheResearchParams) FormulaID() string { return "shopify-niche-research" }

// SocialSentimentParams are the parameters of Social Sentiment (`social-sentiment`): How often the tickers and
// brands you follow are mentioned on social sources, in what tone, and how both moved since the last run.
//
// Terms is required. Any other field left nil or empty is not sent, and the formula's default applies; [Ptr] makes
// the pointer of a value.
type SocialSentimentParams struct {
	// Tickers or brands: Each is searched as you type it and counted where a post names it as a whole word, such
	// as Nvidia or $NVDA. Mastodon is searched by hashtag, so only one-word terms are read there. 1 to 10 items.
	// Required.
	Terms []string `json:"terms"`

	// Sources: Where to look. Telegram, X and Bluesky are read through the channels and accounts named below. Each
	// one of "hackernews", "mastodon", "youtube", "telegram", "x", "bluesky", "stackoverflow". At least 1 item.
	// Default: ["hackernews", "mastodon", "youtube", "telegram", "x", "bluesky"].
	Sources []string `json:"sources,omitempty"`

	// Mentions of the last (days): Only mentions posted in the last this many days are counted. YouTube comments
	// carry no date and are not cut. From 1 to 30. Default: 7.
	WindowDays *int `json:"window_days,omitempty"`

	// Posts read per term on Mastodon and Telegram: The most posts read for one term on Mastodon and in each
	// Telegram channel, newest first. Hacker News and Stack Overflow are read 100 deep. From 5 to 100. Default:
	// 40.
	PerSourceLimit *int `json:"per_source_limit,omitempty"`

	// Telegram channels: Public channels whose posts are searched for each term, by name (WatcherGuru) or address.
	// At most 10 items. Default: ["WatcherGuru", "financialjuice"].
	TelegramChannels []string `json:"telegram_channels,omitempty"`

	// X accounts: Accounts whose latest posts (X shows the newest five of each) are read for mentions. At most 20
	// items. Default: ["DeItaone", "unusual_whales", "zerohedge", "StockMKTNewz", "WatcherGuru"].
	XAccounts []string `json:"x_accounts,omitempty"`

	// Bluesky accounts: Accounts by handle, such as reuters.com, whose latest posts are read for mentions. At most
	// 20 items. Default: ["bloomberg.com", "reuters.com", "wsj.com"].
	BlueskyAccounts []string `json:"bluesky_accounts,omitempty"`

	// Bluesky posts read per account: The newest posts read from each Bluesky account. From 5 to 100. Default: 30.
	PostsPerAccount *int `json:"posts_per_account,omitempty"`

	// YouTube videos per term: How many of the videos a YouTube search for the term finds first have their
	// comments read. From 1 to 5. Default: 2.
	VideosPerTerm *int `json:"videos_per_term,omitempty"`

	// Comments read per video: The newest comments read on each video. From 5 to 100. Default: 20.
	CommentsPerVideo *int `json:"comments_per_video,omitempty"`

	// Mastodon server: The server whose view of a hashtag is read; a large one sees the posts of many others.
	// Default: "mastodon.social".
	MastodonServer *string `json:"mastodon_server,omitempty"`

	// Positive words: A mention that uses one of these as a whole word, and none of the negative ones, counts as
	// positive. At most 100 items. Default: a list of 44 values, which `get` on the formula shows.
	PositiveWords []string `json:"positive_words,omitempty"`

	// Negative words: A mention that uses one of these as a whole word, and none of the positive ones, counts as
	// negative. At most 100 items. Default: a list of 52 values, which `get` on the formula shows.
	NegativeWords []string `json:"negative_words,omitempty"`

	// Most rows: The run stops at this many rows, the biggest moves first. From 1 to 100. Default: 100.
	MaxResults *int `json:"max_results,omitempty"`
}

// FormulaID is "social-sentiment".
func (SocialSentimentParams) FormulaID() string { return "social-sentiment" }

// StaleListingsParams are the parameters of Stale Listings (`stale-listings`): Homes for sale that have sat on
// Zillow for months or been cut in price, with list price, Zestimate and the listing agent.
//
// Locations is required. Any other field left nil or empty is not sent, and the formula's default applies; [Ptr]
// makes the pointer of a value.
type StaleListingsParams struct {
	// Places: Cities with their state ("Austin, TX"), ZIP codes or counties, as Zillow names them. 1 to 20 items.
	// Required.
	Locations []string `json:"locations"`

	// Least days on Zillow: Homes listed at least this long. A home with a price cut is kept sooner when the next
	// option is on. From 1 to 730. Default: 90.
	MinDaysOnMarket *int `json:"min_days_on_market,omitempty"`

	// A price cut makes a listing stale too: Keep a home whose price was cut, however long it has been listed.
	// Default: true.
	PriceCutCounts *bool `json:"price_cut_counts,omitempty"`

	// Kinds of home: Which kinds of home to search. Each one of "house", "condo", "townhouse", "multi_family",
	// "manufactured", "land". Default: ["house", "townhouse", "multi_family"].
	HomeTypes []string `json:"home_types,omitempty"`

	// Least price (USD): 0 for no floor. At least 0. Default: 0.
	MinPrice *int `json:"min_price,omitempty"`

	// Most price (USD): 0 for no ceiling. At least 0. Default: 0.
	MaxPrice *int `json:"max_price,omitempty"`

	// Listings read per place: Zillow's first page of each place is read, up to this many homes. From 5 to 41.
	// Default: 41.
	ListingsPerPlace *int `json:"listings_per_place,omitempty"`

	// Read each listing's agent and history: The listing agent and brokerage, the Zestimate and the price history.
	// One page a listing. Default: true.
	WithAgent *bool `json:"with_agent,omitempty"`

	// Show a listing again after (days): For a saved formula, how long before a delivered listing may come back.
	// From 1 to 365. Default: 30.
	ResurfaceAfterDays *int `json:"resurface_after_days,omitempty"`

	// Most rows: The run stops at this many listings, longest listed first. From 1 to 1000. Default: 100.
	MaxResults *int `json:"max_results,omitempty"`
}

// FormulaID is "stale-listings".
func (StaleListingsParams) FormulaID() string { return "stale-listings" }

// StorefrontResearchParams are the parameters of Storefront Research (`storefront-research`): What a competing
// Amazon seller sells in your niche, its prices and Buy Box, the sellers it competes with, and where to source it.
//
// SellerIds and Queries are required. Any other field left nil or empty is not sent, and the formula's default
// applies; [Ptr] makes the pointer of a value.
type StorefrontResearchParams struct {
	// Competitors' seller ids: The id in a storefront's address (seller=...). 1 to 10 items. Required.
	SellerIds []string `json:"seller_ids"`

	// Niche search words: Searched on Amazon; the listings found are checked for the competitor's offers. At most
	// 10 items. Required.
	Queries []string `json:"queries"`

	// Best-seller lists: Also check the products of these Best Sellers lists, such as toys-and-games. At most 5
	// items. Default: [].
	Categories []string `json:"categories,omitempty"`

	// Listings per search: Each listing is one read of its offers. From 1 to 48. Default: 20.
	ResultsPerSearch *int `json:"results_per_search,omitempty"`

	// Look for each product at Target and Best Buy. Default: true.
	SourceAtRetailers *bool `json:"source_at_retailers,omitempty"`

	// Amazon referral fee (%). From 0 to 50. Default: 15.
	ReferralFeePct *float64 `json:"referral_fee_pct,omitempty"`

	// Fulfilment and shipping a unit (USD). From 0 to 100. Default: 5.
	FulfilmentCost *float64 `json:"fulfilment_cost,omitempty"`

	// Most rows: The run stops at this many products, Buy Box first. From 1 to 500. Default: 100.
	MaxResults *int `json:"max_results,omitempty"`
}

// FormulaID is "storefront-research".
func (StorefrontResearchParams) FormulaID() string { return "storefront-research" }

// StrUnderwritingParams are the parameters of Short-Term Rental Comps (`str-underwriting`): Airbnb and Vrbo rentals
// in an area for your dates, with nightly rate, rating, reviews and the area's average rate.
//
// Location and Checkin and Checkout are required. Any other field left nil or empty is not sent, and the formula's
// default applies; [Ptr] makes the pointer of a value.
type StrUnderwritingParams struct {
	// Place: Where, as you would type it into Airbnb, with the country, such as "Nashville, Tennessee, United
	// States". Required.
	Location string `json:"location"`

	// Check-in: The first night, YYYY-MM-DD. Rates are for this stay. Required.
	Checkin string `json:"checkin"`

	// Check-out: The day of departure, YYYY-MM-DD. Required.
	Checkout string `json:"checkout"`

	// Guests: Adults staying; rentals too small for them are not offered. From 1 to 16. Default: 2.
	Guests *int `json:"guests,omitempty"`

	// Platforms: Where to look. Each one of "airbnb", "vrbo". At least 1 item. Default: ["airbnb", "vrbo"].
	Platforms []string `json:"platforms,omitempty"`

	// Airbnb currency: The currency Airbnb shows prices in. Vrbo shows the destination's own. One of "USD", "EUR",
	// "GBP", "CAD", "AUD". Default: "USD".
	Currency *string `json:"currency,omitempty"`

	// Least reviews: Leave out rentals with fewer reviews, such as new listings. At least 0. Default: 0.
	MinReviews *int `json:"min_reviews,omitempty"`

	// Most rows: The run stops at this many rentals, most reviewed first. From 1 to 500. Default: 150.
	MaxResults *int `json:"max_results,omitempty"`
}

// FormulaID is "str-underwriting".
func (StrUnderwritingParams) FormulaID() string { return "str-underwriting" }

// SupplierDiscoveryParams are the parameters of Supplier Discovery (`supplier-discovery`): Suppliers of a product
// from IndiaMART, Google Maps, Gelbe Seiten and ThomasNet, each with location, website and email.
//
// Products is required. Any other field left nil or empty is not sent, and the formula's default applies; [Ptr]
// makes the pointer of a value.
type SupplierDiscoveryParams struct {
	// Products: What you buy, as a supplier would list it, such as CNC machined parts or corrugated boxes. Each is
	// searched on its own. 1 to 5 items. Required.
	Products []string `json:"products"`

	// Places: Towns or regions to find suppliers in on Google Maps and Gelbe Seiten, such as Cleveland, OH or
	// Stuttgart. IndiaMART is searched across India whatever is given here. At most 5 items. Default: [].
	Places []string `json:"places,omitempty"`

	// Directories: IndiaMART (India), Google Maps (anywhere, by place), Gelbe Seiten (Germany, by town) and
	// ThomasNet (the profiles you give). Each one of "indiamart", "google_maps", "gelbeseiten", "thomasnet". At
	// least 1 item. Default: ["indiamart", "google_maps"].
	Sources []string `json:"sources,omitempty"`

	// ThomasNet profiles: The links or ids of ThomasNet supplier profiles you have shortlisted, read in full when
	// ThomasNet is chosen. At most 50 items. Default: [].
	ThomasnetProfiles []string `json:"thomasnet_profiles,omitempty"`

	// Map listings per search: How many Google Maps listings are read for each product and place. From 5 to 60.
	// Default: 20.
	PerSearch *int `json:"per_search,omitempty"`

	// Only verified IndiaMART sellers: Keeps the IndiaMART sellers it shows as verified or holding its TrustSEAL.
	// Default: false.
	VerifiedOnly *bool `json:"verified_only,omitempty"`

	// Lowest rating: Leaves out suppliers rated lower by their buyers or reviewers. A supplier with no rating is
	// kept. From 0 to 5. Default: 0.
	MinRating *float64 `json:"min_rating,omitempty"`

	// Suppliers to leave out: Names, such as the suppliers you already buy from. At most 1000 items. Default: [].
	ExcludeCompanies []string `json:"exclude_companies,omitempty"`

	// Websites to leave out: Domains, such as acme.com, of the suppliers you already buy from. At most 1000 items.
	// Default: [].
	ExcludeDomains []string `json:"exclude_domains,omitempty"`

	// Read each IndiaMART product page: Adds the supplier's town and business type, such as manufacturer or
	// exporter. One unit for each IndiaMART supplier. Default: true.
	IndiamartDetails *bool `json:"indiamart_details,omitempty"`

	// Only suppliers with an email: Keeps the suppliers whose own website publishes an email address that passes
	// the check. Default: false.
	RequireEmail *bool `json:"require_email,omitempty"`

	// Most rows: The run stops at this many suppliers. From 1 to 1000. Default: 100.
	MaxResults *int `json:"max_results,omitempty"`
}

// FormulaID is "supplier-discovery".
func (SupplierDiscoveryParams) FormulaID() string { return "supplier-discovery" }

// TalentSourcingParams are the parameters of Talent Sourcing (`talent-sourcing`): People who posted that they want
// to be hired, or show their work, in your role and place, scored on the skills you need.
//
// JobTitle and Location are required. Any other field left nil or empty is not sent, and the formula's default
// applies; [Ptr] makes the pointer of a value.
type TalentSourcingParams struct {
	// Job title: The role you are hiring for, as people would write it about themselves. Required.
	JobTitle string `json:"job_title"`

	// Location: The city, region or country the people should be in. Write Remote for anywhere. Required.
	Location string `json:"location"`

	// Skills: What the role asks for. They score each person, and pick the GitHub projects to look through. At
	// most 10 items. Default: [].
	Skills []string `json:"skills,omitempty"`

	// Where to look: GitHub finds only developers whose profile gives the place, and is searched only when skills
	// are given. Each one of "hackernews", "github". At least 1 item. Default: ["hackernews"].
	Sources []string `json:"sources,omitempty"`

	// Read more posts: Reads the newest 300 posts that name the place instead of 200, reaching further back.
	// Default: false.
	MoreResults *bool `json:"more_results,omitempty"`

	// Posted within (days): Only posts from the last this many days; an older post may no longer be true. From 7
	// to 730. Default: 365.
	PostedWithinDays *int `json:"posted_within_days,omitempty"`

	// GitHub profiles to read: How many developers' profiles to open, from the projects that use the skills. From
	// 1 to 50. Default: 15.
	GithubProfiles *int `json:"github_profiles,omitempty"`

	// Least skills named: Leave out a person who names fewer of the skills than this. From 0 to 10. Default: 0.
	MinSkills *int `json:"min_skills,omitempty"`

	// Show a person again after (days): For a saved formula, how long before a delivered person may come back.
	// From 1 to 365. Default: 90.
	ResurfaceAfterDays *int `json:"resurface_after_days,omitempty"`

	// Most rows: The run stops at this many people, best first. From 1 to 500. Default: 50.
	MaxResults *int `json:"max_results,omitempty"`
}

// FormulaID is "talent-sourcing".
func (TalentSourcingParams) FormulaID() string { return "talent-sourcing" }

// TenderAlertsParams are the parameters of Tender Alerts (`tender-alerts`): New public tenders that match your
// words and codes, from the EU, UK and US, with buyer, deadline, value and link.
//
// Any other field left nil or empty is not sent, and the formula's default applies; [Ptr] makes the pointer of a
// value.
type TenderAlertsParams struct {
	// Words: What you sell, as a buyer would write it, such as cybersecurity or road resurfacing. Each is searched
	// on its own. Leave empty to take every notice of the window, or to match on CPV codes alone. At most 10
	// items. Default: [].
	Keywords []string `json:"keywords,omitempty"`

	// Notice boards: TED (European Union), Find a Tender and Contracts Finder (United Kingdom), Grants.gov, New
	// York City's City Record, Chicago's and the MTA's contracts (awards only), and SAM.gov, which is read where
	// it is offered. Each one of "ted", "find_a_tender", "contracts_finder", "grants_gov", "new_york_city",
	// "chicago", "ny_mta", "sam_gov". At least 1 item. Default: ["ted", "find_a_tender", "contracts_finder",
	// "grants_gov", "new_york_city", "sam_gov"].
	Sources []string `json:"sources,omitempty"`

	// Notice stages: Pre-RFP notices (prior information, sources sought, presolicitation, forecasted grants and UK
	// notices of planned procurements), open tenders, or awards. Each one of "planning", "tender", "award". At
	// least 1 item. Default: ["tender"].
	NoticeTypes []string `json:"notice_types,omitempty"`

	// Published within (days): Only notices published in the last this many days, today included. Run daily, 3
	// covers a weekend. From 1 to 30. Default: 3.
	PublishedWithinDays *int `json:"published_within_days,omitempty"`

	// CPV codes: Eight-digit CPV codes, such as 72000000 for IT services. TED is searched by them; a UK notice is
	// kept when its main code is one of them. A notice matching your words or your codes is kept. At most 10
	// items. Default: [].
	CpvCodes []string `json:"cpv_codes,omitempty"`

	// Buyers' countries: Three-letter codes, such as DEU, FRA, GBR or USA. Leave empty for every country of the
	// boards you chose. At most 30 items. Default: [].
	Countries []string `json:"countries,omitempty"`

	// NAICS codes (SAM.gov): Only SAM.gov notices in these industry codes, 2 to 6 digits, such as 541512. At most
	// 10 items. Default: [].
	NaicsCodes []string `json:"naics_codes,omitempty"`

	// Set-asides (SAM.gov): Only SAM.gov notices set aside for these, by SAM.gov's codes, such as SBA, 8A, HZC,
	// SDVOSBC or WOSB. At most 10 items. Default: [].
	SetAsides []string `json:"set_asides,omitempty"`

	// Places of performance (SAM.gov): Only SAM.gov notices for work in these US states, as two letters. At most
	// 10 items. Default: [].
	USStates []string `json:"us_states,omitempty"`

	// At least this many days to respond: Leaves out notices due sooner. A notice past its deadline is always left
	// out; one that states no deadline is kept. From 0 to 365. Default: 0.
	MinDaysLeft *int `json:"min_days_left,omitempty"`

	// Due within (days): Only notices due within this many days. 0 keeps every deadline, and notices that state
	// none. From 0 to 730. Default: 0.
	DueWithinDays *int `json:"due_within_days,omitempty"`

	// Leave out titles with: Notices whose title holds one of these words are left out, such as cleaning or
	// catering. At most 30 items. Default: [].
	ExcludeWords []string `json:"exclude_words,omitempty"`

	// Read each grant in full: Adds a Grants.gov opportunity's funding, largest award, eligibility and
	// description. One more unit for each grant. Default: true.
	GrantDetails *bool `json:"grant_details,omitempty"`

	// Most rows: The run stops at this many notices, best matches first. From 1 to 2000. Default: 200.
	MaxResults *int `json:"max_results,omitempty"`
}

// FormulaID is "tender-alerts".
func (TenderAlertsParams) FormulaID() string { return "tender-alerts" }

// TenderChangeWatchParams are the parameters of Tender Change Watch (`tender-change-watch`): The tenders, grants
// and awards you watch, read again each run, with a row only when a deadline, status or value changed.
//
// Watch is required. Any other field left nil or empty is not sent, and the formula's default applies; [Ptr] makes
// the pointer of a value.
type TenderChangeWatchParams struct {
	// Notices to watch: A Grants.gov opportunity's id or link, a USAspending.gov award's id or link, or a UK
	// procurement's Open Contracting id (ocds-…, shown on its Find a Tender or Contracts Finder notice). Anything
	// else is left out. 1 to 150 items. Required.
	Watch []string `json:"watch"`

	// UK notices read (days): For UK procurements, how many days of new notices are read for their updates. Run
	// daily, 3 is enough. From 1 to 30. Default: 7.
	UKDays *int `json:"uk_days,omitempty"`

	// Most rows: The run stops at this many changed notices. From 1 to 1000. Default: 200.
	MaxResults *int `json:"max_results,omitempty"`
}

// FormulaID is "tender-change-watch".
func (TenderChangeWatchParams) FormulaID() string { return "tender-change-watch" }

// UnlinkedMentionsParams are the parameters of Unlinked Mentions (`unlinked-mentions`): Pages that name your brand
// but do not link to your site, each read to be sure, with the email the site publishes.
//
// Brand and BrandDomain are required. Any other field left nil or empty is not sent, and the formula's default
// applies; [Ptr] makes the pointer of a value.
type UnlinkedMentionsParams struct {
	// Brand name: The name as people write it. A name that is also a common word finds pages using the word as
	// well. Required.
	Brand string `json:"brand"`

	// Your domain: The domain a link should lead to, such as example.com. Pages on it are left out. Required.
	BrandDomain string `json:"brand_domain"`

	// Searches: Words searched with the name, each its own search; an empty one searches the name alone. 1 to 10
	// items. Default: ["", "review", "alternatives", "pricing", "integration"].
	SearchTerms []string `json:"search_terms,omitempty"`

	// Include news stories: Also reads the stories the news search lists for the name. Default: true.
	IncludeNews *bool `json:"include_news,omitempty"`

	// Your other domains: Other domains of yours; a link to any of them counts as a link, and pages on them are
	// left out. At most 20 items. Default: [].
	OtherDomains []string `json:"other_domains,omitempty"`

	// Pages to read: How many of the pages found are read, news first, then the best placed in search. From 5 to
	// 150. Default: 40.
	PagesToRead *int `json:"pages_to_read,omitempty"`

	// Results pages per search: How many pages of results to read for each search, about ten results a page. From
	// 1 to 3. Default: 1.
	SearchPages *int `json:"search_pages,omitempty"`

	// Sites to leave out: Domains never shown, such as social networks and review sites where no editor adds
	// links. At most 2000 items. Default: a list of 22 values, which `get` on the formula shows.
	ExcludeDomains []string `json:"exclude_domains,omitempty"`

	// Search country: The two-letter country whose search results are read, such as us, gb or de. Default: "us".
	Country *string `json:"country,omitempty"`

	// Search language: The language of the results, such as en or de. Default: "en".
	Language *string `json:"language,omitempty"`

	// Only sites with an email: Keeps the pages whose site publishes an email address that passes the check.
	// Default: false.
	RequireEmail *bool `json:"require_email,omitempty"`

	// Show a page again after (days): For a saved formula, how long before a delivered page may come back. From 1
	// to 365. Default: 180.
	ResurfaceAfterDays *int `json:"resurface_after_days,omitempty"`

	// Most rows: The run stops at this many pages, news and pages naming the brand in their title first. From 1 to
	// 1000. Default: 100.
	MaxResults *int `json:"max_results,omitempty"`
}

// FormulaID is "unlinked-mentions".
func (UnlinkedMentionsParams) FormulaID() string { return "unlinked-mentions" }

// UsedCarArbitrageParams are the parameters of Used-Car Arbitrage (`used-car-arbitrage`): Copart lots selling well
// below what the same make, model and year fetch at retail near you, ranked by the spread, with bid moves.
//
// Make and Models and Zip are required. Any other field left nil or empty is not sent, and the formula's default
// applies; [Ptr] makes the pointer of a value.
type UsedCarArbitrageParams struct {
	// Make: The make, such as Toyota or Ford. Required.
	Make string `json:"make"`

	// Models: The models of that make, named as Copart names them, such as Camry or RAV4. 1 to 10 items. Required.
	Models []string `json:"models"`

	// Your ZIP code: Where you would sell the cars. Retail comparables are the listings around it. Required.
	Zip string `json:"zip"`

	// Model years: The model years to search. Each year of each model is one search for retail comparables.
	// Default: {"min": 2018, "max": 2022}.
	Years *IntRange `json:"years,omitempty"`

	// Lots to watch: Copart lot numbers or lot pages you follow, read on every run whether or not a search finds
	// them. At most 50 items. Default: [].
	WatchLots []string `json:"watch_lots,omitempty"`

	// Yard states: Two-letter codes of the states whose yards you buy from, such as TX. Leave empty for all. At
	// most 60 items. Default: [].
	States []string `json:"states,omitempty"`

	// Most miles: Lots with more miles on the odometer are left out. 0 for any mileage. From 0 to 500000. Default:
	// 0.
	MaxOdometer *int `json:"max_odometer,omitempty"`

	// Leave out titles: Words of the titles to leave out. By default the cars that can never be titled for the
	// road again, which no retail price fits. At most 20 items. Default: ["NON-REPAIRABLE", "DESTRUCTION",
	// "SCRAP", "JUNK", "PARTS ONLY"].
	ExcludeTitles []string `json:"exclude_titles,omitempty"`

	// Clean titles only: Leaves out salvage and other branded titles. Default: false.
	CleanTitleOnly *bool `json:"clean_title_only,omitempty"`

	// Runs and drives only: Keeps the lots Copart marks as running and driving. Default: false.
	RunsAndDrivesOnly *bool `json:"runs_and_drives_only,omitempty"`

	// Leave out damage: Words of the main damage to leave out, such as FLOOD, BURN or ALL OVER. At most 20 items.
	// Default: [].
	ExcludeDamage []string `json:"exclude_damage,omitempty"`

	// Least current bid ($): Leaves out lots whose bidding has barely begun, whose spread on a first small bid
	// says little. A lot with a buy-it-now price is kept. 0 keeps every lot. From 0 to 100000. Default: 1000.
	MinBid *int `json:"min_bid,omitempty"`

	// A bid that counts (% of retail): A lot whose current bid is under this share of its retail median ranks
	// after the others, since bidding usually goes well past it. A buy-it-now price always counts. From 0 to 90.
	// Default: 25.
	SeriousBidPct *int `json:"serious_bid_pct,omitempty"`

	// Include lots selling today: Lots whose sale is on the day of the run, which may be over when the rows
	// arrive. Days are whole calendar days (UTC). A lot you watch is kept on its sale day. Default: false.
	IncludeSaleDay *bool `json:"include_sale_day,omitempty"`

	// Lots read per model: How many of Copart's lots to read for each model, the soonest sales first. One of 20,
	// 50, 100. Default: 100.
	LotsPerModel *int `json:"lots_per_model,omitempty"`

	// Least spread ($): Lots whose retail median is less than this above their price are left out. From 0 to
	// 100000. Default: 2000.
	MinSpread *int `json:"min_spread,omitempty"`

	// Retail comparables from: The sites whose dealer listings price each lot. A car listed on both counts once.
	// Each one of "cargurus", "autotrader". At least 1 item. Default: ["cargurus", "autotrader"].
	CompSources []string `json:"comp_sources,omitempty"`

	// Comparables within (miles): How far from your ZIP code a retail listing may be. One of 25, 50, 75, 100, 200,
	// 500. Default: 100.
	Radius *int `json:"radius,omitempty"`

	// Similar mileage (miles either way): Comparables within this many miles of the lot's odometer are preferred.
	// 0 compares every mileage. From 0 to 100000. Default: 30000.
	MileageBand *int `json:"mileage_band,omitempty"`

	// Price on any mileage when few are similar: When too few comparables have a similar mileage, price the lot on
	// every mileage instead of leaving it out. A lot with far more miles than most listings then looks better than
	// it is. Default: false.
	AnyMileageFallback *bool `json:"any_mileage_fallback,omitempty"`

	// Least comparables: A lot priced on fewer retail listings than this is left out. From 1 to 20. Default: 3.
	MinComps *int `json:"min_comps,omitempty"`

	// Comparables searches: The most retail searches for each site, one for each make, model and year found. From
	// 1 to 100. Default: 40.
	CompSearches *int `json:"comp_searches,omitempty"`

	// Most rows: The run stops at this many lots, widest spread first. From 1 to 2000. Default: 200.
	MaxResults *int `json:"max_results,omitempty"`
}

// FormulaID is "used-car-arbitrage".
func (UsedCarArbitrageParams) FormulaID() string { return "used-car-arbitrage" }

// VcDealSourcingParams are the parameters of VC Deal Sourcing (`vc-deal-sourcing`): Young startups gaining traction
// on GitHub, Hacker News and Y Combinator, one row each with its signals and website.
//
// Any other field left nil or empty is not sent, and the formula's default applies; [Ptr] makes the pointer of a
// value.
type VcDealSourcingParams struct {
	// Topics: What you invest in, as words a startup would describe itself with (agents, database, robotics).
	// Leave empty for everything. At most 10 items. Default: [].
	Topics []string `json:"topics,omitempty"`

	// Within the last (days): How new a repository, a Show HN or a Y Combinator launch may be. From 1 to 90.
	// Default: 30.
	WithinDays *int `json:"within_days,omitempty"`

	// Where to look: GitHub's new and trending repositories, Show HN on Hacker News, and Y Combinator's directory.
	// Each one of "github", "hacker_news", "ycombinator". At least 1 item. Default: ["github", "hacker_news",
	// "ycombinator"].
	Sources []string `json:"sources,omitempty"`

	// Fewest GitHub stars: A repository with fewer stars is left out. From 0 to 100000. Default: 200.
	MinStars *int `json:"min_stars,omitempty"`

	// Fewest Show HN points: A Show HN with fewer points is left out. From 0 to 10000. Default: 50.
	MinHnPoints *int `json:"min_hn_points,omitempty"`

	// Leave out: GitHub accounts, domains or company names to leave out; by default, large companies whose
	// repositories trend. At most 1000 items. Default: a list of 47 values, which `get` on the formula shows.
	Exclude []string `json:"exclude,omitempty"`

	// Results read per search: How many repositories, Show HN posts and Y Combinator companies to read for each
	// topic. From 5 to 100. Default: 30.
	ResultsPerSearch *int `json:"results_per_search,omitempty"`

	// Most rows: The run stops at this many startups, strongest signals first. From 1 to 2000. Default: 200.
	MaxResults *int `json:"max_results,omitempty"`
}

// FormulaID is "vc-deal-sourcing".
func (VcDealSourcingParams) FormulaID() string { return "vc-deal-sourcing" }

// WebRankTrackerParams are the parameters of Web Rank Tracker (`web-rank-tracker`): Where the websites you follow
// rank in Similarweb's most-visited rankings, by country and category, and how they moved since last run.
//
// Domains is required. Any other field left nil or empty is not sent, and the formula's default applies; [Ptr]
// makes the pointer of a value.
type WebRankTrackerParams struct {
	// Websites: The domains to follow, such as etsy.com. 1 to 30 items. Required.
	Domains []string `json:"domains"`

	// Countries: The countries whose rankings to read, or worldwide. Each one of "worldwide", "US", "GB", "DE",
	// "FR", "CA", "AU", "JP", "IN", "BR", "MX", "IT", "ES", "NL", "SE", "CH", "KR", "SG", "AE", "SA", "ZA", "TR",
	// "PL", "ID", "PH", "TH", "VN", "IL", "IE", "NZ". 1 to 4 items. Default: ["US"].
	Countries []string `json:"countries,omitempty"`

	// Categories: The rankings to read in each country. All is the ranking of every category together, 50 places
	// deep; a category's ranking shows 50 places worldwide and 5 in one country. Each one of "all",
	// "ai-chatbots-and-tools", "arts-and-entertainment", "business-and-consumer-services",
	// "community-and-society", "computers-electronics-and-technology", "crypto-and-blockchain",
	// "e-commerce-and-shopping", "finance", "food-and-drink", "gambling", "games", "health",
	// "heavy-industry-and-engineering", "hobbies-and-leisure", "home-and-garden", "jobs-and-career",
	// "law-and-government", "lifestyle", "news-and-media", "pets-and-animals", "reference-materials",
	// "science-and-education", "sports", "travel-and-tourism", "vehicles". 1 to 4 items. Default: ["all"].
	Categories []string `json:"categories,omitempty"`

	// Most rows: The run stops at this many rows, the biggest moves first. From 1 to 500. Default: 500.
	MaxResults *int `json:"max_results,omitempty"`
}

// FormulaID is "web-rank-tracker".
func (WebRankTrackerParams) FormulaID() string { return "web-rank-tracker" }

// WholesaleCatalogCheckParams are the parameters of Wholesale Catalogue Check (`wholesale-catalog-check`): Your
// supplier's price list checked against Amazon line by line, with Amazon's price, reviews and your margin after
// fees.
//
// Catalog is required. Any other field left nil or empty is not sent, and the formula's default applies; [Ptr]
// makes the pointer of a value.
type WholesaleCatalogCheckParams struct {
	// Supplier catalogue: One line a product, from the supplier's sheet - its UPC or EAN (code), your cost, the
	// name, and the case pack when the cost is for a case. 1 to 500 items. Required.
	Catalog []map[string]any `json:"catalog"`

	// Amazon marketplace: Where you sell. One of "us", "uk", "de", "it", "es", "ca". Default: "us".
	Marketplace *string `json:"marketplace,omitempty"`

	// Currency of your costs: Costs are compared only with Amazon prices in this currency. One of "USD", "GBP",
	// "EUR", "CAD". Default: "USD".
	CostCurrency *string `json:"cost_currency,omitempty"`

	// Least return on cost (%): Lines below this are left out. 0 keeps every line that makes money; -100 keeps
	// every line found. From -100 to 1000. Default: 0.
	MinRoiPct *float64 `json:"min_roi_pct,omitempty"`

	// Amazon referral fee (%): The share of the sale price Amazon keeps. From 0 to 50. Default: 15.
	ReferralFeePct *float64 `json:"referral_fee_pct,omitempty"`

	// Fulfilment and shipping a unit: Your FBA fee, prep and inbound shipping for one unit, in the currency of
	// your costs. From 0 to 100. Default: 4.
	FulfilmentCost *float64 `json:"fulfilment_cost,omitempty"`

	// Read each listing's page: Its Buy Box price, whether Amazon sells it, and stock. One page a line found.
	// Default: true.
	CheckOffer *bool `json:"check_offer,omitempty"`

	// Most rows: The run stops at this many lines, best return first. From 1 to 500. Default: 500.
	MaxResults *int `json:"max_results,omitempty"`
}

// FormulaID is "wholesale-catalog-check".
func (WholesaleCatalogCheckParams) FormulaID() string { return "wholesale-catalog-check" }
