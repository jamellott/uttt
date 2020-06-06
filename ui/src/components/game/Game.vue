<template>
  <div class="game">
    <div>Game ID: {{ gameID }}</div>
    <h2>
      <span :class="classForID(playerID)">You</span> vs
      <span :class="classForID(opponentID)">{{ nameForID(opponentID) }}</span>
    </h2>
    <div>
      <grid :game="game"></grid>
    </div>
  </div>
</template>

<script>
import Grid from "./Grid.vue";

export default {
  name: "Game",
  props: ["gameID"],
  components: {
    Grid,
  },
  computed: {
    game() {
      return this.$store.state.games[this.gameID];
    },
    playerID() {
      return this.$store.state.playerID;
    },
    opponentID() {
      if (this.game.playerX == this.playerID) {
        return this.game.playerO;
      }

      return this.game.playerX;
    },
  },
  methods: {
    nameForID(id) {
      if (this.game.playerX == id) {
        return this.game.playerXName;
      }
      return this.game.playerOName;
    },
    classForID(id) {
      if (id == this.game.playerX) {
        return "player-x";
      }
      return "player-o";
    },
  },
};
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>
h3 {
  margin: 40px 0 0;
}
ul {
  list-style-type: none;
  padding: 0;
}
li {
  display: inline-block;
  margin: 0 10px;
}
a {
  color: #42b983;
}

h2 {
  padding: 4px;
}
h2 span {
  padding: 4px;
}
.player-o {
  background-color: lightskyblue;
}

.player-x {
  background-color: lightcoral;
}
</style>
