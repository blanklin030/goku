+ run development
```
bee run -downdoc=true -gendoc=true
```
+ get api wiki
> visit http://localhost:8080/swagger/

+ run test
```
make test
```
> visit http://localhost:8888
+ run product
```
make release 
make push
# at product host enter the following command
GOKU_RUNMODE=prod docker-compose run -d
```