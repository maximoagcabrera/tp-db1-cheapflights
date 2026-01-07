package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/fatih/color"
	_ "github.com/lib/pq"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	//"strings"
	"time"
	"go.etcd.io/bbolt"
)

//Interfaz del usuario
func viewMenu() {
	azul := color.New(color.FgHiBlue).SprintFunc()  // Azul
	amarillo := color.New(color.FgHiYellow).SprintFunc()  // Amarillo
	fmt.Println(
		azul(`--------MENU---------`) + "\n" +
		amarillo(`1- Crear base de datos
2- Borrar base de datos
3- Conectar a la base de datos
4- Crear tablas
5- Crear PKs y FKs
6- Elimina PKs y FKs
7- Cargar datos
8- Crear SP y Triggers
9- Iniciar Pruebas
10- Cargar datos en BoltBD
0- Cerrar menu `) + "\n" +
		azul(`---------------------`) + "\n",  // Añadimos un salto de línea al final
	)
}

//Creamos la base de datos
func createDatabase() {
	db, err := sql.Open("postgres", "user=postgres host=localhost dbname=postgres sslmode=disable")
	if err != nil {
		log.Fatal("Error conectando a postgres:", err)
	}
	defer db.Close()

	_, err = db.Exec(`create database cabral_cabrera_fauda_goni_db1;`)
	if err != nil {
		log.Fatal("Error creando base de datos:", err)
	}
	fmt.Println("Base de datos creada con éxito.")
}

//Establecemos conexion con la base de datos ya creada
func conectionDB() *sql.DB {
	newDB, err := sql.Open("postgres", "user=postgres host=localhost dbname=cabral_cabrera_fauda_goni_db1 sslmode=disable")
	if err != nil {
		log.Fatal("Error conectando a la base de datos",err)
	}
	fmt.Println("Conexión establecida correctamente.")
	return newDB
}

//Limpieza de la base de datos
func eraseDatabase() {
	db, err := sql.Open("postgres", "user=postgres host=localhost dbname=postgres sslmode=disable")
	if err != nil {
		log.Fatal("Error conectando a postgres:", err)
	}
	defer db.Close()

	_, err = db.Exec(`drop database if exists cabral_cabrera_fauda_goni_db1;`)
	if err != nil {
		log.Fatal("Error eliminando base de datos:", err)
	}
	fmt.Println("Base de datos eliminada con éxito.")
}

//Lectura de archivos .sql externos
func readArchive(ubication string) string {
		result, err := ioutil.ReadFile(ubication)
		if err != nil {
			log.Fatal(err)
		}
		return string(result)
}

func createTables(db *sql.DB) {
	_, err := db.Exec(readArchive("sql/tablas/remove_tablas.sql")) //Si existen las tablas, las borra para no generar error.
	if err!= nil{
		log.Fatal("Error borrando las tablas",err)
	}

	_, err = db.Exec(readArchive("sql/tablas/clientes.sql")) //Crea la tabla "cliente".
	if err != nil {
		log.Fatal("Error creando la tabla 'clientes'",err)
	}
	fmt.Println("Tabla 'cliente' creada exitosamente.")

	_, err = db.Exec(readArchive("sql/tablas/aeropuerto.sql")) //Crea la tabla "aeropuerto"
	if err != nil {
		log.Fatal("Error creando la tabla 'aeropuerto'",err)
	}
	fmt.Println("Tabla 'aeropuerto' creada exitosamente.")

	_, err = db.Exec(readArchive("sql/tablas/ruta.sql")) //Crea la tabla "ruta"
	if err != nil {
		log.Fatal("Error creando la tabla 'ruta'",err)
	}
	fmt.Println("Tabla 'ruta' creada exitosamente.")

	_, err = db.Exec(readArchive("sql/tablas/vuelo.sql")) //Crea la tabla "vuelo"
	if err != nil {
		log.Fatal("Error creando la tabla 'vuelo'",err)
	}
	fmt.Println("Tabla 'vuelo' creada exitosamente.")

	_, err = db.Exec(readArchive("sql/tablas/reserva_pasaje.sql")) //Crea la tabla "reserva_pasaje"
	if err != nil {
		log.Fatal("Error creando la tabla 'reserva_pasaje'",err)
	}
	fmt.Println("Tabla 'reserva_pasaje' creada exitosamente.")

	_, err = db.Exec(readArchive("sql/tablas/error.sql")) //Crea la tabla "error"
	if err != nil {
		log.Fatal("Error creando la tabla 'error'",err)
	}
	fmt.Println("Tabla 'error' creada exitosamente.")

	_, err = db.Exec(readArchive("sql/tablas/envio_email.sql")) //Crea la tabla "envio_email"
	if err != nil {
		log.Fatal("Error creando la tabla 'envio_email'",err)
	}
	fmt.Println("Tabla 'envio_email' creada exitosamente.")

	_, err = db.Exec(readArchive("sql/tablas/datos_de_prueba.sql")) //Crea la tabla "datos_de_prueba"
	if err != nil {
		log.Fatal("Error creando la tabla 'datos_de_prueba'",err)
	}
	fmt.Println("Tabla 'datos_de_prueba' creada exitosamente.")
}


