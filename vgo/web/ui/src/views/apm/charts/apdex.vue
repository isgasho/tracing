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
      this.chart = echarts.init(document.getElementById(this.id))
        var timeData = ['0:00', '0:30',
      '1:00', '1:30',
      '2:00', '2:30',
      '3:00', '3:30',
      '4:00', '4:30',
      '5:00', '5:30',
      '6:00', '6:30',
      '7:00', '7:30',
      '8:00', '8:30',
      '9:00', '9:30',
      '10:00', '10:30',
      '11:00', '11:30',
      '12:00', '12:30',
      '13:00', '13:30',
      '14:00', '14:30',
      '15:00', '15:30',
      '16:00', '16:30',
      '17:00', '17:30',
      '18:00', '18:30',
      '19:00', '19:30',
      '20:00', '20:30',
      '21:00', '21:30',
      '22:00', '22:30',
      '23:00', '23:30'
  ];
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
              data: timeData
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
              data: [0.9, 0.7, 0.6, 0.34, 0.97, 0.99, 0.85, 0.65, 1, 1, 1, 222, 223, 312, 223, 222, 223, 222, 232, 262, 232, 232, 223, 222, 223, 332, 223, 232, 223, 322, 123, 222, 231, 322, 233, 122, 223, 232, 232, 222, 223, 232, 232, 222, 232, 132, 123, 212]

          }]
      };
      this.chart.setOption(option)
    }
  }
}
</script>
