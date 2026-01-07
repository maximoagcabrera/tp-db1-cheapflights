alter table cliente drop constraint if exists pk_cliente;
alter table aeropuerto drop constraint if exists pk_aeropuerto;
alter table ruta drop constraint if exists pk_ruta;
alter table vuelo drop constraint if exists pk_vuelo;
alter table reserva_pasaje drop constraint if exists pk_reserva_pasaje;
alter table error drop constraint if exists pk_error;
alter table envio_email drop constraint if exists pk_envio_email;
alter table reserva_pasaje drop constraint if exists pk_datos_de_prueba;
