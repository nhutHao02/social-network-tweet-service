CREATE TABLE IF NOT EXISTS `tweetcomment` (
  `ID` INT(11) NOT NULL AUTO_INCREMENT,
  `Description` VARCHAR(255) NOT NULL ,
  `UserID` INT(11) NOT NULL,
  `CreatedAt` DATETIME NOT NULL DEFAULT current_timestamp(),
  `UpdatedAt` DATETIME NOT NULL DEFAULT current_timestamp() ON UPDATE current_timestamp(),
  `DeletedAt` DATETIME,
  `TweetID` INT(11),
  PRIMARY KEY (`ID`),
  FOREIGN KEY (`TweetID`) REFERENCES `tweet`(`ID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;