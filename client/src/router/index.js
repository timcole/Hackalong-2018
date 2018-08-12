import Vue from 'vue'
import Router from 'vue-router'

import Landing from '@/screens/landing'
import Room from '@/screens/room'

import Shortkey from 'vue-shortkey'
Vue.use(Shortkey)

Vue.use(Router)

export default new Router({
  mode: 'history',
  routes: [
    {
      path: '/',
      name: 'landing',
      component: Landing
    },
    {
      path: '/room',
      name: 'room',
      component: Room
    }
  ]
})
