CREATE FUNCTION reserva_pasaje(i_id_vuelo int, i_id_cliente int) RETURNS int --id_reserva
 LANGUAGE PLPGSQL AS $$
declare
    var_id_reserva int;
begin

	--Valido id_vuelo
	if not exists(
		select 1
		from vuelo v
		where v.id_vuelo = i_id_vuelo
		) then
			insert into error(
   				operacion,
				id_vuelo,
				id_cliente,
				f_error,
				motivo
    		)
			VALUES (
            	'reservar',
				i_id_vuelo,
				i_id_cliente,
				now(),
				'id de vuelo no valido.'
            );
			return -1;
	end if;

	--Valido id_cliente
	if not exists(
        select 1
        from cliente c
        where c.id_cliente  = i_id_cliente
   		 ) then
            insert into error(
                operacion,
				id_vuelo,
				id_cliente,
				f_error,
				motivo
            )
            VALUES (
                'reservar',
				i_id_vuelo,
				i_id_cliente,
				now(),
				'id de cliente no valido.'
            );
            return -1;
    end if;

	--Valido nro_asientos_disponibles
	if exists(
        select 1
        from vuelo v
        where v.id_vuelo = i_id_vuelo
			and v.nro_asientos_disponibles < 1
    	) then
            insert into error(
                operacion,
				id_vuelo,
				id_cliente,
				nro_asientos_totales,
				f_error,
				motivo
            )
            values (
                'reservar',
				i_id_vuelo,
				i_id_cliente,
				(select nro_asientos_disponibles
					from vuelo
					where id_vuelo = i_id_vuelo
				),
				now(),
				'El vuelo ya esta completo'
			);
			return -1;
	else
		--Resto 1 a  nro_asientos_disponibles
		update vuelo
		set nro_asientos_disponibles = nro_asientos_disponibles - 1
		where id_vuelo = i_id_vuelo;
	end if;

	--Guardo la reserva
    insert into reserva_pasaje (
        id_vuelo,
		id_cliente,
		f_reserva,
		estado
    )
    values (
        i_id_vuelo,
		i_id_cliente,
		now(),
		'reservado'
    )
	returning id_reserva into var_id_reserva; --Guarda el valor de id_reserva en var_id_reserva

	return var_id_reserva;
  end;
  $$;
