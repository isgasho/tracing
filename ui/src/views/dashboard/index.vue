<template>
  <div>
    <Row
      class="subnav"
      style="background-color:#595959;border-bottom: 1px solid #474747;color:#eaeaea;padding: 2px 10px;padding-bottom:1px;padding-left:40px;vertical-align:middle;font-size:13px;margin-top:0px;height:36px"
    >
      <span class="item hover-cursor app-map" @click="navShow(1)">
        <span class="count" :class="{'bg-second': showItem==1}">应用地图</span>
      </span>
      <span class="item hover-cursor" @click="navShow(2)">
        <span class="count" :class="{'bg-second': showItem==2}">应用列表</span>
      </span>

      <span class="topology-filter">
        <Select
          v-model="selTopologyDate"
          style="width:100px;margin-left:40px"
          placeholder="选择时间"
          @on-change="selDate"
        >
          <Option value="3">最近3分钟</Option>
          <Option value="10">最近10分钟</Option>
          <Option value="30">最近30分钟</Option>
          <Option value="60">最近1小时</Option>
          <Option value="360">最近6小时</Option>
          <Option value="1440">最近1天</Option>
          <Option value="4320">最近3天</Option>
        </Select>
        <Checkbox
          v-model="topologyError"
          class="margin-left-10"
          size="small"
          v-show="showItem==1"
          @on-change="setMapError"
        >仅显示错误</Checkbox>
        <Checkbox
          v-model="topologyRealtime"
          size="small"
          v-show="showItem==1"
          @on-change="setMapRefresh"
        >实时刷新</Checkbox>
        <Select
          class="select-app"
          v-model="topologyHighlighted"
          filterable
          multiple
          style="width:250px;border-right:.5px solid #aaa;padding-right:10px"
          placeholder="特定应用高亮显示"
          :max-tag-count="1"
          v-show="showItem==1"
          @on-change="mapHighlightChange"
        >
          <Option v-for="item in appNames" :value="item" :key="item">{{ item }}</Option>
        </Select>

        <Input
          class="margin-left-10"
          v-model="mapErrorFilter"
          placeholder="自定义错误阈值 e.g. count>10,error>30,duration>300 "
          style="width: 350px;border:none;"
          @on-blur="setMapErrorFilter"
          @on-enter="setMapErrorFilter"
          v-show="showItem==1"
        />
      </span>

      <span style="float:right;margin-right:-50px;"   v-show="showItem==2">
        <span class="item">
          应用总数：
          <span class="count bg-second">{{appList.length}}</span>
        </span>
        <span class="item">
          不健康应用数：
          <span class="count bg-second">0</span>
        </span>
      </span>
    </Row>
    <div v-show="showItem==1" class="no-border">
      <span style="float:right;margin-right:20px">
        <Tag style="width:30px;height:15px;border:none" :style="{background:tagColor(0)}"></Tag>
        <span class="font-size-10">普通应用</span>
        <Tag
          style="width:30px;height:15px;border:none"
          class="margin-left-10"
          :style="{background:tagColor(1)}"
        ></Tag>
        <span class="font-size-10">数据库/中间件</span>
        <!-- <Tag
          style="width:30px;height:15px;border:none"
          class="margin-left-10"
          :style="{background:tagColor(3)}"
        ></Tag>
        <span class="font-size-10">高亮应用</span>-->
      </span>

      <div
        :id="id"
        class="app-service-map"
        style="width:calc(100vw - 180px);height:calc(100vh - 100px)"
      ></div>
    </div>
    <div class="app-container" v-show="showItem==2">
      <Row style="padding:0 10px;" class="split-border-bottom no-border">
        <Col span="17" class="split-border-right">
          <span class="padding-bottom-5 font-size-18">应用列表</span>
          <Select
            v-model="selApps"
            filterable
            multiple
            style="width:350px;border:none;float:right;margin-right:20px"
            placeholder="过滤应用"
            :max-tag-count="2"
          >
            <Option v-for="app in appList" :value="app.name" :key="app.name">{{ app.name }}</Option>
          </Select>
        </Col>
        <Col span="6">
          <span class="padding-bottom-5 font-size-18 margin-left-10">总体动态</span>
        </Col>
      </Row>
      <Row style="padding:0 10px">
        <Col span="17" class="split-border-right no-border" style="padding:8px 10px;">
          <Table
            :columns="appLabels"
            :data="showAppList()"
            class="margin-top-15"
            :row-class-name="rowClassName"
            @on-row-click="gotoApp"
          ></Table>
          <Page :current="1" :total="totalApps" size="small" class="margin-top-15" simple/>
        </Col>
        <Col span="6" style="padding:8px 10px;padding-left:20px">
          <div class="margin-top-10 card-tab">
            <Button type="primary" ghost>告警通知</Button>
            <Button>事件日志</Button>
          </div>
          <div>
            <Icon
              type="ios-happy-outline"
              class="margin-top-20 color-primary2 margin-left-20"
              style="font-size:60px"
            />
          </div>
          <div class="margin-top-10 font-size-18 margin-left-5">恭喜，当前没有任何告警</div>
        </Col>
      </Row>
    </div>
  </div>
