#!/bin/bash

# Proyecto Sistemas Operativos 1
# VM1 - containerd / nerdctl
# Carnet: 202112092

set -e

echo "======================================"
echo " Iniciando API1 y API2 en VM1"
echo "======================================"

# Crear la red si todavía no existe
if ! sudo nerdctl network inspect red-apis-vm1 >/dev/null 2>&1; then
    echo "Creando red red-apis-vm1..."
    sudo nerdctl network create red-apis-vm1
fi

# Eliminar contenedores anteriores si existen
sudo nerdctl rm -f api1 >/dev/null 2>&1 || true
sudo nerdctl rm -f api2 >/dev/null 2>&1 || true

echo "Iniciando API1..."

sudo nerdctl run -d \
  --name api1 \
  --restart=unless-stopped \
  --network red-apis-vm1 \
  -p 8081:8080 \
  -e CARNET=202112092 \
  -e API2_URL="http://api2:8080" \
  -e API3_URL="http://192.168.122.22:8083" \
  api1-202112092:v1

echo "Iniciando API2..."

sudo nerdctl run -d \
  --name api2 \
  --restart=unless-stopped \
  --network red-apis-vm1 \
  -p 8082:8080 \
  -e CARNET=202112092 \
  -e API1_URL="http://api1:8080" \
  -e API3_URL="http://192.168.122.22:8083" \
  api2-202112092:v1

echo
echo "Contenedores activos:"
sudo nerdctl ps
