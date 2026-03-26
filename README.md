# Sistemas Distribuidos - TP0: Docker + Comunicaciones + Concurrencia
En el presente trabajo se resuelve una guía de ocho ejercicios incrementales que abarcan el uso de *Docker* como herramienta de virtualización, comunicación de procesos y concurrencia.  

El repositorio consta de un cliente y un servidor escritos en *Golang* y *Python*, respectivamente. En cada una de las ramas (`ej<núm>`) se encuentra el desarrollo de cada uno de los ejercicios propuestos.

## Instrucciones de uso
El repositorio cuenta con un **Makefile** que incluye distintos comandos en forma de targets. Los targets se ejecutan mediante la invocación de:  **make \<target\>**. Los target imprescindibles para iniciar y detener el sistema son **docker-compose-up** y **docker-compose-down**, siendo los restantes targets de utilidad para el proceso de depuración.

Los targets disponibles son:

| target  | accion  |
|---|---|
|  `docker-compose-up`  | Inicializa el ambiente de desarrollo. Construye las imágenes del cliente y el servidor, inicializa los recursos a utilizar (volúmenes, redes, etc) e inicia los propios containers. |
| `docker-compose-down`  | Ejecuta `docker-compose stop` para detener los containers asociados al compose y luego  `docker-compose down` para destruir todos los recursos asociados al proyecto que fueron inicializados. Se recomienda ejecutar este comando al finalizar cada ejecución para evitar que el disco de la máquina host se llene de versiones de desarrollo y recursos sin liberar. |
|  `docker-compose-logs` | Permite ver los logs actuales del proyecto. Acompañar con `grep` para lograr ver mensajes de una aplicación específica dentro del compose. |
| `docker-image`  | Construye las imágenes a ser utilizadas tanto en el servidor como en el cliente. Este target es utilizado por **docker-compose-up**, por lo cual se lo puede utilizar para probar nuevos cambios en las imágenes antes de arrancar el proyecto. |
| `build` | Compila la aplicación cliente para ejecución en el _host_ en lugar de en Docker. De este modo la compilación es mucho más veloz, pero requiere contar con todo el entorno de Golang y Python instalados en la máquina _host_. |


## Comentarios sobre el desarrollo de los ejercicios
A continuación se detallan los aspectos más importantes de la solución provista para cada uno de los ejercicios.

### Ejercicio N°1:
**Definir un script de bash `generar-compose.sh` que permita crear una definición de Docker Compose con una cantidad configurable de clientes.  El nombre de los containers deberá seguir el formato propuesto: client1, client2, client3, etc.**  

Se implementó un script en Python (que se llama desde otro bash script, puesto que las pruebas ignoran el shebang) que define los containers para el server y para tantos clientes como se pidan; configurando sus variables de entorno y la Docker network.  
El script puede ejecutarse con
```sh
./generar-compose.sh <nombre del archivo de salida> <cantidad de clientes>`
```
o con
```sh
./generar-compose <nombre del archivo de salida> <cantidad de clientes>` # python
```


### Ejercicio N°2:
**Modificar el cliente y el servidor para lograr que realizar cambios en el archivo de configuración no requiera reconstruír las imágenes de Docker para que los mismos sean efectivos. La configuración a través del archivo correspondiente (`config.ini` y `config.yaml`, dependiendo de la aplicación) debe ser inyectada en el container y persistida por fuera de la imagen (hint: `docker volumes`).**  

Se agregaron los archivos de configuración como volúmenes de los containers de los clientes y del servidor.


### Ejercicio N°3:
**Crear un script de bash `validar-echo-server.sh` que permita verificar el correcto funcionamiento del servidor utilizando el comando `netcat` para interactuar con el mismo. Dado que el servidor es un echo server, se debe enviar un mensaje al servidor y esperar recibir el mismo mensaje enviado.**  

