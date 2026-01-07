create function anular_reserva(p_id_reserva int, p_id_cliente int) returns boolean language plpgsql as $$
declare
	v_estado char(10);
	v_id_vuelo int;
begin
	--Validamos que la reserva exista y pertenezca al cliente
	select estado, id_vuelo into v_estado, v_id_vuelo
	from reserva_pasaje
	where id_reserva = p_id_reserva and id_cliente = p_id_cliente;

	--Si no encuentro ninguna fila, v_estado y v_id_vuelo == null
	--Inserto error de que la reserva no existe
	if v_id_vuelo is null then
		insert into error(operacion, id_reserva, id_cliente, f_error, motivo)
		values ('anulacion', p_id_reserva, p_id_cliente, now(), 'id de reserva no valido');
		return false;
	end if;

	--Valido que el estado sea reservado
	if v_estado != 'reservado' then
	insert into error(operacion, id_reserva, id_cliente, f_error, motivo)
		values ('anulacion', p_id_reserva, p_id_cliente, now(), 'No se puede anular una reserva ya confirmada');
		return false;
	end if;

	--Si todo es valido, actualizo resrva a anulado
	update reserva_pasaje
	set estado = 'anulado'
	where id_reserva = p_id_reserva;

	--Aumento los asientos disponibles
	update vuelo
	set nro_asientos_disponibles = nro_asientos_disponibles + 1
	where id_vuelo = v_id_vuelo;

	return true;
end;
$$;
