package models

import "regexp"

// predefined category rules using regex patterns
var CategoryRules = map[string]*regexp.Regexp{
	"Groceries":          regexp.MustCompile(`(?i)(kaufland|carrefour|lidl|profi|mega|supermarket|walmart|tesco|sainsbury|asda|aldi|grocery|market|food\sstore)`),
	"Dining":             regexp.MustCompile(`(?i)(kfc|mcdo|mcdonald|pizza|restaurant|bistro|grill|burger|cafe|coffee|bar|pub|lunch|dinner|takeout|fast\s?food|steakhouse|diner)`),
	"Transportation":     regexp.MustCompile(`(?i)(uber|bolt|taxi|lyft|stb|metrorex|transport|gas|petrol|fuel|parking|train|bus|metro|subway|tram|airplane|flight|aviation|toll|highway)`),
	"Utilities":          regexp.MustCompile(`(?i)(enel|engie|electrica|utility|utilities|electricity|power|gas|water|internet|fiber|broadband|wifi|phone\s?bill|mobile\s?bill|telecom|cable\s?tv|bill|invoice)`),
	"Healthcare":         regexp.MustCompile(`(?i)(pharmacy|doctor|hospital|medical|clinic|health|dentist|medicine|prescription|drugstore|optical|therapy|physio|dermatology|cardiology)`),
	"Entertainment":      regexp.MustCompile(`(?i)(cinema|movie|theater|concert|festival|spotify|netflix|hbo|disney|prime\s?video|games|gaming|steam|playstation|xbox|music|show|event|subscription)`),
	"Shopping":           regexp.MustCompile(`(?i)(amazon|mall|shop|store|fashion|clothing|shoes|nike|adidas|zara|hm|h&m|retail|electronics|tech|appliance|ikea|decathlon|walmart|target)`),
	"Fitness":            regexp.MustCompile(`(?i)(gym|fitness|sport|yoga|swimming|pool|training|exercise|athlete|workout|membership|fitbit|protein|supplement)`),
	"Education":          regexp.MustCompile(`(?i)(school|university|college|course|book|tuition|education|library|learning|udemy|coursera|edx|training|seminar|workshop)`),
	"Travel":             regexp.MustCompile(`(?i)(hotel|airbnb|booking|accommodation|flight|airline|travel|vacation|holiday|trip|tour|hostel|expedia|agoda)`),
	"Insurance":          regexp.MustCompile(`(?i)(insurance|premium|coverage|policy|health\sinsurance|life\sinsurance|car\sinsurance|home\sinsurance|deductible)`),
	"Banking":            regexp.MustCompile(`(?i)(bank|fee|transfer|wire|payment|charge|commission|atm|withdrawal|deposit|overdraft|card\sfee)`),
	"Home & Maintenance": regexp.MustCompile(`(?i)(furniture|repair|maintenance|home\simprovement|tools|plumbing|electrician|cleaning|appliances|hardware\store)`),
	"Subscriptions":      regexp.MustCompile(`(?i)(subscription|monthly|membership|plan|renewal|premium|software|license|saas|cloud|storage)`),
	"Pets":               regexp.MustCompile(`(?i)(pet|vet|veterinary|dog|cat|petfood|grooming|petshop)`),
	"Gifts & Donations":  regexp.MustCompile(`(?i)(gift|present|donation|charity|ngo|fundraiser|tip)`),
	"Salary":             regexp.MustCompile(`(?i)(salary|paycheck|wage|income|monthly\sincome|payroll|compensation|earning|bonus|commission\sincome)`),
	"Other Income":       regexp.MustCompile(`(?i)(refund|rebate|cashback|interest|dividend|royalty|payout|income)`),
	"Other":              regexp.MustCompile(`.*`),
}
