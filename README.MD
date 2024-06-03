# Restaurant-Service

This is a Go application for managing restaurant service.  The Restaurant Service offers robust features designed to streamline restaurant management
## Features

- Restaurant Registration: Create, update, and retrieve restaurant .
- Menu Management: Seamlessly manage restaurant menus, including adding, editing, and removing menu items, as well as updating prices and descriptions.
- Order Processing: Facilitates the efficient processing of customer orders, including order placement, tracking, and fulfillment,
- Availability Tracking: Real-time monitoring of restaurant availability status
- Integration Support: Seamless integration with other services such as the Order Service enabling smooth communication.
## Getting Started with Restaurant-Service

### Requirements

- A working Go environment - [https://go.dev/dl/](https://go.dev/dl/)
- Check the go version with command: go version.
- One should also be familiar with the Golang syntax. [Golang Tour](https://tour.golang.org/) has an excellent guided tour and highly recommended.

### Installation

## GOFR as dependency used for migrations

- To get the GOFR as a dependency, use the command:
  `go get gofr.dev`

- Then use the command `go mod tidy`, to download the necessary packages.


### To Run Server

Run `go run main.go` command in CLI.

## Usage

The application provides the following RESTful endpoints:
1. POST /restaurants: Registers a new restaurant.

2. GET /restaurants: Retrieves restaurant details by ID.

3. PATCH /restaurants/{id}: Updates restaurant details for the specified ID.

4. DELETE /restaurants/{id}: Deletes a restaurant by ID.

5. POST /menu: Creates a new menu item for a restaurant.

6. GET /menu/{restaurant_id}: Retrieves the menu items for a specific restaurant.
7. PUT /restaurant/{restaurant_id}/order/{order_id}: Create mapping for order and restaurant

## Dependencies

The application uses the following dependencies:

- `gofr.dev/pkg/gofr`: A Go web framework used for handling HTTP requests.
- `Restaurant-service/handlers`: Handlers package for handling HTTP requests related to restaurants.
- `Restauanct-service/services`: Services package for business logic related to restaurants.
- `Restaurant-service/stores`:Store package for handling db operations related to restaurants

For any information please reach out to me via rishimanudeepg@gmail.com