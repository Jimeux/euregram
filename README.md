# euregram

### Next-generation image sharing service.

<img src="ui.jpg" width="600">

---

## Architecture

<img src="architecture.png">

---

## Frontend

Vue SPA hosted on S3 behind CloudFront. 

### Run

```
npm install
npm run dev
```

### Deploy

_Note: S3 bucket names and CloudFront distribution ID values must be set in `sls/Makefile`.
(Requires backend deploy be completed first.)_

Run below command in `sls/Makefile`. 
```
make deploy-frontend
```

---

## Backend

Serverless API with API Gateway + Lambda running behind CloudFront. 
DynamoDB as the datastore, and S3 for image files.

Authentication is done with a combination of Google OAuth and a Lambda authorizer.

### Deploy 

_Note: SSM parameters must be set first. (See details below.)_

Run below command found in `sls/Makefile`.

```
make deploy
```

---

## Image Upload Sequence 

```mermaid
sequenceDiagram
  participant User
  participant API
  participant CloudFront
  participant Lambda
  participant UploadBucket
  participant ImageBucket
  
  User->>API: POST /presign 
  API->>User: https://euregram.com<br>?presign_query 
  User->>CloudFront: PUT https://euregram.com<br>?presign_query

  Note over User,CloudFront: 1. The URL from POST /presign must be used as-is.<br>Only pre-signed headers can be included in the request.<br>2. The image file is included as binary data in the request body.

  CloudFront-->>Lambda: Trigger
  Lambda->>Lambda: Validate file type
  Lambda-->>CloudFront: Return request as-is
  CloudFront->>UploadBucket: Upload image

  User->>API: POST /persist
  API->>ImageBucket: Copy from UploadBucket to ImageBucket
  API->>User: https://euregram.com/image/123

  User->>CloudFront: GET https://euregram.com/image/123

  alt First access
    CloudFront-->>Lambda: Trigger
    Lambda->>ImageBucket:  Fetch object
    Lambda->>Lambda: Resize image
    Lambda-->>CloudFront: Return base64 data
    CloudFront->>CloudFront: Cache
  end

  CloudFront->>User: Return image data
```
---

## Authentication Sequence

_Note: The authentication flow is for demonstration purposes only, and 
was implemented out of pure convenience. Use at your own risk._

```mermaid
sequenceDiagram
  participant User
  participant Frontend
  participant API
  participant Authorizer
  participant Google

  User->>Frontend: Access /
  Frontend->>API: GET /list
  API-->>Authorizer: Intercept
  Authorizer-->>API: HTTP 403
  API->>Frontend: HTTP 403
  
  Frontend->>API: GET /auth/init
  API->>API: Save state value
  API->>Frontend: Return Google URL
  Frontend->>Google: Redirect to Google URL
  
  User->>Google: Log in
  Google->>Frontend: Redirect to /?params
  Frontend->>API: POST /auth/confirm {params}
  
  API->>API: Retreive state value
  API->>Google: Verify params & state
  Google->>API: Return Google token
  API->>Frontend: Return JWT token
  
  Frontend->>API: GET /list
  Note over Frontend,API: Authorization header includes JWT token
  API-->>Authorizer: Intercept
  Authorizer-->>API: HTTP 200
  API->>Frontend: Return image list
```

---

### SSM Parameters

<details>
<summary>Parameter command list</summary>

### Strings

```
aws ssm put-parameter \
--name "euregram-dev-domain" \
--value "my-domain" \
--type String \
--region "us-east-1"  

aws ssm put-parameter \
--name "euregram-dev-hosted-zone" \
--value "my-hosted-zone-id" \
--type String \
--region "us-east-1"  

aws ssm put-parameter \
--name "euregram-dev-google-redirect-url" \
--value "https://euregram.jimeux.com" \
--type String \
--region "us-east-1"  
```

### SecureStrings

```
aws ssm put-parameter \
--name "euregram-dev-jwt-secret" \
--value "my-secret" \
--type "SecureString" \
--region "us-east-1"

aws ssm put-parameter \
--name "euregram-dev-google-client-id" \
--value "my-google-client-id" \
--type "SecureString" \
--region "us-east-1" 

aws ssm put-parameter \
--name "euregram-dev-google-client-secret" \
--value "my-google-secret" \
--type "SecureString" \
--region "us-east-1" 

```

</details>
