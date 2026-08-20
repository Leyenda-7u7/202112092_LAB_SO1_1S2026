UNIVESIDAD DE SAN CARLOS DE GUATEMALA

FACULTAD DE INGENIERIA

ESCUELA DE CIENCAS Y SISTEMAS

LABORATORIO SISTEMAS OPERATIVOS 1 

SECCIÓN P

SEGUNDO SEMESTRE 2026

AUX. JOSÉ DANIEL LORENZANA MEDINA




<p align="center"> MANUAL TECNICO </p>



BRANDON EDUARDO PABLO GARCIA

202112092

Guatemala


---

# Introduccion

El presente manual técnico describe la arquitectura, configuración y funcionamiento del Proyecto 1. La solución fue implementada utilizando virtualización mediante KVM y tres máquinas virtuales Ubuntu Server.

Cada máquina virtual utiliza una tecnología diferente para la administración de contenedores. Además, se desarrollaron tres APIs REST en Go que se comunican entre sí y se implementó un registro privado OCI mediante Zot.

---

# Objetivos

* Brindar la información necesaria para poder  representar la funcionalidad técnica de la estructura, diseño y definición del aplicativo.

* Describir las herramientas utilizadas para el diseño y desarrollo del prototipo

---
# Requerimientos de funcion


|          Requerimientos      |     Descripcion |                                      
|----------------|-------------------------------|
|Visual Studio Code            |Se recomienda el uso de Visual Studio Code que fue la versión donde se programó el sistema de información. |       
| Virtual Manager        | Se utilizara un creador de maquinas virtuales prefriblemente VM |            |            |


---

# Arquitectura general

La infraestructura está compuesta por tres máquinas virtuales:

| Máquina | Dirección IP | Tecnología | Servicios |
|---|---|---|---|
| VM1 | 192.168.122.101 | containerd + nerdctl + BuildKit | API1 y API2 |
| VM2 | 192.168.122.22 | Podman | API3 |
| VM3 | 192.168.122.239 | Docker Engine | Zot Registry |

Las tres máquinas se encuentran conectadas a la red virtual `default` de KVM/libvirt mediante NAT.

---

# Recursos asignados

Cada máquina virtual fue configurada con los siguientes recursos:

- 2 CPU virtuales
- 2048 MiB de memoria RAM
- 15 GiB de almacenamiento
- Ubuntu Server 24.04 LTS
- Adaptador de red VirtIO

Esta configuración permite ejecutar los servicios requeridos utilizando una cantidad moderada de recursos del equipo anfitrión.

---

## VM1

### Función

VM1 aloja API1 y API2.

### Tecnologías

- Ubuntu Server 24.04 LTS
- containerd
- nerdctl
- BuildKit
- Go

### Dirección de red

`192.168.122.101`

### API1

Puerto publicado:

`8081`

Endpoints:

- `GET /health`
- `GET /api1/202112092/call-api2`
- `GET /api1/202112092/call-api3`

### API2

Puerto publicado:

`8082`

Endpoints:

- `GET /health`
- `GET /api2/202112092/call-api1`
- `GET /api2/202112092/call-api3`

### Red interna

API1 y API2 utilizan la red de contenedores:

`red-apis-vm1`

Esto permite que ambas APIs puedan comunicarse utilizando sus nombres de contenedor.

---

## VM2

### Función

VM2 aloja API3 utilizando Podman.

### Tecnologías

- Ubuntu Server 24.04 LTS
- Podman
- Go
- systemd
- Podman Quadlet

### Dirección IP

`192.168.122.22`

### API3

Puerto publicado:

`8083`

Endpoints:

- `GET /health`
- `GET /api3/202112092/call-api1`
- `GET /api3/202112092/call-api2`

### Inicio automático

API3 utiliza un archivo Quadlet llamado:

`api3.container`

Este archivo permite que systemd administre automáticamente el contenedor de API3.

Después de reiniciar VM2 se verificó que:

- `api3.service` se encuentra activo.
- El contenedor API3 inicia automáticamente.
- El endpoint `/health` continúa funcionando.

---

## VM3

### Función

VM3 aloja el registro privado de imágenes OCI.

### Tecnologías

- Ubuntu Server 24.04 LTS
- Docker Engine
- Zot Registry

### Dirección IP

`192.168.122.239`

### Puerto del registro

`5000`

Dirección del registro:

`192.168.122.239:5000`

### Persistencia

Los datos del registro Zot se almacenan de manera persistente utilizando un volumen montado desde:

`~/zot-registry/data`

hacia:

`/var/lib/registry`

Esto permite conservar las imágenes aunque el contenedor Zot sea reiniciado.

---

## Imágenes almacenadas

El registro privado contiene las siguientes imágenes:

- `api1-202112092:v1`
- `api2-202112092:v1`
- `api3-202112092:v1`

Se verificó el almacenamiento mediante:

`GET /v2/_catalog`

y los tags mediante:

`GET /v2/<nombre-imagen>/tags/list`

---

## Comunicación entre APIs

La arquitectura permite seis comunicaciones REST:

1. API1 → API2
2. API1 → API3
3. API2 → API1
4. API2 → API3
5. API3 → API1
6. API3 → API2

Todas las comunicaciones fueron verificadas exitosamente y retornaron:

`"connection": true`

---

## Direcciones y puertos

| Servicio | Dirección |
|---|---|
| API1 | `192.168.122.101:8081` |
| API2 | `192.168.122.101:8082` |
| API3 | `192.168.122.22:8083` |
| Zot | `192.168.122.239:5000` |

---

## Reservas DHCP

Para evitar cambios de dirección IP después de reinicios se configuraron reservas DHCP mediante libvirt.

| VM | MAC | IP |
|---|---|---|
| VM1 | `52:54:00:b0:c8:76` | `192.168.122.101` |
| VM2 | `52:54:00:bf:45:e0` | `192.168.122.22` |
| VM3 | `52:54:00:f7:1c:e2` | `192.168.122.239` |

---

## Registro privado

Las imágenes desarrolladas en VM1 y VM2 fueron publicadas en Zot.

Se verificaron operaciones:

- Push de API1
- Push de API2
- Push de API3
- Pull desde containerd
- Pull desde Podman
- Consulta del catálogo
- Consulta de tags
- Persistencia después de reiniciar Zot

---

## Pruebas realizadas

Se realizaron las siguientes pruebas:

- Verificación de conectividad entre máquinas virtuales.
- Verificación de servicios SSH.
- Verificación de containerd.
- Verificación de Podman.
- Verificación de Docker.
- Verificación de Zot.
- Health check de API1.
- Health check de API2.
- Health check de API3.
- Prueba de las seis comunicaciones REST.
- Push y Pull desde el registro privado.
- Persistencia de Zot.
- Reinicio de máquinas virtuales.
- Inicio automático de los servicios.

---

## Estructura del código

```text
codigo/
├── api1/
│   ├── Dockerfile
│   ├── go.mod
│   └── main.go
├── api2/
│   ├── Dockerfile
│   ├── go.mod
│   └── main.go
└── api3/
    ├── Dockerfile
    ├── go.mod
    └── main.go
