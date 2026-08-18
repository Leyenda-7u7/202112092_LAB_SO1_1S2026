# Proyecto 1 - Sistemas Operativos 1

**Carnet:** 202112092

## Descripción

Implementación de una infraestructura virtualizada utilizando KVM, tres máquinas virtuales, diferentes runtimes de contenedores, tres APIs REST desarrolladas en Go y un registro privado de imágenes OCI mediante Zot.

## Infraestructura

### VM1
- Ubuntu Server 24.04 LTS
- containerd
- nerdctl
- BuildKit
- API1
- API2
- IP: 192.168.122.101

### VM2
- Ubuntu Server 24.04 LTS
- Podman
- API3
- Podman Quadlet / systemd
- IP: 192.168.122.22

### VM3
- Ubuntu Server 24.04 LTS
- Docker Engine
- Zot Registry
- IP: 192.168.122.239

## Puertos

- API1: 8081
- API2: 8082
- API3: 8083
- Zot: 5000

## Imágenes

- api1-202112092:v1
- api2-202112092:v1
- api3-202112092:v1

## Comunicación REST

- API1 -> API2
- API1 -> API3
- API2 -> API1
- API2 -> API3
- API3 -> API1
- API3 -> API2

Las seis comunicaciones fueron verificadas exitosamente.
