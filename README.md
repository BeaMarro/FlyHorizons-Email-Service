# 📧 FlyHorizons - Email Service

This is the **Email Service** for **FlyHorizons**, an enterprise-grade airline booking system. The service is responsible for sending confirmation emails to passengers, including all their flight details and a QR code for scanning at the airport gate.

---

## 🚀 Overview

This microservice handles the **delivery of transactional emails** related to flight bookings. Upon receiving booking confirmation events via **RabbitMQ**, it composes and sends a fully formatted email with flight itinerary and a **scannable QR code**.

Built with **Go** and the **Gin** framework, the service can be extended to use templated emails and external email delivery APIs (e.g., SendGrid, SMTP, or Mailgun).

---

## 🛠️ Tech Stack

- **Language**: Go (Golang)
- **Framework**: Gin
- **Messaging**: RabbitMQ
- **Email Delivery**: SMTP or external provider (e.g., SendGrid)
- **QR Code Generation**: QR libraries (e.g., `skip2/go-qrcode`)
- **Architecture**: Microservices

---

## 📦 Features

- 📬 **Sends Email Confirmations** for flight bookings
- 🧾 **Includes Flight Details** in email body (passenger, seat, route, time, etc.)
- 🔳 **Generates QR Codes** for airport gate scanning
- 📥 **Listens to Booking Confirmation Events**
- 🔗 **Integrates with external email providers**
- 🔄 **Event-driven communication** via RabbitMQ
- ⚠️ **Centralized error handling** for failed deliveries

---

## 📄 License
This project is shared for educational and portfolio purposes only. Commercial use, redistribution, or modification is not allowed without explicit written permission. All rights reserved © 2025 Beatrice Marro.

## 👤 Author
Beatrice Marro GitHub: https://github.com/beamarro
