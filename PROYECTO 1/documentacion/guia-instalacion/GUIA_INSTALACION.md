UNIVESIDAD DE SAN CARLOS DE GUATEMALA

FACULTAD DE INGENIERIA

ESCUELA DE CIENCAS Y SISTEMAS

LABORATORIO SISTEMAS OPERATIVOS 1 

SECCIÓN P

SEGUNDO SEMESTRE 2026

AUX. JOSÉ DANIEL LORENZANA MEDINA




<p align="center"> GUIA DE INSTALACION </p>



BRANDON EDUARDO PABLO GARCIA

202112092

Guatemala


---



# Objetivo

La presente guía describe el procedimiento necesario para instalar, configurar y ejecutar la infraestructura desarrollada para el Proyecto 1.

La solución utiliza KVM para virtualización y está formada por tres máquinas virtuales. Cada una utiliza una tecnología diferente para el manejo de contenedores:

- **VM1:** containerd + nerdctl + BuildKit.
- **VM2:** Podman.
- **VM3:** Docker Engine + Zot Registry.

Además, se ejecutan tres APIs REST desarrolladas en Go que se comunican entre sí.

---

# Arquitectura del proyecto

La infraestructura implementada es la siguiente:

| Máquina | Dirección IP | Tecnología principal | Servicios |
|---|---|---|---|
| VM1 | `192.168.122.101` | containerd + nerdctl + BuildKit | API1 y API2 |
| VM2 | `192.168.122.22` | Podman | API3 |
| VM3 | `192.168.122.239` | Docker Engine | Zot Registry |

Los puertos utilizados son:

| Servicio | Puerto |
|---|---:|
| API1 | `8081` |
| API2 | `8082` |
| API3 | `8083` |
| Zot Registry | `5000` |

La comunicación general queda de la siguiente forma:

```text
VM1 - 192.168.122.101
├── API1 :8081
└── API2 :8082
        ↕
        ↕ REST/HTTP
        ↕
VM2 - 192.168.122.22
└── API3 :8083

VM3 - 192.168.122.239
└── Zot Registry :5000
    ├── api1-202112092:v1
    ├── api2-202112092:v1
    └── api3-202112092:v1
```

---

# Requisitos del equipo anfitrión

Se requiere:

- Sistema operativo Linux.
- Procesador con soporte de virtualización.
- KVM.
- QEMU.
- libvirt.
- Virt-Manager.
- Git.
- Cliente SSH.

Para verificar si el procesador soporta virtualización:

```bash
grep -Eoc '(vmx|svm)' /proc/cpuinfo
```

Un resultado mayor a `0` indica que existe soporte para virtualización.

También puede verificarse mediante:

```bash
lscpu | grep -i virtualization
```

Verificar la existencia del dispositivo KVM:

```bash
ls -l /dev/kvm
```
---

# Instalación de KVM, libvirt y Virt-Manager

Actualizar los repositorios:

```bash
sudo apt update
```

Instalar los paquetes necesarios:

```bash
sudo apt install qemu-kvm libvirt-daemon-system libvirt-clients bridge-utils virt-manager -y
```

Agregar el usuario actual a los grupos `libvirt` y `kvm`:

```bash
sudo usermod -aG libvirt $USER
sudo usermod -aG kvm $USER
```

Cerrar sesión e iniciar nuevamente para aplicar los cambios.

Verificar los grupos:

```bash
groups
```

Verificar libvirt:

```bash
virsh --version
```

Verificar Virt-Manager:

```bash
virt-manager --version
```

Verificar la red virtual:

```bash
virsh net-list --all
```

La red:

```text
default
```

debe encontrarse activa.

Si estuviera inactiva:

```bash
virsh net-start default
```

Para habilitarla automáticamente:

```bash
virsh net-autostart default
```

---

# Creación de las máquinas virtuales

Se deben crear tres máquinas virtuales utilizando Virt-Manager.

Los nombres utilizados son:

```text
VM1
VM2
VM3
```

Cada máquina utiliza:

