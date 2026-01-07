package main

//----------------------Bibliotecas--------------------

import (
	"os"
	"fmt"
	"log"
	"time"
	"encoding/json"
	"go.etcd.io/bbolt"
)

//----------------------Structs para cargar-------------------

type Vuelo struct {
	NroRuta				int
	FSalidaVuelo		string
	NroAsientosTotales	int
}

type Reserva struct {
	IDVuelo				int
	IDCliente			int
	FReserva			string
	Estado				string
}

type Error struct {
	descripcion			string
}

//----------------------Conectar con DB-----------------------
// OpenBoltDB?
func Open_boltDB() *bbolt.DB {
	const boltDB_var = "boltDB_data.db"
	// 1) Borrar la db si existe
	borrarBaseDeDatos()

	// 2) Crear
	db, err := bbolt.Open(boltDB_var, 0600, nil)
	if err != nil {
		log.Fatal("Error al borrar el archivo %s: %w", boltDB_var, err)
	}
	fmt.Println("BoltDB inicializado correctamente")

	// 3) Orden para ejecutar las funciones
	err = exec_boltDB(db)
	if err != nil {
        db.Close()
        log.Fatal("Error ejecutando:", err)
	}
	return db
}

func borrarBaseDeDatos() error {
    const boltDB_var = "boltDB_data.db"

    //Verifica si el archivo existe
    _, err := os.Stat(boltDB_var)
    if os.IsNotExist(err) {
        fmt.Println("La base de datos ya está borrada.")
        return nil
    }

    //Si existe, lo borra
    err = os.Remove(boltDB_var)
    if err != nil {
        return fmt.Errorf("Error al borrar el archivo %s: %w", boltDB_var, err)
    }

	fmt.Println("Base de datos borrada")
    return nil
}


//--------------------Ejecucion----------------------

func exec_boltDB(db *bbolt.DB) error {

	//Variables
	var aeropuertoB []Aeropuerto
	var clienteB []Cliente
	var rutaB []Ruta
	var pruebaB []DatoPrueba

	//Creamos buckets
	err := crearBuckets(db)
	if err != nil {
		return err
	}

	//Cargamos datos Json
	//Aeropuerto
	err = cargarJSON("data/aeropuertos.json", &aeropuertoB)
	if err != nil{
		return err
	}
	err = cargarAeropuertos(aeropuertoB, db)
    if err != nil {
		return err
	}

	//Cliente
	err = cargarJSON("data/clientes.json", &clienteB)
	if err != nil{
		return err
	}
	err = cargarClientes(clienteB, db)
    if err != nil {
		return err
	}

	//Ruta
	err = cargarJSON("data/rutas.json", &rutaB)
	if err != nil{
		return err
	}
	err = cargarRutas(rutaB, db)
    if err != nil {
		return err
	}

	//Datos de prueba
	err = cargarJSON("data/datos_de_prueba.json", &pruebaB)
	if err != nil{
		return err
	}
	err = cargarDatosDePrueba(pruebaB, db)
    if err != nil {
		return err
	}

    fmt.Println("Ejecución completada")
	return nil
}

//---------------------Buckets-------------------------

func crearBuckets(db *bbolt.DB) error {
	var table_name = []string{"Cliente", "Aeropuerto",
							  "Ruta", "Vuelo",
							  "Reserva Pasaje", "Envìo mail",
							  "Error", "Datos de prueba"}

	//Abro una transaccion, si falla volvemos para atras
	err := db.Update(func(tx *bbolt.Tx) error {

		//Creo los bucckets si no existen
		for _, name := range table_name {
			_, err := tx.CreateBucketIfNotExists([]byte(name))

			if err != nil {
				return err
			}
		}
		fmt.Println("Buckets creados")
		return nil
	})
    return err
}

//----------------Cargar JSON-----------------

func cargarJSON(archive string, destine interface{}) error {
	data, err := os.ReadFile(archive)
	if err != nil {
		log.Fatal("Error abriendo el archivo %s: %v", archive, err)
	}
	//Unmarshal convierte bytes a Structs
	if err := json.Unmarshal(data, destine); err != nil {
		log.Fatal("Error decodificando el JSON %v", err)
	}
	return nil
}

//---------------Insertar datos JSON----------------------

