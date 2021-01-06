# AireLibre Bot

A bot to tweet data from the AireLibre Paraguay effort

## Credentials
You need to get credentials from a Twitter developer account in order to run this project.

For local development set this values into a .env file, use `go/.env.samble` as an example.

## Local development/run
Go to the `go` folder, then do 
```
$ go run main.go
```
Make sure you have valid credentials in the `.env` file, otherwise you need to set this environment variables before runnig
```ini
ACCESS_TOKEN
ACCESS_TOKEN_SECRET
CONSUMER_KEY
CONSUMER_SECRET
API_URL
```

## Build Docker image
Change **padiazg** for your docker hub username, or for whatever suits your needs
```bash
$ docker build -t padiazg/airelibre-bot:0.1.0 .
```

### Run the container
```bash
$ docker run --rm \
    --name airelibre-bot \
    -e ACCESS_TOKEN= \
    -e ACCESS_TOKEN_SECRET= \
    -e CONSUMER_KEY= \
    -e CONSUMER_SECRET= \
    -e API_URL=https://rald-dev.greenbeep.com/api/v1/aqi \
    padiazg/airelibre-bot:0.1.0
```
