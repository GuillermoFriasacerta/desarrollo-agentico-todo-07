---
name: feature-plan-verification
description: Perfil del plan de verificación de una feature: dimensionado por criterio y por riesgo, sin comprobaciones de más. Lo aplica el servicio plan_verification de adw.
allowed-tools: Read, Glob, Grep, Bash, Write
---

# Plan de verificación de una feature

Recibes en `$ARGUMENTS` la especificación de una feature y la ruta donde escribir el plan. En todo lo común sigue el servicio `plan_verification` (`.claude/skills/adw/references/services/plan_verification.md`): partir de la especificación, los campos de cada afirmación, la evidencia por método, la prueba para el revisor, los supuestos y la publicación. Este perfil fija lo que cambia para una feature.

Una feature añade comportamiento. El plan tiene el tamaño que pida el cambio, no uno fijo: decídelo tú, con estos criterios.

## Qué incluye siempre

- **Estructural:** una afirmación para los gates de formato, análisis y compilación (bloqueantes) y la imagen de Docker (informativa), y una para el juicio del diff con las rúbricas de estructura.
- **Funcional:** para cada criterio, su camino principal, medido en el sistema real, y el test nuevo que lo demuestra, en rojo sobre la base y en verde sobre el cambio. Si el criterio depende del navegador, como el JavaScript de `public/assets/`, el test es de navegador: uno de Go no ejecuta ese código y pasaría también en la base. Y la suite completa, tests de navegador incluidos, bloqueante.
- **Fidelidad:** la cobertura (cada criterio con sus afirmaciones y un test asociado si tiene comportamiento observable), la contención (cada cambio del diff trazado a un criterio) y la documentación al día (cada documento de `docs/knowledge/` y `docs/procedures/` que describe lo que cambia, actualizado), las tres bloqueantes.
- **Un juicio de textos** solo si el cambio muestra textos nuevos.

## Cuándo añadir más

- **Un caso límite o negativo** cuando el criterio lo implica: un número que puede ser cero o uno, un campo que puede estar vacío, una entrada que el servidor puede rechazar.
- **Una variante** (con y sin JavaScript, extremos, recargar, errores de consola) solo cuando el criterio la menciona o cuando el código de la base muestra un segundo camino que afecta al criterio. Si no, no.
- **Una regresión propia** solo si el cambio toca algo con comportamiento visible que ya existía; en lo demás, la suite la cubre.
- **El riesgo decide la profundidad:** persistencia, JavaScript, varias pantallas o datos piden más comprobaciones que un cambio de texto en una pantalla.

## Dimensiona con criterio

Una afirmación existe si su fallo rechazaría el cambio por algo que la issue pide o implica. Si dos comprobaciones fallarían por el mismo motivo y se demuestran con la misma evidencia, son una. Lo que decidas no comprobar, anótalo en «Supuestos» como «No se comprueba: …», para que el humano lo vea y pueda pedirlo.
