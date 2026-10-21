---
name: bug-plan-verification
description: Perfil del plan de verificación de un bug: reproducirlo sobre la base y demostrar que el cambio lo arregla, acotado al bug. Lo aplica el servicio plan_verification de adw.
allowed-tools: Read, Glob, Grep, Bash, Write
---

# Plan de verificación de un bug

Recibes en `$ARGUMENTS` la especificación de un bug y la ruta donde escribir el plan. En todo lo común sigue el servicio `plan_verification` (`.claude/skills/adw/references/services/plan_verification.md`): partir de la especificación, los campos de cada afirmación, la evidencia por método, la prueba para el revisor, los supuestos y la publicación. Este perfil fija lo que cambia para un bug.

Un bug es algo que no funciona como debería. La prueba es el bug reproducido y después arreglado, con los mismos pasos, y nada más. Leyendo el código de la base, concreta los pasos que reproducen el fallo y su causa.

## Qué incluye siempre

- **Funcional:**
  - **El bug se reproduce sobre la base.** Los pasos de la issue, con el valor medido que muestra el fallo y una captura, o un vídeo corto si es un recorrido. Si no se reproduce, es un hallazgo: dilo, no lo supongas.
  - **Un test de regresión nuevo,** en rojo sobre la base y en verde sobre el cambio, obligatorio. De navegador si el bug depende del navegador, como el JavaScript de `public/assets/`: uno de Go no ejecuta ese código y pasaría también en la base.
  - **El bug está arreglado.** Los mismos pasos sobre el cambio, con el valor medido esperado y la misma captura o el mismo vídeo. Las dos evidencias, antes y después, son la prueba para el revisor.
  - **La suite completa,** tests de navegador incluidos, bloqueante. Cubre las regresiones.
- **Estructural:** una afirmación para los gates de formato, análisis y compilación (bloqueantes) y la imagen de Docker (informativa), y una para el juicio del diff con las rúbricas de estructura.
- **Fidelidad:** la cobertura (el comportamiento esperado de la issue, con el test asociado), la contención (el diff arregla este bug y nada más: el riesgo típico es arreglar «de paso» otra cosa) y la documentación al día, las tres bloqueantes.

## Qué no entra

- Variantes que el bug no menciona (con y sin JavaScript, extremos, recargar) y comprobaciones del comportamiento cercano que el arreglo no toca.
- Casos límite que no sean la causa del bug. Los que sí lo son, como el valor exacto en el que falla, entran.
- Juicio de textos, salvo que el arreglo cambie un mensaje.
- Más capturas o vídeos que el antes y el después de los mismos pasos.

## Dimensiona con criterio

Una afirmación existe si su fallo rechazaría el cambio por algo que la issue pide. Si dos comprobaciones fallarían por el mismo motivo y se demuestran con la misma evidencia, son una. Lo que decidas no comprobar, anótalo en «Supuestos» como «No se comprueba: …».
