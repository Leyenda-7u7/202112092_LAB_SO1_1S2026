#!/bin/bash

# Proyecto Sistemas Operativos 1
# VM2 - Podman Quadlet / systemd
# Carnet: 202112092

set -e

echo "======================================"
echo " Configurando API3 con Podman Quadlet"
echo "======================================"

mkdir -p "$HOME/.config/containers/systemd"

cp api3.container "$HOME/.config/containers/systemd/api3.container"

# Permitir servicios del usuario aunque no exista sesión SSH abierta
sudo loginctl enable-linger "$USER"

systemctl --user daemon-reload
systemctl --user start api3.service

echo
echo "Estado de API3:"
systemctl --user status api3.service --no-pager

echo
echo "Contenedores:"
podman ps
