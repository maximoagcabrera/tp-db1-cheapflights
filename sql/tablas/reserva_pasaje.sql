create table reserva_pasaje (
		id_reserva serial,
		id_vuelo int,
		id_cliente int,
		f_reserva timestamp,
		nro_asiento int,
		f_check_in timestamp,
		estado CHAR(10)
	);
