#!/bin/bash

mysqld

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

# make sure we can reach the mysqldb
# see http://tldp.org/LDP/abs/html/devref1.html for description of this syntax.
while ! exec 6<>/dev/tcp/boulder-mysql/3306; do
    echo "$(date) - still trying to connect to mysql at 0.0.0.0:3306"
    sleep 1 || exit
done

exec 6>&-
exec 6<&-

# create the database
source $DIR/create_db.sh
