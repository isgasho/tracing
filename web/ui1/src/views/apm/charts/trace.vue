<template>
    <div></div>
</template>

<script>
import { formatTime } from '@/utils/tools'
export default {
    name: '',
    props: ['graphData'],
    data () {
        return {
            tracesChart: null,
            currentNode: null
        }
    },
    mounted () {
        this.initTraceChart()
        let self = this;
        setTimeout(function () {
            // 解析出具体的data
            for (var i=0;i< self.graphData.value.agentNodeModels.length;i++) {
                if (self.graphData.value.agentNodeModels[i].key == self.graphData.value.nowNodeKey) {
                    self.traceChart(self.graphData.value.agentNodeModels[i])
                    break
                }
            }
        }, 100);
    },
    watch: {
    },
    computed: {},
    methods: {
         traceChart: function (node) {
            this.currentNode = node
            this.tracesChart.series = node.agentSources.traceHis.traceSeries;
            // serries颜色，错误: color: 'rgba(223, 83, 83, .5)'  成功 rgb(18, 147, 154)
            this.tracesChart.xAxis.tickPositions = node.timeXticks;
            this.tracesChart.subtitle.text = node.agentSources.traceHis.subTitle;
            Highcharts.chart(this.$el, this.tracesChart);
        },
        initTraceChart: function () {
            Highcharts.setOptions({
                global: {
                    useUTC: false
                }
            });
            this.tracesChart = {
                chart: {
                    type: 'scatter',
                    zoomType: 'xy',
                    events: {
                        selection: selectPointsByDrag
                    }
                },
                title: {
                    text: "链路",
                    y:10
                },
                subtitle: {
                    verticalAlign: 'top',
                    align: 'center',
                    y: 26
                },
                xAxis: {
                    title: {
                        enabled: false
                    },
                    type: 'datetime',
                    labels: {
                        format: '{value: %m-%d %H:%M:%S}',
                        setp: 1
                    },
                    tickPositions: [],
                    startOnTick: false,
                    endOnTick: false,
                    showLastLabel: false,
                    tickWidth: 0,
                    gridLineWidth: 0
                },
                credits: {
                    enabled: false
                },
                yAxis: {
                    title: {
                        enabled: true,
                        align: 'high',
                        text: '(ms)',
                        rotation: 0,
                        margin: 0,
                        x: 30,
                        y: -10
                    },
                    minorGridLineWidth: 1,
                    gridLineWidth: 1
                },
                legend: {
                    layout: 'vertical',
                    align: 'left',
                    verticalAlign: 'top',
                    x: 55,
                    y: 13,
                    floating: true,
                    backgroundColor: (Highcharts.theme && Highcharts.theme.legendBackgroundColor) || '#FFFFFF'
                },
                plotOptions: {
                    scatter: {
                        marker: {
                            radius: 5,
                            states: {
                                hover: {
                                    enabled: true,
                                    // lineColor: 'rgb(100,100,100)'
                                }
                            }
                        },
                        states: {
                            hover: {
                                marker: {
                                    enabled: true
                                }
                            }
                        },
                        tooltip: {
                            pointFormat: '{point.x: %m-%d %H:%M:%S}, {point.y}ms '
                        },
                        enableMouseTracking: false,
                        turboThreshold: "disable"
                    }
                },
                series: []
            };
            
            var self = this
            function selectPointsByDrag(e) {
                if (e.xAxis && e.yAxis) {
                    var traces = [];
                    if (self.tracesChart.series[0].visible === undefined || self.tracesChart.series[0].visible === true) {
                        var succesdata = self.currentNode.agentSources.traceHis.traceSeries[0].data;
                        for (var i = 0; i < succesdata.length; i++) {
                            var point = succesdata[i];
                            if (point.x >= e.xAxis[0].min && point.x <= e.xAxis[0].max &&
                                point.y >= e.yAxis[0].min && point.y <= e.yAxis[0].max) {
                                // traces += point.traceId + ":" + point.agentId + ":" + point.startTime + ","
                                var trace = {
                                    traceId: point.traceId,
                                    agentId: point.agentId,
                                    elapsed: point.y,
                                    errCode: 0,
                                    url: point.url,
                                    showTime: formatTime(point.x),
                                    startTime: point.startTime,
                                    traceIp: point.traceIp
                                };
                                traces.push(trace);
                              
                            }
                        }
                    }
                    if (self.tracesChart.series[1].visible === undefined || self.tracesChart.series[1].visible === true) {
                        var errordata = self.currentNode.agentSources.traceHis.traceSeries[1].data;
                        for (var j = 0; j < errordata.length; j++) {
                            var point2 = errordata[j];
                            if (point2.x >= e.xAxis[0].min && point2.x <= e.xAxis[0].max &&
                                point2.y >= e.yAxis[0].min && point2.y <= e.yAxis[0].max) {
                                // traces += point2.traceId + ":" + point2.agentId + ":" + point2.startTime + ","
                                var trace2 = {
                                    traceId: point2.traceId,
                                    agentId: point2.agentId,
                                    elapsed: point2.y,
                                    errCode: 1,
                                    url: point2.url,
                                    showTime: formatTime(point2.x),
                                    startTime: point2.startTime,
                                    traceIp: point2.traceIp
                                };
                                traces.push(trace2);
                            }
                        }
                    }

                    if (traces !== "") {
                        self.$emit("selTraces", traces) 
                    }
                }
                return false
            }
        },
    }
}
</script>

<style>

</style>
