[English](README.md) | [简体中文](README.zh-CN.md) | [繁體中文](README.zh-TW.md) | [日本語](README.ja.md) | [한국어](README.ko.md) | [Deutsch](README.de.md) | [Français](README.fr.md) | [Español](README.es.md)

# Thoughts

[Installation](INSTALL.md) · [Documentation](docs/README.md) · [Faire un don](DONATE.md) · [Code signing policy](docs/CODE_SIGNING_POLICY.md)

Un outil de notes gratuit et open source pour Windows x64, qui privilégie le stockage local. Confiez la gestion de vos fichiers Markdown à votre assistant IA et consultez-les ou modifiez-les dans une page web locale indépendante.

**Aperçu du code source — 2026-11-04 :** ce dépôt publie le code révisé de `0.3.0-preview.20` sous licence MIT, avec un nouvel historique public. Les 59 tests du backend ont réussi nativement sous Windows. **Aucun programme d’installation Windows n’est encore approuvé pour le téléchargement :** SAC a bloqué le Setup.exe final non signé (erreur 4551). Le code de l’installateur est inclus, mais la diffusion des binaires attend une signature de confiance et la validation finale. Voir la [préparation de la distribution](docs/RELEASE_READINESS_20261104.md), l’[installation](INSTALL.md) et les [tests Windows](docs/WINDOWS_NATIVE_TESTS_20261104.md).

## Vos notes, sur votre ordinateur

Thoughts ne propose ni compte, ni base de données cloud, ni service d’envoi de notes, ni télémétrie, ni IA intégrée, ni thèmes payants, ni abonnement, ni intégration de paiement, ni extension de navigateur. Vos fichiers Markdown et vos images locales vous appartiennent et peuvent s’ouvrir dans d’autres éditeurs.

L’exécutable Go embarque la page React. **Les versions précompilées n’ont besoin ni de Node.js, ni de Go, ni de Python, ni de Docker à l’exécution.** Le client IA que vous choisissez a ses propres exigences d’installation et de compte.

## Comment l’installation sera proposée

