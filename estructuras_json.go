package main

//----------------- Estrucutras para JSON's -----------------------

type Cliente struct {
	IDCliente       int    `json:"id_cliente"`
	Nombre          string `json:"nombre"`
	Apellido        string `json:"apellido"`
	DNI             int    `json:"dni"`
	FechaNacimiento string `json:"fecha_nacimiento"`
	Telefono        string `json:"telefono"`
	Email           string `json:"email"`
}

type Aeropuerto struct {
	IDAeropuerto string `json:"id_aeropuerto"` // char(3)
	Nombre       string `json:"nombre"`
	Localidad    string `json:"localidad"`
	Provincia    string `json:"provincia"`
}

type Ruta struct {
	NroRuta             int    `json:"nro_ruta"`
	IDAeropuertoOrigen  string `json:"id_aeropuerto_origen"`  // char(3)
	IDAeropuertoDestino string `json:"id_aeropuerto_destino"` // char(3)
	Duracion            string `json:"duracion"`              // interval representado como texto
}

type DatoPrueba struct {
	IDOrden            int    `json:"id_orden"`
	Operacion          string `json:"operacion"`
	NroRuta            int    `json:"nro_ruta"`
	FSalidaVuelo       string `json:"f_salida_vuelo"`
	NroAsientosTotales int    `json:"nro_asientos_totales"`
	IDVuelo            int    `json:"id_vuelo"`
	IDCliente          int    `json:"id_cliente"`
	IDReserva          int    `json:"id_reserva"`
	NroAsiento         int    `json:"nro_asiento"`
}
