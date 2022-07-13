# Digital Dexterity

![Builds Successfully](https://github.com/Techypanda/Digital_Dexterity/actions/workflows/ci.yml/badge.svg)
![Deploys Successfully](https://github.com/Techypanda/Digital_Dexterity/actions/workflows/cd.yml/badge.svg)

## See The Project
- Frontend: [https://digitaldexterity.techytechster.com](https://digitaldexterity.techytechster.com)  
- Rest API: [https://api.digitaldexterity.techytechster.com](https://api.digitaldexterity.techytechster.com)  

## Purpose

This is a application you could feasibly deploy within a organization to have a metric called 'digital dexterity', this metric could be used to evaluate if your fellow colleagues are 'digitally dexterious' when compared against everyone else.  
This is a webapp written for [Planetscale Hackathon](https://townhall.hashnode.com/planetscale-hackathon), it uses planetscale for database, golang for backend and react for frontend (Typescript).

## Branching Model
[I am using a form of trunk based development](https://cloud.google.com/architecture/devops/devops-tech-trunk-based-development)

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