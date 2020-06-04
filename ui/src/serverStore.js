class WSMessage {
  constructor(type, payload) {
    this.messageType = type;
    this.payload = payload;
  }
}

function createLoginRequestVerifier(socket, resolve, reject) {
  let handler = (ev) => {
    socket.removeEventListener("message", handler);
    console.log(ev);
    let data = JSON.parse(ev.data);
    if (data.messageType === "LoginSuccess") {
      resolve({ payload: data.payload, socket });
    } else if (data.messageType === "LoginFailure") {
      console.error(data.payload.message);
      reject("login failed");
    } else {
      reject("Unexpected MessageType" + data.messageType);
    }
  };

  socket.addEventListener("message", handler);
}

function sendLoginRequest(socket, username) {
  let msg = new WSMessage("LoginRequest", {
    loginID: username,
  });
  socket.send(JSON.stringify(msg));
  return new Promise((resolve, reject) => {
    createLoginRequestVerifier(socket, resolve, reject);
  });
}

const webSocketHandler = {
  socket: null,
  store: null,
  openRequests: {},
  requestCounter: 1,
  installFunc() {
    let handler = this;
    return (store) => {
      handler.store = store;
    };
  },
  setSocket(socket) {
    if (this.socket !== null) {
      // todo: support reconnections? Or just refresh
      console.error("socket being replaced!");
    }
    this.socket = socket;
    let handler = this;
    socket.addEventListener("message", (ev) => {
      handler.handleMessage(JSON.parse(ev.data));
    });
  },
  handleMessage(msg) {
    console.debug(msg);
    switch (msg.messageType) {
      case "GameUpdate":
        this.store.commit("gameUpdate", msg.payload);
        break;
      case "UserLookup":
        this.store.commit("addLookupResult", msg.payload);
        break;
      default:
        console.error("unknown websocket message type: " + msg.messageType);
        break;
    }
    let id = msg.requestID.toString();
    if (this.openRequests[id] != undefined) {
      this.openRequests[id]();
      this.openRequests[id] = undefined;
    }
  },
  sendMessage(msg) {
    this.sendMessagePromise(msg);
  },
  sendMessagePromise(msg) {
    msg.requestID = this.requestCounter++;
    console.debug(msg);
    return new Promise((resolve) => {
      this.openRequests[msg.requestID.toString()] = resolve;
      this.socket.send(JSON.stringify(msg));
    });
  },
};

const store = {
  state: {
    username: null,
    playerID: null,
    games: [],
    usernameMap: {},
    playerIDMap: {},
  },
  mutations: {
    setUser(state, { username, playerID }) {
      state.username = username;
      state.playerID = playerID;
    },
    initGames(state, games) {
      state.games = games;
    },
    gameUpdate(state, game) {
      let idx = state.games.find((g) => g.id == game.id);
      if (idx == -1) {
        state.games.push(game);
      } else {
        state.games[idx] = game;
      }
    },
    addLookupResult(state, result) {
      let usernameMap = state.usernameMap;
      let playerIDMap = state.playerIDMap;
      usernameMap[result.username] = result.playerID;
      playerIDMap[result.playerID] = result.username;
      state.usernameMap = usernameMap;
      state.playerIDMap = playerIDMap;
    },
  },
  plugins: [webSocketHandler.installFunc()],
  actions: {
    login(context, username) {
      let url = process.env.VUE_APP_SERVER_SOCKET_URL;
      console.log("connecting to " + url);
      let socket = new WebSocket(url);
      return new Promise((resolve, reject) => {
        socket.addEventListener("open", () => {
          sendLoginRequest(socket, username)
            .then(({ payload, socket }) => {
              console.log("successful connection");
              webSocketHandler.setSocket(socket);
              context.commit("setUser", payload);
              resolve();
            })
            .catch((reason) => {
              console.error(reason);
              reject(reason);
            });
        });
      });
    },
    playMove(context, { gameID, move }) {
      let message = new WSMessage("PlayMove", { gameID, move });
      webSocketHandler.sendMessage(message);
    },
    newGame(context, opponent) {
      let message = new WSMessage("NewGame", { opponent });
      webSocketHandler.sendMessage(message);
    },
    lookupOpponent(context, opponent) {
      let message = new WSMessage("UserLookup", { username: opponent });
      return webSocketHandler.sendMessagePromise(message); // TODO: set up better promise mechanism
    },
  },
  // getters: {
  //   getUsername(state) {
  //     return state.username;
  //   },
  // },
};

export default store;
