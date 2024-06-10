# Billing System with Temporal & Encore

This project implements a billing system using Temporal for workflow orchestration in Golang. The system provides APIs to create bills, add line items, and close bills, ensuring reliable and scalable billing processes.

### Overview
The billing system uses Temporal to manage workflows for creating, updating, and closing bills. It interacts with a PostgreSQL database to persist bill and line item data. Temporal ensures that the billing process is reliable and scalable by managing workflow state and retries.

### Features
- Create a new bill
- Add line items to an existing open bill
- Close a bill
- Query open bills
- Query closed bills

### Requirements
- Golang (1.18 or higher)
- Temporal server
- Docker

### Installation

1. Start temporalite - you can learn more about installation of Temporalite [here](https://github.com/temporalio/temporalite)

``` bash
temporalite start --namespace default
```

2. Go mod Tidy
``` bash
go mod tidy
```

2. Run the encore app
``` bash
encore run
```

### Testing 
- Enter into the folder and run the following command

``` bash
encore test ./...
```
