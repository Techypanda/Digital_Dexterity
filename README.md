# Digital Dexterity

## Purpose

This is a webapp written for [Planetscale Hackathon](https://townhall.hashnode.com/planetscale-hackathon), it uses planetscale for database and golang for backend.

## Environment Setup
- db_password: The Password To Database
- db_username: The Username To Database
- db_address: The Databsae Address
- secret_key: A Secret Key -> Deprecated

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