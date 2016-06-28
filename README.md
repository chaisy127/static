# 静态集群
static是一个分布式文件系统，主要存储小文件。由Golang语言实现，存储是基于mongodb的GridFS。</br>

####功能列表：</br>
  支持基本的上传、下载和删除接口</br>
  支持bucket</br>
</br>
####接口详情：</br>
#####上传文件：</br>
    POST /api/v1/storage?fname=$fname&fid=$fid&bucketname=$bucketname</br>
#####请求参数：</br>
    fname       文件名</br>
    fid         文件id，是该文件的md5值</br>
    bucketname  bucket名</br>
    post body   文件的内容</br>
#####响应参数：</br>
    errno       错误码，成功时为10000</br>
    errmsg      错误描述信息，在errno非10000的情况下为非空字符串</br>
</br>
#####下载文件：</br>
    GET /api/v1/storage?fid=$fid&bucketname=$bucketname</br>
#####请求参数：</br>
    fid         文件id，是该文件的md5值</br>
    bucketname  bucket名</br>
#####响应参数：</br>
    errno       错误码，成功时为10000</br>
    errmsg      错误描述信息，在errno非10000的情况下为非空字符串</br>
    data        文件内容</br>
#####删除文件：</br>
    DELETE /api/v1/storage?fid=$fid</br>
#####请求参数：</br>
    fid         文件id，是该文件的md5值</br>
#####响应参数：</br>
    errno       错误码，成功时为10000</br>
    errmsg      错误描述信息，在errno非10000的情况下为非空字符串</br>
</br>
#####创建bucket：</br>
    POST /api/v1/bucket?bucketname=$bucketname</br>
#####请求参数：</br>
    bucketname  bucket名</br>
#####响应参数：</br>
    errno       错误码，成功时为10000</br>
    errmsg      错误描述信息，在errno非10000的情况下为非空字符串</br>
#####删除bucket：</br>
    DELETE /api/v1/bucket?bucketname=$bucketname</br>
#####请求参数：</br>
    bucketname  bucket名</br>
#####响应参数：</br>
    errno       错误码，成功时为10000</br>
    errmsg      错误描述信息，在errno非10000的情况下为非空字符串</br>
