package main

var EntityWords = map[string]bool{
	// People / Accounts
	"User":     true,
	"Account":  true,
	"Customer": true,
	"Client":   true,
	"Member":   true,
	"Admin":    true,
	"Owner":    true,
	"Author":   true,

	// Business Objects
	"Order":        true,
	"Invoice":      true,
	"Payment":      true,
	"Transaction":  true,
	"Product":      true,
	"Item":         true,
	"Cart":         true,
	"Subscription": true,

	// Auth / Security
	"Session":    true,
	"Token":      true,
	"Credential": true,
	"Profile":    true,
	"Role":       true,
	"Permission": true,

	// System / Resources
	"File":    true,
	"Image":   true,
	"Message": true,
	"Event":   true,
	"Record":  true,
	"Log":     true,

	// Generic objects
	"Data":    true,
	"Info":    true,
	"Config":  true,
	"Setting": true,
}

// Use when need to suggestion.
var AttributeWords = map[string]bool{
	// Identification
	"ID":   true,
	"Uuid": true,
	"Guid": true,
	"Key":  true,
	"Code": true,

	// Human labels
	"Name":     true,
	"Title":    true,
	"Label":    true,
	"Alias":    true,
	"Nickname": true,

	// Contact / Profile
	"Email":    true,
	"Phone":    true,
	"Username": true,
	"Password": true,
	"Avatar":   true,

	// Status / Meta
	"Status":   true,
	"State":    true,
	"Type":     true,
	"Category": true,
	"Level":    true,
	"Rank":     true,
	"Score":    true,
	"Priority": true,

	// Time-based attributes
	"Date":      true,
	"Time":      true,
	"Timestamp": true,
	"CreatedAt": true,
	"UpdatedAt": true,
	"ExpiresAt": true,

	// Countable numbers
	"Count":    true,
	"Number":   true,
	"Quantity": true,
	"Total":    true,

	// Data-like attributes
	"Value":  true,
	"Amount": true,
	"Size":   true,
	"Length": true,
	"Hash":   true,
}