| Recurso | Valor |
|---|---|
| Memoria RAM | 2048 MiB |
| CPU | 2 vCPU |
| Disco | 15 GiB |
| Sistema operativo | Ubuntu Server 24.04 LTS |
| Red | `default` NAT |
| Modelo de red | VirtIO |

Durante la instalación de Ubuntu Server se deben utilizar las siguientes opciones generales:

- Ubuntu Server normal, no minimizado.
- DHCP para configuración inicial de red.
- Proxy vacío.
- Mirror predeterminado.
- Utilizar todo el disco virtual.
- Ubuntu Pro: `Skip for now`.
- Instalar OpenSSH Server.
- No instalar Featured Server Snaps.

El usuario utilizado dentro de las máquinas es:

```text
pablo
```

Los nombres de host son:

```text
VM1 → vm1
VM2 → vm2
VM3 → vm3
```

---

# Configuración de SSH

Después de instalar Ubuntu Server, iniciar sesión en cada VM.

Verificar SSH:

```bash
systemctl is-active ssh
```

Verificar que esté habilitado:

```bash
systemctl is-enabled ssh
```

Los resultados esperados son:

```text
active
enabled
```

Si SSH se encuentra inactivo:

```bash
sudo systemctl enable --now ssh
```

Verificar nuevamente:

```bash
systemctl is-active ssh
systemctl is-enabled ssh
```

Consultar la dirección IP:

```bash
ip -br addr
```

---

# Direcciones IP reservadas

La infraestructura utiliza las siguientes direcciones:

```text
VM1 → 192.168.122.101
VM2 → 192.168.122.22
VM3 → 192.168.122.239
```

Las MAC utilizadas son:

| Máquina | MAC | IP |
|---|---|---|
| VM1 | `52:54:00:b0:c8:76` | `192.168.122.101` |
| VM2 | `52:54:00:bf:45:e0` | `192.168.122.22` |
| VM3 | `52:54:00:f7:1c:e2` | `192.168.122.239` |

El repositorio incluye el script:

```text
scripts/reservar-ips.sh
```

Este script debe ejecutarse únicamente en el equipo anfitrión.

También pueden configurarse manualmente.

VM1:

```bash
virsh -c qemu:///system net-update default add ip-dhcp-host \
"<host mac='52:54:00:b0:c8:76' name='vm1' ip='192.168.122.101'/>" \
--live --config
```

VM2:

```bash
virsh -c qemu:///system net-update default add ip-dhcp-host \
"<host mac='52:54:00:bf:45:e0' name='vm2' ip='192.168.122.22'/>" \
--live --config
```

VM3:

```bash
virsh -c qemu:///system net-update default add ip-dhcp-host \
"<host mac='52:54:00:f7:1c:e2' name='vm3' ip='192.168.122.239'/>" \
--live --config
```

Verificar las reservas:

```bash
virsh -c qemu:///system net-dumpxml default | grep "<host"
```

---

# Verificación de conectividad

Desde el host se pueden comprobar las tres VMs:

```bash
ping -c 2 192.168.122.101
ping -c 2 192.168.122.22
ping -c 2 192.168.122.239
```

También se debe verificar SSH:

```bash
ssh pablo@192.168.122.101
ssh pablo@192.168.122.22
ssh pablo@192.168.122.239
```

---

# VM1 - containerd

Conectarse a VM1:

```bash
ssh pablo@192.168.122.101
```

Actualizar el sistema:

```bash
sudo apt update
sudo apt upgrade -y
```

Verificar containerd después de su instalación:

```bash
containerd --version
```

El servicio debe encontrarse activo:

```bash
systemctl is-active containerd
```

También debe estar habilitado:

```bash
systemctl is-enabled containerd
```

---

# VM1 - nerdctl

Verificar nerdctl:

```bash
nerdctl --version
```

Realizar una prueba descargando Alpine:

```bash
sudo nerdctl pull docker.io/library/alpine:latest
```

Ejecutar:

```bash
sudo nerdctl run --rm docker.io/library/alpine:latest \
echo "containerd funciona correctamente en VM1"
```

Resultado esperado:

