<template>
  <span :class="className" :id="id"></span>
</template>

<script>
import echarts from "echarts";
export default {
  props: {
    className: {
      type: String,
      default: "chart"
    },
    id: {
      type: String,
      default: "chart"
    }
  },
  data() {
    return {
      chart: null,
      nodes : [
        {
          name: "H1",
          span_count: 50,
          error_count: 30
        },
        {
          name: "mysql",
          span_count: 0,
          error_count: 0
        },
        {
          name: "redis",
          span_count: 0,
          error_count: 0
        },
        {
          name: "外部运营商",
          span_count: 0,
          error_count: 0
        },
        {
          name: "party",
          span_count: 200,
          error_count: 30
        },
        {
          name: "payGateway",
          span_count: 150,
          error_count: 20
        }
      ],
    links : [
        {
          source: "H1",
          target: "mysql",
          access_count: 30,
          error_count: 3,
          avg: 10
        },
        {
          source: "H1",
          target: "redis",
           access_count: 30,
          error_count: 3,
          avg: 2
        },
        {
          source: "H1",
          target: "外部运营商",
          status: "30/10%/300ms",
          access_count: 30,
          error_count: 15,
          avg: 1000
        },
        {
          source: "party",
          target: "H1",
          access_count: 20,
          error_count: 5,
          avg: 10
        },
        {
          source: 'party',
          target: 'mysql',
          access_count: 20,
          error_count: 1,
          avg: 25
        },
        {
          source: "payGateway",
          target: "party",
          access_count: 80,
          error_count: 30,
          avg: 30
        }
    ],
    lineLength: 200,
    primaryNodeSize: 50,
    smallNodeSize: 35,
    lineLabelSize: 13,
    repulsion: 500
    };
  },
 
  mounted() {
    this.initChart();
  },
  beforeDestroy() {
    if (!this.chart) {
      return;
    }
    this.chart.dispose();
    this.chart = null;
  },
  methods: {
    formatLinkLabel(link) {
      var error = 0
      if (link.access_count > 0) {
        error = (link.error_count / link.access_count) * 100
      }
      return link.access_count + '/' + error + '%/' + link.avg + 'ms'
    },
    calcSize() {
        var l = this.nodes.length 
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
    initChart() {
      this.chart = echarts.init(document.getElementById(this.id));
      for (var j = 0; j < this.nodes.length; j++) {
        this.nodes[j].symbolSize = this.smallNodeSize;
        this.nodes[j].itemStyle = {
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
      }
      this.calcSize()
      for (var i = 0; i < this.links.length; i++) {
        var color = "#12b5d0";
        if (this.links[i].access_count > 0) {
          if ((this.links[i].error_count / this.links[i].access_count) > 0.2 ) {
            color = "#ff0000";
          }
        }

        this.links[i].label = {
          normal: {
            show: true,
            formatter: this.formatLinkLabel(this.links[i]),
            fontSize: this.lineLabelSize
          }
        };

        this.links[i].lineStyle = {
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
              edgeLength: this.lineLength,
              layoutAnimation: false
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
            data: this.nodes,
            links: this.links,
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
      // 用于告警的动态效果
      setTimeout(() => {
        var dataI = [];
        for (var n = 0; n < this.nodes.length; n++) {
          // 节点错误率超过一个值，则添加特殊显示
          if (this.nodes[n].span_count > 0) {
            if ((this.nodes[n].error_count / this.nodes[n].span_count) > 0.2){
              option.series[0].data[n].label = {
                normal: {
                  color: "#ff0000"
                }
              };
              dataI.push(n);
            }
          }
        
          if (this.nodes[n].name == this.$store.state.apm.appName) {
            option.series[0].data[n].symbolSize = this.primaryNodeSize;
            option.series[0].data[n].itemStyle = {
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
        this.chart.setOption(option);
      }, 500);

      this.chart.setOption(option);

      this.chart.group = this.group;

      function nodeOnClick(params) {
        
      }
      this.chart.on("click", nodeOnClick);
      //'click'、'dblclick'、'mousedown'、'mousemove'、'mouseup'、'mouseover'、'mouseout'
    }
  }
};
</script>