# yunfs

**Yunfs** is a cloud based FUSE file system.

Features
===========
* It's designed for achieving large scale of files.
* Everything on cloud, one can get his file system as long as he can access the cloud.
* Supporting multiple cloud service provides (Aliyun OSS, Amazon S3, ...).
* All the files are encrypted.
* Advanced cache system.
* Support Linux & OSX.

Dependency
==========
* github.com/hanwen/go-fuse/fuse
* github.com/aliyun/aliyun-oss-go-sdk


Usage
=======
```
git clone https://github.com/Jiajun-Fan/Yunfs
go build github.com/Jiajun-Fan/yunfs
./yunfs # for the first time running yunfs, it creates a configuration template.
vim ~/.yunfs/yunfs.json # edit configuration
./yunfs
```

configuration
===============
```
{
    "oss": {
        "type": "",     # cloud storage type ("aliyun", "s3", ...)
        "key": "",      # access key
        "secret": "",   # access secret
        "bucket": "",   # bucket name
        "end_point": "" # cloud storage end point URL
    },
    "encrypt": {
        "type": "",     # encrypt type ("aes" ...)
        "key": ""       # encrypt key
    },
    "file_system": {
        "block_size": 1024, # entry number in each meta file
        "cache_size": 0,    # max entry number in local cache
        "meta_prefix": "",  # prefix for storing meta file
        "mount_point": ""   # mount point
    }
}
```

Milestones
===========
| Milestones | Description | Status |
|:----------:| ----------- |:------:|
|1|Read only file system with encryption & decryption|Done|
|2|Support Aliyun OSS|Done|
|3|Support writable file system|In progress|
|4|Local cache system|TODO|
|5|Support more cloud service provides|TODO|
|6|Performance optimazition|TODO|

