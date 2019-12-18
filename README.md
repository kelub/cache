# caches
缓存策略通用接口及其实现，包括常用的:
>* Cache Aside 旁路缓存模式
>* Read/Write Through 读穿/写穿模式
>* Write Back 写回策略



## Cache Aside

### read

read from cache
hit return date
miss read from DB, write to cache

### write

write to DB
delete cache

## Read/Write Through

### Read Through

read from cache
hit return date
else read from DB and write cache
return value

### Write Through
check isExist from cache

if exist Update cache and update DB
else 
write cache and write DB or write DB

## WriteBack

### Write
read cache
    hit write cache and  mark cache isDirty
    
    miss 
    if isDirty write DB
    read DB
    write cache
    mark cache isDirty

### Read
read cache
    hit return
if isDirty write DB
read cache
mark no dirty