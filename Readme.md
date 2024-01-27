# Running the Application
To run the application, follow these instructions:

# 1. Preparation
Before running the server and client applications, you need to prepare the environment.

```bash
make run-prepare
```
This command will:

Start the server environment using Docker Compose.
Populate the database with initial data.
Vendor the Go modules.
Install the required Node.js packages for the client.

# 2. Running the Server
To start the server, run the following command:

```bash
make run-server
```
This command will run the Go server application.

# 3. Running the Client
To start the client application, open new terminal and run:

```bash
make run-client
```
This command will start the development server for the client application.

## Additional Notes
Ensure that Docker is installed and running on your system before executing the run-prepare command.
Make sure that Go and Node.js are installed and properly configured on your machine.
Adjust any paths or configurations as necessary in the Makefile and associated scripts.
With these steps, you should be able to run both the server and client applications smoothly.