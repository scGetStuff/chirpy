# /api/users

### POST

create new user  
expects JSON request body in the form of

```json
{
  "password": "word",
  "email": "email@domain.com"
}
```

### GET

dev PLATFORM only

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

# /api/chirps/{chirpID}

### GET

### DELETE

# /api/refresh

### POST

### GET

# /api/revoke

### POST

# /admin/metrics

### GET

# /admin/reset

### POST

dev PLATFORM only  
deletes everything

# /api/healthz

### GET

server status, always 200  
there was a comment about add in 503, but nothing

# /api/polka/webhooks

### POST