</template>

<script>
import request from "@/utils/request";
import echarts from "echarts";
export default {
  name: "appList",
  data() {
    return {
      id: "service-map-id",
      chart: null,
      lineLength: 200,
      primaryNodeSize: 50,
      smallNodeSize: 35,
      lineLabelSize: 13,
      repulsion: 500,
      mapLinkColor: "#12b5d0",
      mapLinkErrorColor: "#ff0000",
      appNames: [],
      selApps: [],
      mapRefreshTimerID: "",
      mapErrorFilter:  this.$store.state.apm.errorFilterNav,
      mapErrorFilterRes: {},
      appLabels: [
        {
          title: "应用名",
          key: "name"
        },
        {
          title: "请求数",
          key: "count"
        },
        {
          title: "响应时间(ms)",
          key: "average_elapsed"
        },
        {
          title: "错误率(%)",
          key: "error_percent"
        },
        {
          title: "Apdex",
          key: "apdex"
        }
      ],
      appList: [],

      totalApps: 3,

      showItem: this.$store.state.apm.dashNav,

      selTopologyDate: this.$store.state.apm.dashSelDate,
      topologyError: false,
      topologyRealtime: false,
      topologyHighlighted: [],

      mapOption: {}
    };
  },
  computed: {},
  beforeDestroy() {
    this.destroyChart();
    clearInterval(this.mapRefreshTimerID);
  },
  methods: {
    // 应用地图，设置错误阈值
    setMapErrorFilter() {
      this.parseErrorFilter()
      this.chart.setOption(this.mapOption)
    },
    parseErrorFilter() {
      if (this.mapErrorFilter == '') {
         this.mapErrorFilterRes = {}
         for (var i = 0; i < this.mapOption.series[0].links.length; i++) {
            var link = this.mapOption.series[0].links[i];
            link.error = false
            var color = this.mapLinkColor;
            if (link.access_count > 0) {
                if (link.error_count / link.access_count > 0.2) {
                    link.error = true;
                    color = this.mapLinkErrorColor;
                }
            }

            link.lineStyle.normal.color = color
            // 仅显示错误处理
            if (this.topologyError) {
                if (link.error) {
                    link.lineStyle.normal.color = "transparent";
                }
            }
        }

        this.$store.dispatch("setErrorFilterNav", "");
          return 
      }
         // 解析设定字符串，格式：count<10 error>30 duration>300
      var splited = [];
      var filter = this.mapErrorFilter.trim();
      var s = filter.split(",");
      for (var i = 0; i < s.length; i++) {
        splited.push(s[i]);
      }

      var temp = {};
      for (var i = 0; i < splited.length; i++) {
        var s = splited[i];
        s = s.trim();
        if (
          s.indexOf("count") == 0 ||
          s.indexOf("error") == 0 ||
          s.indexOf("duration") == 0
        ) {
          var s1 = s.split("<");
          if (s1.length == 2) {
            // 判断是否是数字
            var n = parseInt(s1[1]);
            if (isNaN(n)) {
              this.$Message.warning("输入的错误阈值字符串不合法2");
              return;
            }

            temp[s1[0]] = {
              compare: "<",
              value: n
            };
            continue;
          }

          var s2 = s.split(">");
          if (s2.length == 2) {
            // 判断是否是数字
            var n = parseInt(s2[1]);
            if (isNaN(n)) {
              this.$Message.warning("输入的错误阈值字符串不合法4");
              return;
            }

            temp[s2[0]] = {
              compare: ">",
              value: n
            };
            continue;
          }

          this.$Message.warning("输入的错误阈值字符串不合法3");
          return;
        } else {
          console.log(s);
          this.$Message.warning("输入的错误阈值字符串不合法5");
          return;
        }
      }
      this.mapErrorFilterRes = temp;
      console.log(this.mapErrorFilterRes);

      // 根据设置好的条件，重新标示错误
      for (var i = 0; i < this.mapOption.series[0].links.length; i++) {
        var link = this.mapOption.series[0].links[i];

        var countFilter = this.mapErrorFilterRes['count'];

        var isError = false;
        if (countFilter != undefined) {
          if (countFilter.compare == ">") {
            if (link.access_count >= countFilter.value) {
              isError = true;
            }
          } else {
            if (link.access_count < countFilter.value) {
              isError = true;
            }
          }
        }

        if (!isError) {
          var filter = this.mapErrorFilterRes['error'];
          if (filter != undefined) {
            if (filter.compare == ">") {
                if (link.error_count >= filter.value) {
                    isError = true;
                }
            } else {
                if (link.error_count < filter.value) {
                    isError = true;
                }
            }
          }
        }

         if (!isError) {
          var filter = this.mapErrorFilterRes['duration'];
          if (filter != undefined) {
            if (filter.compare == ">") {
                if (link.avg >= filter.value) {
                    isError = true;
                }
            } else {
                if (link.avg < filter.value) {
                    isError = true;
                }
            }
          }
        }


        if (isError) {
          link.error = true;
          link.lineStyle.normal.color = this.mapLinkErrorColor;
        } else {
          link.error = false;
          link.lineStyle.normal.color = this.mapLinkColor;
          // 仅显示错误处理
          if (this.topologyError) {
            link.lineStyle.normal.color = "transparent";
          }
        }


        // 提交到store
        this.$store.dispatch("setErrorFilterNav", this.mapErrorFilter);
      }
    },
    // 设定后，每60秒刷新一次应用地图
    setMapRefresh(refresh) {
      var _this = this;
      if (refresh) {
        _this.initServiceMap();
        _this.$Message.success("应用地图实时刷新！");

        _this.mapRefreshTimerID = setInterval(function() {
          _this.initServiceMap();
          _this.$Message.success("应用地图实时刷新！");
        }, 60000);
      } else {
        clearInterval(this.mapRefreshTimerID);
      }
    },
    setMapError(error) {
      for (var i = 0; i < this.mapOption.series[0].links.length; i++) {
        console.log(this.mapOption.series[0].links[i].source,this.mapOption.series[0].links[i].target,this.mapOption.series[0].links[i].error)
        if (!this.mapOption.series[0].links[i].error) {
          if (error) {
            this.mapOption.series[0].links[i].lineStyle.normal.color =
              "transparent";
          } else {
            this.mapOption.series[0].links[
              i
            ].lineStyle.normal.color = this.mapLinkColor;
          }
        } 
      }
      this.chart.setOption(this.mapOption);
    },
    showAppList() {
      var apps = [];
      if (this.selApps.length == 0) {
        apps = this.appList;
      } else {
        for (var i = 0; i < this.appList.length; i++) {
          for (var j = 0; j < this.selApps.length; j++) {
            if (this.selApps[j] == this.appList[i].name) {
              apps.push(this.appList[i]);
            }
          }
        }
      }

      return apps;
    },
    mapHighlightChange(v) {
      for (var i = 0; i < this.mapOption.series[0].data.length; i++) {
        var node = this.mapOption.series[0].data[i];
        node.symbolSize = this.smallNodeSize;
        for (var j = 0; j < v.length; j++) {
          if (v[j] == node.name) {
            node.symbolSize = this.primaryNodeSize;
          }
        }
      }

      this.chart.setOption(this.mapOption);
    },
    selDate() {
      this.$store.dispatch("setDashSelDate", this.selTopologyDate);
      this.initServiceMap();
      this.initAppList();
    },
    navShow(v) {
      this.showItem = v;
      this.$store.dispatch("setDashNav", v);
      if (v == 1) {
        console.log(this.chart);
        this.chart.setOption(this.mapOption);
      }
    },
    rowClassName(row, index) {
      if (row.error_percent >= 10 || row.apdex < 0.8) {
        return "error-trace";
      } else {
        return "success-trace";
      }
    },
    gotoApp(app) {
    //   this.$store.dispatch("setAPPID", app.id);
      this.$store.dispatch("setAPPName", app.name);
      this.$router.push("/ui/apm");
    },
    destroyChart() {
      if (this.chart) {
        this.chart.dispose();
        this.chart = null;
      }
    },
    tagColor(tp) {
      switch (tp) {
        case 0: //普通应用
          return "linear-gradient(to right, #01acca , #5adbe7)";
          break;
        case 1: // 数据库中间件
          return "linear-gradient(to right, #ffb402 , #ffdc84)";
        case 3: // 当前应用
          return "linear-gradient(to right, #157eff , #35c2ff)";
        default:
          break;
      }
    },
    initServiceMap() {
      this.destroyChart();
      this.$Loading.start();
      request({
        url: "/web/serviceMap",
        method: "GET",
        params: {
          start: this.selTopologyDate
        }
      })
        .then(res => {
          this.$Loading.finish();
          if (res.data.data.nodes.length == 0) {
            this.$Message.info({
              content: "没有查询到数据",
              duration: 3
            });
          } else {
            this.initChart(res.data.data.nodes, res.data.data.links);
            for (var i = 0; i < res.data.data.nodes.length; i++) {
              this.appNames.push(res.data.data.nodes[i].name);
            }
          }
        })
        .catch(error => {
          this.$Loading.error();
        });
    },
    formatLinkLabel(link) {
      var error = 0;
      if (link.access_count > 0) {
        error = (link.error_count * 100) / link.access_count;
        error = error.toFixed(1);
      }
      return link.access_count + "/" + error + "%/" + link.avg + "ms";
    },
    calcSize(nodes) {
      var l = nodes.length;
      if (l < 20) {
        this.lineLength = 200;
        this.primaryNodeSize = 50;
        this.smallNodeSize = 30;
        this.lineLabelSize = 11;
        this.repulsion = 500;
        return;
      }

      if (l < 40) {
        this.lineLength = 200;
        this.primaryNodeSize = 40;
        this.smallNodeSize = 20;
        this.lineLabelSize = 10;
        this.repulsion = 5500;
        return;
      }

      if (l < 80) {
        this.lineLength = 200;
        this.primaryNodeSize = 30;
        this.smallNodeSize = 15;
        this.lineLabelSize = 9;
        this.repulsion = 500;
        return;
      }

      this.lineLength = 100;
      this.primaryNodeSize = 30;
      this.smallNodeSize = 10;
      this.lineLabelSize = 9;
      this.repulsion = 500;
    },
    initChart(nodes, links) {
      this.chart = echarts.init(document.getElementById(this.id));
      for (var j = 0; j < nodes.length; j++) {
        // 设置node的样式
        nodes[j].symbolSize = this.smallNodeSize;
        nodes[j].itemStyle = {
          normal: {
            color: new echarts.graphic.LinearGradient(0, 0, 1, 0, [
              {
                offset: 0,
                color: "#01acca"
              },
              {
                offset: 1,
                color: "#5adbe7"
              }
            ])
          }
        };
        // node错误率超过一个值，则添加特殊显示
        if (nodes[j].span_count > 0) {
          if (nodes[j].error_count / nodes[j].span_count > 0.2) {
            nodes[j].label = {
              normal: {
                color: "#ff0000"
              }
            };
          }
        }

        // 对于数据库/中间件node进行特殊展示
        if (nodes[j].category == 1) {
          nodes[j].itemStyle = {
            normal: {
              color: new echarts.graphic.LinearGradient(0, 0, 1, 0, [
                {
                  offset: 0,
                  color: "#ffb402"
                },
                {
                  offset: 1,
                  color: "#ffdc84"
                }
              ])
            }
          };
        }

        // // 高亮显示
        for (var k = 0; k < this.topologyHighlighted.length; k++) {
          if (this.topologyHighlighted[k] == nodes[j].name) {
            nodes[j].symbolSize = this.primaryNodeSize;
          }
        }
      }

      this.calcSize(nodes);
      for (var i = 0; i < links.length; i++) {
        links[i].label = {
          normal: {
            show: true,
            formatter: this.formatLinkLabel(links[i]),
            fontSize: this.lineLabelSize
          }
        };
        links[i].lineStyle = {
          normal: {
            width: 1,
            curveness: 0
          }
        };
      }
      var option = {
        animationDurationUpdate: 500,
        series: [
          {
            type: "graph",
            layout: "force",
            force: {
              repulsion: 2000,
              edgeLength: this.lineLength,
              layoutAnimation: true
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
      this.mapOption = option;

      // 设置错误高亮显示
      this.parseErrorFilter()
      
      this.chart.setOption(this.mapOption);

      var _this = this
      function nodeOnClick(params) {
         if (params.data.category == 0 && params.name != 'UNKNOWN') {
            // 只有普通类型的应用才可以前往应用页面
            _this.$store.dispatch("setAPPName", params.name);
            _this.$router.push("/ui/apm");
         } else {
             _this.$Message.warning('非普通类型应用不可以前往应用详情页面')
         }
      }
      this.chart.on("click", nodeOnClick);
      //'click'、'dblclick'、'mousedown'、'mousemove'、'mouseup'、'mouseover'、'mouseout'
    },

    initAppList() {
      this.$Loading.start();
      // 加载APPS
      request({
        url: "/web/appListWithSetting",
        method: "GET",
        params: {
          start: this.selTopologyDate
        }
      })
        .then(res => {
          this.appList = res.data.data;
          this.$Loading.finish();
        })
        .catch(error => {
          this.$Loading.error();
        });
    }
  },
  mounted() {
    this.initServiceMap();
    this.initAppList();
  }
};
</script>

<style lang="less">
@import "../../theme/gvar.less";
.ivu-table .error-trace td {
  background-color: rgba(223, 83, 83, 0.5);
  color: #333;
}

.topology-filter {
  .ivu-select-arrow {
    // display:none;
  }
  .ivu-select-selected-value {
    font-size: 12px !important;
  }
  .ivu-checkbox-wrapper {
    font-size: 12px !important;
  }
  input {
    font-size: 12px !important;
    padding-top: 0px;
    padding-bottom: 0px;
  }
  .ivu-select-selection {
    background: transparent;
    border: none;
    // border-bottom: 1px solid white;
    color: white;
  }

  input {
    background: transparent;
    color: white;
  }
}

.ivu-checkbox-wrapper-checked {
  color: @primary-color;
}
.ivu-checkbox {
  display: none;
}
</style>

<style lang="less" scoped>
@import "../../theme/gvar.less";
.subnav {
  .item {
    margin-left: 10px;
    .count {
      display: inline-block;
      padding: 1px 7px;
      padding-bottom: 3px;
      border-radius: 4px;
      text-shadow: 0 1px 2px rgba(0, 0, 0, 0.2);
    }
  }
}
</style>

