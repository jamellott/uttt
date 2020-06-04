<template>
  <div class="row w-100 h-100">
    <div class="col-md-4 col-12 sidebar">
      <ul class="border border-secondary">
        <li class="border border-secondary rounded"><profile-display /></li>
        <li v-for="game in games" :key="game.id">
          <game-list-item />
        </li>
        <router-link to="/app/new"
          ><li class="border border-primary rounded">
            + New Game
          </li></router-link
        >
      </ul>
    </div>
    <div class="col-12 col-md-8">
      <router-view></router-view>
    </div>
  </div>
</template>

<script>
import GameListItem from "./sidebar/GameListItem.vue";
import ProfileDisplay from "./sidebar/ProfileDisplay.vue";

export default {
  name: "App",
  components: {
    GameListItem,
    ProfileDisplay,
  },
  data() {
    return {
      username: "",
    };
  },
  computed: {
    games() {
      return this.$store.state.games;
    },
  },
  methods: {},
  mounted() {
    if (this.$store.state.username === null) {
      this.$router.push("/login");
    }
  },
};
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>
.sidebar ul {
  list-style-type: none;
  padding: 0;
  height: 100%;
}
.sidebar li {
  padding: 20px;
  margin: 20px;
  width: 80%-20px;
  min-height: 80px;
}
</style>
