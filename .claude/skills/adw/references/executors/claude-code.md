# Ejecutor: Claude Code

Cómo se lanza en Claude Code cada rol del [catálogo](../roles.md), con qué modelo y qué impone el arnés por construcción.

## Niveles

| Nivel | Modelo |
| --- | --- |
| Capaz | `opus` |
| Intermedio | `sonnet` |
| Ligero | `haiku` |

## Agentes

Cada rol, salvo `coordinator`, es un agente del proyecto en `.claude/agents/<rol>.md`. Su definición fija el modelo de su nivel, sus herramientas y sus skills, y su prompt dice lo estable del rol: qué es, qué no recibe y qué no toca.

| Rol | Herramientas | Sin | Precarga | Usa con `Skill` |
| --- | --- | --- | --- | --- |
| `planner` | Read, Glob, Grep, Bash, Write, Skill | Edit | — | `feature`, `bug`, `chore` |
| `verification-planner` | Read, Glob, Grep, Bash, Write | Edit | — | — |
| `implementer` | Read, Glob, Grep, Bash, Write, Edit, Skill | — | `implement` | — |
| `checker` | Read, Glob, Grep, Bash, Write, Skill | Edit | `check` | — |
| `evidence-collector` | Read, Glob, Grep, Bash, Write, Skill | Edit | `verify`, `playwright-cli` | — |
| `judge` | Read, Glob, Grep, Bash, Write, Skill | Edit | `judge` | — |
| `documenter` | Read, Glob, Grep, Bash, Write, Edit, Skill | — | `document` | — |

- **Precarga** pone la skill entera en el contexto del agente, que la relee en cada turno: es para la que el rol siempre aplica.
- **Usa con `Skill`** es para las que dependen del caso, como la del tipo de la petición. El rol solo carga la que toca, y el encargo dice cuál.

Ninguno tiene la herramienta `Agent`: un rol no lanza otros. Si necesitas apoyo para explorar o diagnosticar, lánzalo tú.

## Cómo se lanza

- Con la herramienta `Agent` y `subagent_type: <rol>`. Sin el parámetro `model`: lo pone el agente.
- Pasa `model` solo para subir el nivel de un rol en un caso concreto, y apunta el motivo en el panel. El de la llamada gana al de la definición.
- El encargo lleva solo lo del caso: la issue, la rama, las afirmaciones, el directorio de evidencias, la base. Lo del rol ya está en el agente. Recuérdale que nadie le va a responder y que devuelva sus dudas y bloqueos en su informe.
- Las skills precargadas llegan sin argumentos: donde la skill dice «recibe en `$ARGUMENTS`», el encargo hace de argumento. Si el rol puede usar varias con `Skill`, como `planner`, dile cuál aplicar según el tipo de la petición.
- Puedes lanzar en paralelo lo independiente y esperar a todos, como dice la skill.

## Qué impone Claude Code y qué no

- **Impone:**
  - el modelo de cada rol;
  - las skills que trae cargadas;
  - qué herramientas tiene;
  - por encima de todo, los permisos de `.claude/settings.json`, que valen también para los agentes: lo que el proyecto deniega sigue denegado.
- **No impone** que un rol no cambie el código. Un agente sin `Edit` puede reescribir un fichero con `Write`, y cualquiera con `Bash` puede escribir con la shell. `tools` admite nombres de herramientas, no patrones de rutas ni de comandos. Por eso quitar `Edit` es una señal de lo que se espera del rol, no una garantía. La garantía es la comprobación de la huella del árbol que hace el coordinador tras cada rol que no cambia el código.

## Coordinador

- Es la sesión de `claude -p` que lanza el trigger (`bin/adw`), con `--model` si `config/adw.yml` fija `model`. Si no, el modelo por defecto de Claude Code.
- En una petición directa, es la sesión del humano, con su modelo.

## Logs

Los hooks `SubagentStart` y `SubagentStop` reciben el nombre del agente en `agent_type`, y `log.sh` los guarda en el log del caso: ahí se ve qué rol hizo cada trabajo.
