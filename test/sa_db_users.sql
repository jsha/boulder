--
-- Copyright 2015 ISRG.  All rights reserved
-- This Source Code Form is subject to the terms of the Mozilla Public
-- License, v. 2.0. If a copy of the MPL was not distributed with this
-- file, You can obtain one at http://mozilla.org/MPL/2.0/.
--
-- This file defines the default users for the primary database, used by
-- all the parts of Boulder except the Certificate Authority module, which
-- utilizes its own database.
--

-- Create users for each component with the appropriate permissions. We want to
-- drop each user and recreate them, but if the user doesn't already exist, the
-- drop command will fail. So we grant the dummy `USAGE` privilege to make sure
-- the user exists and then drop the user.

-- Storage Authority
GRANT USAGE ON *.* TO 'sa'@'localhost';
DROP USER 'sa'@'localhost';
CREATE USER 'sa'@'localhost';
GRANT SELECT,INSERT,UPDATE ON authz TO 'sa'@'localhost';
GRANT SELECT,INSERT,UPDATE,DELETE ON pendingAuthorizations TO 'sa'@'localhost';
GRANT SELECT,INSERT ON certificates TO 'sa'@'localhost';
GRANT SELECT,INSERT,UPDATE ON certificateStatus TO 'sa'@'localhost';
GRANT SELECT,INSERT ON issuedNames TO 'sa'@'localhost';
GRANT SELECT ON sctReceipts TO 'sa'@'localhost';
GRANT SELECT,INSERT ON deniedCSRs TO 'sa'@'localhost';
GRANT INSERT ON ocspResponses TO 'sa'@'localhost';
GRANT SELECT,INSERT,UPDATE ON registrations TO 'sa'@'localhost';
GRANT SELECT,INSERT,UPDATE ON challenges TO 'sa'@'localhost';

-- OCSP Responder
GRANT USAGE ON *.* TO 'ocsp_resp'@'localhost';
DROP USER 'ocsp_resp'@'localhost';
CREATE USER 'ocsp_resp'@'localhost';
GRANT SELECT ON ocspResponses TO 'ocsp_resp'@'localhost';

-- OCSP Generator Tool (Updater)
GRANT USAGE ON *.* TO 'ocsp_update'@'localhost';
DROP USER 'ocsp_update'@'localhost';
CREATE USER 'ocsp_update'@'localhost';
GRANT INSERT ON ocspResponses TO 'ocsp_update'@'localhost';
GRANT SELECT ON certificates TO 'ocsp_update'@'localhost';
GRANT SELECT,UPDATE ON certificateStatus TO 'ocsp_update'@'localhost';

-- Revoker Tool
GRANT USAGE ON *.* TO 'revoker'@'localhost';
DROP USER 'revoker'@'localhost';
CREATE USER 'revoker'@'localhost';
GRANT SELECT ON registrations TO 'revoker'@'localhost';
GRANT SELECT ON certificates TO 'revoker'@'localhost';
GRANT SELECT,INSERT ON deniedCSRs TO 'revoker'@'localhost';

-- External Cert Importer
GRANT USAGE ON *.* TO 'importer'@'localhost';
DROP USER 'importer'@'localhost';
CREATE USER 'importer'@'localhost';
GRANT SELECT,INSERT,UPDATE,DELETE ON identifierData TO 'importer'@'localhost';
GRANT SELECT,INSERT,UPDATE,DELETE ON externalCerts TO 'importer'@'localhost';

-- Expiration mailer
GRANT USAGE ON *.* TO 'mailer'@'localhost';
DROP USER 'mailer'@'localhost';
CREATE USER 'mailer'@'localhost';
GRANT SELECT ON certificates TO 'mailer'@'localhost';
GRANT SELECT,UPDATE ON certificateStatus TO 'mailer'@'localhost';

-- Cert checker
GRANT USAGE ON *.* TO 'cert_checker'@'localhost';
DROP USER 'cert_checker'@'localhost';
CREATE USER 'cert_checker'@'localhost';
GRANT SELECT ON certificates TO 'cert_checker'@'localhost';
