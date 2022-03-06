#!/bin/bash
RED='\033[1;31m'
YELLOW='\033[1;33m'
GREEN='\033[1;32m'
NC='\033[0m'
# check if pgloader is installed
if ! command -v pgloader &> /dev/null 
then
    echo -e "${YELLOW}[WARNING] - unable to find pgloader - installing...${NC}"
    apt-get install --yes pgloader
else
    echo -e "${GREEN}[INFO] - pgloader found${NC}"
fi

# if the installation failed exit
if ! command -v pgloader &> /dev/null 
then
    echo -e "${RED}[ERROR] something went wrong in the installation of pgloader${NC}"
    exit
fi

# check if postgres is installed
if ! command -v psql &> /dev/null 
then
    echo -e "${YELLOW}[WARNING] - unable to find postgresql - installing...${NC}"
    apt-get install --yes postgresql-10 
else
    echo -e "${GREEN}[INFO] - postgresql found${NC}"
fi
# if the installation failed exit
if ! command -v psql &> /dev/null 
then
    echo -e "${RED}[ERROR] something went wrong in the installation of postgresql${NC}"
    exit
fi

POSTGRES_USER=postgres
POSTGRES_PASS=password
POSTGRES_DB_NAME=minitwit
# check if postgres database exists
# l = list all databases, q = remove header from table, t = turn into tuples
# cut = split on charater | (splits the table from psql), -f = select first item
# grep -w match entire word, -q supress output
if ! psql -U ${POSTGRES_USER} -lqt | cut -d \| -f 1 | grep -qw minitwit
then
    echo -e "${RED}[ERROR] postgres database does not exists${NC}"
    exit
else
    echo -e "${GREEN}[INFO] - postgres database found${NC}"
fi

#############################################################
#                                                           #
#                   migrate the database                    #
#                                                           #
#############################################################
echo -e "${GREEN}[INFO] - starting migration...${NC}"
# copy database to this server (does not seem that you can directly do it in the pgloader command)
scp -i ~/.ssh/id_rsa.pub root@thomsen-it.dk:/tmp/minitwit.db /tmp/minitwit.db
# migrate the data only from the database
pgloader -v --with 'data only' /tmp/minitwit.db postgresql:///${POSTGRES_DB_NAME}?user=${POSTGRES_USER}&password=${POSTGRES_PASS}
