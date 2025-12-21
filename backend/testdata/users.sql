-- phpMyAdmin SQL Dump
-- version 5.2.2
-- https://www.phpmyadmin.net/
--
-- Host: mariadb:3306
-- Erstellungszeit: 27. Aug 2025 um 20:34
-- Server-Version: 11.4.8-MariaDB-ubu2404
-- PHP-Version: 8.2.27

SET SQL_MODE = "NO_AUTO_VALUE_ON_ZERO";
START TRANSACTION;
SET time_zone = "+00:00";


/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!40101 SET NAMES utf8mb4 */;

--
-- Datenbank: `bamort`
--

-- --------------------------------------------------------

--
-- Tabellenstruktur für Tabelle `users`
--

CREATE TABLE `users` (
  `user_id` bigint(20) UNSIGNED NOT NULL,
  `username` varchar(191) DEFAULT NULL,
  `password_hash` longtext DEFAULT NULL,
  `email` varchar(191) DEFAULT NULL,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `reset_pw_hash` varchar(191) DEFAULT NULL,
  `reset_pw_hash_expires` datetime(3) DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

--
-- Daten für Tabelle `users`
--

INSERT INTO `users` (`user_id`, `username`, `password_hash`, `email`, `created_at`, `updated_at`, `reset_pw_hash`, `reset_pw_hash_expires`) VALUES
(1, 'bebe', 'ee8b1834c791b857b5dd2b746a403000', 'frank@wuenscheonline.de', '2025-01-04 09:01:44.911', '2025-08-13 06:28:09.989', NULL, NULL),
(4, 'testuser', '179ad45c6ce2cb97cf1029e212046e81', '', '2025-07-20 21:00:28.601', '2025-07-20 21:00:28.601', NULL, NULL);

--
-- Indizes der exportierten Tabellen
--

--
-- Indizes für die Tabelle `users`
--
ALTER TABLE `users`
  ADD PRIMARY KEY (`user_id`),
  ADD UNIQUE KEY `uni_users_username` (`username`),
  ADD UNIQUE KEY `uni_users_email` (`email`),
  ADD KEY `idx_users_reset_pw_hash` (`reset_pw_hash`);

--
-- AUTO_INCREMENT für exportierte Tabellen
--

--
-- AUTO_INCREMENT für Tabelle `users`
--
ALTER TABLE `users`
  MODIFY `user_id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=6;
COMMIT;

/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
