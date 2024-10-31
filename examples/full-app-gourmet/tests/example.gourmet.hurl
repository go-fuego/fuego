# Gourmet is UP
GET {{GOURMET_BASE_URL}}
[Options]
retry: 10   # Wait for server to start
retry-interval: 200ms
HTTP 200

# List Recipes
GET {{GOURMET_BASE_URL}}/recipes
HTTP 200
[Asserts]
header "Content-Type" == "text/html; charset=utf-8"
body contains "Dahl de lentilles"

# List Recipes as JSON
GET {{GOURMET_BASE_URL}}/recipes
Accept: application/json
HTTP 200
[Asserts]
header "Content-Type" contains "application/json"
body contains "Dahl de lentilles"

# Recipe page is accessible
GET {{GOURMET_BASE_URL}}/recipes/dahl-de-lentilles
HTTP 200
[Asserts]
body contains "Dahl de lentilles"

# Recipe page does not accept JSON for the moment
GET {{GOURMET_BASE_URL}}/recipes/dahl-de-lentilles
Accept: application/json
HTTP 406
[Asserts]
body contains "Cannot serialize"

# Admin is not accessible without Basic Auth credentials
GET {{GOURMET_BASE_URL}}/admin
HTTP 401

# Admin is accessible with Basic Auth credentials
GET {{GOURMET_BASE_URL}}/admin
Authorization: Basic {{GOURMET_BASIC_AUTH}}
HTTP 200
