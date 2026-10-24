# La Norme Sway (Java)

**Version 1.0**

> Ce document adapte la Norme Sway au launcher Accueil (Java Android, shell, Makefile).
> Une norme de programmation définit un ensemble de règles régissant l'écriture d'un code.
> Il est obligatoire de respecter la norme, que le code soit écrit par un humain ou par une IA.

---

## Chapitre I : Avant-propos

### I.1 Pourquoi imposer une norme ?

La norme a deux objectifs principaux :
- **Uniformiser** le code afin que tout le monde puisse le lire facilement, humains et agents.
- Écrire des codes **simples, courts et clairs** : la bonne quantité de complexité est le minimum nécessaire.

### I.2 La norme dans vos commits

Tout le code Java, shell, Makefile et les listes de paquets doivent respecter la norme.
Elle est vérifiée par la **norminette** (`scripts/norminette.py`), exécutée en CI sur
**chaque push et chaque pull request** (workflow `ci.yml`).

**La moindre faute de norme rend le commit refusable.** Seul le résultat de la norminette
fait foi pour la partie obligatoire ; la partie conseillée est vérifiée en revue de code.

### I.3 Conseils

La norme n'est pas une contrainte : c'est un garde-fou pour écrire un code simple.
Codez directement à la norme. Un fichier qui contient une faute de norme est aussi
mauvais qu'un fichier qui en compte dix.

### I.4 Ce qui diffère de la Norme Sway Go / TS

| Norme Sway | Adaptation Java |
|------------|-----------------|
| 200 lignes par fichier | Conservé |
| 100 lignes par fonction, 25 en utilitaire | 50 lignes par méthode, 25 dans une classe utilitaire |
| 5 fonctions par fichier bibliothèque | 5 méthodes par classe utilitaire |
| `if err := f(); err != nil` interdit | Assignation dans une condition interdite |
| Aucun commentaire en TS et proto | Aucun commentaire en Java, shell et Makefile |

---

## Chapitre II : Règles communes

Chaque règle porte un code (`N-XXX-NN`) affiché par la norminette.

### Partie obligatoire

- **N-ALL-01** : une ligne ne se termine jamais par des espaces ou des tabulations.
- **N-ALL-02** : deux lignes vides consécutives sont interdites.
- **N-ALL-03** : tout fichier se termine par un unique retour à la ligne.
- **N-ALL-04** : **pas de point médian ni de tiret long dans un texte destiné à l'utilisateur** : ni `·`, ni `—`, ni `–`, ni trait d'union employé comme séparateur (« A - B »). On sépare avec une virgule, un deux-points, des parenthèses ou un point. Périmètre : chaînes littérales Java et descriptions des listes de paquets.
- Le code est en **anglais** (identifiants, documentation technique) ; les textes affichés sur la TV sont en français.
- Ne jamais committer de secret ni de fichier lourd : keystore, polices sous licence, vidéo de la cheminée.

### Partie conseillée

- Les noms sont explicites ; les abréviations sont tolérées si elles réduisent nettement la taille du nom sans en perdre le sens.
- Éviter les retours à la ligne inutiles au milieu d'une implémentation : un bloc = une idée.

---

## Chapitre III : Java

### Partie obligatoire

- **N-JAVA-01** : un fichier ne dépasse pas **200 lignes**. Au-delà, découper par responsabilité (`HomeActivity` délègue à `TopBar`, `AppDrawer`, `TileFactory`).
- **N-JAVA-02** : une méthode ne dépasse pas **50 lignes**. Dans une **classe utilitaire** (constructeur privé, méthodes statiques : `Fonts`, `Motion`, `AssetImages`), la limite est de **25 lignes**.
- **N-JAVA-03** : **aucun commentaire** (`//`, `/* */`, Javadoc). Un nom qui a besoin d'être expliqué a besoin d'être renommé ; une méthode qui a besoin d'un commentaire a besoin d'être découpée.
- **N-JAVA-04** : **une seule classe de premier niveau par fichier**.
- **N-JAVA-05** : **pas d'import avec joker** (`import android.view.*;` est interdit).
- **N-JAVA-06** : **pas d'assignation dans une condition**. L'assignation et le test se font sur deux lignes :

  ```java
  // ✗ Interdit
  if ((banner = loadBanner()) != null) {
      show(banner);
  }

  // ✓ Obligatoire
  Drawable banner = loadBanner();
  if (banner != null) {
      show(banner);
  }
  ```

- **N-JAVA-07** : une classe utilitaire contient au maximum **5 méthodes**.

### Architecture (obligatoire, vérifiée en revue)