```text
containerd funciona correctamente en VM1
```

---

# VM1 - BuildKit

La configuración utilizada durante el proyecto se encuentra en:

```text
configuracion/vm1/buildkitd.toml
configuracion/vm1/buildkit.service
```

Crear el directorio:

```bash
sudo mkdir -p /etc/buildkit
```

Copiar la configuración:

```bash
sudo cp buildkitd.toml /etc/buildkit/buildkitd.toml
```

Copiar el servicio:

```bash
sudo cp buildkit.service /etc/systemd/system/buildkit.service
```

Recargar systemd:

```bash
sudo systemctl daemon-reload
```

Habilitar BuildKit:

```bash
sudo systemctl enable --now buildkit
```

Verificar:

```bash
systemctl is-active buildkit
systemctl is-enabled buildkit
```

Resultado esperado:

```text
active
enabled
```

Verificar los workers:

```bash
sudo buildctl debug workers
```

---

# Código de API1 y API2

El código fuente se encuentra en:

```text
codigo/api1/
codigo/api2/
```

Cada API posee:

```text
Dockerfile
go.mod
main.go
```

---

# Construcción de API1

En VM1, ingresar al directorio de API1:

```bash
cd codigo/api1
```

Construir la imagen:

```bash
sudo nerdctl build -t api1-202112092:v1 .
```

Verificar:

```bash
sudo nerdctl images
```

Debe aparecer:

```text
api1-202112092:v1
```

---

# Construcción de API2

Ingresar al directorio:

```bash
cd codigo/api2
```

Construir:

```bash
sudo nerdctl build -t api2-202112092:v1 .
```

Verificar:

```bash
sudo nerdctl images
```

Debe aparecer:

```text
api2-202112092:v1
```

---

# Crear red de API1 y API2

Crear:

```bash
sudo nerdctl network create red-apis-vm1
```

Verificar:

```bash
sudo nerdctl network ls
```

---

# Ejecutar API1

Ejecutar:

```bash
sudo nerdctl run -d \
  --name api1 \
  --restart=unless-stopped \
  --network red-apis-vm1 \
  -p 8081:8080 \
  -e CARNET=202112092 \
  -e API2_URL="http://api2:8080" \
  -e API3_URL="http://192.168.122.22:8083" \
  api1-202112092:v1
```

---

# Ejecutar API2

```bash
sudo nerdctl run -d \
  --name api2 \
  --restart=unless-stopped \
  --network red-apis-vm1 \
  -p 8082:8080 \
  -e CARNET=202112092 \
  -e API1_URL="http://api1:8080" \
  -e API3_URL="http://192.168.122.22:8083" \
  api2-202112092:v1
```

También puede utilizarse el script:

```text
configuracion/vm1/ejecutar-apis.sh
```

Dar permisos:

```bash
chmod +x ejecutar-apis.sh
```

Ejecutar:

```bash
./ejecutar-apis.sh
```

Verificar:

```bash
sudo nerdctl ps
```

---

# Pruebas de API1 y API2

API1:

```bash
curl http://localhost:8081/health
```

API2:

```bash
curl http://localhost:8082/health
```

Ambas deben devolver:

```text
"status":"UP"
```

---

# VM2 - Instalación de Podman

Conectarse:

```bash
ssh pablo@192.168.122.22
```

Actualizar:

```bash
sudo apt update
```

Instalar Podman:

```bash
sudo apt install podman -y
```

Verificar:

```bash
podman --version
```

Consultar información:

```bash
podman info
```

---

# Prueba de Podman

Descargar Alpine:

```bash
podman pull docker.io/library/alpine:latest
```

Ejecutar:

```bash
podman run --rm docker.io/library/alpine:latest \
echo "Podman funciona correctamente en VM2"
```

Comprobar acceso HTTP:

```bash
podman run --rm docker.io/library/alpine:latest \
sh -c 'wget -qO- http://example.com >/dev/null && echo "Internet HTTP funciona correctamente"'
```

---

# Construcción de API3

El código se encuentra en:

```text
codigo/api3/
```

Ingresar:

