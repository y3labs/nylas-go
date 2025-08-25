# nylas-go

An unofficial, Go client for the \[Nylas v3 API]. It provides a small, well-typed core with resource helpers (Calendars, Messages, Threads, Folders, Drafts, Events, Webhooks, Scheduler, etc.)

> **Note**: This library is not maintained by Nylas, Inc. Y3 Labs is in no way affiliated with Nylas, we are just a fan of their product as well as Go and wanted to share the Go SDK we built for others that may share our interests.

---

## About this SDK

This Go SDK was designed to follow the same patterns and functionality found in the official \[Nylas Python SDK]. The goal is to provide Go developers with an interface while maintaining parity in resource coverage, request/response handling, and overall developer experience. If you are already familiar with the Python SDK, you will find the API surface and usage of this library very similar, making it straightforward to adopt across different projects.  This SDK is implemented purely with the Go 1.24 standard library and so has no external module dependencies.

---

## Features

* **Simple client construction** with functional options for region, base URL, custom `http.Client`, user agent, and timeouts.&#x20;
* **Resource helpers** for first-class endpoints (applications, attachments, auth, availability, calendars, contacts, connectors, drafts, events, folders, grants, messages, notetakers, redirect URIs, scheduler, threads, webhooks).&#x20;
* **Typed responses** (`Response[T]`, `ListResponse[T]`) with header propagation and `request_id` convenience.&#x20;
* **Query encoding** from Go structs to URL params via `EncodeQuery`, including common fields like `limit` and `select`.&#x20;
* **Error handling** that distinguishes OAuth vs. API errors; helpers `IsOAuthError`, `IsAPIError`, `IsTimeoutError`.&#x20;
* **Retries with backoff** on `429/5xx` and transport/timeout errors; configurable policy.&#x20;
* **Multipart uploads** (files + fields) with content-type preservation.&#x20;

---

## Requirements

* Go **1.24+** (per `go.mod`).&#x20;

---

## Installation

```bash
go get github.com/y3labs/nylas-go@latest
```

Then import the SDK package:

```go
import "github.com/y3labs/nylas-go/nylas"
```
---

## Quick start

### Create a client

```go
client := nylas.NewClient(os.Getenv("NYLAS_API_KEY"))
// or customize:
client := nylas.NewClient(
  os.Getenv("NYLAS_API_KEY"),
  nylas.WithRegion(nylas_models.RegionEU),     // sets base URL if none given (Defaults to models.RegionUS
  nylas.WithUserAgent("myapp/1.0"),
  nylas.WithTimeout(60*time.Second),
)
```

The client supports `WithRegion`, `WithServerURL`, `WithHTTPClient`, `WithUserAgent`, `WithTimeout`.&#x20;

### List calendars

```go
ctx := context.Background()
grantID := "your-grant-id"

cal, err := client.Calendars().List(ctx, grantID, nil)
if err != nil {
  log.Fatal(err)
}
for _, c := range cal.Data {
  fmt.Println(c.ID, c.Name)
}
fmt.Println("request-id:", cal.Headers.Get("x-request-id"))
```

This maps to `GET /v3/grants/{identifier}/calendars`. Use query structs for filters; fields are encoded to URL params automatically (`limit`, `select`, etc.). &#x20;

### Look up a calendar

```go
one, err := client.Calendars().Get(ctx, grantID, "primary", nil) // "primary" targets the primary calendar
if err != nil {
  log.Fatal(err)
}
fmt.Println(one.Data.ID, one.Data.Name)
```

---

## Working with responses and headers

All helpers return typed responses and **propagate headers**, so you can read `x-request-id` and rate-limit headers:

```go
page, hdrs, err := nylas.DoJSON[any](client, ctx, "GET", "/v3/grants/"+grantID+"/messages", nil, nil, nil)
fmt.Println("request-id:", hdrs.Get("x-request-id"))
```

&#x20;

---

## Errors

The SDK distinguishes between:

* **OAuth errors** (HTTP 401/403 with `error`, `error_description`, `error_code`), surfaced via `IsOAuthError(err)`.&#x20;
* **API errors** (Nylas JSON error envelope) with fields like `message`, `type`, **provider error** details, `request_id`, `status_code`, surfaced via `IsAPIError(err)`.&#x20;

Example of handling:

```go
_, _, err := nylas.DoJSON[any](client, ctx, "GET",
  "/v3/grants/"+grantID+"/calendars/does-not-exist/events", nil, nil, nil)
if err != nil {
  if e, ok := nylas.IsAPIError(err); ok {
    log.Println("APIError:", e.Message, e.ProviderErrorString(), e.RequestID, e.StatusCode)
  } else if oe, ok := nylas.IsOAuthError(err); ok {
    log.Println("OAuthError:", oe.ErrorDescription, oe.ErrorCode)
  } else {
    log.Println("Other error:", err)
  }
}
```

&#x20;

### Timeouts

Transport timeouts are wrapped as `SDKTimeoutError` and detectable via `IsTimeoutError(err)`.&#x20;

---

## Retries & rate limits

The client retries on `429`, `502`, `503`, `504`, and on transient transport errors, using exponential backoff. The default policy caps attempts and respects `X-RateLimit-Reset` when present.&#x20;

```go
// Internal flow (simplified):
//  - If ShouldRetry(status) and attempts remain, sleep(backoff) and retry
//  - Otherwise parse and return API error
```

See `DoJSON` / `doStream` for the retry loop. &#x20;

The SDK parses rate-limit headers into a `RateLimitInfo` alongside errors when available.&#x20;

---

## Multipart upload

You can post multipart requests (files + fields). Content types are preserved:

```go
resp, err := client.
  // internal helper shown in tests; resource-level helpers may wrap this.
  // See Upload endpoints or use nylas.DoJSON for JSON APIs.
  // ...
  // doMultipart(ctx, "POST", "/v3/upload", fields, files)
  // (example verifies content types and parts)
  // …
_ = resp
```



---

## Regions & base URLs

By default, the client points to **US**. Use `WithRegion(RegionEU)` or `WithServerURL` to override:

```go
c := nylas.NewClient(apiKey, nylas.WithRegion(nylas_models.RegionEU)) // sets EU base URL
// or
c := nylas.NewClient(apiKey, nylas.WithServerURL("https://api.us.nylas.com"))
```

Defaults: `https://api.us.nylas.com` (US), `https://api.eu.nylas.com` (EU). &#x20;

---

## Examples

Run any example with your credentials:

```bash
export NYLAS_API_KEY=...
export NYLAS_GRANT_ID=...
# optional:
export NYLAS_API_URI=...        # override base URL
export NYLAS_CALENDAR_ID=primary

go run ./examples/list_messages
go run ./examples/list_calendar_events
go run ./examples/folders_single_level
go run ./examples/response_headers
go run ./examples/provider_error
go run ./examples/import_events
```

Examples demonstrate header access (`x-request-id`), query params (`limit`, `select`), calendar/event flows, provider error handling, and environment variable usage.   &#x20;

---

## Testing

Run unit tests:

```bash
go test ./...
```

Tests cover client construction, region defaults, resource accessors, query encoding, OAuth/API error parsing, timeouts, multipart uploads, and scheduler flows. &#x20;

---

## Contributing

Issues and PRs are welcome. Please include reproducible cases and, when relevant, add or update tests alongside code changes.

---

## License

This project is licensed under the **MIT License**.
See the [LICENSE](./LICENSE) file for details.
