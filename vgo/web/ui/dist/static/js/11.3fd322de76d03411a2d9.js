webpackJsonp([11],{eP2N:function(t,e,a){"use strict";Object.defineProperty(e,"__esModule",{value:!0});var s=a("woOf"),i=a.n(s),n=a("mvHQ"),r=a.n(n),l=a("M+K7"),o=a("KI+6"),c=a("qorP"),p=a("vLgD"),m={name:"strategy",cmponents:{clipboard:o.a,clip:c.a},data:function(){return{services:[],selectedService:"",strategies:[],dialogStatus:"",dialogTitle:"",tempStrategy:{content:{}},bwlistVisible:!1,tempBWKey:"",tempBWVal:"",selbwType:1,apiShow:!1,isCretaeSubmit:!1}},methods:{handleSelService:function(t){this.$store.dispatch("setService",t),this.selectedService=t,this.loadStrategy(t)},handleCopy:function(t,e){var a=r()(t);Object(c.a)(a,e),this.$message({message:"Copied",type:"success",duration:2e3})},changeStatus:function(t){var e=this,a="";a=0==t.status?'You will start using strategy："'+t.name+'"！When started, all the apis using will  start using this strategy！':'You will stop using strategy："'+t.name+'"！When stopped, all the apis will stop using this strategy',this.$confirm(a,"Warning",{dangerouslyUseHTMLString:!0,cancelButtonText:"Cancel",confirmButtonText:"Submit",type:"info"}).then(function(){var a={target_app:"juzManage",target_path:"/manage/strategy/change",id:t.id,name:t.name,status:t.status};Object(l.a)("POST",a).then(function(a){0==t.status?t.status=1:t.status=0,e.$message({message:"Status changed ok",type:"success",duration:3e3,center:!0})})})},delStrategy:function(t){var e=this;this.$confirm('You will delete this strategy："'+t.name+'"! When deleted, all the apis will detele this strategy',"Warning",{dangerouslyUseHTMLString:!0,cancelButtonText:"Cancel",confirmButtonText:"Submit",type:"info"}).then(function(){var a={target_app:"juzManage",target_path:"/manage/strategy/delete",id:t.id,name:t.name,service:t.service,type:t.type};Object(l.a)("POST",a).then(function(t){e.loadStrategy(e.selectedService),e.$message({message:"Delete strategy ok",type:"success",duration:3e3,center:!0})})})},apiSet:function(t){this.apiShow=!0,this.tempStrategy=t},submitEdit:function(){var t=this;""==this.tempStrategy.name&&this.$message({message:"Strategy name cant be empty",type:"warning",duration:3e3,center:!0}),3==this.tempStrategy.type&&(this.tempStrategy.content.param=this.tempStrategy.content.param.trim());var e=i()({},this.tempStrategy);e.content=r()(e.content);var a={target_app:"juzManage",target_path:"/manage/strategy/update",strategy:e};Object(l.a)("POST",a).then(function(e){1!=t.tempStrategy.type&&(t.bwlistVisible=!1),t.loadStrategy(t.selectedService),t.$message({message:"Submit ok",type:"success",duration:3e3,center:!0})})},handleEdit:function(t){this.dialogStatus="edit",this.dialogTitle=this.$t("juz.editStrategy"),this.tempStrategy=i()({},t),this.tempStrategy.content=JSON.parse(this.tempStrategy.content),this.bwlistVisible=!0,this.tempBWKey="ip",this.isCretaeSubmit=!1},removeBW:function(t){for(var e=0;e<this.tempStrategy.content.length;e++)t.key==this.tempStrategy.content[e].key&&t.val==this.tempStrategy.content[e].val&&this.tempStrategy.content.splice(e,1)},addBW:function(){if(""!=this.tempBWKey&&""!=this.tempBWVal){for(var t=0;t<this.tempStrategy.content.length;t++)if(this.tempBWKey==this.tempStrategy.content[t].key&&this.tempBWVal==this.tempStrategy.content[t].val)return void this.$message({message:"List item already exist",type:"warning",duration:3e3,center:!0});this.tempStrategy.content.push({key:this.tempBWKey,val:this.tempBWVal}),1==this.selbwType?this.tempBWKey="ip":this.tempBWKey="",this.tempBWVal=""}else this.$message({message:"key、val cant be empty",type:"warning",duration:3e3,center:!0})},submitStrategy:function(){var t=this;""==this.tempStrategy.name&&this.$message({message:"Strategy name cant empty",type:"warning",duration:3e3,center:!0}),3==this.tempStrategy.type&&(this.tempStrategy.content.param=this.tempStrategy.content.param.trim());var e=i()({},this.tempStrategy);e.content=r()(e.content);var a={target_app:"juzManage",target_path:"/manage/strategy/create",strategy:e};Object(l.a)("POST",a).then(function(e){t.isCretaeSubmit=!1,1!=t.tempStrategy.type&&(t.bwlistVisible=!1),t.loadStrategy(t.selectedService),t.$message({message:"Submit ok",type:"success",duration:3e3,center:!0})})},selBwType:function(t){this.tempBWKey=1==t?"ip":""},selType:function(t){switch(t){case 1:this.tempStrategy={name:this.tempStrategy.name,type:1,sub_type:1,content:[],service:this.selectedService};break;case 2:this.tempStrategy={name:this.tempStrategy.name,type:2,sub_type:0,service:this.selectedService,content:{req_timeout:15,retry_times:0,retry_interval:3}};break;case 3:this.tempStrategy={name:this.tempStrategy.name,type:3,sub_type:0,service:this.selectedService,content:{qps:-1,concurrent:-1,param:"",span:2,times:1,fuse_error:0,fuse_error_count:20,fuse_percent:50,fuse_recover:25,fuse_recover_count:10}}}},handleCreate:function(){""!=this.selectedService?(this.dialogStatus="create",this.dialogTitle=this.$t("juz.createStrategy"),this.tempBWKey="ip",this.bwlistVisible=!0,this.isCretaeSubmit=!0,this.tempStrategy={type:1,sub_type:1,content:[],name:"",service:this.selectedService}):this.$message({message:"select a service first",type:"warning",duration:3e3,center:!0})},selStrategy:function(t){this.selectedStrategy=t},loadStrategy:function(t){var e=this,a={target_app:"juzManage",target_path:"/manage/strategy/load",service:this.selectedService,type:0};Object(l.a)("POST",a).then(function(t){e.strategies=t.data.data})},calcService:function(){return this.selectedService||this.$store.getters.service},loadServices:function(){var t=this;Object(p.a)({url:"/ops/service/query",method:"GET",params:{}}).then(function(e){t.services=e.data.data})}},created:function(){this.loadServices(),this.selectedService=this.$store.getters.service,""!=this.selectedService&&this.loadStrategy(this.selectedService)},destroyed:function(){}},u={render:function(){var t=this,e=t.$createElement,a=t._self._c||e;return a("div",{staticClass:"app-container"},[a("div",{staticClass:"filter-container"},[a("el-tag",[t._v(" "+t._s(t.$t("common.service")))]),t._v(" "),a("el-select",{staticClass:"filter-item",staticStyle:{width:"200px"},attrs:{clearable:"",value:t.calcService(),placeholder:"select a service"},on:{change:t.handleSelService}},t._l(t.services,function(t){return a("el-option",{key:t.name,attrs:{label:t.name,value:t.name}})})),t._v(" "),a("el-button",{staticClass:"filter-item",staticStyle:{"margin-left":"10px"},attrs:{type:"success",icon:"el-icon-edit"},on:{click:t.handleCreate}},[t._v(t._s(t.$t("juz.addStrategy")))])],1),t._v(" "),a("div",{staticClass:"table"},[a("el-table",{staticStyle:{width:"100%"},attrs:{data:t.strategies,fit:"","highlight-current-row":""}},[a("el-table-column",{attrs:{align:"left",label:t.$t("common.type"),width:"150",prop:"name"},scopedSlots:t._u([{key:"default",fn:function(e){return[1==e.row.type?a("span",[t._v(t._s(t.$t("juz.bwList")))]):t._e(),t._v(" "),2==e.row.type?a("span",[t._v(t._s(t.$t("juz.timoutRetry")))]):t._e(),t._v(" "),3==e.row.type?a("span",[t._v(t._s(t.$t("juz.trafficControl")))]):t._e()]}}])}),t._v(" "),a("el-table-column",{attrs:{align:"left",label:t.$t("common.name"),width:"200",prop:"name"},scopedSlots:t._u([{key:"default",fn:function(e){return[a("span",[t._v(t._s(e.row.name))])]}}])}),t._v(" "),a("el-table-column",{attrs:{width:"180",align:"left",label:t.$t("common.updateDate")},scopedSlots:t._u([{key:"default",fn:function(e){return[a("span",[t._v(t._s(e.row.modify_date))])]}}])}),t._v(" "),a("el-table-column",{attrs:{width:"150",align:"left",label:t.$t("common.status")},scopedSlots:t._u([{key:"default",fn:function(e){return[0==e.row.status?a("span",[a("el-tag",{staticStyle:{border:"none"},attrs:{type:"warning"}},[t._v("Off")])],1):a("span",[a("el-tag",{staticStyle:{border:"none"},attrs:{type:"success"}},[t._v("On")])],1)]}}])}),t._v(" "),a("el-table-column",{attrs:{align:"center",label:t.$t("common.operate"),"class-name":"small-padding fixed-width"},scopedSlots:t._u([{key:"default",fn:function(e){return[a("span",{staticClass:"table-op-btn",on:{click:function(a){t.handleEdit(e.row)}}},[t._v(t._s(t.$t("common.edit")))]),t._v(" "),a("span",{staticClass:"table-op-btn",on:{click:function(a){t.handleCopy(e.row,a)}}},[t._v(t._s(t.$t("common.copyConfig")))]),t._v(" "),a("span",{staticClass:"table-op-btn",on:{click:function(a){t.changeStatus(e.row)}}},[0==e.row.status?a("span",[t._v(t._s(t.$t("common.start")))]):a("span",[t._v(t._s(t.$t("common.stop")))])]),t._v(" "),a("span",{staticClass:"table-op-btn",on:{click:function(a){t.delStrategy(e.row)}}},[t._v("Delete")])]}}])})],1)],1),t._v(" "),a("el-dialog",{staticClass:"mf-dialog",attrs:{title:t.dialogTitle,visible:t.bwlistVisible},on:{"update:visible":function(e){t.bwlistVisible=e}}},[a("el-form",{staticStyle:{width:"650px","margin-left":"50px"},attrs:{"label-position":"left","label-width":"120px",size:"mini"}},[a("div",{staticClass:"form-block"},[a("span",[t._v(t._s(t.$t("juz.basic")))]),t._v(" "),a("el-form-item",{staticStyle:{"margin-top":"10px"},attrs:{label:t.$t("common.name")}},[a("el-input",{staticStyle:{width:"300px"},attrs:{placeholder:""},model:{value:t.tempStrategy.name,callback:function(e){t.$set(t.tempStrategy,"name",e)},expression:"tempStrategy.name"}})],1),t._v(" "),a("el-form-item",{staticStyle:{"margin-top":"10px"},attrs:{label:t.$t("common.type")}},[a("el-radio-group",{attrs:{disabled:!t.isCretaeSubmit},on:{change:t.selType},model:{value:t.tempStrategy.type,callback:function(e){t.$set(t.tempStrategy,"type",e)},expression:"tempStrategy.type"}},[a("el-radio",{attrs:{label:1}},[t._v(t._s(t.$t("juz.bwList")))]),t._v(" "),a("el-radio",{attrs:{label:2}},[t._v(t._s(t.$t("juz.timoutRetry")))]),t._v(" "),a("el-radio",{attrs:{label:3}},[t._v(t._s(t.$t("juz.trafficControl")))])],1)],1)],1),t._v(" "),1==t.tempStrategy.type?a("div",{staticClass:"form-block"},[a("span",[t._v(t._s(t.$t("juz.bwList")))]),t._v(" "),a("el-form-item",{staticStyle:{"margin-top":"10px"},attrs:{label:t.$t("common.type")}},[a("el-radio-group",{model:{value:t.tempStrategy.sub_type,callback:function(e){t.$set(t.tempStrategy,"sub_type",e)},expression:"tempStrategy.sub_type"}},[a("el-radio",{attrs:{label:1}},[t._v("Black")]),t._v(" "),a("el-radio",{attrs:{label:2}},[t._v("White")])],1)],1),t._v(" "),1==t.tempStrategy.sub_type||2==t.tempStrategy.sub_type?a("el-form-item",{staticStyle:{"margin-top":"10px"},attrs:{label:t.$t("common.add")}},[a("el-radio-group",{on:{change:t.selBwType},model:{value:t.selbwType,callback:function(e){t.selbwType=e},expression:"selbwType"}},[a("el-radio",{attrs:{label:1}},[t._v("IP")]),t._v(" "),a("el-radio",{attrs:{label:2}},[t._v(t._s(t.$t("common.param")))])],1),t._v(" "),a("div",{staticStyle:{}},[a("el-input",{staticStyle:{width:"80px"},attrs:{placeholder:"Key",disabled:1==t.selbwType},model:{value:t.tempBWKey,callback:function(e){t.tempBWKey=e},expression:"tempBWKey"}}),t._v(" "),a("el-input",{staticStyle:{width:"160px"},attrs:{placeholder:"Val"},model:{value:t.tempBWVal,callback:function(e){t.tempBWVal=e},expression:"tempBWVal"}}),t._v(" "),a("el-button",{staticStyle:{"margin-left":"10px"},attrs:{size:"mini",icon:"el-icon-plus",circle:""},on:{click:t.addBW}})],1)],1):t._e(),t._v(" "),1==t.tempStrategy.sub_type||2==t.tempStrategy.sub_type?a("el-form-item",{attrs:{label:t.$t("common.currentList")}},t._l(t.tempStrategy.content,function(e){return a("div",{key:e.key,staticStyle:{"margin-top":"-2px"}},[a("el-tag",{staticStyle:{width:"70px",border:"none"},attrs:{type:"success",size:"large"}},[t._v(t._s(e.key))]),t._v(" "),a("el-tag",{staticStyle:{width:"200px",border:"none"},attrs:{type:"info",size:"large"}},[t._v(t._s(e.val))]),t._v(" "),a("el-button",{staticStyle:{"margin-left":"10px"},attrs:{type:"text",size:"mini",icon:"el-icon-minus"},nativeOn:{click:function(a){t.removeBW(e)}}})],1)})):t._e()],1):t._e(),t._v(" "),2==t.tempStrategy.type?a("div",{staticClass:"form-block"},[a("span",[t._v(t._s(t.$t("juz.timoutRetry")))]),t._v(" "),a("el-form-item",{staticStyle:{width:"200px","margin-top":"10px"},attrs:{label:t.$t("juz.timeout")}},[a("el-tooltip",{attrs:{content:"0<X<=60",placement:"top"}},[a("el-input-number",{attrs:{min:1,max:60},model:{value:t.tempStrategy.content.req_timeout,callback:function(e){t.$set(t.tempStrategy.content,"req_timeout",e)},expression:"tempStrategy.content.req_timeout"}})],1)],1),t._v(" "),a("el-form-item",{staticStyle:{width:"200px"},attrs:{label:t.$t("juz.retryTimes")}},[a("el-tooltip",{attrs:{content:t.$t("juz.retryTimesTips"),placement:"top"}},[a("el-input-number",{attrs:{min:0,max:5},model:{value:t.tempStrategy.content.retry_times,callback:function(e){t.$set(t.tempStrategy.content,"retry_times",e)},expression:"tempStrategy.content.retry_times"}})],1)],1),t._v(" "),a("el-form-item",{staticStyle:{width:"200px"},attrs:{label:t.$t("juz.retryIntv")}},[a("el-tooltip",{attrs:{content:t.$t("juz.retryIntvTips"),placement:"top"}},[a("el-input-number",{attrs:{min:1,max:30},model:{value:t.tempStrategy.content.retry_interval,callback:function(e){t.$set(t.tempStrategy.content,"retry_interval",e)},expression:"tempStrategy.content.retry_interval"}})],1)],1)],1):t._e(),t._v(" "),3==t.tempStrategy.type?a("div",[a("div",{staticClass:"form-block"},[a("span",[t._v(t._s(t.$t("juz.trafficControl")))]),t._v(" "),a("el-form-item",{staticStyle:{"margin-top":"10px"},attrs:{label:"QPS"}},[a("el-input-number",{attrs:{min:0,max:1e4},model:{value:t.tempStrategy.content.qps,callback:function(e){t.$set(t.tempStrategy.content,"qps",e)},expression:"tempStrategy.content.qps"}}),t._v(" "),a("el-alert",{attrs:{title:t.$t("juz.qpsTips"),closable:!1,type:"success"}})],1),t._v(" "),a("el-form-item",{attrs:{label:t.$t("juz.concurrents")}},[a("el-input-number",{attrs:{min:0,max:1e4},model:{value:t.tempStrategy.content.concurrent,callback:function(e){t.$set(t.tempStrategy.content,"concurrent",e)},expression:"tempStrategy.content.concurrent"}}),t._v(" "),a("el-alert",{attrs:{title:t.$t("juz.concurrentsTips"),closable:!1,type:"success"}})],1)],1),t._v(" "),a("div",{staticClass:"form-block"},[a("span",[t._v(t._s(t.$t("juz.userQuota")))]),t._v(" "),a("el-form-item",{staticStyle:{"margin-top":"10px",width:"400px"},attrs:{label:t.$t("common.param")}},[a("el-input",{attrs:{placeholder:t.$t("juz.emptyNoLimit")},model:{value:t.tempStrategy.content.param,callback:function(e){t.$set(t.tempStrategy.content,"param",e)},expression:"tempStrategy.content.param"}})],1),t._v(" "),a("el-form-item",{attrs:{label:t.$t("juz.timeSpan")}},[a("el-input-number",{attrs:{min:2,max:2592e3},model:{value:t.tempStrategy.content.span,callback:function(e){t.$set(t.tempStrategy.content,"span",e)},expression:"tempStrategy.content.span"}}),t._v(t._s(t.$t("common.second"))+"\n                ")],1),t._v(" "),a("el-form-item",{attrs:{label:t.$t("juz.times")}},[a("el-input-number",{attrs:{min:1,max:1024e4},model:{value:t.tempStrategy.content.times,callback:function(e){t.$set(t.tempStrategy.content,"times",e)},expression:"tempStrategy.content.times"}}),t._v(" "),a("el-alert",{attrs:{title:t.$t("juz.userQuotaTips"),closable:!1,type:"success"}})],1)],1)]):t._e()]),t._v(" "),a("div",{staticClass:"dialog-footer",attrs:{slot:"footer"},slot:"footer"},[a("el-button",{on:{click:function(e){t.bwlistVisible=!1}}},[t._v(t._s(t.$t("common.cancel")))]),t._v(" "),t.isCretaeSubmit?a("el-button",{attrs:{type:"primary"},on:{click:t.submitStrategy}},[a("span",[t._v(t._s(t.$t("common.submit")))])]):a("el-button",{attrs:{type:"primary"},on:{click:t.submitEdit}},[a("span",[t._v(t._s(t.$t("common.submit")))])])],1)],1)],1)},staticRenderFns:[]},y=a("VU/8")(m,u,!1,null,null,null);e.default=y.exports}});