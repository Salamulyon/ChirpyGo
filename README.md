# ChirpyGo
Building a Test Server in Go that mimics Twitter. It has users and "chirps" which are equivalent to tweets.

## .env keys
"DB_URL" - A url that is a link to the database. I used postgres

"PLATFORM" - Initialy set to "dev". This allows for resetting the number of hits.Can be left blank if not in dev mode

"secret" - A Bearer Authorization token for authentication. The key should be completely random and the best choice is to use a command to generate a random text like "openssl rand -base64 64" in cmd

"POLKA_KEY" - An API key for a webhook from "Polka",a payment service provider that sends a webhook when a user has subscribed to "Chirpy Red"

## API Endpoints
GET /api/healthz - A Get request that checks if the server is online and working

GET /admin/metrics - A Get request that checks how many times the main page has been visited

POST /admin/reset - A POST request that resets the number of hits to the main page. Requires a "dev" setting in PLATFORM in .env

POST /api/chirps - A POST request that creates a chirp. Requires a Bearer Authorization token and a secret key that is set in .env under "secret".

GET /api/chirps - A GET request that gets all chirps. it also accepts 2 optional query paramters:
    author_id - if an author id is included,it only gets chirps from that author
    sort - if no sort query or sort = asc,it sorts the chirps in ascending order. if sort = desc,it sorts the query in descending order

GET /api/chirps/{chirpID} - A GET request that gets a single chirp by it's id.

DELETE /api/chirps/{chirpID} - A DELETE request. Requires an authentication token and deletes a chirp by it's chirp id

POST /api/users - A POST request. Creates a user based on email and password. Stores the hashed password in the database

PUT /api/users - A PUT request that updates a user's email and password. Requires an authentication token.

POST /api/login - A POST request that logs in a user. Creates an access and a refresh token for the specific user and logs both in the database. The accesss token has  an expiry date of an hour and the refresh of 60 days

POST /api/refresh - A POST request thats creates an access token for the user

POST /api/revoke - A POST request that revokes the current refresh token 

POST /api/polka/webhooks - A POST request for a webhook from the POLKA payment service. accepts a webhook from the POLKA server that tells us to upgrade a user to chirpy red. Will fail if the POLKA key in the .env is ddifferent or missing

/app/ - A link to the file server that serves assets in the assest folder

