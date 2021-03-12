# Engauge

## Concept

Track user interactions in your apps and products in real-time and see the corresponding stats in a simple dashboard.

Allows for custom instrumentation of applications/products by utilizing a simple and consistent data format.

<!-- Interactions format -->
```go
type Interaction struct {
	// how
	Action *string // required

	// what
	EntityType *string
	EntityID   *string

	// where
	OriginType *string
	OriginID   *string

    // who
	UserType *string
	UserID   *string // required

	DeviceType *string
	DeviceID   *string

	// when
	Timestamp  *string

	// metadata
	Properties map[string]interface{}
}
```

The *only required fields* are `action` and `userID`. All other fields are optional to be sent in an interaction.

## Features

- Automatic Session Detection
- All-Time Statistics
- Interval-Based Statistics
  - Hourly
  - Daily
  - Weekly
  - Monthly
- Unit Metrics

## deployment

Deploying the engauge service is super simple. Use `make build` for a linux binary build. The binary will be placed in a created `bin/` directory inside the repository root directory.

The binary is completely self-contained. Just copy it to your favorite VPS instance, set your environment variables and start it up with a simple `./engauge` command.

You can use `make windows` for windows binary builds.

### environment variables

- `ENGAUGE_HTTPS` can be used to specify if Engauge should use HTTPS (RECOMMENDED to be set to true, defaults to false)
- `ENGAUGE_BASEPATH` is used to specify the name of the root directory in the local filesystem to store data.
- `ENGAUGE_TIMEZONE` specifies the default timezone for Engauge, defaults to the local timezone of the Engauge service instance.
- `ENGAUGE_SESSIONDELAY` specifies, in minutes, how long to wait after the last seen interaction for a user before considering that user's session to be completed.
- `ENGAUGE_USER` is the admin username
- `ENGAUGE_PASSWORD` is the admin password
- `ENGAUGE_JWT` is the JWT secret key
- `ENGAUGE_APIKEY` is the API key for sending interactions to the Engauge service

## Special Values

- `action`s:
  - `conversion`s will be counted towards conversion counts, and will be searched for an `amount` property key for unit metric analytics
- Property Keys:
  - `amount` key, expected to be a numerical value

## Sending Interactions

After you have found the locations within your app or product that you would like to gather data from (we call these "Endpoints"), you will simply need to send the data via HTTP/JSON.

Just send a `POST` to the `/api/interaction` route with your json-encoded interaction object in the request body (and your api key in the `api-key` header).

Easy as pie.

For example:

```sh
curl --location --request POST 'https://example.com/api/interaction' \
--header 'api-key: my-api-key' \
--header 'Content-Type: application/json' \
--data-raw '{
    "action": "close",
    "entityType": "list",
    "entityID": "list-2",
    "originType": "page",
    "originID": "page-2",
    "userType": "customer",
    "userID": "c9449105-9b01-4385-aab5-b45ce4948b96",
    "deviceType": "desktop-web",
    "deviceID": "d396bded-70a6-46ac-bff6-7e09394c43e0",
    "properties": {
        "browser": "safari",
        "browserVersion": "0.8.6",
        "country": "germany",
        "language": "it",
        "mobile": "false",
        "os": "macos",
        "path": "/page-2",
        "ref": "ddg",
        "screenHeight": 200,
        "screenWidth": 600,
        "userAgent": "b0a63af7-4431-4355-9143-cd75c5a6852a"
    }
}'
```

or

```http
POST /api/interaction HTTP/1.1
Host: example.com
api-key: my-api-key
Content-Type: application/json
Content-Length: 686

{
    "action": "close",
    "entityType": "list",
    "entityID": "list-2",
    "originType": "page",
    "originID": "page-2",
    "userType": "customer",
    "userID": "c9449105-9b01-4385-aab5-b45ce4948b96",
    "deviceType": "desktop-web",
    "deviceID": "d396bded-70a6-46ac-bff6-7e09394c43e0",
    "properties": {
        "browser": "safari",
        "browserVersion": "0.8.6",
        "country": "germany",
        "language": "it",
        "mobile": "false",
        "os": "macos",
        "path": "/page-2",
        "ref": "ddg",
        "screenHeight": 200,
        "screenWidth": 600,
        "userAgent": "b0a63af7-4431-4355-9143-cd75c5a6852a"
    }
}
```

where you will need to replace `example.com` with your Engauge instance domain.

## Roadmap

Many more features for Engauge are currently in progress and/or in planning.

Stay tuned!
