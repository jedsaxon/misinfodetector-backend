package validation

import "net/url"

// ValidateGetPostsRequest validates the page number and result amount. It finds "pageNumber"
// and "resultAmount' from the URL values, and returns them as numbers. This function returns
// a map of the queries, to their errors. This means that you should perform your error checking
// with `if len(errors) > 0 { log.Fatalf("errors ...") }`
func ValidateGetPostsRequest(query url.Values) (int64, int64, map[string]string) {
	errors := make(map[string]string)

	pageNumber, err := ValidatePageNumber(query.Get("pageNumber"))
	if err != nil {
		errors["pageNumber"] = err.Error()
	}

	resultAmount, err := ValidateResultAmount(query.Get("resultAmount"))
	if err != nil {
		errors["resultAmount"] = err.Error()
	}

	return pageNumber, resultAmount, errors
}
