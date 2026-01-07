create table error (
		id_error serial,
		operacion char(10),
		nro_ruta int,
		f_salida_vuelo timestamp,
		nro_asientos_totales int,
		id_vuelo int,
		id_cliente int,
		id_reserva int,
		nro_asiento int,
		f_error timestamp,
		motivo varchar(80)
	);
