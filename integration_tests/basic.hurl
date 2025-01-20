# Test the basic endpoints
# Basic POST endpoint test with all required headers
POST http://localhost:8088/
Content-Type: application/json
test: test
{
    "name": "John"
}
HTTP 200
[Asserts]
header "Content-Type" == "application/json"
header "X-Hello" == "World"
jsonpath "$.message" == "HELLO, JOHN"
jsonpath "$.best" == "Fuego!"

# Test validation error
POST http://localhost:8088/
Content-Type: application/json
{
    "name": ""
}
HTTP 400

# Test header validation
POST http://localhost:8088/
Content-Type: application/json
{
    "name": "John"
}
HTTP 500

# Test standard handler
GET http://localhost:8088/std
HTTP 200
[Asserts]
body == "Hello, World!" 