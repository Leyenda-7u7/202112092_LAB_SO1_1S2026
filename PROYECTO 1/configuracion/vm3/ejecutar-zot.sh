#!/bin/bash

# Proyecto Sistemas Operativos 1
# VM3 - Docker + Zot Registry
# Carnet: 202112092

set -e

echo "======================================"
echo " Iniciando registro privado Zot"
echo "======================================"

mkdir -p "$HOME/zot-registry/data"

sudo docker pull ghcr.io/project-zot/zot:latest

# Eliminar una instancia anterior si existe
sudo docker rm -f zot >/dev/null 2>&1 || true

sudo docker run -d \
  --name zot \
  --restart unless-stopped \
  -p 5000:5000 \
  -v "$HOME/zot-registry/data:/var/lib/registry" \
  ghcr.io/project-zot/zot:latest

echo
echo "Contenedor Zot:"
sudo docker ps --filter name=zot

echo
echo "Catálogo:"
curl -s http://localhost:5000/v2/_catalog
echo
