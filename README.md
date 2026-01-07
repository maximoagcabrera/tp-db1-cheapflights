##CheapFlights: Sistema de Gestión de Vuelos

##Descripción del Proyecto

CheapFlights es una plataforma robusta diseñada para la administración de operaciones aeroportuarias. El núcleo del sistema reside en una arquitectura de persistencia 
híbrida que garantiza la integridad relacional mediante PostgreSQL y la agilidad de almacenamiento local con BoltDB (Key-Value), todo orquestado por una aplicación central desarrollada en Go.

##Características Principales (Core Features)
- Persistencia Híbrida: Sincronización de datos entre una base de datos relacional (SQL) y una NoSQL orientada a documentos (BoltDB).
- Lógica de Negocio en Base de Datos: Implementación de reglas críticas mediante Stored Procedures que gestionan la apertura de vuelos, reservas y check-in.
- Automatización con Triggers: Sistema de "notificaciones" reactivo que detecta cambios en las reservas y genera automáticamente registros de auditoría/envío de emails.
- Sistema de Auditoría de Errores: Tabla dedicada para el loggeo automático de errores de negocio (vuelos vencidos, falta de asientos, IDs inválidos).
- Motor de Pruebas Automatizado: Ingesta de archivos JSON para simular cargas de datos reales y ejecutar flujos de prueba secuenciales.

##Stack Tecnológico
- Backend: Go (Golang).
- Relational DB: PostgreSQL (SQL).
- NoSQL DB: BoltDB.
- Data Format: JSON.

##Estructura del Repositorio
sql/: Colección de scripts modulares para la base de datos:  
  tablas/: Definición de esquemas para aeropuertos, rutas, clientes y vuelos.  
  stored_procedures/: Lógica programada en PLpgSQL (apertura_vuelo, anular_reserva, check_in_asiento).  
  triggers/: Automatización de tareas post-transacción.  
data/: Datasets JSON con aeropuertos de Argentina, clientes y escenarios de prueba.  
main.go: Interfaz de consola interactiva y orquestador de conexiones.  
BoltDB.go: Capa de persistencia NoSQL y manejo de transacciones atómicas.  

##Ejemplo de Lógica Implementada
El sistema no solo guarda datos, sino que los valida en tiempo real:
"Antes de confirmar una reserva, el sistema verifica en milisegundos la existencia del cliente, la vigencia del vuelo y la disponibilidad de asientos, disparando un 
registro de error automático si alguna condición falla".

##Instalación y Ejecución
1. Clonar el repositorio.
2. Configurar PostgreSQL localmente.
4. Compilar y correr el orquestador:
5. Ejecutar con Bash

##Equipo de Desarrollo
Cabral Tobias  
Cabrera Maximo  
Fauda Matias  
Inti Goñi  
