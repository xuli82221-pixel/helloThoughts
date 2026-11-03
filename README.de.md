[English](README.md) | [简体中文](README.zh-CN.md) | [繁體中文](README.zh-TW.md) | [日本語](README.ja.md) | [한국어](README.ko.md) | [Deutsch](README.de.md) | [Français](README.fr.md) | [Español](README.es.md)

# Thoughts

[Installation](INSTALL.md) · [Dokumentation](docs/README.md) · [Spenden](DONATE.md) · [Code signing policy](docs/CODE_SIGNING_POLICY.md)

Kostenlose, quelloffene Notizen für Windows x64 mit lokaler Datenhaltung. Lass deinen KI-Assistenten Markdown-Dateien verwalten und lies und bearbeite sie auf einer eigenständigen lokalen Webseite.

**Quellcode-Vorschau — 2026-11-04:** Dieses Repository veröffentlicht den geprüften Quellcode von `0.3.0-preview.20` unter der MIT-Lizenz mit einer neuen öffentlichen Git-Historie. Alle 59 Backend-Tests wurden nativ unter Windows bestanden. **Es gibt noch keinen freigegebenen Windows-Installer zum Herunterladen:** SAC hat die endgültige unsignierte Setup.exe blockiert (Fehler 4551). Der Installer-Quellcode ist enthalten; Binärveröffentlichungen erfordern noch eine vertrauenswürdige Signatur und die abschließende Abnahme. Siehe [Veröffentlichungsvorbereitung](docs/RELEASE_READINESS_20261104.md), [Installation](INSTALL.md) und [Windows-Tests](docs/WINDOWS_NATIVE_TESTS_20261104.md).

## Deine Notizen auf deinem Computer

Thoughts hat keine Konten, Cloud-Datenbank, Upload-Dienste für Notizen, Telemetrie, integrierte KI, kostenpflichtigen Designs, Abonnements, Zahlungsintegration oder Browsererweiterung. Markdown-Dateien und lokale Bilder gehören dir und lassen sich auch mit anderen Editoren öffnen.

Die Go-Anwendung enthält die React-Webseite. **Vorkompilierte Pakete benötigen zur Laufzeit weder Node.js noch Go, Python oder Docker.** Für deinen gewählten KI-Client gelten dessen eigene Installations- und Kontovoraussetzungen.

## So wird die Installation möglich sein

