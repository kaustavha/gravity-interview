# Getting started

3 main components, front-end, server, fakeiot data generator. and a DB.  

1. postgres setup
https://docs.docker.com/samples/library/postgres/
Note: The pass is base64 mysqlpassword
```
docker run -d -p 5432:5432 --name my-postgres2 -e POSTGRES_PASSWORD=bXlzcWxwYXNzd29yZAo= postgres

psql -h localhost -p 5432 -U postgres -W
CREATE DATABASE iotdb;
\l
\q

```

1. Go server
```
cd $GOPATH
// git clone this repo into src or install it similiar to fakeiot against the github url via go get

// install deps
go get github.com/jinzhu/gorm
go install github.com/jinzhu/gorm

// cd into the repo
go install && $GOPATH/bin/gravitational_interview
```

3. Frontend
```
cd ui/
npm i
npm start
```

4. FakeIOT:
Startup the data generator
```
cd $GOPATH
go get github.com/gravitational/fakeiot
go install github.com/gravitational/fakeiot

$GOPATH/bin/fakeiot --token="YmVhcmVydG9rZW5wYXNzd29yZAo=" --url="https://127.0.0.1:8443" --ca-cert=./fixtures/ca-cert.pem test

$GOPATH/bin/fakeiot --token="YmVhcmVydG9rZW5wYXNzd29yZAo=" --url="https://127.0.0.1:8443" --ca-cert=./fixtures/ca-cert.pem run --period=100s --freq=1s --users=100 --account-id="5a28fa21-c70d-4bf3-b4c4-c4b109d5d269"
```
Note: To recreate the password
```
echo bearertokenpassword | base64  
```  
Likewise the `defaultSigningKey` is `echo mysecretsigningkey | base64`
If the server is running you should see some output, also if youre logging into a DB you should see it update.  

---

# TLS / CACERT notes

Using the mock cert and keys from gravitational/fakeiot/fixtures. 
We'll need to copy over the fixtures dir into this projects root.   
```
$GOPATH/bin/fakeiot --token="shmoken" --url="https://127.0.0.1:8443" --ca-cert=./fixtures/ca-cert.pem test

$GOPATH/bin/fakeiot --token=shmoken --url="https://127.0.0.1:8443" --ca-cert=./fixtures/ca-cert.pem run --period=10s --freq=1s --users=10
```

---

# Basic Architecture
The front-end is scaffolded using `create-react-app` and then the interview html and css is added and minimally reactified. 
Front-end runs on `:3000`.  
Using `react-router-dom` for routing. 
Haven't dug into the internals of `create-react-app` but ES6/7 features work fine :).
Also using regular CSS.

The `create-react-app` development server has a proxy setting set in `package.json` to forward unknown routes to through to the backend running on `:8000`.  

We store the admin account and its usercount, upgrade status etc in the DB  
The server listens for posts from metrics and stores them in a table, and updates the admin account every time  

The server is written in Go and used gorm  

---

# Testing

For server test: `go test`  
For front-end tests: `cd ui && npm test`  

---

# TODO
- [] Security:
  - Env vars instead of hardcode
  - More secure secrets and passwords
- front-end:
  - [x] clean up old `create-react-app` scaffold code
    - [x] auth storage
        - [x] store in local storage - check if there and add to x-session-token header when sending reqs to server/ on first app load; on fail or logout selete the token from local storage or try cookies
    - [x] securely send login form data to server? in header...
    - [x] auto-redirect to db if valid auth in storage
    - [x] websocket or quickpoll soln for dashboard progress updates
    - [x] reactively showing prompts and alerts
    - [x] 1 unit test

- back-end
    - [x] auth handling with bearer tokens in middleware and finished login route
    - [x] dashboard route with(out) websocket handling
    - [x] post handler from fakeiot data generator
    - [x] database & orm for storing fakeiot data
    - [x] edge case/bad data handling from fakeiot generator
    - [x] figure out proper use of bearer tokens and CA certs for fakeiot
    - [x] 1 unit test


# Assumptions:
- Frontend db polls the backend every 1s instead of keeping a socket open or any other solns - this is fragile and may stop updating the FE DB if we change tabs/windows and come back
- backend - dashboard getinfo - doesnt hit database if we have an active account

- Multiple admin users are allowed to be logged in at the same time with different sessions but the same credentials
- Auth based redirects are handled clientside not server side
- we console.log any not ok responses on the frontend for debugging
- Not sure if cross dependency between go packages is good, i.e. pkg 1 requires apis from pkg 2 and an interface of pkg 2 is instantiated when we init pkg 1
- Ideally should use a FOSS auth handler like authboss but this is a good excercise in how these things work
- authenticator has more users logic than i think is good, but i think its an acceptable tradeoff right now
- if the iotdata generator posts 2 users with the same accid, userid, we assume it was a mistake and dont create a new user but update the timestamp
- if the timestamp incoming is empty or equals to Time.IsZero() then we reject the metric - this case is not part of the iot data gen tests
- metrics generator > we create a new item in the db every time, time out and batch write may be faster. We also update the active users details at this time if the acc id matches