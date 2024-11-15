CREATE TABLE IF NOT EXISTS `tweet` (
  `ID` INT(11) NOT NULL AUTO_INCREMENT,
  `Content` VARCHAR(255) NOT NULL ,
  `CreatedAt` DATETIME NOT NULL DEFAULT current_timestamp(),
  `UpdatedAt` DATETIME NOT NULL DEFAULT current_timestamp() ON UPDATE current_timestamp(),
  `DeletedAt` DATETIME,
  `UserID` INT(11),
  PRIMARY KEY (`ID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;
