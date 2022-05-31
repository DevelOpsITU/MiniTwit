# The dark times..

* We had to migrate our database from SQLite to Postgres 

* Took a long time and was difficult

# Hope ahead..

* Things appeared to run smooth after the migration

# We jinxed it..

* However, multiple follower entries in the DB was zero

* This happend as a consequence of our initial downtime (we missed a bunch of user registrations) 

* Basically, we came to the conclusion that the simulater initially registered all the users used in the simulation

# The hack

* Since we missed these initial users we came up with a brilliant plan to get them back... 

* We decided to introduce a **Dirrty hack**, where we would register users that tried to follow, and thus repopulating the database with the lost users



-----------------------------------------------------------
Initial release by Kaare, second version written by Malthe.
