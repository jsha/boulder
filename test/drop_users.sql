-- Before setting up any privileges, we revoke existing ones to make sure we
-- start from a clean slate.
REVOKE ALL PRIVILEGES, GRANT OPTION FROM 'policy'@'localhost';
REVOKE ALL PRIVILEGES, GRANT OPTION FROM 'sa'@'localhost';
REVOKE ALL PRIVILEGES, GRANT OPTION FROM 'ocsp_resp'@'localhost';
REVOKE ALL PRIVILEGES, GRANT OPTION FROM 'ocsp_update'@'localhost';
REVOKE ALL PRIVILEGES, GRANT OPTION FROM 'revoker'@'localhost';
REVOKE ALL PRIVILEGES, GRANT OPTION FROM 'importer'@'localhost';
REVOKE ALL PRIVILEGES, GRANT OPTION FROM 'mailer'@'localhost';
REVOKE ALL PRIVILEGES, GRANT OPTION FROM 'cert_checker'@'localhost';

