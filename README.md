# Coke Day

## Business problem
Two companies, COKE and PEPSI, are sharing an office building but as they are
competitors, they don’t trust each other. Tomorrow is COLA day (for one day), that the
two companies are celebrating. They are gathering a number of business partners in
the building. In order to optimize space utilization, they have decided to set-up a joint
booking system where any user can book one of the 20 meeting rooms available, 10
from each company (C01, C02, ..., C10 and P01, P02, ...., P10).

The booking system has the following functionalities:
* Users can see meeting rooms availability
* Users can book meeting rooms by the hour (first come first served)
* Users can cancel their own reservations

## Tech Stack
TLDR: REST API built in Go and deployed in AWS. Usage of Lambda functions, DynamoDB, Serverless, Cryptography, Validators,
CloudFormation, Makefile, Gingko, JWT Auth, among others.

## How to deploy
To deploy it just run `make` in terminal (serverless must be installed and configured)

## TODO:
* Add logging
* Add/Improve unit testing
* Configure a CORS policy and API throttling in the API Gateway.
* Avoid putting salt in serverless.yml  
* Add local testing
* Organize the project better, avoid duplicate code
