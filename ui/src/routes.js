import * as components from "./components/all.js";

console.log(components);
const routes = [
  { path: "/", redirect: "login" },
  { path: "/login", component: components.Login },
  {
    path: "/app",
    component: components.App,
    children: [
      { path: "new", component: components.NewGameMenu },
      { path: "game/:gameID", component: components.Game, props: true },
    ],
  },
];

export default routes;
