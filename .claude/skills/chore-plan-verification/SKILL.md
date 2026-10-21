---
name: chore-plan-verification
description: Perfil del plan de verificación de un chore, una tarea de mantenimiento que no cambia lo que hace la aplicación. Lo aplica el servicio plan_verification de adw.
allowed-tools: Read, Glob, Grep, Bash, Write
---

# Plan de verificación de un chore

Recibes en `$ARGUMENTS` la especificación de un chore y la ruta donde escribir el plan. En todo lo común sigue el servicio `plan_verification` (`.claude/skills/adw/references/services/plan_verification.md`): partir de la especificación, los campos de cada afirmación, la evidencia por método, la prueba para el revisor, los supuestos y la publicación. Este perfil fija lo que cambia para un chore.

Un chore no cambia lo que hace la aplicación. Verificarlo es demostrar dos cosas: que no ha cambiado nada observable y que el cambio es el que la issue dice. Casi todo se comprueba con pruebas automáticas.

## Qué incluye siempre

- **Estructural:** una afirmación para los gates de formato, análisis y compilación (bloqueantes) y la imagen de Docker (informativa), y una para el juicio del diff con las rúbricas de estructura.
- **Funcional:** una afirmación para la suite completa, tests de navegador incluidos, bloqueante. Es lo que detectaría que el mantenimiento ha roto algo.
- **Cada criterio de la issue,** con una afirmación automatizada que lo demuestre: un comando, una búsqueda en el repositorio o un gate. Por ejemplo, «la versión de PostgreSQL de `compose.yaml` y del CI es la 17.7». Un juicio solo si no hay forma de automatizarlo.
- **Fidelidad:** la contención (cada cambio del diff responde a lo que dice la issue: subir una dependencia no trae de paso un refactor) y la documentación al día (cada documento de `docs/knowledge/` y `docs/procedures/` que cuenta lo que cambia, actualizado), las dos bloqueantes.

## Qué no entra

- Capturas, vídeos ni medidas del DOM, salvo que la issue diga que algo visible debe quedar igual: entonces, una afirmación con la medida antes y después.
- Un test nuevo en rojo sobre la base y en verde sobre el cambio: no hay comportamiento nuevo que demostrar.
- Afirmaciones propias de regresión, variantes de navegador y juicio de textos: la suite ya cubre lo existente.

## Si no es un chore

Si la issue pide algo que se ve o se comporta distinto, no es un chore. No lo planifiques como tal: devuélvelo como duda funcional.

## Dimensiona con criterio

Una afirmación existe si su fallo rechazaría el cambio por algo que la issue pide. Si dos comprobaciones fallarían por el mismo motivo y se demuestran con la misma evidencia, son una. Lo que decidas no comprobar, anótalo en «Supuestos» como «No se comprueba: …».