```bash
cd codigo/api3
```

Construir:

```bash
podman build -t api3-202112092:v1 .
```

Verificar:

```bash
podman images
```

Debe aparecer:

```text
api3-202112092:v1
```

---

# Configuración automática de API3

API3 utiliza Podman Quadlet y systemd.

El archivo se encuentra en:

```text
configuracion/vm2/api3.container
```

Crear directorio:

```bash
mkdir -p ~/.config/containers/systemd
```

Copiar:

```bash
cp api3.container ~/.config/containers/systemd/api3.container
```

Permitir que los servicios del usuario puedan iniciar aunque no exista una sesión SSH abierta:

```bash
sudo loginctl enable-linger pablo
```

Verificar:

```bash
loginctl show-user pablo -p Linger
```

Resultado esperado:

```text
Linger=yes
```

Recargar systemd:

```bash
systemctl --user daemon-reload
```

Iniciar:

```bash
systemctl --user start api3.service
```

Verificar:

```bash
systemctl --user is-active api3.service
```

Resultado esperado:

```text
active
```

Comprobar Podman:

```bash
podman ps
```

Debe aparecer API3 exponiendo:

```text
0.0.0.0:8083->8080/tcp
```

También puede utilizarse:

```text
configuracion/vm2/configurar-api3.sh
```

---

# Prueba de API3

```bash
curl http://localhost:8083/health
```

Debe devolver información similar a:

```json
{
  "status": "UP",
  "message": "API3 is Ready",
  "VM": "VM2",
  "carnet": "202112092"
}
```

---

# VM3 - Instalación de Docker Engine

Conectarse:

```bash
ssh pablo@192.168.122.239
```

Actualizar:

```bash
sudo apt update
```

Instalar dependencias:

```bash
sudo apt install ca-certificates curl -y
```

Crear el directorio de llaves:

```bash
sudo install -m 0755 -d /etc/apt/keyrings
```

Descargar la llave:

```bash
sudo curl -fsSL https://download.docker.com/linux/ubuntu/gpg \
-o /etc/apt/keyrings/docker.asc
```

Asignar permisos:

```bash
sudo chmod a+r /etc/apt/keyrings/docker.asc
```

Agregar el repositorio:

```bash
sudo tee /etc/apt/sources.list.d/docker.sources > /dev/null <<EOF
Types: deb
URIs: https://download.docker.com/linux/ubuntu
Suites: $(. /etc/os-release && echo "${UBUNTU_CODENAME:-$VERSION_CODENAME}")
Components: stable
Architectures: $(dpkg --print-architecture)
Signed-By: /etc/apt/keyrings/docker.asc
EOF
```

Actualizar:

```bash
sudo apt update
```

Instalar Docker:

```bash
sudo apt install docker-ce docker-ce-cli containerd.io docker-buildx-plugin docker-compose-plugin -y
```

---

# Verificación de Docker

Ver versión:

```bash
sudo docker --version
```

Verificar servicio:

```bash
systemctl is-active docker
systemctl is-enabled docker
```

Resultados esperados:

```text
active
enabled
```

Ejecutar prueba:

```bash
sudo docker run hello-world
```

Debe aparecer:

```text
Hello from Docker!
```

---

# Instalación de Zot Registry

Crear directorio persistente:

```bash
mkdir -p ~/zot-registry/data
```

Descargar Zot:

```bash
sudo docker pull ghcr.io/project-zot/zot:latest
```

Ejecutar:

```bash
sudo docker run -d \
  --name zot \
  --restart unless-stopped \
  -p 5000:5000 \
  -v "$HOME/zot-registry/data:/var/lib/registry" \
  ghcr.io/project-zot/zot:latest
```

También puede utilizarse:

```text
configuracion/vm3/ejecutar-zot.sh
```

Verificar:

```bash
sudo docker ps
```

Debe aparecer:

```text
zot
```

con el puerto:

```text
0.0.0.0:5000->5000/tcp
```

---

# Verificación de Zot

Probar:

```bash
curl http://localhost:5000/v2/
```

Consultar catálogo:

