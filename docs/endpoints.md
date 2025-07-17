# /app

### GET

root of the site, just a message  
also increments a counter variable, see [metrics](#adminmetrics)

# /api/users

### POST

create new user  
expects JSON request body in the form of

```json
{
  "password": "password",
  "email": "email@domain.com"
}
```

### GET

dev PLATFORM only  
returns all user records

### PUT

update existing user  
same body as [users](#apiusers) POST

# /api/login

### POST

authenticate as user  
same body as [users](#apiusers) POST

# /api/chirps

### POST

must be authenticated [login](#apilogin)  
expects JSON request body in the form of

```json
{
  "body": "words"
}
```

### GET

takes optional query string `author_id`
returns chirp records

# /api/chirps/{chirpID}

### GET

returns a single record for the given `chirpID`

### DELETE

must be authenticated [login](#apilogin)
delete chirp record for the given `chirpID`

# /api/refresh

### POST

creates a new access token for the owner of the refresh token passed in the `Authorization` header  
one hour time out

### GET

returns all refresh token records

# /api/revoke

### POST

revokes the refresh token passed in the `Authorization` header

# /admin/reset

### POST

dev PLATFORM only  
deletes everything

# /admin/metrics

### GET

display variable that counts how many times the [app](#app) was hit

# /api/healthz

### GET

server status, always 200  
there was a comment about add in 503, but nothing

# /api/polka/webhooks

### POST

some webhook thing
