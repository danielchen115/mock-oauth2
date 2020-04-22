# Mock OAuth2
A self-contained, framework-agnostic OAuth 2 mock server to develop and test authorization without worrying about authentication. Spin up permissions logic and token handling without commiting to the dependencies and complexities of a third party auth server.


## Usage
* Have docker installed
* Clone repo
* `docker-compose up -d`
* Create a `config.yml` file similar to that of `config_example.yml`. Note that if changes are made to database credentials, `docker-compose.yml` must reflect that as well.
* Import user pool by running 
`docker-compose exec import-users go run ./cmd/mock-oauth <your_users.json>`
Refer to `users_example.json` for formatting of user pool.
* (Optional) Install MongoDB Compass to interact with your user pool from the auth server side.

## Endpoints
#### GET `/current-user`
##### Parameters
`id`: the id of the current user

Set the current user that your app expects to be "logged in." Its scope and properties are those passed in through the user pool.

#### GET `/authorize`
##### Parameters
`redirect_uri`: the URI redirected to after authorization is complete

Mocks the authorization step in the authorization code grant. This will redirect to `redirect_uri` with an authorization code for the set current user.

#### POST `/token`
##### Parameters
`grant_type`: Either _authorization_code_ for the initial tokens or _refresh_token_ to refresh access token.
`code` or `refresh_token`: for _authorization_code_ and _refresh_token_ respectively
`client_id`: Client id specified in config
`client_secret`: Client secret specified in config

Mocks token exchange with JWTs and refresh token flow. Tokens last for the length of time set in config. 


NOTE: This is only to be used for development and do not have the security features present in a production-ready OAuth server.
