# Roles

Los servicios dicen qué trabajo se hace, no quién lo hace. Los roles son ese quién: cada trabajo que el coordinador delega lo hace un rol de este catálogo. El rol fija lo que el trabajo necesita y lo que no debe tener, con independencia del arnés que lo ejecute.

- **Nivel de modelo.** Qué capacidad necesita el trabajo:
  - **capaz**, para decidir qué se construye y cómo se comprueba: planificar;
  - **intermedio**, para ejecutar un plan o una rúbrica ya fijados: escribir código o documentos, juzgar con una rúbrica y coordinar;
  - **ligero**, para lo mecánico: ejecutar comandos y extraer material.
- **Puede.** Lo que el trabajo necesita leer, escribir o ejecutar.
- **No recibe.** El contexto que no se le pasa, para que su resultado sea independiente: un juez que conoce la justificación del implementador ya no juzga solo el cambio.

Cómo se lanza cada rol, con qué modelo concreto y qué impone el arnés por construcción lo dice el ejecutor de cada arnés (ver la skill).

## Catálogo

| Rol | Trabajo | Nivel | Puede | No recibe |
| --- | --- | --- | --- | --- |
| `coordinator` | Interpreta el evento, reconstruye el estado, compone el workflow, delega y proyecta. Hace él mismo `clarify`, `setup_environment`, `teardown_environment`, `publish`, `answer_review` y `report_failure`, y decide si una revisión pide cambios | Intermedio | Todo lo que permita el proyecto | — |
| `planner` | `write_plan`, con las skills `feature`, `bug` y `chore` | Capaz | Leer el código y la documentación; escribir el plan | — |
| `verification-planner` | `plan_verification` | Capaz | Leer la especificación y el código de la base; escribir el plan de verificación | El plan de implementación |
| `implementer` | `implement`, con la skill `implement` | Intermedio | Cambiar el código, los tests y la documentación del sistema; ejecutar los tests | — |
| `checker` | En `validate`, los gates (skill `check`) y la extracción del material que se juzga: el diff, la lista de ficheros, las claves nuevas | Ligero | Ejecutar los gates y leer el diff; escribir evidencias | El informe del implementador |
| `evidence-collector` | En `validate`, las afirmaciones automatizadas con su evidencia (skill `verify`): tests en rojo y en verde, HTTP, valores del DOM, capturas, vídeos y base de datos | Intermedio | Ejecutar los tests, el navegador y las consultas; escribir evidencias | El informe del implementador |
| `judge` | En `validate`, el juicio con rúbrica, o el intento de refutar un hallazgo (skill `judge`) | Intermedio | Leer el diff y el material; escribir su acta | El informe del implementador y, al refutar, el razonamiento del juicio |
| `documenter` | `document`, con la skill `document` | Intermedio | Escribir la historia del caso: su documento de cierre y sus decisiones | — |

## Reglas

- **Cada trabajo delegado va a su rol.** El coordinador puede usar otro rol o subir el nivel en un caso concreto, por ejemplo un `implementer` capaz tras un fallo que no entiende, pero lo apunta en el panel del caso con su motivo.
- **Cada instancia empieza sin contexto.** El rol no se reutiliza entre trabajos: cada verificación y cada refutación es una instancia nueva.
- **Un verificador por dimensión** es una combinación de roles: en la estructural, `checker` y `judge`; en la funcional, `evidence-collector` y `judge`; en la fidelidad, `checker` y `judge`.
- **Quién no cambia el código.** `planner`, `verification-planner`, `checker`, `evidence-collector` y `judge` no cambian el código: lo suyo va a ficheros ignorados por git, salvo la copia de trabajo del plan, que se borra al terminar. El coordinador lo comprueba en lo observable después de cada uno, porque no todos los arneses pueden impedirlo.
