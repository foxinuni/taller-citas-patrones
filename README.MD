## Para Desarrolladores
Pues realmente que no se necesita mucha documentacion, todos los scripts estan en el archivo de npm `package.json`. Se utilizan algunas dependencias de Node como Tailwind CLI asi que es necesario que se instalen esas dependencias antes de correr cualquier script.

En adicion a esto es necesario tener los siguientes ejecutables por aparte:

Instalar Air (opcional):
```
go install github.com/air-verse/air@latest
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