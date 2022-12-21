# quick-cache

## Description 
This is a tool is a lightweight caching tool in golang. It is not met for caching large amounts of data but more for 
reducing how often you need to call a database. If you are looking for large amount of data storage
try to take a look at [rockdb](http://rocksdb.org/)

## How to use

### Setup Timelimit
```golang
//TimeLimit can take up to 3 input for how long you want something to stay in cache
//Seconds in int
//Minutes in int
//Hours in int
time := cache.TimeLimit{
Seconds: 0,
Minutes: 1,
Hours: 2,
}
```

### Initilize cache
This must be done before using cache
This example show if you want to run an auto cleanup to free up memory on schedule
```golang 
cache.InitCache(&cache.TimeLimit{Minutes: 30},cache.TimeLimit{Minutes: 30})
```

This shows how to not automcatic run cleanups but it will still cleanup if the item is found and the key has expired
```golang
cache.InitCache(nil,cache.TimeLimit{Minutes: 30})
```

### Get
```golang
val, ok := cache.Get(key)
if ok {
  //You will choose your what interface you wish to return.      
  targetFoo = val.(targetFoo.type)
}
else{
//Grab your data
}
```

### Add
```golang
//key is a string
//value is whatever you wish it to be
cache.Add(key, value)
```

