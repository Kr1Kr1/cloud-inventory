import Vue from 'vue'
import axios from 'axios';
import App from './App.vue'
import vuetify from './plugins/vuetify'

Vue.config.productionTip = false
axios.defaults.baseURL = document.baseURI;

new Vue({
  vuetify,
  render: h => h(App)
}).$mount('#app')
