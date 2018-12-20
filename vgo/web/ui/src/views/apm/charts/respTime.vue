<template>
  <div :class="className" :id="id" :style="{height:height,width:width}"></div>
</template>

<script>
import echarts from 'echarts'
 
export default {
  props: {
    className: {
      type: String,
      default: 'chart'
    },
    id: {
      type: String,
      default: 'chart'
    },
    width: {
      type: String,
      default: '200px'
    },
    height: {
      type: String,
      default: '200px'
    }
  },
  data() {
    return {
      chart: null
    }
  },
  mounted() {
    this.initChart()
  },
  beforeDestroy() {
    if (!this.chart) {
      return
    }
    this.chart.dispose()
    this.chart = null
  },
  methods: {
    initChart() {
        var dateList = ["2018-12-06 12:45:00","2018-12-06 12:47:00","2018-12-06 12:49:00","2018-12-06 12:51:00",
      "2018-12-06 12:53:00","2018-12-06 12:55:00","2018-12-06 12:57:00","2018-12-06 12:59:00","2018-12-06 13:01:00",
      "2018-12-06 13:03:00","2018-12-06 13:05:00","2018-12-06 13:07:00","2018-12-06 13:09:00","2018-12-06 13:11:00",
      "2018-12-06 13:13:00","2018-12-06 13:15:00","2018-12-06 13:17:00","2018-12-06 13:19:00","2018-12-06 13:21:00",
      "2018-12-06 13:23:00","2018-12-06 13:25:00","2018-12-06 13:27:00","2018-12-06 13:29:00","2018-12-06 13:31:00",
      "2018-12-06 13:33:00","2018-12-06 13:35:00","2018-12-06 13:37:00","2018-12-06 13:39:00","2018-12-06 13:41:00",
      "2018-12-06 13:43:00"]
      var valueList = [356,355,341,373,349,362,342,363,298,332,325,353,334,355,345,268,378,336,424,371,364,364,347,283,
      302,306,349,316,358,360]

      this.chart = echarts.init(document.getElementById(this.id))
      var option = {
        backgroundColor: '#fff',
        title: {
            text: '应用响应时间',
            textStyle: {
                fontWeight: 'normal',
                fontSize: 16
            },
            left: 'center'
        },
        tooltip: {
            trigger: 'axis', 
            axisPointer: {
            }
        },
        grid: {
            left: '4%',
            right: '2%',
            bottom: '8%',
            top:'14%',
            containLabel: true
        },
        xAxis: [{
            type: 'category',
            boundaryGap: false,
            axisLine: {
            },
            data: dateList
        }],
        yAxis: [{
            type: 'value',
             name: '单位（毫秒）',
            axisTick: {
                show: false
            },
            // axisLine: {
            //     lineStyle: {
            //         color: '#57617B'
            //     }
            // },
            // axisLabel: {
            //     textStyle: {
            //         fontSize: 12
            //     }
            // },
            splitLine: {
                show: false
            }
        }],
        series: [{
            name: '',
            type: 'line',
            smooth: true,
            lineStyle: {
                normal: {
                    width: 2
                }
            },
            areaStyle: {
                normal: {
                    color: new echarts.graphic.LinearGradient(0, 0, 0, 1, [{
                        offset: 0,
                        color: 'rgba(82, 191, 255, 0.3)'
                    }, {
                        offset: 0.8,
                        color: 'rgba(82, 191, 255, 0)'
                    }], false),
                    shadowColor: 'rgba(228, 139, 76, 0.1)',
                    shadowBlur: 10
                }
            },
            symbolSize:4,  
            itemStyle: {
                normal: {
                    color: 'rgb(82, 191, 255)',
                    borderColor:'#e48b4c'
                },
            },
            data: valueList,
        } ]
    };
      this.chart.setOption(option)
    }
  }
}
</script>
