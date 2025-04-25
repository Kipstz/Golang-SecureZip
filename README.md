
# SecureZip

SecureZip est une application de chiffrement et déchiffrement de fichiers, développée en Go avec une interface graphique basée sur Fyne. Elle utilise l'algorithme AES-GCM pour garantir la sécurité des données.

## Fonctionnalités

- **Chiffrement sécurisé** : Protégez vos fichiers avec un mot de passe.
- **Déchiffrement facile** : Récupérez vos fichiers en toute simplicité.
- **Interface intuitive** : Une interface graphique conviviale pour une utilisation rapide.
- **Métadonnées intégrées** : Conserve des informations sur le fichier original.

## Prérequis

- **Go** (version 1.18 ou supérieure)
- **Fyne** (version 2.0 ou supérieure)
- Un environnement **Windows**, **macOS** ou **Linux**.

## Installation

1. Clonez le dépôt :
   ```bash
   git clone https://github.com/Kipstz/Golang-SecureZip.git
   cd Golang-SecureZip
   ```

2. Installez les dépendances :
   ```bash
   go mod tidy
   ```

3. Compilez l'application pour votre système d'exploitation :
   - Pour Windows :
     ```bash
     fyne package -os windows -icon icon.png
     ```
   - Pour macOS :
     ```bash
     fyne package -os darwin -icon icon.png
     ```
   - Pour Linux :
     ```bash
     fyne package -os linux -icon icon.png
     ```

## Utilisation

1. Lancez l'application :
   ```bash
   go run main.go
   ```

2. **Chiffrement** :
   - Sélectionnez un fichier à chiffrer.
   - Entrez un mot de passe (minimum 6 caractères).
   - Cliquez sur "Chiffrer". Le fichier chiffré sera généré avec l'extension `.enc`.

3. **Déchiffrement** :
   - Sélectionnez un fichier `.enc`.
   - Entrez le mot de passe utilisé pour le chiffrement.
   - Cliquez sur "Déchiffrer". Le fichier original sera restauré.

## Structure du projet

```
securezip/
├── main.go          # Point d'entrée de l'application
├── ui/              # Interface utilisateur (Fyne)
│   └── ui.go
├── zipcrypto/       # Logique de chiffrement et déchiffrement
│   └── zipcrypto.go
└── icon.png         # Icône de l'application
```

## Contribution

Les contributions sont les bienvenues ! Veuillez ouvrir une issue ou soumettre une pull request pour proposer des améliorations.

---

**Note** : Si vous avez des questions, n'hésitez pas à ouvrir une issue ou à contacter l'auteur du projet.
