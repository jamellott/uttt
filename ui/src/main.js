import Vue from "vue";
import Vuex from "vuex";
import VueRouter from "vue-router";
import serverStore from "./serverStore.js";
import routes from "./routes.js";
import * as components from "./components/all.js";

Vue.use(Vuex);
Vue.use(VueRouter);

Vue.config.productionTip = false;
const store = new Vuex.Store(serverStore);
const router = new VueRouter({ routes });

new Vue({
  render: (h) => h(components.Main),
  components,
  router,
  store,
}).$mount("#app");
