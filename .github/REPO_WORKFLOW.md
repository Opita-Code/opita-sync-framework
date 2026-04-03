# Repo workflow

## Objetivo

Mantener una governance mínima y proporcional para un repo trabajado principalmente por una sola persona, sin perder trazabilidad básica.

## Reglas mínimas

- cambios chicos pueden ir directos a `main`
- cambios medianos pueden usar branch sin issue obligatorio
- cambios grandes, arquitectónicos o que cambian convenciones deberían usar issue + PR

## Pull requests

El template de PR pide:

- issue relacionado
- tipo de cambio
- resumen
- tabla de cambios
- validación realizada

El workflow actual valida solo dos cosas mínimas:

1. que el body del PR enlace un issue (`Closes/Fixes/Resolves #N`)
2. que el título del PR siga formato de Conventional Commits

## Filosofía

No buscamos burocracia enterprise.

Buscamos:

- contexto suficiente
- trazabilidad básica
- menos fricción futura
- mejor continuidad mientras crece el framework

## Boundary mental

- este repo es **OSF** como framework/kernel reusable
- el repo hermano `opita-sync` conserva el source of truth del producto
- la governance del repo no debe mezclar concerns de producto con concerns del framework
