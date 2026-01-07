alter table ruta add constraint fk_id_aeropuerto_origen foreign key (id_aeropuerto_origen) references aeropuerto (id_aeropuerto); --fk ruta
alter table ruta add constraint fk_id_aeropuerto_destino foreign key (id_aeropuerto_destino) references aeropuerto (id_aeropuerto);

alter table vuelo add constraint fk_nro_ruta foreign key (nro_ruta) references ruta (nro_ruta); --fk vuelo

alter table reserva_pasaje add constraint fk_id_vuelo foreign key (id_vuelo) references vuelo (id_vuelo); --fk pasaje
alter table reserva_pasaje add constraint fk_id_cliente foreign key (id_cliente) references cliente (id_cliente);

