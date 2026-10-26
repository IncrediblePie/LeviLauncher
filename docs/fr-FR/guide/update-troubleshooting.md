# Mises à jour & Dépannage

Lorsque le comportement d'installation, de mise à jour ou de lancement de LeviLauncher ne correspond pas à vos attentes, veuillez consulter cette page en priorité.

## Problèmes courants

### Gaming Services manquant ou endommagé

- Utilisez le guide fourni par le lanceur pour installer ou réparer.
- Si nécessaire, retournez dans le Microsoft Store pour vérifier l'état d'installation de Minecraft Bedrock.

### GameInput manquant

- Installez le redistribuable requis selon les instructions de LeviLauncher.
- Après l'installation, redémarrez le lanceur si nécessaire.

### Chemin d'installation non accessible en écriture

- Dans **Paramètres → Paramètres généraux**, modifiez le chemin du contenu géré pour le placer sur un emplacement accessible en écriture.
- Si le lanceur demande des privilèges d'administrateur pour l'installation ou la mise à jour automatique, autorisez-le uniquement si vous faites confiance à la source.

### LeviLauncher se ferme immédiatement ou n'affiche jamais de fenêtre

- Sous Windows, LeviLauncher démarre désormais WebView2 dans un mode de rendu logiciel plus sûr afin de réduire les plantages au démarrage avec certaines combinaisons de GPU et de pilotes.
- Si le lanceur échoue toujours avant l'affichage de l'interface, consultez `%APPDATA%\levilauncher.exe\logs\startup.log`.
- Réessayez une fois après un lancement échoué. Si le problème persiste, joignez `startup.log` lorsque vous ouvrez une Issue sur GitHub.

### Une version ne démarre pas

- Testez d'abord si la version peut démarrer sans Mods
- Vérifiez que les composants Windows requis sont toujours présents
- Confirmez que la version est entièrement installée
- Essayez de reproduire le problème dans un environnement de test isolé et propre

### Téléchargement très lent ou échecs répétés

- Changez de miroir
- Réessayez plus tard
- Si possible, utilisez une source de paquet local

## Étapes de récupération plus sûres

1. Sauvegardez d'abord vos mondes importants.
2. Supprimez les Mods ou packs de ressources ajoutés récemment.
3. Effectuez un test de comparaison avec une version isolée propre.
4. Si nécessaire, réinstallez la version concernée.

## À propos de la mise à jour automatique

LeviLauncher peut vérifier, télécharger et installer ses propres mises à jour. Si le répertoire d'installation n'est pas accessible en écriture, certaines opérations pourraient nécessiter des privilèges d'administrateur.

## Quand créer une Issue

Il est recommandé de soumettre un problème sur GitHub lorsque les conditions suivantes sont remplies :

- Vous pouvez reproduire le problème de manière stable
- Le problème persiste après avoir appliqué les étapes de récupération de base
- Le problème semble provenir de LeviLauncher lui-même, et non d'un problème général de l'environnement Windows

Lors de la soumission, veuillez joindre autant que possible :

- La version de Windows
- La version de LeviLauncher
- Vos étapes précises pour reproduire le problème
- Des captures d'écran ou des journaux (logs) si disponibles

Lien pour signaler un problème : [GitHub Issues](https://github.com/LiteLDev/LeviLauncher/issues)