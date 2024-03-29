// The Vue build version to load with the `import` command
// (runtime-only or standalone) has been set in webpack.base.conf with an alias.
import Vue from 'vue'
import App from './components/layouts/App'
import routes from './router/index'
import VueRouter from 'vue-router'
import 'bootstrap/dist/css/bootstrap.css'
import axios from 'axios'

Vue.config.productionTip = false
Vue.prototype.$axios = axios
Vue.use(VueRouter)
Vue.use(axios)

const router = new VueRouter({
  routes,
  mode: 'history'
})

/* eslint-disable no-new */
new Vue({
  el: '#app',
  router: router,
  components: { App },
  template: '<App/>'
})
