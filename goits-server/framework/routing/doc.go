// Package routing contains the api routing engine and api handlers of common crud functionality.
//
// There are three main components: engine, resource and handlers.
// The part engine contains the third party dependent code and isolates these code as much as possible.
// The part resource just a wrapper or a semi abstraction representing a rest resource.
// The part handlers contain actual rest api handlers. Handlers ideally should not contain any third party dependent code
// and use third party code via abstractions provided.
package routing
