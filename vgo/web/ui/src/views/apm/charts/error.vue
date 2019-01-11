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
    backgroundColor: '#fff',
    title: {
            text: '请求错误率',
            textStyle: {
                fontWeight: 'normal',
                fontSize: 16
            },
            left: 'center'
        },
    tooltip: {
        trigger: 'axis'
    },
    xAxis: [
        {
        type: 'category',
        boundaryGap: false,
        axisLine: {
         
        },
        axisLabel: {
            margin: 10
        },
        axisTick: {
            show: false
        },
        data: this.dateList
    }],
    grid: {
            left: '4%',
            right: '2%',
            bottom: '8%',
            top:'14%',
            containLabel: true
        },
    yAxis: [{
        type: 'value',
        name: '单位（%）',
        axisTick: {
            show: false
        },
		max:100,
        axisLabel: {
            margin: 10
        },
        splitLine: {
            show: false,
            lineStyle: {
                color: '#57617B'
            }
        }
    }],
    series: [ {
        name: '',
        type: 'line',
        stack: '总量',
        smooth: true,
        symbol: 'circle',
        symbolSize: 5,
        showSymbol: false,
        animationDelay: 0,
        animationDuration: 1000,
    
        lineStyle: {
            normal: {
                width: 1,
                color: {
                    type: 'linear',
                    x: 0,
                    y: 0,
                    x2: 1,
                    y2: 0,
                    colorStops: [{
                        offset: 0, color: 'red' // 0% 处的颜色
                    }, {
                        offset: 1, color: 'yellowgreen' // 100% 处的颜色
                    }],
                    globalCoord: false // 缺省为 false
                },
                opacity: 0.9
            }
        },
        areaStyle: {
            normal: {
                color: new echarts.graphic.LinearGradient(0, 0, 0, 1, [{
                    offset: 0,
                    color: 'rgba(219, 50, 51, 0.3)'
                }, {
                    offset: 0.8,
                    color: 'rgba(219, 50, 51, 0)'
                }], false),
                shadowColor: 'rgba(0, 0, 0, 0.1)',
                shadowBlur: 10
            }
        },
        data: this.valueList
    }, ]
};
      this.chart.setOption(option)
    }
  }
}
</script>
