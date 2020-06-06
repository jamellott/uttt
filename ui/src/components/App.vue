<template>
  <div class="row w-100 h-100">
    <nav class="col-md-4 col-12 sidebar">
      <ul class="border border-secondary font-weight-bold">
        <li>
          <div class="sidebar-item border border-secondary rounded">
            <profile-display />
          </div>
        </li>
        <li v-for="game in games" :key="game.gameID">
          <router-link :to="'/app/game/' + game.gameID">
            <div class="sidebar-item border border-primary rounded">
              <span class="sidebar-text"> vs {{ getOpponent(game) }}</span>
            </div>
          </router-link>
        </li>

        <li>
          <router-link to="/app/new">
            <div
              class="sidebar-item border border-primary rounded sidebar-text"
            >
              + New Game
            </div>
          </router-link>
        </li>
      </ul>
    </nav>
    <div class="mt-5 col-12 col-md-8">
      <h1>Ultimate Tic Tac Toe</h1>
      <router-view></router-view>
    </div>
  </div>
</template>

<script>
import ProfileDisplay from "./sidebar/ProfileDisplay.vue";

export default {
  name: "App",
  components: {
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
  methods: {
    getOpponent(game) {
      if (game.playerX == this.$store.state.playerID) {
        return game.playerOName;
      }

      return game.playerXName;
    },
  },
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
  margin: 20px;
  width: 80%-20px;
  min-height: 80px;
}
.sidebar-item {
  padding: 20px;
  width: 100%;
  height: 100%;
}
.sidebar-text {
  color: blue;
}
</style>
