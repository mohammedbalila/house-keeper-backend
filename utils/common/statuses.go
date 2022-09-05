package common

var CREATED = "created"
var PENDING = "pending"
var APPROVED = "approved"
var REJECTED = "rejected"

// Statuses is a map of all the purchase subscription statuses
var Statuses = map[string]int64{
	CREATED:  1,
	PENDING:  2,
	APPROVED: 3,
	REJECTED: 4,
}
