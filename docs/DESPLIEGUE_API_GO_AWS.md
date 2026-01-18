# 🚀 Despliegue y acceso a API Go en AWS usando Bastion + SSH Tunnel

Este documento describe paso a paso cómo:
- Compilar y ejecutar una API en Go en una EC2 privada
- Acceder a la API desde macOS mediante un túnel SSH
- Mantener la API aislada y sin afectar aplicaciones existentes (Laravel)

---

## 📌 Prerrequisitos

- macOS
- Acceso a AWS Console
- Bastion EC2 con IP pública
- EC2 privada donde corre la API Go
- Archivo `bastion_migration.pem`
- Go instalado en la EC2 del API
- Proyecto ubicado en `~/JVAIRV2`
- API configurada para escuchar en el puerto `8080`

---

## 🧩 PARTE A — Levantar la API Go en AWS

### 1️⃣ Prender la instancia
Desde **AWS Console → EC2 → Instances**, verificar que la instancia del API esté en estado **Running**.

### 2️⃣ Conectarse al Bastion (desde macOS)

```bash
ssh -i bastion_migration.pem ec2-user@IP_PUBLICA_BASTION
```

Ejemplo:

```bash
ssh -i bastion_migration.pem ec2-user@44.204.136.89
```

### 3️⃣ Entrar a la carpeta del proyecto

```bash
cd ~/JVAIRV2
```

### 5️⃣ Compilar la aplicación Go

```bash
GOMAXPROCS=1 go build -o api ./cmd/api
```

### 6️⃣ Ejecutar la API

```bash
./api
```

---

## 🌐 PARTE B — Crear el túnel SSH desde macOS

### 7️⃣ Abrir una nueva terminal en macOS

### 8️⃣ Crear el túnel SSH

```bash
ssh -N -L 8080:10.0.1.148:8080 -i bastion_migration.pem ec2-user@44.204.136.89
```

---

## 🧪 PARTE C — Consumir la API desde macOS

```bash
curl http://localhost:8080
```

---

## 🛑 Cerrar conexiones

```bash
Ctrl + C
```

---

## 🧠 Resumen

La API corre en una EC2 privada y se accede únicamente mediante un túnel SSH desde macOS usando el Bastion.
