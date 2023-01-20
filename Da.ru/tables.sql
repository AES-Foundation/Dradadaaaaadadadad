-- --------------------------------------------------------
-- Хост:                         176.124.212.74
-- Версия сервера:               10.1.48-MariaDB-0ubuntu0.18.04.1 - Ubuntu 18.04
-- Операционная система:         debian-linux-gnu
-- HeidiSQL Версия:              12.3.0.6589
-- --------------------------------------------------------

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET NAMES utf8 */;
/*!50503 SET NAMES utf8mb4 */;
/*!40103 SET @OLD_TIME_ZONE=@@TIME_ZONE */;
/*!40103 SET TIME_ZONE='+00:00' */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;

-- Дамп структуры для таблица laefye_test.heads
CREATE TABLE IF NOT EXISTS `heads` (
  `id` varchar(36) NOT NULL DEFAULT '',
  `user` varchar(36) DEFAULT NULL,
  `skin` varchar(255) DEFAULT NULL,
  `head` varchar(255) DEFAULT NULL,
  `params` varchar(255) DEFAULT NULL,
  `createdAt` bigint(20) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- Экспортируемые данные не выделены.

-- Дамп структуры для таблица laefye_test.payments
CREATE TABLE IF NOT EXISTS `payments` (
  `id` varchar(36) NOT NULL DEFAULT '',
  `user` varchar(36) DEFAULT NULL,
  `paysystem` varchar(50) DEFAULT NULL,
  `isPaid` tinyint(4) DEFAULT NULL,
  `amount` double DEFAULT NULL,
  `createdAt` bigint(20) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- Экспортируемые данные не выделены.

-- Дамп структуры для таблица laefye_test.sessions
CREATE TABLE IF NOT EXISTS `sessions` (
  `id` varchar(36) NOT NULL DEFAULT '',
  `key` varchar(64) NOT NULL DEFAULT '',
  `user` varchar(36) NOT NULL DEFAULT '',
  `expires` bigint(20) NOT NULL DEFAULT '0',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- Экспортируемые данные не выделены.

-- Дамп структуры для таблица laefye_test.users
CREATE TABLE IF NOT EXISTS `users` (
  `id` varchar(36) NOT NULL DEFAULT '',
  `email` varchar(64) NOT NULL DEFAULT '',
  `password` varchar(64) NOT NULL DEFAULT '',
  `displayName` varchar(64) NOT NULL DEFAULT '',
  `isVerified` tinyint(4) NOT NULL DEFAULT '0',
  `isHiddenInTopDonaters` tinyint(4) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- Экспортируемые данные не выделены.

/*!40103 SET TIME_ZONE=IFNULL(@OLD_TIME_ZONE, 'system') */;
/*!40101 SET SQL_MODE=IFNULL(@OLD_SQL_MODE, '') */;
/*!40014 SET FOREIGN_KEY_CHECKS=IFNULL(@OLD_FOREIGN_KEY_CHECKS, 1) */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40111 SET SQL_NOTES=IFNULL(@OLD_SQL_NOTES, 1) */;