func createPkFk(db *sql.DB) {

	_, err := db.Exec(readArchive("sql/pkfk/primary_key.sql")) //Crea las pks
	if err != nil {
		log.Fatal("Error creando las Primary Keys",err)
	}
	fmt.Println("Primary Keys creadas exitosamente.")

	_, err = db.Exec(readArchive("sql/pkfk/foreign_key.sql"))//Crea las fks
	if err != nil {
		log.Fatal("Error creando las Foreign Keys",err)
	}
	fmt.Println("Foreign Keys creadas exitosamente.")
}

func removeClaves(db *sql.DB) {
	_, err := db.Exec(readArchive("sql/pkfk/remove_fk.sql")) //Borra las fks
	if err != nil {
		log.Fatal("Error eliminando las FKs.", err)
	}

	_, err = db.Exec(readArchive("sql/pkfk/remove_pk.sql")) //Borra las pks.
	if err != nil {
		log.Fatal("Error eliminando las PKs.", err)
	}
	fmt.Println("Todas las FKS y PKS se borraron exitosamente.")
}


//Carga de datos Json, convierte el .json en estructuras de go
func loadDataJson(archive string, destine interface{}) error {
	file, err := os.Open(archive)
	if err != nil {
		log.Fatal("Error abriendo el archivo %s: %v", archive, err)
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	if err := decoder.Decode(destine); err != nil {
		log.Fatal("Error decodificando el JSON %v", err)
	}
	return nil
}

func loadData(route string, list interface{}) {
	err := loadDataJson(route, list)
	if err != nil {
		log.Fatal("Error al cargar datos del archivo %s: %v", route, err)
	}
}

//Insertamos los datos en sql
func insertAeropuertos(list []Aeropuerto, db *sql.DB) {
	for _, a := range list {
		_, err := db.Exec("insert into aeropuerto (id_aeropuerto, nombre, localidad, provincia) values ($1, $2, $3, $4)", a.IDAeropuerto, a.Nombre, a.Localidad, a.Provincia)
		if err != nil {
			log.Fatal("Error insertando datos a 'aeropuerto'",err)
		}
	}
	fmt.Println("Aeropuertos ingresados correctamente.")
}

func insertClientes(list []Cliente, db *sql.DB) {
	for _, c := range list {
		fechaNacimiento, err := time.Parse("2006-01-02", c.FechaNacimiento)
		if err != nil {
			log.Fatal("Error al parsear FechaNacimiento:", err)
		}
		_, err = db.Exec("insert into cliente (id_cliente, nombre, apellido, dni, fecha_nacimiento, telefono, email) VALUES ($1, $2, $3, $4, $5, $6, $7)",
			c.IDCliente, c.Nombre, c.Apellido, c.DNI, fechaNacimiento, c.Telefono, c.Email)
		if err != nil {
			log.Fatal("Error insertando datos a 'cliente'", err)
		}
	}
	fmt.Println("Clientes ingresados correctamente.")
}

func insertRutas(list []Ruta, db *sql.DB){
	for _, r := range list {
		_, err := db.Exec("insert into ruta (nro_ruta, id_aeropuerto_origen, id_aeropuerto_destino, duracion) values ($1, $2, $3, $4)",
			r.NroRuta, r.IDAeropuertoOrigen, r.IDAeropuertoDestino, r.Duracion)
		if err != nil {
			log.Fatal("Error insertando datos a 'ruta'",err)
		}
	}
	fmt.Println("Rutas ingresadas correctamente.")
}

func insertPruebas(list []DatoPrueba, db *sql.DB) {
	for _, p := range list {
		if p.FSalidaVuelo != "" { //Veo si la fecha existe.
			Fsalidavuelo, err := time.Parse("2006-01-02 15:04", p.FSalidaVuelo)
			if err != nil {
				log.Fatal("Error al parsear Fecha de salida del vuelo:", err)
			}

			_, err = db.Exec("insert into datos_de_prueba (id_orden, operacion, nro_ruta, f_salida_vuelo, nro_asientos_totales, id_vuelo, id_cliente, id_reserva, nro_asiento) values ($1, $2, $3, $4, $5, $6, $7, $8, $9)",
				p.IDOrden, p.Operacion, p.NroRuta, Fsalidavuelo, p.NroAsientosTotales, p.IDVuelo, p.IDCliente, p.IDReserva, p.NroAsiento)
			if err != nil {
				log.Fatal("Error insertando datos a 'datos_de_prueba",err)
			}
		} else { //Si no existe, que no la ponga y listo.
			_, err := db.Exec("insert into datos_de_prueba (id_orden, operacion, nro_ruta, nro_asientos_totales, id_vuelo, id_cliente, id_reserva, nro_asiento) values ($1, $2, $3, $4, $5, $6, $7, $8)",
				p.IDOrden, p.Operacion, p.NroRuta, p.NroAsientosTotales, p.IDVuelo, p.IDCliente, p.IDReserva, p.NroAsiento)
			if err != nil {
				log.Fatal("Error insertando datos a 'datos_de_prueba",err)
			}
		}
	}
	fmt.Println("Datos de prueba ingresados correctamente.")
}

//Logica de negocio, lee archivos .sql que contienen funciones y triggers y compila
func crearStoredProcedures(db *sql.DB) error {
	_, err := db.Exec(readArchive("sql/stored_procedures/remove_sp_trigger.sql")) //Borra los stored procedures y trigger que existan.
	if err != nil {
		log.Fatal("Error eliminando los Stored Procedures y Triggers.", err)
	}
	fmt.Println("Stored procedures y Triggers eliminados correctamente.")

	_, err = db.Exec(readArchive("sql/stored_procedures/apertura_vuelo.sql")) //
	if err != nil {
		log.Fatal("Error ejecutando stored procedure apertura_vuelo", err)
	}
	fmt.Println("apertura_vuelo iniciado correctamente.")

	_, err = db.Exec(readArchive("sql/stored_procedures/reserva_pasaje.sql")) //
	if err != nil {
		log.Fatal("Error ejecutando stored procedure reserva_pasaje", err)
	}
	fmt.Println("reserva_pasaje iniciado correctamente.")

	_, err = db.Exec(readArchive("sql/stored_procedures/check_in_asiento.sql")) //
	if err != nil {
		log.Fatal("Error ejecutando stored procedure check_in_asiento", err) //
	}
	fmt.Println("check_in_asiento iniciado correctamente.")

	_, err = db.Exec(readArchive("sql/stored_procedures/anular_reserva.sql")) //
	if err != nil {
		log.Fatal("Error ejecutando stored procedure anular_reserva", err) //
	}
	fmt.Println("anular_reserva iniciado correctamente.")

	_, err = db.Exec(readArchive("sql/stored_procedures/envio_email_reserva.sql")) //
	if err != nil {
		log.Fatal("Error ejecutando trigger envio_email_reserva.sql", err) //
	}
	fmt.Println("envio_email_reserva iniciado correctamente.")
	return nil
}

//Funcion que testea
func iniciar_pruebas(db *sql.DB) error{
	//Llamos a readArchive, que solo devuelve el contenido del archivo como un string
	sqlScript := readArchive("sql/stored_procedures/iniciar_pruebas.sql")

	//Luego puedes realizar la ejecucion del script sin necesidad de capturar un error aqui
    _, err := db.Exec(sqlScript)
	if err != nil {
		log.Fatal("Error creando la funcion en la base datos: %v", err)
	}

	_,err = db.Exec("SELECT iniciar_pruebas();")
	if err != nil{
		log.Fatal("Error ejecutando la funcion iniciar_pruebas: %v", err)
	}
	fmt.Println("Pruebas ejecutadas correctamente.")
	return nil
}

//Funcion para depurar BoltDB
func dumpBoltDB(db *bbolt.DB) {
    db.View(func(tx *bbolt.Tx) error {
        fmt.Println("=== Contenido de BoltDB ===")

        return tx.ForEach(func(bucketName []byte, b *bbolt.Bucket) error {
            fmt.Printf("Bucket: %s\n", bucketName)

            // Recorrer claves/valores
            b.ForEach(func(k, v []byte) error {
                fmt.Printf("  %s -> %s\n", k, v)
                return nil
            })

            fmt.Println()
            return nil
        })
    })
}


func readOption() int {
	var option int
	_, err := fmt.Scan(&option)
	if err != nil {
		log.Fatal(err)
		return -1
	}
	return option
}

func clearScreen() {
	var cmd *exec.Cmd
	cmd = exec.Command("clear")
	cmd.Stdout = os.Stdout
	cmd.Run()
}
func closeMenu(){
	fmt.Println("Menu cerrado correctamente :)")
}

func main() {
	clearScreen()
	var db *sql.DB
	var aeropuertos []Aeropuerto
	var clientes []Cliente
	var rutas []Ruta
	var pruebas []DatoPrueba
	var boltDB *bbolt.DB
	for {
		viewMenu()
		option := readOption()
		switch option {
		case 0:
			clearScreen()
			closeMenu()
			return
		case 1:
			clearScreen()
			createDatabase()
		case 2:
			clearScreen()
			eraseDatabase()
		case 3:
			clearScreen()
			db = conectionDB()
		case 4:
			clearScreen()
			createTables(db)
		case 5:
			clearScreen()
			createPkFk(db)
		case 6:
			clearScreen()
			removeClaves(db)
		case 7:
			clearScreen()
			loadData("data/aeropuertos.json", &aeropuertos)
			loadData("data/clientes.json", &clientes)
			loadData("data/rutas.json", &rutas)
			loadData("data/datos_de_prueba.json", &pruebas)
			insertAeropuertos(aeropuertos, db)
			insertClientes(clientes,db)
			insertRutas(rutas,db)
			insertPruebas(pruebas,db)
		case 8:
			clearScreen()
			crearStoredProcedures(db)
		case 9:
            clearScreen()
            iniciar_pruebas(db)
        case 10:
			clearScreen()
			boltDB = Open_boltDB()
			dumpBoltDB(boltDB)
			defer boltDB.Close()
		}
	}
}
