import './styles/style.css';

import { Greet } from '../wailsjs/go/main/App';
import { LogInfo } from "../wailsjs/runtime/runtime";

import "./apps/platform/platform";

window.addEventListener("load", (event) => {
  LogInfo(">> LOADED: ");
});

