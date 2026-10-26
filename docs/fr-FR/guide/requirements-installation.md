# Configuration requise & Installation

Cette page décrit les conditions nécessaires pour que LeviLauncher puisse installer et gérer normalement les versions de Minecraft Bedrock (GDK / UWP).

Les versions UWP nécessitent le mode développeur de Windows ainsi que les dépendances du framework UWP déclarées par le paquet (comme Microsoft.VCLibs). Les exigences Gaming Services et GameInput ci-dessous s'appliquent à GDK. L'isolation des instances fonctionne pour GDK et pour les paquets UWP x64 de version 1.19.70.2 ou ultérieure ; les instances UWP non isolées partagent le contenu `LocalState` de l'application Windows. Les instances UWP x64 prennent aussi en charge la console et les Mods DLL `preload-native` compatibles, mais l'installation automatique de LeviLamina n'est disponible que pour les versions GDK prises en charge.

## Exigences système

| Élément | Exigence |
| --- | --- |
| Système d'exploitation | Windows 10 ou Windows 11 |
| Version du jeu | Minecraft Bedrock Edition (GDK / UWP) |
| Licence | Licence officielle liée à un compte Microsoft |
| Réseau | Nécessaire pour télécharger les versions, obtenir les métadonnées, tester la vitesse des miroirs et vérifier les mises à jour |

## Composants Windows nécessaires

Avant le premier lancement ou l'installation, LeviLauncher pourrait vous guider pour installer les composants manquants.

- **Microsoft Gaming Services** (GDK)
- **Microsoft GameInput** (GDK)
- **WebView2 Runtime**
- **Microsoft Visual C++ Redistributable (x64)**, requis par le chargeur de Mods natif ; LeviLauncher vérifie sa présence et propose de l'installer

La présence ou non de ces composants dépend de l'état de votre environnement Windows.

### Utiliser un WebView2 Fixed Version Runtime personnalisé

Si le WebView2 Runtime Evergreen n'est pas installé mais que vous gérez vous-même un WebView2 Fixed Version Runtime, démarrez LeviLauncher avec l'une des méthodes suivantes.

Indiquez le répertoire en ligne de commande :

```powershell
LeviLauncher.exe --webview2-runtime-dir "D:\WebView2\FixedRuntime"
```

Ou créez un fichier `webview2.json` à côté de `LeviLauncher.exe` :

```json
{
  "browserExecutableFolder": "D:\\WebView2\\FixedRuntime"
}
```

Les chemins relatifs sont résolus à partir du répertoire contenant `webview2.json`. Le répertoire choisi doit contenir directement `msedgewebview2.exe`. L'option de ligne de commande est prioritaire sur le fichier de configuration. Si aucune des deux n'est configurée, LeviLauncher continue d'utiliser le WebView2 Runtime installé sur le système.

## Avant d'installer une version

Veuillez compléter cette liste de vérification :

1. Avoir installé Minecraft Bedrock au moins une fois depuis le Microsoft Store.
2. Si l'état du Store est anormal, lancez le jeu une fois pour confirmer que l'installation est complète.
3. Fermez le jeu avant d'utiliser LeviLauncher pour installer ou gérer des versions.

## Installer LeviLauncher lui-même

### Option A : Page GitHub Releases

Convient aux utilisateurs souhaitant télécharger directement depuis la page officielle de LeviLauncher et consulter l'historique des mises à jour.

1. Ouvrez la page [GitHub Releases](https://github.com/LiteLDev/LeviLauncher/releases) de LeviLauncher.
2. Téléchargez le programme d'installation.
3. Exécutez-le et suivez l'assistant d'installation.

### Option B : Miroir Lanzou Cloud

Si vous avez une connexion lente à GitHub depuis votre région, ce lien est généralement plus pratique.

1. Ouvrez [Lanzou Cloud](https://levimc.lanzoue.com/b016ke39hc).
2. Saisissez le mot de passe `levi`.
3. Téléchargez le fichier et exécutez le programme d'installation localement.

## Installer une première version gérée

1. Dans LeviLauncher, ouvrez **Télécharger**.
2. Choisissez la build Minecraft **Version stable** (Release) ou **Aperçu** (Preview) que vous souhaitez installer, ainsi que le type de paquet **GDK** ou **UWP**.
3. Sélectionnez l'entrée de la version cible.
4. Décidez d'activer ou non l'isolation.
5. Lancez l'installation et attendez sa fin.

## Stratégie d'installation recommandée

### Quand choisir la version Release (officielle)

- Vous souhaitez un environnement de jeu quotidien plus stable
- Vous gérez des mondes à long terme
- Vous préférez que les Mods et les packs de ressources changent moins souvent

### Quand choisir la version Preview (Aperçu)

- Vous voulez tester les futures fonctionnalités en avance
- Vous acceptez les instabilités ou les changements de compatibilité
- Vous êtes prêt à séparer clairement votre environnement Preview de votre environnement de jeu quotidien

::: tip Pratique recommandée pour la majorité des joueurs
Commencez par créer une **version Release (officielle) isolée**. N'ajoutez une **version Preview (Aperçu)** que lorsque vous avez clairement besoin de tester le contenu à venir.
:::

## Si l'installation ne peut pas se poursuivre

Les problèmes suivants peuvent être résolus en consultant la page [Mise à jour et résolution de problèmes](./update-troubleshooting) :

- Manque de Gaming Services
- Manque de GameInput
- État incomplet de la licence ou de l'installation via le Store
- Chemin de destination non accessible en écriture
- Échec du téléchargement ou du miroir