# HLFA API Documentation

## Authentication

All API endpoints require JWT authentication via the `Authorization` header:

```bash
Authorization: Bearer <jwt_token>
```

## Endpoints

### Scans

#### Create Scan

```http
POST /api/v1/scans
Content-Type: application/json

{
  "workspace_id": "uuid",
  "targets": [
    {
      "url": "https://example.com",
      "type": "web",
      "headers": {"User-Agent": "..."},
      "cookies": {"session": "..."}
    }
  ],
  "scanner_types": ["web", "network"],
  "rate_limit": 100
}
```

**Response:**
```json
{
  "scan_id": "uuid",
  "status": "running",
  "progress": 0,
  "created_at": "2026-07-25T12:00:00Z"
}
```

#### Get Scan Status

```http
GET /api/v1/scans/{scan_id}
```

**Response:**
```json
{
  "id": "uuid",
  "status": "running",
  "progress": 45,
  "findings_count": 12,
  "started_at": "2026-07-25T12:00:00Z",
  "estimated_completion": "2026-07-25T12:15:00Z"
}
```

#### Get Scan Findings

```http
GET /api/v1/scans/{scan_id}/findings?page=1&limit=20&severity=HIGH
```

**Response:**
```json
{
  "findings": [
    {
      "id": "uuid",
      "title": "IDOR in User Profile",
      "severity": "HIGH",
      "category": "IDOR",
      "status": "OPEN",
      "created_at": "2026-07-25T12:05:00Z"
    }
  ],
  "total": 45,
  "page": 1,
  "limit": 20
}
```

### Findings

#### Get Finding Details

```http
GET /api/v1/findings/{finding_id}
```

**Response:**
```json
{
  "id": "uuid",
  "title": "IDOR in User Profile",
  "description": "User profiles are accessible without authorization checks",
  "severity": "HIGH",
  "category": "IDOR",
  "target": "https://example.com/api/users/{id}",
  "evidence": {
    "request": {
      "method": "GET",
      "url": "https://example.com/api/users/123",
      "headers": {"Authorization": "Bearer ..."}
    },
    "response": {
      "status": 200,
      "body": "{...user data...}"
    },
    "baseline_response": {
      "status": 403,
      "body": "Unauthorized"
    }
  },
  "reproduction_steps": [
    "Login as user A",
    "Request /api/users/999",
    "Observe unauthorized access"
  ],
  "poc_curl": "curl -H 'Authorization: Bearer ...' https://example.com/api/users/999",
  "status": "OPEN",
  "assigned_to": "user@example.com"
}
```

#### Update Finding

```http
PUT /api/v1/findings/{finding_id}
Content-Type: application/json

{
  "status": "CONFIRMED",
  "assigned_to": "user@example.com",
  "comment": "Verified vulnerability, high priority"
}
```

### Heuristic Analysis

#### Set Baseline

```http
POST /api/v1/differential/baseline
Content-Type: application/json

{
  "status": 200,
  "headers": {"Content-Type": "application/json"},
  "body": "{...}",
  "response_time_ms": 150
}
```

#### Analyze Response

```http
POST /api/v1/differential/analyze
Content-Type: application/json

{
  "status": 200,
  "headers": {"Content-Type": "application/json"},
  "body": "{...modified...}",
  "response_time_ms": 200,
  "mutation_type": "parameter_modification"
}
```

**Response:**
```json
{
  "status": "success",
  "analysis": {
    "status_changed": false,
    "content_hash_changed": true,
    "content_length_diff": 150,
    "response_time_diff_ms": 50,
    "headers_changed": false,
    "similarity_score": 0.85,
    "vulnerabilities": [
      {
        "type": "ResponseDiff",
        "severity": "MEDIUM",
        "description": "Response content differs from baseline",
        "confidence": 0.8
      }
    ]
  }
}
```

## Error Responses

### 400 Bad Request
```json
{
  "error": "Invalid request",
  "details": "Missing required field: target_url"
}
```

### 401 Unauthorized
```json
{
  "error": "Unauthorized",
  "details": "Invalid or expired token"
}
```

### 404 Not Found
```json
{
  "error": "Not found",
  "details": "Scan with id 'xyz' not found"
}
```

### 500 Internal Server Error
```json
{
  "error": "Internal server error",
  "details": "An unexpected error occurred"
}
```

## Rate Limiting

- API endpoints are rate-limited to 1000 requests/minute per API key
- Rate limit headers:
  - `X-RateLimit-Limit`: Total requests allowed
  - `X-RateLimit-Remaining`: Requests remaining
  - `X-RateLimit-Reset`: Unix timestamp of rate limit reset

## WebSocket Events

### Real-time Scan Updates

```javascript
socket.on('scan:progress', (data) => {
  console.log(`Scan ${data.scanId}: ${data.progress}%`);
});

socket.on('finding:discovered', (data) => {
  console.log(`New finding: ${data.title}`);
});

socket.on('scan:completed', (data) => {
  console.log(`Scan completed with ${data.findingsCount} findings`);
});
```
