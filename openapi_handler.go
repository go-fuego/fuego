package fuego

import "net/http"

func DefaultOpenAPIHandler(specURL string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_, _ = w.Write([]byte(DefaultOpenAPIHTML(specURL)))
	})
}

func DefaultOpenAPIHTML(specURL string) string {
	return `<!doctype html>
<html lang="en">
<head>
	<meta charset="utf-8" />
	<meta name="referrer" content="same-origin" />
	<meta name="viewport" content="width=device-width, initial-scale=1, shrink-to-fit=no">
	<link rel="icon" type="image/svg+xml" href="https://go-fuego.github.io/fuego/img/logo.svg">
	<title>OpenAPI specification</title>
	<script src="https://unpkg.com/@stoplight/elements/web-components.min.js"></script>
	<link rel="stylesheet" href="https://unpkg.com/@stoplight/elements/styles.min.css" />
</head>
<body style="height: 100vh;">
	<elements-api
		apiDescriptionUrl="` + specURL + `"
		layout="responsive"
		router="hash"
		logo="https://go-fuego.github.io/fuego/img/logo.svg"
		tryItCredentialsPolicy="same-origin"
	/>
</body>
</html>`
}
