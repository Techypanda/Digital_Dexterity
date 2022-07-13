# Digital Dexterity

## Purpose

This is a webapp written for [Planetscale Hackathon](https://townhall.hashnode.com/planetscale-hackathon), it uses planetscale for database and golang for backend.

## Environment Setup
### API
- db_password: The Password To Database
- db_username: The Username To Database
- db_address: The Databsae Address
- secret_key: A Secret Key -> Deprecated
### Infrastructure
- db_username: The username to database
- db_password: The password to database
- db_address: The address to database
- gh_username: The github username to authenticate with to Docker
- gh_token: The github token (PAT) to authenticate with to Docker
- app_name: What you are calling the application in kube context
- image: The image to use for application
- secrets_store: The name of secrets store in kube context
- cors_list: A CSV list of all acceptable CORS
- secret_key: A Secret Key -> Deprecated?
## Running Locally
### Prerequistes
Please install:  
- Docker
Please setup your environment as corresponding in Environment Setup

### Commands
#### Run API (Golang)
```sh
docker compose up api
```
#### Run Web (React App)
```sh
docker compose up web
```
#### Bring Up Entire Environment
```sh
docker compose up
```
#### Debug A Container
```sh
docker compose run api sh # If you want to connect to api
docker compose run web sh # If you want to connect to web
```