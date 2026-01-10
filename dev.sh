#!/bin/bash

# Detener cualquier proceso que esté usando el puerto 8090
echo "Verificando si hay procesos usando el puerto 8090..."
PID=$(lsof -ti:8090)
if [ ! -z "$PID" ]; then
  echo "Deteniendo proceso $PID que está usando el puerto 8090..."
  kill $PID
fi

# Ejecutar la aplicación con Air
echo "Iniciando JVAIRV2 con hot-reload en puerto 8090..."
/Users/eduardo/go/bin/air
