# Voting Result App

## Description

This is the result application for the voting system. It displays the live results of votes cast in the voting application. The app uses Node.js, Express, Socket.IO for real-time updates, and PostgreSQL for data storage.

## Prerequisites

- Node.js
- PostgreSQL database
- Environment variables for DB connection

## Installation

1. Clone the repository.
2. Navigate to the result folder.
3. Run `npm install` to install dependencies.

## Running the App

Set the following environment variables:

- POSTGRES_USER
- POSTGRES_PASSWORD
- POSTGRES_HOST
- POSTGRES_PORT
- POSTGRES_DB

Then run:

```
node server.js
```

The app will run on port 4000 by default.

## Features

- Real-time vote result updates using WebSockets.
- Web interface to view vote counts.