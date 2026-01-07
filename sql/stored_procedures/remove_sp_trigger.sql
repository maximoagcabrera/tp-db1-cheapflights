-- Borro funciones si ya existen

DROP FUNCTION IF EXISTS apertura_vuelo(int, timestamp, int);


DROP FUNCTION IF EXISTS reserva_pasaje(int, int);


DROP FUNCTION IF EXISTS check_in_asiento(int, int, int);


DROP FUNCTION IF EXISTS anular_reserva(int, int);


DROP FUNCTION IF EXISTS envio_email_reserva();


DROP TRIGGER IF EXISTS trg_envio_email_reserva ON reserva_pasaje;


DROP FUNCTION IF EXISTS iniciar_pruebas();