func cargarAeropuertos(list []Aeropuerto, db *bbolt.DB) error {
    b_name := "Aeropuerto"

	//Abro una transaccion
    err := db.Update(func(tx *bbolt.Tx) error {
        b := tx.Bucket([]byte(b_name))
        if b == nil {
            return fmt.Errorf("Bucket '%s' no existe", b_name)
        }

		//Recorro aeropuerto
        for _, a := range list {
            value, err := json.Marshal(a)
            if err != nil {
                return fmt.Errorf("Error con leyendo aeropuerto %s: %w", a.IDAeropuerto, err)
            }
            key := []byte(a.IDAeropuerto)

            //Inserto Key y Value
            err = b.Put(key, value)
            //fmt.Println("Key: %s Value: %v\n", string(key), value)
            if err != nil {
                return fmt.Errorf("Error con insertando Aeropeurto %s: %w", a.IDAeropuerto, err)
            }
        }
        fmt.Println("Aeropuertos cargados correctamente")
        return nil
    })
    return err
}

func cargarClientes(list []Cliente, db *bbolt.DB) error {
    b_name := "Cliente"

	//Abro una transaccion
    err := db.Update(func(tx *bbolt.Tx) error {
        b := tx.Bucket([]byte(b_name))
        if b == nil {
            return fmt.Errorf("Bucket '%s' no existe", b_name)
        }

		//Recorro cliente
        for _, c := range list {
            value, err := json.Marshal(c)
            if err != nil {
                return fmt.Errorf("Error leyendo cliente %s: %w", c.IDCliente, err)
            }
            key := []byte(fmt.Sprintf("%d", c.IDCliente)) //Convierto int en byte

            //Inserto Key y Value
            err = b.Put(key, value)
            //fmt.Println("Key: %s Value: %v\n", string(key), value)
            if err != nil {
                return fmt.Errorf("Error insertando Cliente %s: %w", c.IDCliente, err)
            }
        }
        fmt.Println("Clientes cargados correctamente")
        return nil
    })
    return err
}

func cargarRutas(list []Ruta, db *bbolt.DB) error {
    b_name := "Ruta"

	//Abro una transaccion
    err := db.Update(func(tx *bbolt.Tx) error {
        b := tx.Bucket([]byte(b_name))
        if b == nil {
            return fmt.Errorf("Bucket '%s' no existe", b_name)
        }

		//Recorro ruta
        for _, r := range list {
            value, err := json.Marshal(r)
            if err != nil {
                return fmt.Errorf("Error leyendo ruta %s: %w", r.NroRuta, err)
            }
            key := []byte(fmt.Sprintf("%d", r.NroRuta))

            //Inserto Key y Value
            err = b.Put(key, value)
            //fmt.Println("Key: %s Value: %v\n", string(key), value)
            if err != nil {
                return fmt.Errorf("Error insertando Ruta %s: %w", r.NroRuta, err)
            }
        }
        fmt.Println("Rutas cargadas correctamente")
        return nil
    })
    return err
}


//Procesa las operaciones que necesitan validacion logica
func cargarDatosDePrueba(list []DatoPrueba, db *bbolt.DB) error {
	b_name := "Datos de prueba"

	//Abro una transaccion, el bucle corre dentro de una transaccion grande
    err := db.Update(func(tx *bbolt.Tx) error {
        b := tx.Bucket([]byte(b_name))
        if b == nil {
            return fmt.Errorf("Bucket '%s' no existe", b_name)
        }

        //Recorro datos de prueba
        for _, d := range list {
			//Vuelo
            if d.Operacion == "apertura" {
                err := insertarVuelo(d.NroRuta, d.FSalidaVuelo, d.NroAsientosTotales, tx)
                if err != nil {
                    return err
                }
            }

            //Reserva
			if d.Operacion == "reserva" {
					err := insertarReserva(d.IDVuelo, d.IDCliente, tx)
					if err != nil {
						return err
					}
			}
        }

		fmt.Println("Datos de prueba procesados correctamente")
		return nil
	})
	return err
}

