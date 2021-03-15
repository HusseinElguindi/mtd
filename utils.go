package mtd

import "net/http"

func setHeaders(req *http.Request, headers http.Header) {
	if headers == nil {
		return
	}

	for k, v := range headers {
		req.Header.Set(k, v[0])
	}
}
