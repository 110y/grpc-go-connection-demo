## Demo

### grpc-go v1.41.0

```
> make up
> make show-caller-channel

ID	Name        	State	Channel	SubChannel	Calls	Success	Fail	LastCall
1	callee:5001 	READY	0      	1         	0     	0     	0     	none

> docker compose stop callee
> make show-caller-channel

ID	Name        	State	Channel	SubChannel	Calls	Success	Fail	LastCall
1	callee:5001 	IDLE	0      	1         	0     	0     	0     	none

> docker compose up --detach callee
> make show-caller-channel

ID	Name        	State	Channel	SubChannel	Calls	Success	Fail	LastCall
1	callee:5001 	IDLE	0      	1         	0     	0     	0     	none
```

### grpc-go v1.40.0

-   need to execute `go get google.golang.org/grpc@v1.40.0`

```
> make up
> make show-caller-channel

ID	Name        	State	Channel	SubChannel	Calls	Success	Fail	LastCall
1	callee:5001 	READY	0      	1         	0     	0     	0     	none

> docker compose stop callee
> make show-caller-channel

ID	Name        	State	            Channel	SubChannel	Calls	Success	Fail	LastCall
1	callee:5001 	TRANSIENT_FAILURE	0      	1         	0     	0     	0     	none

> docker compose up --detach callee
> make show-caller-channel

ID	Name        	State	Channel	SubChannel	Calls	Success	Fail	LastCall
1	callee:5001 	READY	0      	1         	0     	0     	0     	none
```