```bash
curl http://localhost:5000/v2/_catalog
```

Antes de subir imágenes puede obtenerse:

```json
{"repositories":[]}
```

---

# Publicar API1 en Zot

Desde VM1:

```bash
sudo nerdctl tag \
docker.io/library/api1-202112092:v1 \
192.168.122.239:5000/api1-202112092:v1
```

Publicar:

```bash
sudo nerdctl --insecure-registry push \
192.168.122.239:5000/api1-202112092:v1
```

---

# Publicar API2 en Zot

Etiquetar:

```bash
sudo nerdctl tag \
docker.io/library/api2-202112092:v1 \
192.168.122.239:5000/api2-202112092:v1
```

Publicar:

```bash
sudo nerdctl --insecure-registry push \
192.168.122.239:5000/api2-202112092:v1
```

---

# Publicar API3 en Zot

Desde VM2:

```bash
podman tag \
api3-202112092:v1 \
192.168.122.239:5000/api3-202112092:v1
```

Publicar:

```bash
podman push --tls-verify=false \
192.168.122.239:5000/api3-202112092:v1
```

---

# Verificación del catálogo

Desde cualquier máquina con acceso a VM3:

```bash
curl http://192.168.122.239:5000/v2/_catalog
```

Resultado esperado:

```json
{
  "repositories": [
    "api1-202112092",
    "api2-202112092",
    "api3-202112092"
  ]
}
```

El orden de los repositorios puede variar.

---

# Verificar tags

API1:

```bash
curl http://192.168.122.239:5000/v2/api1-202112092/tags/list
```

API2:

```bash
curl http://192.168.122.239:5000/v2/api2-202112092/tags/list
```

API3:

```bash
curl http://192.168.122.239:5000/v2/api3-202112092/tags/list
```

Las tres deben mostrar:

```text
v1
```

---

# Prueba de distribución mediante Pull

Para demostrar que Zot puede distribuir imágenes entre diferentes runtimes se realizan las siguientes pruebas.

Desde VM1 descargar API3:

```bash
sudo nerdctl --insecure-registry pull \
192.168.122.239:5000/api3-202112092:v1
```

Verificar:

```bash
sudo nerdctl images | grep api3
```

Desde VM2 descargar API1:

```bash
podman pull --tls-verify=false \
192.168.122.239:5000/api1-202112092:v1
```

Verificar:

```bash
podman images | grep api1
```

Estas pruebas demuestran interoperabilidad entre:

```text
containerd ↔ Zot ↔ Podman
```

---

# Prueba de persistencia de Zot

En VM3:

```bash
sudo docker restart zot
```

Después:

```bash
curl http://localhost:5000/v2/_catalog
```

Las imágenes deben continuar apareciendo.

Esto demuestra que el almacenamiento:

```text
~/zot-registry/data
```

es persistente.

---

# Comunicación API1 hacia API2

Desde VM1:

```bash
curl http://localhost:8081/api1/202112092/call-api2
```

Debe devolver:

```text
"connection":true
```

---

# Comunicación API1 hacia API3

```bash
curl http://localhost:8081/api1/202112092/call-api3
```

Debe devolver:

```text
"connection":true
```

---

# Comunicación API2 hacia API1

```bash
curl http://localhost:8082/api2/202112092/call-api1
```

Resultado:

```text
"connection":true
```

---

# Comunicación API2 hacia API3

```bash
curl http://localhost:8082/api2/202112092/call-api3
```

Resultado:

```text
"connection":true
```

---

# Comunicación API3 hacia API1

Desde VM2:

```bash
curl http://localhost:8083/api3/202112092/call-api1
```

Resultado:

```text
"connection":true
```

---

# Comunicación API3 hacia API2

```bash
curl http://localhost:8083/api3/202112092/call-api2
```

Resultado:

```text
"connection":true
```

---

# Script automático de comunicaciones

El repositorio incluye:

```text
scripts/pruebas-comunicacion.sh
```

Desde el equipo anfitrión, ingresar a la carpeta `PROYECTO 1` y ejecutar:

