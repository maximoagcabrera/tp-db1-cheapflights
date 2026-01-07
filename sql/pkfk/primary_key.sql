alter table cliente add constraint pk_cliente primary key (id_cliente);
alter table aeropuerto add constraint pk_aeropuerto primary key (id_aeropuerto);
alter table ruta add constraint pk_ruta primary key (nro_ruta);
alter table vuelo add constraint pk_vuelo primary key (id_vuelo);
alter table reserva_pasaje add constraint pk_reserva_pasaje primary key (id_reserva);
alter table error add constraint pk_error primary key (id_error);
alter table envio_email add constraint pk_envio_email primary key (id_email);
alter table datos_de_prueba add constraint pk_datos_de_prueba primary key (id_orden);
