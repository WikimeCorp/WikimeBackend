# WikimeBackend
Backend servis for Wikime

## Environment variables and Config 
To successfully start the server, the following settings are required:

  ##### These parameters must be specified either in *Environment variables* or in *config file(.env)*:
  - `ADDR` is the address to start the server 
  - `PORT` is the port on which the server will start
  - `MONGO_URL` is the mongo database address
  - `DB_NAME` is the database name
 ##### These parameters are passed through *command line arguments*:
  - `--configPath` is path to congig file
  - `--addr` is similarly to ADDR, optional if the address from the config does not suit you
  - `--port` is similarly to PORT, optional if the port from the config does not suit you