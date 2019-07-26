## fin-essay

有规律的抓取并展示各网站的金融咨询，追踪热点，提高阅读效率。

## 数据源
- [X]亿欧
获取最近7天数据。
- [X]网贷之家
获取最近7天数据。
- [X]虎嗅
获取最近7天的数据
- []36氪
当前已经没有了金融归类，需要通过其他的方式处理
- [X]微信公众号（金融行业网，互联网金融，独角金融，新流财经，也谈钱，未央网weuyangx，第一财经YiMagazine）
当前支持`独角金融`，`程序员的金融圈`，`互联网金融`，`也谈钱`数据获取，获取的数据为最近10次发表的文章。

## 项目部署
- [X] 直接运行
```sh
go run main.go
```
- [ ] docker
- [ ] k8s

## 后台任务设计

当前使用`chromedp`模拟用户抓取网页的时候，总是会出现超时的情况；另外也不能每次访问接口，都去第三方抓取网页，所以需要设计本地存储的方式。
如何判断是否是同一篇文章：**原则上使用文章名称+文章链接的哈希的方式判断是否是同一篇文章，当前直接使用文章的名称来判断**。

文件中内容格式：名称 链接 作者 时间 哈希码（当前默认xxxx），单个数据源存储在单个文件中，每行一条信息。

考虑到前期数据量较少，对于修改文件的情况采取了比较粗暴的方法，直接将全部内容取出，修改完成之后重新写入。

## 涉及技术

- [gin](https://github.com/gin-gonic/gin)
- [goquery](https://github.com/PuerkitoBio/goquery)
- [gjson](https://github.com/tidwall/gjson)
- [chromewd](https://github.com/chromedp/chromedp)

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

## 问题

代码需持续的优化，加强健壮性。

## 公众号
关注公众号`“程序员的金融圈”`，加入`探讨技术、金融、赚钱的小圈子`。

![](https://user-gold-cdn.xitu.io/2019/6/9/16b39674126fc0f0?imageView2/0/w/1280/h/960/format/webp/ignore-error/1)