//Crear vuelo
func insertarVuelo(nro_ruta int, f_salida_vuelo string, nro_asientos_totales int, tx *bbolt.Tx) error {
	b_name := "Vuelo"
	b_ruta := "Ruta"

	bVuelo := tx.Bucket([]byte(b_name))
	bRuta := tx.Bucket([]byte(b_ruta))

	if bVuelo == nil || bRuta == nil {
		return addError("Bucket Vuelo o Error no existen en la transacción", tx)
	}

	//Validacion 1: ver si la ruta existe
	if bRuta.Get([]byte(fmt.Sprintf("%d", nro_ruta))) == nil {
		msgError := fmt.Sprintf("Numero de ruta invalido: %d no encontrado", nro_ruta)
		return addError(msgError, tx)
	}

	//Validacion 2: formato de fecha
	formato := "2006-01-02 15:04"
    fechaSalida, err := time.Parse(formato, f_salida_vuelo) //Parse convierte string a time
	if err != nil {
		return addError(fmt.Sprintf("Formato de fecha de salida invalido: %v", err), tx)
	}

	now := time.Now()
	if fechaSalida.Before(now) {
		errorMsg := fmt.Sprintf("No se permite abrir un nuevo vuelo con retroactividad. Fecha: %s", f_salida_vuelo)
		return addError(errorMsg, tx)
	}

	//Validacion 3: asientos
	if nro_asientos_totales <= 0 {
		return addError("No se permite abrir un vuelo sin asientos disponibles", tx)
	}

	//Generamos key
	idVuelo, _ := bVuelo.NextSequence() // NextSequence() --> serial de go
	key := []byte(fmt.Sprintf("%d", idVuelo))

	//Value
	data := Vuelo{
		NroRuta:          		nro_ruta,
		FSalidaVuelo:     		f_salida_vuelo,
		NroAsientosTotales: 	nro_asientos_totales,
	}

	//Serializacion a json
	value, err := json.Marshal(data) //json.Marshal --> convierte data en bytes
	if err != nil {
		return fmt.Errorf("Error transformando value a bytes")
	}

	//Insert
	if err := bVuelo.Put(key, value); err != nil {
		return fmt.Errorf("Error al insertar el vuelo con ID %d: %w", idVuelo, err)
	}

	fmt.Printf("Vuelo con ID %d insertado correctamente.\n", idVuelo)

	return nil
}

func insertarReserva(id_vuelo int, id_cliente int, tx *bbolt.Tx) error {
	b_name := "Reserva Pasaje"
	b_cliente := "Cliente"
	b_vuelo := "Vuelo"

	bReserva := tx.Bucket([]byte(b_name))
	bCliente := tx.Bucket([]byte(b_cliente))
	bVuelo := tx.Bucket([]byte(b_vuelo))

	//Validaciones de existencia de tablas
	if bVuelo == nil || bCliente == nil || bReserva == nil {
		return addError("Bucket Vuelo, CLiente o Reserva Pasaje no existen en la transacción", tx)
	}

	//Validacion 1: existencia vuelo
	if bVuelo.Get([]byte(fmt.Sprintf("%d", id_vuelo))) == nil {
		msgError := fmt.Sprintf("IDVuelo: %d no disponible", id_vuelo)
		return addError(msgError, tx)
	}

	//Validacion 2: existencia cliente
	if bCliente.Get([]byte(fmt.Sprintf("%d", id_cliente))) == nil {
		msgError := fmt.Sprintf("IDCliente: %d no encontrado", id_cliente)
		return addError(msgError, tx)
	}

	//Fecha
	formato := "2006-01-02 15:04"
	fecha := time.Now()
    fechaString := fecha.Format(formato)

	//Estado
	state := "reservado"

	//Key
	idReserva, _ := bReserva.NextSequence() // NextSequence() --> serial de go
	key := []byte(fmt.Sprintf("%d", idReserva))

	//Construccion del objeto
	data := Reserva{
		IDVuelo:		id_vuelo,
		IDCliente:		id_cliente,
		FReserva:		fechaString,
		Estado:			state,
	}
	value, err := json.Marshal(data) //json.Marshal --> convierte data en bytes
	if err != nil {
		return fmt.Errorf("Error transformando value a bytes")
	}

	//Insert
	if err := bReserva.Put(key, value); err != nil {
		return fmt.Errorf("Error al insertar el reserva %d: %w", idReserva, err)
	}

	fmt.Printf("Reserva con ID %d insertado correctamente.\n", idReserva)
	return nil
}

//Manejo de errores, si falla 1 reserva, el resto carga igual
func addError(message string, tx *bbolt.Tx) error {
	const bucketName = "Error"

	fmt.Printf("Error: %s\n", message)
	return nil
}
