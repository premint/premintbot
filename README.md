# premintbot

A Discord bot for Premint.

## How to run locally

Setup gcloud credentials (only need to do this once):

```sh
gcloud iam service-accounts create local-dev
gcloud projects add-iam-policy-binding premint-343516 --member="serviceAccount:local-dev@premint-343516.iam.gserviceaccount.com" --role="roles/owner"
gcloud iam service-accounts keys create credentials.json --iam-account=local-dev@premint-343516.iam.gserviceaccount.com
```

Never commit or share credentials.json.

Make sure env variables are set in `~/.zprofile`:

```sh
export PREMINTBOT_DISCORDAUTHTOKEN="REDACTED"
export PREMINTBOT_GOOGLECLOUDPROJECT="premint-343516"
```

Run the app:

```sh
make dev
```

## How to run tests

```sh
make test
```

## How to install the bot

**https://discord.com/oauth2/authorize?client_id=950933570564800552&scope=bot%20applications.commands&permissions=268438552**

TODO: Make sure the permissions are correct

## How to deploy

Run:

```sh
make ship
```

## How to set env variables in Cloud Run

```sh
gcloud run services update premintbot --set-env-vars PREMINTBOT_DISCORDAUTHTOKEN=REDACTED,PREMINTBOT_GOOGLECLOUDPROJECT=premint-343516
```

## Database (Firestore)

https://console.cloud.google.com/firestore/data/guilds?referrer=search&project=premint-343516

## Bot Pinger

Cloud Scheduler link, runs every 5 min: https://console.cloud.google.com/cloudscheduler?project=premint-343516

## Endpoints

- `GET /health` - Health check endpoint, used to keep the bot alive

## Slash Commands

- `/premint` - Check if a user is in the Premint list

## Legacy Commands

- `!help` - TODO: Update to Slash command
- `!setup` - TODO: Update to Slash command
- `!set-premint` - TODO: Update to Slash command
- `!nuke` - TODO: Update to Slash command