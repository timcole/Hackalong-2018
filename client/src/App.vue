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

<style lang="scss">
@import 'assets/global.scss';

#app {
  font-family: Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  text-align: center;
  color: #2c3e50;
}


</style>
