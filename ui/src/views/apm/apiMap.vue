<template>
  <div>
    <span style="float:right;margin-right:20px">
        <Tag style="width:30px;height:15px;border:none" :style="{background:tagColor(0)}"></Tag>
        <span class="font-size-10">来源应用</span> 
         <Tag style="width:30px;height:15px;border:none" class="margin-left-10" :style="{background:tagColor(3)}"></Tag>
        <span class="font-size-10">当前应用接口</span> 
    </span>

    <div :id="id" class="app-service-map" style="width:calc(100vw - 180px);height:calc(100vh - 100px)"></div>
  </div>

  
</template>

<script>
import request from '@/utils/request' 
import echarts from "echarts";
export default {
  data() {
    return {
      id : 'app-service-map-id',
      chart: null,
      lineLength: 200,
      primaryNodeSize: 50,
      smallNodeSize: 35,
      lineLabelSize: 13,
      repulsion: 500
    };
  },
  watch: {
    "$store.state.apm.selDate"() {
      this.initServiceMap()
    },
    "$store.state.apm.appName"() {
      this.initServiceMap()
    }
  },
  mounted() {
    this.initServiceMap()
  },
  beforeDestroy() {
    this.destroyChart()
  },
  methods: {
    destroyChart() {
       if (this.chart) {
        this.chart.dispose();
        this.chart = null;
      }
    },
    tagColor(tp) {
      switch (tp) {
        case 0: //普通应用
           return "linear-gradient(to right, #01acca , #5adbe7)"
          break;
        case 1: // 数据库中间件
          return  "linear-gradient(to right, #ffb402 , #ffdc84)"
        case 3 : // 当前应用
           return  "linear-gradient(to right, #157eff , #35c2ff)"
        default:
          break;
      }
     
      
    },
    initServiceMap() {
      this.destroyChart()
      this.$Loading.start();
        request({ 
            url: '/web/apiMap',
            method: 'GET',
            params: {
              app_name: this.$store.state.apm.appName,
              start: JSON.parse(this.$store.state.apm.selDate)[0],
              end: JSON.parse(this.$store.state.apm.selDate)[1],
            }
        }).then(res => {
          this.$Loading.finish();
           if (res.data.data.nodes.length == 0) {
              this.$Message.info({
                content: '没有查询到数据',
                duration: 3 
              })
          } else {
            console.log(res.data.data)
            this.initChart(res.data.data.nodes,res.data.data.links)
          }
        }).catch(error => {
            this.$Loading.error();
          })
      },
    formatLinkLabel(link) {
      var error = 0
      if (link.access_count > 0) {
        error = (link.error_count / link.access_count) * 100
        error = error.toFixed(1);
      }
      return link.access_count + '/' + error + '%/' + link.avg + 'ms'
    },
    calcSize(nodes) {
        var l = nodes.length 
        if (l < 20) {
            this.lineLength = 200
            this.primaryNodeSize = 50
            this.smallNodeSize = 35
            this.lineLabelSize = 13
            this.repulsion = 500
            return
        }

        if (l < 40) {
           this.lineLength = 200
            this.primaryNodeSize = 40
            this.smallNodeSize = 20
            this.lineLabelSize = 12
            this.repulsion = 5500
            return
        }

        if (l < 80) {
            this.lineLength = 200
            this.primaryNodeSize = 30
            this.smallNodeSize = 15
            this.lineLabelSize = 10
            this.repulsion = 500
            return
        }

         this.lineLength = 100
        this.primaryNodeSize = 30
        this.smallNodeSize = 10
        this.lineLabelSize = 9
        this.repulsion = 500
    },
    initChart(nodes,links) {
      console.log(nodes,links)
      this.chart = echarts.init(document.getElementById(this.id));
      for (var j = 0; j < nodes.length; j++) {
        // 设置node的样式
        nodes[j].symbolSize = this.smallNodeSize;
        nodes[j].itemStyle = {
          normal: {
            color: new echarts.graphic.LinearGradient(0, 0, 1, 0, [{
                offset: 0,
                color: '#01acca'
            }, {
                offset: 1,
                color: '#5adbe7'
            }]),
          }
        };
        // node错误率超过一个值，则添加特殊显示
        if (nodes[j].span_count > 0) {
          if ((nodes[j].error_count / nodes[j].span_count) > 0.2){
            nodes[j].label = {
              normal: {
                color: "#ff0000"
              }
            };
          }
        }

        // 当前应用接口
        if (nodes[j].category == 2) {
          nodes[j].itemStyle = {
                normal: {
                   color:   new echarts.graphic.LinearGradient(0, 0, 1, 0, [{
                    offset: 0,
                    color: '#157eff'
                }, {
                    offset: 1,
                    color: '#35c2ff'
                }]),
                }
              }; 
        } 
      }


      this.calcSize(nodes)
      for (var i = 0; i < links.length; i++) {
        var color = "#12b5d0";
        if (links[i].access_count > 0) {
          if ((links[i].error_count / links[i].access_count) > 0.2 ) {
            color = "#ff0000";
          }
        }

        links[i].label = {
          normal: {
            show: true,
            formatter: this.formatLinkLabel(links[i]),
            fontSize: this.lineLabelSize
          }
        };

        links[i].lineStyle = {
          normal: {
            color: color,
            width: 1,
            curveness: 0
          }
        };
      }
      var option = {
        series: [
          {
            type: "graph",
            layout: "force",
            force: {
              repulsion: 2000,
              edgeLength: this.lineLength
            //   layoutAnimation: false
            },
            name: "应用",
            roam: true,
            draggable: true,
            focusNodeAdjacency: true,
            symbolSize: 20,
            label: {
              normal: {
                show: true,
                position: "bottom",
                color: "#12b5d0"
                // fontSize:10
              }
            },
            edgeSymbol: ["none", "arrow"],
            lineStyle: {
              normal: {
                width: 1,
                shadowColor: "none"
              }
            },
            edgeSymbolSize: 8,
            data: nodes,
            links: links,
            itemStyle: {
              normal: {
                label: {
                  show: true,
                  formatter: function(item) {
                    return item.data.name;
                  }
                }
              }
            }
          }
        ]
      };
      this.chart.setOption(option);


      function nodeOnClick(params) {
        
      }
      this.chart.on("click", nodeOnClick);
      //'click'、'dblclick'、'mousedown'、'mousemove'、'mouseup'、'mouseover'、'mouseout'
    }
  }
};
</script>