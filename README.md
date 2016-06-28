# 静态集群
static是一个分布式文件系统，主要存储小文件。由Golang语言实现，存储是基于mongodb的GridFS。

功能列表：
  支持基本的上传、下载和删除接口
  支持bucket
  
接口详情：
  上传文件：
    POST /api/v1/storage?fname=$fname&fid=$fid&bucketname=$bucketname
  请求参数：
    fname       文件名
    fid         文件id，是该文件的md5值
    bucketname  bucket名
    post body   文件的内容
  响应参数：
    errno       错误码，成功时为10000
    errmsg      错误描述信息，在errno非10000的情况下为非空字符串
    
  下载文件：
    GET /api/v1/storage?fid=$fid&bucketname=$bucketname
  请求参数：
    fid         文件id，是该文件的md5值
    bucketname  bucket名
  响应参数：
    errno       错误码，成功时为10000
    errmsg      错误描述信息，在errno非10000的情况下为非空字符串
    data        文件内容
  删除文件：
    DELETE /api/v1/storage?fid=$fid
  请求参数：
    fid         文件id，是该文件的md5值
  响应参数：
    errno       错误码，成功时为10000
    errmsg      错误描述信息，在errno非10000的情况下为非空字符串
    
  创建bucket：
    POST /api/v1/bucket?bucketname=$bucketname
  请求参数：
    bucketname  bucket名
  响应参数：
    errno       错误码，成功时为10000
    errmsg      错误描述信息，在errno非10000的情况下为非空字符串
  删除bucket：
    DELETE /api/v1/bucket?bucketname=$bucketname
  请求参数：
    bucketname  bucket名
  响应参数：
    errno       错误码，成功时为10000
    errmsg      错误描述信息，在errno非10000的情况下为非空字符串
