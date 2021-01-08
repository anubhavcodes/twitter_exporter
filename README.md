# Twitter Exporter

An lightweight exporter written in go that exposes very simple metrics like `twitter_followers` and `twitter_following` as gauge values for a prometheus server to scrape.

This is a project very early in development so expect breaking changes.

### Usage

Set up the following environment variables 

```
TWITTER_TOKEN=<your twitter bearer token>
TWITTER_HANDLE=<account of the metrics you wish to scrape>
```

and then run 

```
go run main.go
```

### Use Cases
A typical use case would require you to have a prometheus server scraping metrics for your twitter exporter at every X interval which can then be used to graph on dashboards like grafana.
