-- phpMyAdmin SQL Dump
-- version 5.2.2
-- https://www.phpmyadmin.net/
--
-- Host: mariadb:3306
-- Erstellungszeit: 27. Aug 2025 um 21:37
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
-- Tabellenstruktur für Tabelle `char_char_creation_session`
--

CREATE TABLE `char_char_creation_session` (
  `id` varchar(191) NOT NULL,
  `user_id` bigint(20) UNSIGNED NOT NULL,
  `name` longtext DEFAULT NULL,
  `rasse` longtext DEFAULT NULL,
  `typ` longtext DEFAULT NULL,
  `herkunft` longtext DEFAULT NULL,
  `glaube` longtext DEFAULT NULL,
  `attributes` text DEFAULT NULL,
  `derived_values` text DEFAULT NULL,
  `skills` text DEFAULT NULL,
  `spells` text DEFAULT NULL,
  `skill_points` text DEFAULT NULL,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `expires_at` datetime(3) DEFAULT NULL,
  `current_step` bigint(20) DEFAULT NULL,
  `geschlecht` longtext DEFAULT NULL,
  `stand` longtext DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

--
-- Daten für Tabelle `char_char_creation_session`
--

INSERT INTO `char_char_creation_session` (`id`, `user_id`, `name`, `rasse`, `typ`, `herkunft`, `glaube`, `attributes`, `derived_values`, `skills`, `spells`, `skill_points`, `created_at`, `updated_at`, `expires_at`, `current_step`, `geschlecht`, `stand`) VALUES
('char_create_1_1756326371', 1, 'wergw5 z ', 'Mensch', 'Priester Streiter', 'Aran', '', '{\"st\":89,\"gs\":64,\"gw\":77,\"ko\":71,\"in\":87,\"zt\":44,\"au\":87}', '{\"pa\":33,\"wk\":27,\"lp_max\":11,\"ap_max\":14,\"b_max\":26,\"resistenz_koerper\":11,\"resistenz_geist\":11,\"resistenz_bonus_koerper\":0,\"resistenz_bonus_geist\":0,\"abwehr\":11,\"abwehr_bonus\":0,\"ausdauer_bonus\":11,\"angriffs_bonus\":0,\"zaubern\":11,\"zauber_bonus\":0,\"raufen\":8,\"schadens_bonus\":3,\"sg\":0,\"gg\":0,\"gp\":0}', '[{\"name\":\"Klettern\",\"level\":0,\"category\":\"Alltag\",\"cost\":1},{\"name\":\"Reiten\",\"level\":0,\"category\":\"Alltag\",\"cost\":1},{\"name\":\"Sprache\",\"level\":0,\"category\":\"Alltag\",\"cost\":1},{\"name\":\"Athletik\",\"level\":0,\"category\":\"Kampf\",\"cost\":2},{\"name\":\"Spießwaffen\",\"level\":0,\"category\":\"Waffen\",\"cost\":2},{\"name\":\"Stielwurfwaffen\",\"level\":0,\"category\":\"Waffen\",\"cost\":2},{\"name\":\"Waffenloser Kampf\",\"level\":0,\"category\":\"Waffen\",\"cost\":2},{\"name\":\"Stichwaffen\",\"level\":0,\"category\":\"Waffen\",\"cost\":2},{\"name\":\"Heilkunde\",\"level\":0,\"category\":\"Wissen\",\"cost\":2},{\"name\":\"Naturkunde\",\"level\":0,\"category\":\"Wissen\",\"cost\":2}]', '[{\"name\":\"Göttlicher Schutz v. d. Bösen\",\"cost\":1},{\"name\":\"Erkennen der Aura\",\"cost\":1},{\"name\":\"Heiliger Zorn\",\"cost\":1},{\"name\":\"Blutmeisterschaft\",\"cost\":1}]', '{}', '2025-08-27 20:26:11.255', '2025-08-27 21:36:58.718', '2025-09-10 20:26:11.255', 5, 'Männlich', 'Mittelschicht');

--
-- Indizes der exportierten Tabellen
--

--
-- Indizes für die Tabelle `char_char_creation_session`
--
ALTER TABLE `char_char_creation_session`
  ADD PRIMARY KEY (`id`),
  ADD KEY `idx_char_char_creation_session_user_id` (`user_id`);
COMMIT;

/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
