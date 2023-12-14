Requests:
 1. Signup: `curl -X POST localhost:8080/auth/signup -H "Content-Type: application/json" -d '{"username": "tester", "password": "pass"}'`
 2. Login: `curl -X POST localhost:8080/auth/login -H "Content-Type: application/json" -d '{"username": "tester", "password": "pass"}'`
 3. Authenticate: `curl localhost:8080/auth -H "Authorization: 6410c89785ec4e1894a225b396432601" -i`

 Database:
 1. Print tokens: `psql -c "select * from session_token" auth_ranking`
 1. Print users: `psql -c "select * from players" auth_ranking`


