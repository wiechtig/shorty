
POST {{host}}/realms/shorty/protocol/openid-connect/token
[Form]
client_id: shorty
username: user@test.com
password: Test123!
grant_type: password
scope: openid profile email
HTTP 200
[Captures]
access_token: jsonpath "$['access_token']"


POST {{shorty}}/api/short-codes
Authorization: Bearer {{access_token}}
Content-Type: application/json
{
    "shortCode": "test",
    "originalUrl": "https://www.example.com"
}
HTTP 200

