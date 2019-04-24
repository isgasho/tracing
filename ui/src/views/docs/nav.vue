<template>
  <div class="mf-doc-nav">
      <Row style="height:100%">
         <Col :xs="{span:0}" :sm="{span:0}" :md="{span: 6}" :lg="{ span: 6}"   class="padding-top-20  mf-doc-nav-left">
            <div class="left-items">
              <div v-for="i in items" :key="i" class="item item-1"  v-if="level[i]==1" >
                  {{names[i]}}
              </div>
              <div class="item hover-cursor item-2" :class="{'selected': selItem==i}" @click="selectItem(i)" v-else>
                {{names[i]}}
              </div>
            </div>
         </Col>

         <Col :xs="{span:24}" :sm="{span:24}" :md="{span: 18}" :lg="{ span: 18}" class="mf-doc-nav-right" style="padding-bottom:30px">
            <div class="page-header align-center">
              <Poptip trigger="hover" placement="right-start"  width="260" class="margin-left-20 display-in-small" style="float:left;text-align:left">
                     <span class="hover-cursor hover-color font-size-18"><Icon type="md-list" /></span>
                      <div slot="content">
                        <!-- <img src="../../../assets/mafanr_logo.png" class="hover-cursor margin-top-10" style="width:30px;height:30px;margin-left:30px;" @click="goHome" /> -->
                        <div v-for="i in items" :key="i" class="item item-1"  v-if="level[i]==1" >
                            {{names[i]}}
                        </div>
                        <div class="item hover-cursor item-2" :class="{'selected': selItem==i}" @click="selectItem(i)" v-else>
                          {{names[i]}}
                        </div>
                    </div>
              </Poptip>
              <span class="font-size-24 doc-name">监控文档</span>
              <span class="float-right lang-sel">
                <Icon type="ios-home-outline"  class="font-size-24 margin-right-10 hover-cursor" @click="goHome()" />
                <Icon type="logo-github" style="display:none" class="font-size-24 margin-right-10 hover-cursor" @click="gotoGithub()" />
                <Tag>中文</Tag>
                <!-- <Tag v-if="$store.state.misc.language=='en'" @click.native="setLang('zh')">中文</Tag>
                <Tag v-else @click.native="setLang('en')">EN</Tag> -->
              </span>
            </div>
            <Row class="margin-top-5">
              <Col :xs="{span:24}" :sm="{span:24}" :md="{span: 18, offset: 3}" :lg="{ span: 18, offset: 3}">
                <!-- 渲染页面title -->
                <router-view class="page-title" />
                <!-- 渲染页面内容 -->
                <div class="page-body margin-top-10">
                  <mavon-editor class="mf-markdown" :editable=false ref=md :value="$store.state.misc.content" :subfield=false defaultOpen='preview' :scrollStyle=true   :toolbarsFlag=false> </mavon-editor>
                </div>
              </Col>
            </Row>
            <div class="next-page margin-top-20">
              <span class="prev hover-cursor" v-if="prevItem!=''" @click="selectItem(prevItem)">Prev: <span class="font-weight-500 margin-left-3">{{names[prevItem]}}</span></span>
              <span class="next float-right hover-cursor" v-if="nextItem!=''"  @click="selectItem(nextItem)">Next: <span class="font-weight-500 margin-left-3">{{names[nextItem]}}</span></span>
            </div>      
         </Col>
      </Row>
  </div>
</template>