Lorsqu’une version révisée sera disponible, téléchargez **`thoughts-<version>-windows-amd64-setup.exe`** depuis [GitHub Releases](https://github.com/xuli82221-pixel/hello-Thoughts-/releases), vérifiez sa somme de contrôle et la signature de l’éditeur, puis ouvrez-le par double-clic et choisissez **Install**. L’installateur initialise le dossier de base partagé et crée les raccourcis du bureau et du menu Démarrer sans droits administrateur. Les données existantes sont conservées. La désinstallation retire le lanceur et les raccourcis, mais garde le dossier partagé pour les autres clients. Une archive ZIP portable est également prévue pour l’installation manuelle dans cette version révisée ; aucun des deux paquets n’est actuellement un téléchargement approuvé.

Les utilisateurs de Claude Desktop pourront aussi importer l’extension `.mcpb` de cette version. Pour Codex et Claude Code, l’exécutable du ZIP s’enregistre comme serveur MCP local sur stdio. La page web fonctionne aussi sans client IA. Les commandes figurent dans le [guide de démarrage pour les utilisateurs](INSTALL.md#ordinary-windows-users). Sur GitHub, **Code → Download ZIP** télécharge le code source, pas l’application précompilée.

Les clients IA lancent le même adaptateur MCP local utilisant le protocole standard sur stdio. Il appelle le cœur local authentifié, dont la page web écoute uniquement sur `127.0.0.1`. Si vous demandez à un client IA en ligne de lire une note, il en reçoit le contenu selon sa propre politique de traitement des données. Privilégier le stockage local ne garantit pas que les résultats d’outils IA demandés restent sur votre ordinateur. Thoughts ne chiffre pas les fichiers.

## Clients visés

| Client | Installation |
| --- | --- |
| Claude Desktop | Format officiel d’extension `.mcpb` |
| Claude Code sous Windows | MCP local sur stdio, ZIP portable Windows |
| OpenAI Codex sous Windows | MCP local sur stdio, même ZIP portable |

Il s’agit des intégrations visées ; leur état réel de validation figure dans la revue de distribution. Les plugins web ordinaires de ChatGPT, macOS, Linux, les appareils mobiles, MCP distant et la synchronisation cloud sont hors du périmètre de cette phase.

## Fonctionnalités

- Créer, lire, modifier, rechercher, renommer et déplacer des fichiers Markdown locaux.
- Arborescence compacte des projets, notes rapides, archivage et désarchivage des projets, corbeille récupérable pour les notes et les projets.
- Ouverture des notes existantes en aperçu mis en forme ; passage explicite en mode édition.
- Images locales PNG/JPEG/WebP, aperçu et affichage agrandi.
- Thèmes gratuits et importation de thèmes locaux avec validation.
- Interface en anglais, coréen, japonais, chinois simplifié, chinois traditionnel, allemand, français et espagnol, avec notices locales de confidentialité et d’utilisation correspondantes.
- Notifications de changements de fichiers externes et protection contre les conflits par SHA256.
- Une bibliothèque configurée commune à tous les clients ; la réinitialisation conserve la sélection existante.
- Mises à niveau locales du cœur et retour à une version logicielle précédente, vérifiés par somme de contrôle, avec protection des modifications non enregistrées.
- Onglets de notes et positions de lecture enregistrés côté serveur, restaurables même après un changement de port du service.

La protection contre les conflits rejette les écritures fondées sur une ancienne version ; l’historique automatique des révisions n’est pas implémenté. Sauvegardez votre bibliothèque séparément. Un onglet ouvert peut conserver un brouillon en cas de conflit ou de déconnexion, mais sa récupération après fermeture ou redémarrage de Windows n’est pas garantie.

![Thoughts 0.3.0-preview.20 avec des notes de démonstration fictives](docs/assets/preview-0.3.0-preview.20.jpg)

Interface réelle de `0.3.0-preview.20`, capturée le 2026-11-04 dans une nouvelle bibliothèque temporaire isolée avec des notes fictives. Aucune note personnelle, image réelle, chemin propre à la machine ou donnée d’authentification n’est visible. La [capture de l’ancienne interface](docs/assets/preview.png) est conservée comme document historique et ne représente pas la version actuelle. La vérification de la capture ne vaut pas validation de l’installateur ou des clients.

## Installation

[INSTALL.md](INSTALL.md) explique les vérifications SHA256, les commandes à examiner avant exécution et la configuration des trois clients. Depuis preview.8, les nouvelles installations utilisent par défaut le même dossier Thoughts du compte Windows dans chaque client, quel que soit l’ordre d’installation. Laissez le réglage avancé du dossier de base sur `auto`. Un chemin absolu `THOUGHTS_HOME` n’est nécessaire que pour personnaliser volontairement ce dossier :

```text
<Thoughts home>/
  library/     Markdown, images and recoverable Trash
  config/      Configuration, local sessions, preferences and logs
  app/         Managed executable versions (optional)
  backups/     Previous runtime pins; not a backup of note bodies
```

Utilisez normalement le dossier partagé par défaut. L’ancien réglage explicite `MEMO_CONFIG_DIR` reste disponible ; cette version ne migre ni ne fusionne automatiquement une bibliothèque existante. Si vous personnalisez le dossier de base, tous les clients devant partager le contenu doivent utiliser le même dossier physique.

Les installations et mises à jour doivent préserver ce dossier. L’état de lecture stocke les chemins et positions de défilement, pas le texte non enregistré. Les mises à jour du cœur exigent actuellement un exécutable local révisé et son SHA256, fournis explicitement ; aucune mise à jour réseau ne s’effectue en arrière-plan. La mise à jour de l’extension cliente et celle du cœur partagé sont deux opérations distinctes. Voir les [instructions de mise à niveau](INSTALL.md#managed-core-updates-and-rollback).

SmartScreen, Smart App Control ou une politique de contrôle des applications de votre organisation peuvent bloquer les programmes ou scripts Windows non signés. SHA256 vérifie la concordance avec une somme de contrôle, pas l’identité de l’éditeur. Ne désactivez pas la protection et ne contournez pas les avertissements pour lancer un aperçu. Publier le code source ne résout pas la question de la signature.

## Espace de travail

Voir l’[organisation des dossiers](docs/WORKSPACE_LAYOUT.md) et l’[index de documentation](docs/README.md). Les anciens besoins et rapports sont préservés dans une archive locale. Les chemins internes Memo sont des noms de compatibilité, pas des produits supplémentaires ni des fichiers à supprimer librement.

## Compilation

Il faut Windows x64, une version de Go correspondant à go.mod, Node.js 24, Python 3.12+ et le [compilateur officiel Inno Setup 7](https://jrsoftware.org/isdl.php) signé. Les utilisateurs finaux n’ont besoin d’aucun de ces outils. Voir [CONTRIBUTING.md](CONTRIBUTING.md).

```powershell
Set-Location web
npm ci
Set-Location ..
pwsh -File scripts/build-desktop-adapters.ps1 -Python 'C:\Tools\Python\python.exe' -Go 'C:\Tools\Go\bin\go.exe' -InnoCompiler 'C:\Tools\InnoSetup7\ISCC.exe'
```

Remplacez les chemins des outils par leurs emplacements réels. Le processus teste et compile l’interface et le cœur Go, génère les notices, crée Setup.exe/ZIP/MCPB et leurs sommes de contrôle, puis valide les paquets dans un nouveau dossier de sortie. Les résultats précédents ne sont jamais écrasés. Aucun envoi n’est automatique. Le mode candidat public exige des signatures de confiance et de vraies métadonnées d’éditeur ; une signature ne remplace pas les validations sur une machine propre ou dans les clients.

## Sécurité et licence

Voir [SECURITY.md](SECURITY.md), [PRIVACY.md](PRIVACY.md), [CHANGELOG.md](CHANGELOG.md), les [notices des composants tiers](THIRD_PARTY_NOTICES.txt) et la [revue des dépendances](docs/DEPENDENCIES.md). Le code du projet est sous licence MIT ; les dépendances embarquées conservent leurs propres licences, notamment les polices sous OFL.

Ne commitez jamais de vraies notes, captures personnelles, données d’authentification, fichiers de session ou journaux. Signalez les vulnérabilités en privé à hello@spxipo.com. Les paramètres proposent un soutien facultatif à l’adresse Base mainnet du propriétaire. Le QR contient uniquement l’adresse ; choisissez explicitement Base dans votre portefeuille. Thoughts ne connecte pas de portefeuille, ne demande pas de clé privée, ne sélectionne ni actif ni montant et ne traite aucun paiement. Toutes les fonctionnalités restent gratuites.

## Code signing policy

**En attente — la demande auprès de la SignPath Foundation n’est pas encore approuvée.** Tous les programmes Windows actuels restent non signés. En cas d’approbation, la signature de code gratuite sera fournie par SignPath.io, et le certificat par SignPath Foundation. La [Code signing policy](docs/CODE_SIGNING_POLICY.md) et la [préparation de la demande](docs/SIGNPATH_APPLICATION.md) précisent le périmètre, les responsabilités réelles et les exigences restantes. Aucune demande n’a été envoyée dans le cadre de cette préparation.

## Faire un don

Le soutien est volontaire et ne débloque aucune fonctionnalité. La [page dédiée aux dons](DONATE.md) présente l’adresse Base existante du propriétaire, le QR contenant uniquement l’adresse et les limites de confidentialité. Vous pouvez aussi aider par des signalements, des traductions ou de la documentation. Thoughts ne connecte pas de portefeuille et ne traite aucun paiement.
