CREATE TABLE emprunts (
    id_emprunt SERIAL PRIMARY KEY,         
    utilisateur_id INT NOT NULL,           -- ID de l'utilisateur
    livre_id INT NOT NULL,                 -- ID du livre
    date_emprunt TIMESTAMP NOT NULL,       -- Date d'emprunt
    date_retour_prevu TIMESTAMP NOT NULL,  -- Date prévue de retour
    date_retour_effectif TIMESTAMP NULL,   -- Date effective de retour (NULL tant que non retourné)
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    deleted_at TIMESTAMP NULL              -- Soft delete
);

CREATE TABLE penalites (
    id_penalite SERIAL PRIMARY KEY,        -- Identifiant unique
    utilisateur_id INT NOT NULL,           -- ID de l'utilisateur
    emprunt_id INT NOT NULL,               -- ID de l'emprunt
    montant DECIMAL(10, 2) NOT NULL,       -- Montant de la pénalité
    date_paiement TIMESTAMP NULL,          -- Date de paiement
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    deleted_at TIMESTAMP NULL,             -- Soft delete
    FOREIGN KEY (emprunt_id) REFERENCES emprunts (id_emprunt) ON DELETE CASCADE
);

-- Table des emprunts
INSERT INTO emprunts (utilisateur_id, livre_id, date_emprunt, date_retour_prevu, date_retour_effectif)
VALUES
(1, 101, '2024-11-01 10:00:00', '2024-11-15 10:00:00', NULL), -- Emprunt en cours
(2, 102, '2024-10-20 10:00:00', '2024-11-03 10:00:00', '2024-11-02 10:00:00'), -- Emprunt rendu à temps
(3, 103, '2024-10-15 10:00:00', '2024-10-29 10:00:00', '2024-11-01 10:00:00'), -- Emprunt rendu en retard
(4, 104, '2024-10-01 10:00:00', '2024-09-15 10:00:00', NULL), -- Emprunt en cours (en retard)
(1, 105, '2024-11-10 10:00:00', '2024-11-25 10:00:00', NULL); -- Nouvel emprunt en cours

-- Table des pénalités
INSERT INTO penalites (utilisateur_id, emprunt_id, montant, date_paiement, created_at, updated_at)
VALUES
(3, 3, 5.00, '2024-11-05 10:00:00', NOW(), NOW()), -- Pénalité payée
(4, 4, 10.00, NULL, NOW(), NOW()); -- Pénalité non payée pour un retard en cours