```bash
chmod +x scripts/pruebas-comunicacion.sh
./scripts/pruebas-comunicacion.sh
```

El script realiza:

```text
API1 → API2
API1 → API3

API2 → API1
API2 → API3

API3 → API1
API3 → API2
```

Todas las respuestas deben contener:

```json
"connection": true
```

---

# Script de verificación de Zot

El repositorio también contiene:

```text
scripts/pruebas-registro.sh
```

Ejecutar:

```bash
chmod +x scripts/pruebas-registro.sh
./scripts/pruebas-registro.sh
```

Este script comprueba:

- Disponibilidad de Zot.
- Catálogo.
- API1.
- API2.
- API3.
- Tags almacenados.

---

# Verificación después de reiniciar VM1

Reiniciar VM1.

Después:

```bash
sudo nerdctl ps
```

API1 y API2 deben aparecer activos.

Verificar:

```bash
curl http://localhost:8081/health
curl http://localhost:8082/health
```

Ambas APIs deben continuar funcionando.

---

# Verificación después de reiniciar VM2

Reiniciar VM2.

No iniciar API3 manualmente.

Ejecutar:

```bash
systemctl --user is-active api3.service
```

Debe devolver:

```text
active
```

Verificar:

```bash
podman ps
```

API3 debe aparecer automáticamente.

Health:

```bash
curl http://localhost:8083/health
```

Debe devolver:

```text
"status":"UP"
```

---

# Verificación después de reiniciar VM3

Reiniciar VM3.

Ejecutar:

```bash
sudo docker ps
```

Zot debe aparecer activo.

Consultar:

```bash
curl http://localhost:5000/v2/_catalog
```

Las tres imágenes deben continuar almacenadas.

---

# Prueba final de funcionamiento

Con las tres máquinas virtuales encendidas se debe comprobar:

```text
VM1
├── API1 → UP
└── API2 → UP

VM2
└── API3 → UP

VM3
└── Zot → UP
```

Luego ejecutar:

```bash
./scripts/pruebas-comunicacion.sh
```

Las seis comunicaciones deben devolver:

```text
connection = true
```

Y ejecutar:

```bash
./scripts/pruebas-registro.sh
```

Debe encontrarse:

```text
api1-202112092:v1
api2-202112092:v1
api3-202112092:v1
```

---

# Estructura final del Proyecto 1

```text
PROYECTO 1/
├── README.md
│
├── codigo/
│   ├── api1/
│   │   ├── Dockerfile
│   │   ├── go.mod
│   │   └── main.go
│   │
│   ├── api2/
│   │   ├── Dockerfile
│   │   ├── go.mod
│   │   └── main.go
│   │
│   └── api3/
│       ├── Dockerfile
│       ├── go.mod
│       └── main.go
│
├── configuracion/
│   ├── vm1/
│   │   ├── buildkitd.toml
│   │   ├── buildkit.service
│   │   └── ejecutar-apis.sh
│   │
│   ├── vm2/
│   │   ├── api3.container
│   │   └── configurar-api3.sh
│   │
│   └── vm3/
│       └── ejecutar-zot.sh
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

# Resultado final esperado

La instalación se considera exitosa cuando:

- KVM ejecuta correctamente las tres máquinas virtuales.
- VM1 mantiene API1 y API2 ejecutándose mediante containerd.
- VM2 mantiene API3 mediante Podman.
- API3 inicia automáticamente mediante systemd y Quadlet.
- VM3 mantiene Zot mediante Docker.
- Zot conserva las imágenes después de reinicios.
- Las tres imágenes poseen el tag `v1`.
- Zot permite operaciones Push y Pull.
- Las seis comunicaciones REST funcionan correctamente.
- Las direcciones IP permanecen estables mediante reservas DHCP.

La infraestructura final es:

```text
VM1 - 192.168.122.101
│
├── API1 :8081
└── API2 :8082
       ↕
       ↕ REST
       ↕
VM2 - 192.168.122.22
│
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

Con esta configuración, el Proyecto 1 queda instalado y listo para realizar las pruebas de funcionamiento y demostración.

