create function envio_email_reserva() returns trigger language plpgsql as $$
declare
	--Variables tipo record que guardan filas enteras
    v_cliente record;
    v_vuelo record;
    v_origen record;
    v_destino record;
    v_cuerpo text;
    v_asunto text;
    v_operacion_email text;
begin
    --Buscamos los datos de las tablas
    --Recuperamos info del cliente
    select email, nombre, apellido, dni
    into v_cliente
    from cliente
    where id_cliente = new.id_cliente;

	--Recuperamos info del vuelo y ruta para saber los aeropuertos
    select v.id_vuelo, v.fecha_salida, v.hora_salida, r.nro_ruta,
          r.id_aeropuerto_origen, r.id_aeropuerto_destino
    into v_vuelo
    from vuelo v join ruta r on v.nro_ruta = r.nro_ruta
    where v.id_vuelo = new.id_vuelo;

	--Buscamos nombre y ciudad del aeropuerto de origen y destino
    select id_aeropuerto, nombre, localidad, provincia
    into v_origen
    from aeropuerto
    where id_aeropuerto = v_vuelo.id_aeropuerto_origen;

    select id_aeropuerto, nombre, localidad, provincia
    into v_destino
    from aeropuerto
    where id_aeropuerto = v_vuelo.id_aeropuerto_destino;

    --Validaciones
    -- Primer caso, detecta si se esta creando una nueva reserva
    if (TG_OP = 'INSERT') then
        v_operacion_email := 'reserva de pasaje';
        v_asunto := 'CheapFlights - reserva de pasaje';
        v_cuerpo := format(
            'Cliente: %s %s; DNI: %s; id_reserva: %s; id_vuelo: %s; Origen: %s - %s, %s, %s; Destino: %s - %s, %s, %s; Fecha y hora salida: %s %s',
            v_cliente.nombre, v_cliente.apellido,
            v_cliente.dni,
            new.id_reserva,
            v_vuelo.id_vuelo,
            v_origen.id_aeropuerto, v_origen.nombre, v_origen.localidad, v_origen.provincia,
            v_destino.id_aeropuerto, v_destino.nombre, v_destino.localidad, v_destino.provincia,
            v_vuelo.fecha_salida, v_vuelo.hora_salida
        );

    -- Segundo caso(check-in), de reservado tiene que pasar a confirmado, a traves de UPDATE.
    elsif (TG_OP = 'UPDATE' and new.estado = 'confirmado' and old.estado is distinct from 'confirmado') then
        v_operacion_email := 'check-in de asiento';
        v_asunto := 'CheapFlights - check-in de asiento';
        v_cuerpo := format(
            'Cliente: %s %s; DNI: %s; id_reserva: %s; id_vuelo: %s; Asiento: %s; Origen: %s - %s, %s, %s; Destino: %s - %s, %s, %s; Fecha y hora salida: %s %s',
            v_cliente.nombre, v_cliente.apellido,
            v_cliente.dni,
            new.id_reserva,
            v_vuelo.id_vuelo,
            new.nro_asiento,
            v_origen.id_aeropuerto, v_origen.nombre, v_origen.localidad, v_origen.provincia,
            v_destino.id_aeropuerto, v_destino.nombre, v_destino.localidad, v_destino.provincia,
            v_vuelo.fecha_salida, v_vuelo.hora_salida
        );

    -- Tercer caso (anulacion), de reservado tiene que pasar a anulado, a traves de UPDATE.
    elsif (TG_OP = 'UPDATE' and new.estado = 'anulado' and old.estado is distinct from 'anulado') then
        v_operacion_email := 'anulación de reserva';
        v_asunto := 'Cheapflights - anulación de reserva';
        v_cuerpo := format(
            'Cliente: %s %s; DNI: %s; id_reserva: %s; id_vuelo: %s; Origen: %s - %s, %s, %s; Destino: %s - %s, %s, %s; Fecha y hora salida: %s %s',
            v_cliente.nombre, v_cliente.apellido,
            v_cliente.dni,
            new.id_reserva,
            v_vuelo.id_vuelo,
            v_origen.id_aeropuerto, v_origen.nombre, v_origen.localidad, v_origen.provincia,
            v_destino.id_aeropuerto, v_destino.nombre, v_destino.localidad, v_destino.provincia,
            v_vuelo.fecha_salida, v_vuelo.hora_salida
        );
    else
        return new;
    end if;

    --Se inserta en la tabla envio_email cualquier novedad. (RESERVADO, CONFIRMADO, ANULADO)
    insert into envio_email(f_generacion, email_cliente, asunto, cuerpo, estado)
    values (now(), v_cliente.email, v_asunto, v_cuerpo, 'pendiente');

    return new;
end;
$$;


create trigger trg_envio_email_reserva after
insert
or
update on reserva_pasaje
for each row execute function envio_email_reserva();
