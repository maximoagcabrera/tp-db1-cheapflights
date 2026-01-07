CREATE FUNCTION apertura_vuelo(p_nro_ruta int, p_f_salida_vuelo timestamp, p_nro_asientos_totales int) RETURNS int LANGUAGE PLPGSQL AS $$
	declare
	var_id_vuelo int;
	begin
	--Valido que la ruta exista--
	if not exists(
		select 1
		from ruta
		where nro_ruta = p_nro_ruta
		)
		then
		insert into error(
		operacion, nro_ruta, f_salida_vuelo, nro_asientos_totales, f_error,motivo
		)
		values(
		'apertura', p_nro_ruta, p_f_salida_vuelo, p_nro_asientos_totales, now(), 'numero de ruta no valido'
		);
		return -1;
	end if;

	--Valido fecha y hora posterior a la actual--
	if p_f_salida_vuelo <= now()
	then
	insert into error(
	operacion, nro_ruta, f_salida_vuelo, nro_asientos_totales, f_error,motivo
	)
	values(
    'apertura', p_nro_ruta, p_f_salida_vuelo, p_nro_asientos_totales, now(), 'no se permite abrir un vuelo con retroactividad'
	);
	return -1;
end if;

	--Validacion de asientos--
	if p_nro_asientos_totales <= 0
	then
	insert into error(
    operacion, nro_ruta, f_salida_vuelo, nro_asientos_totales, f_error,motivo
    )
    values(
    'apertura', p_nro_ruta, p_f_salida_vuelo, p_nro_asientos_totales, now(), 'no se permite abrir un vuelo sin asientos disponibles'
    );
    return -1;
    end if;

    --Con las validaciones hechas, inserto vuelo--
    insert into vuelo (
    nro_ruta,
    fecha_salida,
    hora_salida,
    nro_asientos_totales,
    nro_asientos_disponibles
    )
    values(
    p_nro_ruta,
    p_f_salida_vuelo::date,
    p_f_salida_vuelo::time,
    p_nro_asientos_totales,
    p_nro_asientos_totales
    )
    returning id_vuelo into var_id_vuelo;
    return var_id_vuelo;
  end;
  $$;
