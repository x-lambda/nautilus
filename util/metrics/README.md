# metrics

1. 安装`prometheus`
```shell
docker run --name prometheus -d -p 9091:9090 quay.io/prometheus/prometheus
```

2. 安装`grafana`
```shell
docker run -d -p 3000:3000 grafana/grafana
```

3. 进入`prometheus`容器修改配置
```shell
vi /etc/prometheus/prometheus.yml
```

4. 添加`target`

![修改配置](https://github.com/x-lambda/note/blob/master/golang/imgs/prometheus_yaml.png)
注意: target中填当前服务的ip

5. 重启`prometheus`
```shell
docker restart prometheus
```

6. 查看

浏览器访问 `http://localhost:9091/classic/targets#job-prometheus`
![target](https://github.com/x-lambda/note/blob/master/golang/imgs/prometheus_unstart.png)   
   
7. 启动应用服务

浏览器访问 `http://localhost:9091/classic/targets#job-prometheus`
![target](https://github.com/x-lambda/note/blob/master/golang/imgs/prometheus_start.png)
可以看到新加的target已经up了

8. 配置`grafana`
![grafana](https://github.com/x-lambda/note/blob/master/golang/imgs/prometheus_show.png)