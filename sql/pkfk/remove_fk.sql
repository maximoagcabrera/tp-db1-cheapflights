alter table ruta drop constraint if exists fk_id_aeropuerto_origen;
alter table ruta drop constraint if exists fk_id_aeropuerto_destino;
alter table vuelo drop constraint if exists fk_nro_ruta;
alter table reserva_pasaje drop constraint if exists fk_id_vuelo;
alter table reserva_pasaje drop constraint if exists fk_id_cliente;
