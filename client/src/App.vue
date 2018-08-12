<template>
  <div id="app">
    <div class="stats_top">
      <div class="stat"><i class="fa fa-user"/> {{stats.players}} Debaters</div>
      <div class="stat"><i class="fa fa-comments"/> {{stats.channels}} Rooms</div>
      </div> 
    <router-view/>
   </div>
</template>

<script>
export default {
  name: 'App',
  data() {
    return {
      stats: {}
    }
  },
  mounted() {

    setInterval(() => {

      fetch('https://api.changemymind.io/stats')
      .then((response) => {
        return response.json();
      })
      .then((x) => {
        console.log(x);
        this.stats = x
      });

    },5000)
  }
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
