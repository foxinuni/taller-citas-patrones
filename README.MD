## Para Desarrolladores
Pues realmente que no se necesita mucha documentacion, todos los scripts estan en el archivo de npm `package.json`. Se utilizan algunas dependencias de Node como Tailwind CLI asi que es necesario que se instalen esas dependencias antes de correr cualquier script.

En adicion a esto es necesario tener los siguientes ejecutables por aparte:

Instalar Air (opcional):
```
go install github.com/air-verse/air@latest
```

Instalar Go migrate:
```
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
```

Instalar Templ:
```
go install github.com/a-h/templ/cmd/templ@latest
```

Instalar Wire:
```
go install github.com/google/wire/cmd/wire@latest
```

El archivo Docker de igual manera contine las instrucciones para crear una maquina que pueda correr el sistema.


Correr con docker (importante seleccionar el folder de salida "data" como un folder real):
```
docker run -d \
  --name taller-1-patrones \
  -v ./data:/app/data \
  -p 8080:8080 \
  $(docker build -q https://github.com/foxinuni/taller-1-patrones.git)
```