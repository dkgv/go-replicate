# go-replicate

go-replicate is a Go client library for [Replicate](https://replicate.com).

## Installation

```bash
go get github.com/dkgv/go-replicate
```

## Usage

Get started by [creating an API token](https://replicate.com/account/api-tokens) on Replicate and instantiating a new client:

```go
client := replicate.NewClient("REPLICATE_TOKEN")
```

### Models

#### Retrieving a model

```go
model, err := client.Models.Get(ctx, "OWNER", "MODEL_NAME")
```

### Predictions

#### Creating a prediction

```go
var input any
// ...
prediction, err := client.Predictions.Create(ctx, "MODEL_ID", input)
```

#### Creating a prediction with a webhook

```go
webhook := replicate.Webook{
    CallbackURL: "https://example.com/callback",
    Events: []replicate.PredictionEvent{
        replicate.EventPredictionStarted,
        replicate.EventPredictionCompleted,
    },
}
prediction, err := client.Predictions.CreateWithWebhook(ctx, "MODEL_ID", input, webhook)
```

#### Getting a prediction

```go
prediction, err := client.Predictions.Get(ctx, "PREDICTION_ID")
```

#### Awaiting a prediction

```go
type Destination struct {
    // ...
}

var destination Destination
err := client.Predictions.Await(ctx, "PREDICTION_ID", &destination)
```

#### Cancelling a prediction

```go
err := client.Predictions.Cancel(ctx, "PREDICTION_ID")
```
