create table envio_email (
		id_email serial,
		f_generacion timestamp,
		email_cliente text,
		asunto text,
		cuerpo text,
		f_envio timestamp,
		estado char(10)
	);