Sobald eine geprüfte Veröffentlichung verfügbar ist, lade **`thoughts-<version>-windows-amd64-setup.exe`** von [GitHub Releases](https://github.com/xuli82221-pixel/hello-Thoughts-/releases) herunter, prüfe Prüfsumme und Herausgebersignatur, öffne die Datei per Doppelklick und wähle **Install**. Das Setup initialisiert das gemeinsame Basisverzeichnis und erstellt Desktop- und Startmenüverknüpfungen ohne Administratorrechte. Vorhandene Daten bleiben erhalten. Die Deinstallation entfernt Programmstarter und Verknüpfungen; das gemeinsame Basisverzeichnis bleibt für andere Clients bestehen. Für diese Veröffentlichung ist auch ein portables ZIP zur manuellen Installation vorgesehen. Beide Pakete sind derzeit noch keine freigegebenen Downloads.

Claude-Desktop-Nutzer können dann alternativ die `.mcpb`-Erweiterung der Veröffentlichung importieren. Codex und Claude Code registrieren die ausführbare Datei aus dem ZIP als lokalen stdio-MCP-Server. Die Webseite funktioniert auch ohne KI-Client. Die Befehle stehen im [Schnellstart für Anwender](INSTALL.md#ordinary-windows-users). **Code → Download ZIP** auf GitHub lädt Quellcode herunter, keine vorkompilierte Anwendung.

Die KI-Clients starten denselben lokalen Adapter für das standardisierte stdio-MCP-Protokoll. Er spricht den authentifizierten lokalen Kern an; die Webseite lauscht ausschließlich auf `127.0.0.1`. Bittest du einen Online-KI-Client, eine Notiz zu lesen, erhält dieser ihren Inhalt und verarbeitet ihn nach seinen eigenen Datenrichtlinien. Lokale Datenhaltung bedeutet nicht, dass angeforderte KI-Werkzeugergebnisse den Computer nie verlassen. Thoughts verschlüsselt die Dateien nicht.

## Vorgesehene Clients

| Client | Installation |
| --- | --- |
| Claude Desktop | Offizielles `.mcpb`-Erweiterungsformat |
| Claude Code unter Windows | Lokales stdio MCP, portables Windows-ZIP |
| OpenAI Codex unter Windows | Lokales stdio MCP, dasselbe portable ZIP |

Dies sind die vorgesehenen Integrationen. Der tatsächliche Abnahmestand ist im Veröffentlichungsbericht dokumentiert. Gewöhnliche ChatGPT-Web-Plugins, macOS, Linux, Mobilgeräte, Remote-MCP und Cloud-Synchronisation gehören nicht zu dieser Phase.

## Funktionen

- Lokale Markdown-Dateien erstellen, lesen, bearbeiten, durchsuchen, umbenennen und verschieben.
- Kompakter Projektbaum, Schnellnotizen, Projekte archivieren und wieder aktivieren sowie wiederherstellbarer Papierkorb für Notizen und Projekte.
- Vorhandene Notizen öffnen in der formatierten Vorschau; die Bearbeitung wird ausdrücklich aktiviert.
- Lokale PNG-/JPEG-/WebP-Bilder, Vorschau und vergrößerte Ansicht.
- Kostenlose Designs und geprüfter Import lokaler Designs.
- Oberfläche auf Englisch, Koreanisch, Japanisch, vereinfachtem und traditionellem Chinesisch, Deutsch, Französisch und Spanisch mit passenden lokalen Datenschutz- und Nutzungshinweisen.
- Benachrichtigungen bei externen Dateiänderungen und SHA256-Konfliktschutz.
- Eine gemeinsam konfigurierte Bibliothek für alle Clients; erneute Initialisierung behält die bestehende Auswahl bei.
- Prüfsummengeprüfte lokale Kernaktualisierung und Software-Rollback mit Schutz ungespeicherter Änderungen.
- Serverseitig gespeicherte Notiz-Tabs und Lesepositionen, die auch nach einem Wechsel des Dienstports wiederhergestellt werden.

Der Konfliktschutz weist Schreibzugriffe auf veralteter Grundlage zurück; eine automatische Versionshistorie gibt es noch nicht. Sichere deine Bibliothek separat. Ein offener Tab kann bei Konflikten oder Verbindungsabbruch einen Entwurf behalten. Seine Wiederherstellung nach dem Schließen oder einem Windows-Neustart ist nicht garantiert.

![Thoughts 0.3.0-preview.20 mit fiktiven Beispielnotizen](docs/assets/preview-0.3.0-preview.20.jpg)

Tatsächliche Oberfläche von `0.3.0-preview.20`, aufgenommen am 2026-11-04 in einer neuen, isolierten temporären Bibliothek mit fiktiven Notizen. Persönliche Notizen, echte Bilder, lokale Rechnerpfade und Zugangsdaten sind nicht zu sehen. Der [Screenshot des früheren Layouts](docs/assets/preview.png) bleibt als historisches Material erhalten und stellt nicht die aktuelle Version dar. Die Screenshot-Prüfung ersetzt keine Installer- oder Client-Abnahme.

## Installation

[INSTALL.md](INSTALL.md) enthält SHA256-Prüfung, vor der Ausführung zu prüfende Befehle und Anleitungen für alle drei Clients. Seit preview.8 verwenden Neuinstallationen unabhängig von der Reihenfolge standardmäßig dasselbe Thoughts-Verzeichnis des Windows-Benutzers. Lass die erweiterte Home-Einstellung auf `auto`. Ein absoluter `THOUGHTS_HOME`-Pfad ist nur für ein bewusst abweichendes Basisverzeichnis erforderlich:

```text
<Thoughts home>/
  library/     Markdown, images and recoverable Trash
  config/      Configuration, local sessions, preferences and logs
  app/         Managed executable versions (optional)
  backups/     Previous runtime pins; not a backup of note bodies
```

Verwende normalerweise den gemeinsamen Standard. Das ausdrücklich gesetzte ältere `MEMO_CONFIG_DIR` bleibt verfügbar; diese Version migriert oder vereinigt bestehende Bibliotheken nicht automatisch. Bei einem eigenen Basisverzeichnis müssen Clients, die Inhalte teilen sollen, denselben physischen Ordner verwenden.

Installation und Aktualisierungen müssen dieses Verzeichnis erhalten. Der Lesestatus speichert Pfade und Scrollpositionen, keinen ungespeicherten Text. Kernaktualisierungen benötigen derzeit eine ausdrücklich angegebene, geprüfte lokale Programmdatei samt SHA256; es gibt keinen automatischen Netzwerk-Updater im Hintergrund. Client-Erweiterung und gemeinsamer Kern werden getrennt aktualisiert. Siehe [Aktualisierung und Rollback](INSTALL.md#managed-core-updates-and-rollback).

Unsignierte Windows-Programme oder Skripte können von SmartScreen, Smart App Control oder Richtlinien der Organisation blockiert werden. SHA256 belegt die Übereinstimmung mit einer Prüfsumme, nicht die Identität des Herausgebers. Deaktiviere keine Schutzmaßnahmen und umgehe keine Warnungen, um eine Vorschau auszuführen. Die Veröffentlichung des Quellcodes löst das Signaturproblem nicht.

## Arbeitsbereich

Siehe [Verzeichnisstruktur](docs/WORKSPACE_LAYOUT.md) und [Dokumentationsindex](docs/README.md). Frühere Anforderungen und Berichte bleiben in einem lokalen Archiv erhalten. Interne Memo-Pfade sind Kompatibilitätsnamen, keine zusätzlichen Produkte oder entbehrlichen Dateien.

## Selbst bauen

Benötigt werden Windows x64, Go passend zu go.mod, Node.js 24, Python 3.12+ und der signierte offizielle [Inno-Setup-7-Compiler](https://jrsoftware.org/isdl.php). Anwender benötigen diese Werkzeuge nicht. Siehe [CONTRIBUTING.md](CONTRIBUTING.md).

```powershell
Set-Location web
npm ci
Set-Location ..
pwsh -File scripts/build-desktop-adapters.ps1 -Python 'C:\Tools\Python\python.exe' -Go 'C:\Tools\Go\bin\go.exe' -InnoCompiler 'C:\Tools\InnoSetup7\ISCC.exe'
```

Ersetze die Werkzeugpfade durch deine Installationspfade. Der Ablauf testet und baut Oberfläche und Go-Kern, erzeugt Hinweise, erstellt Setup.exe/ZIP/MCPB samt Prüfsummen und prüft die Pakete in einem neuen Ausgabeverzeichnis. Vorhandene Ausgaben werden nicht überschrieben. Es erfolgt kein automatischer Upload. Der Modus für öffentliche Kandidaten verlangt vertrauenswürdige Signaturen und echte Herausgebermetadaten. Signaturen ersetzen weder Tests auf einem sauberen System noch die Client-Abnahme.

## Sicherheit und Lizenz

Siehe [SECURITY.md](SECURITY.md), [PRIVACY.md](PRIVACY.md), [CHANGELOG.md](CHANGELOG.md), [Drittanbieterhinweise](THIRD_PARTY_NOTICES.txt) und [Abhängigkeitsprüfung](docs/DEPENDENCIES.md). Der Projektcode steht unter MIT; enthaltene Abhängigkeiten behalten ihre eigenen Lizenzen, einschließlich der OFL-Schriftarten.

Committe niemals echte Notizen, persönliche Screenshots, Zugangsdaten, Sitzungsdateien oder Logs. Melde Sicherheitslücken vertraulich an hello@spxipo.com. Die Einstellungen zeigen eine freiwillige Unterstützungsmöglichkeit über die Base-Mainnet-Adresse des Eigentümers. Der QR-Code enthält nur die Adresse; wähle Base im Wallet ausdrücklich aus. Thoughts verbindet keine Wallets, fragt keine privaten Schlüssel ab, legt weder Asset noch Betrag fest und verarbeitet keine Zahlungen. Alle Funktionen bleiben kostenlos.

## Code signing policy

**Ausstehend — der Antrag bei der SignPath Foundation ist noch nicht genehmigt.** Alle aktuellen Windows-Programme sind weiterhin unsigniert. Bei Genehmigung stellt SignPath.io die kostenlose Codesignierung bereit, das Zertifikat die SignPath Foundation. Die [Code signing policy](docs/CODE_SIGNING_POLICY.md) und [Antragsvorbereitung](docs/SIGNPATH_APPLICATION.md) nennen Umfang, tatsächliche Zuständigkeiten und offene Anforderungen. Im Rahmen dieser Vorbereitung wurde kein Antrag eingereicht.

## Spenden

Unterstützung ist freiwillig und schaltet keine Funktionen frei. Die eigene [Spendenseite](DONATE.md) enthält die bestehende Base-Adresse des Eigentümers, den reinen Adress-QR-Code und die Datenschutzgrenzen. Du kannst auch mit Fehlermeldungen, Übersetzungen und Dokumentation helfen. Thoughts verbindet keine Wallets und verarbeitet keine Zahlungen.
