<template>
  <div class="container" id="app">

    <div class="logo">ChangeMyMind<span>.io</span></div>
    <!-- <router-view/> -->

    <div class="dropshadow biginput">
      <input type="text" placeholder="Type something controversial"/>
    </div>

    <div class="landing_actions">
      <div class="button dropshadow skew red"><span class="noselect"> Create Room </span></div>
      <h3 class="or">OR</h3>
      <div class="button dropshadow skew blue"><span class="noselect"> Join Random Room </span></div>
   </div>

   <div class="landing_info">
     <p>Enter a controversial statment and get matched up with random peoeple for a debate, as the conversation progresses, we’ll throw more randoms into your room to make things more intresting.</p>
   </div>

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
  font-family: 'Avenir', Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  text-align: center;
  color: #2c3e50;
  margin-top: 60px;
}

.landing_actions {
  display: flex;
  justify-content: center;
  margin-top: 20px;
  .or {
    margin: 0 40px;
    margin-top: 34px;
    text-align: center;
    font-size: 1.5rem;
  }
}
.landing_info {
  width: 600px;
  margin: auto;
  margin-top: 70px;
}

</style>