<script>
export default {
  name: 'docNav',
  data () {
    return {
      doc : '',
      path : '',

      items: [],
      level: {},

      selItem : '',
      prevItem: '',
      nextItem: '',

      initMd : ''
    }
  },
  watch: {
    $route() {
      this.initItem()
    }
  },
  computed: {
 
  },
  methods: {
    gotoGithub() {
      window.location.href= 'https://github.com/mafanr/juz'
    },
    setLang(lang) {
        this.$store.dispatch('setLanguage', lang)
        this.$i18n.locale = lang
        location.reload()
    },
    goHome() {
        window.location.href = window.location.origin + '/apm/ui/list'
    },
    selectItem(i) {
      this.$router.push('/apm/ui/docs/'+ i)
    },
    genPrev() {
        var index = 0
       for (var i=0;i< this.items.length;i++) {
         if (this.selItem == this.items[i]) {
           index = i
         }
       }

       if (i==0) {
         this.prevItem = ''
         return 
       }
       
       for (var i=index-1;i>=0;i--) {
         if (this.level[this.items[i]] == 2) {
            this.prevItem  =  this.items[i]
           return 
         } 
       }
       
        this.prevItem = ''
       return 
    },
    genNext() {
      var index = 0
      for (var i=0;i< this.items.length;i++) {
         if (this.selItem == this.items[i]) {
           index = i
         }
       }

       if (i==0) {
         this.nextItem = ''
         return 
       }
       
       for (var i=index+1;i < this.items.length ;i++) {
         if (this.level[this.items[i]] == 2) {
           this.nextItem  =  this.items[i]
           return 
         } 
       }
       
       this.nextItem = ''
       return 
    },
    initItem() {
      this.path = window.location.pathname

      this.items = ['introduce','about','deploy','install']

      this.level = {introduce: 1,'about':2, deploy:1,install:2}

      this.names = {introduce: this.$t('pageName.introduce'),'about': this.$t('pageName.about'),
      deploy: this.$t('pageName.deploy'),install:this.$t('pageName.install')}

      this.selItem = this.path.split('/')[4]

      this.genPrev()
      this.genNext() 
    }
  },
  mounted() {
    // set nav items
    this.initItem()
  }
}
</script>

<style lang="less">


</style>

<style lang="less" scoped> 
@import "../../theme/doc.less";
.mf-doc-nav {
    .mf-doc-nav-left {
      // height:100vh;
      // background-color: rgb(59, 46, 42);
      position: relative;
      .left-items {
        // position: fixed;
        // border-right:.5px solid #ccc;
        padding-left: 90px;
        // padding-right: 90px;
        margin-top:20px;
        // overflow-y: scroll;
      }
      .site {
        margin-left: 30px;
      }
    }

    .mf-doc-nav-right {
      padding: 0 40px;
      // background:#e1e1db;
      // height:100vh;
      // position: relative;
      // overflow-y: scroll;
      .page-header {  
        padding:12px 0;
        .doc-name {
          color: rgb(115, 116, 128)
        }
        .lang-sel {
          margin-top:5px
        }
      }
  }

    .item-2 {
          border-radius: 4px;
          transition: background-color .3s ease-in-out;
          padding : 7px 20px;
          padding-left: 45px !important
      }

      .item-1 {
        padding: 7px 30px;
        color: @text-light-color
      }

      .item.selected {
        color: rgb(230, 159, 103)
      }
      .item.selected:hover {
        background: transparent
      }
}

.display-in-small {
  display: none;
}
.next-page {
  color: @text-dark-color;
  .hover-cursor:hover {
    color: #333;
  }
}

.top{
        padding: 10px;
        background: rgba(0, 153, 229, .7);
        color: #fff;
        text-align: center;
        border-radius: 2px;
        z-index:999;
    }
@media only screen and (max-width: 992px) {
    .display-in-small {
      display: inherit !important;
    }
    .mf-doc-nav-right {
      padding: 0px 30px !important;
      padding-bottom: 20px !important;
    }
    .page-header {
      position: fixed;
      z-index:999;
      background : white;
      width: 100%;
      border-bottom: .5px solid #ccc;
      .lang-sel {
        margin-right:40px;
      }
    }
    .page-title {
      margin-top: 90px !important;
    }
}
@media only screen and (max-width: 768px) {
    .mf-doc-nav-right {
      padding: 0 5px !important;
      padding-bottom: 20px !important;
    }
    .lang-sel {
      margin-right: 20px !important;
    }
}
</style>
