## Overview

This project is a high-concurrency, production-grade banking service designed as a core component of a larger UPI (Unified Payments Interface) ecosystem. It serves as the primary ledger and account management system that communicates with an NPCI (National Payments Corporation of India) switch to facilitate real-time peer-to-peer transactions.

---

## Key Features

- **Pessimistic Locking**  
  Prevents "double-spending" and race conditions during concurrent transaction requests by using `FOR UPDATE` SQL locks on account records.

- **Idempotency Protection**  
  A Redis-backed middleware ensures that retried requests do not result in duplicate debits or credits.

- **On-Read Reconciliation**  
  Every balance inquiry triggers a real-time integrity check that verifies the stored account balance against the sum of all ledger entries.

- **Audit Logging**  
  A structured logging system tracks authentication successes, failures, and critical system events for security auditing.

- **Secure Authentication**  
  API Key-based security with SHA-256 hashing, expiration tracking, and CIDR-based IP whitelisting.

- **ACID Compliant Transactions**  
  Ensures atomicity across account updates and ledger entries within a single database transaction block.

---

## Architecture

![Architecture](./architecture.png)

The project follows a clean Handler-Service-Repository pattern to ensure a separation of concerns and ease of testing.

---

## Tech Stack

- **Language:** Go (Golang)  
- **Database:** PostgreSQL (Primary Store)  
- **Cache:** Redis (Idempotency Store)  
- **Web Framework:** Gin Gonic  
- **Database Tools:** sqlc (Type-safe SQL), goose (Migrations)
