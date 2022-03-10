# Migration from SQLite to Postgres

This document will describe the strategy for how to migrate from the running SQLite database to a new Postgres database.

It will also link to the issues/ discussions that have occurred in the process.
- [Mayor migration problem #109](https://github.com/DevelOpsITU/MiniTwit/issues/109)
- [Migrate data from SQLite to PostgreSQL #92](https://github.com/DevelOpsITU/MiniTwit/issues/92)
- [Handle Database migration/setup #37 ](https://github.com/DevelOpsITU/MiniTwit/issues/37)
- [Add repository pattern for database access #35](https://github.com/DevelOpsITU/MiniTwit/issues/35)
- [Feature/#37 handle database migration #89](https://github.com/DevelOpsITU/MiniTwit/pull/89)

## Migration plan

In this section will the migration plan be described and there will be a couple of different docker containers and servers in play.

- Application server
- Database server
- Minitwit-main (groupddevops/minitwit-go:237df00) = The minitwit application currently running on the application server
- Minitwit-Postgres = The newest application that can talk with both SQLite and Postgres databases. It will be this instance that in the end will run on the application server.



1. Start postgres server 13.5 on the database server
2. Create the database "minitwit"
3. Connect the Minitwit-Postgres container to the database to create table with GORM.
4. Stop Minitwit-main on the application server
5. Make a local copy of the database as a backup.
6. Change the data in the pwhash column to be a new valid option, maybe taken from the test?
7. Use Pgloader to move the data from SQLite to Postgres
8. Run the sql qurry: `UPDATE minitwit.public."user"
   SET pw_hash='pbkdf2:sha256:50000' || '$' || 'CCWW6o8F' || '$' || 'c3f9294679a99b7da156d4f267be5dcee37afebba25e3c893c8dd67e78513cb9'
   where 1=1` Because the values was still messed up.
9. Start the Minitwit-Postgres with the connection string to the postgres database.
10. Verify that it works. 
11. If does not work then gather information of why and then start the old minitwit-main application again.


## Test migration

- [x] Locally deploy groupddevops/minitwit-go:237df00 with a clean database
- [x] Run simulation 0-1000 on the groupddevops/minitwit-go:237df00
- [x] Start the postgres 13.5 on the database server
- [x] Create the database "minitwit" on the database server.
- [x] Connect the local Minitwit-Postgres to the database server. (Create the tables)
- [x] Create a test user and note the pw_hash entry (pbkdf2:sha256:50000$CCWW6o8F$c3f9294679a99b7da156d4f267be5dcee37afebba25e3c893c8dd67e78513cb9) = a
- [x] Stop all services
- [x] Change the data in pw_hash column of the sqlite database. `./scripts/updateSQLitePw_hash.sh`
- [x] Use Pgloader to transfer the data to Postgres
- [x] run the SQL stamtement to update all columns again.
- [x] Test logging in to a user with the password created for the test user 
- [x] run simulation 1001-2000 and verify that it looks okay.


