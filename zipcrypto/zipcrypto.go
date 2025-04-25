package zipcrypto

import (
    "encoding/binary"
    "crypto/aes"
    "crypto/cipher"
    "crypto/rand"
    "encoding/json"
    "fmt"
    "io/ioutil"
    "os"
    "time"
    "golang.org/x/crypto/scrypt"
)

type metadata struct {
    OriginalPath string
    Timestamp    time.Time
    Comment      string
}

func EncryptFile(zipFilePath, outPath, password string) error {
    zipData, err := ioutil.ReadFile(zipFilePath)
    if err != nil {
        return fmt.Errorf("impossible de lire le fichier ZIP : %v", err)
    }

    salt := make([]byte, 16)
    _, err = rand.Read(salt)
    if err != nil {
        return fmt.Errorf("échec de la génération du sel : %v", err)
    }

    key, err := scrypt.Key([]byte(password), salt, 32768, 8, 1, 32)
    if err != nil {
        return fmt.Errorf("échec de la génération de la clé : %v", err)
    }

    block, err := aes.NewCipher(key)
    if err != nil {
        return fmt.Errorf("échec de la création du chiffreur AES : %v", err)
    }

    gcm, err := cipher.NewGCM(block)
    if err != nil {
        return fmt.Errorf("échec de la création du chiffreur GCM : %v", err)
    }

    nonce := make([]byte, gcm.NonceSize())
    _, err = rand.Read(nonce)
    if err != nil {
        return fmt.Errorf("échec de la génération du nonce : %v", err)
    }

    ciphertext := gcm.Seal(nonce, nonce, zipData, nil)

    meta := metadata{
        OriginalPath: zipFilePath,
        Timestamp:    time.Now(),
        Comment:      "Fichier chiffré avec SecureZip",
    }
    metaBytes, err := json.Marshal(meta)
    if err != nil {
        return fmt.Errorf("échec de la sérialisation des métadonnées : %v", err)
    }

    metaSize := make([]byte, 4)
    binary.BigEndian.PutUint32(metaSize, uint32(len(metaBytes)))

    ciphertext = append(ciphertext, metaBytes...)
    ciphertext = append(ciphertext, metaSize...)

    outFile, err := os.Create(outPath)
    if err != nil {
        return fmt.Errorf("échec de la création du fichier de sortie : %v", err)
    }
    defer outFile.Close()

    _, err = outFile.Write(append(salt, ciphertext...))
    if err != nil {
        return fmt.Errorf("échec de l'écriture des données chiffrées : %v", err)
    }

    return nil
}

func DecryptFile(encryptedPath, outputPath, password string) error {
    encFile, err := os.Open(encryptedPath)
    if err != nil {
        return fmt.Errorf("failed to open encrypted file: %v", err)
    }
    defer encFile.Close()

    encData, err := ioutil.ReadAll(encFile)
    if err != nil {
        return fmt.Errorf("failed to read encrypted file: %v", err)
    }

    if len(encData) < 16 {
        return fmt.Errorf("invalid encrypted file: missing salt")
    }

    salt := encData[:16]
    ciphertext := encData[16:]

    key, err := scrypt.Key([]byte(password), salt, 32768, 8, 1, 32)
    if err != nil {
        return fmt.Errorf("failed to generate key: %v", err)
    }

    block, err := aes.NewCipher(key)
    if err != nil {
        return fmt.Errorf("failed to create AES cipher: %v", err)
    }

    gcm, err := cipher.NewGCM(block)
    if err != nil {
        return fmt.Errorf("failed to create GCM cipher: %v", err)
    }

    nonceSize := gcm.NonceSize()
    if len(ciphertext) < nonceSize {
        return fmt.Errorf("invalid encrypted file: missing nonce")
    }

    nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]

    if len(ciphertext) < 4 {
        return fmt.Errorf("invalid encrypted file: missing metadata size")
    }

    metaSize := binary.BigEndian.Uint32(ciphertext[len(ciphertext)-4:])
    if len(ciphertext) < int(metaSize)+4 {
        return fmt.Errorf("invalid encrypted file: metadata size mismatch")
    }

    metaBytes := ciphertext[len(ciphertext)-4-int(metaSize) : len(ciphertext)-4]
    var meta metadata
    err = json.Unmarshal(metaBytes, &meta)
    if err != nil {
        return fmt.Errorf("failed to extract metadata: %v", err)
    }

    filesData := ciphertext[:len(ciphertext)-4-int(metaSize)]

    decrypted, err := gcm.Open(nil, nonce, filesData, nil)
    if err != nil {
        return fmt.Errorf("failed to decrypt data: %v", err)
    }

    outFile, err := os.Create(outputPath)
    if err != nil {
        return fmt.Errorf("failed to create output file: %v", err)
    }
    defer outFile.Close()

    _, err = outFile.Write(decrypted)
    if err != nil {
        return fmt.Errorf("failed to write decrypted content: %v", err)
    }

    return nil
}