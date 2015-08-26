
-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied

CREATE TABLE `sctReceipts` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `sctVersion` tinyint(3) NOT NULL,
  `logID` varchar(255) NOT NULL,
  `timestamp` bigint(20) NOT NULL,
  `extensions` mediumblob,
  `signature` mediumblob,
  `certificateSerial` varchar(255) NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back

DROP TABLE `sctReceipts`;
