#!/bin/bash

sqlite3 scan.db <<END_SQL
SELECT  * FROM obmenkadata;
SELECT count(*) FROM obmenkadata;
END_SQL

ls -al scan.db
