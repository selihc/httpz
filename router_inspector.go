package httpz

// A RouterInspector is essentially used to walk all registered routes
// and inspect their data to generate documentation and/or clients
// from the results
type RouterInspector interface {
	Inspect()
}
