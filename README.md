## lima_beta_player

A player side program of the board game "Loot Of Lima" in virtual version. 

By the team Beta of CSCI630 Spring 2021 at CSU, Chico.

## Run program

```
go run beta_player.go
```
```
go build beta_player.go
./beta_player
```

# Parameters in command line

[-h]  : To know needed parameters

[-n=x]  : Player number for the game - options from 1 to 4 for x (default "2")

[-pn=name] : Prefixed name for named pipe - needed to be same as server's named pipe (default "all")

[-a=bool]  : Autopilot mode - boolean value (default "true")

ex.
go run
```
go run beta_player.go -h
```
```
go run beta_player.go -n=1 -pn=all -a=true
```

go build
```
go build beta_player.go
./beta_player -h
```
```
./beta_player -n=1 -pn=all -a=true
```