import Vue from "vue";
import Vuex from "vuex";
import VueRouter from "vue-router";
import { BootstrapVue, IconsPlugin } from "bootstrap-vue";
import "bootstrap";
import "bootstrap/dist/css/bootstrap.min.css";

import serverStore from "./serverStore.js";
import routes from "./routes.js";
import { Main } from "./components/all.js";

Vue.use(Vuex);
Vue.use(VueRouter);
Vue.use(BootstrapVue);
Vue.use(IconsPlugin);

Vue.config.productionTip = false;
const store = new Vuex.Store(serverStore);
const router = new VueRouter({ routes });

new Vue({
  render: (h) => h(Main),
  router,
  store,
}).$mount("#main");
