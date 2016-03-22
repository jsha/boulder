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
GRANT SELECT,INSERT,UPDATE ON authz TO 'sa'@'boulder2';
GRANT SELECT,INSERT,UPDATE,DELETE ON pendingAuthorizations TO 'sa'@'boulder2';
GRANT SELECT(id,Lockcol) ON pendingAuthorizations TO 'sa'@'boulder2';
GRANT SELECT,INSERT ON certificates TO 'sa'@'boulder2';
GRANT SELECT,INSERT,UPDATE ON certificateStatus TO 'sa'@'boulder2';
GRANT SELECT,INSERT ON issuedNames TO 'sa'@'boulder2';
GRANT SELECT,INSERT ON sctReceipts TO 'sa'@'boulder2';
GRANT SELECT,INSERT ON deniedCSRs TO 'sa'@'boulder2';
GRANT INSERT ON ocspResponses TO 'sa'@'boulder2';
GRANT SELECT,INSERT,UPDATE ON registrations TO 'sa'@'boulder2';
GRANT SELECT,INSERT,UPDATE ON challenges TO 'sa'@'boulder2';
GRANT SELECT,INSERT on fqdnSets TO 'sa'@'boulder2';

-- OCSP Responder
GRANT SELECT ON certificateStatus TO 'ocsp_resp'@'boulder2';
GRANT SELECT ON ocspResponses TO 'ocsp_resp'@'boulder2';

-- OCSP Generator Tool (Updater)
GRANT INSERT ON ocspResponses TO 'ocsp_update'@'boulder2';
GRANT SELECT ON certificates TO 'ocsp_update'@'boulder2';
GRANT SELECT,UPDATE ON certificateStatus TO 'ocsp_update'@'boulder2';
GRANT SELECT ON sctReceipts TO 'ocsp_update'@'boulder2';

-- Revoker Tool
GRANT SELECT ON registrations TO 'revoker'@'boulder2';
GRANT SELECT ON certificates TO 'revoker'@'boulder2';
GRANT SELECT,INSERT ON deniedCSRs TO 'revoker'@'boulder2';

-- External Cert Importer
GRANT SELECT,INSERT,UPDATE,DELETE ON identifierData TO 'importer'@'boulder2';
GRANT SELECT,INSERT,UPDATE,DELETE ON externalCerts TO 'importer'@'boulder2';

-- Expiration mailer
GRANT SELECT ON certificates TO 'mailer'@'boulder2';
GRANT SELECT,UPDATE ON certificateStatus TO 'mailer'@'boulder2';

-- Cert checker
GRANT SELECT ON certificates TO 'cert_checker'@'boulder2';

-- Name set table backfiller
GRANT SELECT ON certificates to 'backfiller'@'boulder2';
GRANT INSERT,SELECT ON fqdnSets to 'backfiller'@'boulder2';

-- Test setup and teardown
GRANT ALL PRIVILEGES ON * to 'test_setup'@'boulder2';
