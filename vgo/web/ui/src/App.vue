<template>
  <div id="app">
    <router-view/>
  </div>
</template>

<script>
export default {
  name: 'App'
}
</script>

<style>
#app {
  font-family: 'Avenir', Helvetica, Arial, sans-serif;
}
</style>


<style lang="less" scoped> 
#app::before {
    content: '';
    position: absolute;
    display: block;
    top: 0;
    left: 0;
    right: 0;
    margin: 0;
    padding: 0;
    height: 6px;
    background-color: #F2BC56;
    z-index: 1019;
}
#app{
  // background-color: #f8f8f8;
}
</style>
<style>

</style>