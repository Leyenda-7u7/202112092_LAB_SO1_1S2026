# Proyecto 1 - Sistemas Operativos 1

**Nombre:** BRANDON EDUARDO PABLO GARCIA  
**Carnet:** 202112092  
**Sistema operativo utilizado en las VMs:** Ubuntu Server 24.04 LTS

---

## Descripción

Este proyecto implementa una infraestructura virtualizada utilizando KVM y tres máquinas virtuales.

Cada máquina utiliza una tecnología diferente para la administración de contenedores. Sobre esta infraestructura se ejecutan tres APIs REST desarrolladas en Go y un registro privado OCI mediante Zot.

El objetivo principal es demostrar virtualización, contenerización, comunicación entre servicios y distribución de imágenes mediante diferentes runtimes.

---

## Arquitectura

| Máquina | IP | Tecnología | Servicios |
|---|---|---|---|
| VM1 | `192.168.122.101` | containerd + nerdctl + BuildKit | API1 y API2 |
| VM2 | `192.168.122.22` | Podman | API3 |
| VM3 | `192.168.122.239` | Docker Engine | Zot Registry |

### Puertos

| Servicio | Puerto |
|---|---:|
| API1 | `8081` |
| API2 | `8082` |
| API3 | `8083` |
| Zot | `5000` |

---

## Diagrama general

```text
VM1 - 192.168.122.101
├── containerd
├── BuildKit
├── API1 :8081
└── API2 :8082
       ↕
       ↕ REST / HTTP
       ↕
VM2 - 192.168.122.22
├── Podman
└── API3 :8083


VM1 / containerd ──────────┐
                           │
                           ▼
                  VM3 - 192.168.122.239
                  Docker + Zot :5000
                           ▲
                           │
VM2 / Podman ──────────────┘
```

---

## API1

Ubicada en VM1.

### Endpoints

```text
GET /health
GET /api1/202112092/call-api2
GET /api1/202112092/call-api3
```

Imagen:

```text
api1-202112092:v1
```

---

## API2

Ubicada en VM1.

### Endpoints

```text
GET /health
GET /api2/202112092/call-api1
GET /api2/202112092/call-api3
```

Imagen:

```text
api2-202112092:v1
```

---

## API3

Ubicada en VM2.

### Endpoints

```text
GET /health
GET /api3/202112092/call-api1
GET /api3/202112092/call-api2
```

Imagen:

```text
api3-202112092:v1
```

API3 utiliza Podman Quadlet y systemd para iniciar automáticamente después de reiniciar VM2.

---

## Comunicación REST

Se validaron las seis comunicaciones posibles entre las tres APIs:

```text
API1 → API2
API1 → API3

API2 → API1
API2 → API3

API3 → API1
API3 → API2
```

Todas las pruebas retornaron:

```json
"connection": true
```

---

## Registro privado Zot

Zot se encuentra desplegado en VM3 mediante Docker Engine.

Dirección:

```text
192.168.122.239:5000
```

Imágenes almacenadas:

```text
api1-202112092:v1
api2-202112092:v1
api3-202112092:v1
```

Se verificaron:

- Push de las tres imágenes.
- Pull mediante containerd/nerdctl.
- Pull mediante Podman.
- Consulta del catálogo.
- Consulta de tags.
- Persistencia después del reinicio.

---

## Persistencia

Los datos de Zot se almacenan mediante:

```text
~/zot-registry/data
```

montado en:

```text
/var/lib/registry
```

Esto permite conservar las imágenes incluso después de reiniciar o recrear el contenedor Zot.

---

## Direcciones IP

Las direcciones fueron reservadas mediante DHCP de libvirt:

| VM | MAC | IP |
|---|---|---|
| VM1 | `52:54:00:b0:c8:76` | `192.168.122.101` |
| VM2 | `52:54:00:bf:45:e0` | `192.168.122.22` |
| VM3 | `52:54:00:f7:1c:e2` | `192.168.122.239` |

---

## Estructura del Proyecto

```text
PROYECTO 1/
├── codigo/
│   ├── api1/
│   ├── api2/
│   └── api3/
│
├── configuracion/
│   ├── vm1/
│   ├── vm2/
│   └── vm3/
│
├── scripts/
│   ├── pruebas-comunicacion.sh
│   ├── pruebas-registro.sh
│   └── reservar-ips.sh
│
└── documentacion/
    ├── capturas/
    ├── guia-instalacion/
    │   └── GUIA_INSTALACION.md
    └── manual-tecnico/
        └── MANUAL_TECNICO.md
```

---

## Scripts incluidos

### Pruebas de comunicación

```bash
./scripts/pruebas-comunicacion.sh
```

Ejecuta las seis comunicaciones REST.

### Pruebas del registro

```bash
./scripts/pruebas-registro.sh
```

Comprueba Zot, catálogo y tags.

### Reservas DHCP

```bash
./scripts/reservar-ips.sh
```

Configura las reservas de VM1, VM2 y VM3 en la red virtual de libvirt.

---

## Documentación

La documentación completa se encuentra en:

### Manual técnico

```text
documentacion/manual-tecnico/MANUAL_TECNICO.md
```

### Guía de instalación

```text
documentacion/guia-instalacion/GUIA_INSTALACION.md
```

### Evidencias

```text
documentacion/capturas/
```

---

## Pruebas realizadas

Se comprobó correctamente:

- Funcionamiento de KVM.
- SSH en las tres máquinas.
- containerd.
- nerdctl.
- BuildKit.
- Podman.
- Docker Engine.
- Zot Registry.
- API1.
- API2.
- API3.
- Seis comunicaciones REST.
- Push de imágenes.
- Pull de imágenes.
- Persistencia de Zot.
- Reservas DHCP.
- Reinicio de las máquinas virtuales.
- Inicio automático de los servicios.

---

## Resultado

La infraestructura permite ejecutar las tres APIs en diferentes tecnologías de contenerización, mantener comunicación REST entre ellas y almacenar/distribuir sus imágenes mediante un registro privado OCI.

Las pruebas finales fueron realizadas satisfactoriamente.
