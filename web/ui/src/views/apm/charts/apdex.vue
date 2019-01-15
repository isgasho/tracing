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
    },
     dateList: {
        type: Array,
        default: []
    },
    valueList: {
        type: Array,
        default: []
    }
  },
  data() {
    return {
      chart: null
    }
  },
  watch: {
    dateList(val) {
      this.initChart()
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
      this.chart = echarts.init(document.getElementById(this.id))
  var option = {
          title: {
              text: 'Apdex健康指标',
              textStyle: {
                fontWeight: 'normal',
                fontSize: 16
            },
              x: 'center',
          },
          tooltip: {
              trigger: 'axis',

              axisPointer: {
                  animation: false
              }
          },
          legend: {
              data: ['流量'],
              x: 'left'
          },
          // toolbox: {
          //     feature: {
          //         saveAsImage: {}
          //     }
          // },
          axisPointer: {
              link: {
                  xAxisIndex: 'all'
              }
          },
          grid: [{
              left: 40,
              right: 40,
          }, {
              left: 40,
              right: 40,
          }],
          xAxis: [{


              type: 'category',
              boundaryGap: false,
              axisLine: {
                  onZero: true
              },
              data: this.dateList
          }, {
              gridIndex: 1
          }],

          yAxis: [{

              type: 'value',
              max: 1,
              min: 0,
              interval: 0.2,


          }, {
              gridIndex: 1
          }],
          series: [{
              name: '数值',
              type: 'line',
              smooth: true,
              symbol: 'circle',
              symbolSize: 9,
              showSymbol: false,
              lineStyle: {
                  normal: {
                      width: 1.5
                  }
              },
              itemStyle: {
                  normal: {
                      color: '#847ecc'
                  }
              },
              markPoint: {
                  data: [{
                      type: 'max',
                      name: '最大值'
                  }, {
                      type: 'min',
                      name: '最小值'
                  }]
              },
              markArea: {
                  silent: true,
                  label: {
                      normal: {
                          position: ['10%', '50%']
                      }
                  },
                  data: [
                      [{
                          name: '优',
                          yAxis: 1,
                          itemStyle: {
                              normal: {
                                  color: 'rgba(183,234,209,0.7)'
                              }
                          },
                      }, {
                          yAxis: 0.8
                      }],
                      [{
                          name: '良',
                          yAxis: 0.8,
                          itemStyle: {
                              normal: {
                                  color: 'rgba(175,214,254,0.7)'
                              }
                          },
                      }, {
                          yAxis: 0.6,
                      }],
                      [{
                          name: '差',
                          yAxis: 0.6,
                          itemStyle: {
                              normal: {
                                  color: 'rgba(244,228,199,0.7)'
                              }
                          }
                      }, {
                          yAxis: 0,
                      }]
                  ]
              },
              data: this.valueList

          }]
      };
      this.chart.setOption(option)
    }
  }
}
</script>
