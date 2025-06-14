# Wagecloud Server

[![wakatime](https://wakatime.com/badge/github/wagecloud/wagecloud-server.svg)](https://wakatime.com/badge/github/wagecloud/wagecloud-server)
[![Lines of Code](http://khoakomlem-internal.ddns.net:9000/api/project_badges/measure?project=wagecloud_wagecloud-server_f5966681-fd0e-4ccc-8136-47fc5dc4d892&metric=ncloc&token=sqb_6883ff21a3b534a5b3a481cc52544a04cca6c6b0)](http://khoakomlem-internal.ddns.net:9000/dashboard?id=wagecloud_wagecloud-server_f5966681-fd0e-4ccc-8136-47fc5dc4d892)
[![Quality Gate Status](http://khoakomlem-internal.ddns.net:9000/api/project_badges/measure?project=wagecloud_wagecloud-server_f5966681-fd0e-4ccc-8136-47fc5dc4d892&metric=alert_status&token=sqb_6883ff21a3b534a5b3a481cc52544a04cca6c6b0)](http://khoakomlem-internal.ddns.net:9000/dashboard?id=wagecloud_wagecloud-server_f5966681-fd0e-4ccc-8136-47fc5dc4d892)

A Go-based server application for managing virtual machines.
[Figma](https://www.figma.com/design/qvTI6W2NxgFB9JEMOTv724/WEB-WAGE?node-id=0-1&p=f)

## Overview

Wagecloud Server is a backend service that provides VM management capabilities with features including:

- User account management with role-based access (Admin/User)
- Virtual machine provisioning and management
- Network management for VMs
- Support for different OS and architectures

## Tech Stack

- **Language**: Go
- **Database**: PostgreSQL with Prisma
- **Authentication**: JWT-based
- **Payment Integration**: VNPay
- **Monitoring**: Sentry
- **Caching**: In-memory cache implementation
- **Configuration**: YAML-based

## Prerequisites

- Go 1.18 or higher
- PostgreSQL
- Node.js (for Prisma migrations)
- libvirt (for VM management)

## Configuration

Copy `config/config.example.yml` to `config/config.yml` and adjust the settings:

- Database configuration
- S3 storage settings
- JWT secrets
- VNPay integration details
- Sentry configuration
- Logger settings

## Development

Start the development server:

```bash
make dev
```

### Database Management

Initialize database migrations:

```bash
make init-migrate
```

Generate SQL:

```bash
make sqlc
```

### VM Management

The project includes commands for VM management:

```bash
# Create cloud-init ISO
make cloudinit

# Install a new VM
make install

# Check VM IP addresses
make ip

# Remove a VM
make remove
```

## Project Structure

WIP

## Environment Variables

- `APP_STAGE`: Set to "production" for production environment, otherwise defaults to development

## License

ISC License
