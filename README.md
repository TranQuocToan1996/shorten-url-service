## Project
   URL shortening system using Golang and front end side use NextJs.

## Raw requirements: 
[Raw Requirement md](raw_requiment.md), [AssignmentHub_raw.pdf](AssignmentHub_raw.pdf)

## How to run project
1. In root folder: Change the name of ".env copy" to ".env". Fix the Postgres and Redis if needed.
2. In ./frontend/fe_shorten_url folder: Change the name of ".env.local" to ".env". Fix the Postgres and Redis if needed. Fix the API service link if needed.
3. In the root folder: run command "Docker compose up -d"

## How to test
1. Default setting for Backend API: "http://localhost:8080"
2. Default setting for Frontend: "http://localhost:3000"
3. For API doc: http://localhost:3000/swagger


## Base62 encode:
[Base62 algo](/docs/base62.md)

## Assume
1. Deploy to cloud AWS
2. Client is other systems
3. Design for scale in the future
4. Shortened URLs cannot be deleted or updated
5. Domain can change
6. No authenticate/authorization for simple.
7. Skip the analysis for simple
8. Some protection in the infra side (Rate limit by IP, DDOS,...)
9. URL is keep even the origin URL die
10. No expiration time for shortened URLs 
11. Skip Monitoring and logging for simple

## Flow Diagram
![Diagram](/docs/drawing/shorten_url.png)

## Security
1. Password/Secret save in secret manager AWS.

## Improvements
1. Webhook add HMAC sign.
2. Add more notify type (FCM, websocket,...), and handle it in background
3. Deploy a version to AWS
4. Swap redis stream -> Kafka, Postgres -> ScyllaDB
5. Add application and infrastructure monitoring (e.g. Prometheus + Grafana, OpenTelemetry)
6. Implement per-user (or per-IP) rate limiting on API endpoints for abuse prevention
7. Add input validation for URL schemes and block private/internal IPs to mitigate SSRF risks
8. Support custom alias/vanity URLs and allow users to manage their links (requires auth)
9. Add OAuth2/JWT authentication and RBAC for user-level permissions
10. Improve tests: add integration tests using Docker Compose + Go testing for end-to-end coverage
11. API usage analytics (counts, geolocation, traffic statistics), possibly using ClickHouse
12. Add paginated listing and search for created links, both via API and optional admin dashboard
13. Automatic HTTPS (Let’s Encrypt support with certbot/cert-manager in production)
14. CLI tool or admin panel for operational tasks (purge, re-encode, debug webhooks, etc)
15. Seamless blue/green or canary deployment support in CI/CD (Argo Rollouts, GitHub Actions)

