# Voting App

## Description

This is the voting application built with Spring Boot. It provides a web interface for users to cast votes between two options (A and B). Votes are sent to a Kafka topic for processing by the worker service.

## Prerequisites

- Java 22
- Maven
- Kafka cluster
- Environment variables for configuration

## Installation

1. Clone the repository.
2. Navigate to the votes folder.
3. Run `mvn clean install` to build the application.

## Running the App

Set the following environment variables:

- OPTION_A: Name of the first voting option
- OPTION_B: Name of the second voting option
- KAFKA_BOOTSTRAP_SERVERS: Kafka bootstrap servers (e.g., localhost:9092)

Then run:

```
mvn spring-boot:run
```

Alternatively, package the app with `mvn package` and run the JAR:

```
java -jar target/vote-0.0.1-SNAPSHOT.jar
```

The app will run on port 8080 by default.

## Features

- Web interface for casting votes using Thymeleaf templates
- Asynchronous sending of votes to Kafka topic "votes"
- Cookie-based voter identification to prevent multiple votes