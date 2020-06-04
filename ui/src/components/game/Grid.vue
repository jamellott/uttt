<template>
  <div class="container">
    <div v-for="rowIdx in range" :class="`row game-row-${rowIdx + 1}`">
      <div v-for="colIdx in range" :class="`col-4 game-col-${colIdx + 1}`">
        <div
          v-for="subRowIdx in range"
          :class="`row game-row-${subRowIdx + 1}`"
        >
          <div
            v-for="subColIdx in range"
            :class="`col-4 game-col-${subColIdx + 1}`"
          >
            <square
              v-bind:square="rows[rowIdx][colIdx][subRowIdx][subColIdx]"
            />
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  name: "Grid",
  props: {
    gameID: String,
  },
  data() {
    return {
      username: "",
    };
  },
  computed: {
    rows() {
      return this.$store.dispatch("getGame", this.gameID);
    },
  },
  methods: {},
};
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>
/*
.row {
  border-radius: 2px;
}
*/

.game-row-1 {
  border-bottom: 4px solid black;
}

.game-row-2 {
  border-bottom: 4px solid black;
}

.game-row-3 {
}

.game-col-1 {
  border-right: 4px solid black;
}

.game-col-2 {
  border-right: 4px solid black;
}

.game-col-3 {
}
</style>
