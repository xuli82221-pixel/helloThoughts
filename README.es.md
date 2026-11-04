[English](README.md) | [简体中文](README.zh-CN.md) | [繁體中文](README.zh-TW.md) | [日本語](README.ja.md) | [한국어](README.ko.md) | [Deutsch](README.de.md) | [Français](README.fr.md) | [Español](README.es.md)

# Thoughts

[Instalación](INSTALL.md) · [Documentación](docs/README.md) · [Donar](DONATE.md) · [Code signing policy](docs/CODE_SIGNING_POLICY.md)

Una herramienta de notas gratuita y de código abierto para Windows x64 que prioriza el almacenamiento local. Deja que tu asistente de IA gestione archivos Markdown y léelos o edítalos en una página web local independiente.

**Vista previa sin firmar — 2026-11-04:** `0.3.0-preview.21` unifica los nombres del código fuente, el ejecutable y las configuraciones de clientes como Thoughts. El responsable ha autorizado una [versión preliminar experimental sin firmar](https://github.com/xuli82221-pixel/hello-Thoughts-/releases/tag/v0.3.0-preview.21); los archivos estarán disponibles al terminar la compilación y las comprobaciones. **No se afirma que tenga una firma de editor de confianza ni validación como versión estable.** La aprobación de SignPath Foundation sigue pendiente. Si Windows SAC/SmartScreen bloquea la ejecución, detente y mantén la protección activada. Consulta la [preparación de la publicación](docs/RELEASE_READINESS_20261104.md), la [instalación](INSTALL.md) y las [pruebas nativas de Windows anteriores](docs/WINDOWS_NATIVE_TESTS_20261104.md).

## Tus notas, en tu ordenador

Thoughts no tiene cuentas, base de datos en la nube, servicio de subida de notas, telemetría, IA integrada, temas de pago, suscripciones, integración de pagos ni extensión de navegador. Los archivos Markdown y las imágenes locales son tuyos y se pueden abrir con otros editores.

El ejecutable de Go incorpora la página React. **Las versiones precompiladas no necesitan Node.js, Go, Python ni Docker para ejecutarse.** El cliente de IA que elijas tiene sus propios requisitos de instalación y de cuenta.

## Cómo se podrá instalar

Esta vista previa experimental ofrece **`thoughts-<version>-windows-amd64-setup.exe`**, un ZIP portátil y un MCPB de Claude Desktop en [GitHub Releases](https://github.com/xuli82221-pixel/hello-Thoughts-/releases). El programa se llama **`thoughts.exe`**. Comprueba SHA256 antes de usarlo. Los archivos actuales **no tienen firma Authenticode**; una suma coincidente no autentica al editor. Si Windows bloquea la ejecución, detente. La instalación por usuario inicializa el directorio compartido y crea accesos directos en el escritorio y el menú Inicio sin permisos de administrador. Los datos existentes se conservan; la desinstalación quita lanzadores y accesos directos y mantiene el directorio compartido para otros clientes. Los paquetes estables firmados siguen pendientes.

Los usuarios de Claude Desktop también podrán importar la extensión `.mcpb` de esa versión. Para Codex y Claude Code, el ejecutable del ZIP se registra como servidor MCP local por stdio. La página web puede funcionar sin un cliente de IA. Consulta los comandos de la [guía rápida para usuarios](INSTALL.md#ordinary-windows-users). **Code → Download ZIP** de GitHub descarga código fuente, no la aplicación precompilada.

Los clientes de IA inician el mismo adaptador MCP local con el protocolo estándar por stdio. Este llama al núcleo local autenticado, cuya página web escucha solo en `127.0.0.1`. Cuando pides a un cliente de IA en línea que lea una nota, recibe el contenido y lo trata según sus propias políticas de datos. Priorizar el almacenamiento local no significa que los resultados de las herramientas de IA que solicites nunca salgan del ordenador. Thoughts no cifra los archivos.

## Clientes previstos

| Cliente | Instalación |
| --- | --- |
| Claude Desktop | Formato oficial de extensión `.mcpb` |
| Claude Code en Windows | MCP local por stdio, ZIP portátil para Windows |
| OpenAI Codex en Windows | MCP local por stdio, el mismo ZIP portátil |

Estas son las integraciones previstas; su estado real de validación figura en la revisión de la publicación. Los plugins web habituales de ChatGPT, macOS, Linux, móviles, MCP remoto y la sincronización en la nube quedan fuera de esta fase.

## Funciones

- Crear, leer, editar, buscar, renombrar y mover archivos Markdown locales.
- Árbol de proyectos compacto, notas rápidas, archivo y desarchivo de proyectos y papelera recuperable para notas y proyectos.
- Las notas existentes se abren en vista previa con formato; la edición se activa de forma explícita.
- Imágenes locales PNG/JPEG/WebP, vista previa y vista ampliada.
- Temas gratuitos e importación validada de temas locales.
- Interfaz en inglés, coreano, japonés, chino simplificado, chino tradicional, alemán, francés y español, con avisos locales de privacidad y uso en los mismos idiomas.
- Notificaciones de cambios externos en los archivos y protección frente a conflictos mediante SHA256.
- Una biblioteca configurada común a todos los clientes; al repetir la inicialización se conserva la selección existente.
- Actualizaciones locales del núcleo y reversión de la versión del software, verificadas mediante sumas de comprobación y con protección de las ediciones sin guardar.
- Pestañas de notas y posiciones de lectura guardadas en el servidor, recuperables aunque cambie el puerto del servicio.

La protección frente a conflictos rechaza escrituras basadas en versiones antiguas; no hay un historial automático de revisiones. Haz copias de seguridad de la biblioteca por separado. Una pestaña abierta puede conservar un borrador durante un conflicto o una desconexión, pero no se garantiza recuperarlo tras cerrar la pestaña o reiniciar Windows.

![Thoughts 0.3.0-preview.20 con notas de demostración ficticias](docs/assets/preview-0.3.0-preview.20.jpg)

Interfaz real de `0.3.0-preview.20`, capturada el 2026-11-04 en una biblioteca temporal nueva y aislada con notas ficticias. No aparecen notas personales, imágenes reales, rutas locales del equipo ni credenciales. La [captura del diseño anterior](docs/assets/preview.png) se conserva como material histórico y no representa la versión actual. La comprobación de la captura no equivale a validar el instalador ni los clientes.

## Instalación

[INSTALL.md](INSTALL.md) incluye la verificación SHA256, los comandos que debes revisar antes de ejecutarlos y los pasos para los tres clientes. Desde preview.8, las instalaciones nuevas usan por defecto el mismo directorio Thoughts del usuario de Windows en todos los clientes, independientemente del orden de instalación. Mantén el ajuste avanzado del directorio base en `auto`. Solo hace falta una ruta absoluta `THOUGHTS_HOME` si deseas personalizarlo expresamente:

```text
<Thoughts home>/
  library/     Markdown, images and recoverable Trash
  config/      Configuration, local sessions, preferences and logs
  app/         Managed executable versions (optional)
  backups/     Previous runtime pins; not a backup of note bodies
```

En las instalaciones normales, usa el directorio compartido predeterminado. La configuración explícita antigua la configuración anterior sigue disponible; esta versión no migra ni combina automáticamente bibliotecas existentes. Si personalizas el directorio base, usa la misma carpeta física en los clientes que deban compartir contenido.

Las instalaciones y actualizaciones deben conservar ese directorio. El estado de lectura guarda rutas y posiciones de desplazamiento, no texto sin guardar. Actualmente las actualizaciones del núcleo requieren un ejecutable local revisado y su SHA256 indicados expresamente; no hay actualizaciones de red en segundo plano. La extensión del cliente y el núcleo compartido se actualizan por separado. Consulta las [instrucciones de actualización](INSTALL.md#managed-core-updates-and-rollback).

SmartScreen, Smart App Control o las políticas de control de aplicaciones de una organización pueden bloquear programas o scripts de Windows sin firmar. SHA256 confirma la coincidencia con una suma de verificación, no la identidad del editor. No desactives la protección ni eludas las advertencias para ejecutar una versión preliminar. Publicar el código fuente no resuelve la firma de código.

## Espacio de trabajo

Consulta la [organización del espacio de trabajo](docs/WORKSPACE_LAYOUT.md) y el [índice de documentación](docs/README.md). Los requisitos e informes anteriores se conservan en un archivo local. El código y los ejecutables usan Thoughts; se conservan los identificadores de datos anteriores por compatibilidad.

## Compilación

Se necesitan Windows x64, Go de acuerdo con go.mod, Node.js 24, Python 3.12+ y el [compilador oficial Inno Setup 7](https://jrsoftware.org/isdl.php) firmado. Los usuarios finales no necesitan estas herramientas. Consulta [CONTRIBUTING.md](CONTRIBUTING.md).

```powershell
Set-Location web
npm ci
Set-Location ..
pwsh -File scripts/build-desktop-adapters.ps1 -Python 'C:\Tools\Python\python.exe' -Go 'C:\Tools\Go\bin\go.exe' -InnoCompiler 'C:\Tools\InnoSetup7\ISCC.exe'
```

Sustituye las rutas de las herramientas por las de tu instalación. El proceso prueba y compila la interfaz y el núcleo Go, genera los avisos, crea Setup.exe/ZIP/MCPB y sus sumas de verificación y valida los paquetes en un directorio de salida nuevo. No sobrescribe resultados anteriores ni sube archivos automáticamente. El modo de candidato público exige firmas de confianza y metadatos reales del editor; la firma no sustituye la validación en un equipo limpio ni en los clientes.

## Seguridad y licencia

Consulta [SECURITY.md](SECURITY.md), [PRIVACY.md](PRIVACY.md), [CHANGELOG.md](CHANGELOG.md), los [avisos de terceros](THIRD_PARTY_NOTICES.txt) y la [revisión de dependencias](docs/DEPENDENCIES.md). El código del proyecto tiene licencia MIT; las dependencias incluidas conservan sus propias licencias, incluidas las fuentes OFL.

Nunca incluyas en commits notas reales, capturas personales, credenciales, archivos de sesión ni registros. Informa de las vulnerabilidades de forma privada a hello@spxipo.com. Los ajustes ofrecen apoyo voluntario mediante la dirección Base mainnet del propietario. El QR solo contiene la dirección; selecciona Base expresamente en tu monedero. Thoughts no conecta monederos, no solicita claves privadas, no elige activos ni importes y no procesa pagos. Todas las funciones siguen siendo gratuitas.

## Code signing policy

**Pendiente — la solicitud a SignPath Foundation todavía no se ha aprobado.** Todos los programas actuales de Windows siguen sin firmar. Si se aprueba, SignPath.io proporcionará la firma de código gratuita y SignPath Foundation el certificado. La [Code signing policy](docs/CODE_SIGNING_POLICY.md) y la [preparación de la solicitud](docs/SIGNPATH_APPLICATION.md) detallan el alcance, los responsables reales y los requisitos pendientes. Esta preparación no ha enviado ninguna solicitud.

## Donar

El apoyo es voluntario y no desbloquea funciones. La [página específica de donaciones](DONATE.md) muestra la dirección Base existente del propietario, el QR que contiene solo la dirección y los límites de privacidad. También puedes colaborar con incidencias, traducciones o documentación. Thoughts no conecta monederos ni procesa pagos.
