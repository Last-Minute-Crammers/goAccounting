# 问题求解实战项目开发 -- goAccounting server端

### docker 运行注意

为了避免网络问题, 请尽量挂着clash的梯子进行

```bash
docker-compose up --build 
```

如果出现问题, 请看下你的clash verge代理端口是否是7897, 不一样的话就把`docker-compose.yml`里面的http_proxy等相关配置的端口改成clash verge的代理端口



然后访问 (详细见下面的接口文档)

```url
http://localhost:8080/api/${other routers}
```

即可

### 管理员相关的api

请见:[记账APP的接口文档](https://documenter.getpostman.com/view/43153095/2sB3WpS1bD)

