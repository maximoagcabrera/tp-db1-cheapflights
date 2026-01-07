CREATE FUNCTION iniciar_pruebas() RETURNS void LANGUAGE PLPGSQL AS $$
	declare
    fila record;
begin
    raise notice 'Iniciando pruebas:';

    for fila in (select * from datos_de_prueba order by id_orden)
    loop

        if fila.operacion = 'apertura' then
            perform apertura_vuelo(
                fila.nro_ruta,
                fila.f_salida_vuelo,
                fila.nro_asientos_totales
            );

        elsif fila.operacion = 'reserva' then
            perform reserva_pasaje(
                fila.id_vuelo,
                fila.id_cliente
            );

        elsif fila.operacion = 'check-in' then
            perform check_in_asiento(
                fila.id_reserva,
                fila.id_cliente,
                fila.nro_asiento
            );

        elsif fila.operacion = 'anulacion' then
            perform anular_reserva(
                fila.id_reserva,
                fila.id_cliente
            );

        end if;

    end loop;

    raise notice 'Pruebas completadas';
end;
$$;
