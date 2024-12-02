# go-crud-auth
simple crud operation with authentication in golang

# Prerequisites
- Ensure Docker is installed on your computer.
You can download and install Docker from the official Docker website

# Running the Application

- Open the project folder in your terminal:
```sh
sudo cd /path/to/your/project
```
- Run the following command to start the application using Docker Compose:
```sh
sudo docker-compose up
```
- Once the command executes, the application will be up and running.

# Notes

- If you want to run the application in detached mode (background), use:
```sh
docker-compose up -d
```
- Stop the application by pressing Ctrl+C or running:
```sh
docker-compose down
```

# API Documentation

## Authentication Endpoints

### POST `/auth/signup/`

This endpoint allows a new user to sign up by providing their details.
- `role` field is used as `normal` or `admin`
- only `admin` has privledge to see active tokens and revoke tokens 

#### Request

```sh
curl -X POST http://localhost:8080/auth/signup/ \
-H "Content-Type: application/json" \
-d '{
  "email": "swapnilbarai1889@gmail.com",
  "mobileNo": "9880887678",
  "username": "swapnil",
  "password": "swapnil@123",
  "role": "admin"
}'
```

### POST `/auth/signin`

- This endpoint allow user to log in and get back the access token and refresh token in header
- Return access token in `Authorization` custom header and  refresh token in`Refresh-Authorization` custom header


#### Request
```sh
curl -i POST http://localhost:8080/auth/signin \
-d "username=swapnil" \
-d "password=barai"
```

### GET `/auth/refresh`

- This endpoint allow user to get new refresh, access token if refresh token still valid
- Return access token in `Authorization` custom header and  refresh token in`Refresh-Authorization` custom header
- Refresh token have longer expiry duration (default it 60 days ,but can be reduce by changing env variable in docker-compose).
- If user find it's access token expire, then they can ask for new access token,refresh token using older refresh token
- This way user doesn't have to logged in again.

### Request
```sh
curl -i POST http://localhost:8080/auth/refresh \
-H "Refresh-Authorization: Bearer your-refresh-jwt-token"
```

## GET `/user/details/:username`

- This endpoint allow user to get the information about other user 
- Access token is required as this route is protected

```sh
curl -X GET http://localhost:8080/user/details/swapnil \
-H "Authorization: Bearer your-access-jwt-token"
```


## GET `/user/active/tokens`

- This endpoint retrieves a list of active tokens for the authenticated user.
- The user making request must be admin

```sh
curl -X GET http://localhost:8080/user/active/tokens \
-H "Authorization: Bearer your-access-jwt-token"

```


## GET `/user/revoke/:tokenID`

- This endpoint revoke active tokens for particular tokenID.
- If tokentype is refresh and if there is any active access token for given  refresh token then access token is also going to revoke as logout mechanism

```sh
curl -X GET http://localhost:8080/user/revoke/12345 \
-H "Authorization: Bearer your-access-jwt-token"

```

