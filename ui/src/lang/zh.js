export default {
   common: {
      test : '测试'
   },
   pageName: {introduce: '介绍','about': '我们能做什么',tech:'技术解密',deploy: '部署',install:'接入使用'},
   pages: {
about: `
### 我们能做什么
##### 先来看看监控金字塔
![](https://upload-images.jianshu.io/upload_images/8245841-85a846be3f1cd84d.png?imageMogr2/auto-orient/strip%7CimageView2/2/w/1240)

图有点复杂，简单来说就是：监控从来都不是孤立存在的，出现的问题也从来都不是孤立存在，因此我们需要一个统一的监控平台，能从各个层面把问题关联起来。

因此，OpenAPM2.0版本的规划之初，就把监控的全面性和关联性放在首要位置，举几个点简单说明下：
1. 以往虚拟机的监控指标CPU、内存都是在zabbix中的，用户需要在OpenAPM看应用监控，再去zabbix看虚拟机监控，麻烦不说，而且数据往往难以精确关联，
请问你怎么在一个月内，把这两个数据精确的关联到一起？
2. JVM的指标跟虚拟机指标能否关联？答案是可以
3. 用户的每次请求其实都是一条链路，从入口应用开始一直到请求结束A ->B ->C ->B ->A，在此过程中，存在监控采集到的数据(全链路数据)、用户打的业务日志等，
这些数据以往都是存在各个系统中的，切无法关联(只能想办法通过时间比对，但是很不准确，也没法做进一步数据分析)，在我们新的APM中，这些数据都将关联起来

监控分为几大类：机器监控、应用监控、数据库监控、中间件监控、业务监控，我们希望在未来，能把这些监控全部在数据层面打通，实现纵向(请求链路)、横向(数据关联)的全方位监控

**因此之前问题的答案已经呼之欲出**
- 为用户提供全方位的监控功能
- 问题发生时，用户能及时、精确的收到告警
- 查问题时，用户看到的数据不再是孤立存在的
- 用户可以从任何一个角度切入，进而在全局进行定位
- 提供丰富的数据统计和查询接口
`,


tech: `
### 技术难点
监控平台有一大特点，数据量非常非常大，每台虚拟机每秒上报的数据包在5 * QPS左右，同时虚拟机还要上报其他监控指标的数据包，因此对于监控平台而言，每秒处理的数据包量
轻轻松松就可以上几万。
因此我们需要在设计上重点考虑以下几点：
1. 高并发数据处理
2. 如何节省硬件资源(这个节省的成本是很大的)
3. 数据库的高并发存储和查询

下面我们一起来看看这些问题怎么解决

### 架构图
![](https://upload-images.jianshu.io/upload_images/8245841-d1b118a0065cdc59.png?imageMogr2/auto-orient/strip%7CimageView2/2/w/1240)

### 架构解析
1. 数据采集

   数据采集分为几个场景后端应用、移动APP以及浏览器，其中后端应用涉及多种语言，传化是以Java语言为主，因此我们对Java语言做了完善的支持。包括以下几点：
   - 引入pinpoint的采集端，通过字节码增强+插件的方式实现无侵入采集
   - 对log4j/log4j2/logback日志组件进行增强，将业务日志跟监控链路ID进行关联

   同时，为了更好的接入其他语言和客户端，我们未来将提供Java/Go/Javascript的采集SDK，实现真正的全面监控部署

2. 监控后台

   我们的监控后台全部采用Go语言开发，优点在于性能高、节省资源、天生高并发，整个后台有四大块组成：
   - Collector: 数据清洗，把外部数据清洗为内部的统一格式；数据统计，将原始数据计算为统计指标以及全链路数据，用于UI展示和告警；数据存储，定期将数据存储至Casaandra数据库
   - 告警服务：将Collector初步处理的数据按照告警设置，进行告警计算，最终通过短信、邮件等方式进行告警通知
   - 监控后台/UI: 提供页面展示，对统计数据和全链路数据进行处理和展示，只跟cassandra数据库交互

3. 数据库

   在OpenAPM1.0中，我们使用了hbase和mysql两大数据库，但是在其中发现了一些问题，例如
   - 资源占用过高
   - 性能不够好
   - 查询延迟较高
   - hbase很难维护

   在APM2.0设计阶段，我们调研了很多数据库，SQL方面包含Mysql、Tidb、Cockroach;No SQL方面包含Cassandra、hbase、mongodb、foundationdb、scyllaDB等，
   最后我们发现Cassandra和ScyllaDB特别适合，特别是这两个数据库都是基于同样的Cql语法，其中Cassandra是java实现，scylladb是c++实现，前者的优点在于发展更早，更稳定，功能更全，
   后者的优势在于性能非常高(吞吐是cassandra的10倍左右)，同时拥有极其稳定的延迟(得益于c++的无gc)，缺点是比较新，对cql的支持还不够完善。
   
   因此我们最终决定了先使用cassandra，以后再无缝替换为scylladb，就cassandra本身，性能也是hbase的3-6倍。

   在使用了新数据库后，吞吐提升了、延迟降低了，两个数据库统一为了一个，维护复杂度大大降低，同时硬件的成本也大幅降低。


### 如何解决技术难点
1. 高并发数据接收/处理
由于使用Golang语言，因此高并发和高性能都不是问题

2. 告警实时计算
最开始我们用的是spark，因为监控全部是滑动窗口计算，因此资源占用非常非常高，最终我们出于以下几点考虑，自己实现了实时计算部分：
   - 监控用的计算很纯粹，无需复杂的开源平台
   - 自己实现，可以降低机器的使用到之前的5分之一，并且延迟控制在1秒之内
   - 代码并不会很复杂，只要做好一致性hash，保证好服务的可用性即可

3. 数据库高并发存储和查询
之前也提到了，我们选型了cassandra，本身作为分布式nosql，性能就是很不错的，同时我们也做了以下优化：
   - 数据异步、批量存表
   - 为热数据建立反向索引表和影子表，提升查询性能
   - 优化服务器的存储参数和网络参数
   - 大量查询文献，优化服务器参数

4. 图表性能
对于监控来说，图表性能也是很重要的，因为加载数据都很多，因此我们做了优化：
   - 选型echarts作为图表组件
   - 在数据过滤层面限制数据的展示量，让用户更灵活的查询自己的数据
   - 数据分页展示
   - 采用了vuejs作为前端框架，性能大幅提升
`,
install: `
`




  }
}
