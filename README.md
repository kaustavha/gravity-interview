# Getting started

3 main components, front-end, server, fakeiot data generator. 

- FakeIOT:
Startup the data generator
```
cd $GOPATH
go get github.com/gravitational/fakeiot
go install github.com/gravitational/fakeiot
$GOPATH/bin/fakeiot
```

- Go server
```
cd $GOPATH
// git clone this repo into src or install it similiar to fakeiot against the github url
// cd into the repo
go install && $GOPATH/bin/gravitational_interview
```

- Frontend
```
cd ui/
npm i
npm start
```

# TLS / CACERT notes
https://github.com/denji/golang-tls
Generate private key (.key)
# Key considerations for algorithm "RSA" ≥ 2048-bit
openssl genrsa -out server.key 2048

# Key considerations for algorithm "ECDSA" (X25519 || ≥ secp384r1)
# https://safecurves.cr.yp.to/
# List ECDSA the supported curves (openssl ecparam -list_curves)
openssl ecparam -genkey -name secp384r1 -out server.key
Generation of self-signed(x509) public key (PEM-encodings .pem|.crt) based on the private (.key)
openssl req -new -x509 -sha256 -key server.key -out server.crt -days 3650

# Basic Architecture
The front-end is scaffolded using `create-react-app` and then the interview html and css is added and minimally reactified. 
Front-end runs on `:3000`.  
Using `react-router-dom` for routing. 
Haven't dug into the internals of `create-react-app` but ES6/7 features work fine :).
Also using regular CSS.

The `create-react-app` development server has a proxy setting set in `package.json` to forward unknown routes to through to the backend running on `:8000`.  

The server is written in Go and will use:
- gorm as the orm


# TODO
- front-end:
    - clean up old `create-react-app` scaffold code
    - auth storage
        - store in local storage - check if there and add to x-session-token header when sending reqs to server/ on first app load; on fail or logout selete the token from local storage
    - securely send login form data to server? in header...
    - auto-redirect to db if valid auth in storage
    - websocket or quickpoll soln for dashboard progress updates
    - reactively showing prompts and alerts
    - 1 unit test

- back-end
    - auth handling with bearer tokens in middleware and finished login route
    - dashboard route with websocket handling
    - post handler from fakeiot data generator
    - database & orm for storing fakeiot data
    - edge case/bad data handling from fakeiot generator
    - figure out proper use of bearer tokens and CA certs for fakeiot
    -  1 unit test
