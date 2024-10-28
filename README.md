
# Word of Wisdom TCP Server

This repository contains a "Word of Wisdom" TCP server, designed to provide quotes to clients after they successfully complete a Proof of Work (PoW) challenge. This setup is built to mitigate the risk of DDoS attacks by requiring clients to solve a computational challenge before they can receive a quote. The project includes Docker support for easy deployment of both the server and a client capable of solving the PoW challenge.

## Table of Contents
- [Overview](#overview)
- [Features](#features)
- [Architecture](#architecture)
- [Proof of Work Algorithm](#proof-of-work-algorithm)
- [Setup and Installation](#setup-and-installation)
- [Usage](#usage)
- [License](#license)

## Overview

The "Word of Wisdom" TCP server requires clients to complete a PoW challenge, following a challenge-response protocol, before receiving a quote from a pre-defined list. This PoW mechanism helps prevent DDoS attacks by ensuring that each client must perform some work before consuming server resources.

## Features

- **DDoS Protection**: Uses Proof of Work (PoW) as a challenge-response protocol to prevent abuse.
- **Quote Distribution**: Sends one random quote from a collection upon successful PoW validation.
- **Dockerized**: Docker setup for easy deployment of both server and client.
- **Configurable**: Server settings, such as the difficulty of PoW, can be adjusted for scalability and testing.

## Architecture

The solution consists of two main components:
1. **TCP Server**:
    - Listens for client connections on a specified port.
    - Generates and sends a PoW challenge to each client upon connection.
    - Validates PoW responses from clients before sending a random quote.
2. **Client**:
    - Connects to the server and receives a PoW challenge.
    - Solves the challenge and sends the solution back to the server.
    - Receives and displays the quote upon successful validation.

## Proof of Work Algorithm

### Choice of PoW Algorithm
The PoW challenge involves finding a solution that meets certain computational criteria, typically requiring hash computations that are easy to verify but moderately costly to compute. For this project, a hash-based algorithm was selected for the following reasons:

- **Efficiency**: Hash computations are simple to verify, making it fast for the server to validate responses.
- **Customizable Difficulty**: The difficulty can be adjusted by changing the leading zero requirement on the hash solution, allowing the server to increase or decrease work as needed.
- **Widely Used**: Hash-based PoW algorithms are well-documented and used in various applications, making it easier to implement and extend if needed.

### How It Works
1. Upon connection, the server sends the client a "challenge" consisting of a string and a target number of leading zeroes required in the hash result.
2. The client appends a nonce (a random number) to the challenge string and computes the hash.
3. If the resulting hash meets the target criteria (e.g., starts with the required number of leading zeroes), the client sends the solution (nonce) back to the server.
4. The server verifies the solution and, upon successful verification, sends a random quote to the client.

## Setup and Installation

### Prerequisites
- Docker and Docker Compose should be installed.

### Steps

1. **Clone the repository**:
   ```bash
   git clone https://github.com/Desst/powtcptest.git
   cd powtcptest
   ```

2. **Build Docker images**:
   Build images for both the server and the client using Docker Compose.
   ```bash
   docker-compose build
   ```

3. **Run the server**:
   Start the TCP server container using Docker Compose.
   ```bash
   docker-compose up server
   ```

4. **Run the client**:
   Open another terminal and start the client container to connect to the server and retrieve a quote.
   ```bash
   docker-compose up client
   ```

## Usage

Once both the server and the client are running:

- The server will start listening for incoming TCP connections.
- The client will connect to the server, receive a PoW challenge, solve it, and send back the response.
- Upon successful verification of the PoW, the server sends a quote from the "Word of Wisdom" collection to the client.

### Example Workflow

1. The client connects to the server.
2. The server sends a challenge: a random string with instructions on how many leading zeroes are needed in the hash.
3. The client iterates with different nonces until it finds a hash that meets the criteria.
4. The client sends the correct nonce to the server.
5. The server verifies the solution and sends back a quote.