Se implementó un script de Python (al que también lo llama un bash scrip) que corre un container con la imagen `subfuzion/netcat` ([docker hub](https://hub.docker.com/r/subfuzion/netcat)) y que usa `netcat` para mandar un mensaje al (en esta instancia del trabajo) *echo* server, y checkear posteriormente que el mensaje que éste rebota coincide con el que se mandó. La dirección del servidor se extrae de su archivo de configuración `server/config.ini`.  
El script puede ejecutarse con `./validar-echo-server.sh` o `./validar-echo-server`.  
El output del script es 
```
action: test_echo_server | result: success
```
o
```
action: test_echo_server | result: fail`
```
en caso de error.

### Ejercicio N°4:
**Modificar servidor y cliente para que ambos sistemas terminen de forma _graceful_ al recibir la signal SIGTERM. Terminar la aplicación de forma _graceful_ implica que todos los _file descriptors_ (entre los que se encuentran archivos, sockets, threads y procesos) deben cerrarse correctamente antes que el thread de la aplicación principal muera. Loguear mensajes en el cierre de cada recurso (hint: Verificar que hace el flag `-t` utilizado en el comando `docker compose down`).**

Se agregó un signal handler para `SIGTERM` y `SIGINT` que cierra el welcoming socket del server y detiene su loop de ejecución.

### Ejercicio N°5:
**Modificar la lógica de negocio tanto de los clientes como del servidor para nuestro nuevo caso de uso.**

- **Cliente: Emulará a una _agencia de quiniela_ que participa del proyecto. Existen 5 agencias. Deberán recibir como variables de entorno los campos que representan la apuesta de una persona: nombre, apellido, DNI, nacimiento, numero apostado (en adelante 'número'). Ej.: `NOMBRE=Santiago Lionel`, `APELLIDO=Lorca`, `DOCUMENTO=30904465`, `NACIMIENTO=1999-03-17` y `NUMERO=7574` respectivamente. Los campos deben enviarse al servidor para dejar registro de la apuesta. Al recibir la confirmación del servidor se debe imprimir por log: `action: apuesta_enviada | result: success | dni: ${DNI} | numero: ${NUMERO}`.**
- **Servidor: Emulará a la _central de Lotería Nacional_. Deberá recibir los campos de la cada apuesta desde los clientes y almacenar la información mediante la función `store_bet(...)` para control futuro de ganadores. La función `store_bet(...)` es provista por la cátedra y no podrá ser modificada por el alumno. Al persistir se debe imprimir por log: `action: apuesta_almacenada | result: success | dni: ${DNI} | numero: ${NUMERO}`.**

Se implementó un módulo de comunicación (`net/`, y `comms` en el cliente; `net` ya existe en Golang) entre el cliente y el servidor, donde se maneja el envío y la recepción de los mensajes que éstos se envían: `Bet`, `ACK`/`NACK`. La aceptación de conexiones entrantes en el servidor se modularizó a la clase `Rendezvous`. El correcto manejo de los sockets se modularizó en la clase `Conn`(`ection`) tanto en el cliente como en el servidor. Y se definió un protocolo de comunicación, con su correspondiente serialización y deserialización de mensajes:
```py
# en el servidor
BYTE_ORDER = "big"

LEN_TYPE = 1
LEN_STR_SIZE = 1  # preappended string length comes in one byte
LEN_BET_NUMBER = 8

TYPE_BET = b"\x00"
TYPE_RESPONSE_ACK = b"\x01"
TYPE_RESPONSE_NACK = b"\x02"
```
Los atributos de los mensajes `Bet` se tratan como `string`s y se envían seguidos de su largo, luego del byte con el tipo de mensaje. Los mensajes de ACK se mandan como un sólo byte.  

A medida que el servidor recibe apuestas, las persiste e imprime
```
action: apuesta_almacenada | result: success | dni: <documento del apostante> | numero: <número de la apuesta>
```

### Ejercicio N°6:
**Modificar los clientes para que envíen varias apuestas a la vez (modalidad conocida como procesamiento por _chunks_ o _batchs_).**

Se inyectaron los archivos de apuestas de cada cliente usando con docker volumes para sus containers. Para cargar los chunks en la memoria del cliente, se implementó la clase `Storage`, que recibe como parámetro la cantidad máxima de apuestas a leer en cada iteración del main loop del cliente. A su vez, `Storage` mantiene la memoria consumida por las apuestas cappeada a 8 kB.  
A medida que el cliente levanta batches de apuestas, los envía al servidor, que imprime `action: apuesta_recibida | result: success | cantidad: <cantidad de apuestas>` si todas fueron procesadas correctamente, o `action: apuesta_recibida | result: fail | cantidad: <cantidad de apuestas>` en caso contrario.  
La cantidad máxima de apuestas de cada batch puede ser establecida en el archivo de configuración del cliente `client/config.yaml`.


### Ejercicios N°7 y N°8
La resolución del ejercicio siete bastó para resolver el ejercicio ocho.

#### Ejercicio N°7:
**Modificar los clientes para que notifiquen al servidor al finalizar con el envío de todas las apuestas y así proceder con el sorteo.
Inmediatamente después de la notificacion, los clientes consultarán la lista de ganadores del sorteo correspondientes a su agencia.
Una vez el cliente obtenga los resultados, deberá imprimir por log: `action: consulta_ganadores | result: success | cant_ganadores: ${CANT}`.**

#### (Parte 3: Repaso de Concurrencia) Ejercicio N°8:
**Modificar el servidor para que permita aceptar conexiones y procesar mensajes en paralelo. En caso de que el alumno implemente el servidor en Python utilizando _multithreading_.**

Se extendió el para que los clientes envíen un mensaje `Fin` al llegar al `EOF` de sus archivos de apuestas. Por su parte, el servidor ahora pasa a levantar threads que manejan el procesamiento de los mensajes entrantes de los clientes; los handles de los threads se guardan para llamar a `join()` cuando termina el server.  

El servidor usa una *condition variable* `pending` para evaluar si todos los clientes que espera (obtenidos del parámetro `<cantidad de clientes>` generación del script de generación del docker compose) terminaron.  
Cuando uno de los threads `client_handler` recibe un mensaje de `Fin`, decrementa la variable `pending` en uno.  
Cuando uno de los threads `client_handler` recibe un mensaje de `Query`, se bloquea esperando a que la condition variable `pending == 0` se cumpla.  
Cuando `pending` llega a cero, el thread que realizó el último decremento, levanta el archivo de apuestas y arma un diccionario con los ids de los clientes y los documentos de sus correspondientes ganadores (en caso de no haber ganadores, el mensaje se manda igual pero vacío) y hace un `notify_all()` sobre la condvar, de manera tal de que los handlers de `Query` bloqueados se despierten y envíen los resultados a sus clientes.  

Por último, el servidor usa un *mutex* para manejar el acceso de escritura concurrente de los handlers de mensajes `Bets`.

## Testing
La cátedra proveyó [pruebas automáticas](https://github.com/7574-sistemas-distribuidos/tp0-tests) de caja negra.
