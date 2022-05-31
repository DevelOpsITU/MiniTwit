# Introduction to the bug (See issue #109)

* We had to migrate our database from SQLite to Postgres, as we thought it might be a bottleneck for our application.

* It took a long time and was difficult, which included researching a solution and any tools that might be of help.

* Things appeared to run smooth after the migration, but then after a night of sleep we woke up to a Devopsers worst nightmare.

# Incident response

* When we looked at the database, multiple follower entries in the DB was zero. This had happend as a consequence of our initial downtime (we missed a bunch of user registrations). Basically, we came to the conclusion that the simulater initially registered all the users used in the simulation, so we would never be able to get those users registered.

# The hack

* Since we missed these initial users we came up with a brilliant plan to get them back. We decided to introduce a **Dirrty hack**, where we would register users that tried to follow, and thus repopulating the database with the lost users. After some time, the group decided that they wanted to know how often the hack was used. So the group implemented a counter in the code and in the Grafana dashboard, so a counter would go each time the hack was used.


