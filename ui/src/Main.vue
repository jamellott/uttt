<template>
  <main class="container w-100 h-100 p-0 m-0">
    <router-view></router-view>
  </main>
</template>

<script>
export default {
  name: "Main",
};
</script>

<style>
#main {
  font-family: Avenir, Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  text-align: center;
  margin: 0px;
  padding: 0px;
  width: 100%;
  height: 100%;
}
body {
  margin: 0px;
  padding: 0px;
  width: 100%;
  height: 100%;
}
:root {
  margin: 0px;
  padding: 0px;
  width: 100%;
  height: 100%;
}
#main div {
}
</style>
