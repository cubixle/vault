# Vault
To build the vault you'll need to use the makefile. From the root of this project run ```make build```.
To run the app ```make start```
To stop the app ```make stop```


### Routes

POST "/"
```
{
	"data": "Test Data",
	"ttl": 30
}
```

POST "/decrypt"
```
{
  "vault": "3HXCBDkcT7R2ub39FXluykb_SZmC2udY09R-F1UuEnnTaekT60T6LbhUf_llovadxRg3w1ZL_krFsyoHodpvqeNpuXGsdMoHAVxJMhZpBOzH",
  "key": "efc927b61eab7f9cc3bd274b40330956"
}
```