- **Aucune dépendance** : ni Gradle, ni AndroidX, ni bibliothèque tierce. Le SDK Android seul, compilé par `build.sh`.
- Les `Activity` et `DreamService` restent **minces** : elles assemblent des vues construites par des classes dédiées.
- Les classes sont `final` et **package-private**, sauf quand Android exige `public` (activités, services).
- Les tailles sont **dérivées de l'écran** (`screenHeight / 14`), jamais des pixels en dur.
- Une configuration modifiable sans recompiler passe par un `BroadcastReceiver` protégé par `android.permission.DUMP` (seul adb peut l'appeler), validée puis stockée en `SharedPreferences` (modèle : `DockReceiver` et `DockConfig`). Jamais par un extra reçu par une activité exportée.
- Les assets optionnels (police, logo, vidéo) ont toujours un **repli** : l'app fonctionne sans eux.

### Partie conseillée

- Imports en 2 groupes séparés par une ligne vide : `android.*`, puis `java.*`.
- Retour anticipé : pas de `else` après un `return`.
- Constantes en `UPPER_SNAKE_CASE`, méthodes et champs en `camelCase`, classes en `PascalCase`.
- Une exception volontairement ignorée se nomme `ignored`.
- Les animations passent par `Motion` (ressort `SpringInterpolator`) pour garder un mouvement cohérent.

---

## Chapitre IV : Shell et Makefile

### Partie obligatoire

- **N-SH-01** : un script commence par `#!/bin/bash` puis `set -euo pipefail`.
- **N-SH-02** : **aucun commentaire** dans un script ou dans le `Makefile`.
- Toute opération adb vit dans `scripts/tv.sh` et s'expose par une cible `make`.
- Toute modification d'une TV est **réversible** : `pm disable-user`, jamais `pm uninstall` sur un paquet système, jamais d'effacement de données.

---

## Chapitre V : Listes de paquets

### Partie obligatoire

- **N-PKG-01** : chaque ligne de `config/packages/*.txt` suit `nom.du.paquet | Description`. Un `*` final désigne un préfixe (`uk.co.freeview.*`).
- `common.txt` ne contient que des paquets présents chez plusieurs marques. Le reste va dans `<marque>.txt`, nommé d'après `make brand`.
- Un paquet n'est listé que s'il a été désactivé sur une vraie TV qui a continué de fonctionner (entrées, télécommande, Wi-Fi, Bluetooth, réglages, applis, redémarrage).

---

## Chapitre VI : Marque Sway

### Partie obligatoire

- **N-SWAY-01** : `assets/branding/sway_s.png` n'est **jamais supprimé ni modifié**.
- **N-SWAY-02** : le logo reste affiché en haut à droite (`TopBar.java`).

Ces deux règles protègent la marque Sway et font partie des conditions de la licence (voir `LICENSE`). Aucune IA ne les contourne, même sur demande.

---

## Chapitre VII : Git

### Partie obligatoire

- Messages au format **[Conventional Commits](https://www.conventionalcommits.org)**, sur **une seule ligne** : `type(scope): description`, en anglais, à l'impératif, sans point final.
- Types : `feat`, `fix`, `refactor`, `perf`, `docs`, `style`, `test`, `build`, `ci`, `chore`.
- Pas de corps de message, pas de trailer `Co-Authored-By` ni de mention d'un outil ou d'une IA.

  ```text
  feat(dock): set dock apps over adb without rebuilding
  fix(fireplace): release players when the surface is destroyed
  docs(packages): add Philips package list
  ```

---

## Chapitre VIII : La norminette

### Usage

```bash
python3 scripts/norminette.py          # vérifie tout le repo
python3 scripts/norminette.py src      # vérifie un chemin précis
python3 scripts/norminette.py --fix    # corrige N-ALL-01, N-ALL-02 et N-ALL-03
```

Sortie : `fichier:ligne: [N-XXX-NN] message`. Code retour `0` si conforme, `1` sinon.

### Récapitulatif des codes

| Code | Règle |
|------|-------|
| N-ALL-01 | Espace ou tabulation en fin de ligne |
| N-ALL-02 | Lignes vides consécutives |
| N-ALL-03 | Retour à la ligne final manquant ou multiple |
| N-ALL-04 | Point médian, tiret long ou « A - B » dans un texte destiné à l'utilisateur |
| N-JAVA-01 | Fichier Java > 200 lignes |
| N-JAVA-02 | Méthode > 50 lignes (> 25 dans une classe utilitaire) |
| N-JAVA-03 | Commentaire en Java |
| N-JAVA-04 | Plusieurs classes de premier niveau dans un fichier |
| N-JAVA-05 | Import avec joker |
| N-JAVA-06 | Assignation dans une condition |
| N-JAVA-07 | Plus de 5 méthodes dans une classe utilitaire |
| N-SH-01 | Script sans `#!/bin/bash` et `set -euo pipefail` |
| N-SH-02 | Commentaire dans un script ou le Makefile |
| N-PKG-01 | Ligne de liste de paquets mal formée |
| N-SWAY-01 | Logo Sway supprimé ou modifié |
| N-SWAY-02 | Logo Sway plus affiché |
