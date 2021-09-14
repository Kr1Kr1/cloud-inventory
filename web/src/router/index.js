import Vue from 'vue';
import VueRouter from 'vue-router';
import Servers from '../views/Servers.vue';

Vue.use(VueRouter)

const routes = [
  {
      path: '/servers', 
      component: Servers,
  },
];

const router = new VueRouter({
    mode: 'history',
    routes,
});

export default router;