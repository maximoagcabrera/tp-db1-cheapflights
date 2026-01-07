CREATE FUNCTION check_in_asiento(p_id_reserva int, p_id_cliente int, p_nro_asiento int) RETURNS boolean LANGUAGE PLPGSQL AS $$
declare
	v_id_vuelo int;
	v_estado_actual varchar(20);
	v_asientos_totales int;
begin
	--Intentamos recuperar el vuelo y estado actual
	select rp.id_vuelo, rp.estado
	into v_id_vuelo, v_estado_actual
	from reserva_pasaje rp
	where rp.id_reserva = p_id_reserva and rp.id_cliente = p_id_cliente;

	--Si la consulta anterior no trae nada, los datos son incorrectos
	if not found then
		insert into error (operacion, id_reserva, id_cliente, f_error, motivo)
		values ('check-in', p_id_reserva, p_id_cliente, now(), 'id de reserva no válido');
		return false;
	end if;
	
	--Solo permitimos check-in si el estado es reservado
	if v_estado_actual != 'reservado' then
		insert into error (operacion, id_reserva, id_cliente, f_error, motivo)
		values ('check-in', p_id_reserva, p_id_cliente, now(), 'check-in ya realizado para el id de reserva');
		return false;
	end if;

	--Obtenemos la capacidad del avion para validar el numero de asiento
	select nro_asientos_totales into v_asientos_totales
	from vuelo where id_vuelo = v_id_vuelo;

	--El asiento tiene que estar dentro del rango
	if p_nro_asiento <= 0 OR p_nro_asiento > v_asientos_totales then
		insert into error (operacion, id_reserva, id_cliente, nro_asiento, f_error, motivo)
		values ('check-in', p_id_reserva, p_id_cliente, p_nro_asiento, now(), 'número de asiento inexistente');
		return false;
	end if;

	--Verifico si alguien mas ya confirmo ese asiento en este vuelo
	if exists (
		select 1
		from reserva_pasaje rp_ocupada
		where rp_ocupada.id_vuelo = v_id_vuelo
		and rp_ocupada.nro_asiento = p_nro_asiento
		and rp_ocupada.estado = 'confirmado'
	) then
		insert into error (operacion, id_reserva, id_cliente, nro_asiento, f_error, motivo)
		values ('check-in', p_id_reserva, p_id_cliente, p_nro_asiento, now(), 'numero de asiento ya ocupado');
		return false;
	end if;


	--Si paso todo lo anterior, realizamos el update
	update reserva_pasaje
	set
		nro_asiento = p_nro_asiento,
		f_check_in = now(),
		estado = 'confirmado'
	where id_reserva = p_id_reserva;

	return true;
end;
$$;
