## fin-essay
定时抓取并展示各网站的金融咨询。

## 数据源
虎嗅、亿欧、36氪、网贷之家、微信公众号（金融行业网，互联网金融，独角金融，新流财经，也谈钱，未央网weuyangx，第一财经YiMagazine）

## 展示方式

### 每日数据

### 今日最火

### 本周最火

## 项目运行及部署
- [X] 直接运行
```sh
go run main.go
```
- [ ] docker
- [ ] k8s

## 接口文档
```sh
# 获取网贷之家本周数据
curl -XGET http://localhost:8080/fin/wdzj
# 获取亿欧本周数据
curl -XGET http://localhost:8080/fin/iyiou
# 获取虎嗅本周数据
curl -XGET http://localhost:8080/fin/huxiu
# 获取微信公众号数据，当前支持独角金融、互联网金融、也谈钱、程序员的金融圈
curl -XGET http://localhost:8080/fin/wx/:wxname
```

## 公众号
关注公众号`“程序员的金融圈”`，加入`探讨技术、金融、赚钱的小圈子`。

![](https://user-gold-cdn.xitu.io/2019/6/9/16b39674126fc0f0?imageView2/0/w/1280/h/960/format/webp/ignore-error/1)
